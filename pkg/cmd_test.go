package pkg_test

import (
	server "github.com/jenkins-zh/mirror-proxy/pkg"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"net/url"
)

//var _ = Describe("server cmd test", func() {
//	It("", func() {
//		rootCmd := server.GetRootCmd()
//		rootCmd.SetArgs([]string{""})
//
//		buf := new(bytes.Buffer)
//		rootCmd.SetOutput(buf)
//
//		_, err := rootCmd.ExecuteC()
//
//		Expect(err).To(BeNil())
//		Expect(buf.String()).To(Equal("open : no such file or directory\n"))
//	})
//})

var _ = Describe("GetAndCacheURL", func() {
	var (
		opt       *server.ServerOptions
		query     server.UpdateCenterQuery
		cachedURL *url.URL
		err       error
	)

	BeforeEach(func() {
		opt = &server.ServerOptions{}
	})

	JustBeforeEach(func() {
		cachedURL, err = opt.GetAndCacheURL(query)
	})

	Context("experimental case", func() {
		BeforeEach(func() {
			query = server.UpdateCenterQuery{
				Experimental: true,
			}
		})

		It("get experimental URL", func() {
			expectURL, expectErr := url.Parse("https://updates.jenkins.io/experimental/update-center.json")
			Expect(cachedURL).To(Equal(expectURL))
			Expect(err).To(BeNil())
			Expect(expectErr).To(BeNil())
		})
	})

	Context("get from cache file", func() {
	})

	Context("can not get from cache file", func() {
	})
})
