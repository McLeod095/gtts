package main

import (
//	"log"
//	"gtts/tts"
//	"fmt"
	"gtts/cache"
	"fmt"
	"crypto/md5"
	"encoding/hex"
)

func NameConvert(text string) string {
	sum := md5.Sum([]byte(text))
	data := sum[:]
	fname := hex.EncodeToString(data)
	return fmt.Sprintf("%s.mp3", fname)
}

func main() {
	text := "Привет, у нас проблема! PROBLEM at app01.tema: OPR.VK.CPA.171400.SubmitSMresp.LagTooLong (1s 276ms)"
	//t, err := tts.New(text,1)
	//if err != nil {
	//	log.Panicln(err)
	//}
	//filename, err := t.ToSpeech()
	//if err != nil {
	//	log.Panicln(err)
	//}
	fcache, err := cache.New("/var/tmp/asterisk", NameConvert)
	if err != nil {
		fmt.Println(err)
	}

	fmt.Println(fcache.Get(text))
	fmt.Println(text)
}
