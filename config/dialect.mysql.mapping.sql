-- 数据库方言SQL --

-- 库、表所属者相关SQL --

--  owner create sql start --
CREATE DATABASE [IF NOT EXISTS] ${ownerName}
[[DEFAULT] CHARACTER SET ${characterSetName}]
[[DEFAULT] COLLATE ${collationName}]
--  owner create sql end --

--  owners select sql start --
SELECT
    SCHEMA_NAME name,
    DEFAULT_CHARACTER_SET_NAME characterSetName,
    DEFAULT_COLLATION_NAME collationName
FROM information_schema.schemata
ORDER BY SCHEMA_NAME
--  owners select sql end --

--  owner select sql start --
SELECT
    SCHEMA_NAME name,
    DEFAULT_CHARACTER_SET_NAME characterSetName,
    DEFAULT_COLLATION_NAME collationName
FROM information_schema.schemata
WHERE SCHEMA_NAME='${ownerName}'
--  owner select sql end --

--  owner delete sql start --
DROP DATABASE IF EXISTS ${ownerName}
--  owner delete sql end --

--  owner delete sql start --
DROP DATABASE IF EXISTS ${ownerName}
--  owner delete sql end --

-- 表相关SQL --

--  tables select sql start --
SELECT
    TABLE_NAME name,
    TABLE_COMMENT comment,
    TABLE_SCHEMA ownerName
FROM information_schema.tables
WHERE TABLE_SCHEMA='${ownerName}'
ORDER BY TABLE_NAME
--  tables select sql end --

--  table select sql start --
SELECT
  TABLE_NAME name,
  TABLE_COMMENT comment,
  TABLE_SCHEMA ownerName
FROM information_schema.tables
WHERE TABLE_SCHEMA='${ownerName}'
  AND TABLE_NAME='${tableName}'
--  table select sql end --

--  table create sql start --
CREATE TABLE [${ownerName}.]${tableName}(
${tableCreateColumns}
)
--  table create sql end --

--  table create column sql start --
${columnName} ${columnType} [CHARACTER SET ${characterSetName}] [DEFAULT ${default}] [NOT NULL]
--  table create column sql end --

--  table comment sql start --
ALTER TABLE [${ownerName}.]${tableName} COMMENT '${tableComment}'
--  table comment sql end --

--  table rename sql start --
ALTER TABLE [${ownerName}.]${oldTableName} RENAME AS ${newTableName}
--  table rename sql end --

--  table delete sql start --
DROP TABLE IF EXISTS [${ownerName}.]${tableName}
--  table delete sql end --

-- 字段相关SQL --

--  columns select sql start --
SELECT
    COLUMN_NAME columnName,
    COLUMN_COMMENT columnComment,
    COLUMN_DEFAULT columnDefault,
    EXTRA columnExtra,
    TABLE_NAME tableName,
    TABLE_SCHEMA ownerName,
    CHARACTER_SET_NAME characterSetName,
    IS_NULLABLE isNullable,
    COLUMN_TYPE columnType,
    DATA_TYPE dataType
FROM information_schema.columns
WHERE TABLE_SCHEMA='${ownerName}'
  AND TABLE_NAME='${tableName}'
--  columns select sql end --

--  column select sql start --
SELECT
    COLUMN_NAME columnName,
    COLUMN_COMMENT columnComment,
    COLUMN_DEFAULT columnDefault,
    EXTRA columnExtra,
    TABLE_NAME tableName,
    TABLE_SCHEMA ownerName,
    CHARACTER_SET_NAME characterSetName,
    IS_NULLABLE isNullable,
    COLUMN_TYPE columnType,
    DATA_TYPE dataType
FROM information_schema.columns
WHERE TABLE_SCHEMA='${ownerName}'
  AND TABLE_NAME='${tableName}'
  AND TABLE_NAME='${tableName}'
  AND COLUMN_NAME='${columnName}'
--  column select sql end --

--  column add sql start --
ALTER TABLE [${ownerName}.]${tableName} ADD COLUMN ${columnName} ${columnType} [CHARACTER SET ${characterSetName}] [DEFAULT ${columnDefault}] [NOT NULL] [COMMENT ${columnComment}]
--  column add sql end --

--  column comment sql start --
ALTER TABLE [${ownerName}.]${tableName} CHANGE COLUMN ${columnName} ${columnName} ${columnType} [CHARACTER SET ${characterSetName}] [DEFAULT ${columnDefault}] [NOT NULL] [COMMENT ${columnComment}]
--  column comment sql end --

--  column rename sql start --
ALTER TABLE [${ownerName}.]${tableName} CHANGE COLUMN ${oldColumnName} ${newColumnName} ${columnType} [CHARACTER SET ${characterSetName}] [DEFAULT ${columnDefault}] [NOT NULL] [COMMENT ${columnComment}]
--  column rename sql end --

--  column update sql start --
ALTER TABLE [${ownerName}.]${tableName} CHANGE COLUMN ${columnName} ${columnName} ${columnType} [CHARACTER SET ${characterSetName}] [DEFAULT ${columnDefault}] [NOT NULL] [COMMENT ${columnComment}] [AFTER ${columnAfter}]
--  column update sql end --

--  column delete sql start --
ALTER TABLE [${ownerName}.]${tableName} DROP COLUMN ${columnName}
--  column delete sql end --

-- 主键相关SQL --

--  primary keys select sql start --
SELECT
    key_column_usage.COLUMN_NAME columnName,
    table_constraints.TABLE_NAME tableName,
    table_constraints.TABLE_SCHEMA ownerName
FROM information_schema.table_constraints
JOIN information_schema.key_column_usage USING (CONSTRAINT_NAME,TABLE_SCHEMA,TABLE_NAME)
WHERE table_constraints.TABLE_SCHEMA='${ownerName}'
  AND table_constraints.TABLE_NAME='${tableName}'
  AND table_constraints.CONSTRAINT_TYPE='PRIMARY KEY'
--  primary keys select sql end --

--  primary key add sql start --
ALTER TABLE [${ownerName}.]${tableName} ADD PRIMARY KEY (${columnNames})
--  primary key add sql end --

--  primary key delete sql start --
ALTER TABLE [${ownerName}.]${tableName} DROP PRIMARY KEY
--  primary key delete sql end --

-- 索引相关SQL --

--  indexes select sql start --
SELECT
    INDEX_NAME indexName,
    COLUMN_NAME columnName,
    INDEX_COMMENT indexComment,
    NON_UNIQUE nonUnique,
    TABLE_NAME tableName,
    TABLE_SCHEMA ownerName
FROM information_schema.statistics
WHERE TABLE_SCHEMA='${ownerName}'
  AND TABLE_NAME='${tableName}'
  AND TABLE_NAME='${tableName}'
  AND INDEX_NAME NOT IN(
    SELECT table_constraints.CONSTRAINT_NAME
    FROM information_schema.table_constraints
    JOIN information_schema.key_column_usage USING (CONSTRAINT_NAME,TABLE_SCHEMA,TABLE_NAME)
    WHERE table_constraints.TABLE_SCHEMA='${ownerName}'
      AND table_constraints.TABLE_NAME='${tableName}'
      AND table_constraints.CONSTRAINT_TYPE='PRIMARY KEY'
)

--  indexes select sql end --

--  index add sql start --
ALTER TABLE [${ownerName}.]${tableName} ADD [PRIMARY KEY | UNIQUE | FULLTEXT | INDEX] ${indexName} (${columnNames}) [COMMENT ${columnComment}]
--  index add sql end --

--  index delete sql start --
ALTER TABLE [${ownerName}.]${tableName} DROP INDEX
--  index delete sql end --
