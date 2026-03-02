package integration_test

import (
	"fmt"
	"io"
	"net/http"
	"net/url"
	"os/exec"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
	"github.com/onsi/gomega/gbytes"
	"github.com/onsi/gomega/gexec"
)

const (
	peerPort      = 7801
	peerOverlay   = "kustomize/peer"
	peerNamespace = "tmpbbs-peer"
)

var _ = Describe("peer", Ordered, func() {
	peerURL := fmt.Sprintf("http://localhost:%d", peerPort)

	BeforeAll(func() {
		command := exec.Command("kubectl", "apply", "--kustomize", peerOverlay)
		session, err := gexec.Start(command, GinkgoWriter, GinkgoWriter)
		Expect(err).NotTo(HaveOccurred())
		Eventually(session).Should(gexec.Exit(0))
		DeferCleanup(func() {
			deleteCommand := exec.Command("kubectl", "delete", "--kustomize", peerOverlay)
			deleteSession, deleteErr := gexec.Start(deleteCommand, GinkgoWriter, GinkgoWriter)
			Expect(deleteErr).NotTo(HaveOccurred())
			Eventually(deleteSession, "30s").Should(gexec.Exit(0))
		})

		command = exec.Command("kubectl", "rollout", "status", "statefulset/tmpbbs", "--namespace", peerNamespace)
		session, err = gexec.Start(command, GinkgoWriter, GinkgoWriter)
		Expect(err).NotTo(HaveOccurred())
		Eventually(session, "1m").Should(gexec.Exit(0))

		command = exec.Command("kubectl", "port-forward", "service/tmpbbs-http", "--namespace", peerNamespace,
			fmt.Sprintf("%d:8080", peerPort))
		peerPortForwardSession, err := gexec.Start(command, GinkgoWriter, GinkgoWriter)
		Expect(err).NotTo(HaveOccurred())
		Eventually(peerPortForwardSession, "10s").Should(gbytes.Say("Forwarding from"))
		DeferCleanup(func() {
			peerPortForwardSession.Terminate()
		})
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
			peerResp, peerErr := http.Get(peerURL)
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
