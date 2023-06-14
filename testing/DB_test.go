package testing

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"testing"

	_ "github.com/go-sql-driver/mysql"
	"github.com/lulouis/marisfrolg_utils"
	_ "github.com/mattn/go-oci8"
)

func TestExecuteNonQueryByTran(t *testing.T) {
	var (
		db      *sql.DB
		SqlList []string
	)
	SqlList = append(SqlList, "INSERT INTO TableName VALUER ('1')")
	SqlList = append(SqlList, "INSERT INTO TableName VALUER ('2')")
	r := marisfrolg_utils.ExecuteNonQueryByTran(db, SqlList)
	fmt.Println(r)
}

func TestAssemblyParameters(t *testing.T) {
	r := marisfrolg_utils.AssemblyParameters("'987654321','123456789'", "ABC")
	fmt.Println(r)
}

func TestGetSqlList(t *testing.T) {
	r := marisfrolg_utils.GetSqlList("'123456789'", "ABC")
	fmt.Println(r)
}
func TestStringToRuneArr(t *testing.T) {
	r := marisfrolg_utils.StringToRuneArr("'123456789'")
	fmt.Println(r)
}

func TestGetDataBySQL(t *testing.T) {
	db, _ := sql.Open("oci8", "username/password@M4DEV")
	defer db.Close()
	data, _ := marisfrolg_utils.GetDataBySQL("select * from TABLENAME where ROWNUM<10 ", db)
	for _, v := range data {
		str, _ := json.Marshal(v)
		fmt.Println(string(str))
	}

}

func TestGetMysqlDataBySql(t *testing.T) {
	//refs := make([]interface{}, 10)
	//for i,v := range refs {
	//	fmt.Println(i,v)
	//}
	db, err := sql.Open("mysql", "username:password(id:port)/yourDbName?charset=utf8&parseTime=true&loc=Local")
	if err != nil {
		fmt.Println(err)
		return
	}
	defer db.Close()
	if err = db.Ping(); err != nil {
		fmt.Println(err)
		return
	}
	data, err := marisfrolg_utils.GetDataByMysql("select * from PRICE where MAT_ID =139835", db)
	if err != nil {
		fmt.Println(err)
		return
	}
	for _, v := range data {
		fmt.Println(v)
		for k, v1 := range v {
			fmt.Println(k, v1)
		}
	}

}
