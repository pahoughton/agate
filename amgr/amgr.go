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

	"github.com/pahoughton/agate/config"
	"github.com/pahoughton/agate/ticket"
	"github.com/pahoughton/agate/remed"
	"github.com/pahoughton/agate/db"

	pmod "github.com/prometheus/common/model"
	proma "github.com/prometheus/client_golang/prometheus/promauto"
	promp "github.com/prometheus/client_golang/prometheus"

)

const (
	ATimeFmt = "2006-01-02T15:04:05.000000000-07:00"
)

type Metrics struct {
    Recvd	*promp.CounterVec
	Alerts	*promp.CounterVec
	Errors	promp.Counter
}


type Amgr struct {
	debug			bool
	db				*db.DB
	ticket			*ticket.Ticket
	remed			*remed.Remed
	qmgr			*Manager
	respq			chan
	metrics			Metrics
}

func New(c *config.Config,dataDir string,dbg bool) *Amgr {

	adb, err := db.Open(dataDir, 0664, c.MaxDays);
	if err != nil {
		panic(err)
	}
	am := &Amgr{
		debug:		dbg,
		db:			adb,
		ticket:		ticket.New(c.Ticket,dbg),
		remed:		remed.New(c.Global,dbg),

		metrics: Metrics{
			Recvd: proma.NewCounter(
				promp.CounterOpts{
					Namespace: "agate",
					Name:      "agroup_received_total",
					Help:      "number of alert groups received",
				}),
			Alerts: proma.NewCounterVec(
				promp.CounterOpts{
					Namespace: "agate",
					Name:      "alerts_received_total",
					Help:      "number of alerts received",
				}, []string{
					"name",
					"node",
				}),
			Errors: proma.NewCounter(
				promp.CounterOpts{
					Namespace: "agate",
					Name:      "amgr_errors_total",
					Help:      "number of amgr errors",
				}),
		},
	}
	am.qmgr = NewManager()
	am.respq = make(chan uint64)
	go am.Manage()
	return am
}
/*
func (h *Handler)AlertQueueManager() {

	for {
		// wait for next alert and double check queue every 10 min
		select {
		case recvd := <- h.alerts:
		case <- time.After(10 * time.Minute):
		}

		for {
			id, err := h.db.AlertNext()
			if err != nil {
				panic(err)
			}
			if id == nil {
				break
			}
			h.proc <- id
		}
	}
}

func (*h Hander)AlertProc(id string) {
	for {
		id := <- h.proc

		var agrp model.AlertGroup
		if err := json.Unmarshal(b, &agrp); err != nil {
			panic(fmt.Sprintf("json.Unmarshal alert: %s\n%v",
				err.Error(),b))
		}

		if h.Debug {
			var dbgbuf bytes.Buffer
			if err := json.Indent(&dbgbuf, b, " ", "  "); err != nil {
				fmt.Printf("DEBUG json.Indent: ",err.Error())
			}
			fmt.Println("DEBUG req body\n",dbgbuf.String())
		}

		h.AlertGroupsRecvd.With(
			promp.Labels{
				"status": agrp.Status,
			}).Inc()

		if agrp.Status == "firing" {

			remed := false

			for _, a := range agrp.Alerts {

				aname := a.Name()
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

				ardir := path.Join(h.Proc.PlaybookDir,"roles",aname)
				finfo, err := os.Stat(ardir)
				if err == nil && finfo.IsDir() {
					remed = true
					break
				}

				sfn := path.Join(h.Proc.ScriptsDir,aname)
				finfo, err = os.Stat(sfn)
				if err == nil && (finfo.Mode() & 0111) != 0 {
					remed = true
					break;
				}
			}
			if remed {
				agrp.ComAnnots['remediation'] = "pending"
			} else {
				agrp.ComAnnots['remediation'] = "none"
			}
			FIX TICKET DATA STRUCT - ticket id by alert group := resolve updates ticket,
			all resolved to close! The rabit hole .... BIG FUN.



			if a.Status == "firing" {

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

			if h.CloseResolved || a.Annotations["close_resolved"] == "true" {

				if err = h.Ticket.Close(a,tid,tcom); err != nil {
					return fmt.Errorf("ticket close: %s",err)
				}
			} else {
				if err = h.Ticket.AddComment(a,tid,tcom); err != nil {
					return fmt.Errorf("ticket comment: %s\n%s",err,tcom)
				}
			}

			if err = h.Adb.Delete(a.StartsAt, aKey); err != nil {
				return err
			}
		}
	}



			id, err = h.db.AlertNext()
		}
		if err != nil {
			bad bad bad
		}
	}
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
	err := h.db.AlertAdd(b)
	if err != nil {
		select {
		case h.alerts <- true:
		case <- time.After(1):
		}
	}
	return err
}
*/
