package server

import (
	"net/http"

	"github.com/AceDarkknight/GoProxyCollector/storage"
	
    "github.com/cihub/seelog"
)

var s storage.Storage

// NewServer will start a new server.
func NewServer(storage storage.Storage) {
	if storage != nil {
		s = storage
	}
	
	defer func() {
	    if r:=recover();r!=nil{
	        seelog.Critical(r)
        }
    }()

	http.HandleFunc("/get", getIp)
	http.HandleFunc("/delete", deleteIp)
	err := http.ListenAndServe(":8090", nil)
	if err != nil {
		panic(err)
	}
}

// getIp will get a random Ip.
// Sample usage: http://localhost:8090/get
func getIp(w http.ResponseWriter, r *http.Request) {
	if r.Method == "GET" {
		w.Header().Add("content-type", "application/json")
		if s == nil {
			w.WriteHeader(http.StatusInternalServerError)
		}

		_, result := s.GetRandomOne()
		if len(result) > 0 {
			w.Write(result)
		} else {
			w.WriteHeader(http.StatusNotFound)
		}
	} else {
		w.WriteHeader(http.StatusMethodNotAllowed)
	}
}

// deleteIp will delete the given ip. Return 200 if succeed.
// Sample usage: http://localhost:8090/delete?ip=0.0.0.0
func deleteIp(w http.ResponseWriter, r *http.Request) {
	if r.Method == "GET" {
		values := r.URL.Query()
		if len(values["ip"]) > 1 {
			w.WriteHeader(http.StatusInternalServerError)
		}

		if s.Delete(values["ip"][0]) {
			w.WriteHeader(http.StatusOK)
		}
	} else {
		w.WriteHeader(http.StatusMethodNotAllowed)
	}
}
