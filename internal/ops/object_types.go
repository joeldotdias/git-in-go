package ops

type Object interface {
	Serialize() []byte
	Format() string
	Data() []byte
}

type Blob struct {
	contents []byte
}

func (blob *Blob) Serialize() []byte {
	return blob.contents
}
func (blob *Blob) Format() string {
	return "blob"
}
func (blob *Blob) Data() []byte {
	return blob.contents
}
