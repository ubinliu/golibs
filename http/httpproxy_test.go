package httpproxy

import (
	"fmt"
    "testing"
)

func TestGet(t *testing.T) {
	httpProxy := &HttpProxy{ConnectTimeout: 300, ReadWriteTimeout: 1000}
	urlstr := "http://www.baidu.com/s"
	queryMap := make(map[string]string)
	queryMap["wd"] = "baiduyun"
	
	err := httpProxy.Get(urlstr, queryMap)
	if err != nil {
		fmt.Println(err.Error())
		return
	}
	fmt.Println(httpProxy.ResponseStatusCode)
	fmt.Println(httpProxy.ResponseHeader)
	fmt.Println(httpProxy.ResponseBody)
}

func TestCookie(t *testing.T) {
	httpProxy := &HttpProxy{ConnectTimeout: 300, ReadWriteTimeout: 1000}
	urlstr := "http://pan.baidu.com/api/list?channel=chunlei&clienttype=0&web=1&num=100&page=1&dir=%2F&order=time&desc=1&showempty=0&_=1439836009675&bdstoken=54d0c0937ca0bee02d67649b146854d7&channel=chunlei&clienttype=0&web=1&app_id=250528"
	
	httpProxy.AddCookie("BDUSS","TZZT1g3OUcxY35NQUNyd2dBRFVVaWxqWFNNLWNDRVRBUXdsOGlwOUlOV2tULXBWQVFBQUFBJCQAAAAAAAAAAAEAAADfh6UFeW91YmlubGl1AAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAKTCwlWkwsJVc")
	err := httpProxy.Get(urlstr, nil)
	if err != nil {
		fmt.Println(err.Error())
		return
	}
	
	fmt.Println(httpProxy.ResponseBody)
}


