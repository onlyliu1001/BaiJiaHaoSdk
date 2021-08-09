package BaiJiaHaoSdk

import (
	"crypto/sha1"
	"fmt"
	"io"
)

const (
SUCCESS = 0	//成功
ERR_2 = 2 //参数错误
ERR_4 = 4
ERR_5 = 5    // 4或5	数据库错误
ERR_10002 = 10002	//无权限调用相关接口
ERR_10003 =10003	//数据库错误
ERR_10005 =10005	//百家号内部服务异常
ERR_10009 = 10009	//非作者文章
ERR_27012 = 27012	//文章修改次数达到上限
ERR_301001 = 301001	//无有效数据
ERR_301002 = 301002	//今日上传图片次数耗尽
ERR_301003 = 301003	//图片重复入库
ERR_301004 = 301004	//要删除的图片不存在
ERR_300016 = 300016	//目前只有图文类型支持修改
ERR_60001001 = 60001001	//授权校验失败
ERR_60001003 = 60001003	//撤回或修改时文章状态错误
ERR_60000005 = 60000005	//帐号未审核通过或在禁言期
ERR_60000006 = 60000006	//百家号内部服务异常
ERR_60000009 = 60000009	//帐号不存在
ERR_60000020 = 60000020	//百家号内部服务异常
ERR_60001008 = 60001008	//当天发文篇数校验失败
ERR_60001009 = 60001009	//发文篇数用尽
ERR_80000107 = 80000107	//文章不存在
ERR_800000201 = 800000201//	文章状态异常
)

const(
	article_publish =  `https://baijiahao.baidu.com/builderinner/open/resource/article/publish`
	article_gallery =  `https://baijiahao.baidu.com/builderinner/open/resource/article/gallery`
	video_publish   =  `https://baijiahao.baidu.com/builderinner/open/resource/video/publish`
	image_publish  =  `https://baijiahao.baidu.com/builderinner/open/resource/image/pushPic`
	image_delPic = `https://baijiahao.baidu.com/builderinner/open/resource/image/delPic`

	article_withdraw = `https://baijiahao.baidu.com/builderinner/open/resource/article/withdraw`
	article_republish = `https://baijiahao.baidu.com/builderinner/open/resource/article/republish`

	query_articleListall = `https://baijiahao.baidu.com/builderinner/open/resource/query/articleListall`

	query_articlestatus = `https://baijiahao.baidu.com/builderinner/open/resource/query/status`


)

type Lib_Util_AesEncrypt struct {
	Signature	string	`json:"signature"`//是	对传参的加密签名值，用来确认请求是来自开放平台的必要校验
	Timestamp	string	`json:"timestamp"`//是	时间戳
	Nonce	string	`json:"nonce"`//是	随机数
	Encrypt	string	`json:"encrypt"`//是	加密内容（消息内容，为传送消息的xml或者json字符串）
}

type BaijiaResponse_Public struct {
	Errno  int    `json:"errno"`            //错误码
	Errmsg string `json:"errmsg"`            //错误信息
	Data   interface{}  `json:"data,string"`            //返回数据
}

/////////////////////////////////////////发布图文//////////////////////////////
type Article_publish_cover_images struct{
	Src string 	`json:"src"`
}

type Article_publish struct {
	App_id	string	`json:"app_id"`//是	作者帐号ID
	App_token	string	`json:"app_token"`//是	授权密钥
	Title	string	`json:"title"`//是	文章标题，限定5-40个中英文字符以内
	Content	string	`json:"content"`//是	正文内容，限制20000个中英文字符内，富文本
	Origin_url	string	`json:"origin_url"`//是	原文地址，相同URL的文章会被认为是同一篇文章，禁止提交
	Cover_images	string `json:"cover_images,omitempty"`//json	否	文章封面图片地址url, 0-3张封面图，封面图尺寸不小于218*146，可以为空，没有封面图的内容将会进入草稿
	Is_original	int	`json:"is_original,omitempty"`//否	标定是否原创，1 为原创，0 为非原创
	Is_split_article	int	 `json:"is_split_article,omitempty"`//否	仅媒体类型百家号可用。标识此篇文章是否单独发布视频内容，1为拆分，0为不拆分，不传此字段时，应用百家号后台默认配置
	Video_title	string	`json:"video_title,omitempty"`//否	拆分发布子视频标题，限定 8-40 个中英文字符以内
	Video_cover_images	string	 `json:"video_cover_images,omitempty"`//否	拆分发布子视频封面图片地址 url, 目前只支持 1 张图片作为封面，封面图尺寸不小于660*370
}

