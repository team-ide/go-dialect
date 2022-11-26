package dialect

import "strings"

var (
	postgresqlColumnTypeList []*ColumnTypeInfo
)

func appendPostgresqlColumnType(columnType *ColumnTypeInfo) {
	postgresqlColumnTypeList = append(postgresqlColumnTypeList, columnType)
}
func init() {
	appendPostgresqlColumnType(&ColumnTypeInfo{Name: "SMALLINT", Format: "SMALLINT", IsNumber: true})
	appendPostgresqlColumnType(&ColumnTypeInfo{Name: "INTEGER", Format: "INTEGER", IsNumber: true})
	appendPostgresqlColumnType(&ColumnTypeInfo{Name: "BIGINT", Format: "BIGINT", IsNumber: true})
	appendPostgresqlColumnType(&ColumnTypeInfo{Name: "DECIMAL", Format: "DECIMAL", IsNumber: true})
	appendPostgresqlColumnType(&ColumnTypeInfo{Name: "NUMERIC", Format: "NUMERIC", IsNumber: true})
	appendPostgresqlColumnType(&ColumnTypeInfo{Name: "REAL", Format: "REAL", IsNumber: true})
	appendPostgresqlColumnType(&ColumnTypeInfo{Name: "DOUBLE PRECISION", Format: "DOUBLE PRECISION", IsNumber: true})
	appendPostgresqlColumnType(&ColumnTypeInfo{Name: "SMALLSERIAL", Format: "SMALLSERIAL", IsNumber: true})
	appendPostgresqlColumnType(&ColumnTypeInfo{Name: "SERIAL", Format: "SERIAL", IsNumber: true})
	appendPostgresqlColumnType(&ColumnTypeInfo{Name: "BIGSERIAL", Format: "BIGSERIAL", IsNumber: true})

	appendPostgresqlColumnType(&ColumnTypeInfo{Name: "MONEY", Format: "MONEY", IsNumber: true})

	appendPostgresqlColumnType(&ColumnTypeInfo{Name: "CHARACTER VARYING", Format: "CHARACTER VARYING($l)", IsString: true})
	appendPostgresqlColumnType(&ColumnTypeInfo{Name: "VARCHAR", Format: "VARCHAR($l)", IsString: true})
	appendPostgresqlColumnType(&ColumnTypeInfo{Name: "CHARACTER", Format: "CHARACTER($l)", IsString: true})
	appendPostgresqlColumnType(&ColumnTypeInfo{Name: "CHAR", Format: "CHAR($l)", IsString: true})
	appendPostgresqlColumnType(&ColumnTypeInfo{Name: "TEXT", Format: "TEXT($l)", IsString: true})

	appendPostgresqlColumnType(&ColumnTypeInfo{Name: "TIMESTAMP", Format: "TIMESTAMP", IsDateTime: true})
	appendPostgresqlColumnType(&ColumnTypeInfo{Name: "TIMESTAMP WITHOUT TIME ZONE", Format: "TIMESTAMP WITHOUT TIME ZONE", IsDateTime: true})
	appendPostgresqlColumnType(&ColumnTypeInfo{Name: "TIMESTAMP WITH TIME ZONE", Format: "TIMESTAMP WITH TIME ZONE", IsDateTime: true})
	appendPostgresqlColumnType(&ColumnTypeInfo{Name: "DATE", Format: "DATE", IsDateTime: true})
	appendPostgresqlColumnType(&ColumnTypeInfo{Name: "TIME WITHOUT TIME ZONE", Format: "TIME WITHOUT TIME ZONE", IsDateTime: true})
	appendPostgresqlColumnType(&ColumnTypeInfo{Name: "TIME WITH TIME ZONE", Format: "TIME WITH TIME ZONE", IsDateTime: true})
	appendPostgresqlColumnType(&ColumnTypeInfo{Name: "INTERVAL", Format: "INTERVAL", IsDateTime: true})

	appendPostgresqlColumnType(&ColumnTypeInfo{Name: "BOOLEAN", Format: "BOOLEAN", IsString: true})
	appendPostgresqlColumnType(&ColumnTypeInfo{Name: "ENUM", IsString: true, IsEnum: true,
		FullColumnByColumnType: func(columnType string, column *ColumnModel) (err error) {
			if strings.Contains(columnType, "(") {
				setStr := columnType[strings.Index(columnType, "(")+1 : strings.Index(columnType, ")")]
				setStr = strings.ReplaceAll(setStr, "'", "")
				column.ColumnEnums = strings.Split(setStr, ",")
			}
			return
		},
	})

}

var (
	postgresqlIndexTypeList []*IndexTypeInfo
)

func appendPostgresqlIndexType(indexType *IndexTypeInfo) {
	postgresqlIndexTypeList = append(postgresqlIndexTypeList, indexType)
}

func init() {

}
