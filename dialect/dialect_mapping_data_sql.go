package dialect

import (
	"errors"
	"strings"
)

func (this_ *mappingDialect) DataListInsertSql(param *ParamModel, ownerName string, tableName string, columnList []*ColumnModel, dataList []map[string]interface{}) (sqlList []string, valuesList [][]interface{}, err error) {
	if len(dataList) == 0 {
		return
	}
	var keys []string
	for _, column := range columnList {
		if column.PrimaryKey {
			keys = append(keys, column.ColumnName)
		}
	}

	for _, data := range dataList {

		var values []interface{}
		insertColumns := ""
		insertValues := ""
		for _, column := range columnList {
			value, valueOk := data[column.ColumnName]
			if !valueOk {
				continue
			}

			insertColumns += this_.ColumnNamePack(param, column.ColumnName) + ", "
			this_.AppendSqlValue(param, &insertValues, column, value, &values)
			insertValues += ", "
		}
		insertColumns = strings.TrimSuffix(insertColumns, ", ")
		insertValues = strings.TrimSuffix(insertValues, ", ")

		sql := "INSERT INTO "

		if ownerName != "" {
			sql += this_.OwnerNamePack(param, ownerName) + "."
		}
		sql += this_.TableNamePack(param, tableName)
		if insertColumns != "" {
			sql += "(" + insertColumns + ")"
		}
		if insertValues != "" {
			sql += " VALUES (" + insertValues + ")"
		}

		sqlList = append(sqlList, sql)
		valuesList = append(valuesList, values)
	}
	return
}

func (this_ *mappingDialect) DataListUpdateSql(param *ParamModel, ownerName string, tableName string, columnList []*ColumnModel, dataList []map[string]interface{}, dataWhereList []map[string]interface{}) (sqlList []string, valuesList [][]interface{}, err error) {
	if len(dataList) == 0 {
		return
	}
	if len(dataList) != len(dataWhereList) {
		err = errors.New("更新数据与更新条件数量不一致")
		return
	}
	var keyColumnList []*ColumnModel
	for _, column := range columnList {
		if column.PrimaryKey {
			keyColumnList = append(keyColumnList, column)
		}
	}

	for index, data := range dataList {
		dataWhere := dataWhereList[index]
		if len(dataWhere) == 0 {
			err = errors.New("更新数据条件丢失")
			return
		}

		sql := "UPDATE "
		var values []interface{}

		if ownerName != "" {
			sql += this_.OwnerNamePack(param, ownerName) + "."
		}
		sql += this_.TableNamePack(param, tableName)
		sql += " SET "

		for _, column := range columnList {
			value, valueOK := data[column.ColumnName]
			if !valueOK {
				continue
			}

			sql += "" + this_.ColumnNamePack(param, column.ColumnName) + "="
			this_.AppendSqlValue(param, &sql, column, value, &values)
			sql += ", "
		}
		sql = strings.TrimSuffix(sql, ", ")

		sql += " WHERE "
		whereColumnList := keyColumnList
		if len(keyColumnList) == 0 {
			whereColumnList = columnList
		} else {
			for _, column := range whereColumnList {
				sql += "" + this_.ColumnNamePack(param, column.ColumnName) + "="
				this_.AppendSqlValue(param, &sql, column, dataWhere[column.ColumnName], &values)
				sql += " AND "
			}
		}
		sql = strings.TrimSuffix(sql, " AND ")

		sqlList = append(sqlList, sql)
		valuesList = append(valuesList, values)
	}
	return
}

