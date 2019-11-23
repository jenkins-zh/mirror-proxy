package pkg_test

import (
	cache "github.com/jenkins-zh/mirror-proxy/pkg"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"os"
)

var _ = Describe("cache test", func() {
	var (
		cacheFile   string
		cacheServer cache.CacheServer
		saveErr     error

		key   string
		value string
	)

	BeforeEach(func() {
		cacheFile = ""
		key = "key"
		value = "value"
	})

	JustBeforeEach(func() {
		cacheServer = &cache.FileSystemCacheServer{FileName: cacheFile}
		saveErr = cacheServer.Save(key, value)
	})

	Context("FileSystemCacheServer", func() {
		It("get from an empty cache file", func() {
			notExists := cacheServer.Load("fake")
			Expect(notExists).To(Equal(""))

			notExists = cacheServer.Load("")
			Expect(notExists).To(Equal(""))
		})

		It("save cache without given file", func() {
			Expect(saveErr).NotTo(BeNil())
		})

		Context("FileSystemCacheServer with file", func() {
			BeforeEach(func() {
				cacheFile = "cache.yaml"
			})

			AfterEach(func() {
				err := os.Remove(cacheFile)
				Expect(err).To(BeNil())
			})

			It("save cache with given file", func() {
				Expect(saveErr).To(BeNil())
			})

			It("get cache item", func() {
				val := cacheServer.Load(key)
				Expect(val).To(Equal(value))
			})
		})
	})
})
