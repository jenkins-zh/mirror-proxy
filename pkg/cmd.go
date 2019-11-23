package pkg

import (
	"fmt"
	"github.com/jenkins-zh/mirror-proxy/pkg/helper"
	"github.com/mitchellh/go-homedir"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"log"
	"net/http"
	"net/url"
	"os"
)

// ServerOptions represents the options for a server
type ServerOptions struct {
	Config string

	DefaultProvider   string
	DefaultJSONServer string

	Host    string
	Port    int
	PortLTS int

	CertFile string
	KeyFile  string
}

var serverOptions ServerOptions

var rootCmd = &cobra.Command{
	Use:   "mirror-proxy",
	Short: "mirror-proxy is the proxy of Jenkins Update Center",
	Run: func(cmd *cobra.Command, args []string) {
		err := serverOptions.Run(cmd, args)
		helper.CheckErr(err)
	},
}

func init() {
	cobra.OnInitialize(initConfig)

	rootCmd.Flags().StringVar(&serverOptions.Config, "config", "", "config file (default is $HOME/.mirror-proxy.yaml)")
	rootCmd.Flags().StringVarP(&serverOptions.DefaultProvider, "default-provider", "", "tsinghua",
		"The default provider of the update center mirror")
	rootCmd.Flags().StringVarP(&serverOptions.DefaultJSONServer, "default-json-server", "", "https://gitee.com/jenkins-zh/update-center-mirror/raw/master",
		"The default JSON server of the update center mirror")

	rootCmd.Flags().StringVarP(&serverOptions.Host, "host", "", "127.0.0.1",
		"The host of the server")
	rootCmd.Flags().IntVarP(&serverOptions.Port, "port", "", 7070,
		"The port of the server")
	rootCmd.Flags().IntVarP(&serverOptions.Port, "port-lts", "", 7071,
		"The LTS port of the server")

	rootCmd.Flags().StringVarP(&serverOptions.CertFile, "cert", "", "",
		"The cert file of the server")
	rootCmd.Flags().StringVarP(&serverOptions.KeyFile, "key", "", "",
		"The key file of the server")

	viper.BindPFlag("default-provider", rootCmd.PersistentFlags().Lookup("default-provider"))
	viper.BindPFlag("default-json-server", rootCmd.PersistentFlags().Lookup("default-json-server"))
	viper.BindPFlag("cert", rootCmd.PersistentFlags().Lookup("cert"))
	viper.BindPFlag("key", rootCmd.PersistentFlags().Lookup("key"))
}

func initConfig() {
	if serverOptions.Config != "" {
		// Use config file from the flag.
		viper.SetConfigFile(serverOptions.Config)
	} else {
		// Find home directory.
		home, err := homedir.Dir()
		if err != nil {
			log.Fatal(err)
		}

		// Search config in home directory with name ".mirror-proxy" (without extension).
		viper.AddConfigPath(home)
		viper.SetConfigName(".mirror-proxy")
	}

	viper.AutomaticEnv()

	if err := viper.ReadInConfig(); err == nil {
		fmt.Println("Using config file:", viper.ConfigFileUsed())
	}
}

// GetProviderURL get the update center URL from a provider
func (o *ServerOptions) GetProviderURL(official *url.URL, query UpdateCenterQuery) (targetURL string) {
	jsonServer, provider := query.JSONServer, query.Provider
	if provider == "" {
		provider = o.DefaultProvider
	}

	jsonServer, ok := GetJSONServers()[jsonServer]
	if !ok {
		jsonServer = o.DefaultJSONServer
	}

	targetURL = fmt.Sprintf("%s/%s%s", jsonServer, provider, official.RequestURI())
	return
}

// UpdateCenterQuery holds the info for query a update center
type UpdateCenterQuery struct {
	Version      string
	Provider     string
	JSONServer   string
	Experimental bool
}

// QuerySource which contains the methods to query
type QuerySource interface {
	Get(key string) string
}

// Run startup a server
func (o *ServerOptions) Run(cmd *cobra.Command, args []string) (err error) {
	mux := http.NewServeMux()

	mux.Handle("/update-center.json", AddContext(http.HandlerFunc(HandleUpdateCenter), o))
	mux.Handle("/json-servers", AddContext(http.HandlerFunc(HandleJSONServers), o))
	mux.Handle("/providers", AddContext(http.HandlerFunc(HandleProviders), o))
	mux.Handle("/providers/default", AddContext(http.HandlerFunc(HandleDefaultProvider), o))
	go func() {
		err := http.ListenAndServe(fmt.Sprintf("%s:%d", o.Host, o.Port), mux)
		if err != nil {
			log.Println(err)
		}
	}()
	err = http.ListenAndServeTLS(fmt.Sprintf("%s:%d", o.Host, o.PortLTS),
		o.CertFile, o.KeyFile, nil)
	return
}

// GetURL get the real URL from the official site
func (o *ServerOptions) GetURL(version string) (targetURL *url.URL, err error) {
	var (
		request  *http.Request
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
func (o *ServerOptions) GetAndCacheURL(query UpdateCenterQuery) (targetURL *url.URL, err error) {
	var cacheErr error

	if query.Experimental {
		return url.Parse("https://updates.jenkins.io/experimental/update-center.json")
	}

	version := query.Version
	cacheServer := FileSystemCacheServer{FileName: "cache.yaml"}
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

// Execute will execute the command
func Execute() {
	if err := rootCmd.Execute(); err != nil {
		os.Exit(1)
	}
}
