package pkg

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/jenkins-zh/mirror-proxy/pkg/helper"
	"io"
	"net/http"
	"net/url"
	"strings"
)

// AddContext add context inject all handlers
func AddContext(next http.Handler, option *ServerOptions) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		ctx := context.WithValue(r.Context(), context.TODO(), *option)
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}

// GetUpdateCenterQuery returns the query object
func GetUpdateCenterQuery(querySources ...QuerySource) (query UpdateCenterQuery) {
	for _, querySource := range querySources {
		if version := querySource.Get("version"); version != "" {
			query.Version = version
		}

		if provider := querySource.Get("mirror-provider"); provider != "" {
			query.Provider = provider
		}

		if jsonServer := querySource.Get("mirror-jsonServer"); jsonServer != "" {
			query.JSONServer = jsonServer
		}

		if experimental := querySource.Get("mirror-experimental"); experimental == "true" {
			query.Experimental = true
		}
	}
	return
}

// HandleUpdateCenter handle GET /update-center.json
func HandleUpdateCenter(w http.ResponseWriter, r *http.Request) {
	o := r.Context().Value(context.TODO()).(ServerOptions)
	query := GetUpdateCenterQuery(r.URL.Query(), r.Header)

	o.WorkPool.AddTask(Task{
		TaskFunc: func(_ interface{}) {
			pluginDownloadCounter := &GitPluginDownloadCounter{
				Path: o.DataFilePath,
			}
			if err := pluginDownloadCounter.RecordUpdateCenterVisitData(); err != nil {
				fmt.Println(err)
			}
		},
	})

	var err error
	var targetURL *url.URL
	if targetURL, err = o.GetAndCacheURL(query); err == nil {
		w.Header().Set("Location", o.GetProviderURL(targetURL, query))
		w.WriteHeader(http.StatusMovedPermanently)
	} else {
		w.WriteHeader(http.StatusNotFound)

		_, err = w.Write([]byte(fmt.Sprintf("%v", err)))
	}
	helper.CheckErr(o.Printer, err)
}

// HandleJSONServers handle /json-servers
func HandleJSONServers(w http.ResponseWriter, r *http.Request) {
	o := r.Context().Value(context.TODO()).(ServerOptions)
	servers := GetJSONServers()
	if _, ok := servers["default"]; ok {
		servers["default"] = o.DefaultJSONServer
	} else {
		match := false
		for _, val := range servers {
			if val == o.DefaultJSONServer {
				match = true
				break
			}
		}

		if !match {
			servers["default"] = o.DefaultJSONServer
		}
	}

	data, err := json.Marshal(GetJSONServers())
	if err == nil {
		_, err = w.Write(data)
	}
	helper.CheckErr(o.Printer, err)
}

// HandleProviders handle /providers
func HandleProviders(w http.ResponseWriter, r *http.Request) {
	o := r.Context().Value(context.TODO()).(ServerOptions)

	providers := GetProviders()
	includeDefaultProvider := false
	for _, provider := range providers {
		if provider == o.DefaultProvider {
			includeDefaultProvider = true
			break
		}
	}

	if !includeDefaultProvider {
		providers = append(providers, o.DefaultProvider)
	}

	var writeErr error
	if data, err := json.Marshal(providers); err == nil {
		_, writeErr = w.Write(data)
	} else {
		w.WriteHeader(500)
		_, writeErr = io.WriteString(w, fmt.Sprintf("%v", err))
	}
	helper.CheckErr(o.Printer, writeErr)
}

// HandleDefaultProvider handle /providers/default
func HandleDefaultProvider(w http.ResponseWriter, r *http.Request) {
	var writeErr error
	o := r.Context().Value(context.TODO()).(ServerOptions)
	_, writeErr = io.WriteString(w, o.DefaultProvider)
	helper.CheckErr(o.Printer, writeErr)
}

// HandlePluginDownload as a proxy of plugin download
func HandlePluginDownload(w http.ResponseWriter, r *http.Request) {
	queryValues := r.URL.Query()

	var providerHost string
	provider := queryValues.Get("provider")
	switch provider {
	default:
		provider = "tsinghua"
		providerHost = "https://mirrors.tuna.tsinghua.edu.cn"
	}

	uri := r.RequestURI
	uri = strings.Split(uri, "?")[0]

	o := r.Context().Value(context.TODO()).(ServerOptions)
	o.WorkPool.AddTask(Task{
		TaskFunc: func(_ interface{}) {
			pluginDownloadCounter := &GitPluginDownloadCounter{
				Path: o.DataFilePath,
			}

			index := strings.LastIndex(uri, "/")
			pluginName := uri[index+1:]
			pluginName = strings.Split(pluginName, ".")[0]

			if err := pluginDownloadCounter.RecordPluginDownloadData(pluginName, provider); err != nil {
				fmt.Println(err)
			}
		},
	})

	w.Header().Set("Location", fmt.Sprintf("%s%s", providerHost, uri))
	w.WriteHeader(http.StatusMovedPermanently)
}

// HandlePluginsData returns the data of a plugin
func HandlePluginsData(w http.ResponseWriter, r *http.Request) {
	queryValues := r.URL.Query()

	year := queryValues.Get("year")
	name := queryValues.Get("name")

	// use current year as the default
	if year == "" {
		year = GetCurrentYear()
	}

	// make sure we can return a data
	if name == "" {
		name = "update-center"
	}

	fmt.Println("plugin", name, "year", year)

	// get plugin data
	o := r.Context().Value(context.TODO()).(ServerOptions)
	pluginDownloadCounter := &GitPluginDownloadCounter{
		Path: o.DataFilePath,
	}

	pluginData, err := pluginDownloadCounter.FindPluginData(year, name)

	responseData := ResponseData{
		Data: pluginData,
		Error: err,
	}

	var data []byte
	if data, err = json.Marshal(responseData); err != nil {
		fmt.Println(err)
	}

	w.Write(data)
}

func HandlePluginsDataList(w http.ResponseWriter, r *http.Request) {
	queryValues := r.URL.Query()

	year := queryValues.Get("year")

	// use current year as the default
	if year == "" {
		year = GetCurrentYear()
	}

	// get plugin data
	o := r.Context().Value(context.TODO()).(ServerOptions)
	pluginDownloadCounter := &GitPluginDownloadCounter{
		Path: o.DataFilePath,
	}

	responseData := ResponseData {}
	if downloadData, err := pluginDownloadCounter.FindByYear(year); err == nil {
		plugins := make([]string, 0)
		for key, _ := range downloadData.Plugins {
			plugins = append(plugins, key)
		}

		responseData.Data = plugins
	} else {
		responseData.Error = err
	}

	var data []byte
	var err error
	if data, err = json.Marshal(responseData); err != nil {
		fmt.Println(err)
	}

	w.Write(data)
}

type ResponseData struct {
	Data interface{}
	Error error
}
