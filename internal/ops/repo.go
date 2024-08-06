package ops

import (
	"errors"
	"fmt"
	"os"
	"path/filepath"

	"github.com/joeldotdias/gat/internal/helpers"
	"gopkg.in/ini.v1"
)

// more fields will be added when needed
type Config struct {
	defaultBranch string
}

func makeCfg() Config {
	homedir, _ := os.UserHomeDir()
	cfg, err := ini.Load(filepath.Join(homedir, ".gitconfig"))
	if err != nil {
		fmt.Fprintf(os.Stderr, "Couldn't read global config file: %v", err)
	}

	return Config{
		defaultBranch: cfg.Section("init").Key("defaultBranch").String(),
	}
}

type Repository struct {
	worktree string
	gitdir   string
	conf     Config
	rootDir  string
}

func Repo(path string) *Repository {
	return &Repository{
		worktree: path,
		gitdir:   filepath.Join(path, ".git"),
		conf:     makeCfg(),
	}
}

func (repo *Repository) makePath(paths ...string) string {
	var rootDir string
	if repo.rootDir != "" {
		rootDir = filepath.Join(repo.rootDir, ".git")
	} else {
		rootDir = repo.gitdir
	}

	parts := append([]string{rootDir}, paths...)
	return filepath.Join(parts...)
}

func findRepoRoot(path string) (string, error) {
	if path == "" {
		path = "."
	}

	absPath, err := filepath.Abs(path)
	if err != nil {
		return "", err
	}

	if helpers.IsDir(filepath.Join(absPath, ".git")) {
		return absPath, nil
	}

	parent := filepath.Dir(absPath)
	if parent == absPath {
		return "", errors.New("Not in a repository. Run gat init to make one.")
	}

	return findRepoRoot(parent)
}
