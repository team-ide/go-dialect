package main

import (
	"fmt"
	"github.com/team-ide/go-dialect/dialect"
	"github.com/team-ide/go-dialect/worker"
	"os"
	"testing"
)

func loadSql(name string) (srcSql string) {
	bs, err := os.ReadFile(name)
	if err != nil {
		panic(err)
	}
	srcSql = string(bs)
	return
}

func saveSql(destSql string, name string) {
	err := os.WriteFile(name, []byte(destSql), 0777)
	if err != nil {
		panic(err)
	}
	return
}

func TestSqlParse(t *testing.T) {
	var err error
	var convertParser *worker.ConvertParser

	srcSql := loadSql(`temp/sql_test.sql`)

	convertParser = worker.NewConvertParser(srcSql, dialect.Mysql)
	err = convertParser.Parse()
	if err != nil {
		panic(err)
	}
	saveSql(convertParser.GetDestSql(), "temp/sql_mysql.sql")

	convertParser = worker.NewConvertParser(srcSql, dialect.Oracle)
	err = convertParser.Parse()
	if err != nil {
		panic(err)
	}
	saveSql(convertParser.GetDestSql(), "temp/sql_oracle.sql")

	convertParser = worker.NewConvertParser(srcSql, dialect.ShenTong)
	err = convertParser.Parse()
	if err != nil {
		panic(err)
	}
	saveSql(convertParser.GetDestSql(), "temp/sql_shentong.sql")

	convertParser = worker.NewConvertParser(srcSql, dialect.KinBase)
	err = convertParser.Parse()
	if err != nil {
		panic(err)
	}
	saveSql(convertParser.GetDestSql(), "temp/sql_kinbase.sql")

	convertParser = worker.NewConvertParser(srcSql, dialect.DaMen)
	err = convertParser.Parse()
	if err != nil {
		panic(err)
	}
	saveSql(convertParser.GetDestSql(), "temp/sql_damen.sql")

	convertParser = worker.NewConvertParser(srcSql, dialect.Sqlite)
	err = convertParser.Parse()
	if err != nil {
		panic(err)
	}
	saveSql(convertParser.GetDestSql(), "temp/sql_sqlite.sql")

}

func TestSqlSplit(t *testing.T) {
	//var err error
	var sqlInfo = `
(238, 'MATCH AGAINST', 16, 'Syntax:
MATCH (col1,col2,...) AGAINST (expr [search_modifier])

MySQL has support for full-text indexing and searching:

o A full-text index in MySQL is an index of type FULLTEXT.

o Full-text indexes can be used only with InnoDB or MyISAM tables, and
  can be created only for CHAR, VARCHAR, or TEXT columns.

o MySQL provides a built-in full-text ngram parser that supports
  Chinese, Japanese, and Korean (CJK), and an installable MeCab
  full-text parser plugin for Japanese. Parsing differences are
  outlined in
  https://dev.mysql.com/doc/refman/5.7/en/fulltext-search-ngram.html,
  and
  https://dev.mysql.com/doc/refman/5.7/en/fulltext-search-mecab.html.

o A FULLTEXT index definition can be given in the CREATE TABLE
  statement when a table is created, or added later using ALTER TABLE
  or CREATE INDEX.

o For large data sets, it is much faster to load your data into a table
  that has no FULLTEXT index and then create the index after that, than
  to load data into a table that has an existing FULLTEXT index.

Full-text searching is performed using MATCH() AGAINST syntax. MATCH()
takes a comma-separated list that names the columns to be searched.
AGAINST takes a string to search for, and an optional modifier that
indicates what type of search to perform. The search string must be a
string value that is constant during query evaluation. This rules out,
for example, a table column because that can differ for each row.

There are three types of full-text searches:

o A natural language search interprets the search string as a phrase in
  natural human language (a phrase in free text). There are no special
  operators, with the exception of double quote (") characters. The
  stopword list applies. For more information about stopword lists, see
  https://dev.mysql.com/doc/refman/5.7/en/fulltext-stopwords.html.

  Full-text searches are natural language searches if the IN NATURAL
  LANGUAGE MODE modifier is given or if no modifier is given. For more
  information, see
  https://dev.mysql.com/doc/refman/5.7/en/fulltext-natural-language.htm
  l.

o A boolean search interprets the search string using the rules of a
  special query language. The string contains the words to search for.
  It can also contain operators that specify requirements such that a
  word must be present or absent in matching rows, or that it should be
  weighted higher or lower than usual. Certain common words (stopwords)
  are omitted from the search index and do not match if present in the
  search string. The IN BOOLEAN MODE modifier specifies a boolean
  search. For more information, see
  https://dev.mysql.com/doc/refman/5.7/en/fulltext-boolean.html.

o A query expansion search is a modification of a natural language
  search. The search string is used to perform a natural language
  search. Then words from the most relevant rows returned by the search
  are added to the search string and the search is done again. The
  query returns the rows from the second search. The IN NATURAL
  LANGUAGE MODE WITH QUERY EXPANSION or WITH QUERY EXPANSION modifier
  specifies a query expansion search. For more information, see
  https://dev.mysql.com/doc/refman/5.7/en/fulltext-query-expansion.html
  .

URL: https://dev.mysql.com/doc/refman/5.7/en/fulltext-search.html

', 'mysql> SELECT id, body, MATCH (title,body) AGAINST
    (''Security implications of running MySQL as root''
    IN NATURAL LANGUAGE MODE) AS score
    FROM articles WHERE MATCH (title,body) AGAINST
    (''Security implications of running MySQL as root''
    IN NATURAL LANGUAGE MODE);
+----+-------------------------------------+-----------------+
| id | body                                | score           |
+----+-------------------------------------+-----------------+
|  4 | 1. Never run mysqld as root. 2. ... | 1.5219271183014 |
|  6 | When configured properly, MySQL ... | 1.3114095926285 |
+----+-------------------------------------+-----------------+
2 rows in set (0.00 sec)
', 'https://dev.mysql.com/doc/refman/5.7/en/fulltext-search.html'),;
`
	fmt.Println(dialect.Mysql.IsSqlEnd(sqlInfo))

}
