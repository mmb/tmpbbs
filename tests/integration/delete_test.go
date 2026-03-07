package integration_test

import (
	"context"
	"fmt"
	"time"

	"github.com/chromedp/chromedp"
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
)

var _ = Describe("deleting posts", func() {
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

	It("allows users to delete their own posts", func() {
		post(mainTab, testRootURL, "delete test title", "mikami1#secret", "delete test body")
		Eventually(func() string {
			return get(checkTab, testRootURL)
		}, "5s").Should(ContainSubstring("delete test body"))
		postURL := mostRecentReplyURL(checkTab, testRootURL)

		post(mainTab, postURL, "", "mikami1#secret", "!delete")
		Eventually(func() string {
			return get(checkTab, postURL)
		}, "5s").Should(SatisfyAll(
			ContainSubstring("deleted"),
			Not(ContainSubstring("delete test title")),
		))
	})

	It("does not allow deletion unless the delete reply tripcode matches the original post", func() {
		post(mainTab, testRootURL, "delete test title", "mikami#secret", "delete test body")
		Eventually(func() string {
			return get(checkTab, testRootURL)
		}, "5s").Should(ContainSubstring("delete test body"))
		postURL := mostRecentReplyURL(checkTab, testRootURL)

		post(mainTab, postURL, "", "mikami1#other", "!delete")
		Eventually(func() string {
			return get(checkTab, postURL)
		}, "5s").Should(ContainSubstring("!delete"))
	})

	It("allows superusers to delete any post", func() {
		post(mainTab, testRootURL, "delete test title", "mikami#secret", "delete test body")
		Eventually(func() string {
			return get(checkTab, testRootURL)
		}, "5s").Should(ContainSubstring("delete test body"))
		postURL := mostRecentReplyURL(checkTab, testRootURL)

		post(mainTab, postURL, "", "superuser#secret", "!delete")
		Eventually(func() string {
			return get(checkTab, postURL)
		}, "5s").Should(SatisfyAll(
			ContainSubstring("deleted"),
			Not(ContainSubstring("delete test title")),
			Not(ContainSubstring("mikami")),
		))
	})
})
