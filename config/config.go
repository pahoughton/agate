/* 2019-01-07 (cc) <paul4hough@gmail.com>
   agate configuration
*/
package config

import (
	"io/ioutil"
	"path"
	"gopkg.in/yaml.v2"
)

type Config struct {
	ListenAddr			string	`yaml:"listen-addr"`
	TicketDefaultSys	string	`yaml:"ticket-default-sys"`
	TicketDefaultGrp	string	`yaml:"ticket-default-grp"`
	CfgScriptsDir		string	`yaml:"scripts-dir,omitempty"`
	CfgPlaybookDir		string	`yaml:"playbook-dir,omitempty"`
	MaxDays				uint	`yaml:"max-days,omitempty"`
	CloseResolved		bool	`yaml:"close-resolved,omitempty"`
	// EmailSmtp			string	`yaml:"email-smtp,omitempty"`
	// EmailFrom			string	`yaml:"email-from,omitempty"`
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
	ScriptsDir			string
	PlaybookDir			string
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
	baseDir := path.Dir(fn)

	if len(cfg.CfgPlaybookDir) > 0 {
		cfg.PlaybookDir = cfg.CfgPlaybookDir
	} else {
		cfg.PlaybookDir = "playbook"
	}
	if ! path.IsAbs(cfg.PlaybookDir) {
		cfg.PlaybookDir = path.Join(baseDir,cfg.PlaybookDir)
	}
	if len(cfg.CfgScriptsDir) > 0 {
		cfg.ScriptsDir = cfg.CfgScriptsDir
	} else {
		cfg.ScriptsDir = "scripts"
	}
	if ! path.IsAbs(cfg.ScriptsDir) {
		cfg.ScriptsDir = path.Join(baseDir,cfg.ScriptsDir)
	}

	return cfg, nil
}
