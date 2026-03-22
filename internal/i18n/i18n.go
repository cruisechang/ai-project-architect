package i18n

import (
	"os"
	"strings"
)

var currentLang = "en"
var messages = map[string]map[string]string{}

func register(lang string, msgs map[string]string) {
	messages[lang] = msgs
}

// Init detects language from APA_LANG env or --lang flag in os.Args, then sets it.
// Must be called before building cobra commands so help text uses the right language.
func Init() {
	SetLang(detect())
}

func detect() string {
	if lang := os.Getenv("APA_LANG"); lang != "" {
		return normalize(lang)
	}
	args := os.Args[1:]
	for i, arg := range args {
		if arg == "--lang" && i+1 < len(args) {
			return normalize(args[i+1])
		}
		if v, ok := strings.CutPrefix(arg, "--lang="); ok {
			return normalize(v)
		}
	}
	return "en"
}

func normalize(lang string) string {
	switch strings.ToLower(strings.ReplaceAll(lang, "_", "-")) {
	case "zh-tw":
		return "zh-TW"
	case "zh-cn":
		return "zh-CN"
	case "ja":
		return "ja"
	case "ko":
		return "ko"
	case "de":
		return "de"
	case "es":
		return "es"
	case "fr":
		return "fr"
	default:
		return "en"
	}
}

// SetLang sets the active language.
func SetLang(lang string) { currentLang = lang }

// Lang returns the active language code.
func Lang() string { return currentLang }

// T returns the translation for key in the active language, falling back to English.
func T(key string) string {
	if msgs, ok := messages[currentLang]; ok {
		if s, ok := msgs[key]; ok && s != "" {
			return s
		}
	}
	if currentLang != "en" {
		if msgs, ok := messages["en"]; ok {
			if s, ok := msgs[key]; ok {
				return s
			}
		}
	}
	return key
}
