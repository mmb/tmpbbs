package integration_test

import (
	"context"
	"fmt"
	"time"

	"github.com/chromedp/chromedp"
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
)

var _ = Describe("markdown", func() {
	var (
		mainTab     context.Context
		checkTab    context.Context
		testRootURL string
	)

	BeforeEach(func() {
		var cancel context.CancelFunc

		mainTab, cancel = chromedp.NewContext(browser)
		DeferCleanup(cancel)
		checkTab, cancel = chromedp.NewContext(browser)
		DeferCleanup(cancel)

		testID := fmt.Sprintf("%d", time.Now().UnixNano())
		post(mainTab, tmpbbsURL, testID, "", "")
		Eventually(func() string {
			return get(checkTab, tmpbbsURL)
		}, "5s").Should(ContainSubstring(testID))
		testRootURL = mostRecentReplyURL(checkTab, tmpbbsURL)
	})

	DescribeTable("markdown rendering",
		func(input string, output string) {
			post(mainTab, testRootURL, "", "", input)
			Eventually(func() string {
				return get(checkTab, testRootURL)
			}, "5s").Should(ContainSubstring(output))
		},
		Entry("definition list", "test term\n:   test definition", "<dl>\n<dt>test term</dt>\n<dd>test definition</dd>\n</dl>"),
		Entry("fenced code blocks", "```test\ntest = 1\n```", "<pre><code class=\"language-test\">test = 1\n</code></pre>"),
		Entry("linkify", "http://test.test/", `<a href="http://test.test/" rel="nofollow noreferrer">http://test.test/</a>`),
		Entry("strikethrough", "~strikethrough test~", "<del>strikethrough test</del>"),
		Entry("table", "| test column 1 | test column 2 |\n --- | --- |\n| test 1 | test 2 |", "<table>\n<thead>\n<tr>\n<th>test column 1</th>\n<th>test column 2</th>\n</tr>\n</thead>\n<tbody>\n<tr>\n<td>test 1</td>\n<td>test 2</td>\n</tr>\n</tbody>\n</table>"),
	)
})
