package integration_test

import (
	"context"

	"github.com/chromedp/chromedp"
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
)

var _ = Describe("emoji suggestions", func() {
	var mainTab context.Context

	BeforeEach(func() {
		var cancel context.CancelFunc

		mainTab, cancel = chromedp.NewContext(browser)
		DeferCleanup(cancel)
	})

	It("suggests emoji completions", func() {
		var suggestions string

		Expect(chromedp.Run(mainTab,
			chromedp.Navigate(mainURL),
			chromedp.WaitVisible("#body"),
			chromedp.SendKeys("#body", ":sku"),
			chromedp.Poll(`document.querySelectorAll('#emoji-suggestions > *').length == 4`, nil),
			chromedp.Text("#emoji-suggestions", &suggestions),
		)).To(Succeed())

		Expect(suggestions).To(SatisfyAll(
			ContainSubstring("💀"),
			ContainSubstring("☠️"),
			ContainSubstring("☠"),
			ContainSubstring("🦨"),
		))
	})
})
