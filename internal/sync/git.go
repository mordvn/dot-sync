package sync

import (
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"time"
)

type GitManager struct {
	repoPath string
	repoURL  string
}

func NewGitManager(repoPath, repoURL string) *GitManager {
	return &GitManager{
		repoPath: repoPath,
		repoURL:  repoURL,
	}
}

func (g *GitManager) CommitAndPush(message string) error {
	gitDir := filepath.Join(g.repoPath, ".git")
	if _, err := os.Stat(gitDir); os.IsNotExist(err) {
		if err := g.initRepo(); err != nil {
			return fmt.Errorf("failed to init repo: %v", err)
		}
	}

	if err := g.runCmd("git", "add", "."); err != nil {
		return fmt.Errorf("failed to git add: %v", err)
	}

	statusOutput, err := g.cmdOutput("git", "status", "--porcelain")
	if err != nil || statusOutput == "" {
		return nil
	}

	if message == "" {
		message = fmt.Sprintf("backup from %s", time.Now().Format("2006-01-02 15:04:05"))
	}

	if err := g.runCmd("git", "commit", "-m", message); err != nil {
		return fmt.Errorf("failed to git commit: %v", err)
	}

	if err := g.runCmd("git", "push", "origin", "main"); err != nil {
		if err := g.runCmd("git", "push", "origin", "master"); err != nil {
			return fmt.Errorf("failed to git push: %w", err)
		}
	}

	return nil
}

func (g *GitManager) GetStatus() (string, error) {
	output, err := g.cmdOutput("git", "status", "--short")
	if err != nil {
		return "", err
	}
	if output == "" {
		return "✓ clean", nil
	}
	return fmt.Sprintf("⚠ %d changes", len(output)/3), nil
}

func (g *GitManager) initRepo() error {
	if err := g.runCmd("git", "init"); err != nil {
		return err
	}

	if g.repoURL != "" {
		if err := g.runCmd("git", "remote", "add", "origin", g.repoURL); err != nil {
			return err
		}
	}

	if err := g.runCmd("git", "config", "user.email", "dotfiles@local"); err != nil {
		return err
	}

	if err := g.runCmd("git", "config", "user.name", "dotfiles"); err != nil {
		return err
	}

	if err := g.runCmd("git", "add", "."); err != nil {
		return err
	}

	return g.runCmd("git", "commit", "-m", "init commit")
}

func (g *GitManager) runCmd(name string, args ...string) error {
	cmd := exec.Command(name, args...)
	cmd.Dir = g.repoPath
	cmd.Stdout = nil
	cmd.Stderr = nil
	return cmd.Run()
}

func (g *GitManager) cmdOutput(name string, args ...string) (string, error) {
	cmd := exec.Command(name, args...)
	cmd.Dir = g.repoPath
	output, err := cmd.Output()
	return string(output), err
}
