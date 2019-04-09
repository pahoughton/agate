/* 2019-02-17 (cc) <paul4hough@gmail.com>
   implements ticket.AlertTid interface
*/
package nid

type Nid interface {
	Id() string
	Sys() uint
	Bytes() []byte
}

type data struct {
	id		string
	sys		uint8
}

func NewBytes(id []byte) *data {
	return &data{ id: string(id[:len(id)-1]), sys: id[len(id)-1], }
}
func NewString(nsys uint8,id string) *data {
	return &data{id,nsys}
}

func (t *data) Id() string {
	return t.id
}
func (t *data) Sys() uint {
	return uint(t.sys)
}
func (t *data) Bytes() []byte {
	return append([]byte(t.id),byte(t.sys))
}
