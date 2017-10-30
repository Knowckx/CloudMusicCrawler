package crawlerpkg

import (
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"log"
	"math/rand"
	"net/http"
	"net/url"
	"strings"
	"time"
)

//proxyClient 1.返回代理的Client
func proxyClient(proxyIP string, timeOut int) *http.Client {
	var proxy func(*http.Request) (*url.URL, error)
	if proxyIP != "127.0.0.1" {
		proxy = func(_ *http.Request) (*url.URL, error) {
			return url.Parse(fmt.Sprintf("http://%s", proxyIP))
		}
	} else {
		proxy = nil
	}

	t := &http.Transport{Proxy: proxy} //RoundTripper  设置代理
	c := http.Client{Timeout: time.Duration(timeOut) * time.Second, Transport: t}
	return &c
}

//SetGetReq  2.返回get时的req
func SetGetReq(UID int) (*http.Request, error) {
	urlStr := fmt.Sprintf("http://music.163.com/user/home?id=%d", UID)
	request, err := http.NewRequest("GET", urlStr, strings.NewReader(""))
	SetHeader(request)
	return request, err
}

//SetPostReq  2.返回post时的req
func SetPostReq(songID string, page uint32) (*http.Request, error) {
	//RequestParam，一个结构体，对应post时的结构
	params := RequestParam{Offset: page * 100, Limit: 100, Rid: songID, CsrfToken: ""} //uint64(page) * 100,
	body, err := Encrypt(&params)                                                      //把结构体通过AES和Rsa加密了
	if err != nil {
		log.Fatalln("参数加密出错")
	}
	urlStr := fmt.Sprintf("http://music.163.com/weapi/v1/resource/comments/R_SO_4_%s", songID)

	v := url.Values{}
	v.Set("params", body.Params)
	v.Add("encSecKey", body.EncSecKey)
	request, err := http.NewRequest("POST", urlStr, strings.NewReader(v.Encode()))
	SetHeader(request)
	return request, err
}

//SetHeader  2.1  为request 设置对应的Header
func SetHeader(request *http.Request) {
	userAgent := randomUserAgent() //随机userAgent
	// fmt.Println(userAgent)
	header := map[string]string{
		"Accept":          "*/*",
		"Accept-Language": "zh-CN,zh;q=0.8,gl;q=0.6,zh-TW;q=0.4",
		"Connection":      "keep-alive",
		"Content-Type":    "application/x-www-form-urlencoded",
		"Referer":         "http://music.163.com",
		"Host":            "music.163.com",
		"Cookie":          "",
		"User-Agent":      userAgent,
	}
	//为req设定header
	for k, v := range header {
		request.Header.Set(k, v)
	}
}

//randomUserAgent  2.1  随机Agent
func randomUserAgent() string {
	userAgentList := []string{
		"Mozilla/5.0 (Macintosh; Intel Mac OS X 10_12_5) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/59.0.3071.115 Safari/537.36",
		"Mozilla/5.0 (Macintosh; Intel Mac OS X 10_12_5) AppleWebKit/603.2.4 (KHTML, like Gecko) Version/10.1.1 Safari/603.2.4",
		"Mozilla/5.0 (Macintosh; Intel Mac OS X 10.12; rv:46.0) Gecko/20100101 Firefox/46.0",
		"Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/51.0.2704.103 Safari/537.36",
		"Mozilla/5.0 (Windows NT 10.0; Win64; x64; rv:46.0) Gecko/20100101 Firefox/46.0",
		"Mozilla/4.0 (compatible; MSIE 7.0; Windows NT 6.0)",
		"Mozilla/4.0 (compatible; MSIE 8.0; Windows NT 6.0; Trident/4.0)",
		"Mozilla/5.0 (compatible; MSIE 9.0; Windows NT 6.1; Trident/5.0)",
		"Mozilla/5.0 (compatible; MSIE 10.0; Windows NT 6.2; Win64; x64; Trident/6.0)",
		"Mozilla/5.0 (Windows NT 6.3; Win64, x64; Trident/7.0; rv:11.0) like Gecko",
		"Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/42.0.2311.135 Safari/537.36 Edge/13.10586",
	}
	r := rand.New(rand.NewSource(time.Now().UnixNano()))
	pos := r.Intn(len(userAgentList))
	return userAgentList[pos]
}

//doReq 3 发起req,处理,返回原始数据[]byte
func doReq(Req *http.Request, proxyIP string) ([]byte, error) {
	Client := proxyClient(proxyIP, 5)
	SetHeader(Req)
	resp, err := Client.Do(Req)
	if err != nil {
		return nil, err
	}
	resqHost := resp.Request.Host //有的代理IP被DNS劫持，不干净
	if !strings.Contains(resqHost, "163") {
		return nil, errors.New("Request is error")
	}
	if resp == nil || resp.StatusCode != 200 {
		return nil, fmt.Errorf(resp.Status)
	}
	//---检查完毕
	defer func() {
		if resp != nil {
			resp.Body.Close()
		}
	}()
	p, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}
	return p, err
}

//commentReq 4 某一页评论的查询
func commentReq(proxyIP string, page uint32) (*commentMusic, error) {
	Req, err := SetPostReq(songID, page)
	if err != nil {
		return nil, err
	}
	respComments, err := doReq(Req, proxyIP)
	if err != nil {
		return nil, err
	}
	var comment commentMusic
	err = json.Unmarshal(respComments, &comment)
	// fmt.Printf("%+v", comment)
	if err != nil {
		return nil, err
	}
	return &comment, nil
}

//userReq 5 查询用户信息
func getuserinfo(UID int, proxyIP string) {
	defer func() { //线程错误恢复
		if err := recover(); err != nil {
		}
		wgRquest.Done()
	}()
	userInfo, err := userinfoReq(UID, proxyIP)
	if err != nil {
		dealuserinfoerr(UID, err)
		return
	}
	rw.Lock()
	dataUserInfo.Listuser = append(dataUserInfo.Listuser, userInfo)
	rw.Unlock()
}

func userinfoReq(UID int, proxyIP string) (*UserInfo, error) {
	Req, err := SetGetReq(UID)
	if err != nil {
		return nil, err
	}
	useinfoPage, err := doReq(Req, proxyIP)
	if err != nil {
		return nil, err
	}
	userInfo := userinfoparse(useinfoPage, UID)
	// fmt.Printf("用户信息:\n%#v\n", *userInfo)
	return userInfo, nil
}
func dealuserinfoerr(UID int, err error) {
	fmt.Printf("用户信息%d没有获取到\n", UID)
	fmt.Printf("error %s\n", err)
	rw.Lock()
	errosUser = append(errosUser, UID) //错误记录
	rw.Unlock()
}
