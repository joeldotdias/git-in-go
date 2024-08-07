package ops

import (
	"fmt"
	"os"
	"path/filepath"

	"gopkg.in/ini.v1"
)

// Initializes a repository in the current working dir
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

func (repo *Repository) CatFile(objFormat string, object string) {
	var err error
	repo.rootDir, err = findRepoRoot("")
	if err != nil {
		fmt.Fprintf(os.Stderr, err.Error())
	}

	sha, err := repo.findObject(object)
	if err != nil {
		fmt.Println(err.Error())
	}

	obj, err := repo.makeObject(sha)
	if err != nil {
		fmt.Println(err.Error())
	}

	switch objFormat {
	case "blob":
		if blob, ok := obj.(*Blob); ok {
			fmt.Print(string(blob.data))
		} else {
			fmt.Println("Not a blob")
		}
	case "commit", "tree":
		fmt.Print(string(obj.Serialize()))
	}

}

func (repo *Repository) HashObject(write bool, objFormat string, path string) {
	var err error
	if write {
		repo.rootDir, err = findRepoRoot("")
		if err != nil {
			fmt.Fprintf(os.Stderr, err.Error())
		}
	}

	file, err := os.Open(path)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Couldn't open file: %v\n", err)
	}
	defer file.Close()

	sha, err := repo.makeObjectHash(file, objFormat)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Couldn't hash object: %v\n", err)
	}

	fmt.Println(sha)
}

func (repo *Repository) TopLsTree(tree string, recursive bool) error {
	var err error
	repo.rootDir, err = findRepoRoot("")
	if err != nil {
		fmt.Fprintf(os.Stderr, err.Error())
	}

	return repo.lsTree(tree, recursive, "")
}

func (repo *Repository) lsTree(ref string, recursive bool, prefix string) error {
	sha, _ := repo.findObject(ref)
	obj, err := repo.makeObject(sha)
	if err != nil {
		return err
	}

	var tree *Tree
	switch o := obj.(type) {
	case *Commit:
		treeSha, err := o.getField("tree")
		if err != nil {
			return err
		}
		treeObj, err := repo.makeObject(treeSha)
		if err != nil {
			return err
		}
		var ok bool
		tree, ok = treeObj.(*Tree)
		if !ok {
			return fmt.Errorf("Couldn't get tree object from this commit")
		}
	case *Tree:
		tree = o
	default:
		return fmt.Errorf("Object %s is neither a tree nor a commit. It's type is %s", sha, obj.GetType())
	}

	for _, leaf := range tree.leaves {
		var typeStr string
		modeStr := string(leaf.mode)

		switch modeStr {
		case "40000":
			typeStr = "tree"
		case "100644", "100664", "100755":
			typeStr = "blob"
		case "120000":
			typeStr = "blob" // but it's a symlink
		case "160000":
			typeStr = "commit"
		default:
			return fmt.Errorf("unknown mode %s", modeStr)
		}

		if !recursive || (recursive && typeStr == "blob") {
			fmt.Printf("%06s %s %s\t%s\n",
				modeStr,
				typeStr,
				leaf.sha,
				filepath.Join(prefix, leaf.path))
		}

		if recursive && typeStr == "tree" {
			err := repo.lsTree(leaf.sha, recursive, filepath.Join(prefix, leaf.path))
			if err != nil {
				return err
			}
		}
	}

	return nil
}
