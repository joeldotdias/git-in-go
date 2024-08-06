package ops

import (
	"bytes"
	"encoding/hex"
)

type TreeLeaf struct {
	mode []byte
	path string
	sha  string
}

func treeParseSingle(raw []byte, start int) (int, *TreeLeaf) {
	x := bytes.IndexByte(raw[start:], ' ')
	if x == -1 {
		panic("Malormed tree entry: no space found")
	}

	mode := raw[start : start+x]
	y := bytes.IndexByte(raw[start+x+1:], 0)
	if y == -1 {
		panic("Invalid tree entry: no null byte found")
	}
	path := raw[start+x+1 : start+x+1+y]

	// Extract SHA
	// should be 20 bytes after the null byte
	shaStart := start + x + 1 + y + 1
	if len(raw) < shaStart+20 {
		panic("Invalid tree entry: not enough bytes for SHA")
	}
	sha := hex.EncodeToString(raw[shaStart : shaStart+20])
	return shaStart + 20, &TreeLeaf{
		mode: mode,
		path: string(path),
		sha:  sha,
	}
}

func treeParseEntirety(raw []byte) []*TreeLeaf {
	pos := 0
	max := len(raw)
	parsed := []*TreeLeaf{}
	for pos < max {
		var data *TreeLeaf
		pos, data = treeParseSingle(raw, pos)
		parsed = append(parsed, data)
	}

	return parsed
}

func sortLeafByKey(leaf *TreeLeaf) string {
	if bytes.HasPrefix(leaf.mode, []byte("40")) {
		return leaf.path
	}

	// if not a dir, append separator
	return leaf.path + "/"
}
