package integration_test

import (
	"context"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"strconv"
	"strings"
	"testing"
	"time"

	"github.com/chromedp/cdproto/runtime"
	"github.com/chromedp/chromedp"
	"github.com/chromedp/chromedp/remote"
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
	"github.com/onsi/gomega/gbytes"
	"github.com/onsi/gomega/gexec"
)

const (
	chromeDebuggingPort = "9222"
	basePort            = 7800
	namespacePrefix     = "tmpbbs-test-"
	chromeTimeout       = 1 * time.Minute
)

var (
	tmpbbsURL          string
	browser            context.Context
	chromeWebSocketURL = "http://localhost:" + chromeDebuggingPort
	mainTab            context.Context
	checkTab           context.Context
)

var _ = SynchronizedBeforeSuite(
	func(ctx SpecContext) {
		var browserCancel context.CancelFunc

		if os.Getenv("TMPBBS_BUILD_IMAGE") == "true" {
			command := exec.CommandContext(ctx, "docker", "build", "../..", "--tag", "kind-registry:5000/tmpbbs:test")
			session, err := gexec.Start(command, GinkgoWriter, GinkgoWriter)
			Expect(err).NotTo(HaveOccurred())
			Eventually(session, "1m").Should(gexec.Exit(0))
		}

		execAllocator, execAllocatorCancel := chromedp.NewExecAllocator(context.Background(),
			append(chromedp.DefaultExecAllocatorOptions[:],
				chromedp.Flag("disable-dev-shm-usage", true),
				chromedp.Flag("ignore-certificate-errors", true),
				chromedp.Flag("remote-debugging-port", chromeDebuggingPort),
				chromedp.NoSandbox,
				chromedp.WSURLReadTimeout(1*time.Minute),
				remote.WebSocket,
			)...)
		DeferCleanup(execAllocatorCancel)

		browser, browserCancel = chromedp.NewContext(execAllocator)
		DeferCleanup(browserCancel)
		Expect(chromedp.Do(browser)).To(Succeed())
	},
	func(ctx SpecContext) {
		var browserCancel context.CancelFunc

		name := strconv.Itoa(GinkgoParallelProcess())
		overlayPath := filepath.Join("kustomize", name)
		Expect(os.RemoveAll(overlayPath)).To(Succeed())
		Expect(os.Mkdir(overlayPath, 0o755)).To(Succeed())
		DeferCleanup(os.RemoveAll, overlayPath)

		kustomizationYaml := fmt.Appendf(nil, "namespace: %s%s\nresources: [../base]", namespacePrefix, name)
		Expect(os.WriteFile(filepath.Join(overlayPath, "kustomization.yaml"), kustomizationYaml, 0o644)).To(Succeed())

		tmpbbsURL = deployOverlay(ctx, name, basePort+GinkgoParallelProcess())

		remoteAllocator, remoteAllocatorCancel := remote.NewAllocator(context.Background(), chromeWebSocketURL)
		DeferCleanup(remoteAllocatorCancel)

		browser, browserCancel = chromedp.NewContext(remoteAllocator)
		DeferCleanup(browserCancel)
	})

var _ = SynchronizedAfterSuite(func() {}, func(ctx SpecContext) {
	command := exec.CommandContext(ctx, "kubectl", "delete", "namespace", "--selector", "tmpbbs-test=true")
	session, err := gexec.Start(command, GinkgoWriter, GinkgoWriter)
	Expect(err).NotTo(HaveOccurred())
	Eventually(session, "1m").Should(gexec.Exit(0))
})

var _ = BeforeEach(func() {
	mainTab = newTab()
	checkTab = newTab()
})

func TestIntegration(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "Integration Suite")
}

func deployOverlay(ctx context.Context, name string, port int) string {
	command := exec.CommandContext(ctx, "kubectl", "apply", "--kustomize", filepath.Join("kustomize", name))
	session, err := gexec.Start(command, GinkgoWriter, GinkgoWriter)
	Expect(err).NotTo(HaveOccurred())
	Eventually(session, "5s").Should(gexec.Exit(0))

	namespace := namespacePrefix + name

	command = exec.CommandContext(ctx, "kubectl", "rollout", "status", "statefulset/tmpbbs", "--namespace", namespace)
	session, err = gexec.Start(command, GinkgoWriter, GinkgoWriter)
	Expect(err).NotTo(HaveOccurred())
	Eventually(session, "30s").Should(gexec.Exit(0))

	command = exec.Command("kubectl", "port-forward", "service/tmpbbs-http", //nolint:noctx // needs to keep running
		"--namespace", namespace, fmt.Sprintf("%d:8080", port))
	portForwardSession, err := gexec.Start(command, GinkgoWriter, GinkgoWriter)
	Expect(err).NotTo(HaveOccurred())
	Eventually(portForwardSession, "10s").Should(gbytes.Say("Forwarding from"))
	DeferCleanup(portForwardSession.Terminate)

	command = exec.Command("kubectl", "logs", "--follow", "--namespace", namespace, //nolint:noctx // needs to keep running
		"--prefix", "--selector", "app=tmpbbs")
	logSession, err := gexec.Start(command, GinkgoWriter, GinkgoWriter)
	Expect(err).NotTo(HaveOccurred())
	DeferCleanup(logSession.Terminate)

	scheme := "http"
	if strings.Contains(name, "tls-server") {
		scheme = "https"
	}

	return fmt.Sprintf("%s://localhost:%d/", scheme, port)
}

func newTab() context.Context {
	tab, tabCancel := chromedp.NewContext(browser)
	DeferCleanup(tabCancel)

	tab, tabCancelTimeout := context.WithTimeout(tab, chromeTimeout)
	DeferCleanup(tabCancelTimeout)

	exceptions := chromedp.Events(tab, runtime.ExceptionThrown)
	go func() {
		for exception, err := range exceptions {
			if err != nil {
				return
			}

			GinkgoWriter.Printf("javascript exception: %s\n",
				&chromedp.ExceptionError{ExceptionDetails: exception.ExceptionDetails})
		}
	}()

	return tab
}

func post(ctx context.Context, url string, title string, author string, body string) {
	Expect(chromedp.Do(ctx,
		chromedp.Navigate(url),
		chromedp.WaitVisible(`input[type="submit"]`),
		chromedp.SendKeys("#title", title),
		chromedp.SendKeys("#author", author),
		chromedp.SendKeys("#body", body),
		chromedp.Click(`input[type="submit"]`),
	)).To(Succeed())
}

func get(ctx context.Context, url string) string {
	Expect(chromedp.Do(ctx, chromedp.Navigate(url))).To(Succeed())

	html, err := chromedp.Run(ctx, chromedp.OuterHTML("html"))
	Expect(err).NotTo(HaveOccurred())

	return html
}

func mostRecentReplyURL(ctx context.Context, parentURL string) string {
	Expect(chromedp.Do(ctx,
		chromedp.Navigate(parentURL),
		chromedp.WaitVisible("#replies-start + li a"),
	)).To(Succeed())

	replyURL, err := chromedp.Run(ctx,
		chromedp.Evaluate[string]("document.querySelector('#replies-start + li a').href"),
	)
	Expect(err).NotTo(HaveOccurred())

	return replyURL
}
