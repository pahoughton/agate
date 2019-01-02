/* 2018-12-31 (cc) <paul4hough@gmail.com>
   FIXME what is this for?
*/
package main


type Ticket struct {
	Title		string
	Node		string
	State		string
	Worker		string
	Desc		string
	Comments	[]string
}

type ApiTicket struct {
	ID		string `json:"id,omitempty"`
	Title	string `json:"title,omitempty"`
	Node	string `json:"node,omitempty"`
	State	string `json:"state,omitempty"`
	Worker	string `json:"worker,omitempty"`
	Desc	string `json:"desc,omitempty"`
	Comment	string `json:"comment,omitempty"`
}
