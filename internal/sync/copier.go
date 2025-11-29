package sync

import (
	"fmt"
	"os"
	"os/user"
	"path/filepath"
	"strings"

	"github.com/mordvn/dotfiles-sync/internal/config"
)

type Copier struct {
	basePath string
}

func NewCopier(basePath string) *Copier {
	return &Copier{
		basePath: basePath,
	}
}

func (c *Copier) Copy(pathCfg config.PathConfig) error {
	expandedSrc := expandPath(pathCfg.Source)

	srcInfo, err := os.Stat(expandedSrc)
	if err != nil {
		return fmt.Errorf("source not found: %s (%v)", pathCfg.Source, err)
	}

	isDir := pathCfg.Type == "directory" || srcInfo.IsDir()

	dest := filepath.Join(c.basePath, strings.TrimSuffix(pathCfg.Name, " Config"))
	if isDir {
		dest = filepath.Join(dest, "config")
	}

	destParent := filepath.Dir(dest)
	if err := os.MkdirAll(destParent, 0o755); err != nil {
		return fmt.Errorf("failed to create dest dir: %v", err)
	}

	if isDir && exists(dest) {
		if err := os.RemoveAll(dest); err != nil {
			return fmt.Errorf("failed to remove existing dir: %v", err)
		}
	}

	if isDir {
		return copyDir(expandedSrc, dest)
	}
	return copyFile(expandedSrc, dest)
}

func copyFile(src, dest string) error {
	srcData, err := os.ReadFile(src)
	if err != nil {
		return fmt.Errorf("failed to read source file: %v", err)
	}

	if err := os.MkdirAll(filepath.Dir(dest), 0o755); err != nil {
		return fmt.Errorf("failed to create dest dir: %v", err)
	}

	if err := os.WriteFile(dest, srcData, 0o644); err != nil {
		return fmt.Errorf("failed to write dest file: %v", err)
	}

	return nil
}

func copyDir(src, dest string) error {
	return filepath.Walk(src, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}

		relPath, err := filepath.Rel(src, path)
		if err != nil {
			return err
		}

		destPath := filepath.Join(dest, relPath)

		if info.IsDir() {
			return os.MkdirAll(destPath, info.Mode())
		}

		return copyFile(path, destPath)
	})
}

func expandPath(path string) string {
	if strings.HasPrefix(path, "~") {
		usr, err := user.Current()
		if err != nil {
			return path
		}
		return filepath.Join(usr.HomeDir, path[1:])
	}
	return os.ExpandEnv(path)
}

func exists(path string) bool {
	_, err := os.Stat(path)
	return err == nil
}
