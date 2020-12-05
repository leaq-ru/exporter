package xlsx

import (
	"fmt"
	"github.com/google/uuid"
	"github.com/nnqq/scr-proto/codegen/go/parser"
	"github.com/tealeg/xlsx/v3"
	"os"
	"strings"
)

func Create(ch chan *parser.FullCompanyV2) (s3XlsxURL string, err error) {
	file := xlsx.NewFile(xlsx.UseDiskVCellStore)
	sh, err := file.AddSheet("Компании")
	if err != nil {
		return
	}

	const managersCount = 5

	header := sh.AddRow()
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
	for _, val := range headerVals {
		header.AddCell().SetValue(val)
	}

	for comp := range ch {
		row := sh.AddRow()

		var techCats []string
		for _, tc := range comp.GetTechnologyCategories() {
			for _, t := range tc.GetTechnologies() {
				techCats = append(techCats, strings.Join([]string{
					tc.GetName(),
					t.GetName(),
				}, " - "))
			}
		}

		var managers []interface{}
		for _, ppl := range comp.GetPeople() {
			managers = append(
				managers,
				ppl.GetVkId(),
				ppl.GetFirstName(),
				ppl.GetLastName(),
				makeRuBool(ppl.GetVkIsClosed()),
				makeRuSex(ppl.GetSex()),
				ppl.GetPhoto_200(),
				ppl.GetPhone(),
				ppl.GetEmail(),
				ppl.GetDescription(),
			)
		}

		vals := []interface{}{
			comp.GetUrl(),
			comp.GetCategory().GetTitle(),
			comp.GetSlug(),
			comp.GetTitle(),
			comp.GetEmail(),
			comp.GetPhone(),
			comp.GetDescription(),
			makeRuBool(comp.GetOnline()),
			comp.GetInn(),
			comp.GetKpp(),
			comp.GetOgrn(),
			comp.GetDomain().GetAddress(),
			comp.GetDomain().GetRegistrar(),
			comp.GetDomain().GetRegistrationDate(),
			comp.GetAvatar(),
			comp.GetLocation().GetCity().GetTitle(),
			comp.GetLocation().GetAddress(),
			comp.GetLocation().GetAddressTitle(),
			comp.GetApp().GetAppStore().GetUrl(),
			comp.GetApp().GetGooglePlay().GetUrl(),
			comp.GetSocial().GetVk().GetGroupId(),
			comp.GetSocial().GetVk().GetName(),
			comp.GetSocial().GetVk().GetScreenName(),
			parser.IsClosed_name[int32(comp.GetSocial().GetVk().GetIsClosed())],
			comp.GetSocial().GetVk().GetDescription(),
			comp.GetSocial().GetVk().GetMembersCount(),
			comp.GetSocial().GetVk().GetPhoto_200(),
			comp.GetSocial().GetInstagram().GetUrl(),
			comp.GetSocial().GetTwitter().GetUrl(),
			comp.GetSocial().GetYoutube().GetUrl(),
			comp.GetSocial().GetFacebook().GetUrl(),
			comp.GetUpdatedAt(),
			strings.Join(techCats, ","),
			comp.GetPageSpeed(),
			makeRuBool(comp.GetVerified()),
			makeRuBool(comp.GetPremium()),
		}
		vals = append(vals, managers...)

		for _, val := range vals {
			row.AddCell().SetValue(val)
		}
	}

	u, err := uuid.NewRandom()
	if err != nil {
		return
	}

	s3XlsxURL = "tmp/" + u.String() + ".xlsx"
	err = file.Save(s3XlsxURL)
	if err != nil {
		return
	}

	return uploadToS3(s3XlsxURL)
}

func uploadToS3(path string) (s3URL string, err error) {
	os.Open(path)
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
