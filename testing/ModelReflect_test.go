package testing

import (
	"fmt"
	"testing"

	"github.com/lulouis/marisfrolg_utils"
)

type User struct {
	Id  int
	Age int
}

// go test -v -run GetModelProperty .\testing\ModelReflect_test.go
func TestGetModelProperty(t *testing.T) {
	//对象实体
	var user1 = new(User)
	user1.Id = 123
	user1.Age = 100
	result1, _ := marisfrolg_utils.GetModelProperty(*user1, "Age")
	fmt.Println(result1)
	//值类型
	user2 := User{}
	user2.Id = 123
	user2.Age = 88
	result2, _ := marisfrolg_utils.GetModelProperty(user2, "Age")
	fmt.Println(result2)
}