func (this_ *mappingDialect) AppendSqlValue(param *ParamModel, sqlInfo *string, column *ColumnModel, value interface{}, args *[]interface{}) {
	if param != nil && param.AppendSqlValue != nil && *param.AppendSqlValue {
		*sqlInfo += this_.SqlValuePack(param, column, value)
	} else {
		*sqlInfo += "?"
		*args = append(*args, value)
	}
}
func (this_ *mappingDialect) DataListDeleteSql(param *ParamModel, ownerName string, tableName string, columnList []*ColumnModel, dataWhereList []map[string]interface{}) (sqlList []string, valuesList [][]interface{}, err error) {
	if len(dataWhereList) == 0 {
		return
	}
	var keyColumnList []*ColumnModel
	for _, column := range columnList {
		if column.PrimaryKey {
			keyColumnList = append(keyColumnList, column)
		}
	}

	for _, dataWhere := range dataWhereList {
		if len(dataWhere) == 0 {
			err = errors.New("更新数据条件丢失")
			return
		}

		sql := "DELETE FROM "
		var values []interface{}

		if ownerName != "" {
			sql += this_.OwnerNamePack(param, ownerName) + "."
		}
		sql += this_.TableNamePack(param, tableName)

		sql += " WHERE "
		whereColumnList := keyColumnList
		if len(keyColumnList) == 0 {
			whereColumnList = columnList
		} else {
			for _, column := range whereColumnList {
				sql += "" + this_.ColumnNamePack(param, column.ColumnName) + "="
				this_.AppendSqlValue(param, &sql, column, dataWhere[column.ColumnName], &values)
				sql += " AND "
			}
		}
		sql = strings.TrimSuffix(sql, " AND ")

		sqlList = append(sqlList, sql)
		valuesList = append(valuesList, values)
	}
	return
}

type Where struct {
	Name                    string `json:"name"`
	Value                   string `json:"value"`
	Before                  string `json:"before"`
	After                   string `json:"after"`
	CustomSql               string `json:"customSql"`
	SqlConditionalOperation string `json:"sqlConditionalOperation"`
	AndOr                   string `json:"andOr"`
}

type Order struct {
	Name    string `json:"name"`
	AscDesc string `json:"ascDesc"`
}

func (this_ *mappingDialect) DataListSelectSql(param *ParamModel, ownerName string, tableName string, columnList []*ColumnModel, whereList []*Where, orderList []*Order) (sql string, values []interface{}, err error) {
	selectColumns := ""
	for _, column := range columnList {
		selectColumns += this_.ColumnNamePack(param, column.ColumnName) + ","
	}
	selectColumns = strings.TrimSuffix(selectColumns, ",")
	if selectColumns == "" {
		selectColumns = "*"
	}
	sql = "SELECT " + selectColumns + " FROM "

	if ownerName != "" {
		sql += this_.OwnerNamePack(param, ownerName) + "."
	}
	sql += this_.TableNamePack(param, tableName)

	//构造查询用的finder
	if len(whereList) > 0 {
		sql += " WHERE"
		for index, where := range whereList {
			sql += " " + this_.ColumnNamePack(param, where.Name)
			value := where.Value
			switch where.SqlConditionalOperation {
			case "like":
				sql += " LIKE ?"
				values = append(values, "%"+value+"%")
			case "not like":
				sql += " NOT LIKE ?"
				values = append(values, "%"+value+"%")
			case "like start":
				sql += " LIKE ?"
				values = append(values, ""+value+"%")
			case "not like start":
				sql += " NOT LIKE ?"
				values = append(values, ""+value+"%")
			case "like end":
				sql += " LIKE ?"
				values = append(values, "%"+value+"")
			case "not like end":
				sql += " NOT LIKE ?"
				values = append(values, "%"+value+"")
			case "is null":
				sql += " IS NULL"
			case "is not null":
				sql += " IS NOT NULL"
			case "is empty":
				sql += " = ?"
				values = append(values, "")
			case "is not empty":
				sql += " <> ?"
				values = append(values, "")
			case "between":
				sql += " BETWEEN ? AND ?"
				values = append(values, where.Before, where.After)
			case "not between":
				sql += " NOT BETWEEN ? AND ?"
				values = append(values, where.Before, where.After)
			case "in":
				sql += " IN (?)"
				values = append(values, value)
			case "not in":
				sql += " NOT IN (?)"
				values = append(values, value)
			default:
				sql += " " + where.SqlConditionalOperation + " ?"
				values = append(values, value)
			}
			// params_ = append(params_, where.Value)
			if index < len(whereList)-1 {
				sql += " " + where.AndOr + " "
			}
		}
	}
	if len(orderList) > 0 {
		sql += " ORDER BY"
		for index, order := range orderList {
			sql += " " + this_.ColumnNamePack(param, order.Name)
			if order.AscDesc != "" {
				sql += " " + order.AscDesc
			}
			// params_ = append(params_, where.Value)
			if index < len(orderList)-1 {
				sql += ","
			}
		}

	}
	return
}
