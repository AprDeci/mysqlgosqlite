package main

import (
	"bytes"
	"fmt"
	"log"
	"os"
	"os/exec"

	"github.com/jarvanstack/mysqldump"
)

func main() {

	dsn := "root:Luchen1122@tcp(localhost:3306)/test?charset=utf8mb4&parseTime=true&loc=Asia%2FShanghai"
	exportSqlFile := "dump.sql"

	exportSql(dsn, exportSqlFile)
	convert2sqliteFile(exportSqlFile, "dump_sqlite.sql")
	Convert2Sqlite(exportSqlFile, "testDB.db")

}

func exportSql(dsn, outputFile string) {

	f, _ := os.Create(outputFile)

	_ = mysqldump.Dump(
		dsn,                       // DSN
		mysqldump.WithDropTable(), // Option: Delete table before create (Default: Not delete table)
		mysqldump.WithData(),      // Option: Dump Data (Default: Only dump table schema)
		mysqldump.WithWriter(f),   // Option: Writer (Default: os.Stdout)
	)
}

func convert2sqliteFile(inputFile, outputFile string) error {
	if _, err := os.Stat(inputFile); err != nil {
		return fmt.Errorf("输入文件不存在: %s", inputFile)
	}

	scriptPath, err := ensureMysql2sqliteExecutable()
	if err != nil {
		return fmt.Errorf("准备mysql2sqlite脚本失败: %w", err)
	}

	// 创建输出文件
	outFile, err := os.Create(outputFile)
	if err != nil {
		return fmt.Errorf("创建输出文件失败: %v", err)
	}
	defer outFile.Close()

	// 执行mysql2sqlite命令
	cmd := exec.Command(scriptPath, inputFile)
	cmd.Stdout = outFile

	// 捕获标准错误
	var stderr bytes.Buffer
	cmd.Stderr = &stderr

	// 执行命令
	if err := cmd.Run(); err != nil {
		return fmt.Errorf("转换失败: %v, 错误信息: %s", err, stderr.String())
	}

	log.Printf("成功将 %s 转换为 SQLite 格式: %s", inputFile, outputFile)
	return nil
}

func Convert2Sqlite(inputFile, dbFile string) error {
	// 检查输入文件是否存在
	if _, err := os.Stat(inputFile); os.IsNotExist(err) {
		return fmt.Errorf("输入文件不存在: %s", inputFile)
	}

	scriptPath, err := ensureMysql2sqliteExecutable()
	if err != nil {
		return fmt.Errorf("准备mysql2sqlite脚本失败: %w", err)
	}

	if _, err := os.Stat(dbFile); os.IsNotExist(err) {
		return fmt.Errorf("数据库文件不存在: %s", dbFile)
	}

	// 创建第一个命令：mysql2sqlite
	cmd1 := exec.Command(scriptPath, inputFile)

	// 创建第二个命令：sqlite3
	cmd2 := exec.Command("sqlite3", dbFile)

	// 创建管道连接两个命令
	pipe, err := cmd1.StdoutPipe()
	if err != nil {
		return fmt.Errorf("创建管道失败: %v", err)
	}
	cmd2.Stdin = pipe

	// 捕获第二个命令的输出和错误
	var output bytes.Buffer
	var stderr bytes.Buffer
	cmd2.Stdout = &output
	cmd2.Stderr = &stderr

	// 启动第一个命令
	if err := cmd1.Start(); err != nil {
		return fmt.Errorf("启动mysql2sqlite失败: %v", err)
	}

	// 启动第二个命令
	if err := cmd2.Start(); err != nil {
		return fmt.Errorf("启动sqlite3失败: %v", err)
	}

	// 等待第一个命令完成
	if err := cmd1.Wait(); err != nil {
		return fmt.Errorf("mysql2sqlite执行失败: %v", err)
	}

	// 等待第二个命令完成
	if err := cmd2.Wait(); err != nil {
		return fmt.Errorf("sqlite3导入失败: %v, 错误信息: %s", err, stderr.String())
	}

	log.Printf("成功将 %s 导入到 SQLite 数据库: %s", inputFile, dbFile)
	return nil
}
