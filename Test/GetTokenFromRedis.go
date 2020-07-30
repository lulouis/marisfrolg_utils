package main

import (
	"encoding/json"
	"fmt"
	"github.com/go-redis/redis"
	"time"
)
type UserInfo struct {
	ClientID string `json:"ClientID"`
	Access string `json:"Access"`
}

func main() {
	//读取db2所有token值
	redisConn := "192.168.2.119:6379"
	client := redis.NewClient(&redis.Options{
		Addr:        redisConn,
		Password:    "", // no password set
		DB:          2,  // use default DB
		ReadTimeout: 240 * time.Second,
	})
	defer client.Close()
	//删除00221
	cmd:=client.Keys("*")
	keyList,err:=cmd.Result()
	if err!=nil{
		panic(err)
	}
	for _,k:=range keyList{
		data, err := client.Get(k).Result()
		if err!=nil{
			continue
		}
		u:=new(UserInfo)
		err = json.Unmarshal([]byte(data),&u)
		if err==nil{
			if u.ClientID=="00275" {//|| u.ClientID=="00275"
				fmt.Println("删除",u.Access)
				client.Del(u.Access)
			}
		}


	}

	//for _,v:=range keyList{
	//	//根据这个值来删除
	//	data, err := client.Get(v).Result()
	//}


}
