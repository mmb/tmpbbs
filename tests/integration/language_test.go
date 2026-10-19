package integration_test

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/chromedp/cdproto/network"
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
	"golang.org/x/text/language"
)

var _ = Describe("language", func() {
	localesPath := filepath.Join("..", "..", "internal", "tmpbbs", "locales")
	languages, err := os.ReadDir(localesPath)
	Expect(err).NotTo(HaveOccurred())

	entries := make([]TableEntry, len(languages))
	for i, lang := range languages {
		entries[i] = Entry("sets the language and direction for "+lang.Name(), lang.Name())
	}

	DescribeTable("language and direction",
		func(lang string) {
			dir := "ltr"

			script, _ := language.MustParse(lang).Script()
			if script.String() == "Arab" {
				dir = "rtl"
			}

			Expect(getWithHeaders(checkTab, tmpbbsURL, network.Headers{"Accept-Language": lang})).To(
				ContainSubstring(fmt.Sprintf("html lang=%q dir=%q", lang, dir)))
		}, entries,
	)
})
