package integration_test

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/chromedp/cdproto/network"
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
)

var _ = Describe("language", func() {
	localesPath := filepath.Join("..", "..", "internal", "tmpbbs", "locales")
	languages, err := os.ReadDir(localesPath)
	Expect(err).NotTo(HaveOccurred())

	entries := make([]TableEntry, len(languages))
	for i, language := range languages {
		entries[i] = Entry("sets the language and direction for "+language.Name(), language.Name())
	}

	DescribeTable("language and direction",
		func(language string) {
			dir := "ltr"
			if language == "ar" {
				dir = "rtl"
			}

			Expect(getWithHeaders(checkTab, tmpbbsURL, network.Headers{"Accept-Language": language})).To(
				ContainSubstring(fmt.Sprintf("html lang=%q dir=%q", language, dir)))
		}, entries,
	)
})
