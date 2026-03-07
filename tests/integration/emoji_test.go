package integration_test

import (
	"context"

	"github.com/chromedp/chromedp"
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
)

var _ = Describe("emoji", func() {
	var mainTab context.Context
	var checkTab context.Context

	BeforeEach(func() {
		var cancel context.CancelFunc

		mainTab, cancel = chromedp.NewContext(browser)
		DeferCleanup(cancel)
		checkTab, cancel = chromedp.NewContext(browser)
		DeferCleanup(cancel)
	})

	It("suggests emoji completions", func() {
		var suggestions string

		Expect(chromedp.Run(mainTab,
			chromedp.Navigate(tmpbbsURL),
			chromedp.WaitVisible("#body"),
			chromedp.SendKeys("#body", ":sku"),
			chromedp.Poll("document.querySelectorAll('#emoji-suggestions > *').length == 4", nil),
			chromedp.Text("#emoji-suggestions", &suggestions),
		)).To(Succeed())

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
