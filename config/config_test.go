/* 2019-01-07 (cc) <paul4hough@gmail.com>
   config validation
*/
package config

import (
//	"fmt"
//	"strings"
	"testing"
	"gopkg.in/yaml.v2"

	"github.com/stretchr/testify/assert"
)
func TestNewConfig(t *testing.T) {
	c := New()
	got, err := yaml.Marshal(c)
	assert.Nil(t,err)
	// print(string(got))

	exp := `global:
  listen: "6101"
  data-age: 15
  scripts-dir: scripts
  playbook-dir: playbook
  scriptsdir: ""
  playbookdir: ""
ticket-sys:
  default: mock
  close-resolved: true
  systems:
    gitlab:
      url: https://gitlab.com
      repo: ""
      token: ""
    mock:
      url: http://localhost:6102
`
	assert.Equal(t,exp,string(got))
}

func TestLoadMissing(t *testing.T) {
	c, err := Load("not-a-file")
	assert.Error(t,err)
	assert.Nil(t,c)
}
func TestLoadBad(t *testing.T) {
	c, err := Load("testdata/bad.yml")
	assert.Error(t,err)
	assert.Nil(t,c)
}
func TestLoadBadCont(t *testing.T) {
	c, err := Load("testdata/bad-cont.yml")
	assert.Error(t,err)
	assert.Nil(t,c)
}
func TestLoadMin(t *testing.T) {
	got, err := Load("testdata/good-min.yml")
	assert.Nil(t,err)
	assert.NotNil(t,got)
	exp := New()
	// load sets dirs
	exp.Global.ScriptsDir = "testdata/scripts"
	exp.Global.PlaybookDir = "testdata/playbook"

	assert.Equal(t,exp,got)
}
func TestLoadFull(t *testing.T) {
	exp :=  &Config{
		Global: Global{
			Listen: "127.0.0.1:4464",
			DataAge: 30,
			CfgScriptsDir: "/sdiff",
			CfgPlaybookDir: "/pdiff",
			ScriptsDir: "/sdiff",
			PlaybookDir: "/pdiff",
		},
		Email: Email{
			Smtp: "localhost:25",
			From: "agate@nowhere",
			To: "invalid",
		},
		Ticket: Ticket{
			Resolved: true,
			Default: "gitlab",
			Sys: TicketSys{
				Gitlab: TSysGitlab{
					Url: "https://mylab",
					Group: "paul",
					Token: "secret-sauce",
				},
				Hpsm: TSysHpsm{
					Url:		"https://myhpsm/api",
					User:		"paul",
					Pass:		"secret-sauce",
					CreateEp:	"create",
					UpdateEp:	"update",
					CloseEp:	"close",
					Group:		"team",
					Defaults:	map[string]string{
						"urgency": "now",
						"assignee": "you",
					},
				},
				Mock: TSysMock{
					Url: "http://cbed:1234/abc",
				},
			},
		},
	}
	if _, err := yaml.Marshal(exp); err == nil {
		// print(string(yml))
	} else {
		assert.Nil(t,err)
	}
	got, err := Load("testdata/good-full.yml")
	assert.Nil(t,err)
	assert.NotNil(t,got)

	assert.Equal(t,exp,got)
}
