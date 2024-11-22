package tmpbbs_test

import (
	"github.com/mmb/tmpbbs/internal/tmpbbs"
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
)

var _ = Describe("poststore", func() {
	Describe("LoadYAML", func() {
		postStore := tmpbbs.NewPostStore("test-title")
		var path string

		Context("when reading the file returns an error", func() {
			BeforeEach(func() {
				path = "/does/not/exist.yml"
			})

			It("returns the error", func() {
				err := postStore.LoadYAML(path, nil)
				Expect(err).To(MatchError(ContainSubstring("no such file or directory")))
			})
		})

		Context("when the YAML is invalid", func() {
			BeforeEach(func() {
				path = "fixtures/invalid.yml"
			})

			It("returns the error", func() {
				err := postStore.LoadYAML(path, nil)
				Expect(err).To(MatchError(ContainSubstring("unmarshal errors")))
			})
		})

		Context("when the YAML is valid", func() {
			BeforeEach(func() {
				path = "fixtures/posts.yml"
			})

			It("loads the posts in the YAML file", func() {
				tripCoder, err := tmpbbs.NewTripCoder("test-salt", nil)
				Expect(err).ToNot(HaveOccurred())
				err = postStore.LoadYAML(path, tripCoder)
				Expect(err).ToNot(HaveOccurred())
			})
		})
	})
})
