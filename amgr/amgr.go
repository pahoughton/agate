/* 2018-12-25 (cc) <paul4hough@gmail.com>
   Prometheus AlertManager Alerts Body
*/

package amgr

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"os"
	"path"
	"sort"
	"strings"
	"time"

	"github.com/pahoughton/agate/config"
	"github.com/pahoughton/agate/ticket"
	"github.com/pahoughton/agate/proc"
	"github.com/pahoughton/agate/db"

	"github.com/prometheus/common/model"
	proma "github.com/prometheus/client_golang/prometheus/promauto"
	promp "github.com/prometheus/client_golang/prometheus"

)

const (
	ATimeFmt = "2006-01-02T15:04:05.000000000-07:00"
)

type Handler struct {
	Debug				bool
	Adb					*db.AlertDB
	Ticket				*ticket.Ticket
	Proc				*proc.Proc
	CloseResolved		bool
	AlertGroupsRecvd	*promp.CounterVec
	AlertsRecvd			*promp.CounterVec
	AlertDups			*promp.CounterVec
	Errors				promp.Counter
}

func New(c *config.Config) *Handler {

	adb, err := db.Open(path.Join(c.BaseDir, "data"), 0664, c.MaxDays);
	if err != nil {
		fmt.Println("FATAL: open db - ",err.Error())
		os.Exit(1)
	}

	h := &Handler{
		Debug:			c.Debug,
		Adb:			adb,
		Ticket:			ticket.New(c),
		CloseResolved:	c.CloseResolved,

		AlertGroupsRecvd: proma.NewCounterVec(
			promp.CounterOpts{
				Namespace: "agate",
				Name:      "alert_group_received_total",
				Help:      "number of alert groups received",
			}, []string{
				"status",
			}),
		AlertsRecvd: proma.NewCounterVec(
			promp.CounterOpts{
				Namespace: "agate",
				Name:      "alerts_received_total",
				Help:      "number of alerts received",
			}, []string{
				"name",
				"node",
				"status",
			}),
		AlertDups: proma.NewCounterVec(
			promp.CounterOpts{
				Namespace: "agate",
				Name:      "alert_dups_total",
				Help:      "number of duplicate alerts received",
			}, []string{
				"name",
				"node",
			}),
		Errors: proma.NewCounter(
			promp.CounterOpts{
				Namespace: "agate",
				Name:      "errors_total",
				Help:      "number of errors",
			}),
	}

	h.Proc = proc.New(c,h.Ticket)

	return h

}

func (h *Handler)ServeHTTP(w http.ResponseWriter,r *http.Request) {
	if err := h.AlertGroup(w,r); err != nil {
		fmt.Println("ERROR: ",err.Error())
		h.Errors.Inc()
    }
}



type Alert struct {

	Annotations map[string]string `json:"annotations,omitempty"`

	StartsAt time.Time `json:"startsAt"`

	EndsAt time.Time `json:"endsAt,omitempty"`

	GeneratorURL string `json:"generatorURL"`

	Labels map[string]string `json:"labels"`

	Status string `json:"status"`
}

type AlertGroup struct {

	Alerts []Alert `json:"alerts"`

	CommonAnnotations map[string]string `json:"commonAnnotations,omitempty"`

	CommonLabels map[string]string `json:"commonLabels,omitempty"`

	ExternalURL string `json:"externalURL"`

	GroupKey string `json:"groupKey"`

	GroupLabels map[string]string `json:"groupLabels,omitempty"`

	Receiver string `json:"receiver"`

	Status string `json:"status"`

	Version string `json:"version"`
}

func (a *Alert)Key() uint64 {

	labs := model.LabelSet{}

	for k, v := range a.Labels {
		labs[model.LabelName(k)] = model.LabelValue(v)
	}
	return uint64(labs.Fingerprint())
}

