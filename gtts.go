package main

import (
	"gtts/tts"
	"log"
)

func main() {
	text := "test"
	t, err := tts.New(text, "ru", 1)
	if err != nil {
		log.Panicln(err)
	}
	t.ToSpeech()
}
