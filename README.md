# BaiJiaHaoSdk
目前写了这些功能，后续会慢慢增加。
一 内容发布

+发布图文

+发布图集

+发布视频

+发布图片


二 内容管理

+ 撤回内容

+ 修改图文

三 文章查询

+ 获取文章状态


四 评论管理

五 数据查询

六 用户管理


使用示例：

var baijiahao BaiJiaHaoSdk.ContentPublish

//设置APPID和TOKEN

baijiahao.SetAppid("appid "," apptoken")

//发布文章 返回参数看函数吧。

baijiahao.ArticlePublish("正文标题", `<div class="article_show_body"><p>正文内容</p></div>`, "原网址",nil,nil,nil,nil,nil)

//查询状态

baijiahao.Query_Status("文章ID")
