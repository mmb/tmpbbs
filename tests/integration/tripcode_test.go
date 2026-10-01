package integration_test

import (
	"strconv"
	"time"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
)

var _ = Describe("tripcode", func() {
	var testRootURL string

	BeforeEach(func() {
		testID := strconv.FormatInt(time.Now().UnixNano(), 10)
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
