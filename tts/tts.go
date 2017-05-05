package tts

import (
	"fmt"
	"gtts/token"
	"strings"
	"net/http"
	"strconv"
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

type TTS struct {
	Lang string
	Speed float64
	Text string
}

func New(text string, lang string, speed float64) (*TTS, error) {
	if _, ok := LANGS[strings.ToLower(lang)]; ! ok {
		return nil, fmt.Errorf("Language [%s] not supported", lang)
	}
	if speed < 0 || speed > 2 {
		return nil, fmt.Errorf("Speed [%f] not supported", speed)
	}
	t := &TTS{Text: text, Lang: lang, Speed: speed}
	return t, nil
}

func (t *TTS) ToSpeech() error {
	req, err := http.NewRequest("GET", GOOGLE_TTS_URL, nil)
	if err != nil {
		return err
	}
	tkn := token.New()
	query := req.URL.Query()
	query.Add("q", t.Text)
	query.Add("ie", "UTF-8")
	query.Add("tl", t.Lang)
	query.Add("ttsspeed", strconv.FormatFloat(t.Speed, 'g', -1, 64))
	query.Add("total", "1")
	query.Add("idx", "0")
	query.Add("client", "tw-ob")
	query.Add("textlen", string(len(t.Text)))
	query.Add("tk", tkn.Calculate_token(t.Text))
	req.URL.RawQuery = query.Encode()
	req.Header.Add("User-Agent", "Mozilla/5.0 (Windows NT 10.0; WOW64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/47.0.2526.106 Safari/537.36")
	req.Header.Add("Referer", token.GttsUrl)
	fmt.Println(req)
	client := &http.Client{}
	resp, err := client.Do(req)
	fmt.Println(resp)
	return nil
}