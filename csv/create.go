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

		vals := []string{
			comp.GetUrl(),
			comp.GetCategory().GetTitle(),
			comp.GetSlug(),
			comp.GetTitle(),
			comp.GetEmail(),
			strconv.Itoa(int(comp.GetPhone())),
			comp.GetDescription(),
			makeRuBool(comp.GetOnline()),
			strconv.Itoa(int(comp.GetInn())),
			strconv.Itoa(int(comp.GetKpp())),
			strconv.Itoa(int(comp.GetOgrn())),
			comp.GetDomain().GetAddress(),
			comp.GetDomain().GetRegistrar(),
			comp.GetDomain().GetRegistrationDate(),
			comp.GetAvatar(),
			comp.GetLocation().GetCity().GetTitle(),
			comp.GetLocation().GetAddress(),
			comp.GetLocation().GetAddressTitle(),
			comp.GetApp().GetAppStore().GetUrl(),
			comp.GetApp().GetGooglePlay().GetUrl(),
			strconv.Itoa(int(comp.GetSocial().GetVk().GetGroupId())),
			comp.GetSocial().GetVk().GetName(),
			comp.GetSocial().GetVk().GetScreenName(),
			parser.IsClosed_name[int32(comp.GetSocial().GetVk().GetIsClosed())],
			comp.GetSocial().GetVk().GetDescription(),
			strconv.Itoa(int(comp.GetSocial().GetVk().GetMembersCount())),
			comp.GetSocial().GetVk().GetPhoto_200(),
			comp.GetSocial().GetInstagram().GetUrl(),
			comp.GetSocial().GetTwitter().GetUrl(),
			comp.GetSocial().GetYoutube().GetUrl(),
			comp.GetSocial().GetFacebook().GetUrl(),
			comp.GetUpdatedAt(),
			strings.Join(techCats, ";"),
			strconv.Itoa(int(comp.GetPageSpeed())),
			makeRuBool(comp.GetVerified()),
			makeRuBool(comp.GetPremium()),
		}
		vals = append(vals, managers...)

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