func (h *Handler)AlertGroup(w http.ResponseWriter,r *http.Request ) error {

	b, err := ioutil.ReadAll(r.Body)
	if err != nil {
		return fmt.Errorf("ioutil.ReadAll - %s",err.Error())
	}
	defer r.Body.Close()

	if h.Debug {
		var dbgbuf bytes.Buffer
		if err := json.Indent(&dbgbuf, b, " ", "  "); err != nil {
			return fmt.Errorf("json.Indent: ",err.Error())
		}
		fmt.Println("DEBUG req body\n",dbgbuf.String())
	}

	var agrp AlertGroup
	if err := json.Unmarshal(b, &agrp); err != nil {
		return fmt.Errorf("json.Unmarshal alert: %s\n%v",err.Error(),b)
    }

	h.AlertGroupsRecvd.With(
		promp.Labels{
			"status": agrp.Status,
		}).Inc()

	for _, alert := range agrp.Alerts {

		aname := alert.Labels["alertname"]
		node := strings.Split(alert.Labels["instance"],":")[0]
		aKey := alert.Key()
		tsys := alert.Annotations["ticket"]

		if len(tsys) < 1 {
			tsys = h.Ticket.DefaultSys
		}

		h.AlertsRecvd.With(
			promp.Labels{
				"name": aname,
				"node": node,
				"status": agrp.Status,
			}).Inc()

		if alert.Status == "firing" {

			var (
				ok		bool
				title	string
				desc	string
				tid		string
			)

			// dup prevention
			tid, err = h.Adb.GetTicket(alert.StartsAt, aKey)
			if err == nil && len(tid) > 0 {
				h.AlertDups.With(
					promp.Labels{
						"name": aname,
						"node": node,
					}).Inc()
				if h.Debug {
					fmt.Printf("dup alert: %v",alert.Labels)
				}
				return nil
			}

			if _, ok = alert.Annotations["title"]; ok {
				title = alert.Annotations["title"]
			} else if  _, ok = alert.Annotations["subject"]; ok {
				title = alert.Annotations["subject"]
			} else {
				title = aname + " on " + node
			}

			desc = "from: " + alert.GeneratorURL + "\n"
			desc += "when: " + alert.StartsAt.String() + "\n"

			desc += "\nAnnotations:\n"
			ankeys := make([]string, 0, len(alert.Annotations))
			for ak, _ := range alert.Annotations {
				ankeys = append(ankeys, ak)
			}
			sort.Strings(ankeys)
			for _, ak := range ankeys {
				desc += ak + ": " +  alert.Annotations[ak]  + "\n"
			}

			desc  += "\nLabels:\n"
			lbkeys := make([]string, 0, len(alert.Labels))
			for lk, _ := range alert.Labels {
				lbkeys = append(lbkeys, lk)
			}
			sort.Strings(lbkeys)
			for _, lk := range lbkeys {
				desc += lk + ": " + alert.Labels[lk] + "\n"
			}

			aremed := false
			sremed := false

			ardir := path.Join(h.Proc.PlaybookDir,"roles",aname)
			finfo, err := os.Stat(ardir)
			if err == nil && finfo.IsDir() {
				aremed = true
				desc += "\nansible remediation pending\n"
			}

			sfn := path.Join(h.Proc.ScriptsDir,aname)
			finfo, err = os.Stat(sfn)
			if err == nil && (finfo.Mode() & 0111) != 0 {
				sremed = true;
				desc += "\nscript remediation pending\n"
			}

			if aremed == false && sremed == false  {
				desc += "\nno remediation available\n"
			}

			tsub := alert.Annotations["ticket_group"]
			if len(tsub) < 1 {
				tsub = h.Ticket.DefaultGrp
			}

			tid, err = h.Ticket.Create(tsys,tsub,title,desc);

			if err != nil {
				return fmt.Errorf("ticket.Create: %s",err.Error())
			}

			if err = h.Adb.AddTicket(alert.StartsAt,aKey,tid); err != nil {
				return err
			}
			if aremed {
				err := h.Proc.Ansible(node,alert.Labels,tsys,tid)
				if err != nil {
					return err
				}
			}
			if sremed {
				err := h.Proc.Script(node,alert.Labels,tsys,tid)
				if err != nil {
					return err
				}
			}

		} else if alert.Status == "resolved" {

			tid, err := h.Adb.GetTicket(alert.StartsAt, aKey)
			if err != nil {
				fmt.Printf("WARN resolved not found: %v",alert.Labels)
				return nil
			}

			tcom := fmt.Sprintf("resolved at %v",alert.EndsAt)

			if err = h.Ticket.AddComment(tsys,tid,tcom); err != nil {
				return fmt.Errorf("ticket comment: %s",err)
			}

			if h.CloseResolved ||
				alert.Annotations["close_resolved"] == "true" {

				if err = h.Ticket.Close(tsys,tid); err != nil {
					return fmt.Errorf("ticket close: %s",err)
				}
			}
			if err = h.Adb.Delete(alert.StartsAt, aKey); err != nil {
				return err
			}
		}
	}
	return nil
}
