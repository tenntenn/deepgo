package toolutil

import (
	"fmt"
	"os"

	"golang.org/x/tools/txtar"
)

func CopyTxtar(dir, txtarStr string) error {
	ar := txtar.Parse([]byte(txtarStr))

	var files []txtar.File
	for _, file := range ar.Files {
		if len(file.Data) > 0 {
			files = append(files, file)
		}
	}
	ar.Files = files

	fsys, err := txtar.FS(ar)
	if err != nil {
		return fmt.Errorf("failed to convert txtar format string to fs.FS: %w", err)
	}

	if err := os.MkdirAll(dir, 0o644); err != nil {
		return fmt.Errorf("failed to create directories %s: %w", dir, err)
	}

	if err := os.CopyFS(dir, fsys); err != nil {
		return fmt.Errorf("failed to copy files from txtar to %s: %w", dir, err)
	}

	return nil
}
