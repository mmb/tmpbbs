package integration_test

import (
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
)

var _ = Describe("peering with TLS", Ordered, func() {
	var (
		peerServerURL string
		peerClientURL string
	)

	BeforeAll(func(ctx SpecContext) {
		peerServerURL = deployOverlay(ctx, "peer-tls-server", 7903)
		peerClientURL = deployOverlay(ctx, "peer-tls-client", 7904)
	})

	It("pulls a post from a peer", func() {
		post(mainTab, peerServerURL, "test title", "test author#tripcode", "test body")

		Eventually(func() string {
			return get(checkTab, peerClientURL)
		}, "5s").Should(SatisfyAll(
			ContainSubstring("test title"),
			ContainSubstring("test author"),
			ContainSubstring("!a24ebe09a9"),
			ContainSubstring("test body"),
		))
	})
})
