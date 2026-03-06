package integration_test

import (
	"context"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"strconv"
	"testing"

	"github.com/chromedp/chromedp"
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
	"github.com/onsi/gomega/gbytes"
	"github.com/onsi/gomega/gexec"
)

const (
	namespacePrefix = "tmpbbs-test-"
	basePort        = 7800
)

var (
	tmpbbsURL string
	browser   context.Context
)

var _ = SynchronizedBeforeSuite(
	func() {},
	func() {
		name := strconv.Itoa(GinkgoParallelProcess())
		overlayPath := filepath.Join("kustomize", name)
		Expect(os.Mkdir(overlayPath, 0o755)).To(Succeed())
		DeferCleanup(os.RemoveAll, overlayPath)

		kustomizationYaml := fmt.Appendf(nil, "namespace: %s%s\nresources: [../base]", namespacePrefix, name)
		Expect(os.WriteFile(filepath.Join(overlayPath, "kustomization.yaml"), kustomizationYaml, 0o644)).To(Succeed())

		tmpbbsURL = deployOverlay(name, basePort+GinkgoParallelProcess())

		execAllocator, cancel := chromedp.NewExecAllocator(context.Background(),
			append(chromedp.DefaultExecAllocatorOptions[:],
				chromedp.Flag("disable-dev-shm-usage", true),
				chromedp.Flag("no-sandbox", true),
			)...)
		DeferCleanup(cancel)
		browser, cancel = chromedp.NewContext(execAllocator)
		DeferCleanup(cancel)
	})

func TestIntegration(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "Integration Suite")
}

func deployOverlay(name string, port int) string {
	path := filepath.Join("kustomize", name)
	command := exec.Command("kubectl", "apply", "--kustomize", path)
	session, err := gexec.Start(command, GinkgoWriter, GinkgoWriter)
	Expect(err).NotTo(HaveOccurred())
	Eventually(session, "5s").Should(gexec.Exit(0))
	DeferCleanup(func() {
		deleteCommand := exec.Command("kubectl", "delete", "--kustomize", path)
		deleteSession, deleteErr := gexec.Start(deleteCommand, GinkgoWriter, GinkgoWriter)
		Expect(deleteErr).NotTo(HaveOccurred())
		Eventually(deleteSession, "30s").Should(gexec.Exit(0))
	})

	namespace := namespacePrefix + name
	command = exec.Command("kubectl", "rollout", "status", "statefulset/tmpbbs", "--namespace", namespace)
	session, err = gexec.Start(command, GinkgoWriter, GinkgoWriter)
	Expect(err).NotTo(HaveOccurred())
	Eventually(session, "10s").Should(gexec.Exit(0))

	command = exec.Command("kubectl", "port-forward", "service/tmpbbs-http", "--namespace", namespace,
		fmt.Sprintf("%d:8080", port))
	portForwardSession, err := gexec.Start(command, GinkgoWriter, GinkgoWriter)
	Expect(err).NotTo(HaveOccurred())
	Eventually(portForwardSession, "10s").Should(gbytes.Say("Forwarding from"))
	DeferCleanup(portForwardSession.Terminate)

	return fmt.Sprintf("http://localhost:%d/", port)
}

func post(ctx context.Context, url string, title string, author string, body string) {
	Expect(chromedp.Run(ctx,
		chromedp.Navigate(url),
		chromedp.WaitVisible(`input[type="submit"]`),
		chromedp.SendKeys("#title", title),
		chromedp.SendKeys("#author", author),
		chromedp.SendKeys("#body", body),
		chromedp.Click(`input[type="submit"]`),
	)).To(Succeed())
}

func get(ctx context.Context, url string) string {
	var html string
	Expect(chromedp.Run(ctx, chromedp.Navigate(url), chromedp.OuterHTML("html", &html))).To(Succeed())

	return html
}
