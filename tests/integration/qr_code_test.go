package integration_test

import (
	"github.com/chromedp/chromedp"
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
)

var _ = Describe("QR code dialog", func() {
	BeforeEach(func() {
		Expect(chromedp.Do(mainTab,
			chromedp.Navigate(tmpbbsURL),
			chromedp.WaitVisible("#show-qr"),
			chromedp.Click(`#show-qr`),
		)).To(Succeed())
	})

	It("shows the QR code image", func() {
		Expect(chromedp.Do(mainTab,
			chromedp.Poll[chromedp.Void](`document.querySelector('img[src^="/qr"]')?.naturalWidth > 0`)),
		).To(Succeed())
	})

	It("hides the QR code dialog", func() {
		Expect(chromedp.Do(mainTab,
			chromedp.WaitVisible("dialog#qr button"),
			chromedp.Click("dialog#qr button"),
			chromedp.Poll[chromedp.Void](`!document.querySelector("#qr").open`),
		)).To(Succeed())
	})
})
