package helm

import (
	"os"
	"path/filepath"
	"strings"

	"github.com/otiai10/copy"
)

func (h *Helm) sync(src, dest string) error {
	if err := os.MkdirAll(dest, os.FileMode(0700)); err != nil {
		return err
	}

	err := filepath.Walk(src, func(path string, sfi os.FileInfo, err error) error {
		relPath := strings.TrimPrefix(path, src)
		destPath := filepath.Join(dest, relPath)

		if sfi.IsDir() {
			info, err := os.Stat(path)
			if err != nil {
				return err
			}
			if err := os.MkdirAll(destPath, info.Mode()); err != nil {
				return err
			}
			return nil
		}

		if err := copy.Copy(path, destPath); err != nil {
			return err
		}
		if err := os.Chmod(destPath, sfi.Mode()); err != nil {
			return err
		}

		return nil
	})
	return err
}
