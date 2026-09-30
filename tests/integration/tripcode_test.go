package integration_test

import (
	"fmt"
	"time"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
)

var _ = Describe("tripcode", func() {
	var testRootURL string

	BeforeEach(func() {
		testID := fmt.Sprintf("%d", time.Now().UnixNano())
		post(mainTab, tmpbbsURL, testID, "", "")
		Eventually(func() string {
			return get(checkTab, tmpbbsURL)
		}, "5s").Should(ContainSubstring(testID))
		testRootURL = mostRecentReplyURL(checkTab, tmpbbsURL)
	})

	It("calculates a tripcode", func() {
		post(mainTab, testRootURL, "", "user#tripcode", "")
		Eventually(func() string {
			return get(checkTab, testRootURL)
		}, "5s").Should(ContainSubstring("user"))

		Expect(get(checkTab, testRootURL)).To(ContainSubstring("!b57220e925"))
	})
})
