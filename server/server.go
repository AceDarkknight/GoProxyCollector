package server

import (
	"net/http"

	"github.com/AceDarkkinght/GoProxyCollector/storage"
)

var s storage.Storage

func NewServer(storage storage.Storage) {
	if storage != nil {
		s = storage
	}

	http.HandleFunc("/ip", getIp)
	err := http.ListenAndServe(":8090", nil)
	if err != nil {
		panic(err)
	}
}

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
