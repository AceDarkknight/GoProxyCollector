package collector

import (
	"fmt"
	"strings"

	_ "github.com/AceDarkkinght/GoProxyCollector/util"
	"github.com/parnurzeal/gorequest"
)

func CollectXici(url string) *[]Result {
	if !strings.HasPrefix(url, "http://www.xicidaili.com") {
		return nil
	}

	resp, body, errs := gorequest.New().
		Get(url).
		End()

	if errs != nil {
		fmt.Println(errs)
		return nil
	}

	if resp.StatusCode != 200 {
		return nil
	}

	fmt.Print(body)

	return nil
}
