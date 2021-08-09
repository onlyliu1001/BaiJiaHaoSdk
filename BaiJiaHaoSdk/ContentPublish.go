package BaiJiaHaoSdk

import (
	"bytes"
	"encoding/json"
	"errors"
	"io/ioutil"
	"math"
	"net/http"
	"regexp"
	"unicode"
)

type ContentPublish struct {
	Appid string
	Apptoken string
}

//SetAppid
func (c *ContentPublish)SetAppid(appid string ,apptoken string ){

	c.Appid = appid
	c.Apptoken = apptoken
}

// GetStrLength 返回输入的字符串的字数，汉字和中文标点算 1 个字数，英文和其他字符 2 个算 1 个字数，不足 1 个算 1个
func (c *ContentPublish)getStrLength(str string) float64 {
	var total float64

	reg := regexp.MustCompile("/·|，|。|《|》|‘|’|”|“|；|：|【|】|？|（|）|、/")

	for _, r := range str {
		if unicode.Is(unicode.Scripts["Han"], r) || reg.Match([]byte(string(r))) {
			total = total + 1
		} else {
			total = total + 0.5
		}
	}

	return math.Ceil(total)
}

//ArticlePublish
//发布图文
//发布图文API介绍
//开发者用户可通过本接口向百家号发布图文形式的内容。通过本接口发布的内容数量受帐号可发布内容篇数的统一限制。
//用户上传的图片/视频素材在异步处理的过程中，可能存在处理失败/超时等情况，含有此类素材的文章将不会显示在文章列表，可以使用文章状态API定位问题。
//百家号仅支持审核通过状态的帐号发布内容，其他状态下不可发布。
//ArticlePublish文章发布接口
//title	string	是	文章标题，限定5-40个中英文字符以内
//content	string	是	正文内容，限制20000个中英文字符内，富文本
//origin_url	string	是	原文地址，相同URL的文章会被认为是同一篇文章，禁止提交
//cover_images	json	否	文章封面图片地址url, 0-3张封面图，封面图尺寸不小于218*146，可以为空，没有封面图的内容将会进入草稿
//is_original	int	否	标定是否原创，1 为原创，0 为非原创
//is_split_article	int	否	仅媒体类型百家号可用。标识此篇文章是否单独发布视频内容，1为拆分，0为不拆分，不传此字段时，应用百家号后台默认配置
//video_title	string	否	拆分发布子视频标题，限定 8-40 个中英文字符以内
//video_cover_images	string	否	拆分发布子视频封面图片地址 url, 目前只支持 1 张图片作为封面，封面图尺寸不小于660*370
func (c *ContentPublish)ArticlePublish(title string ,content string,origin_url string ,cover_images []string ,is_original interface{},is_split_article interface{} ,video_title interface{} ,video_cover_images interface{} )(int,*Article_publish_response_data,error)  {

//	文章标题，限定5-40个中英文字符以内
	lentitle:=c.getStrLength(title)
	if(lentitle<5||lentitle>40){
		return 2,nil,errors.New("文章标题，限定5-40个中英文字符以内")
	}
	lencontent:=c.getStrLength(content)
	if(lencontent<1||lencontent>20000){
		return 2,nil,errors.New("正文内容，限制20000个中英文字符内，富文本")
	}

	lenorigin_url:=c.getStrLength(origin_url)
	if lenorigin_url<1 {
		return 2,nil,errors.New("原文地址不能为空")
	}

	var article Article_publish
	article.App_id = c.Appid
	article.App_token = c.Apptoken;
	article.Title = title
	article.Content = content;
	article.Origin_url = origin_url;
	cover_imageslen:=len(cover_images)
	if (cover_imageslen>0) {
		var datacover_images []Article_publish_cover_images;
		datacover_images = make([]Article_publish_cover_images,cover_imageslen)
		for i,j:=range cover_images{
			datacover_images[i].Src=j
		}
		bdatacover_images, _ :=json.Marshal(datacover_images)
		article.Cover_images = string(bdatacover_images);
	}

	if(is_original!=nil){
		article.Is_original = is_original.(int)
	}
	if(is_split_article!=nil) {
		article.Is_split_article = is_split_article.(int)
	}

	if(video_cover_images!=nil) {
	article.Video_cover_images = video_cover_images.(string)
	}
	if(video_title!=nil) {
		article.Video_title = video_title.(string)
	}
	b, _ :=json.Marshal(article)
	//println(string(b))
	//return 2,nil,errors.New(string(b))


	req, err := http.NewRequest("POST", article_publish, bytes.NewReader(b))
	if err != nil {
		return 2,nil,err
	}
	req.Header.Set("Accept", "application/json")
	//这里的http header的设置是必须设置的.
	req.Header.Set("Content-Type", "application/json;charset=utf-8")


	client := http.Client{}
	resp, _err := client.Do(req)
	if _err != nil {
		return 2,nil,err
	}

	respBytes, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return 2,nil,err
	}
	println(string(respBytes))
	var baijiaresponse BaijiaResponse_Public
	err1:=json.Unmarshal(respBytes,&baijiaresponse)
	if err !=nil {
		return 2,nil,err1
	}

	if baijiaresponse.Errno==0 {
		var responedata Article_publish_response_data
		bVal,_:=json.Marshal(baijiaresponse.Data)

		err2:=json.Unmarshal(bVal,&responedata)
		if err2!=nil {
			return 2,nil,err2
		}
		return baijiaresponse.Errno,&responedata,nil
	}else {
		return baijiaresponse.Errno,nil,errors.New(baijiaresponse.Errmsg)
	}

}

