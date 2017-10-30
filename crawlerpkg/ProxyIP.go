package crawlerpkg

import (
	"fmt"
	"io/ioutil"
	"net/http"
	"os"
	"strings"
	"sync"
)

//拿到可用Ip
func getOkProxyIP() []string {
	ip := getIPfromFile()
	fmt.Println("本地缓存代理IP数:", len(ip))
	okIP := checkIP(ip)
	limitCnt := 10
	if len(okIP) < limitCnt {
		fmt.Printf("代理IP数小于%d，从WEB获取新的代理IP:\n", limitCnt)
		webIP, err := GetProxyIP()
		if err != nil {
			fmt.Println("获取WEB代理IP失败:", err)
			os.Exit(-1)
		}
		fmt.Printf("WEB新代理IP数:%d\n", len(webIP))
		okwebIP := checkIP(webIP)
		fmt.Printf("从WEB获取可用代理IP数:%d\n", len(okwebIP))
		okIP = append(okIP, okwebIP...)
		okIP = RemoveRepByMap(okIP)
		fmt.Printf("去重后，代理IP数:%d\n", len(okIP))
		if len(okIP) < limitCnt {
			fmt.Println("代理IP数量太少！")
			writeIP2file(okIP)
			os.Exit(-1)
		}
	}
	writeIP2file(okIP)
	return okIP
}

//GetProxyIP 去Web获取代理IP
func GetProxyIP() ([]string, error) {
	resp, err := http.Get("http://127.0.0.1:5010/get_all")
	if err != nil {
		return nil, err
	}
	defer func() {
		resp.Body.Close() //http得到的resp必须在最后手动关闭
	}()
	respBody, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}
	strProxyIP := string(respBody)
	strTemp := strProxyIP[1 : len(strProxyIP)-1]
	strTemp = strings.Replace(strTemp, "\n", "", -1)
	strTemp = strings.Replace(strTemp, "\"", "", -1)
	strTemp = strings.Replace(strTemp, " ", "", -1)
	resultSlice := strings.Split(strTemp, ",")
	return resultSlice, err
}

func getIPfromFile() []string {
	filePath := `./proxyIP/IP.json`
	dat, err := ioutil.ReadFile(filePath)
	if err != nil {
		panic("ReadFile IPFile fail")
	}
	strProxyIP := string(dat)
	resultSlice := strings.Split(strProxyIP, ",\n")

	return resultSlice
}

func writeIP2file(okIP []string) {
	reStr := strings.Join(okIP, ",\n")
	filePath := `./proxyIP/IP.json`
	ioutil.WriteFile(filePath, []byte(reStr), 0666)
}

func checkIP(ip []string) []string {
	var okIP []string
	var sn sync.Mutex
	var wg sync.WaitGroup
	for _, v := range ip {
		wg.Add(1)
		go func(http string) {
			defer wg.Done()
			_, err := userinfoReq(55645959, http)
			if err != nil {
				fmt.Printf("IP无效:%s\n", http)
				return
			}
			fmt.Printf("IP可用:%s\n", http)
			sn.Lock()
			okIP = append(okIP, http)
			sn.Unlock()
		}(v)
	}
	wg.Wait()
	return okIP
}

//RemoveRepByMap Slice去重
func RemoveRepByMap(slc []string) []string {
	result := []string{}
	tempMap := map[string]byte{} // 存放不重复主键
	for _, e := range slc {
		l := len(tempMap)
		tempMap[e] = 0
		if len(tempMap) != l { // 加入map后，map长度变化，则元素不重复
			result = append(result, e)
		}
	}
	return result
}

//Slice去重
func lIntRemoveByMap(slc []int) []int {
	result := []int{}
	tempMap := map[int]byte{} // 存放不重复主键
	for _, e := range slc {
		l := len(tempMap)
		tempMap[e] = 0
		if len(tempMap) != l { // 加入map后，map长度变化，则元素不重复
			result = append(result, e)
		}
	}
	return result
}
