package ops

import (
	"fmt"
	"os"
	"path/filepath"

	"gopkg.in/ini.v1"
)

// Initializes a repository in the current workind dir
// TODO: add path param
func (repo *Repository) Init() {
	for _, dir := range []string{".git", ".git/objects", ".git/refs", ".git/branches"} {
		if err := os.MkdirAll(dir, 0755); err != nil {
			fmt.Fprintf(os.Stderr, "Couldn't create directories: %s\n", err)
		}
	}

	toWrite := map[string]string{
		"HEAD":        "ref: refs/heads/" + repo.conf.defaultBranch + "\n",
		"description": "Unnamed repository; edit this file 'description' to name the repository.\n",
	}

	for fname, contents := range toWrite {
		if err := os.WriteFile(filepath.Join(".git", fname), []byte(contents), 0644); err != nil {
			fmt.Fprintf(os.Stderr, "Error writing file: %v\n", err)
		}
	}

	ini.PrettyFormat = false
	ini.PrettyEqual = true
	configContents := ini.Empty()
	defaultConfig := map[string]string{"repositoryformatversion": "0", "filemode": "true", "bare": "false"}
	coreSec, _ := configContents.NewSection("core")
	for k, v := range defaultConfig {
		_, err := coreSec.NewKey(k, v)
		if err != nil {
			fmt.Fprintf(os.Stderr, "Error writing config key '%s': %v\n", k, err)
		}
	}

	err := configContents.SaveToIndent(".git/config", "\t")
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error writing config file: %v\n", err)
	}
}
