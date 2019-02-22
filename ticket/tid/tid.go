/* 2019-02-17 (cc) <paul4hough@gmail.com>
   implements ticket.AlertTid interface
*/
package tid

type Tid struct {
	str		string
	sys	uint8
}


func NewBytes(id []byte) *Tid {
	return &Tid{ str: string(id[:len(id)-1]), sys: id[len(id)-1], }
}
func NewString(tsys uint8,id string) *Tid {
	return &Tid{id,tsys}
}

func (t *Tid) String() string {
	return t.str
}
func (t *Tid) Sys() uint {
	return uint(t.sys)
}
func (t *Tid) Bytes() []byte {
	return append([]byte(t.str),byte(t.sys))
}
