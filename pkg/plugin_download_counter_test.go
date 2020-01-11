package pkg_test

import (
	server "github.com/jenkins-zh/mirror-proxy/pkg"
	"github.com/magiconair/properties/assert"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"os"
	"testing"
)

var _ = Describe("GitPluginDownloadCounter", func() {
	var (
		counter server.PluginDownloadCounter
		year    string
	)

	JustBeforeEach(func() {
		counter = &server.GitPluginDownloadCounter{}
		year = server.GetCurrentYear()
	})

	AfterEach(func() {
		gitCounter := counter.(*server.GitPluginDownloadCounter)
		err := os.RemoveAll(gitCounter.GetDataFilePath(year))
		Expect(err).NotTo(HaveOccurred())
	})

	It("ReadData", func() {
		_, err := counter.ReadData()
		Expect(err).NotTo(HaveOccurred())
	})

	It("FindByYear", func() {
		data, err := counter.FindByYear(year)
		Expect(err).NotTo(HaveOccurred())
		Expect(data).NotTo(BeNil())
	})

	It("Save", func() {
		data := &server.PluginDownloadData{
			Year:    year,
			Plugins: map[string]server.PluginData{},
		}
		data.Plugins["update-center"] = server.PluginData{Data: map[string]int64{
			server.GetDate(): 1,
		}}
		err := counter.Save(data)
		Expect(err).NotTo(HaveOccurred())

		var resultData *server.PluginDownloadData
		resultData, err = counter.FindByYear(year)
		Expect(err).NotTo(HaveOccurred())
		Expect(resultData).NotTo(BeNil(), "cannot found the saved data")
	})
})

func TestGetCurrentYear(t *testing.T) {
	currentYear := server.GetCurrentYear()
	assert.Equal(t, currentYear, "2020")
}

func TestGetDate(t *testing.T) {
	currentYear := server.GetDate()
	assert.Equal(t, currentYear, "2020-01-11")
}
