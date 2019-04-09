/* 2019-02-19 (cc) <paul4hough@gmail.com>
*/
package notify

import (
	"fmt"
	promp "github.com/prometheus/client_golang/prometheus"
	"github.com/pahoughton/agate/notify/nid"
)

func (n *Notify) Group(nsys NSys) string {
	if n.System(nsys) != nil {
		return n.System(nsys).Group()
	} else {
		return "invalid"
	}
}

func (n *Notify) Create(
	nsys	NSys,
	grp		string,
	title	string,
	desc	string,
	remed	bool,
	resolve	bool) (nid.Nid, error) {

	if n.System(nsys) != nil {
		var (
			aclose string
			aremed string
		)
		if resolve {
			aclose = "closes on resolve"
		} else {
			aclose = "manual"
		}
		if remed {
			aremed = "true"
		} else {
			aremed = "false"
		}
		ndesc := fmt.Sprintf(
			"\nauto-close: %s  remediation: %s\n%s",
			aclose,
			aremed,
			desc)

		nid, err := n.System(nsys).Create(grp,title,ndesc)
		if err == nil {
			n.metrics.notes.With(promp.Labels{
				"sys": nsys.String(),
				"grp": grp,
			}).Inc()
			return nid, err
		} else {
			n.metrics.errors.Inc()
			return nid, err
		}
	} else {
		panic(fmt.Sprintf("invalid nsys: %d\n",nsys))
		return nil, nil
	}
}

func (n *Notify) Update(nid nid.Nid, msg string) bool {
	if n.System(NSys(nid.Sys())) != nil {
		err := n.System(NSys(nid.Sys())).Update(nid,msg)
		if err == nil {
			return true
		} else {
			n.metrics.errors.Inc()
		}
	}
	return false
}

func (n *Notify) Close(nid nid.Nid, msg string) bool {
	if n.System(NSys(nid.Sys())) != nil {
		err := n.System(NSys(nid.Sys())).Close(nid,msg)
		if err == nil {
			return true
		} else {
			n.metrics.errors.Inc()
		}
	}
	return false
}
