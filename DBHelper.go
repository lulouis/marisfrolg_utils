package Marisfrolg_utils


import (
	"database/sql"
	"fmt"
)

/// <summary>
/// 批量执行 含事务 True 成功 False 失败
/// </summary>
/// <param name="db">数据连接</param>
/// <param name="SqlList">Sql 列表</param>
/// <returns>说明：MongoDB 禁止使用，其他数据库自行斟酌（目前支持Oracle），如果修改请先联系创建人或者最后修改人，
/// 如果修改此方法请把最后修改人注释更改成你自己并把时间更新</returns>
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
