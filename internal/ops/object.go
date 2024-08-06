package ops

import (
	"bytes"
	"compress/zlib"
	"crypto/sha1"
	"encoding/hex"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"strings"

	"github.com/joeldotdias/gat/internal/helpers"
)

func (repo *Repository) makeObject(sha string) (Object, error) {
	path := repo.makePath("objects", sha[:2], sha[2:])
	file, err := os.Open(path)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	zr, err := zlib.NewReader(file)
	if err != nil {
		return nil, err
	}
	defer zr.Close()

	data, err := io.ReadAll(zr)
	if err != nil {
		return nil, err
	}

	nullByteIndex := bytes.IndexByte(data, 0)
	if nullByteIndex == -1 {
		return nil, fmt.Errorf("malformed object format")
	}

	header := string(data[:nullByteIndex])
	contents := data[nullByteIndex+1:]

	var objType string
	var size int
	_, err = fmt.Sscanf(header, "%s %d", &objType, &size)
	if err != nil {
		return nil, err
	}

	var obj Object
	switch objType {
	case "commit":
		obj = &Commit{klvm: make(map[string][]string)}
	case "tree":
		obj = &Tree{leaves: []*TreeLeaf{}}
	case "blob":
		obj = &Blob{}
	default:
		return nil, fmt.Errorf("Unknown object type: %s", objType)
	}

	obj.Deserialize(contents)
	return obj, nil
}

func (repo *Repository) writeObject(obj Object, write bool) (string, error) {
	data := obj.Serialize()
	res := []byte(fmt.Sprintf("%s %d\x00", obj.GetType(), len(data)))
	res = append(res, data...)
	hash := sha1.Sum(res)
	sha := hex.EncodeToString(hash[:])

	if write {
		path := repo.makePath("objects", sha[:2], sha[2:])
		dir := filepath.Dir(path)
		if err := os.MkdirAll(dir, 0755); err != nil {
			return "", fmt.Errorf("failed to create directories: %w", err)
		}

		if _, err := os.Stat(path); os.IsNotExist(err) {
			file, err := os.Create(path)
			if err != nil {
				return "", fmt.Errorf("failed to create file: %w", err)
			}
			defer file.Close()

			zw := zlib.NewWriter(file)
			_, err = zw.Write(res)
			if err != nil {
				return "", fmt.Errorf("failed to write compressed data: %w", err)
			}
			err = zw.Close()
			if err != nil {
				return "", fmt.Errorf("failed to close zlib writer: %w", err)
			}
		}
	}

	return sha, nil
}

func (repo *Repository) makeObjectHash(file io.Reader, objFormat string) (string, error) {
	contents, err := io.ReadAll(file)
	if err != nil {
		return "", fmt.Errorf("Couldn't read file: %v", err)
	}

	var obj Object
	switch objFormat {
	case "blob":
		obj = &Blob{data: contents}
	case "tree":
		obj = &Tree{leaves: treeParseEntirety(contents)}
	case "commit":
		commit := &Commit{klvm: make(map[string][]string)}
		commit.Deserialize(contents)
		obj = commit
	default:
		return "", fmt.Errorf("%s is not a valid object type", objFormat)
	}

	return repo.writeObject(obj, false)
}

func (repo *Repository) findObject(name string) (string, error) {
	if len(name) == 40 && helpers.IsHex(name) {
		return name, nil
	}

	path := repo.makePath("refs", name)
	if _, err := os.Stat(path); err != nil {
		content, err := os.ReadFile(path)
		if err != nil {
			return "", fmt.Errorf("failed to read ref file: %w", err)
		}
		return strings.TrimSpace(string(content)), nil
	}

	prefix := name[:2]
	path = repo.makePath("objects", prefix)
	if matches, err := filepath.Glob(filepath.Join(path, name[2:]+"*")); err == nil && len(matches) > 0 {
		if len(matches) > 1 {
			return "", fmt.Errorf("ambiguous object name: %s", name)
		}
		return filepath.Base(prefix + matches[0]), nil
	}

	return "", fmt.Errorf("Didn't find object: %s", name)
}