//发布图集
//发布图集API介绍
//开发者用户可通过本接口向百家号发布图集形式的内容。通过本接口发布的内容数量受帐号可发布内容篇数的统一限制。
//用户上传的图片/视频素材在异步处理的过程中，可能存在处理失败/超时等情况，含有此类素材的文章将不会显示在文章列表，可以使用文章状态API定位问题。
//百家号仅支持审核通过状态的帐号发布内容，其他状态下不可发布。
//title	string	是	图集标题，限定5-40个中英文字符以内
//photograph	json	是	至少4张图片，desc描述为0-200个汉字；不支持GIF格式；封面图尺寸不小于400*224
//origin_url	string	是	图集资源原地址url，相同URL的文章会被认为是同一篇文章，禁止提交
func (c *ContentPublish)GalleryPublish(title string ,photographsrc []string,photographdesc []string ,origin_url string )(int,*Article_gallery_response_data,error)  {

	//	文章标题，限定5-40个中英文字符以内
	lentitle:=c.getStrLength(title)
	if(lentitle<5||lentitle>40){
		return 2,nil,errors.New("文章标题，限定5-40个中英文字符以内")
	}

	lenphotographsrc := len(photographsrc)
	lenphotographdesc := len(photographdesc)

	if  (lenphotographsrc==0||lenphotographdesc==0||lenphotographsrc!=lenphotographdesc||lenphotographsrc<4||lenphotographdesc<4){
		return 2,nil,errors.New("至少4张图片，desc描述为0-200个汉字；不支持GIF格式；封面图尺寸不小于400*224")
	}

	for _,j:=range photographdesc  {
		lencontent:=c.getStrLength(j)
		if(lencontent<1||lencontent>200){
			return 2,nil,errors.New("desc描述为0-200个汉字")
		}
	}

	var data Article_gallery
	data.Title = title
	data.App_id = c.Appid
	data.App_token = c.Apptoken;
	data.Origin_url = origin_url

	var cover_images []Article_gallery_cover_images
	cover_images = make([]Article_gallery_cover_images,lenphotographsrc)

	for ii,jj:=range photographsrc{
		cover_images[ii].Src=jj
		cover_images[ii].Desc= photographdesc[ii]

	}

	strcover_images, _ :=json.Marshal(cover_images)

	data.Photograph = string(strcover_images);

	b, _ :=json.Marshal(data)

	req, err := http.NewRequest("POST", article_gallery, bytes.NewReader(b))
	if err != nil {
		return 2,nil,err
	}
	req.Header.Set("Accept", "application/json")
	//这里的http header的设置是必须设置的.
	req.Header.Set("Content-Type", "application/json;charset=utf-8")


	client := http.Client{}
	resp, _err := client.Do(req)
	if _err != nil {
		return 2,nil,err
	}

	respBytes, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return 2,nil,err
	}
	println(string(respBytes))
	var baijiaresponse BaijiaResponse_Public
	err1:=json.Unmarshal(respBytes,&baijiaresponse)
	if err !=nil {
		return 2,nil,err1
	}

	if baijiaresponse.Errno==0 {
		var responedata Article_gallery_response_data
		bVal,_:=json.Marshal(baijiaresponse.Data)

		err2:=json.Unmarshal(bVal,&responedata)
		if err2!=nil {
			return 2,nil,err2
		}
		return baijiaresponse.Errno,&responedata,nil
	}else {
		return baijiaresponse.Errno,nil,errors.New(baijiaresponse.Errmsg)
	}

}


