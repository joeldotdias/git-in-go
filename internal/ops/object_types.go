package ops

import (
	"encoding/hex"
	"fmt"
	"sort"
)

type Object interface {
	GetType() string
	Serialize() []byte
	Deserialize(data []byte)
}

type Blob struct {
	data []byte
}

type Tree struct {
	leaves []*TreeLeaf
}

type Commit struct {
	klvm    map[string][]string
	message string
}

func (b *Blob) GetType() string {
	return "blob"
}

func (b *Blob) Serialize() []byte {
	return b.data
}

func (b *Blob) Deserialize(data []byte) {
	b.data = data
}
func (t *Tree) GetType() string {
	return "tree"
}

func (t *Tree) Serialize() []byte {
	sort.Slice(t.leaves, func(i, j int) bool {
		return sortLeafByKey(t.leaves[i]) < sortLeafByKey(t.leaves[j])
	})

	var serialized []byte
	for _, i := range t.leaves {
		serialized = append(serialized, i.mode...)
		serialized = append(serialized, ' ')
		serialized = append(serialized, []byte(i.path)...)
		serialized = append(serialized, 0)
		sha, _ := hex.DecodeString(i.sha)
		serialized = append(serialized, sha...)
	}
	return serialized
}

func (t *Tree) Deserialize(data []byte) {
	t.leaves = treeParseEntirety(data)
}

func (c *Commit) GetType() string {
	return "commit"
}

func (c *Commit) Serialize() []byte {
	return serializeKlvm(c.klvm, c.message)
}

func (c *Commit) Deserialize(data []byte) {
	c.klvm, c.message = parseKlvm(data)
}

// helper commit function
func (c *Commit) getField(key string) (string, error) {
	values, ok := c.klvm[key]
	if !ok || len(values) == 0 {
		return "", fmt.Errorf("Field %s does not exist", key)
	}
	return values[0], nil
}
