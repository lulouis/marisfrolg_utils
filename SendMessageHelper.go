package marisfrolg_utils

import (
	"fmt"
	"github.com/aliyun/alibaba-cloud-sdk-go/sdk"
	"github.com/aliyun/alibaba-cloud-sdk-go/sdk/requests"
)

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
		panic(err)
	}
	fmt.Print(response.GetHttpContentString())

	return response.GetHttpContentString()
}