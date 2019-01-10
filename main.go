package main

import (
	"flag"
	"fmt"
	"log"

	"github.com/0xAX/notificator"
	"github.com/arrufat/papago"
	"github.com/atotto/clipboard"
)

func main() {
	known := flag.String("k", "en", "the language you already know")
	learn := flag.String("l", "ko", "the language you are learning")
	list := flag.Bool("list", false, "list all possible language codes")
	var primary bool
	if hasPrimary {
		flag.BoolVar(&primary, "p", false, "use primary selection")
	}
	flag.Parse()

	if *list {
		fmt.Println("Supported language codes:")
		for _, lang := range papago.SupportedLanguages() {
			fmt.Printf("%s \t=>   %s\n", lang.Code(), lang)
		}
		return
	}

	// notification setup
	notify := notificator.New(notificator.Options{
		DefaultIcon: "icon/default.png",
		AppName:     "Papago",
	})

	knownLang, err := papago.ParseLanguageCode(*known)
	if err != nil {
		notify.Push("Error", err.Error(), "", notificator.UR_NORMAL)
		log.Fatal(err)
	}
	learnLang, err := papago.ParseLanguageCode(*learn)
	if err != nil {
		notify.Push("Error", err.Error(), "", notificator.UR_NORMAL)
		log.Fatal(err)
	}

	// set clipboard to use selection
	if primary {
		setPrimary(true)
	}

	// read from the clipboard
	text, err := clipboard.ReadAll()
	if err != nil {
		notify.Push("Error", err.Error(), "", notificator.UR_NORMAL)
		log.Fatal(err)
	}
	if text == "" {
		notify.Push("Error", "No text selected", "", notificator.UR_NORMAL)
		return
	}
	log.Println("selected text:", text)

	// detect language from selection
	lang, err := papago.Detect(text)
	if err != nil {
		notify.Push("Error", "Unable to detect the language", "", notificator.UR_NORMAL)
		log.Fatal(err)
	}
	log.Println("detected language:", lang)

	// translate
	var trans string
	switch lang {
	case knownLang:
		trans, err = papago.Translate(text,
			knownLang,
			learnLang)
	default:
		trans, err = papago.Translate(text,
			lang,
			knownLang)
	}
	// handle translation errors
	if err != nil {
		notify.Push("Error", err.Error(), "", notificator.UR_NORMAL)
		log.Fatal(err)
	}
	log.Println("translation:", trans)

	// copy the translation to the clipboard
	if primary {
		setPrimary(false)
	}
	if err := clipboard.WriteAll(trans); err != nil {
		log.Fatal(err)
	}

	// notify
	notify.Push(fmt.Sprintf("Translating from %s: %s", lang, text), trans, "", notificator.UR_CRITICAL)

	return
}
