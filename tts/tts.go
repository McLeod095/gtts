package tts

import (
	"fmt"
	"gtts/token"
	"net/http"
	"strconv"
	"io/ioutil"
	"unicode"
	"gtts/cache"
)

var GOOGLE_TTS_URL = "https://translate.google.com/translate_tts"
var MAX_CHARS = 100
var LANGS = map[string]string{
	"af" : "Afrikaans",
	"sq" : "Albanian",
	"ar" : "Arabic",
	"hy" : "Armenian",
	"bn" : "Bengali",
	"ca" : "Catalan",
	"zh" : "Chinese",
	"zh-cn" : "Chinese (Mandarin/China)",
	"zh-tw" : "Chinese (Mandarin/Taiwan)",
	"zh-yue" : "Chinese (Cantonese)",
	"hr" : "Croatian",
	"cs" : "Czech",
	"da" : "Danish",
	"nl" : "Dutch",
	"en" : "English",
	"en-au" : "English (Australia)",
	"en-uk" : "English (United Kingdom)",
	"en-us" : "English (United States)",
	"eo" : "Esperanto",
	"fi" : "Finnish",
	"fr" : "French",
	"de" : "German",
	"el" : "Greek",
	"hi" : "Hindi",
	"hu" : "Hungarian",
	"is" : "Icelandic",
	"id" : "Indonesian",
	"it" : "Italian",
	"ja" : "Japanese",
	"km" : "Khmer (Cambodian)",
	"ko" : "Korean",
	"la" : "Latin",
	"lv" : "Latvian",
	"mk" : "Macedonian",
	"no" : "Norwegian",
	"pl" : "Polish",
	"pt" : "Portuguese",
	"ro" : "Romanian",
	"ru" : "Russian",
	"sr" : "Serbian",
	"si" : "Sinhala",
	"sk" : "Slovak",
	"es" : "Spanish",
	"es-es" : "Spanish (Spain)",
	"es-us" : "Spanish (United States)",
	"sw" : "Swahili",
	"sv" : "Swedish",
	"ta" : "Tamil",
	"th" : "Thai",
	"tr" : "Turkish",
	"uk" : "Ukrainian",
	"vi" : "Vietnamese",
	"cy" : "Welsh",
}

type Phrase struct {
	Text string
	Lang string
}

type TTS struct {
	Speed float64
	Text	string
	Phrases []Phrase
	Token *token.GttsToken
	Cache cache.Cache
}

func splitText(text string) []Phrase {
	var p []Phrase
	var lang string
	var prevlang string
	previndex := 0
	for index, rune := range text {
		if unicode.IsLetter(rune) || unicode.IsDigit(rune) {
			lang = checkLang(rune)
		}
		if prevlang == "" {
			prevlang = lang
		}
		if prevlang != lang {
			p=append(p, Phrase{Text: text[previndex:index], Lang:prevlang})
			prevlang = lang
			previndex = index
		}
	}
	p=append(p, Phrase{Text: text[previndex:], Lang:prevlang})
	return p
}

func checkLang(letter rune) string {
	if unicode.IsOneOf([]*unicode.RangeTable{unicode.Cyrillic, unicode.Digit}, letter) {
		return "ru"
	}
	return "en"
}

//func checkLang(letter rune) string {
//	ru := map[rune]struct{}{
//		'а': struct{}{},
//		'б': struct{}{},
//		'в': struct{}{},
//		'г': struct{}{},
//		'д': struct{}{},
//		'е': struct{}{},
//		'ё': struct{}{},
//		'ж': struct{}{},
//		'з': struct{}{},
//		'и': struct{}{},
//		'й': struct{}{},
//		'к': struct{}{},
//		'л': struct{}{},
//		'м': struct{}{},
//		'н': struct{}{},
//		'о': struct{}{},
//		'п': struct{}{},
//		'р': struct{}{},
//		'с': struct{}{},
//		'т': struct{}{},
//		'у': struct{}{},
//		'ф': struct{}{},
//		'х': struct{}{},
//		'ц': struct{}{},
//		'ч': struct{}{},
//		'ш': struct{}{},
//		'щ': struct{}{},
//		'ъ': struct{}{},
//		'ы': struct{}{},
//		'ь': struct{}{},
//		'э': struct{}{},
//		'ю': struct{}{},
//		'я': struct{}{},
//		'0': struct{}{},
//		'1': struct{}{},
//		'2': struct{}{},
//		'3': struct{}{},
//		'4': struct{}{},
//		'5': struct{}{},
//		'6': struct{}{},
//		'7': struct{}{},
//		'8': struct{}{},
//		'9': struct{}{},
//	}
//	if _, ok := ru[unicode.ToLower(letter)]; ok {
//		return "ru"
//	}
//	return "en"
//}

func New(text string, speed float64, c cache.Cache) (*TTS, error) {
	if speed < 0 || speed > 2 {
		return nil, fmt.Errorf("Speed [%f] not supported", speed)
	}
	t := &TTS{Text: text, Speed: speed, Phrases: splitText(text), Cache: c}
	tkn, err := token.New()
	if err != nil {
		return nil, err
	}
	t.Token = tkn
	return t, nil
}

func (t *TTS) speech(p Phrase, length int, index int) ([]byte, error) {
	req, err := http.NewRequest("GET", GOOGLE_TTS_URL, nil)
	if err != nil {
		return nil, err
	}
	query := req.URL.Query()
	query.Add("q", p.Text)
	query.Add("ie", "UTF-8")
	query.Add("tl", p.Lang)
	query.Add("ttsspeed", strconv.FormatFloat(t.Speed, 'g', -1, 64))
	query.Add("total", string(length))
	query.Add("idx", string(index))
	query.Add("client", "tw-ob")
	query.Add("textlen", string(len(p.Text)))
	query.Add("tk", t.Token.Calculate_token(p.Text))

	req.URL.RawQuery = query.Encode()
	req.Header.Add("User-Agent", "Mozilla/5.0 (Windows NT 10.0; WOW64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/47.0.2526.106 Safari/537.36")
	req.Header.Add("Referer", token.GttsUrl)
	client := &http.Client{}

	resp, err := client.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}
	return body, nil
}

func (t *TTS) ToSpeech() (string, error) {
	if t.Cache.Exists(t.Text) {
		return t.Cache.GetName(t.Text), nil
	}
	var mp3 []byte
	for index, phrase := range t.Phrases {
		var body []byte
		var err error
		if t.Cache.Exists(phrase.Text) {
			body, err = t.Cache.Get(phrase.Text)
			if err != nil {
				return "", err
			}
		}else{
			body, err = t.speech(phrase, len(t.Phrases), index)
			if err != nil {
				return "", err
			}
			t.Cache.Set(phrase.Text,body)
		}
		mp3 = append(mp3, body...)
	}
	return t.Cache.Set(t.Text,mp3)
}