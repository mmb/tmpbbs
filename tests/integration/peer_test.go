package integration_test

import (
	"context"
	"fmt"

	"github.com/chromedp/chromedp"
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
)

var _ = Describe("peer", Ordered, func() {
	var (
		tmpbbsURL string
		mainTab   context.Context
		peerTab   context.Context
	)

	BeforeAll(func() {
		port := 7801
		deployOverlay("peer", port)
		tmpbbsURL = fmt.Sprintf("http://localhost:%d", port)
	})

	BeforeEach(func() {
		var cancel context.CancelFunc

		mainTab, cancel = chromedp.NewContext(browser)
		DeferCleanup(cancel)
		peerTab, cancel = chromedp.NewContext(browser)
		DeferCleanup(cancel)
	})

	It("pulls a post from main", func() {
		post(mainTab, mainURL, "test title", "test author#tripcode", "test body")

		Eventually(func() string {
			return get(peerTab, tmpbbsURL)
		}, "1m15s").Should(SatisfyAll(
			ContainSubstring("test title"),
			ContainSubstring("test author"),
			ContainSubstring("!a24ebe09a9"),
			ContainSubstring("test body"),
		))
	})
})
