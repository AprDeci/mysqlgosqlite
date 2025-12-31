# mysqlgosqlite

## 依赖

[github.com/jarvanstack/mysqldump](https://github.com/jarvanstack/mysqldump)

基于 [UN1Q-com/mysql2sqlite: Converts MySQL dump to SQLite3 compatible dump](https://github.com/UN1Q-com/mysql2sqlite) 修改

- 所有没有主键的表但存在 id 字段则 id 为主键
- 调整索引转换规则

## 所需环境

- awk
- sqlite3

## 用法

```go
package main

import "log"
import "github.com/AprDeci/mysqlgosqlite"

func main() {
	const dsn = "root:password@tcp(127.0.0.1:3306)/demo?charset=utf8mb4&parseTime=true&loc=Asia%2FShanghai"

	if err := mysqlgosqlite.ExportSQL(dsn, "dump.sql"); err != nil {
		log.Fatalf("导出失败: %v", err)
	}

	if err := mysqlgosqlite.ConvertToSQLiteFile("dump.sql", "dump_sqlite.sql"); err != nil {
		log.Fatalf("转换失败: %v", err)
	}

	// 如需指定自定义 sqlite3 路径，可传入 Option：
	// mysqlgosqlite.ImportToSQLiteDB("dump.sql", "demo.db", mysqlgosqlite.WithSQLite3Path("/usr/local/bin/sqlite3"))
	if err := mysqlgosqlite.ImportToSQLiteDB("dump.sql", "demo.db"); err != nil {
		log.Fatalf("导入失败: %v", err)
	}
}
```

`ExportSQL` 的第三个可选参数可以传入 `mysqldump.DumpOption`（需额外导入 `github.com/AprDeci/mysqlgosqlite/mysqldump`），例如需要保留 DATETIME/TIMESTAMP 的毫秒：

```go
mysqlgosqlite.ExportSQL(dsn, `dump.sql`, mysqldump.WithTimeFormat(`2006-01-02 15:04:05.000`))
```
