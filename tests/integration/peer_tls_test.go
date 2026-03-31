package integration_test

import (
	"context"

	"github.com/chromedp/chromedp"
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
)

var _ = Describe("peering with TLS", Ordered, func() {
	var (
		peerServerURL string
		peerClientURL string
		peerServerTab context.Context
		peerClientTab context.Context
	)

	BeforeAll(func() {
		peerServerURL = deployOverlay("peer-tls-server", 7903)
		peerClientURL = deployOverlay("peer-tls-client", 7904)
	})

	BeforeEach(func() {
		var cancel context.CancelFunc

		peerServerTab, cancel = chromedp.NewContext(browser)
		DeferCleanup(cancel)
		peerClientTab, cancel = chromedp.NewContext(browser)
		DeferCleanup(cancel)
	})

	It("pulls a post from a peer", func() {
		post(peerServerTab, peerServerURL, "test title", "test author#tripcode", "test body")

		Eventually(func() string {
			return get(peerClientTab, peerClientURL)
		}, "5s").Should(SatisfyAll(
			ContainSubstring("test title"),
			ContainSubstring("test author"),
			ContainSubstring("!a24ebe09a9"),
			ContainSubstring("test body"),
		))
	})
})
