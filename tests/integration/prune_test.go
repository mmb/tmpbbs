package integration_test

import (
	"context"
	"fmt"

	"github.com/chromedp/chromedp"
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
)

var _ = Describe("prune", Ordered, func() {
	var (
		tmpbbsURL string
		mainTab   context.Context
		checkTab  context.Context
	)

	BeforeAll(func() {
		port := 7802
		deployOverlay("prune", port)
		tmpbbsURL = fmt.Sprintf("http://localhost:%d", port)
	})

	BeforeEach(func() {
		var cancel context.CancelFunc

		mainTab, cancel = chromedp.NewContext(browser)
		DeferCleanup(cancel)
		checkTab, cancel = chromedp.NewContext(browser)
		DeferCleanup(cancel)
	})

	It("prunes posts", func() {
		post(mainTab, tmpbbsURL, "", "", "prune test")

		Eventually(func() string {
			return get(checkTab, tmpbbsURL)
		}, "5s").Should(ContainSubstring("prune test"))

		Eventually(func() string {
			return get(checkTab, tmpbbsURL)
		}, "10s").ShouldNot(ContainSubstring("prune test"))
	})
})
