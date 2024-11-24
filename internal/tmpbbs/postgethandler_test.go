package tmpbbs_test

import (
	"bytes"
	"net/http"

	"github.com/mmb/tmpbbs/internal/tmpbbs"
	"github.com/mmb/tmpbbs/tmpbbsfakes"
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
	"github.com/onsi/gomega/gbytes"
)

var _ = Describe("postgethandler", func() {
	Describe("ServeHTTP", func() {
		var fakeResponseWriter *tmpbbsfakes.FakeResponseWriter
		var buffer *gbytes.Buffer
		var request *http.Request

		BeforeEach(func() {
			var err error
			request, err = http.NewRequest(http.MethodGet, "", bytes.NewBuffer([]byte{}))
			Expect(err).ToNot(HaveOccurred())

			fakeResponseWriter = &tmpbbsfakes.FakeResponseWriter{}
			fakeResponseWriter.HeaderReturns(http.Header{})
			buffer = gbytes.NewBuffer()
			fakeResponseWriter.WriteCalls(buffer.Write)
		})

		Context("when there is no id in the path", func() {
			It("shows the root post", func() {
				postStore := tmpbbs.NewPostStore("test-title")
				postGetHandler := tmpbbs.NewPostGetHandler(10, []string{}, true, true, postStore)
				postGetHandler.ServeHTTP(fakeResponseWriter, request)

				Eventually(buffer).Should(gbytes.Say(`<a href="/0">test-title</a>`))
			})
		})

		Context("when the id is not an integer", func() {
			BeforeEach(func() {
				request.SetPathValue("id", "invalid")
			})

			It("returns 404 not found", func() {
				postGetHandler := tmpbbs.NewPostGetHandler(10, []string{}, true, true, nil)
				postGetHandler.ServeHTTP(fakeResponseWriter, request)

				Expect(fakeResponseWriter.WriteHeaderCallCount()).To(Equal(1))
				Expect(fakeResponseWriter.WriteHeaderArgsForCall(0)).To(Equal(http.StatusNotFound))
				Eventually(buffer).Should(gbytes.Say("404 page not found"))
			})
		})

		Context("when the id is not found", func() {
			BeforeEach(func() {
				request.SetPathValue("id", "invalid")
			})

			It("returns 404 not found", func() {
				postStore := tmpbbs.NewPostStore("test-title")
				postGetHandler := tmpbbs.NewPostGetHandler(10, []string{}, true, true, postStore)
				postGetHandler.ServeHTTP(fakeResponseWriter, request)

				Expect(fakeResponseWriter.WriteHeaderCallCount()).To(Equal(1))
				Expect(fakeResponseWriter.WriteHeaderArgsForCall(0)).To(Equal(http.StatusNotFound))
				Eventually(buffer).Should(gbytes.Say("404 page not found"))
			})
		})

		Context("when the id is found", func() {
			BeforeEach(func() {
				request.SetPathValue("id", "1")
			})

			It("shows the post", func() {
				postStore := tmpbbs.NewPostStore("test-title")

				postPostHandler := tmpbbs.NewPostPostHandler(10, postStore, nil)
				postRequest, err := http.NewRequest(http.MethodPost, "", bytes.NewBufferString(`title=test-post-1`))
				Expect(err).ToNot(HaveOccurred())
				postRequest.Header.Add("Content-Type", "application/x-www-form-urlencoded")
				postRequest.SetPathValue("id", "0")
				postPostHandler.ServeHTTP(fakeResponseWriter, postRequest)

				postGetHandler := tmpbbs.NewPostGetHandler(10, []string{}, true, true, postStore)
				postGetHandler.ServeHTTP(fakeResponseWriter, request)

				Eventually(buffer).Should(gbytes.Say("test-post-1"))
			})
		})
	})
})
