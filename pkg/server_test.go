package pkg_test

import (
	"context"
	server "github.com/jenkins-zh/mirror-proxy/pkg"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
)

var _ = Describe("test server", func() {
	var (
		option server.ServerOptions
		api    string

		request *http.Request
		reqErr  error

		recorder *httptest.ResponseRecorder

		bodyData []byte
		bodyErr  error

		reqHandler http.HandlerFunc
	)

	BeforeEach(func() {
		option = server.ServerOptions{}
		option.WorkPool = &server.WorkPool{}
		option.WorkPool.InitPool(5)
	})

	JustBeforeEach(func() {
		request, reqErr = http.NewRequest("GET", api, nil)
		Expect(reqErr).To(BeNil())

		recorder = httptest.NewRecorder()

		ctx := request.Context()
		ctx = context.WithValue(ctx, context.TODO(), option)

		request = request.WithContext(ctx)
		reqHandler.ServeHTTP(recorder, request)

		bodyData, bodyErr = ioutil.ReadAll(recorder.Body)
		Expect(bodyErr).To(BeNil())
	})

	Context("HandleUpdateCenter", func() {
		BeforeEach(func() {
			api = "/update-center.json"
			reqHandler = server.HandleUpdateCenter
		})

		It("should success", func() {
			Expect(recorder.Code).To(Equal(http.StatusMovedPermanently))
		})
	})

	Context("HandleJSONServers", func() {
		BeforeEach(func() {
			api = "/json-servers"
			reqHandler = server.HandleJSONServers
		})

		It("should success", func() {
			Expect(recorder.Code).To(Equal(http.StatusOK))
			Expect(string(bodyData)).To(Equal("{}"))
		})

		Context("only with one default JSON server", func() {
			BeforeEach(func() {
				option = server.ServerOptions{
					DefaultJSONServer: "fake",
				}
			})

			It("only default JSON server", func() {
				Expect(string(bodyData)).To(Equal("{}"))
			})
		})
	})

	Context("HandleProviders", func() {
		BeforeEach(func() {
			api = "/providers"
			reqHandler = server.HandleProviders
		})

		It("should success", func() {
			Expect(recorder.Code).To(Equal(http.StatusOK))
			Expect(string(bodyData)).To(Equal(`[""]`))
		})

		Context("with default provider", func() {
			BeforeEach(func() {
				option = server.ServerOptions{
					DefaultProvider: "fake",
				}
			})

			It("only default provider", func() {
				Expect(string(bodyData)).To(Equal(`["fake"]`))
			})
		})
	})

	Context("HandleDefaultProvider", func() {
		BeforeEach(func() {
			api = "/providers/default"
			reqHandler = server.HandleDefaultProvider
		})

		It("should success", func() {
			Expect(recorder.Code).To(Equal(http.StatusOK))
			Expect(string(bodyData)).To(Equal(""))
		})

		Context("with default provider", func() {
			BeforeEach(func() {
				option = server.ServerOptions{
					DefaultProvider: "fake",
				}
			})

			It("only default provider", func() {
				Expect(string(bodyData)).To(Equal("fake"))
			})
		})
	})

	Context("HandleHealthCheck", func() {
		BeforeEach(func() {
			api = "/status"
			reqHandler = server.HandleHealthCheck
		})

		It("should return ok", func() {
			Expect(string(bodyData)).To(Equal("ok"))
		})
	})
})

var _ = Describe("GetUpdateCenterQuery", func() {
	var (
		query        server.UpdateCenterQuery
		querySources []server.QuerySource
	)

	JustBeforeEach(func() {
		query = server.GetUpdateCenterQuery(querySources...)
	})

	It("without input, should match with the default value", func() {
		Expect(query.Experimental).To(BeFalse())
		Expect(query.Version).To(BeEmpty())
		Expect(query.Provider).To(BeEmpty())
		Expect(query.JSONServer).To(BeEmpty())
	})

	Context("give one query", func() {
		BeforeEach(func() {
			querySources = []server.QuerySource{
				&FakeQuery{Key: "version", Value: "fake"},
				&FakeQuery{Key: "mirror-provider", Value: "fake"},
				&FakeQuery{Key: "mirror-jsonServer", Value: "fake"},
				&FakeQuery{Key: "mirror-experimental", Value: "true"},
			}
		})

		It("fake version", func() {
			Expect(query.Version).To(Equal("fake"))
			Expect(query.Provider).To(Equal("fake"))
			Expect(query.JSONServer).To(Equal("fake"))
			Expect(query.Experimental).To(BeTrue())
		})
	})
})

// FakeQuery only for test
type FakeQuery struct {
	Key   string
	Value string
}

// Get only for test
func (f *FakeQuery) Get(key string) (val string) {
	if key == f.Key {
		val = f.Value
	}
	return
}
