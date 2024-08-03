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
	"strconv"
)

func (repo *Repository) DecodeObject(sha string) (Object, error) {
	objPath := repo.repoPath("objects", sha[:2], sha[2:])
	file, err := os.Open(objPath)
	if err != nil {
		return nil, err
	}

	z, err := zlib.NewReader(file)
	if err != nil {
		return nil, err
	}
	defer z.Close()

	dcomp, err := io.ReadAll(z)
	if err != nil {
		return nil, fmt.Errorf("Couldn't decompress file: %v", err)
	}

	x := bytes.IndexByte(dcomp, ' ')
	if x == -1 {
		return nil, fmt.Errorf("Malformed object '%s': didn't find space", sha)
	}
	objType := string(dcomp[:x])

	y := bytes.Index(dcomp[x:], []byte("\x00"))
	if y == -1 {
		return nil, fmt.Errorf("Malformed object '%s': didn't find null byte", sha)
	}
	y = x + y

	size, err := strconv.Atoi(string(dcomp[x+1 : y]))
	if size != len(dcomp)-y-1 {
		return nil, fmt.Errorf("Malformed object '%s': invalid size", sha)
	}

	contents := dcomp[y+1:]

	switch objType {
	case "blob":
		return &Blob{contents}, nil
	case "commit":
		fmt.Println("commit")
	case "tree":
		fmt.Println("tree")
	case "tag":
		fmt.Println("tag")
	}

	return nil, nil
}

func (repo *Repository) WriteObject(obj Object, write bool) (string, error) {
	data := obj.Serialize()
	res := []byte(fmt.Sprintf("%s %d\x00", obj.Format(), len(data)))
	res = append(res, data...)

	hash := sha1.Sum(res)
	sha := hex.EncodeToString(hash[:])
	if write {
		path := repo.repoPath("objects", sha[:2], sha[2:])
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
		obj = &Blob{contents}
	default:
		panic(objFormat + " is not yet implemented")
	}

	return repo.WriteObject(obj, false)
}

func (repo *Repository) objectFind(name string, _ []byte) string {
	return name
}
