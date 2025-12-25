package main

import (
	_ "embed"
	"os"
	"sync"
)

//go:embed mysql2sqlite
var mysql2sqliteScript []byte

var (
	mysql2sqlitePath string
	mysql2sqliteOnce sync.Once
	mysql2sqliteErr  error
)

func ensureMysql2sqliteExecutable() (string, error) {
	mysql2sqliteOnce.Do(func() {
		tmp, err := os.CreateTemp("", "mysql2sqlite-*")
		if err != nil {
			mysql2sqliteErr = err
			return
		}

		if _, err = tmp.Write(mysql2sqliteScript); err != nil {
			tmp.Close()
			_ = os.Remove(tmp.Name())
			mysql2sqliteErr = err
			return
		}

		if err = tmp.Close(); err != nil {
			_ = os.Remove(tmp.Name())
			mysql2sqliteErr = err
			return
		}

		if err = os.Chmod(tmp.Name(), 0o755); err != nil {
			_ = os.Remove(tmp.Name())
			mysql2sqliteErr = err
			return
		}

		mysql2sqlitePath = tmp.Name()
	})

	return mysql2sqlitePath, mysql2sqliteErr
}
