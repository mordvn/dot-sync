package sync

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/mordvn/dotfiles-sync/internal/config"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestNewCopier(t *testing.T) {
	c := NewCopier("/tmp/dotfiles")
	require.NotNil(t, c)
	assert.Equal(t, "/tmp/dotfiles", c.basePath)
}

func TestCopier_Copy_File(t *testing.T) {
	tmpDir := t.TempDir()
	srcFile := filepath.Join(tmpDir, "src.txt")
	require.NoError(t, os.WriteFile(srcFile, []byte("hello"), 0o644))

	destBase := filepath.Join(tmpDir, "dest")
	c := NewCopier(destBase)
	err := c.Copy(config.PathConfig{
		Source: srcFile,
		Type:   "file",
		Name:   "src",
	})
	require.NoError(t, err)

	destPath := filepath.Join(destBase, "src")
	data, err := os.ReadFile(destPath)
	require.NoError(t, err)
	assert.Equal(t, "hello", string(data))
}

func TestCopier_Copy_File_SourceNotFound(t *testing.T) {
	tmpDir := t.TempDir()
	c := NewCopier(tmpDir)
	err := c.Copy(config.PathConfig{
		Source: filepath.Join(tmpDir, "nonexistent"),
		Type:   "file",
		Name:   "x",
	})
	require.Error(t, err)
	assert.Contains(t, err.Error(), "source not found")
}

func TestCopier_Copy_Directory(t *testing.T) {
	tmpDir := t.TempDir()
	srcDir := filepath.Join(tmpDir, "srcdir")
	require.NoError(t, os.MkdirAll(srcDir, 0o755))
	require.NoError(t, os.WriteFile(filepath.Join(srcDir, "a.txt"), []byte("a"), 0o644))
	require.NoError(t, os.MkdirAll(filepath.Join(srcDir, "sub"), 0o755))
	require.NoError(t, os.WriteFile(filepath.Join(srcDir, "sub", "b.txt"), []byte("b"), 0o644))

	destBase := filepath.Join(tmpDir, "dest")
	c := NewCopier(destBase)
	err := c.Copy(config.PathConfig{
		Source: srcDir,
		Type:   "directory",
		Name:   "myname",
	})
	require.NoError(t, err)

	destPath := filepath.Join(destBase, "myname", "config")
	aPath := filepath.Join(destPath, "a.txt")
	bPath := filepath.Join(destPath, "sub", "b.txt")
	aData, err := os.ReadFile(aPath)
	require.NoError(t, err)
	bData, err := os.ReadFile(bPath)
	require.NoError(t, err)
	assert.Equal(t, "a", string(aData))
	assert.Equal(t, "b", string(bData))
}

func TestExpandPath_EnvVar(t *testing.T) {
	const testVal = "/some/test/path"
	t.Setenv("TEST_EXPAND_PATH", testVal)
	got := expandPath("$TEST_EXPAND_PATH")
	assert.Equal(t, testVal, got)
}

func TestExists(t *testing.T) {
	tmpDir := t.TempDir()
	assert.True(t, exists(tmpDir))
	assert.False(t, exists(filepath.Join(tmpDir, "nonexistent")))
}
