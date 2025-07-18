package marisfrolg_utils

import (
	"encoding/json"
	"errors"
	"fmt"
	"os"
	"time"

	"github.com/streadway/amqp"
	"gopkg.in/mgo.v2"
	"gopkg.in/mgo.v2/bson"
)

var conn1 *amqp.Connection
var conn2 *amqp.Connection
var conn3 *amqp.Connection

type SendQueueMessageInput struct {
	CurrentMode          string
	RabbitmqConn         string
	BusinessName         string
	MessageBody          interface{}
	QueueIssueCollection *mgo.Collection
	MongoSession         *mgo.Session
}

func SendQueueMessage(currentMode string, rabbitmq_conn string, businessName string, messageBody interface{}, mongoSession *mgo.Session) (err error) {
	logInfo := NewAuditLogInfo("marisfrolg_utils", GetMethodName(), "00000")
	logInfo.SetRequest(map[string]interface{}{
		"currentMode":   currentMode,
		"rabbitmq_conn": rabbitmq_conn,
		"businessName":  businessName,
		"messageBody":   messageBody,
	})
	defer func() {
		if rec := recover(); rec != nil {
			err = fmt.Errorf("%v", rec)
		}
		if err != nil {
			AuditLog(logInfo, err)
		}
	}()
	c := mongoSession.DB("OC-QUEUE").C("QueueIssue")
	issue := bson.M{}
	issue["_id"] = bson.NewObjectId()
	// 声明一个队列
	queueName := ""
	switch currentMode {
	case "DEV":
		queueName = businessName + "_dev"
	case "TEST":
		queueName = businessName + "_dev"
	case "PRD":
		queueName = businessName + "_prd"
		issue["queueName"] = queueName
		issue["queueDualConn"] = rabbitmq_conn
		issue["queueMessage"] = messageBody
		issue["currentMode"] = currentMode
		issue["createTime"] = time.Now()
		issue["finishedStatus"] = 0 //0未处理,1已完结,2取消处理
	}
	// 连接RabbitMQ服务器
	if conn1 == nil || conn1.IsClosed() {
		conn1, err = amqp.Dial(rabbitmq_conn)
		if err != nil {
			//记录mongo日志
			if currentMode == "PRD" {
				c.Insert(issue)
			}
			return
		}
	}

	// 创建一个channel
	ch, err := conn1.Channel()
	if err != nil {
		//记录mongo日志
		if currentMode == "PRD" {
			c.Insert(issue)
		}
		return
	}
	defer ch.Close()
	//交换机检查x-dead-letter-exchange-all-business
	ch.ExchangeDeclare("x-dead-letter-exchange-all-business", "direct", true, false, false, false, nil)
	// 创建本业务死信
	ch.QueueDeclare(
		businessName+"_dead", // 队列名称
		true,                 // 是否持久化
		false,                // 是否自动删除
		false,                // 是否独占
		false, nil,
	)
	ch.QueueBind(businessName+"_dead", businessName+"_dead", "x-dead-letter-exchange-all-business", false, nil)

	args := make(map[string]interface{}, 0)
	args["x-dead-letter-exchange"] = "x-dead-letter-exchange-all-business"
	args["x-dead-letter-routing-key"] = businessName + "_dead"
	q, err := ch.QueueDeclare(
		queueName, // 队列名称
		true,      // 是否持久化
		false,     // 是否自动删除
		false,     // 是否独占
		false, args,
	)
	if err != nil {
		//记录mongo日志
		if currentMode == "PRD" {
			c.Insert(issue)
		}
		err = errors.New("连接" + queueName + "队列时失败")
		return
	}
	// 发送消息到队列中
	// body := fmt.Sprintf(`{"name":"刘宇辉","id":1}`)
	body, err := json.Marshal(messageBody)
	if err != nil {
		//记录mongo日志
		if currentMode == "PRD" {
			c.Insert(issue)
		}
		err = errors.New("消息序列化失败")
		return
	}

	err = ch.Publish(
		"",     // exchange
		q.Name, // routing key
		false,  // mandatory
		false,  // immediate
		amqp.Publishing{
			ContentType:  "application/json",
			Body:         body,
			DeliveryMode: 2, //持久化
			AppId:        os.Getenv("APP_NAME"),
			Timestamp:    time.Now(),
		},
	)
	if err != nil {
		//记录mongo日志
		if currentMode == "PRD" {
			c.Insert(issue)
		}
		err = errors.New(queueName + "发送队列消息失败")
		return
	}

	return
}

