# go-dialect

SQL方言

## 参数说明

#### 全局参数说明
* -do
  * 操作：export(导出)、import(导入)、sync(同步)
* -sourceDialect
  * 源 数据库 方言 mysql、sqlite3、dm、kingbase、oracle、shentong
* -sourceHost
    * 源 数据库 host
* -sourcePort
    * 源 数据库 port
* -sourceUser
    * 源 数据库 user
* -sourcePassword
    * 源 数据库 password
* -sourceDatabase
    * 源 数据库 连接库（库名、用户名、SID）
* -fileType
    * 文件 类型：sql、excel、txt、csv(sql将导出单个文件，其它每个表导出一个文件)
* -fileDialect
    * 文件 数据库 方言
* -skipOwner
  * 忽略库名，多个使用“,”隔开
  * 示例：d1,d2
* -skipTable
    * 忽略表名，多个使用“,”隔开
    * 示例：t1,t2

#### 导出参数说明
* -exportDialect
    * 导出 数据库 方言
* -exportDir
    * 导出 文件存储目录
* -exportOwner
    * 导出 库（库名、表拥有者），默认全部，多个使用“,”隔开
    * 如果导出重置 库名，使用“库=修改后库名”
    * 示例：x,xx=xx1,xxx=xxx1
    * 如果是sql类型，导出单独一个文件，${ownerName}.sql
    * 其它文件每个表一个文件，${ownerName}/${tableName}.sql
* -exportTable
  * 导出 表，默认全部，多个使用“,”隔开
  * 如果导出重置 表名，使用“表=修改后表名”
  * 示例：x,xx=xx1,xxx=xxx1
* -exportStruct
    * 导出 结构体，默认true，适用于导出类型为sql、excel
* -exportData
    * 导出 数据，默认true
* -exportAppendOwner
  * sql文件类型的sql拼接 连接库（库名、用户名），拼接原库名或重命名后的库名，默认false

## 导入参数说明
* -importOwner
  * 导入 库（库名、表拥有者），并指定文件路径，多个使用“,”隔开
  * 如果是sql类型，指定到sql文件
  * 其它类型，指定到文件目录，以表名为文件文件名
  * 示例：db1=data/db1,db2=data/db2.sql
* -importOwnerCreateIfNotExist
  * 导入 库如果不存在，则创建，默认false
* -importOwnerCreatePassword
  * 导入 库创建的密码，只有库为所属者有效，默认为sourcePassword，如：oracle等数据库
    
## 同步参数说明
* -targetDialect
  * 目标 数据库 方言 mysql、sqlite3、dm、kingbase、oracle、shentong
* -targetHost
  * 目标 数据库 host
* -targetPort
  * 目标 数据库 port
* -targetUser
  * 目标 数据库 user
* -targetPassword
  * 目标 数据库 password
* -targetDatabase
  * 目标 数据库 连接库（库名、用户名、SID）
* -syncOwner
  * 同步 库（库名、表拥有者），默认全部，多个使用“,”隔开
  * 如果同步重置 库名，使用“库=修改后库名”
  * 示例：x,xx=xx1,xxx=xxx1
* -syncStruct
  * 同步 结构体，默认true，适用于导出类型为sql、excel
* -syncData
  * 同步 数据，默认true
* -syncOwnerCreateIfNotExist
  * 同步 库如果不存在，则创建，默认false
* -syncOwnerCreatePassword
  * 同步 库创建的密码，只有库为所属者有效，默认为targetPassword，如：oracle等数据库


