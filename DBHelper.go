package marisfrolg_utils

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"math/big"
	"reflect"
	"strconv"
	"strings"

	"github.com/SAP/go-hdb/driver"
	"github.com/jackc/pgtype"
	"github.com/jackc/pgx"
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
func ExecuteNonQueryByTran(db *sql.DB, SqlList []string) (err error) {
	var (
		tx     *sql.Tx
		result sql.Result
		//aff    int64
	)
	tx, err = db.Begin() //开启事务
	if err != nil {
		goto ERR
	}
	for i := 0; i < len(SqlList); i++ {
		result, err = tx.Exec(SqlList[i])
		if err != nil {
			goto ERR
		}
		if result == nil {
			err = fmt.Errorf(`事务执行出错`)
			goto ERR
		}
		_, err = result.RowsAffected()
		if err != nil {
			goto ERR
		}
	}
	tx.Commit()
	return
ERR:
	//回滚
	tx.Rollback()
	return
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
//Parameters 数据查询参数
//Field 字段名
func GetSqlList(Parameters string, Field string) string {
	a := strings.Index(Parameters, ",")
	if a == -1 {
		var condition string
		Parameters = strings.Replace(Parameters, " ", "", -1)
		Parameters = strings.Replace(Parameters, "“", "", -1)
		Parameters = strings.Replace(Parameters, "”", "", -1)
		Parameters = strings.Replace(Parameters, "'", "", -1)
		Parameters = strings.Replace(Parameters, "‘", "", -1)
		Parameters = strings.Replace(Parameters, "，", "", -1)
		Parameters = strings.Replace(Parameters, "[", "", -1)
		Parameters = strings.Replace(Parameters, "]", "", -1)
		condition = " AND " + Field + " = " + "'" + Parameters + "'"
		return condition
	} else {
		var list []string
		var fz = ""
		var condition = Field + "="
		var bz = " AND ("
		Parameters += ","
		Parameters = strings.Replace(Parameters, " ", "", -1)
		Parameters = strings.Replace(Parameters, "\"", "", -1)
		Parameters = strings.Replace(Parameters, "“", "", -1)
		Parameters = strings.Replace(Parameters, "”", "", -1)
		Parameters = strings.Replace(Parameters, "'", "", -1)
		Parameters = strings.Replace(Parameters, "‘", "", -1)
		Parameters = strings.Replace(Parameters, "，", "", -1)
		Parameters = strings.Replace(Parameters, "[", "", -1)
		Parameters = strings.Replace(Parameters, "]", "", -1)
		for _, a := range Parameters {
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
//Parameters 数据结果列请求参数
func StringToRuneArr(Parameters string) []string {
	Parameters = strings.Replace(Parameters, " ", "", -1)
	Parameters = strings.Replace(Parameters, "\"", "", -1)
	Parameters = strings.Replace(Parameters, "“", "", -1)
	Parameters = strings.Replace(Parameters, "”", "", -1)
	Parameters = strings.Replace(Parameters, "'", "", -1)
	Parameters = strings.Replace(Parameters, "‘", "", -1)
	// s = strings.Replace(s, ",", "", -1)
	Parameters = strings.Replace(Parameters, "，", "", -1)
	Parameters = strings.Replace(Parameters, "[", "", -1)
	Parameters = strings.Replace(Parameters, "]", "", -1)
	Result := strings.Split(Parameters, ",")
	return Result
}

//万能SQL语句查询 专门用于oracle
func GetDataBySQL(SQL string, db *sql.DB) (data []map[string]interface{}, err error) {
	//危险语句检查
	if strings.Contains(strings.ToUpper(SQL), `INSERT `) || strings.Contains(strings.ToUpper(SQL), `UPDATE `) || strings.Contains(strings.ToUpper(SQL), `DELETE `) || strings.Contains(strings.ToUpper(SQL), `TRUNCATE `) || strings.Contains(strings.ToUpper(SQL), `GRANT `) {
		return nil, errors.New("危险语句禁止执行")
	}
	if err != nil {
		return
	}
	rows, err := db.Query(SQL)
	defer rows.Close()
	if err != nil {
		return
	}
	var (
		refs   []interface{}
		cnt    int64 //首行处理
		cols   []string
		indexs []int
	)

	data = make([]map[string]interface{}, 0)
	columns, _ := rows.Columns()
	columnTypes, _ := rows.ColumnTypes()
	for rows.Next() {
		if cnt == 0 {
			indexs = make([]int, 0, len(columns))
			cols = columns
			refs = make([]interface{}, len(cols))
			for i := range refs {
				typeName := columnTypes[i].DatabaseTypeName()
				if typeName == "SQLT_NUM" || typeName == "SQLT_BDOUBLE" || typeName == "SQLT_INT" || typeName == "SQLT_FLT" || typeName == "SQLT_BFLOAT" {
					var ref sql.NullFloat64
					refs[i] = &ref
				} else if typeName == "SQLT_DAT" || typeName == "SQLT_TIMESTAMP" || typeName == "SQLT_TIMESTAMP_TZ" {
					var ref sql.NullTime
					refs[i] = &ref
				} else {
					var ref sql.NullString
					refs[i] = &ref
				}
				indexs = append(indexs, i)
			}
		}

		if err := rows.Scan(refs...); err != nil {
			return nil, err
		}
		params := make(map[string]interface{}, len(cols))
		for _, i := range indexs {
			ref := refs[i]
			typeName := columnTypes[i].DatabaseTypeName()
			if typeName == "SQLT_NUM" || typeName == "SQLT_BDOUBLE" || typeName == "SQLT_INT" || typeName == "SQLT_FLT" || typeName == "SQLT_BFLOAT" {
				value := reflect.Indirect(reflect.ValueOf(ref)).Interface().(sql.NullFloat64)
				if value.Valid {
					params[cols[i]] = value.Float64
				} else {
					params[cols[i]] = nil
				}
			} else if typeName == "SQLT_DAT" || typeName == "SQLT_TIMESTAMP" || typeName == "SQLT_TIMESTAMP_TZ" {
				value := reflect.Indirect(reflect.ValueOf(ref)).Interface().(sql.NullTime)
				if value.Valid {
					params[cols[i]] = value.Time.Format("2006-01-02T15:04:05")
				} else {
					params[cols[i]] = nil
				}
			} else {
				value := reflect.Indirect(reflect.ValueOf(ref)).Interface().(sql.NullString)
				if value.Valid {
					params[cols[i]] = value.String
				} else {
					params[cols[i]] = nil
				}
			}
		}
		data = append(data, params)
		cnt++
	}

	return data, nil
}

//万能查询PG数据库
//liqifeng 解决了读取Numeric类型数据不正常问题
func GetDataByPostgresSql(querySql string, conn *pgx.Conn) (data []map[string]interface{}, err error) {
	//危险语句检查
	if strings.Contains(strings.ToUpper(querySql), `INSERT `) || strings.Contains(strings.ToUpper(querySql), `UPDATE `) ||
		strings.Contains(strings.ToUpper(querySql), `DELETE `) || strings.Contains(strings.ToUpper(querySql), `TRUNCATE `) ||
		strings.Contains(strings.ToUpper(querySql), `GRANT `) {
		return nil, errors.New("危险语句禁止执行")
	}
	rows, err := conn.Query(context.Background(), querySql)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	fields := rows.FieldDescriptions()
	for rows.Next() {
		values, err := rows.Values()
		if err != nil {
			return nil, err
		}
		dataItem := make(map[string]interface{}, 0)
		for i, v := range fields {
			value, ok := values[i].(pgtype.Numeric)
			if ok {
				driverValue, _ := value.Value()
				_value := fmt.Sprintf(`%v`, driverValue)
				dataItem[string(v.Name)], _ = strconv.ParseFloat(_value, 64)
			} else {
				dataItem[string(v.Name)] = values[i]
			}
		}
		data = append(data, dataItem)
	}
	if rows.Err() != nil {
		return nil, rows.Err()
	}
	return
}

//万能hana查询语句
func GetDataByHanaSql(querySql string, db *sql.DB) ([]map[string]interface{}, error) {
	//危险语句检查
	if strings.Contains(strings.ToUpper(querySql), `INSERT `) || strings.Contains(strings.ToUpper(querySql), `UPDATE `) || strings.Contains(strings.ToUpper(querySql), `DELETE `) || strings.Contains(strings.ToUpper(querySql), `TRUNCATE `) || strings.Contains(strings.ToUpper(querySql), `GRANT `) {
		return nil, errors.New("危险语句禁止执行")
	}
	rows, err := db.Query(querySql)
	defer rows.Close()
	if err != nil {
		return nil, err
	}
	var (
		refs   []interface{}
		cnt    int64 //首行处理
		cols   []string
		indexs []int
	)
	data := make([]map[string]interface{}, 0)
	columns, _ := rows.Columns()
	columnTypes, _ := rows.ColumnTypes()
	for rows.Next() {
		if cnt == 0 {
			indexs = make([]int, 0, len(columns))
			cols = columns
			refs = make([]interface{}, len(cols))
			for i := range refs {
				typeName := columnTypes[i].DatabaseTypeName()
				if typeName == "DECIMAL" {
					var dec driver.Decimal
					var ref driver.NullDecimal
					ref.Decimal = &dec
					refs[i] = &ref
				} else if typeName == "DOUBLE" {
					var ref sql.NullFloat64
					refs[i] = &ref
				} else if typeName == "INTERAGE" {
					var ref sql.NullInt64
					refs[i] = &ref
				} else if typeName == "DATE" {
					var ref sql.NullTime
					refs[i] = &ref
				} else {
					var ref sql.NullString
					refs[i] = &ref
				}
				indexs = append(indexs, i)
			}
		}

		if err := rows.Scan(refs...); err != nil {
			return nil, err
		}
		params := make(map[string]interface{}, len(cols))

		for _, i := range indexs {
			ref := refs[i]

			typeName := columnTypes[i].DatabaseTypeName()
			if typeName == "DECIMAL" {
				if nullDecimal := reflect.Indirect(reflect.ValueOf(ref)).Interface().(driver.NullDecimal); nullDecimal.Valid {
					decimal, _ := (*big.Rat)(nullDecimal.Decimal).Float64()
					params[cols[i]] = decimal
				} else {
					params[cols[i]] = nil
				}
			} else if typeName == "DOUBLE" {
				value := reflect.Indirect(reflect.ValueOf(ref)).Interface().(sql.NullFloat64)
				if value.Valid {
					params[cols[i]] = value.Float64
				} else {
					params[cols[i]] = nil
				}
			} else if typeName == "INTERAGE" {
				value := reflect.Indirect(reflect.ValueOf(ref)).Interface().(sql.NullInt32)
				if value.Valid {
					params[cols[i]] = value.Int32
				} else {
					params[cols[i]] = nil
				}
			} else if typeName == "DATE" {
				value := reflect.Indirect(reflect.ValueOf(ref)).Interface().(sql.NullTime)
				if value.Valid {
					params[cols[i]] = value.Time.Format("2006-01-02T15:04:05")
				} else {
					params[cols[i]] = nil
				}
			} else {
				value := reflect.Indirect(reflect.ValueOf(ref)).Interface().(sql.NullString)
				if value.Valid {
					params[cols[i]] = value.String
				} else {
					params[cols[i]] = nil
				}
			}
		}
		data = append(data, params)
		cnt++
	}

	return data, nil
}
