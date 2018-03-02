package verifier

import (
	"strconv"

	"github.com/parnurzeal/gorequest"
)

func VerifierHTTP(ip string, port int) bool {
	if ip == "" {
		return false
	}

	if port <= 0 {
		return false
	}

	proxy := "http://" + ip + ":" + strconv.Itoa(port)
	resp, _, errs := gorequest.New().
		Proxy(proxy).
		Get("http://httpbin.org/get").
		End()

	if errs != nil {
		return false
	}

	if resp.StatusCode != 200 {
		return false
	}

	return true
}
