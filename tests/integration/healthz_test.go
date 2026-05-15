package integration_test

import (
	"github.com/chromedp/chromedp"
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
)

var _ = Describe("healthz", func() {
	It("returns ok", func() {
		var body string

		Expect(chromedp.Run(mainTab,
			chromedp.Navigate(tmpbbsURL+"healthz"),
			chromedp.Text("body", &body),
		)).To(Succeed())

		Expect(body).To(Equal("ok"))
	})
})
