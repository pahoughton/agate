/* 2018-12-25 (cc) <paul4hough@gmail.com>
   Prometheus AlertManager Alerts Body
*/

package amgr

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"net/http"
	"os"
	"path"
	"strings"

	"github.com/pahoughton/agate/model"
	"github.com/pahoughton/agate/config"
	"github.com/pahoughton/agate/ticket"
	"github.com/pahoughton/agate/proc"
	"github.com/pahoughton/agate/db"

	pmod "github.com/prometheus/common/model"
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

func New(c *config.Config, dbg bool) *Handler {

	adb, err := db.Open(path.Join(c.BaseDir, "data"), 0664, c.MaxDays);
	if err != nil {
		fmt.Println("FATAL: open db - ",err.Error())
		os.Exit(1)
	}

	h := &Handler{
		Debug:			dbg,
		Adb:			adb,
		Ticket:			ticket.New(c,dbg),
		Proc:			proc.New(c,dbg),
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


	return h

}

func (h *Handler)ServeHTTP(w http.ResponseWriter,r *http.Request) {
	if err := h.AlertGroup(w,r); err != nil {
		fmt.Println("ERROR: ",err.Error())
		h.Errors.Inc()
    }
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

	var agrp model.AlertGroup
	if err := json.Unmarshal(b, &agrp); err != nil {
		return fmt.Errorf("json.Unmarshal alert: %s\n%v",err.Error(),b)
    }

	h.AlertGroupsRecvd.With(
		promp.Labels{
			"status": agrp.Status,
		}).Inc()

	for _, a := range agrp.Alerts {

		aname := a.Name()
		aKey := a.Key()
		node := "unknown"

		if inst, ok := a.Labels["instance"]; ok {
			node = strings.Split(string(inst),":")[0]
		}

		h.AlertsRecvd.With(
			promp.Labels{
				"name": aname,
				"node": node,
				"status": string(a.Status),
			}).Inc()

		if a.Status == "firing" {

			var (
				tid		string
			)

			// dup prevention
			tid, err = h.Adb.GetTicket(a.StartsAt, aKey)
			if err == nil && len(tid) > 0 {
				h.AlertDups.With(
					promp.Labels{
						"name": aname,
						"node": node,
					}).Inc()
				if h.Debug {
					fmt.Printf("DEBUG: dup alert: %s",a.Labels.String())
				}
				return nil
			}

			pending := ""
			aremed := false
			sremed := false


			ardir := path.Join(h.Proc.PlaybookDir,"roles",aname)
			finfo, err := os.Stat(ardir)
			if err == nil && finfo.IsDir() {
				aremed = true
				pending += "ansible remediation pending\n"
			}

			sfn := path.Join(h.Proc.ScriptsDir,aname)
			finfo, err = os.Stat(sfn)
			if err == nil && (finfo.Mode() & 0111) != 0 {
				sremed = true;
				pending += "script remediation pending\n"
			}

			if aremed == false && sremed == false  {
				pending += "no remediation available\n"
			}

			a.Annotations["pending"] = pmod.LabelValue(pending)
			tid, err = h.Ticket.Create(a)

			if err != nil {
				return fmt.Errorf("ticket.Create: %s",err.Error())
			}

			if err = h.Adb.AddTicket(a.StartsAt,aKey,tid); err != nil {
				return err
			}

			procErr := ""

			if aremed {
				emsg := ""
				out, err := h.Proc.Ansible(node,a.Labels)
				if err != nil {
					emsg = "ERROR: " + err.Error() + "\n"
					procErr += "ansible - " + err.Error() + "\n"
				}
				tcom := "ansible remediation results\n" + emsg + out

				if err = h.Ticket.AddComment(a,tid,tcom); err != nil {
					return fmt.Errorf("ticket add comment: %s\n%s",err,tcom)
				}
			}

			if sremed {
				emsg := ""
				out, err := h.Proc.Script(node,a.Labels)
				if err != nil {
					emsg = "ERROR: " + err.Error() + "\n"
					procErr += "script - " + err.Error() + "\n"
				}
				tcom := "script remediation results\n" + emsg + out

				if err = h.Ticket.AddComment(a,tid,tcom); err != nil {
					return fmt.Errorf("ticket add comment: %s\n%s",err,tcom)
				}
			}
			if len(procErr) > 0 {
				return errors.New(procErr)
			}
		} else if a.Status == "resolved" {

			tid, err := h.Adb.GetTicket(a.StartsAt, aKey)
			if err != nil {
				fmt.Printf("WARN resolved not found: %v",a.Labels)
				return nil
			}

			tcom := fmt.Sprintf("resolved at %v",a.EndsAt)

			if err = h.Ticket.AddComment(a,tid,tcom); err != nil {
				return fmt.Errorf("ticket comment: %s\n%s",err,tcom)
			}

			if h.CloseResolved || a.Annotations["close_resolved"] == "true" {

				if err = h.Ticket.Close(a,tid); err != nil {
					return fmt.Errorf("ticket close: %s",err)
				}
			}
			if err = h.Adb.Delete(a.StartsAt, aKey); err != nil {
				return err
			}
		}
	}
	return nil
}
