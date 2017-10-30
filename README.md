# NeteaseMusicSpider
概述：
---------
go的一个练手项目，爬取网易云音乐指定音乐下所有的评论信息。然后根据评论用户列表爬取用户信息。

最后结果会以json格式写入NeteaseMusicSpider\ResultData目录下。

输入参数处理还没做，要爬啥去begin.go里修改songID变量就好了（跑~

依赖：
---------
[goquery](https://github.com/PuerkitoBio/goquery) 解析HTML，定位元素

代理IP池：
---------
这个问题也困扰我好久，自己买什么的最方便了……
 
我最后用的是[proxy_pool](https://github.com/jhao104/proxy_pool)。你可以自己在本地搞一个

目前代理IP的获取API地址写死了"http://127.0.0.1:5010/get_all" 以后有时间我弄一套去服务器上（跑

