package httpproxy
/**
* http request client wapper
* author liuyoubin@baidu.com
**/
import (
	"fmt"
	"net"
	"net/url"
	"net/http"
	"strings"
	"time"
	"io"
	"io/ioutil"
	"bytes"
)

type HttpError struct{
	Errno int
	Errmsg string
	Err error
}

func (httpError *HttpError) Error() string{
	if httpError.Err != nil {
		return fmt.Sprintf("Errno:%d, Errmsg:%s, Error:%s",
			httpError.Errno, httpError.Errmsg, httpError.Err.Error())
	}else{
		return fmt.Sprintf("Errno:%d, Errmsg:%s, Error:nil",
			httpError.Errno, httpError.Errmsg)
	}
	
}

type HttpProxy struct{
	Url string
	Headers map[string]string
	Cookies map[string]string
	Method string
	PostBody string
	QueryString string
	ConnectTimeout time.Duration
	ReadWriteTimeout time.Duration
	MaxRedirects int
	ResponseStatusCode int
	ResponseHeader string
	ResponseBody string
}

func (httpProxy *HttpProxy) AddHeader(name string, value string){
	if len(httpProxy.Headers) == 0 {
		httpProxy.Headers = make(map[string]string)
	}
	(*httpProxy).Headers[name] = value
}

func (httpProxy *HttpProxy) SetReadWriteTimeout(timeout time.Duration){
	(*httpProxy).ReadWriteTimeout = timeout
}

func (httpProxy *HttpProxy) SetConnectTimeout(timeout time.Duration){
	(*httpProxy).ConnectTimeout = timeout
}

func (httpProxy *HttpProxy) AddCookie(name string, value string){
	if len(httpProxy.Cookies) == 0 {
		httpProxy.Cookies = make(map[string]string)
	}
	(*httpProxy).Cookies[name] = value
}

func (httpProxy *HttpProxy) GetResponseBody() string{
	return httpProxy.ResponseBody
}

func (httpProxy *HttpProxy) GetResponseHeader() string{
	return httpProxy.ResponseHeader
}

func (httpProxy *HttpProxy) GetResponseStatusCode() int{
	return httpProxy.ResponseStatusCode
}

func timeoutDialer(cTimeout time.Duration, rwTimeout time.Duration) func(net, addr string) (c net.Conn, err error) {
    return func(netw, addr string) (net.Conn, error) {
        conn, err := net.DialTimeout(netw, addr, cTimeout)
        if err != nil {
            return nil, err
        }
        conn.SetDeadline(time.Now().Add(rwTimeout))
        return conn, nil
    }
}

func (httpProxy *HttpProxy) DoRequest() (err error){
    var DefaultTransport http.RoundTripper = &http.Transport{
    	Dial: timeoutDialer(httpProxy.ConnectTimeout * time.Millisecond,
    		httpProxy.ReadWriteTimeout * time.Millisecond),
    }
    var DefaultClient = &http.Client{Transport: DefaultTransport}
    
    DefaultClient.CheckRedirect = func(req *http.Request, via[]*http.Request) error{
    	if len(via) > httpProxy.MaxRedirects {
    		return &HttpError{Errno:1, Errmsg:"reach max redirects", Err:nil}
    	}
    	return nil
    }
    
    var bodyReader io.Reader = nil
    if httpProxy.Method == "POST" && httpProxy.PostBody != "" {
    	bodyReader = strings.NewReader(httpProxy.PostBody)
    }
    
    req, err := http.NewRequest(httpProxy.Method, httpProxy.Url, bodyReader)
    if err != nil {
    	return &HttpError{Errno:1,Errmsg:"http new request failed", Err:err}
    }
    
    if len(httpProxy.Headers) > 0 {
    	for k, v := range(httpProxy.Headers) {
    		req.Header.Add(k,v)
    	}
    }
    
    if len(httpProxy.Cookies) > 0 {
    	for k, v := range(httpProxy.Cookies) {
    		req.AddCookie(&http.Cookie{Name:k,Value:v})
    	}
    }
    
    res, err := DefaultClient.Do(req)
    if err != nil {
    	return &HttpError{Errno:1, Errmsg:"request failed", Err:err}
    }
    defer res.Body.Close()
    
    (*httpProxy).ResponseStatusCode = res.StatusCode
    var b bytes.Buffer
	err = res.Header.Write(&b)
	if err != nil {
    	return &HttpError{Errno:1, Errmsg:"get response header failed", Err:err}
    }
    (*httpProxy).ResponseHeader = b.String()
    
    body, err := ioutil.ReadAll(res.Body)
    if err != nil {
    	return &HttpError{Errno:1, Errmsg:"get response body failed", Err:err}
    }
    
	(*httpProxy).ResponseBody = string(body)
	
	return nil
}

func (httpProxy *HttpProxy)Get(_url string,	_rawQueryMap map[string]string) (err error){
	httpProxy.Method = "GET"
	
	queryStr := ""
	
	if _rawQueryMap != nil {
		for k, v := range(_rawQueryMap) {
			queryStr += k+"="+url.QueryEscape(v)+"&"
		}
	}
	
	if queryStr != "" {
		queryStr = strings.TrimRight(queryStr, "&")
		if strings.Index(_url, "?") > 0{
			_url = _url + "&" + queryStr
		}else{
			_url = _url + "?" + queryStr
		}
	}
	
	httpProxy.Url = _url
	
	err = httpProxy.DoRequest()
	
	if err != nil {
		return err
	}
	return nil
}


func (httpProxy *HttpProxy)Post(_url string, _rawPostMap map[string]string) (err error){
	(*httpProxy).Method = "POST"
	(*httpProxy).Url = _url
	(*httpProxy).AddHeader("Content-type", "application/x-www-form-urlencoded")
	postBody := ""
	
	if _rawPostMap != nil {
		for k, v := range(_rawPostMap) {
			postBody += k+"="+url.QueryEscape(v)+"&"
		}
	}
	
	if postBody != "" {
		postBody = strings.TrimRight(postBody, "&")
	}
	
	(*httpProxy).PostBody = postBody
	
	err = httpProxy.DoRequest()
	
	if err != nil {
		return err
	}
	
	return nil
}

