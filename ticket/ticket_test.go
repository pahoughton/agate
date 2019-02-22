/* 2019-02-17 (cc) <paul4hough@gmail.com>
   FIXME what is this for?
*/
package ticket

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/pahoughton/agate/config"
)

func TestNew(t *testing.T) {
	cfg := config.New()
	got := New(cfg.Ticket,false)
	assert.NotNil(t,got)
	got.Close()
}

func TestNewPanic(t *testing.T) {
	assert.Panics(t, func() {
		cfg := config.New()
		cfg.Ticket.Default = "george"
		New(cfg.Ticket,false)
	}, "New did not panic")
}

/* FIXME
func TestAlertTSys(t *testing.T) {
	obj := New(config.New().Ticket,false)

	for k, exp := range tsysmap {

		a := alert.Alert{pmod.Alert{
			Labels: pmod.LabelSet{ "ticket_sys": pmod.LabelValue(k) }},
			"firing",
		}
		got := obj.AlertTSys(a)
		assert.Equal(t,got,exp)
	}

	a := alert.Alert{}
	exp := obj.Default
	got := obj.AlertTSys(a)
	assert.Equal(t,got,exp)

	a.Labels = pmod.LabelSet{ "ticket_sys": "invalid" }
	got = obj.AlertTSys(a)
	assert.Equal(t,got,exp)
	promp.Unregister(obj.MetrTicketsGend)
	promp.Unregister(obj.MetrErrors)
}

func TestAGroupTSys(t  *testing.T) {
	obj := New(config.New().Ticket,false)

	for k, exp := range tsysmap {

		ag := alert.AlertGroup{}
		ag.ComLabels = alert.LabelMap{ "ticket_sys":  pmod.LabelValue(k) }

		got := obj.AgroupTSys(ag)
		assert.Equal(t,got,exp)

		a1 := alert.Alert{pmod.Alert{
			Labels: pmod.LabelSet{ "ticket_sys": pmod.LabelValue(k) }},
			"firing",
		}
		a2 := alert.Alert{pmod.Alert{
			Labels: pmod.LabelSet{ "ticket_sys": pmod.LabelValue("other") }},
			"firing",
		}
		abad := alert.Alert{pmod.Alert{
			Labels: pmod.LabelSet{ "ticket_sys": pmod.LabelValue("invalid") }},
			"firing",
		}
		ag.ComLabels = alert.LabelMap{}
		ag.Alerts = []alert.Alert{ a1, a1, a2, a1, alert.Alert{}, abad }
		got = obj.AgroupTSys(ag)
		assert.Equal(t,got,exp)
	}

	a := alert.Alert{}
	ag := alert.AlertGroup{}
	ag.Alerts = []alert.Alert{ a }
	exp := obj.Default
	assert.Equal(t,obj.AgroupTSys(ag),exp)

	ag.ComLabels = alert.LabelMap{ "ticket_sys": "invalid"  }
	assert.Equal(t,obj.AgroupTSys(ag),exp)

	aDef := alert.Alert{pmod.Alert{
		Labels: pmod.LabelSet{
			"ticket_sys": pmod.LabelValue(obj.Default.String()) }},
		"firing",
	}
	aGit := alert.Alert{pmod.Alert{
		Labels: pmod.LabelSet{
			"ticket_sys": pmod.LabelValue(TSysGitlab.String()) }},
		"firing",
	}
	aHpsm := alert.Alert{pmod.Alert{
		Labels: pmod.LabelSet{
			"ticket_sys": pmod.LabelValue(TSysHpsm.String()) }},
		"firing",
	}
	aMock := alert.Alert{pmod.Alert{
		Labels: pmod.LabelSet{
			"ticket_sys": pmod.LabelValue(TSysMock.String()) }},
		"firing",
	}
	aBad := alert.Alert{pmod.Alert{
		Labels: pmod.LabelSet{
			"ticket_sys": pmod.LabelValue("invalid") }},
		"firing",
	}
	ag.ComLabels = alert.LabelMap{}
	ag.Alerts = []alert.Alert{a,aDef,aGit,aHpsm,aMock,aBad}
	assert.Equal(t,obj.AgroupTSys(ag),exp)

	ag.Alerts = []alert.Alert{a,aHpsm,aGit,aHpsm,aHpsm,aMock,aBad}
	assert.Equal(t,obj.AgroupTSys(ag),TSysHpsm)
	promp.Unregister(obj.MetrTicketsGend)
	promp.Unregister(obj.MetrErrors)
}
*/
