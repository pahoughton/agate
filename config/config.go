/* 2019-01-07 (cc) <paul4hough@gmail.com>
   agate configuration
*/
package config

import (
	"fmt"
	"io/ioutil"
	"os"
	"path"
	"gopkg.in/yaml.v2"
)

type Config struct {
	ListenAddr			string	`yaml:"listen-addr"`
	BaseDir				string	`yaml:"base-dir,omitempty"`
	MaxDays				uint	`yaml:"max-days,omitempty"`
	TicketDefaultSys	string	`yaml:"ticket-default-sys"`
	TicketDefaultGrp	string	`yaml:"ticket-default-grp"`
	CloseResolved		bool	`yaml:"close-resolved,omitempty"`
	EmailSmtp			string	`yaml:"email-smtp,omitempty"`
	EmailFrom			string	`yaml:"email-from,omitempty"`
	GitlabURL			string	`yaml:"gitlab-url,omitempty"`
	GitlabToken			string	`yaml:"gitlab-token,omitempty"`
	HpsmURL				string	`yaml:"hpsm-base-url,omitempty"`
	HpsmCreateEp		string	`yaml:"hpsm-create-ep,omitempty"`
	HpsmUpdateEp		string	`yaml:"hpsm-update-ep,omitempty"`
	HpsmCloseEp			string	`yaml:"hpsm-close-ep,omitempty"`
	HpsmUser			string	`yaml:"hpsm-user,omitempty"`
	HpsmPass			string	`yaml:"hpsm-pass,omitempty"`
	MockURL				string	`yaml:"mock-ticket-url,omitempty"`
	// derived
	DataDir				string
	PlaybookDir			string
	ScriptsDir			string
}

func LoadFile(fn string) (*Config, error) {
	dat, err := ioutil.ReadFile(fn)
	if err != nil {
		return nil, err
	}
	cfg := &Config{}
	err = yaml.UnmarshalStrict(dat, cfg)
	if err != nil {
		return nil, err
	}
	if len(cfg.BaseDir) < 1 {
		cfg.BaseDir = "/var/lib/agate"
	}

	cfg.DataDir = path.Join(cfg.BaseDir,"data")
	if err = os.MkdirAll(cfg.DataDir,0775); err != nil {
		return nil, fmt.Errorf("FATAL: %s - %s",cfg.DataDir,err.Error())
	}

	cfg.PlaybookDir = path.Join(cfg.BaseDir,"playbook")
	rDir := path.Join(cfg.PlaybookDir,"roles")
	if err = os.MkdirAll(rDir,0775); err != nil {
		return nil, fmt.Errorf("FATAL: %s - %s",rDir,err.Error())
	}

	cfg.ScriptsDir = path.Join(cfg.BaseDir,"scripts")
	if err = os.MkdirAll(cfg.ScriptsDir,0775); err != nil {
		return nil, fmt.Errorf("FATAL: %s - %s",cfg.ScriptsDir,err.Error())
	}

	return cfg, nil
}
