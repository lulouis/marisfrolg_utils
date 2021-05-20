package testing

import (
	"fmt"
	"strings"
	"testing"

	"github.com/lulouis/marisfrolg_utils"
)

type User struct {
	Id      int
	Age     int
	PoMoney float64
}

// go test -v -run GetModelProperty .\testing\ModelReflect_test.go
func TestGetModelProperty(t *testing.T) {
	var ori_calc_formula = "(User.Id + User.Age) / User.PoMoney * 0.2"
	var calc_formula = ""

	//对象实体
	var user1 = new(User)
	user1.Id = 4
	user1.Age = 100
	user1.PoMoney = 100.05
	//值类型
	user2 := User{}
	user2.Id = 3
	user2.Age = 88
	user2.PoMoney = 88.05

	//表达式替换
	calc_formula = strings.ReplaceAll(ori_calc_formula, "User.Id", fmt.Sprint(marisfrolg_utils.GetModelProperty(*user1, "Id").(int)))
	calc_formula = strings.ReplaceAll(calc_formula, "User.Age", fmt.Sprint(marisfrolg_utils.GetModelProperty(*user1, "Age").(int)))
	calc_formula = strings.ReplaceAll(calc_formula, "User.PoMoney", fmt.Sprint(marisfrolg_utils.GetModelProperty(*user1, "PoMoney").(float64)))

	total := marisfrolg_utils.Calculate(calc_formula)
	fmt.Printf("user1:%s = %s\n", ori_calc_formula, fmt.Sprint(total.(float64)))

	calc_formula = strings.ReplaceAll(ori_calc_formula, "User.Id", fmt.Sprint(marisfrolg_utils.GetModelProperty(user2, "Id").(int)))
	calc_formula = strings.ReplaceAll(calc_formula, "User.Age", fmt.Sprint(marisfrolg_utils.GetModelProperty(user2, "Age").(int)))
	calc_formula = strings.ReplaceAll(calc_formula, "User.PoMoney", fmt.Sprint(marisfrolg_utils.GetModelProperty(user2, "PoMoney").(float64)))

	total = marisfrolg_utils.Calculate(calc_formula)
	fmt.Printf("user2:%s = %s\n", ori_calc_formula, fmt.Sprint(total.(float64)))

}
