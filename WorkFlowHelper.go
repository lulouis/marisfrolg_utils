package marisfrolg_utils

import (
	"errors"
	"time"

	"gopkg.in/mgo.v2/bson"
	"net/url"
	"encoding/json"
	"fmt"
)

//流程引擎
type WorkFlowCase struct {
	ID bson.ObjectId `bson:"_id"`
	//流程模板代码
	TemplateCode string `bson:"templateCode"`
	//业务内容
	BusinessObject interface{} `bson:"businessObject"`
	//创建人工号
	UserCode string `bson:"usercode"`
	//创建人
	UserName string `bson:"username"`
	//创建日期
	CreateTime time.Time `bson:"createtime"`
	//执行状态0表示未开始，1表示审批中，2表示已完成，3表示已拒绝
	Status int `bson:"status"`
	//失败的原因代码，101创建完调用第三方失败 102执行动作完调用第三方失败 103结束完成时调用第三方失败
	ErrCode string `bson:"errCode"`
	//任务盒子
	BoxList []Box `bson:"boxList"`
	//操作记录
	ActionRecord []Node `bson:"actionRecord"`
}

//动作选项
type Action struct {
	//名称
	ActionName string `bson:"actionName"`
	//描述
	ActionDesc string `bson:"actionDesc"`
}

//审批节点
type Node struct {
	//名称
	Name string `bson:"name"`
	//描述
	Desc string `bson:"desc"`
	//人员标签
	Tag string `bson:"tag"`
	//查找范围，MA,SU,AM,*(表示全集团),空(表示逐层查找)
	Scope string `bson:"scope"`
	//标签分组
	SystemName string `bson:"systemName"`
	//动作选项
	Actions []Action `bson:"actions"`
	//预计审批人
	PlanningLeader string `bson:"planningLeader"`
	//预计审批人名称
	PlanningLeaderName string `bson:"planningLeaderName"`
	//实际操作人
	Operator string `bson:"operator"`
	//实际操作人名称
	OperatorName string `bson:"operatorName"`
	//操作日期
	OperatorTime time.Time `bson:"operatorTime"`
	//最后动作记录
	ActionName string `bson:"actionName"`
	//最后动作描述
	ActionDesc string `bson:"actionDesc"`
	//0未执行，1动作执行，2逻辑上跳过，2逻辑上关闭
	Flag int `bson:"flag"`
	//打回到哪个Node
	BackNodeName string `bson:"BackNodeName"`
	BackNodeDesc string `bson:"BackNodeDesc"`
}

//盒子，图形化对象
type Box struct {
	//节点列表
	Nodes []Node `bson:"nodes"`
	//逻辑与或(默认与)
	AllYes bool `bson:"allYes"`
	//代理人
	Agent []string `bson:"agent"`
}

//判断权限需要的标签类
type UserTag struct {
	ID      string      `bson:"ID"`
	Name    string      `bson:"Name"`
	TagList []MFTagList `bson:"mfTagList"`
}
type MFTagList struct {
	TagIdList []TagId `bson:"tagId"`
}
type TagId struct {
	TId   string `bson:"tId"`
	TName string `bson:"tName"`
}

type ErrorResponse struct {
	Code    int `json:"code"`
	Message string `json:"message"`
}

type ActionRequest struct {
	WorkFlowID   string `json:"workFlowID"`
	Token        string `json:"token"`
	Node        *Node `json:"node"`
	Reason string `json:"reason"`
}
//解析流程引擎返回代办事项
func GetApproverList(WorkFlowCase_ *WorkFlowCase) (list []*Node, err error) {
	if WorkFlowCase_==nil || WorkFlowCase_.ID == "" {
		err = errors.New("流程创建失败")
		return
	}
	var BoxList = WorkFlowCase_.BoxList
	for _, v1 := range BoxList {
		var Node = v1.Nodes
		for _, v2 := range Node {
			if v2.Flag == 0 && v2.PlanningLeader != "" { //未执行并且操作人不为空
				list = append(list, &v2)
			}
		}
		if len(list) > 0 {
			break
		}
	}
	return
}

//创建流程  单据类型、用户id、姓名、工单编号
func CreateWorkFlow( businessObj map[string]interface{},userId , UserName , workFlowTag ,token,createUrl  string) (data string, err error) {
	var(
		businessBytes,responseBytes []byte
	)
	//businessObject集合
	//addition := make(map[string]interface{})
	//addition["addUserId"] = userId
	//addition["billType"] = billType
	//addition["billNo"] = billNo
	//addition["systemName"] = systemName
	//addition["brandCode"] = brand
	//初始化参数
	param := url.Values{}
	//配置请求参数,方法内部已处理urlencode问题,中文参数可以直接传参
	//templateCode   businessObject  UserCode  UserName  流程标签     业务数据   创建人工号   创建人姓名
	businessBytes, _ = json.Marshal(businessObj)
	param.Set("templateCode", workFlowTag)
	param.Set("businessObject", string(businessBytes))
	param.Set("UserCode", userId)
	param.Set("UserName", UserName)
	if responseBytes, err = HttpPostWithReqParamAndToken(createUrl, param, "", token);err!=nil{
		err = fmt.Errorf(`请求失败:%s`,err)
		return "", err
	}
	data = string(responseBytes)
	return
}

//执行审核流程，通过、拒绝、驳回等
func ActionWorkFlow(actionUrl,workFlowID string, node *Node, token string) (data string, err error) {
	var(
		responseBytes,nodeBytes []byte
	)
	param := url.Values{}                //初始化参数
	param.Set("workFlowId", workFlowID)
	nodeBytes, _ = json.Marshal(node)
	responseBytes, err = HttpPostWithReqParamAndToken(actionUrl, param, string(nodeBytes), token)
	if err != nil {
		err = errors.New("请求失败,我方错误信息,data=" + string(data))
	}
	data = string(responseBytes)
	fmt.Println("执行审批流程打印:",data)
	return
}

//解析流程引擎数据
func GetWorkFlow(data string) (WorkFlow *WorkFlowCase, err error) {
	WorkFlow = new(WorkFlowCase)
	err = json.Unmarshal([]byte(data), &WorkFlow)
	if err != nil || WorkFlow.ID == "" {
		err = fmt.Errorf(`创建失败"%s"`,err)
		return
	}
	return
}

//获取下一个审批节点
func GetNextNodeFromWorkFlow(boxList []Box)(node Node){
	node = Node{}
	for _, v1 := range boxList {
		for _, v2 := range v1.Nodes {
			if v2.Flag == 0 && v2.PlanningLeader != "" { //未执行并且操作人不为空
				node = v2
				return
			}
		}

	}
	return
}