//发布视频
//发布视频API介绍
//开发者用户可通过本接口向百家号发布视频形式的内容。通过本接口发布的内容数量受帐号可发布内容篇数的统一限制。
//用户上传的图片/视频素材在异步处理的过程中，可能存在处理失败/超时等情况，含有此类素材的文章将不会显示在文章列表，可以使用文章状态API定位问题。
//百家号仅支持审核通过状态的帐号发布内容，其他状态下不可发布。
//title	string	是	视频标题，限定 8-40 个中英文字符以内
//video_url	string	是	视频原地址，目前支持 mp4 等，不支持 m3u8
//cover_images	string	否	视频封面图片地址 url, 目前只支持 1 张图片作为封面，封面图尺寸不小于660*370
//use_auto_cover	int	否	是否使用自动封面，1为使用自动封面，其余为不使用自动封面
//is_original	int	否	标定是否原创，1 为原创，0 为非原创
//tag	string	否	视频tag，tag之间以半角英文逗号分割，每个tag长度不超过10个字符，最多支持10个tag
func (c *ContentPublish)VideoPublish(title string ,video_url string,cover_images interface{},use_auto_cover interface{} ,is_original interface{},tag interface{})(int,*Video_Publish_response_data,error)  {

	//	文章标题，限定5-40个中英文字符以内
	lentitle:=c.getStrLength(title)
	if(lentitle<5||lentitle>40){
		return 2,nil,errors.New("文章标题，限定5-40个中英文字符以内")
	}

	lenvideo_url:=c.getStrLength(video_url)
	if(lenvideo_url<1){
		return 2,nil,errors.New("视频原地址不能为空")
	}

	var videodata Video_Publish
	videodata.App_id = c.Appid
	videodata.App_token = c.Apptoken;
	videodata.Title = title
	videodata.Video_url = video_url

	if(cover_images!=nil){
		videodata.Cover_images = cover_images.(string)
	}

	if(use_auto_cover!=nil){
		videodata.Use_auto_cover  = use_auto_cover.(int)
	}

	if(is_original!=nil){
		videodata.Is_original = is_original.(int)
	}

	if(tag!=nil){
		videodata.Tag = tag.(string)
	}


	b, _ :=json.Marshal(videodata)

	req, err := http.NewRequest("POST", video_publish, bytes.NewReader(b))
	if err != nil {
		return 2,nil,err
	}
	req.Header.Set("Accept", "application/json")
	//这里的http header的设置是必须设置的.
	req.Header.Set("Content-Type", "application/json;charset=utf-8")


	client := http.Client{}
	resp, _err := client.Do(req)
	if _err != nil {
		return 2,nil,err
	}

	respBytes, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return 2,nil,err
	}
	println(string(respBytes))
	var baijiaresponse BaijiaResponse_Public
	err1:=json.Unmarshal(respBytes,&baijiaresponse)
	if err !=nil {
		return 2,nil,err1
	}

	if baijiaresponse.Errno==0 {
		var responedata Video_Publish_response_data
		bVal,_:=json.Marshal(baijiaresponse.Data)

		err2:=json.Unmarshal(bVal,&responedata)
		if err2!=nil {
			return 2,nil,err2
		}
		return baijiaresponse.Errno,&responedata,nil
	}else {
		return baijiaresponse.Errno,nil,errors.New(baijiaresponse.Errmsg)
	}

}


