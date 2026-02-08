package sync

import (
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestNewGitManager(t *testing.T) {
	g := NewGitManager("/tmp/repo", "https://github.com/example/dotfiles.git")
	require.NotNil(t, g)
	assert.Equal(t, "/tmp/repo", g.repoPath)
	assert.Equal(t, "https://github.com/example/dotfiles.git", g.repoURL)
}

func TestGitManager_GetStatus_NoRepo(t *testing.T) {
	tmpDir := t.TempDir()
	g := NewGitManager(tmpDir, "")
	_, err := g.GetStatus()
	assert.Error(t, err)
}

func TestGitManager_InitRepo_GetStatus(t *testing.T) {
	tmpDir := t.TempDir()
	g := NewGitManager(tmpDir, "")

	err := g.initRepo()
	require.NoError(t, err)

	status, err := g.GetStatus()
	require.NoError(t, err)
	assert.Equal(t, "✓ clean", status)
}

func TestGitManager_GetStatus_WithChanges(t *testing.T) {
	tmpDir := t.TempDir()
	g := NewGitManager(tmpDir, "")
	require.NoError(t, g.initRepo())

	f := filepath.Join(tmpDir, "newfile.txt")
	require.NoError(t, os.WriteFile(f, []byte("new"), 0o644))

	status, err := g.GetStatus()
	require.NoError(t, err)
	assert.True(t, strings.HasPrefix(status, "⚠"), "status should indicate changes: %s", status)
	assert.Contains(t, status, "changes")
}
