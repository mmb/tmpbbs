package integration_test

import (
	"fmt"
	"io"
	"net/http"
	"net/url"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
)

var _ = Describe("prune", Ordered, func() {
	var tmpbbsURL string

	BeforeAll(func() {
		port := 7802
		deployOverlay("prune", port)
		tmpbbsURL = fmt.Sprintf("http://localhost:%d", port)
	})

	It("prunes posts", func() {
		resp, err := http.PostForm(tmpbbsURL, url.Values{"body": []string{"prune test"}})
		Expect(err).NotTo(HaveOccurred())
		Expect(resp.StatusCode).To(Equal(http.StatusOK))

		resp, err = http.Get(tmpbbsURL)
		Expect(err).NotTo(HaveOccurred())
		Expect(resp.StatusCode).To(Equal(http.StatusOK))
		body, bodyErr := io.ReadAll(resp.Body)
		Expect(bodyErr).NotTo(HaveOccurred())
		Expect(body).To(ContainSubstring("prune test"))

		Eventually(func() []byte {
			afterResp, afterErr := http.Get(tmpbbsURL)
			Expect(afterErr).NotTo(HaveOccurred())
			Expect(afterResp.StatusCode).To(Equal(http.StatusOK))
			afterBody, afterBodyErr := io.ReadAll(afterResp.Body)
			Expect(afterBodyErr).NotTo(HaveOccurred())

			return afterBody
		}, "10s").ShouldNot(ContainSubstring("prune test"))
	})
})
