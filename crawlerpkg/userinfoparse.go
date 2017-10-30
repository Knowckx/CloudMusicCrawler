package crawlerpkg

import (
	"bytes"
	"fmt"
	"strconv"
	"time"

	"github.com/PuerkitoBio/goquery"
)

type userinfos struct {
	Listuser []*UserInfo
}

//UserInfo 用户信息结构
type UserInfo struct {
	ID       int
	Nickname string
	Gender   Gender
	age      string
	Area     string
	Desc     string

	Sitelink string
}

// Gender 性别
type Gender int

const (
	//Unknown 不男不女
	Unknown Gender = iota
	// Male --> 1
	Male
	// Female --> 2
	Female
)

// userinfoparse 解析用户信息页面
func userinfoparse(page []byte, UID int) *UserInfo {
	// fmt.Printf("page:%s\n", string(page))
	rpage := bytes.NewReader(page)
	doc, err := goquery.NewDocumentFromReader(rpage)
	if err != nil {
		fmt.Println("userpage parse fail", UID)
		panic(err)
	}
	userinfo := new(UserInfo)
	userinfo.ID = UID

	findStr := `#j-name-wrap > span.tit.f-ff2.s-fc0.f-thide` //定位用户名
	nickname := doc.Find(findStr).Text()
	userinfo.Nickname = nickname

	findStr = `#j-name-wrap > i` //定位性别
	selegender := doc.Find(findStr)
	if selegender.HasClass("u-icn-01") {
		userinfo.Gender = 1
	} else if selegender.HasClass("u-icn-02") {
		userinfo.Gender = 2
	} else {
		userinfo.Gender = Unknown
	}

	//HTML <span class="sep" id="age" data-age="920044800000">年龄：
	findStr = `span#age`                             //定位年龄
	strAge, ok := doc.Find(findStr).Attr("data-age") //Unix 毫秒时间
	if ok {
		intAge, _ := strconv.ParseInt(strAge, 10, 64)
		unixTime := time.Unix(intAge/1000, 0)
		strAge = unixTime.Format("2006-01-02")
		userinfo.age = strAge
	} else {
		// fmt.Println("no data-age")
		userinfo.age = ""
	}

	findStr = `#head-box > dd > div.inf.s-fc3 > span:contains("地区")` //定位地区
	areaStr := doc.Find(findStr).Text()
	if areaStr != "" {
		areaStr = string([]rune(areaStr)[5:])
	}
	userinfo.Area = areaStr

	findStr = `#head-box > dd > div.inf.s-fc3.f-brk` //定位个人介绍
	descStr := doc.Find(findStr).Text()
	if descStr != "" {
		descStr = string([]rune(descStr)[5:])
	}
	userinfo.Desc = descStr

	findStr = `#head-box > dd > div.inf.s-fc3.f-cb > ul > li > a` //定位微博
	linkStr, _ := doc.Find(findStr).Attr("href")
	userinfo.Sitelink = linkStr
	// fmt.Printf("userinfo %+v", userinfo)
	return userinfo
}
