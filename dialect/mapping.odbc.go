package dialect

func NewMappingOdbc() (mapping *SqlMapping) {
	mapping = &SqlMapping{
		dialectType: TypeOdbc,

		OwnerNamePackChar:  "\"",
		TableNamePackChar:  "\"",
		ColumnNamePackChar: "\"",
		SqlValuePackChar:   "'",
		SqlValueEscapeChar: "'",
	}

	for _, one := range postgresqlColumnTypeList {
		mapping.AddColumnTypeInfo(one)
	}

	for _, one := range postgresqlIndexTypeList {
		mapping.AddIndexTypeInfo(one)
	}

	return
}

var (
	odbcIndexTypeList []*IndexTypeInfo
)

func appendOdbcIndexType(indexType *IndexTypeInfo) {
	odbcIndexTypeList = append(odbcIndexTypeList, indexType)
}

func init() {

}
