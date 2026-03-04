package integration_test

import (
	"fmt"
	"os/exec"
	"path/filepath"
	"testing"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
	"github.com/onsi/gomega/gbytes"
	"github.com/onsi/gomega/gexec"
)

const mainPort = 7800

var mainURL = fmt.Sprintf("http://localhost:%d", mainPort)

var _ = BeforeSuite(func() {
	deployOverlay("main", mainPort)
})

func TestIntegration(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "Integration Suite")
}

func deployOverlay(name string, port int) {
	path := filepath.Join("kustomize", name)
	namespace := "tmpbbs-" + name

	command := exec.Command("kubectl", "apply", "--kustomize", path)
	session, err := gexec.Start(command, GinkgoWriter, GinkgoWriter)
	Expect(err).NotTo(HaveOccurred())
	Eventually(session).Should(gexec.Exit(0))
	DeferCleanup(func() {
		deleteCommand := exec.Command("kubectl", "delete", "--kustomize", path)
		deleteSession, deleteErr := gexec.Start(deleteCommand, GinkgoWriter, GinkgoWriter)
		Expect(deleteErr).NotTo(HaveOccurred())
		Eventually(deleteSession, "30s").Should(gexec.Exit(0))
	})

	command = exec.Command("kubectl", "rollout", "status", "statefulset/tmpbbs", "--namespace", namespace)
	session, err = gexec.Start(command, GinkgoWriter, GinkgoWriter)
	Expect(err).NotTo(HaveOccurred())
	Eventually(session, "2m30s").Should(gexec.Exit(0))

	command = exec.Command("kubectl", "port-forward", "service/tmpbbs-http", "--namespace", namespace,
		fmt.Sprintf("%d:8080", port))
	portForwardSession, err := gexec.Start(command, GinkgoWriter, GinkgoWriter)
	Expect(err).NotTo(HaveOccurred())
	Eventually(portForwardSession, "10s").Should(gbytes.Say("Forwarding from"))
	DeferCleanup(func() {
		portForwardSession.Terminate()
	})
}
