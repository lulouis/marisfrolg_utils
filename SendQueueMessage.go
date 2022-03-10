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

var conn *amqp.Connection

func SendQueueMessage(currentMode string, rabbitmq_conn string, businessName string, messageBody interface{}, mongoSession *mgo.Session) (err error) {
	defer func() {
		if err := recover(); err != nil {
			logBody := bson.M{}
			logBody["currentMode"] = currentMode
			logBody["rabbitmq_conn"] = rabbitmq_conn
			logBody["businessName"] = businessName
			logBody["messageBody"] = messageBody
			body, _ := json.Marshal(logBody)
			AddOperationLog("marisfrolg_utils", "SendQueueMessage", fmt.Sprintf("错误详情:%s\n传入参数:%s \n", err, string(body)), "Log")
		}
	}()
	c := mongoSession.DB("ODSAPP").C("QueueIssue")
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
	if conn == nil || conn.IsClosed() {
		conn, err = amqp.Dial(rabbitmq_conn)
		if err != nil {
			//记录mongo日志
			if currentMode == "PRD" {
				c.Insert(issue)
			}
			return
		}
	}

	// 创建一个channel
	ch, err := conn.Channel()
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
