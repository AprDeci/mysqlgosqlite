package mysqlgosqlite

import (
	_ "embed"
	"os"
	"sync"
)

//go:embed mysql2sqlite.awk
var mysql2sqliteScript []byte

var (
	embeddedScriptPath string
	embeddedScriptOnce sync.Once
	embeddedScriptErr  error
)

func ensureEmbeddedMysql2sqliteExecutable() (string, error) {
	embeddedScriptOnce.Do(func() {
		tmp, err := os.CreateTemp("", "mysql2sqlite-*")
		if err != nil {
			embeddedScriptErr = err
			return
		}

		if _, err = tmp.Write(mysql2sqliteScript); err != nil {
			tmp.Close()
			_ = os.Remove(tmp.Name())
			embeddedScriptErr = err
			return
		}

		if err = tmp.Close(); err != nil {
			_ = os.Remove(tmp.Name())
			embeddedScriptErr = err
			return
		}

		if err = os.Chmod(tmp.Name(), 0o755); err != nil {
			_ = os.Remove(tmp.Name())
			embeddedScriptErr = err
			return
		}

		embeddedScriptPath = tmp.Name()
	})

	return embeddedScriptPath, embeddedScriptErr
}