type Article_publish_response_data struct{
	Article_id    string    `json:"article_id"` //文章id
	Nid    string    `json:"nid"`        //文章nid
	Split_article_id  string  `json:"split_article_id"`     //string	拆分发布子视频文章id
	Split_nid    string    `json:"split_nid"`  //拆分发布子视频文章nid
}


/////////////////////////////////////////发布图集//////////////////////////////
type Article_gallery_cover_images struct{
	Src string 	`json:"src"`
	Desc string 	`json:"desc"`
}

type Article_gallery struct{
	App_id	string	`json:"app_id"` //是	作者帐号ID
	App_token	string	`json:"app_token"` //是	授权密钥
	Title	string	`json:"title"` //是	图集标题，限定5-40个中英文字符以内
	Photograph	string	`json:"photograph"` //是	至少4张图片，desc描述为0-200个汉字；不支持GIF格式；封面图尺寸不小于400*224
	Origin_url	string	`json:"origin_url"` //是	图集资源原地址url，相同URL的文章会被认为是同一篇文章，禁止提交
}

type Article_gallery_response_data struct{
	Article_id    string    `json:"article_id"` //文章id
	Nid    string    `json:"nid"`        //文章nid
}



/////////////////////////////////////////发布视频//////////////////////////////
type Video_Publish struct{
	App_id	string	`json:"app_id"` //是	作者帐号ID
	App_token	string	`json:"app_token"` //是	授权密钥
	Title	string	`json:"title"` //是	视频标题，限定 8-40 个中英文字符以内
	Video_url	string	`json:"video_url"` //是	视频原地址，目前支持 mp4 等，不支持 m3u8
	Cover_images	string	`json:"cover_images,omitempty"` //否	视频封面图片地址 url, 目前只支持 1 张图片作为封面，封面图尺寸不小于660*370
	Use_auto_cover	int	`json:"use_auto_cover,omitempty"` //否	是否使用自动封面，1为使用自动封面，其余为不使用自动封面
	Is_original	int	`json:"is_original,omitempty"` //否	标定是否原创，1 为原创，0 为非原创
	Tag	string	`json:"tag,omitempty"` //否	视频tag，tag之间以半角英文逗号分割，每个tag长度不超过10个字符，最多支持10个tag
}

type Video_Publish_response_data struct{
	Article_id    string    `json:"article_id"` //文章id
	Nid    string    `json:"nid"`        //文章nid
}

/////////////////////////////////////////发布图片//////////////////////////////
type Image_Publish_image_list_extra struct {
	Modify	string	`json:"modify"` //否	是否为修改图片，默认为0，非1-否；1-是
}

type Image_Publish_image_list struct {
	Tags	string	`json:"tags"` //否	标签tag（例如体育、明星、动物、壁纸）；多个标签用 # 符号链接
	Edit_time	int	`json:"edit_time"` //是	图片在来源页的发布时间,秒级时间戳格式
	Id	string	`json:"id"` //是	资源方图片唯一id，删除图片使用，格式为数字+大小写字母+英文字符
	Title	string	`json:"title"` //是	图片描述，长度大于8个字符，中英文字符均计算为一个单位长度
	Objurl	string	`json:"objurl"` //是	图片文件地址，http开头，不含空白符，不包含中文
	Fromurl	string	`json:"fromurl"` //是	包含这张图片的来源页地址，http开头，不含空白符，不包含中文
	Copyright	int	`json:"copyright,omitempty"` //否	是否为版权图，默认为0，0-否；非0-是
	Extra	Image_Publish_image_list_extra	`json:"extra,omitempty"` //否	额外信息
}

type Image_Publish struct{
	App_id	string	`json:"app_id"` //是	作者帐号ID
	App_token	string	`json:"app_token"` //是	授权密钥
	Image_list string 	`json:"image_list"`
}

