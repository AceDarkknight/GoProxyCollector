package server

import "net/http"

func NewServer() {
	http.HandleFunc("/httpIp", getIp)
	err := http.ListenAndServe(":8090", nil)
	if err != nil {
		panic(err)
	}
}

func getIp(w http.ResponseWriter, r *http.Request) {
	if r.Method == "GET" {
		w.Header().Add("content-type", "application/json")
	} else {
		w.WriteHeader(http.StatusMethodNotAllowed)
		w.Write([]byte(http.StatusText(http.StatusMethodNotAllowed)))
	}
}