```shell

docker run -itd --name mysql-3306 -m 1024m -p 3306:3306 -e MYSQL_ROOT_PASSWORD=123456 mysql:5.7

# 导出 mysql 数据库的 mysql,information_schema,performance_schema,sys 库 为 mysql 的 sql
go run . -do export -sourceDialect mysql -sourceHost 127.0.0.1 -sourcePort 3306 -sourceUser root -sourcePassword 123456 -fileType sql -exportDir temp/export/test -exportOwner mysql,information_schema,performance_schema,sys -exportDialect mysql

# 导出 mysql 数据库的 mysql,information_schema,performance_schema,sys 库 为 sqlite 的 sql
go run . -do export -sourceDialect mysql -sourceHost 127.0.0.1 -sourcePort 3306 -sourceUser root -sourcePassword 123456 -fileType sql -exportDir temp/export/test -exportOwner mysql,information_schema,performance_schema,sys -exportDialect sqlite

# 导入 mysql 的 sql 到 mysql 的 DB1 库
go run . -do import -sourceDialect mysql -sourceHost 127.0.0.1 -sourcePort 3306 -sourceUser root -sourcePassword 123456 -fileType sql -importOwner VRV_JOB=temp/VRV_JOB.sql  -importOwnerCreateIfNotExist true


# 导入 sqlite 的 sql 到 mysql 的 DB1 库
go run . -do import -sourceDialect mysql -sourceHost 127.0.0.1 -sourcePort 3306 -sourceUser root -sourcePassword 123456 -fileType sql -importOwner DB1=temp/export/test/mysql.sql  -importOwnerCreateIfNotExist true


# 导出 mysql 数据库的 mysql 库 为 sqlite 的 sql
go run . -do export -sourceDialect mysql -sourceHost 127.0.0.1 -sourcePort 3306 -sourceUser root -sourcePassword 123456 -fileType sql -exportDir temp/export/test -exportOwner mysql=sqlite -exportDialect sqlite

# 导入 sqlite 的 sql 到 sqlite
go run . -do import -sourceDialect sqlite -sourceDatabase temp/sqlite.db -fileType sql -importOwner main=temp/export/test/sqlite.sql

# 导出 sqlite 的 main 到 sqlite
go run . -do export -sourceDialect sqlite -sourceDatabase temp/sqlite.mysql -fileType sql -exportDir temp/export/test -exportOwner main=main -exportDialect sqlite

# 同步 mysql 数据库的 mysql 库 到 mysql 数据库的 DB2
go run . -do sync -sourceDialect mysql -sourceHost 127.0.0.1 -sourcePort 3306 -sourceUser root -sourcePassword 123456 -targetDialect mysql -targetHost 127.0.0.1 -targetPort 3306 -targetUser root -targetPassword 123456 -syncOwner mysql=DB2 -syncOwnerCreateIfNotExist true
go run . -do sync -sourceDialect mysql -sourceHost 127.0.0.1 -sourcePort 3306 -sourceUser root -sourcePassword 123456 -targetDialect mysql -targetHost 127.0.0.1 -targetPort 3306 -targetUser root -targetPassword 123456 -syncOwner information_schema=DB3 -syncOwnerCreateIfNotExist true
go run . -do sync -sourceDialect mysql -sourceHost 127.0.0.1 -sourcePort 3306 -sourceUser root -sourcePassword 123456 -targetDialect mysql -targetHost 127.0.0.1 -targetPort 3306 -targetUser root -targetPassword 123456 -syncOwner performance_schema=DB4 -syncOwnerCreateIfNotExist true
go run . -do sync -sourceDialect mysql -sourceHost 127.0.0.1 -sourcePort 3306 -sourceUser root -sourcePassword 123456 -targetDialect mysql -targetHost 127.0.0.1 -targetPort 3306 -targetUser root -targetPassword 123456 -syncOwner sys=DB5 -syncOwnerCreateIfNotExist true

# 同步 mysql 数据库的 mysql 库 到 sqlite
go run . -do sync -sourceDialect mysql -sourceHost 127.0.0.1 -sourcePort 3306 -sourceUser root -sourcePassword 123456 -targetDialect sqlite -targetDatabase temp/sqlite.mysql -syncOwner mysql=main
go run . -do sync -sourceDialect mysql -sourceHost 127.0.0.1 -sourcePort 3306 -sourceUser root -sourcePassword 123456 -targetDialect sqlite -targetDatabase temp/sqlite.mysql -syncOwner information_schema=main
go run . -do sync -sourceDialect mysql -sourceHost 127.0.0.1 -sourcePort 3306 -sourceUser root -sourcePassword 123456 -targetDialect sqlite -targetDatabase temp/sqlite.mysql -syncOwner performance_schema=main
go run . -do sync -sourceDialect mysql -sourceHost 127.0.0.1 -sourcePort 3306 -sourceUser root -sourcePassword 123456 -targetDialect sqlite -targetDatabase temp/sqlite.mysql -syncOwner sys=main


docker run -itd --name oracle-1521 -p 1521:1521 teamide/oracle-xe-11g:1.0

# 同步 mysql 数据库的 mysql 库 到 oracle DB_FROM_MYSQL
go run . -do sync -sourceDialect mysql -sourceHost 127.0.0.1 -sourcePort 3306 -sourceUser root -sourcePassword 123456 -targetDialect oracle -targetHost 127.0.0.1 -targetPort 1521 -targetUser root -targetPassword 123456 -targetDatabase xe -syncOwner mysql=DB_FROM_MYSQL -syncOwnerCreateIfNotExist true

```