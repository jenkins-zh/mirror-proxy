package pkg

import (
	"fmt"
	"github.com/spf13/cobra"
	"log"
	"net/http"
	"net/url"
	"os"
)

// ServerOptions represents the options for a server
type ServerOptions struct {
	Host string
	Port int

	CertFile string
	KeyFile string
}

var serverOptions ServerOptions

var rootCmd = &cobra.Command{
	Use:   "mirror-proxy",
	Short: "mirror-proxy is the proxy of Jenkins Update Center",
	Run: func(cmd *cobra.Command, args []string) {
		err := serverOptions.Run(cmd, args)
		log.Fatal(err)
	},
}

func init()  {
	rootCmd.Flags().StringVarP(&serverOptions.Host, "host", "", "127.0.0.1",
		"The host of the server")
	rootCmd.Flags().IntVarP(&serverOptions.Port, "port", "", 7070,
		"The port of the server")

	rootCmd.Flags().StringVarP(&serverOptions.CertFile, "cert", "", "",
		"The cert file of the server")
	rootCmd.Flags().StringVarP(&serverOptions.KeyFile, "key", "", "",
		"The key file of the server")
}

// Run startup a server
func (o *ServerOptions) Run(cmd *cobra.Command, args []string) (err error) {
	http.HandleFunc("/update-center.json", func(w http.ResponseWriter, r *http.Request) {
		version := r.URL.Query().Get("version")

		var targetURL *url.URL
		var err error
		if targetURL, err = o.GetAndCacheURL(version); err == nil {
			w.Header().Set("Location", fmt.Sprintf("https://jenkins-zh.gitee.io/update-center-mirror/tsinghua%s",
				targetURL.RequestURI()))
			w.WriteHeader(301)
		} else {
			w.WriteHeader(400)

			if _, err = w.Write([]byte(fmt.Sprintf("%v", err))); err != nil {
				log.Println(err)
			}
		}
	})

	err = http.ListenAndServeTLS(fmt.Sprintf("%s:%d", o.Host, o.Port),
		o.CertFile, o.KeyFile, nil)
	return
}

// GetURL get the real URL from the official site
func (o *ServerOptions) GetURL(version string) (targetURL *url.URL, err error) {
	var (
		request *http.Request
		response *http.Response
	)

	api := fmt.Sprintf("https://updates.jenkins.io/update-center.json?version=%s", version)
	request, err = http.NewRequest("GET", api, nil)
	if err == nil {
		client := &http.Client{
			CheckRedirect: func(req *http.Request, via []*http.Request) error {
				return http.ErrUseLastResponse
			},
		}

		response, err = client.Do(request)
		if err == nil {
			targetURL, err = response.Location()
		}
	}
	return
}

// GetAndCacheURL get the real URL, then cache it
func (o *ServerOptions) GetAndCacheURL(version string) (targetURL *url.URL, err error) {
	var cacheErr error
	cacheServer := FileSystemCacheServer{FileName:"cache.yaml"}
	if cacheURL := cacheServer.Load(version); cacheURL != "" {
		targetURL, cacheErr = url.Parse(cacheURL)
	} else {
		if targetURL, err = o.GetURL(version); err == nil {
			if cacheErr = cacheServer.Save(version, targetURL.String()); cacheErr != nil {
				log.Println(cacheErr)
			}
		}
	}

	if cacheErr != nil {
		if targetURL, err = o.GetURL(version); err == nil {
			if cacheErr = cacheServer.Save(version, targetURL.String()); cacheErr != nil {
				log.Println(cacheErr)
			}
		}
	}
	return
}

// Execute will exectue the command
func Execute() {
	if err := rootCmd.Execute(); err != nil {
		os.Exit(1)
	}
}
