package ops

type Object interface {
	Serialize() []byte
	Format() string
	Data() []byte
}

type Commit struct {
	contents []byte
}

func (commit *Commit) Serialize() []byte {
	return commit.contents
}
func (commit *Commit) Format() string {
	return "commit"
}
func (commit *Commit) Data() []byte {
	return commit.contents
}

type Blob struct {
	contents []byte
}

func (blob *Blob) Serialize() []byte {
	return blob.contents
}
func (blob *Blob) Format() string {
	return "commit"
}
func (blob *Blob) Data() []byte {
	return blob.contents
}
