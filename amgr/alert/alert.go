/* 2019-01-19 (cc) <paul.houghton.ywi9@statefarm.com>
   agate models
*/
package alert

import (
	"encoding/binary"
	"fmt"
	"sort"
	"strings"
	// "time"
	pmod "github.com/prometheus/common/model"
	amgrtmpl "github.com/prometheus/alertmanager/template"
)

const (
	TIMEFMT = "2006-01-02 15:04:05.9999 -0700"
	GRP_TITLE_LABEL = "group_title"
)
var (
	NODE_LABELS = []string{"agate_node", "hostname", "node", "instance"}
	TITLE_LABELS = []string{"agate_title", "title", "subject"}
	NOTIFY_LABELS = []string{"notify_sys","notify_grp"}
)

type LabelSet	amgrtmpl.KV
type Alert		amgrtmpl.Alert


func (lm LabelSet)SortedKeys() []string {

	lkeys := make([]string, 0, len(lm))
	for k, _ := range lm {
		skip :=  k == GRP_TITLE_LABEL
		if skip {
			continue
		}
		for _, t := range TITLE_LABELS {
			if k == t {
				skip = true
				break;
			}
		}
		if skip {
			continue
		}
		for _, t := range NOTIFY_LABELS {
			if k == t {
				skip = true
				break;
			}
		}
		if skip {
			continue
		}
		lkeys = append(lkeys, k)
	}

	sort.Strings(lkeys)
	return lkeys
}
func (a Alert) LabelSet() LabelSet {
	return LabelSet(a.Labels)
}

func (a Alert) Key() []byte {
	if b, err := a.StartsAt.MarshalBinary(); err == nil {
		k := make([]byte,len(b),len(b)+binary.MaxVarintLen64)
		copy(k,b)
		pls := make(pmod.LabelSet,len(a.Labels))
		for k,v := range a.Labels {
			pls[pmod.LabelName(k)] = pmod.LabelValue(v)
		}

		fp := pls.Fingerprint()
		bfp := make([]byte,binary.MaxVarintLen64)

		fn := binary.PutUvarint(bfp,uint64(fp))
		return append(k,bfp[:fn]...)
	} else {
		panic(err)
	}
}

func (a Alert) Name() string {
	if v, ok := a.Labels["alertname"]; ok {
		return v
	} else {
		return ""
	}
}

func (a Alert) Node() string {
	keys := NODE_LABELS
	for _, k := range keys {
		if v, ok := a.Labels[k]; ok {
			node := string(v)
			if i := strings.IndexRune(node,':'); i > 0 {
				return node[:i]
			} else {
				return node
			}
		}
	}
	return ""
}


func (a Alert) Title() string {

	for _, k := range TITLE_LABELS {
		if _, ok := a.Annotations[k]; ok {
			return string(a.Annotations[k])
		}
		if _, ok := a.Labels[k]; ok {
			return string(a.Labels[k])
		}
	}
	title := a.Name()
	if len(title) < 1 {
		title = "alert"
	}
	return title + " on " + a.Node()
}

func (a Alert) Desc() string {

	desc := "\nFrom: " + a.GeneratorURL + "\n\n"
	desc += "When: " + a.StartsAt.Format(TIMEFMT) + "\n"

	if len(a.Labels) > 0 {
		desc  += "\nLabels:\n"
		for _, k := range LabelSet(a.Labels).SortedKeys() {
			desc += fmt.Sprintf("%16s: %s\n",k,a.Labels[k])
		}
	}
	if len(a.Annotations) > 0 {
		desc += "\nAnnotations:\n"
		for _, k := range LabelSet(a.Annotations).SortedKeys() {
			desc += fmt.Sprintf("%16s: %s\n",k,a.Annotations[k])
		}
	}
	return desc
}