// 用于推送OC队列与OC故障
func SendQueueMessageV2(input *SendQueueMessageInput) (err error) {
	logInfo := NewAuditLogInfo("marisfrolg_utils", GetMethodName(), "00000")
	logInfo.SetRequest(input)
	defer func() {
		if rec := recover(); rec != nil {
			err = fmt.Errorf("%v", rec)
		}
		if err != nil {
			AuditLog(logInfo, err)
		}
	}()
	c := input.QueueIssueCollection
	issue := bson.M{}
	issue["_id"] = bson.NewObjectId()
	// 声明一个队列
	queueName := ""
	switch input.CurrentMode {
	case "DEV":
		queueName = input.BusinessName + "_dev"
	case "TEST":
		queueName = input.BusinessName + "_dev"
	case "PRD":
		queueName = input.BusinessName + "_prd"
		issue["queueName"] = queueName
		issue["queueDualConn"] = input.RabbitmqConn
		issue["queueMessage"] = input.MessageBody
		issue["currentMode"] = input.CurrentMode
		issue["createTime"] = time.Now()
		issue["finishedStatus"] = 0 //0未处理,1已完结,2取消处理
	}
	// 连接RabbitMQ服务器
	if conn2 == nil || conn2.IsClosed() {
		conn2, err = amqp.Dial(input.RabbitmqConn)
		if err != nil {
			//记录mongo日志
			if input.CurrentMode == "PRD" {
				c.Insert(issue)
			}
			return
		}
	}

	// 创建一个channel
	ch, err := conn2.Channel()
	if err != nil {
		//记录mongo日志
		if input.CurrentMode == "PRD" {
			c.Insert(issue)
		}
		return
	}
	defer ch.Close()
	//交换机检查x-dead-letter-exchange-all-business
	ch.ExchangeDeclare("x-dead-letter-exchange-all-business", "direct", true, false, false, false, nil)
	// 创建本业务死信
	ch.QueueDeclare(
		input.BusinessName+"_dead", // 队列名称
		true,                       // 是否持久化
		false,                      // 是否自动删除
		false,                      // 是否独占
		false, nil,
	)
	ch.QueueBind(input.BusinessName+"_dead", input.BusinessName+"_dead", "x-dead-letter-exchange-all-business", false, nil)

	args := make(map[string]interface{}, 0)
	args["x-dead-letter-exchange"] = "x-dead-letter-exchange-all-business"
	args["x-dead-letter-routing-key"] = input.BusinessName + "_dead"
	q, err := ch.QueueDeclare(
		queueName, // 队列名称
		true,      // 是否持久化
		false,     // 是否自动删除
		false,     // 是否独占
		false, args,
	)
	if err != nil {
		//记录mongo日志
		if input.CurrentMode == "PRD" {
			c.Insert(issue)
		}
		err = errors.New("连接" + queueName + "队列时失败")
		return
	}
	// 发送消息到队列中
	// body := fmt.Sprintf(`{"name":"刘宇辉","id":1}`)
	body, err := json.Marshal(input.MessageBody)
	if err != nil {
		//记录mongo日志
		if input.CurrentMode == "PRD" {
			c.Insert(issue)
		}
		err = errors.New("消息序列化失败")
		return
	}

	err = ch.Publish(
		"",     // exchange
		q.Name, // routing key
		false,  // mandatory
		false,  // immediate
		amqp.Publishing{
			ContentType:  "application/json",
			Body:         body,
			DeliveryMode: 2, //持久化
			AppId:        os.Getenv("APP_NAME"),
			Timestamp:    time.Now(),
		},
	)
	if err != nil {
		//记录mongo日志
		if input.CurrentMode == "PRD" {
			c.Insert(issue)
		}
		err = errors.New(queueName + "发送队列消息失败")
		return
	}

	return
}

