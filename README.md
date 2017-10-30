# NeteaseMusicSpider
概述：
---------
golang的一个练手项目，爬取网易云音乐指定音乐下所有的评论信息。然后根据评论用户列表爬取用户信息。

最后结果会以json格式写入NeteaseMusicSpider\ResultData目录下。目前爬取了[剧中歌](http://music.163.com/#/song?id=31889414)作为示例

输入参数处理还没做，要爬啥去begin.go里修改songID变量就好了（跑~

依赖：
---------
* 解析HTML

解析HTML，定位元素 [goquery](https://github.com/PuerkitoBio/goquery) 

* 代理IP池：
这个问题也困扰我好久，没有代理IP多线程的请求分分钟封IP

当然自己买IP池什么的最方便了……但是好贵啊！
 
我最后用的是[proxy_pool](https://github.com/jhao104/proxy_pool)。你需要在本地搞一个

目前代理IP的获取API地址写死了"http://127.0.0.1:5010/get_all" 

以后有时间我可以弄一套去服务器上，嗯……以后（跑

