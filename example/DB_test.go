package main

import (
	"database/sql"
	"fmt"
	"testing"

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
	db, err := sql.Open("oci8", "username/password@M4DEV")
	defer db.Close()
	data, _ := marisfrolg_utils.GetDataBySQL("select * from brand ", db)
	fmt.Println(data)
}
