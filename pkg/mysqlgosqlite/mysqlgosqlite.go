package mysqlgosqlite

import (
	"bytes"
	"fmt"
	"io"
	"log"
	"os"
	"os/exec"
	"time"

	"github.com/jarvanstack/mysqldump"
)

// Option allows customizing command paths at runtime.
type Option func(*options)

type options struct {
	mysql2sqlitePath string
	sqlite3Path      string
}

// WithMysql2sqlitePath lets callers override the mysql2sqlite executable path.
func WithMysql2sqlitePath(path string) Option {
	return func(o *options) {
		o.mysql2sqlitePath = path
	}
}

// WithSQLite3Path lets callers override the sqlite3 binary path.
func WithSQLite3Path(path string) Option {
	return func(o *options) {
		o.sqlite3Path = path
	}
}

func applyOptions(opts ...Option) options {
	o := options{
		sqlite3Path: "sqlite3",
	}
	for _, opt := range opts {
		opt(&o)
	}
	return o
}

func resolveMysql2sqlitePath(o options) (string, error) {
	if o.mysql2sqlitePath != "" {
		return o.mysql2sqlitePath, nil
	}
	return ensureEmbeddedMysql2sqliteExecutable()
}

// ExportSQL dumps MySQL data into a SQL file.
func ExportSQL(dsn, outputFile string) error {
	f, err := os.Create(outputFile)
	if err != nil {
		return fmt.Errorf("创建输出文件失败: %w", err)
	}
	defer f.Close()

	if err := mysqldump.Dump(
		dsn,
		mysqldump.WithDropTable(),
		mysqldump.WithData(),
		mysqldump.WithWriter(f),
	); err != nil {
		return fmt.Errorf("mysqldump 执行失败: %w", err)
	}

	if err := f.Sync(); err != nil {
		return fmt.Errorf("写入 SQL 文件失败: %w", err)
	}
	return nil
}

// ConvertDumpToSQLiteFile uses mysql2sqlite to convert a MySQL dump into a SQLite-compatible SQL file.
func ConvertToSQLiteFile(inputFile, outputFile string, opts ...Option) error {
	if err := ensureFileReadable(inputFile); err != nil {
		return err
	}

	o := applyOptions(opts...)
	scriptPath, err := resolveMysql2sqlitePath(o)
	if err != nil {
		return fmt.Errorf("准备 mysql2sqlite 脚本失败: %w", err)
	}

	outFile, err := os.Create(outputFile)
	if err != nil {
		return fmt.Errorf("创建输出文件失败: %w", err)
	}
	defer outFile.Close()

	cmd := exec.Command(scriptPath, inputFile)
	cmd.Stdout = outFile

	var stderr bytes.Buffer
	cmd.Stderr = &stderr

	if err := cmd.Run(); err != nil {
		return fmt.Errorf("转换失败: %w, 错误信息: %s", err, stderr.String())
	}

	if err := outFile.Sync(); err != nil {
		return fmt.Errorf("同步 SQLite SQL 文件失败: %w", err)
	}

	log.Printf("%s mysql2sqlite 转换完成\n", time.Now())
	return nil
}

// ImportDumpToSQLiteDB pipes the converted dump into sqlite3 to generate a SQLite database file.
func ImportToSQLiteDB(inputFile, dbFile string, opts ...Option) error {
	if err := ensureFileReadable(inputFile); err != nil {
		return err
	}

	o := applyOptions(opts...)
	scriptPath, err := resolveMysql2sqlitePath(o)
	if err != nil {
		return fmt.Errorf("准备 mysql2sqlite 脚本失败: %w", err)
	}

	if err := ensureSQLiteFile(dbFile); err != nil {
		return err
	}

	cmd1 := exec.Command(scriptPath, inputFile)
	cmd2 := exec.Command(o.sqlite3Path, dbFile)

	pipe, err := cmd1.StdoutPipe()
	if err != nil {
		return fmt.Errorf("创建管道失败: %w", err)
	}
	cmd2.Stdin = pipe
	cmd2.Stdout = io.Discard

	var stderr bytes.Buffer
	cmd2.Stderr = &stderr

	if err := cmd1.Start(); err != nil {
		return fmt.Errorf("启动 mysql2sqlite 失败: %w", err)
	}

	if err := cmd2.Start(); err != nil {
		return fmt.Errorf("启动 sqlite3 失败: %w", err)
	}

	if err := cmd1.Wait(); err != nil {
		return fmt.Errorf("mysql2sqlite 执行失败: %w", err)
	}

	if err := cmd2.Wait(); err != nil {
		return fmt.Errorf("sqlite3 导入失败: %w, 错误信息: %s", err, stderr.String())
	}

	log.Printf("%s mysql2sqlite 导入完成\n", time.Now())
	return nil
}

func ensureFileReadable(path string) error {
	if _, err := os.Stat(path); err != nil {
		if os.IsNotExist(err) {
			return fmt.Errorf("文件不存在: %s", path)
		}
		return fmt.Errorf("无法访问文件 %s: %w", path, err)
	}
	return nil
}

func ensureSQLiteFile(path string) error {
	if _, err := os.Stat(path); err == nil {
		return nil
	} else if !os.IsNotExist(err) {
		return fmt.Errorf("检查 SQLite 文件失败: %w", err)
	}

	f, err := os.Create(path)
	if err != nil {
		return fmt.Errorf("创建 SQLite 文件失败: %w", err)
	}
	return f.Close()
}
