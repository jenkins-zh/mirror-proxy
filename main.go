package main

import (
	"log"
	"net/http"
)

func main()  {
	http.HandleFunc("/update-center.json", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Location", "https://jenkins-zh.gitee.io/update-center-mirror/tsinghua/update-center.json")
		w.WriteHeader(301)
	})

	err := http.ListenAndServeTLS(":7899", "demo.crt", "demo.key", nil)
	log.Fatal(err)
}
