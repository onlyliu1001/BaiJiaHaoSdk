package BaiJiaHaoSdk

import (
	"bytes"
	"encoding/json"
	"errors"
	"io/ioutil"
	"net/http"
	"strings"
)

//获取文章列表
//获取文章列表API介绍
//开发者用户可通过本接口来获取图文、图集、视频类型的文章列表，支持用户按照百家号作者帐号ID、文章类型、时间批量查询。单次最多可查询20篇文章数据。
//用户上传的图片/视频素材在异步处理的过程中，可能存在处理失败/超时等情况，含有此类素材的文章将不会显示在文章列表，可以使用文章状态API定位问题。
//start_time	string	否	支持按照年月日格式（2019-06-01）进行查询，仅支持查询到日维度的数据
//end_time	string	否	支持按照年月日格式（2019-07-01）进行查询，仅支持查询到日维度的数据
//page_no	int	否	查询页码，不传默认为1
//page_size	int	否	查询条数，不能超过20，不传默认为20
//article_type	string	否	文章类型，news-图文、gallery-图集、video-视频，不传默认查询所有支持的文章类型
//collection	string	否	文章状态集，不传默认查询所有支持的文章状态集 draft-草稿、publish-已发布、pre_publish-待发布、withdraw-已撤回、rejected-未通过










//获取文章状态
//获取文章状态API介绍
//开发者用户可通过本接口批量查询文章状态。单次最多可查询20篇文章的状态，且不可使用文章ID与NID混合查询。
//用户上传的图片/视频素材在异步处理的过程中，可能存在处理失败/超时等情况，含有此类素材的文章将不会显示在文章列表，可以使用文章状态API定位问题。
//article_id	string	是	需要查询的文章ID,多个使用英文逗号分隔
func (c *ContentPublish)Query_Status(article_id string)(int,[]Query_Article_Status,error)  {

	lenarticle_id:=c.getStrLength(article_id)
	if(lenarticle_id<1){
		return 2,nil,errors.New("查询的文章ID不能为空")
	}

	var article Query_Article
	article.App_id = c.Appid
	article.App_token = c.Apptoken;
	article.Article_id = article_id

	b, _ :=json.Marshal(article)
	//println(string(b))
	//return 2,nil,errors.New(string(b))
	req, err := http.NewRequest("POST", query_articlestatus, bytes.NewReader(b))
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
	if err1 !=nil {
		return 2,nil,err1
	}
	var resplist []Query_Article_Status
	if baijiaresponse.Errno==0 {
		article_id_list:=strings.Split(article_id,",")
		var require map[string]interface{}//请求参数
		b1,_:=json.Marshal(baijiaresponse.Data)
		if err := json.Unmarshal(b1, &require); err == nil {
			for _,j:=range article_id_list{
				if _, ok := require[j]; ok {
					var temp Query_Article_Status
					datalist:=require[j].(map[string]interface {});
					temp.Article_id = j
					if _, ok := datalist["status"]; ok{
						temp.Status = datalist["status"].(string)
					}
					if _, ok := datalist["url"]; ok{
						temp.Url = datalist["url"].(string)
					}
					if _, ok := datalist["msg"]; ok{
						temp.Msg = datalist["msg"].(string)
					}
					resplist=append(resplist, temp)
				}
			}
		}
		return baijiaresponse.Errno,resplist,nil
	}else {
		return baijiaresponse.Errno,nil,errors.New(baijiaresponse.Errmsg)
	}




}