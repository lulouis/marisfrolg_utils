package marisfrolg_utils

import (
	"encoding/json"
	"fmt"
	"github.com/DeanThompson/jpush-api-go-client"
	"github.com/DeanThompson/jpush-api-go-client/push"
	"github.com/aliyun/alibaba-cloud-sdk-go/sdk"
	"github.com/aliyun/alibaba-cloud-sdk-go/sdk/requests"
	gosocketio "github.com/graarh/golang-socketio"
	"github.com/graarh/golang-socketio/transport"
	"strings"
)

type JgVerifyResponse struct {
	Id float64 `json:"id"`
	Code int `json:"code"`
	Content string `json:"content"`
	ExID string `json:"exID"`
	Phone string `json:"phone"`
}

type Text struct {
	Content string `json:"content"`
}

type Textcard struct {
	Title       string `json:"title"`
	Description string `json:"description"`
	Url         string `json:"url"`
}
type File struct {
	Media_id string `json:"media_id"`
}
type Message struct { //touser、toparty、totag不能同时为空
	Touser   string   `json:"touser"` //员工工号，多人用|隔开。(如：xxxxx|xxxxx)
	Toparty  string   `json:"toparty"`
	Totag    string   `json:"totag"`
	Msgtype  string   `json:"msgtype"`
	Agentid  int      `json:"agentid"`  //企业应用的id
	Safe     int      `json:"safe"`     //表示是否是保密消息
	Text     Text     `json:"text"`     //Msgtype为text时方可使用
	File     File     `json:"file"`     //Msgtype为file时方可使用
	Textcard Textcard `json:"textcard"` //Msgtype为textcard时方可使用
}

/// <summary>
/// 阿里短信发送接口
/// 短信服务SDK简介 https://help.aliyun.com/document_detail/101874.html?spm=a2c4g.11186623.6.592.78ba5f30ZXZjpI
/// </summary>
/// <param name="PhoneNumbers">
/// 短信接收号码,支持以逗号分隔的形式进行批量调用，
/// 批量上限为800个手机号码,批量调用相对于单条调用及时性稍有延迟,
/// 验证码类型的短信推荐使用单条调用的方式</param>
/// <param name="SignName">短信签名</param>
/// <param name="TemplateCode">短信模板ID</param>
/// <param name="TemplateParam">
/// 可选参数,
/// 短信模板变量替换JSON串,友情提示:如果JSON中需要带换行符,
/// 请参照标准的JSON协议对换行符的要求,比如短信内容中包含\r\n的情况在JSON中需要表示成\r\n,
/// 否则会导致JSON在服务端解析失败
/// </param>
/// <returns></returns>
func SendShortMessage(regionId string, accessKeyId string, accessKeySecret string,PhoneNumbers string, SignName string, TemplateCode string, TemplateParam string) string {
	client, err := sdk.NewClientWithAccessKey(regionId, accessKeyId, accessKeySecret)
	if err != nil {
		panic(err)
	}
	request := requests.NewCommonRequest()
	request.Method = "POST"
	request.Scheme = "https" // https | http
	request.Domain = "dysmsapi.aliyuncs.com"
	request.Version = "2017-05-25"
	request.ApiName = "SendSms"
	request.QueryParams["RegionId"] = "cn-hangzhou"
	request.QueryParams["PhoneNumbers"] = PhoneNumbers
	request.QueryParams["SignName"] = SignName
	request.QueryParams["TemplateCode"] = TemplateCode
	request.QueryParams["TemplateParam"] = TemplateParam

	response, err := client.ProcessCommonRequest(request)
	if err != nil {
		fmt.Println(err)
	}
	fmt.Print(response.GetHttpContentString())

	return response.GetHttpContentString()
}

// 发送通知给QRCODE changedType:PERMISSION_CHANGED MENU_CHANGED BADGE_CHANGED
func SendMessageToQrCode(userId, changedType,QRCODE_IP string,QRCODE_PORT int) {
	im, _ := gosocketio.Dial(
		gosocketio.GetUrl(QRCODE_IP, QRCODE_PORT, false),
		transport.GetDefaultWebsocketTransport())
	im.Emit(changedType, userId)
}

