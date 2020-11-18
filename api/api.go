package api

import (
	"fmt"
	"github.com/lxn/walk"
	"github.com/mitchellh/go-homedir"
	"github.com/tidwall/gjson"
	"io/ioutil"
	"kingwrcy/dnspod-client/model"
	"net/http"
	"net/url"
	"os"
	"path"
)

var dnsPodToken = ""
var dnspodFileName = ".dnspod"

func ShowErrMsg(msg string) {
	walk.MsgBox(nil, "错误", msg, walk.MsgBoxIconError)
}

func ShowSuccMsg(msg string) {
	walk.MsgBox(nil, "操作成功", msg, walk.MsgBoxOK)
}

func GetLoginToken() string {
	if dnsPodToken != "" {
		return dnsPodToken
	}
	dir, err := homedir.Dir()
	if err != nil {
		ShowErrMsg("读取用户目录异常")
		return ""
	}
	f, err := os.Open(path.Join(dir, dnspodFileName))
	if err != nil {
		ShowErrMsg("login_token没有配置,请在用户目录下配置`.dnspod`文件,内容格式:`ID,Token`,不要有换行符!")
		return ""
	}
	defer f.Close()

	bb, err := ioutil.ReadAll(f)
	if err != nil {
		ShowErrMsg("读取token异常")
		return ""
	}
	dnsPodToken = string(bb)

	return dnsPodToken
}

func GetDomainList() []model.Domain {
	val := url.Values{}
	val.Add("login_token", GetLoginToken())
	val.Add("format", "json")
	resp, err := http.PostForm("https://dnsapi.cn/Domain.List", val)
	if err != nil {
		walk.MsgBox(nil, "错误", fmt.Sprintf("请求域名列表报错:%s", err.Error()), walk.MsgBoxIconError)
		return nil
	}
	bb, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		walk.MsgBox(nil, "错误", fmt.Sprintf("解析域名列表报错:%s", err.Error()), walk.MsgBoxIconError)
		return nil
	}
	body := string(bb)
	domains := gjson.Get(body, "domains").Array()
	var result []model.Domain
	for _, domain := range domains {
		result = append(result, model.Domain{
			Name:    domain.Get("name").String(),
			ID:      domain.Get("id").Int(),
			Records: domain.Get("records").String(),
		})
	}
	return result
}

func GetRecordList(domainID int64) []model.Record {
	val := url.Values{}
	val.Add("login_token", GetLoginToken())
	val.Add("format", "json")
	val.Add("domain_id", fmt.Sprintf("%d", domainID))
	resp, err := http.PostForm("https://dnsapi.cn/Record.List", val)
	if err != nil {
		walk.MsgBox(nil, "错误", fmt.Sprintf("请求记录列表报错:%s", err.Error()), walk.MsgBoxIconError)
		return nil
	}
	bb, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		walk.MsgBox(nil, "错误", fmt.Sprintf("解析记录列表报错:%s", err.Error()), walk.MsgBoxIconError)
		return nil
	}
	body := string(bb)
	domains := gjson.Get(body, "records").Array()
	var result []model.Record
	for _, domain := range domains {
		result = append(result, model.Record{
			Name:  domain.Get("name").String(),
			Type:  domain.Get("type").String(),
			Value: domain.Get("value").String(),
			ID:    domain.Get("id").Int(),
		})
	}
	return result
}

func SaveRecord(domainId int64, record model.Record) bool {
	val := url.Values{}
	val.Add("login_token", GetLoginToken())
	val.Add("format", "json")
	val.Add("sub_domain", record.Name)
	val.Add("record_type", record.Type)
	val.Add("value", record.Value)
	val.Add("record_line", "默认")
	val.Add("domain_id", fmt.Sprintf("%d", domainId))
	val.Add("mx", "20")
	fmt.Printf("value:%+v\n", val)
	resp, err := http.PostForm("https://dnsapi.cn/Record.Create", val)
	if err != nil {
		ShowErrMsg(fmt.Sprintf("保存记录报错:%s", err.Error()))
		return false
	}
	bb, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		ShowErrMsg(fmt.Sprintf("保存记录报错:%s", err.Error()))
		return false
	}
	body := string(bb)
	code := gjson.Get(body, "status.code").Int()
	if code == 1 {
		ShowSuccMsg("添加成功")
		return true
	} else {
		ShowErrMsg("添加失败:" + gjson.Get(body, "status.message").String())
		return false
	}
}

func ModifyRecord(domainId int64, record model.Record) bool {
	val := url.Values{}
	val.Add("login_token", GetLoginToken())
	val.Add("format", "json")
	val.Add("sub_domain", record.Name)
	val.Add("record_type", record.Type)
	val.Add("value", record.Value)
	val.Add("record_id", fmt.Sprintf("%d", record.ID))
	val.Add("record_line", "默认")
	val.Add("domain_id", fmt.Sprintf("%d", domainId))
	val.Add("mx", "20")
	resp, err := http.PostForm("https://dnsapi.cn/Record.Modify", val)
	if err != nil {
		ShowErrMsg(fmt.Sprintf("更新记录报错:%s", err.Error()))
		return false
	}
	bb, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		ShowErrMsg(fmt.Sprintf("更新记录报错:%s", err.Error()))
		return false
	}
	body := string(bb)
	code := gjson.Get(body, "status.code").Int()
	if code == 1 {
		ShowSuccMsg("更新成功")
		return true
	} else {
		ShowErrMsg("更新失败")
		return false
	}
}

func RemoveRecord(domainId int64, recordId int64) bool {
	val := url.Values{}
	val.Add("login_token", GetLoginToken())
	val.Add("format", "json")
	val.Add("record_id", fmt.Sprintf("%d", recordId))
	val.Add("domain_id", fmt.Sprintf("%d", domainId))
	resp, err := http.PostForm("https://dnsapi.cn/Record.Remove", val)
	if err != nil {
		ShowErrMsg(fmt.Sprintf("删除记录报错:%s", err.Error()))
		return false
	}
	bb, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		ShowErrMsg(fmt.Sprintf("删除记录报错:%s", err.Error()))
		return false
	}
	body := string(bb)
	code := gjson.Get(body, "status.code").Int()
	if code == 1 {
		ShowSuccMsg("删除成功")
		return true
	} else {
		ShowErrMsg("删除失败")
		return false
	}
}
