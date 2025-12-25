package main

import (
	"log"
	"os"

	"mysqlgosqlite/pkg/mysqlgosqlite"
)

func main() {
	if err := os.MkdirAll("./localtest", 0o755); err != nil {
		log.Fatalf("创建测试目录失败: %v", err)
	}

	dsn := "root:Luchen1122@tcp(localhost:3306)/test?charset=utf8mb4&parseTime=true&loc=Asia%2FShanghai"
	exportSQLFile := "./localtest/dump.sql"
	sqliteSQLFile := "./localtest/dump_sqlite.sql"
	sqliteDBFile := "./localtest/testDB.db"

	if err := mysqlgosqlite.ExportSQL(dsn, exportSQLFile); err != nil {
		log.Fatalf("导出 MySQL 失败: %v", err)
	}

	if err := mysqlgosqlite.ConvertToSQLiteFile(exportSQLFile, sqliteSQLFile); err != nil {
		log.Fatalf("转换为 SQLite SQL 失败: %v", err)
	}

	if err := mysqlgosqlite.ImportToSQLiteDB(exportSQLFile, sqliteDBFile); err != nil {
		log.Fatalf("导入 SQLite DB 失败: %v", err)
	}

}
