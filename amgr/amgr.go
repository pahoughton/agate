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
	"strings"
	"time"

	"gitlab.com/pahoughton/agate/config"
	"gitlab.com/pahoughton/agate/ticket"
	"gitlab.com/pahoughton/agate/proc"

	proma "github.com/prometheus/client_golang/prometheus/promauto"
	promp "github.com/prometheus/client_golang/prometheus"
)

const (
	ATimeFmt = "2006-01-02T15:04:05.000000000-07:00"
)

type Handler struct {
	Debug				bool
	Ticket				*ticket.Ticket
	Proc				*proc.Proc
	CloseResolved		bool
	AlertGroupsRecvd	*promp.CounterVec
	AlertsRecvd			*promp.CounterVec
	Errors				promp.Counter
}

func New(c *config.Config) *Handler {
	h := &Handler{
		Debug:	c.Debug,
		Ticket:	ticket.New(c),
		CloseResolved: c.CloseResolved,

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

func (a *Alert)Key() string {
	return a.StartsAt.Format(ATimeFmt) + " " +
		a.Labels["alertname"] + "-" +
		a.Labels["instance"]
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
		node := strings.Split(alert.Labels["instance"],":")[0]

		h.AlertsRecvd.With(
			promp.Labels{
				"name": alert.Labels["alertname"],
				"node": node,
				"status": agrp.Status,
			}).Inc()

		aKey := alert.Key()

		if alert.Status == "firing" {

			node := strings.Split(alert.Labels["instance"],":")[0]
			title := alert.Labels["alertname"] + " on " + node

			desc := "start_time: " + alert.StartsAt.String() + "\n"

			desc += "\nAnnotations:\n"
			for k, v := range alert.Annotations {
				desc += k + ": " + v + "\n"
			}
			desc  += "\nLabels:\n"
			for k, v := range alert.Labels {
				desc += k + ": " + v + "\n"
			}
			desc += "\nfrom: " + alert.GeneratorURL + "\n"

			tsys := alert.Labels["ticket"]
			tsub := alert.Labels[tsys]

			tid, err := h.Ticket.Create(tsys,tsub,aKey,title,desc);

			if err != nil {
				return fmt.Errorf("ticket.Create: %s",err.Error())
			}

			if _, ok := alert.Labels["ansible"]; ok {
				err := h.Proc.Ansible(node,alert.Labels,tsys,tid)
				if err != nil {
					return err
				}
			}
			if _, ok := alert.Labels["script"]; ok {
				err := h.Proc.Script(node,alert.Labels,tsys,tid)
				if err != nil {
					return err
				}
			}
		} else if alert.Status == "resolved" {
			tsys := alert.Labels["ticket"]
			tcom := fmt.Sprintf("resolved at %v",alert.EndsAt)

			if err = h.Ticket.AddKeyComment(tsys,aKey,tcom); err != nil {
				return fmt.Errorf("ticket comment: %s",err)
			}
			// fixme - close resolved label?
			if h.CloseResolved || alert.Labels["close_resolved"] == "true" {
				if err = h.Ticket.Close(tsys,aKey); err != nil {
					return fmt.Errorf("ticket close: %s",err)
				}
			}
			if err = h.Ticket.Delete(aKey); err != nil {
				return err
			}
		}
	}
	return nil
}
