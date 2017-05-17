package main

import (
	"log"
	"gtts/cache"
	"fmt"
	"crypto/md5"
	"encoding/hex"
	"gtts/tts"
	"os"
	"strings"
	"io/ioutil"
	"path"
)

func NameConvert(text string) string {
	sum := md5.Sum([]byte(text))
	data := sum[:]
	fname := hex.EncodeToString(data)
	return fmt.Sprintf("%s.mp3", fname)
}

func main() {
	if len(os.Args) < 3 {
		os.Exit(1)
	}

	number := os.Args[1]
	text := os.Args[2]

	if ! strings.Contains(text, "PROBLEM") {
		os.Exit(1)
	}

	fcache, err := cache.New("/tmp/asterisk/sounds", NameConvert)
	if err != nil {
		log.Fatalln(err)
	}

	t, err := tts.New(text,1, fcache)
	if err != nil {
		log.Fatalln(err)
	}

	filename, err := t.ToSpeech()
	if err != nil {
		log.Fatalln(err)
	}

	tmpcall, err := ioutil.TempFile(os.TempDir(), "zabbix_call")
	if err != nil {
		log.Fatalln(err)
	}
	defer os.Remove(tmpcall.Name())

	cfilename := strings.TrimSuffix(path.Base(filename), path.Ext(filename))
	content := "Channel:Local/" + number + "@office-asterisk\nContext:office-asterisk\nWaitTime:60\nCallerID:zabbix\nApplication:playback\nData:" + cfilename + "&" + cfilename + "&" + cfilename + "\nSet: ALARM_TEXT=\"" + text + "\""

	if _, err := tmpcall.Write([]byte(content)); err != nil {
		log.Fatalln(err)
	}

	if err := tmpcall.Chmod(0666); err != nil {
		log.Fatalln(err)
	}

	if err := tmpcall.Close(); err != nil {
		log.Fatalln(err)
	}
	err = os.Rename(tmpcall.Name(), path.Join("/tmp/asterisk/outgoing", path.Base(tmpcall.Name()) + ".call"))
	if err != nil {
		log.Fatalln(err)
	}
}
