package utils_test

import (
	"os"

	"github.com/golgoth31/release-installer/pkg/utils"
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
)

var _ = Describe("Files", func() {
	Describe("Loading yaml files", func() {
		It("should be ok", func() {
			out, err := utils.Load("test_data/ok.yaml")
			Expect(err).To(BeNil())
			Expect(string(out)).To(Equal("{\"this\":\"is a string\"}"))
		})
		It("not found", func() {
			_, err := utils.Load("test_data/bad.yaml")
			Expect(err).To(MatchError(os.ErrNotExist))
		})
	})
	Describe("Moving files", func() {
		It("should be ok", func() {
			err := utils.MoveFile("test_data/ok.yaml", "test_data/moved.yaml", 0600)
			err1 := utils.MoveFile("test_data/moved.yaml", "test_data/ok.yaml", 0600)
			Expect(err).To(BeNil())
			Expect(err1).To(BeNil())
		})
		It("not found", func() {
			err := utils.MoveFile("test_data/bad.yaml", "test_data/moved.yaml", 0600)
			Expect(err).To(MatchError(os.ErrNotExist))
		})
		It("not found", func() {
			err := utils.MoveFile("test_data/ok.yaml", "test_data1/moved.yaml", 0600)
			Expect(err).To(MatchError(os.ErrNotExist))
		})
	})
})
