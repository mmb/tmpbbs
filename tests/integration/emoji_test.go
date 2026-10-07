package integration_test

import (
	"github.com/chromedp/chromedp"
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
)

var _ = Describe("emoji", func() {
	It("suggests emoji completions", func() {
		Expect(chromedp.Do(mainTab,
			chromedp.Navigate(tmpbbsURL),
			chromedp.WaitVisible("#body"),
			chromedp.SendKeys("#body", ":sku"),
			chromedp.WaitVisible("#emoji-suggestions > *"),
			chromedp.Poll[chromedp.Void]("document.querySelectorAll('#emoji-suggestions > *').length == 4"),
		)).To(Succeed())

		suggestions, err := chromedp.Run(mainTab, chromedp.Text("#emoji-suggestions"))
		Expect(err).NotTo(HaveOccurred())

		Expect(suggestions).To(SatisfyAll(
			ContainSubstring("💀"),
			ContainSubstring("☠️"),
			ContainSubstring("☠"),
			ContainSubstring("🦨"),
		))
	})

	It("substitutes emoji shortcodes in the post title", func() {
		post(mainTab, tmpbbsURL, ":beetle:", "", "")

		Eventually(func() string {
			return get(checkTab, tmpbbsURL)
		}, "5s").Should(ContainSubstring("🪲"))
	})

	It("substitutes emoji shortcodes in the post author", func() {
		post(mainTab, tmpbbsURL, "", ":broken_heart:", "")

		Eventually(func() string {
			return get(checkTab, tmpbbsURL)
		}, "5s").Should(ContainSubstring("💔"))
	})

	It("substitutes emoji shortcodes in the post body", func() {
		post(mainTab, tmpbbsURL, "", "", ":butterfly:")

		Eventually(func() string {
			return get(checkTab, tmpbbsURL)
		}, "5s").Should(ContainSubstring("🦋"))
	})
})
