package tmpbbs_test

import (
	"errors"

	"github.com/mmb/tmpbbs/internal/tmpbbs"
	"github.com/mmb/tmpbbs/internal/tmpbbsfakes"
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
)

var _ = Describe("tripcoder", func() {
	Describe("NewTripCoder", func() {
		var fakeReader *tmpbbsfakes.FakeReader
		var salt string

		BeforeEach(func() {
			fakeReader = &tmpbbsfakes.FakeReader{}
		})

		Context("when a salt is passed in", func() {
			BeforeEach(func() {
				salt = "test-salt"
			})

			It("returns a tripCoder", func() {
				tripCoder, err := tmpbbs.NewTripCoder(salt, fakeReader)
				Expect(err).ToNot(HaveOccurred())
				Expect(tripCoder).ToNot(BeNil())
				Expect(fakeReader.ReadCallCount()).To(Equal(0))
			})
		})

		Context("when no salt is passed in", func() {
			BeforeEach(func() {
				salt = ""
			})

			Context("when the reader read is successful", func() {
				BeforeEach(func() {
					fakeReader.ReadCalls(func(p []byte) (int, error) {
						s := "testtesttesttest"
						for i, c := range s {
							p[i] = byte(c)
						}

						return len(s), nil
					})
				})

				It("returns a tripCoder", func() {
					tripCoder, err := tmpbbs.NewTripCoder(salt, fakeReader)
					Expect(err).ToNot(HaveOccurred())
					Expect(tripCoder).ToNot(BeNil())
					Expect(fakeReader.ReadCallCount()).To(Equal(1))
				})
			})

			Context("when the reader read returns an error", func() {
				var err error

				BeforeEach(func() {
					err = errors.New("read error")
					fakeReader.ReadReturns(0, err)
				})

				It("returns the error", func() {
					_, err := tmpbbs.NewTripCoder("", fakeReader)
					Expect(err).To(Equal(err))
					Expect(fakeReader.ReadCallCount()).To(Equal(1))
				})
			})
		})
	})
})
