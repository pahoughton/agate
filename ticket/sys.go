/* 2019-02-16 (cc) <paul4hough@gmail.com>
   FIXME what is this for?
*/
package ticket

func (t *Ticket)ActionCreate(
	tsys uint64, group, title, desc string) (string, error) {

	switch tsys {
	case TSYS_HPSM:
		return t.hpsm.Create(group,title,desc)
	case TSYS_GITLAB:
		return t.gitlab.Create(group,title,desc)
	case TSYS_MOCK:
		return t.mock.Create(group,title,desc)
	default:
		return false
	}
}

func (t *Ticket)ActionUpdate(tsys uint64, tid, desc string) error {

	switch tsys {
	case TSYS_HPSM:
		return t.hpsm.Create(group,title,desc)
	case TSYS_GITLAB:
		return t.gitlab.Create(group,title,desc)
	case TSYS_MOCK:
		return t.mock.Create(group,title,desc)
	default:
		return false
	}
}

func (t *Ticket)ActionUpdate(tsys uint64, tid, desc string) error {

	switch tsys {
	case TSYS_HPSM:
		return t.hpsm.Create(group,title,desc)
	case TSYS_GITLAB:
		return t.gitlab.Create(group,title,desc)
	case TSYS_MOCK:
		return t.mock.Create(group,title,desc)
	default:
		return false
	}
}
