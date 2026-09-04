package integration_test

import (
	"fmt"
	"net/http"
	"net/url"
	"strings"
	"time"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
)

var _ = Describe("maximum input length", func() {
	var testRootURL string

	BeforeEach(func() {
		testID := fmt.Sprintf("%d", time.Now().UnixNano())
		post(mainTab, tmpbbsURL, testID, "", "")
		Eventually(func() string {
			return get(checkTab, tmpbbsURL)
		}, "5s").Should(ContainSubstring(testID))
		testRootURL = mostRecentReplyURL(checkTab, tmpbbsURL)
	})

	Describe("title", func() {
		DescribeTable("validation",
			func(input string, expectedStatus int) {
				resp, err := http.PostForm(testRootURL, url.Values{"title": {input}})
				Expect(err).NotTo(HaveOccurred())
				defer resp.Body.Close()
				Expect(resp.StatusCode).To(Equal(expectedStatus))
			},
			Entry("30 ASCII characters", strings.Repeat("A", 30), http.StatusOK),
			Entry("31 ASCII characters", strings.Repeat("A", 31), http.StatusBadRequest),

			Entry("30 CJK characters", strings.Repeat("字", 30), http.StatusOK),
			Entry("31 CJK characters", strings.Repeat("字", 31), http.StatusBadRequest),

			Entry("15 4-byte emoji", strings.Repeat("😀", 15), http.StatusOK),
			Entry("16 4-byte emoji", strings.Repeat("😀", 16), http.StatusBadRequest),
		)

		DescribeTable("browser maxlength truncation is valid on backend",
			func(input string, output string) {
				post(mainTab, testRootURL, input, "", "")
				Eventually(func() string {
					return get(checkTab, testRootURL)
				}, "5s").Should(ContainSubstring(output))
			},
			Entry("ASCII characters", strings.Repeat("A", 31), strings.Repeat("A", 30)),
			Entry("CJK characters", strings.Repeat("字", 31), strings.Repeat("字", 30)),
			Entry("4-byte emoji", strings.Repeat("😀", 16), strings.Repeat("😀", 15)))
	})

	Describe("author", func() {
		DescribeTable("validation",
			func(input string, expectedStatus int) {
				resp, err := http.PostForm(testRootURL, url.Values{"author": {input}})
				Expect(err).NotTo(HaveOccurred())
				defer resp.Body.Close()
				Expect(resp.StatusCode).To(Equal(expectedStatus))
			},
			Entry("28 ASCII characters", strings.Repeat("A", 28), http.StatusOK),
			Entry("29 ASCII characters", strings.Repeat("A", 29), http.StatusBadRequest),

			Entry("28 CJK characters", strings.Repeat("字", 28), http.StatusOK),
			Entry("29 CJK characters", strings.Repeat("字", 29), http.StatusBadRequest),

			Entry("14 4-byte emoji", strings.Repeat("😀", 14), http.StatusOK),
			Entry("15 4-byte emoji", strings.Repeat("😀", 15), http.StatusBadRequest),
		)

		DescribeTable("browser maxlength truncation is valid on backend",
			func(input string, output string) {
				post(mainTab, testRootURL, "", input, "")
				Eventually(func() string {
					return get(checkTab, testRootURL)
				}, "5s").Should(ContainSubstring(output))
			},
			Entry("ASCII characters", strings.Repeat("A", 29), strings.Repeat("A", 28)),
			Entry("CJK characters", strings.Repeat("字", 29), strings.Repeat("字", 28)),
			Entry("4-byte emoji", strings.Repeat("😀", 15), strings.Repeat("😀", 14)))
	})

	Describe("body", func() {
		DescribeTable("validation",
			func(input string, expectedStatus int) {
				resp, err := http.PostForm(testRootURL, url.Values{"body": {input}})
				Expect(err).NotTo(HaveOccurred())
				defer resp.Body.Close()
				Expect(resp.StatusCode).To(Equal(expectedStatus))
			},
			Entry("8192 ASCII characters", strings.Repeat("A", 8192), http.StatusOK),
			Entry("8193 ASCII characters", strings.Repeat("A", 8193), http.StatusBadRequest),

			Entry("8192 CJK characters", strings.Repeat("字", 8192), http.StatusOK),
			Entry("8193 CJK characters", strings.Repeat("字", 8193), http.StatusBadRequest),

			Entry("4096 4-byte emoji", strings.Repeat("😀", 4096), http.StatusOK),
			Entry("4097 4-byte emoji", strings.Repeat("😀", 4097), http.StatusBadRequest),
		)

		DescribeTable("browser maxlength truncation is valid on backend", Serial,
			func(input string, output string) {
				post(mainTab, testRootURL, "", "", input)
				Eventually(func() string {
					return get(checkTab, testRootURL)
				}, "5s").Should(ContainSubstring(output))
			},
			Entry("ASCII characters", strings.Repeat("A", 8193), strings.Repeat("A", 8192)),
			Entry("CJK characters", strings.Repeat("字", 8193), strings.Repeat("字", 8192)),
			Entry("4-byte emoji", strings.Repeat("😀", 4097), strings.Repeat("😀", 4096)))
	})
})
