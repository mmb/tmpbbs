package tmpbbs

import (
	"log"

	"golang.org/x/text/feature/plural"
	"golang.org/x/text/language"
	"golang.org/x/text/message"
)

func init() {
	err := message.Set(language.English, "%d replies", plural.Selectf(1, "%d",
		"=1", "%d reply",
		"other", "%d replies",
	))

	if err != nil {
		log.Fatal(err)
	}
}
