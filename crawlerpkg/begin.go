package crawlerpkg

import (
	"encoding/json"

	"fmt"
	"os"
	"runtime"
	"sync"
	"time"
)

var songID string

//结果数据
var dataCmt *commentMusic
var dataUserInfo *userinfos

//UserID 查找用户信息时的ID
var UserID string

var once sync.Once

var isContinue = true

//WgRquest 请求获取评论main同步
var rw sync.Mutex
var wgRquest = sync.WaitGroup{}
var wgDealComment = sync.WaitGroup{}

//错误页 重新获取
var errosPages []uint32
var wgDealErros = sync.WaitGroup{}

var errosUser []int
var wgDealUserErros = sync.WaitGroup{}

//可用的代理IP
var okIPs []string

//临时调试用
func test() {
	// ---------------------------------------测试Userinfo的获取
	// UserID = "148332"
	// proxyIP := "127.0.0.1"
	proxyIP := "61.135.217.7:80"
	_ = proxyIP
	_, err := userinfoReq(55645959, proxyIP)
	if err != nil {
		fmt.Printf("err:%v", err)
	}
	os.Exit(1)
	// ---------------------------------------
	// ---------------------------------------测试Comment的获取
	// songID = "561435" //15~Spring Flag~  19条评论
	// proxyIP := "127.0.0.1"
	// cmt, err := commentReq(proxyIP, 0)
	// if err != nil {
	// 	fmt.Printf("err:%v", err)
	// }
	// fmt.Printf("comment:%v", cmt)
	// return
	// ---------------------------------------
	// isBegin := stdinUserNameAndID()  //参数处理，以后搞
	// if isBegin == false {
	// 	os.Exit(-1)
	// }
}

//命令参数处理  待施工
func stdinUserNameAndID() bool {
	defer func() {
		if err := recover(); err != nil {
			fmt.Println(err)
		}
	}()

	// if len(os.Args) != 5 {
	// 	panic(errors.New("命令有误"))
	// }

	// if os.Args[1] == "-name" {
	// 	findUseName = os.Args[2] //拿到用户名
	// } else { //没输入用户名的时候  第二个参数就是歌曲ID
	// 	secondParam := os.Args[2]
	// 	_, err := strconv.ParseUint(secondParam, 10, 64)
	// 	if err != nil {
	// 		panic(errors.New(fmt.Sprint("提供的歌曲ID格式错误!!!"+"--", err)))
	// 	}
	// 	songID = secondParam
	// }

	// if os.Args[3] == "-ID" {
	// 	secondParam := os.Args[4]
	// 	_, err := strconv.ParseUint(secondParam, 10, 64)
	// 	if err != nil {
	// 		panic(errors.New(fmt.Sprint("提供的歌曲ID格式错误!!!"+"--", err)))
	// 	}
	// 	songID = secondParam
	// } else {
	// 	findUseName = os.Args[4]
	// }
	// if songID == "" || findUseName == "" {
	// 	panic(errors.New("歌曲ID获取或者用户名不能为空"))
	// }
	// fmt.Println("您输入的查找的用户名为:", findUseName, "查找的歌曲ID:", songID)
	// fmt.Println("请确认(y/n)")
	// var correct string
	// fmt.Scanln(&correct) //获得输入值
	// if correct == "y" || correct == "Y" {
	// 	return true
	// }
	return false
}

//Begin 程序入口
func Begin() {
	// test()
	runtime.GOMAXPROCS(runtime.NumCPU())

	songID = "31889414"    //以后把这个参数放配置
	okIPs = getOkProxyIP() //代理IP
	commentOnce()          //预访问
	getCmtBegin()

	for len(errosPages) > 0 {
		dealErrorPage()
	}
	fmt.Println("评论数据查找完毕！！！")

	getUserinfoBegin()

	timer := time.NewTimer(time.Duration(5) * time.Minute) //用户数较多，限时。
	go func() {
		<-timer.C
		fmt.Println("处理错误超时")
		os.Exit(-1)
	}()
	for len(errosUser) > 0 {
		dealErrorUser()
	}
	err := writeData()
	if err != nil {
		print(err)
		panic("信息写入出错，退出！")
	}
}

