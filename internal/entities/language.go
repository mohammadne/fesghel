package entities

type Language string

const (
	LanguageEnglish Language = "en"
	LanguagePersian Language = "fa"

	LanguageDefault Language = LanguageEnglish
)

func ToLanguage(rawLanguage string) Language {
	switch rawLanguage {
	case string(LanguageEnglish):
		return LanguageEnglish
	case string(LanguagePersian):
		return LanguagePersian
	default:
		return LanguageDefault
	}
}