// 与V2类似，大批量推送单号数组
func SendQueueMessageV3(input *SendQueueMessageInput, keys []string) (err error) {
	logInfo := NewAuditLogInfo("marisfrolg_utils", GetMethodName(), "00000")
	logInfo.SetRequest(input)
	defer func() {
		if rec := recover(); rec != nil {
			err = fmt.Errorf("%v", rec)
		}
		if err != nil {
			AuditLog(logInfo, err)
		}
	}()
	c := input.QueueIssueCollection
	issue := bson.M{}
	issue["_id"] = bson.NewObjectId()
	// 声明一个队列
	queueName := ""
	switch input.CurrentMode {
	case "DEV":
		queueName = input.BusinessName + "_dev"
	case "TEST":
		queueName = input.BusinessName + "_dev"
	case "PRD":
		queueName = input.BusinessName + "_prd"
		issue["queueName"] = queueName
		issue["queueDualConn"] = input.RabbitmqConn
		issue["queueMessage"] = input.MessageBody
		issue["currentMode"] = input.CurrentMode
		issue["createTime"] = time.Now()
		issue["finishedStatus"] = 0 //0未处理,1已完结,2取消处理
	}
	// 连接RabbitMQ服务器
	if conn3 == nil || conn3.IsClosed() {
		conn3, err = amqp.Dial(input.RabbitmqConn)
		if err != nil {
			//记录mongo日志
			if input.CurrentMode == "PRD" {
				c.Insert(issue)
			}
			return
		}
	}

	// 创建一个channel
	ch, err := conn3.Channel()
	if err != nil {
		//记录mongo日志
		if input.CurrentMode == "PRD" {
			c.Insert(issue)
		}
		return
	}
	defer ch.Close()
	//交换机检查x-dead-letter-exchange-all-business
	ch.ExchangeDeclare("x-dead-letter-exchange-all-business", "direct", true, false, false, false, nil)
	// 创建本业务死信
	ch.QueueDeclare(
		input.BusinessName+"_dead", // 队列名称
		true,                       // 是否持久化
		false,                      // 是否自动删除
		false,                      // 是否独占
		false, nil,
	)
	ch.QueueBind(input.BusinessName+"_dead", input.BusinessName+"_dead", "x-dead-letter-exchange-all-business", false, nil)

	args := make(map[string]interface{}, 0)
	args["x-dead-letter-exchange"] = "x-dead-letter-exchange-all-business"
	args["x-dead-letter-routing-key"] = input.BusinessName + "_dead"
	q, err := ch.QueueDeclare(
		queueName, // 队列名称
		true,      // 是否持久化
		false,     // 是否自动删除
		false,     // 是否独占
		false, args,
	)
	if err != nil {
		//记录mongo日志
		if input.CurrentMode == "PRD" {
			c.Insert(issue)
		}
		err = errors.New("连接" + queueName + "队列时失败")
		return
	}

	for _, key := range keys {
		// 发送消息到队列中
		// body := fmt.Sprintf(`{"name":"刘宇辉","id":1}`)
		var body []byte
		body, err = json.Marshal(key)
		if err != nil {
			//记录mongo日志
			if input.CurrentMode == "PRD" {
				c.Insert(issue)
			}
			err = errors.New("消息序列化失败")
			return
		}
		err = ch.Publish(
			"",     // exchange
			q.Name, // routing key
			false,  // mandatory
			false,  // immediate
			amqp.Publishing{
				ContentType:  "application/json",
				Body:         body,
				DeliveryMode: 2, //持久化
				AppId:        os.Getenv("APP_NAME"),
				Timestamp:    time.Now(),
			},
		)
		if err != nil {
			//记录mongo日志
			if input.CurrentMode == "PRD" {
				c.Insert(issue)
			}
			err = errors.New(queueName + "发送队列消息失败")
			return
		}
	}

	return
}

