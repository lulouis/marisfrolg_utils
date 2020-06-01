package main

import (
	"database/sql"
	"fmt"
	"github.com/lulouis/marisfrolg_utils"
	"testing"
)

func TestExecuteNonQueryByTran(t *testing.T) {
	var (
		db *sql.DB
		SqlList []string
	)
	SqlList=append(SqlList,"INSERT INTO TableName VALUER ('1')")
	SqlList=append(SqlList,"INSERT INTO TableName VALUER ('2')")
	r := marisfrolg_utils.ExecuteNonQueryByTran(db,SqlList)
	fmt.Println(r)
}

func TestAssemblyParameters(t *testing.T) {
	r := marisfrolg_utils.AssemblyParameters("'987654321','123456789'", "ABC")
	fmt.Println(r)
}
