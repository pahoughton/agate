/* 2019-02-17 (cc) <paul4hough@gmail.com>
   implements ticket.AlertTid interface
*/
package tid

type Tid struct {
	str		string
}


func NewBytes(id []byte) *Tid {
	return &Tid{string(id)}
}
func NewString(id string) *Tid {
	return &Tid{id}
}

func (t Tid) String() string {
	return t.str
}
func (t Tid) Bytes() []byte {
	return []byte(t.str)
}
