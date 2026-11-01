package integration_test

import (
	"fmt"
	"strconv"
	"time"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
)

var _ = Describe("page title", func() {
	DescribeTable("page title",
		func(title string, htmlRegexp string) {
			testAuthor := strconv.FormatInt(time.Now().UnixNano(), 10)
			post(mainTab, tmpbbsURL, title, testAuthor, "")
			Eventually(func() string {
				return get(checkTab, tmpbbsURL)
			}, "5s").Should(ContainSubstring(testAuthor))

			postURL := mostRecentReplyURL(checkTab, tmpbbsURL)

			Eventually(func() string {
				return get(checkTab, postURL)
			}, "5s").Should(MatchRegexp(fmt.Sprintf("<title>%s</title>", htmlRegexp)))
		},
		Entry("empty", "", "[0-9a-z]{4}-[0-9a-z]{2}"),
		Entry("whitespace only", "   ", "[0-9a-z]{4}-[0-9a-z]{2}"),
		Entry("plain", "test", "test"),
		Entry("HTML entity escaping", "1 & 2", "1 &amp; 2"),
		Entry("HTML sanitization", "<b>test</b>", "test"),
		Entry("HTML all stripped", "<script></script>", "[0-9a-z]{4}-[0-9a-z]{2}"),
		Entry("emoji expansion", ":yellow_circle::yellow_square:", "🟡🟨"),
	)
})
