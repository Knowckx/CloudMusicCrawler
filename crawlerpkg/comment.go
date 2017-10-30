package crawlerpkg

import (
	"fmt"
)

type commentMusic struct {
	HotComments []comments
	Comments    []comments
	Total       uint64
}
type comments struct {
	User       user
	LikedCount int
	CommentID  int
	Time       int64
	Content    string
	BeReplied  []beReplied
}

//Suser 用户结构
type user struct {
	Nickname string
	UserID   int
}

type beReplied struct {
	User    user
	Content string
}

//RequestParam post request body
type RequestParam struct {
	Offset    uint32 `json:"offset"`
	Rid       string `json:"rid"`
	Limit     int    `json:"limit"`
	CsrfToken string `json:"csrf_token"`
}

func getPageThenBegin(ch chan uint32, proxyIP string) {
	page := <-ch //取一个，作为本次需要操作的页数

	defer func() { //线程错误恢复
		if err := recover(); err != nil {
		}
		wgRquest.Done()
	}()
	getCommentsCount := uint64((page + 1) * 100) //页数X100就是本次的标识位置

	total := dataCmt.Total
	//尾页判断
	if total != 0 && getCommentsCount > total && getCommentsCount-total > 100 {
		// fmt.Println("uint64(page*100) > total：", page, total)
		isContinue = false //都是改false 不锁
		ch <- page + 1
		return
	}
	ch <- page + 1 //下一个线程可执行
	getComments(page, proxyIP, false)
}

//获取评论，并进行处理
func getComments(page uint32, proxyIP string, isDealErr bool) {
	fmt.Printf("本gro对应的page %d页\n", page)
	if isDealErr {
		defer wgDealErros.Done()
	}
	comments, err := commentReq(proxyIP, page) //指定页评论数据
	if err != nil {
		fmt.Printf("第%d页没有获取到\n", page)
		rw.Lock()
		errosPages = append(errosPages, page) //错误页记录
		rw.Unlock()
		return
	}
	rw.Lock()
	dataCmt.Comments = append(dataCmt.Comments, comments.Comments...)
	rw.Unlock()
}
