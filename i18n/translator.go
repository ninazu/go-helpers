package i18n

var instance extends

type translation struct {
	cache map[string]map[string]string
}
type extends struct {
	*translation
	lang string
	path string
}

func (i *translation) load() {
	i.cache = make(map[string]map[string]string)
}

func Init(sourcePath string, lang string) (error) {
	instance.path = sourcePath
	instance.translation = &translation{}
	instance.load()
	SetLang(lang)
	
	return nil
}

func Translate(text string, params ...interface{}) (string) {
	if val, ok := instance.cache[instance.lang][text]; ok {
		text = val
	}
	
	return text
}

func SetLang(lang string) {
	instance.lang = lang
}

func GetLang() (string) {
	return instance.lang
}
