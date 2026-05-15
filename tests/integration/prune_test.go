package integration_test

import (
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
)

var _ = Describe("prune", Ordered, func() {
	var pruneURL string

	BeforeAll(func() {
		pruneURL = deployOverlay("prune", 7902)
	})

	It("prunes posts", func() {
		post(mainTab, pruneURL, "", "", "prune test")

		Eventually(func() string {
			return get(checkTab, pruneURL)
		}, "5s").Should(ContainSubstring("prune test"))

		Eventually(func() string {
			return get(checkTab, pruneURL)
		}, "15s").ShouldNot(ContainSubstring("prune test"))
	})
})
