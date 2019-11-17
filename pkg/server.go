package pkg

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/jenkins-zh/mirror-proxy/pkg/helper"
	"io"
	"net/http"
	"net/url"
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

	var targetURL *url.URL
	var err error
	if targetURL, err = o.GetAndCacheURL(query); err == nil {
		w.Header().Set("Location", o.GetProviderURL(targetURL, query))
		w.WriteHeader(http.StatusMovedPermanently)
	} else {
		w.WriteHeader(http.StatusNotFound)

		_, err = w.Write([]byte(fmt.Sprintf("%v", err)))
	}
	helper.CheckErr(err)
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
	helper.CheckErr(err)
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
	helper.CheckErr(writeErr)
}

// HandleDefaultProvider handle /providers/default
func HandleDefaultProvider(w http.ResponseWriter, r *http.Request) {
	var writeErr error
	o := r.Context().Value(context.TODO()).(ServerOptions)
	_, writeErr = io.WriteString(w, o.DefaultProvider)
	helper.CheckErr(writeErr)
}