//发布图片
//图片搜索批量传图服务
//服务介绍
//图片搜索是百度生态第一大图片消费场景。目前机构类的原创作者可通过发布单图API，向图片搜索批量提交您的优质图片资源。您通过百家号提交的图片会收录至百度图片搜索资源库，通过搜索机制选出部分优质图片展现至图片搜索结果页以及您的号主页图片tab下。
//图片质量要求
//百度图片搜索优质资源提交标准
//发布单图API介绍
//使用该接口，您可向百度图片搜索提交您的优质图片资源
//image_list	list	是	批量上传的图片信息，每次最多上传20张图片
//-image_list[]	dict	是	单个图片信息
//--tags	string	否	标签tag（例如体育、明星、动物、壁纸）；多个标签用 # 符号链接
//--edit_time	int	是	图片在来源页的发布时间,秒级时间戳格式
//--id	string	是	资源方图片唯一id，删除图片使用，格式为数字+大小写字母+英文字符
//--title	string	是	图片描述，长度大于8个字符，中英文字符均计算为一个单位长度
//--objurl	string	是	图片文件地址，http开头，不含空白符，不包含中文
//--fromurl	string	是	包含这张图片的来源页地址，http开头，不含空白符，不包含中文
//--copyright	int	否	是否为版权图，默认为0，0-否；非0-是
//--extra	dict	否	额外信息
//---modify	string	否	是否为修改图片，默认为0，非1-否；1-是
func (c *ContentPublish)ImagePublish(imgs []Image_Publish_image_list)(int,error)  {

	var imgdata Image_Publish

	imgslen:=len(imgs)
	if  imgslen==0||imgslen>20{
		return  2,errors.New("每次最多上传20张图片")
	}
	imgdata.App_token= c.Apptoken
	imgdata.App_id = c.Appid

	bVal,_:=json.Marshal(imgs)
	imgdata.Image_list = string(bVal)

	b, _ :=json.Marshal(imgdata)

	req, err := http.NewRequest("POST", image_publish, bytes.NewReader(b))
	if err != nil {
		return 2,err
	}
	req.Header.Set("Accept", "application/json")
	//这里的http header的设置是必须设置的.
	req.Header.Set("Content-Type", "application/json;charset=utf-8")


	client := http.Client{}
	resp, _err := client.Do(req)
	if _err != nil {
		return 2,err
	}

	respBytes, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return 2,err
	}
	println(string(respBytes))
	var baijiaresponse BaijiaResponse_Public
	err1:=json.Unmarshal(respBytes,&baijiaresponse)
	if err !=nil {
		return 2,err1
	}

	if baijiaresponse.Errno==0 {
		return baijiaresponse.Errno,nil
	}else {
		return baijiaresponse.Errno,errors.New(baijiaresponse.Errmsg)
	}

}


//批量删除图片接口
//图片删除API介绍
//如您希望删除之前提交的单图资源，可以通过该接口进行资源的删除操作。
//delete_list	list	是	批量删除的图片信息,单次最多删除10张图片
//-delete_list[]	dict	是	删除的单个图片信息
//--type	string	是	删除图片时，key 参数的类型，值为"objurl"、"id"，代表通过图片url删除或资源方id删除
//--key	string	是	删除图片时，图片的唯一标识，内容为上传时的objurl或id
func (c *ContentPublish)ImageDelete(imgs []Image_delPic_delete_list)(int,error)  {

	imgslen:=len(imgs)
	if  imgslen==0||imgslen>10{
		return  2,errors.New("单次最多删除10张图片")
	}
	var imgdata Image_delPic
	imgdata.App_token= c.Apptoken
	imgdata.App_id = c.Appid

	bVal,_:=json.Marshal(imgs)
	imgdata.Delete_list = string(bVal)

	b, _ :=json.Marshal(imgdata)

	req, err := http.NewRequest("POST", image_delPic, bytes.NewReader(b))
	if err != nil {
		return 2,err
	}
	req.Header.Set("Accept", "application/json")
	//这里的http header的设置是必须设置的.
	req.Header.Set("Content-Type", "application/json;charset=utf-8")


	client := http.Client{}
	resp, _err := client.Do(req)
	if _err != nil {
		return 2,err
	}

	respBytes, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return 2,err
	}
	println(string(respBytes))
	var baijiaresponse BaijiaResponse_Public
	err1:=json.Unmarshal(respBytes,&baijiaresponse)
	if err !=nil {
		return 2,err1
	}

	if baijiaresponse.Errno==0 {
		return baijiaresponse.Errno,nil
	}else {
		return baijiaresponse.Errno,errors.New(baijiaresponse.Errmsg)
	}

}


////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////


//撤回内容
//撤回内容API介绍
//开发者用户可通过本接口对已发布的内容进行撤回，撤回后的文章将不再展示。 被个性化推荐的文章暂不支持进行撤回。
func (c *ContentPublish)Article_Withdraw(article_id string)(int,*Article_Withdraw_response_data,error)  {

	//	文章标题，限定5-40个中英文字符以内
	lentitle:=c.getStrLength(article_id)
	if(lentitle<1){
		return 2,nil,errors.New("文章ID不能为空")
	}


	var videodata Article_Withdraw
	videodata.App_id = c.Appid
	videodata.App_token = c.Apptoken;
	videodata.Article_id = article_id


	b, _ :=json.Marshal(videodata)

	req, err := http.NewRequest("POST", article_withdraw, bytes.NewReader(b))
	if err != nil {
		return 2,nil,err
	}
	req.Header.Set("Accept", "application/json")
	//这里的http header的设置是必须设置的.
	req.Header.Set("Content-Type", "application/json;charset=utf-8")


	client := http.Client{}
	resp, _err := client.Do(req)
	if _err != nil {
		return 2,nil,err
	}

	respBytes, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return 2,nil,err
	}
	println(string(respBytes))
	var baijiaresponse BaijiaResponse_Public
	err1:=json.Unmarshal(respBytes,&baijiaresponse)
	if err !=nil {
		return 2,nil,err1
	}

	if baijiaresponse.Errno==0 {
		var responedata Article_Withdraw_response_data
		bVal,_:=json.Marshal(baijiaresponse.Data)

		err2:=json.Unmarshal(bVal,&responedata)
		if err2!=nil {
			return 2,nil,err2
		}
		return baijiaresponse.Errno,&responedata,nil
	}else {
		return baijiaresponse.Errno,nil,errors.New(baijiaresponse.Errmsg)
	}

}