func commentOnce() {
	fmt.Println("------查找的歌曲ID为：", songID)
	dataCmt = new(commentMusic) //第一次
	dataUserInfo = new(userinfos)
	cmt, err := commentReq("127.0.0.1", 1)
	if err != nil {
		panic(err)
	}
	dataCmt.Total = cmt.Total
	if dataCmt.Total == 0 {
		fmt.Println("总评论数为0")
		os.Exit(-1)
	}
}
func getCmtBegin() {
	ipCount := len(okIPs)
	ipIndex := 0
	//用这个ch来保存page。在各gro中进行传递。
	ch := make(chan uint32, ipCount)
	ch <- uint32(0)

	//goroutine数量
	count := 0
	allCount := ipCount * 10

	fmt.Printf("开始查找「%s」下的所有评论用户信息:\n", songID)
	fmt.Println("总评论数:", dataCmt.Total)
	for isContinue { //gro分配IP
		if ipIndex > ipCount-1 {
			ipIndex = 0
		}
		if count == allCount { //最大并发数限制
			wgRquest.Wait()
			count = 0
		}
		wgRquest.Add(1)
		go getPageThenBegin(ch, okIPs[ipIndex])
		count++
		ipIndex++
	}
	fmt.Printf("查找评论信息 Wait\n")
	wgRquest.Wait()
	fmt.Printf("查找评论信息 第一次完毕\n")

}

//处理出错页的评论
func dealErrorPage() {
	ipCount := len(okIPs)
	ipIndex := 0
	copySlice := errosPages[:]
	errosPages = []uint32{}
	fmt.Printf("评论 错误页数量%d\n", len(copySlice))

	for _, v := range copySlice {
		if ipIndex >= ipCount {
			ipIndex = 0
		}
		wgDealErros.Add(1)
		fmt.Printf("正在重新获取第%d页\n", v)
		go getComments(v, okIPs[ipIndex], true)
		ipIndex++
	}
	wgDealErros.Wait()
}
func getUserinfoBegin() {
	lUser := []int{}
	for _, v := range dataCmt.Comments {
		// print(v.Nickname)
		lUser = append(lUser, v.User.UserID)
	}
	lUser = lIntRemoveByMap(lUser)
	fmt.Printf("用户列表数 %d\n", len(lUser))

	ipCount := len(okIPs)
	ipIndex := 0

	for _, v := range lUser {
		if ipIndex >= ipCount {
			ipIndex = 0
		}
		wgRquest.Add(1)
		// fmt.Printf("正在获取用户数据 %d\n", v)
		go getuserinfo(v, okIPs[ipIndex])
		ipIndex++
	}
	wgRquest.Wait()
}

//处理出错用户信息
func dealErrorUser() {
	ipCount := len(okIPs)
	ipIndex := 0
	copySlice := errosUser[:]
	errosUser = []int{}

	fmt.Printf("用户信息 错误数量%d\n", len(copySlice))

	for _, v := range copySlice {
		if ipIndex >= ipCount {
			ipIndex = 0
		}
		wgRquest.Add(1)
		fmt.Printf("正在重新获取用户信息 %d\n", v)
		go getuserinfo(v, okIPs[ipIndex])
		ipIndex++
	}
	wgRquest.Wait()
}
func writeData() error {
	fmt.Printf("获得评论数 %d:\n", len(dataCmt.Comments))
	fmt.Printf("获得用户数据 %d:\n", len(dataUserInfo.Listuser))
	dataCmtJSON, err := json.Marshal(dataCmt)
	if err != nil {
		return err
	}
	dataUsersJSON, err := json.Marshal(dataUserInfo)
	if err != nil {
		return err
	}
	filePath := fmt.Sprintf(`./ResultData/comment_%s.json`, songID)
	filePath2 := fmt.Sprintf(`./ResultData/userinfo_%s.json`, songID)
	newFile, err := os.Create(filePath)
	if err != nil {
		fmt.Print(err)
		panic("create File error!")
	}
	_, err = newFile.Write(dataCmtJSON)
	if err != nil {
		return err
	}
	newFile, err = os.Create(filePath2)
	if err != nil {
		fmt.Print(err)
		panic("write File error")
	}
	_, err = newFile.Write(dataUsersJSON)
	if err != nil {
		return err
	}
	fmt.Printf("数据写入成功\n%s\n%s\n", filePath, filePath2)
	return nil
}
