package worker

import "github.com/team-ide/go-dialect/vitess/sqlparser"

func FormatSqlList(sqlInfo string) (sqlList []string, err error) {
	stmts, err := GetStatements(sqlInfo)
	if err != nil {
		return
	}

	for _, stmt := range stmts {

		buf := sqlparser.NewTrackedBuffer(nil)
		stmt.Format(buf)
		sqlList = append(sqlList, buf.String())
	}
	return
}
