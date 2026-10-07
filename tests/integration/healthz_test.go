package integration_test

import (
	"github.com/chromedp/chromedp"
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
)

var _ = Describe("healthz", func() {
	It("returns ok", func() {
		Expect(chromedp.Do(mainTab, chromedp.Navigate(tmpbbsURL+"healthz"))).To(Succeed())

		body, err := chromedp.Run(mainTab, chromedp.Text("body"))
		Expect(err).NotTo(HaveOccurred())

		Expect(body).To(Equal("ok"))
	})
})