//发送推送消息
func SendJPush(appKey,masterSecret,registrationId, pushPlatform, logTitle, title, notice string, extraMap map[string]interface{}) (err error) {
	var (
		ids                 []string
		platform            *push.Platform
		audience            *push.Audience
		notification        *push.Notification
		androidNotification *push.AndroidNotification
		iosNotification     *push.IosNotification
		option              *push.Options
		payload             *push.PushObject
		result              *push.PushResult
	)
	var client = jpush.NewJPushClient(appKey, masterSecret)
	// platform 对象
	platform = push.NewPlatform()
	// 用 Add() 方法添加具体平台参数，可选: "all", "ios", "android"
	//platform.Add("ios")
	// 或者用 All() 方法设置所有平台
	//platform.All()
	if pushPlatform == "all" {
		platform.Add("ios", "android")
	} else {
		platform.Add(pushPlatform)
	}
	// audience 对象，表示消息受众
	audience = push.NewAudience()
	//audience.SetTag([]string{"广州", "深圳"})   // 设置 tag 并集
	//audience.SetTagAnd([]string{"北京", "女"}) // 设置 tag_and 交集
	// audience.SetRegistrationId([]string{"id1", "id2"})   // 设置 registration_id
	// 和 platform 一样，可以调用 All() 方法设置所有受众
	//audience.All()
	//根据用户工号找所有设备ID
	if registrationId != "all" {
		registrationIdArr := strings.Split(registrationId, ",")
		for _, v := range registrationIdArr {
			ids = append(ids, v)
		}
		audience.SetRegistrationId(ids)
	} else {
		audience.All()
	}

	if len(ids) > 0 || registrationId == "all" {
		// notification 对象，表示 通知，传递 alert 属性初始化
		notification = push.NewNotification(logTitle) //推送日志标题

		androidNotification = push.NewAndroidNotification(notice) // android 平台专有的 notification，用 alert 属性初始化
		iosNotification = push.NewIosNotification(notice)         // iOS 平台专有的 notification，用 alert 属性初始化
		// addExtra方法
		for k, v := range extraMap {
			androidNotification.AddExtra(k, v)
			iosNotification.AddExtra(k, v)
		}
		androidNotification.Title = title
		notification.Android = androidNotification

		iosNotification.Badge = 1
		iosNotification.Sound = "default"
		notification.Ios = iosNotification

		option = push.NewOptions()
		option.ApnsProduction = true

		// 可以调用 AddExtra 方法，添加额外信息
		// message.AddExtra("key", 123)
		//创建PushObject对象
		payload = push.NewPushObject()
		payload.Platform = platform
		payload.Audience = audience
		payload.Notification = notification
		//payload.Options = option

		// 打印查看 json 序列化的结果，也就是 POST 请求的 body
		data, err := json.Marshal(payload)
		if err != nil {
			err = fmt.Errorf("json.Marshal PushObject failed:%s\r\n", err)
			goto ERR
		} else {
			fmt.Println("payload:", string(data))
		}
		//开始推送
		result, err = client.Push(payload)
		if err != nil {
			err = fmt.Errorf("Push failed:%s\r\n", err)
			goto ERR
		} else {
			fmt.Println("Push result:", result)
		}
	}

	return
ERR:
	return
}

//使用go语言发企业微信消息
func SendWorkMessage(msg Message, reqUrl,token string) (result []byte, err error) {
	//reqUrl := "http://ip:端口/v1/Push/PushMessage"
	jsonData, err := json.Marshal(msg)
	if err!=nil{
		return nil, err
	}
	data := string(jsonData)
	//bearer := "Bearer " + token
	result, err = HttpPostOnlyBody(reqUrl,  data, token)
	if err != nil {
		fmt.Printf("PushMessage 请求接口失败 err:%s",err)
		return nil, err
	}
	return
}