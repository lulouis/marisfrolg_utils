package marisfrolg_utils

import (
	"database/sql"
	"fmt"
	"strings"
)

/*
数据库相关操作
*/



/// <summary>
/// 批量执行 含事务
/// </summary>
/// <param name="db">数据连接</param>
/// <param name="SqlList">Sql 列表</param>
/// <returns>说明：MongoDB 禁止使用，其他数据库自行斟酌（目前支持Oracle）
///</returns>
func ExecuteNonQueryByTran(db *sql.DB, SqlList []string) error {
	var (
		err    error
		tx     *sql.Tx
		result sql.Result
		aff    int64
	)
	tx, err = db.Begin() //开启事务
	if err != nil {
		goto ERR
	}
	for _, v := range SqlList {
		result, err = tx.Exec(v)
		if err != nil {
			goto ERR
		}
		if result == nil {
			err = fmt.Errorf(`事务执行出错`)
			goto ERR
		}
		aff, err = result.RowsAffected()
		if err != nil {
			goto ERR
		}
		if aff == 0 {
		}
	}
	tx.Commit()
	return nil
ERR:
	//回滚
	tx.Rollback()
	return err
}


//SQL IN()的查询里不能超过1000列，将大于1000列的以900为间隔分开组装
//IdList:数据列 例如："'987654321','123456789'"
//Field:字段名
func AssemblyParameters(IdList, Field string) (condition string) {
	var ApplyIdList [100]string
	mbct := strings.Trim(IdList, ",")
	nameArr := strings.Split(mbct, ",")
	if len(nameArr) != 0 {
		for index, a := range nameArr {
			i := index / 900
			ApplyIdList[i] += a + ","
		}
		if ApplyIdList[0] != "" {
			if ApplyIdList[0] != "" && ApplyIdList[1] == "" {
				condition += fmt.Sprintf(" AND (%s IN (%s) ) ", Field, strings.TrimRight(ApplyIdList[0], ","))
			} else {
				for index, itm := range ApplyIdList {
					if itm != "" {
						if index == 0 {
							condition += fmt.Sprintf(" AND (%s IN (%s)  ", Field, strings.TrimRight(itm, ","))
						} else {
							condition += fmt.Sprintf(" OR %s IN (%s) ", Field, strings.TrimRight(itm, ","))
						}
					}
				}
				condition += fmt.Sprintf(")")
			}
		}
		return condition
	} else {
		return ""
	}
}


//对数据库查询请求参数的优化
func GetSqlList(s string, ziduanming string) string {
	a := strings.Index(s, ",")
	if a == -1 {
		var condition string
		s = strings.Replace(s, " ", "", -1)
		s = strings.Replace(s, "“", "", -1)
		s = strings.Replace(s, "”", "", -1)
		s = strings.Replace(s, "'", "", -1)
		s = strings.Replace(s, "‘", "", -1)
		s = strings.Replace(s, "，", "", -1)
		s = strings.Replace(s, "[", "", -1)
		s = strings.Replace(s, "]", "", -1)
		condition = " AND " + ziduanming + " = " + "'" + s + "'"
		return condition
	} else {
		var list []string
		var fz = ""
		var condition = ziduanming + "="
		var bz = " AND ("
		s += ","
		s = strings.Replace(s, " ", "", -1)
		s = strings.Replace(s, "\"", "", -1)
		s = strings.Replace(s, "“", "", -1)
		s = strings.Replace(s, "”", "", -1)
		s = strings.Replace(s, "'", "", -1)
		s = strings.Replace(s, "‘", "", -1)
		s = strings.Replace(s, "，", "", -1)
		s = strings.Replace(s, "[", "", -1)
		s = strings.Replace(s, "]", "", -1)
		for _, a := range s {
			if string(a) == "," {
				list = append(list, fz)
				fz = ""
			} else {
				fz = fz + string(a)
			}
		}
		for _, a := range list { //循环加
			bz = bz + condition + "'" + a + "'" + " or "
		}
		content := bz[0 : len(bz)-4]
		content = content + ")"
		return content
	}
}

//对HANA查询结果的特殊符号优化
func StringToRuneArr(s string) []string {
	s = strings.Replace(s, " ", "", -1)
	s = strings.Replace(s, "\"", "", -1)
	s = strings.Replace(s, "“", "", -1)
	s = strings.Replace(s, "”", "", -1)
	s = strings.Replace(s, "'", "", -1)
	s = strings.Replace(s, "‘", "", -1)
	// s = strings.Replace(s, ",", "", -1)
	s = strings.Replace(s, "，", "", -1)
	s = strings.Replace(s, "[", "", -1)
	s = strings.Replace(s, "]", "", -1)
	aa := strings.Split(s, ",")
	return aa
}
