package integration_test

import (
	"fmt"
	"net/url"
	"time"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
)

var _ = Describe("pagination", func() {
	var testRootURL string

	BeforeEach(func() {
		testID := fmt.Sprintf("%d", time.Now().UnixNano())
		post(mainTab, tmpbbsURL, testID, "", "")
		Eventually(func() string {
			return get(checkTab, tmpbbsURL)
		}, "5s").Should(ContainSubstring(testID))
		testRootURL = mostRecentReplyURL(checkTab, tmpbbsURL)
	})

	It("shows page navigation", func() {
		for i := range 100 {
			title := fmt.Sprintf("pagination test post %d", i)
			post(mainTab, testRootURL, title, "", "")
			Eventually(func() string {
				return get(checkTab, testRootURL)
			}, "5s").Should(ContainSubstring(title))
		}

		testRootURLParsed, err := url.Parse(testRootURL)
		Expect(err).NotTo(HaveOccurred())
		path := testRootURLParsed.Path

		Expect(get(checkTab, testRootURL)).To(ContainSubstring(fmt.Sprintf(
			`>page 1 / <a href="%s?p=2#replies-start">page 2</a> / <a href="%s?p=10#replies-start">page 10</a></li>`,
			path, path)))
		Expect(get(checkTab, testRootURL+"?p=2")).To(ContainSubstring(fmt.Sprintf(
			`><a href="%s?p=1#replies-start">page 1</a> / page 2 / <a href="%s?p=3#replies-start">page 3</a> / <a href="%s?p=10#replies-start">page 10</a></li>`,
			path, path, path)))
		Expect(get(checkTab, testRootURL+"?p=3")).To(ContainSubstring(fmt.Sprintf(
			`><a href="%s?p=1#replies-start">page 1</a> / <a href="%s?p=2#replies-start">page 2</a> / page 3 / <a href="%s?p=4#replies-start">page 4</a> / <a href="%s?p=10#replies-start">page 10</a></li>`,
			path, path, path, path)))
		Expect(get(checkTab, testRootURL+"?p=9")).To(ContainSubstring(fmt.Sprintf(
			`><a href="%s?p=1#replies-start">page 1</a> / <a href="%s?p=8#replies-start">page 8</a> / page 9 / <a href="%s?p=10#replies-start">page 10</a></li>`,
			path, path, path)))
		Expect(get(checkTab, testRootURL+"?p=10")).To(ContainSubstring(fmt.Sprintf(
			`><a href="%s?p=1#replies-start">page 1</a> / <a href="%s?p=9#replies-start">page 9</a> / page 10</li>`,
			path, path)))
	})
})
