/* 2019-02-17 (cc) <paul4hough@gmail.com>
   implements ticket.AlertTid interface
*/
package tid

type Tid interface {
	String() string
	Sys() uint
	Bytes() []byte
}

type data struct {
	str		string
	sys	uint8
}


func NewBytes(id []byte) *data {
	return &data{ str: string(id[:len(id)-1]), sys: id[len(id)-1], }
}
func NewString(tsys uint8,id string) *data {
	return &data{id,tsys}
}

func (t *data) String() string {
	return t.str
}
func (t *data) Sys() uint {
	return uint(t.sys)
}
func (t *data) Bytes() []byte {
	return append([]byte(t.str),byte(t.sys))
}
