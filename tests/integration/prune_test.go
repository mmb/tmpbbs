package integration_test

import (
	"context"

	"github.com/chromedp/chromedp"
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
)

var _ = Describe("prune", Ordered, func() {
	var (
		pruneURL string
		mainTab  context.Context
		checkTab context.Context
	)

	BeforeAll(func() {
		pruneURL = deployOverlay("prune", 7902)
	})

	BeforeEach(func() {
		var cancel context.CancelFunc

		mainTab, cancel = chromedp.NewContext(browser)
		DeferCleanup(cancel)
		checkTab, cancel = chromedp.NewContext(browser)
		DeferCleanup(cancel)
	})

	It("prunes posts", func() {
		post(mainTab, pruneURL, "", "", "prune test")

		Eventually(func() string {
			return get(checkTab, pruneURL)
		}, "5s").Should(ContainSubstring("prune test"))

		Eventually(func() string {
			return get(checkTab, pruneURL)
		}, "10s").ShouldNot(ContainSubstring("prune test"))
	})
})
