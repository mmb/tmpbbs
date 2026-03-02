package integration_test

import (
	"fmt"
	"os/exec"
	"testing"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
	"github.com/onsi/gomega/gbytes"
	"github.com/onsi/gomega/gexec"
)

const (
	mainPort      = 7800
	mainOverlay   = "kustomize/main"
	mainNamespace = "tmpbbs"
)

var (
	mainURL                = fmt.Sprintf("http://localhost:%d", mainPort)
	mainPortForwardSession *gexec.Session
)

var _ = BeforeSuite(func() {
	var (
		session *gexec.Session
		err     error
	)

	command := exec.Command("kubectl", "apply", "--kustomize", mainOverlay)
	session, err = gexec.Start(command, GinkgoWriter, GinkgoWriter)
	Expect(err).NotTo(HaveOccurred())
	Eventually(session).Should(gexec.Exit(0))
	DeferCleanup(func() {
		deleteCommand := exec.Command("kubectl", "delete", "--kustomize", mainOverlay)
		deleteSession, deleteErr := gexec.Start(deleteCommand, GinkgoWriter, GinkgoWriter)
		Expect(deleteErr).NotTo(HaveOccurred())
		Eventually(deleteSession, "30s").Should(gexec.Exit(0))
	})

	command = exec.Command("kubectl", "rollout", "status", "statefulset/tmpbbs", "--namespace", mainNamespace)
	session, err = gexec.Start(command, GinkgoWriter, GinkgoWriter)
	Expect(err).NotTo(HaveOccurred())
	Eventually(session, "2m30s").Should(gexec.Exit(0))

	command = exec.Command("kubectl", "port-forward", "service/tmpbbs-http", "--namespace", mainNamespace,
		fmt.Sprintf("%d:8080", mainPort))
	mainPortForwardSession, err = gexec.Start(command, GinkgoWriter, GinkgoWriter)
	Expect(err).NotTo(HaveOccurred())
	Eventually(mainPortForwardSession, "10s").Should(gbytes.Say("Forwarding from"))
	DeferCleanup(func() {
		mainPortForwardSession.Terminate()
	})
})

func TestIntegration(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "Integration Suite")
}