//修改图文
//修改图文API介绍
//开发者用户可通过本接口对已撤回的图文内容进行修改，可修改标题、内容和封面图，修改成功后的文章将再次发布。对草稿状态的文章调用此API时，将会直接发布。
//用户上传的图片/视频素材在异步处理的过程中，可能存在处理失败/超时等情况，含有此类素材的文章将不会显示在文章列表，可以使用文章状态API定位问题。
//百家号仅支持审核通过状态的帐号发布内容，其他状态下不可发布。
//使用说明
//需要修改的文章状态为“已撤回”
//发布草稿时将占用当天发文次数
//目前只能针对图文内容进行修改
//title	string	是	文章标题，限定8-40个中英文字符以内
//content	string	是	正文内容，限制20000个中英文字符内，富文本
//cover_images	json	是	文章封面图片地址url, 0-3张封面图，可以为空，没有封面图的内容将会进入草稿
//article_id	string	是	需要修改的文章ID
//origin_url	string	是	原文地址
func (c *ContentPublish)ArticleRePublish(title string ,content string,origin_url string ,article_id string,cover_images []string  )(int,*Article_republish_response_data,error)  {

	//	文章标题，限定5-40个中英文字符以内
	lentitle:=c.getStrLength(title)
	if(lentitle<5||lentitle>40){
		return 2,nil,errors.New("文章标题，限定5-40个中英文字符以内")
	}
	lencontent:=c.getStrLength(content)
	if(lencontent<1||lencontent>20000){
		return 2,nil,errors.New("正文内容，限制20000个中英文字符内，富文本")
	}

	lenorigin_url:=c.getStrLength(origin_url)
	if lenorigin_url<1 {
		return 2,nil,errors.New("原文地址不能为空")
	}


	lenarticle_id:=c.getStrLength(article_id)
	if lenarticle_id<1 {
		return 2,nil,errors.New("需要修改的文章ID不能为空")
	}

	var article Article_Republish
	article.App_id = c.Appid
	article.App_token = c.Apptoken;
	article.Title = title
	article.Content = content;
	article.Origin_url = origin_url;
	article.Article_id = article_id
	cover_imageslen:=len(cover_images)
	if (cover_imageslen>0) {
		var datacover_images []Article_publish_cover_images;
		datacover_images = make([]Article_publish_cover_images,cover_imageslen)
		for i,j:=range cover_images{
			datacover_images[i].Src=j
		}
		bdatacover_images, _ :=json.Marshal(datacover_images)
		article.Cover_images = string(bdatacover_images);
	}else{
		article.Cover_images = "[]";
	}

	b, _ :=json.Marshal(article)
	//println(string(b))
	//return 2,nil,errors.New(string(b))
	req, err := http.NewRequest("POST", article_republish, bytes.NewReader(b))
	if err != nil {
		return 2,nil,err
	}
	req.Header.Set("Accept", "application/json")
	//这里的http header的设置是必须设置的.
	req.Header.Set("Content-Type", "application/json;charset=utf-8")


	client := http.Client{}
	resp, _err := client.Do(req)
	if _err != nil {
		return 2,nil,err
	}

	respBytes, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return 2,nil,err
	}
	println(string(respBytes))
	var baijiaresponse BaijiaResponse_Public
	err1:=json.Unmarshal(respBytes,&baijiaresponse)
	if err !=nil {
		return 2,nil,err1
	}

	if baijiaresponse.Errno==0 {
		var responedata Article_republish_response_data
		bVal,_:=json.Marshal(baijiaresponse.Data)

		err2:=json.Unmarshal(bVal,&responedata)
		if err2!=nil {
			return 2,nil,err2
		}
		return baijiaresponse.Errno,&responedata,nil
	}else {
		return baijiaresponse.Errno,nil,errors.New(baijiaresponse.Errmsg)
	}

}
