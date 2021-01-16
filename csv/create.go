package csv

import (
	"encoding/csv"
	"fmt"
	"github.com/google/uuid"
	"github.com/nnqq/scr-proto/codegen/go/parser"
	"os"
	"strconv"
	"strings"
)

func Create(ch chan *parser.FullCompanyV2) (csvPath string, err error) {
	u, err := uuid.NewRandom()
	if err != nil {
		return
	}

	csvPath = u.String() + ".csv"
	fd, err := os.Create(csvPath)
	if err != nil {
		return
	}
	defer fd.Close()

	file := csv.NewWriter(fd)
	defer file.Flush()

	const managersCount = 5

	headerVals := []string{
		"Сайт",
		"Категория",
		"Slug https://leaq.ru/company/{slug}",
		"Название",
		"Email",
		"Телефон",
		"Описание",
		"Сайт онлайн",
		"ИНН",
		"КПП",
		"ОГРН",
		"IP сервера",
		"Регистратор домена",
		"Дата регистрации домена",
		"Логотип",
		"Город",
		"Адрес улица",
		"Адрес дом/офис",
		"Приложение в AppStore",
		"Приложение в GooglePlay",
		"VK группа - ID",
		"VK группа - Название",
		"VK группа - Короткий адрес",
		"VK группа - Приватность",
		"VK группа - Описание",
		"VK группа - Кол-во подписчиков",
		"VK группа - Аватар",
		"Instagram страница",
		"Twitter профиль",
		"YouTube канал",
		"Facebook группа",
		"Обновлено",
		"Технологии на сайте",
		"Скорость загрузки сайта в ms",
		"Владелец подтвержден",
		"Приоритетное размещение",
	}
	headerVals = append(headerVals, makeManagers(managersCount)...)
	headerVals = append(headerVals, "DNS")
	err = file.Write(headerVals)
	if err != nil {
		return
	}

	for comp := range ch {
		var techCats []string
		for _, tc := range comp.GetTechnologyCategories() {
			for _, t := range tc.GetTechnologies() {
				techCats = append(techCats, strings.Join([]string{
					tc.GetName(),
					t.GetName(),
				}, " - "))
			}
		}

		var managers []string
		for i, ppl := range comp.GetPeople() {
			if i+1 == managersCount {
				break
			}

			managers = append(
				managers,
				strconv.Itoa(int(ppl.GetVkId())),
				ppl.GetFirstName(),
				ppl.GetLastName(),
				makeRuBool(ppl.GetVkIsClosed()),
				makeRuSex(ppl.GetSex()),
				ppl.GetPhoto_200(),
				strconv.Itoa(int(ppl.GetPhone())),
				ppl.GetEmail(),
				ppl.GetDescription(),
			)
		}

		var dnsNames []string
		for _, dns := range comp.GetDns() {
			dnsNames = append(dnsNames, dns.GetName())
		}

		const sepMultiValue = ";"

		vals := make([]string, len(headerVals))

		vals[0] = comp.GetUrl()
		vals[1] = comp.GetCategory().GetTitle()
		vals[2] = comp.GetSlug()
		vals[3] = comp.GetTitle()
		vals[4] = comp.GetEmail()
		vals[5] = strconv.Itoa(int(comp.GetPhone()))
		vals[6] = comp.GetDescription()
		vals[7] = makeRuBool(comp.GetOnline())
		vals[8] = strconv.Itoa(int(comp.GetInn()))
		vals[9] = strconv.Itoa(int(comp.GetKpp()))
		vals[10] = strconv.Itoa(int(comp.GetOgrn()))
		vals[11] = comp.GetDomain().GetAddress()
		vals[12] = comp.GetDomain().GetRegistrar()
		vals[13] = comp.GetDomain().GetRegistrationDate()
		vals[14] = comp.GetAvatar()
		vals[15] = comp.GetLocation().GetCity().GetTitle()
		vals[16] = comp.GetLocation().GetAddress()
		vals[17] = comp.GetLocation().GetAddressTitle()
		vals[18] = comp.GetApp().GetAppStore().GetUrl()
		vals[19] = comp.GetApp().GetGooglePlay().GetUrl()
		vals[20] = strconv.Itoa(int(comp.GetSocial().GetVk().GetGroupId()))
		vals[21] = comp.GetSocial().GetVk().GetName()
		vals[22] = comp.GetSocial().GetVk().GetScreenName()
		vals[23] = parser.IsClosed_name[int32(comp.GetSocial().GetVk().GetIsClosed())]
		vals[24] = comp.GetSocial().GetVk().GetDescription()
		vals[25] = strconv.Itoa(int(comp.GetSocial().GetVk().GetMembersCount()))
		vals[26] = comp.GetSocial().GetVk().GetPhoto_200()
		vals[27] = comp.GetSocial().GetInstagram().GetUrl()
		vals[28] = comp.GetSocial().GetTwitter().GetUrl()
		vals[29] = comp.GetSocial().GetYoutube().GetUrl()
		vals[30] = comp.GetSocial().GetFacebook().GetUrl()
		vals[31] = comp.GetUpdatedAt()
		vals[32] = strings.Join(techCats, sepMultiValue)
		vals[33] = strconv.Itoa(int(comp.GetPageSpeed()))
		vals[34] = makeRuBool(comp.GetVerified())
		vals[35] = makeRuBool(comp.GetPremium())
		for _, man := range managers {
			vals[35+1] = man
		}
		vals[len(makeManagers(1))*managersCount+36] = strings.Join(dnsNames, sepMultiValue)

		for i, val := range vals {
			vals[i] = strings.ReplaceAll(strings.ToValidUTF8(val, " "), "\n", " ")
		}

		if len(headerVals) > len(vals) {
			vals = append(vals, make([]string, len(headerVals)-len(vals))...)
		}
		err = file.Write(vals)
		if err != nil {
			return
		}
	}
	return
}

func makeManagers(count int) (out []string) {
	for i := 1; i <= count; i += 1 {
		prefix := fmt.Sprintf("Менеджер #%v - ", i)

		out = append(out, []string{
			prefix + "VK ID",
			prefix + "Имя",
			prefix + "Фамилия",
			prefix + "VK закрыт?",
			prefix + "Пол",
			prefix + "Аватар",
			prefix + "Телефон",
			prefix + "Email",
			prefix + "Описание",
		}...)
	}
	return
}

func makeRuBool(b bool) (out string) {
	if b {
		out = "да"
	} else {
		out = "нет"
	}
	return
}

func makeRuSex(sex parser.Sex) (out string) {
	const none = "не указан"
	switch sex {
	case parser.Sex_NONE:
		return none
	case parser.Sex_MALE:
		return "муж"
	case parser.Sex_FEMALE:
		return "жен"
	default:
		return none
	}
}