////////////////////////////////////////批量删除图片接口//////////////////////////////
type Image_delPic_delete_list struct {
	Type string	`json:"type"` //是	删除图片时，key 参数的类型，值为"objurl"、"id"，代表通过图片url删除或资源方id删除
	Key	string	`json:"key"` //是	删除图片时，图片的唯一标识，内容为上传时的objurl或id
}

type Image_delPic struct{
	App_id	string	`json:"app_id"` //是	作者帐号ID
	App_token	string	`json:"app_token"` //是	授权密钥
	Delete_list string 	`json:"delete_list"`
}

type Post_Content struct{
	Signature	string	`json:"signature"` //是	对传参的加密签名值，用来确认请求是来自开放平台的必要校验
	Timestamp	string	`json:"timestamp"` //是	时间戳
	Nonce	string	`json:"nonce"` //是	随机数
	Encrypt	string	`json:"encrypt"` //是	加密内容（消息内容，为传送消息的xml或者json字符串）
}

////////////////////////////////////////////////////内容管理////////////////////////////////////////////////////////////////////
type Article_Withdraw struct {
	App_id	string	`json:"app_id"` //是	作者帐号ID
	App_token	string	`json:"app_token"` //是	授权密钥
	Article_id string 	`json:"article_id"` //需要撤回的文章ID
}

type Article_Withdraw_response_data struct{
	Article_id    string    `json:"article_id"` //文章id
	 }


type Article_Republish struct {
	App_id	string	`json:"app_id"`//是	作者帐号ID
	App_token	string	`json:"app_token"`//是	授权密钥
	Title	string	`json:"title"`//是	文章标题，限定5-40个中英文字符以内
	Content	string	`json:"content"`//是	正文内容，限制20000个中英文字符内，富文本
	Origin_url	string	`json:"origin_url"`//是	原文地址，相同URL的文章会被认为是同一篇文章，禁止提交
	Cover_images	string `json:"cover_images"`//json	否	文章封面图片地址url, 0-3张封面图，封面图尺寸不小于218*146，可以为空，没有封面图的内容将会进入草稿
	Article_id	string	`json:"article_id"`//是	需要修改的文章ID
}


type Article_republish_response_data struct{
	Article_id    string    `json:"article_id"` //文章id
}

//////////////////////////////////////////////////文章查询//////////////////////////////////////////////////////////////////////


type Query_articleListall struct {
	app_token	string	`json:"app_id"`//是	授权密钥
	app_id	string	`json:"app_id"`//是	作者帐号ID
	start_time	string	`json:"start_time,omitempty"`//否	支持按照年月日格式（2019-06-01）进行查询，仅支持查询到日维度的数据
	end_time	string	`json:"end_time,omitempty"`//否	支持按照年月日格式（2019-07-01）进行查询，仅支持查询到日维度的数据
	page_no	int	`json:"page_no,omitempty"`//否	查询页码，不传默认为1
	page_size	int	`json:"page_size,omitempty"`//否	查询条数，不能超过20，不传默认为20
	article_type	string	`json:"article_type,omitempty"`//否	文章类型，news-图文、gallery-图集、video-视频，不传默认查询所有支持的文章类型
	collection	string	`json:"collection,omitempty"`//否	文章状态集，不传默认查询所有支持的文章状态集 draft-草稿、publish-已发布、pre_publish-待发布、withdraw-已撤回、rejected-未通过
}

type Query_Article struct {
	App_id	string	`json:"app_id"`//是	作者帐号ID
	App_token	string	`json:"app_token"`//是	授权密钥
	Article_id	string	`json:"article_id"`//是
}


type Query_Article_Status struct {
	Article_id	string	`json:"article_id"`//是
	Status	string	`json:"status"`//是	授权密钥
	Url	string	`json:"url,omitempty"`//是	作者帐号ID
	Msg	string	`json:"msg,omitempty"`//否	支持按照年月日格式（2019-06-01）进行查询，仅支持查询到日维度的数据
}




func Sha1Msg(data string) string {
	t := sha1.New();
	io.WriteString(t,data);
	return fmt.Sprintf("%x",t.Sum(nil));
}