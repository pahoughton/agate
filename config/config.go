/* 2019-01-07 (cc) <paul4hough@gmail.com>
   agate configuration
*/
package config

import (
	"io/ioutil"
	"path"
	"gopkg.in/yaml.v2"

	"github.com/imdario/mergo"
)
type Global struct {
	Listen				string	`yaml:"listen,omitempty"`
	DataAge				uint	`yaml:"data-age,omitempty"`
	CfgScriptsDir		string	`yaml:"scripts-dir,omitempty"`
	CfgPlaybookDir		string	`yaml:"playbook-dir,omitempty"`
	// derived
	ScriptsDir			string
	PlaybookDir			string
}
type Email struct {
	Smtp	string	`yaml:"smtp,omitempty"`
	From	string	`yaml:"from,omitempty"`
	To		string	`yaml:"to,omitempty"`
}

type TSysGitlab struct {
	Url			string	`yaml:"url,omitempty"`
	Group		string	`yaml:"repo"`
	Token		string	`yaml:"token"`
}
type TSysHpsm struct {
	Url			string	`yaml:"url"`
	User		string	`yaml:"user"`
	Pass		string	`yaml:"pass"`
	CreateEp	string	`yaml:"create-ep"`
	UpdateEp	string	`yaml:"update-ep"`
	CloseEp		string	`yaml:"close-ep"`
	Group		string	`yaml:"workgroup"`
	Defaults	map[string]string `yaml:"defaults,omitempty"`
}
type TSysMock struct {
	Url				string	`yaml:"url"`
}

type TicketSys struct {
	Gitlab	TSysGitlab	`yaml:"gitlab,omitempty"`
	Hpsm	TSysHpsm	`yaml:"hpsm,omitempty"`
	Mock	TSysMock	`yaml:"mock,omitempty"`
}
type Ticket struct {
	Default		string		`yaml:"default,omitempty"`
	Resolved	bool		`yaml:"close-resolved,omitempty"`
	Sys			TicketSys	`yaml:"systems,omitempty"`
}
type Config struct {
	Global		Global		`yaml:"global,omitempty"`
	Email		Email		`yaml:"email,omitempty"`
	Ticket		Ticket		`yaml:"ticket-sys,omitempty"`
}

func New() (*Config) {
	return &Config {
		// defaults
		Global: Global{
			Listen: "6101",
			DataAge: 15,
			CfgScriptsDir: "scripts",
			CfgPlaybookDir: "playbook",
		},
		Ticket: Ticket{
			Default: "mock",
			Resolved: true,
			Sys: TicketSys{
				Gitlab: TSysGitlab{
					Url: "https://gitlab.com",
				},
				Mock: TSysMock{
					Url: "http://localhost:6102",
				},
			},
		},
	}
}

func (cfg *Config)Load(fn string) (*Config, error) {
	ycfg := &Config{}
	if dat, err := ioutil.ReadFile(fn); err == nil {
		if 	err = yaml.UnmarshalStrict(dat, ycfg); err == nil {
			if err = mergo.Merge(cfg,ycfg,mergo.WithOverride); err != nil {
				return nil, err
			}
		} else {
			return nil, err
		}
	} else {
		return nil, err
	}

	baseDir := path.Dir(fn)

	if len(cfg.Global.CfgPlaybookDir) > 0 {
		cfg.Global.PlaybookDir = cfg.Global.CfgPlaybookDir
	} else {
		cfg.Global.PlaybookDir = "playbook"
	}
	if ! path.IsAbs(cfg.Global.PlaybookDir) {
		cfg.Global.PlaybookDir = path.Join(baseDir,cfg.Global.PlaybookDir)
	}
	if len(cfg.Global.CfgScriptsDir) > 0 {
		cfg.Global.ScriptsDir = cfg.Global.CfgScriptsDir
	} else {
		cfg.Global.ScriptsDir = "scripts"
	}
	if ! path.IsAbs(cfg.Global.ScriptsDir) {
		cfg.Global.ScriptsDir = path.Join(baseDir,cfg.Global.ScriptsDir)
	}

	return cfg, nil
}
