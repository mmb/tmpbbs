package integration_test

import (
	"context"

	"github.com/chromedp/chromedp"
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
)

var _ = Describe("peer", Ordered, func() {
	var (
		peer1URL string
		peer2URL string
		peer1Tab context.Context
		peer2Tab context.Context
	)

	BeforeAll(func() {
		peer1URL = deployOverlay("peer-1", 7900)
		peer2URL = deployOverlay("peer-2", 7901)
	})

	BeforeEach(func() {
		var cancel context.CancelFunc

		peer1Tab, cancel = chromedp.NewContext(browser)
		DeferCleanup(cancel)
		peer2Tab, cancel = chromedp.NewContext(browser)
		DeferCleanup(cancel)
	})

	It("pulls a post from a peer", func() {
		post(peer1Tab, peer1URL, "test title", "test author#tripcode", "test body")

		Eventually(func() string {
			return get(peer2Tab, peer2URL)
		}, "5s").Should(SatisfyAll(
			ContainSubstring("test title"),
			ContainSubstring("test author"),
			ContainSubstring("!a24ebe09a9"),
			ContainSubstring("test body"),
		))
	})
})
