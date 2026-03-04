package integration_test

import (
	"fmt"
	"io"
	"net/http"
	"net/url"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
)

var _ = Describe("peer", Ordered, func() {
	var tmpbbsURL string

	BeforeAll(func() {
		port := 7801
		deployOverlay("peer", port)
		tmpbbsURL = fmt.Sprintf("http://localhost:%d", port)
	})

	It("pulls a post from main", func() {
		resp, err := http.PostForm(mainURL, url.Values{
			"title":  []string{"test title"},
			"author": []string{"test author#tripcode"},
			"body":   []string{"test body"},
		})
		Expect(err).NotTo(HaveOccurred())
		Expect(resp.StatusCode).To(Equal(http.StatusOK))

		Eventually(func() string {
			peerResp, peerErr := http.Get(tmpbbsURL)
			Expect(peerErr).NotTo(HaveOccurred())
			Expect(peerResp.StatusCode).To(Equal(http.StatusOK))
			body, bodyErr := io.ReadAll(peerResp.Body)
			Expect(bodyErr).NotTo(HaveOccurred())

			return string(body)
		}, "1m15s").Should(SatisfyAll(
			ContainSubstring("test title"),
			ContainSubstring("test author"),
			ContainSubstring("!a24ebe09a9"),
			ContainSubstring("test body"),
		))
	})
})