// 与V2类似，差异是一次连接，大量推送
func SendQueueMessageV4(input *SendQueueMessageInput, objects []interface{}) (err error) {
	logInfo := NewAuditLogInfo("marisfrolg_utils", GetMethodName(), "00000")
	logInfo.SetRequest(map[string]interface{}{
		"input":   input,
		"objects": objects,
	})
	defer func() {
		if rec := recover(); rec != nil {
			err = fmt.Errorf("%v", rec)
		}
		if err != nil {
			AuditLog(logInfo, err)
		}
	}()
	c := input.QueueIssueCollection
	issue := bson.M{}
	issue["_id"] = bson.NewObjectId()
	// 声明一个队列
	queueName := ""
	switch input.CurrentMode {
	case "DEV":
		queueName = input.BusinessName + "_dev"
	case "TEST":
		queueName = input.BusinessName + "_dev"
	case "PRD":
		queueName = input.BusinessName + "_prd"
		issue["queueName"] = queueName
		issue["queueDualConn"] = input.RabbitmqConn
		issue["queueMessage"] = input.MessageBody
		issue["currentMode"] = input.CurrentMode
		issue["createTime"] = time.Now()
		issue["finishedStatus"] = 0 //0未处理,1已完结,2取消处理
	}
	// 连接RabbitMQ服务器
	if conn3 == nil || conn3.IsClosed() {
		conn3, err = amqp.Dial(input.RabbitmqConn)
		if err != nil {
			//记录mongo日志
			if input.CurrentMode == "PRD" {
				c.Insert(issue)
			}
			return
		}
	}

	// 创建一个channel
	ch, err := conn3.Channel()
	if err != nil {
		//记录mongo日志
		if input.CurrentMode == "PRD" {
			c.Insert(issue)
		}
		return
	}
	defer ch.Close()
	//交换机检查x-dead-letter-exchange-all-business
	ch.ExchangeDeclare("x-dead-letter-exchange-all-business", "direct", true, false, false, false, nil)
	// 创建本业务死信
	ch.QueueDeclare(
		input.BusinessName+"_dead", // 队列名称
		true,                       // 是否持久化
		false,                      // 是否自动删除
		false,                      // 是否独占
		false, nil,
	)
	ch.QueueBind(input.BusinessName+"_dead", input.BusinessName+"_dead", "x-dead-letter-exchange-all-business", false, nil)

	args := make(map[string]interface{}, 0)
	args["x-dead-letter-exchange"] = "x-dead-letter-exchange-all-business"
	args["x-dead-letter-routing-key"] = input.BusinessName + "_dead"
	q, err := ch.QueueDeclare(
		queueName, // 队列名称
		true,      // 是否持久化
		false,     // 是否自动删除
		false,     // 是否独占
		false, args,
	)
	if err != nil {
		//记录mongo日志
		if input.CurrentMode == "PRD" {
			c.Insert(issue)
		}
		err = errors.New("连接" + queueName + "队列时失败")
		return
	}

	for _, obj := range objects {
		var body []byte
		// 判断是否为字符串
		if str, ok := obj.(string); ok {
			body = []byte(str)
		} else {
			body, _ = json.Marshal(obj)
		}
		err = ch.Publish(
			"",     // exchange
			q.Name, // routing key
			false,  // mandatory
			false,  // immediate
			amqp.Publishing{
				ContentType:  "application/json",
				Body:         body,
				DeliveryMode: 2, //持久化
				AppId:        os.Getenv("APP_NAME"),
				Timestamp:    time.Now(),
			},
		)
		if err != nil {
			//记录mongo日志
			if input.CurrentMode == "PRD" {
				c.Insert(issue)
			}
			err = errors.New(queueName + "发送队列消息失败")
			return
		}
	}
	return
}
