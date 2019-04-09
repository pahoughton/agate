/* 2019-01-07 (cc) <paul4hough@gmail.com>
   agate configuration
*/
package config

import (
	"io/ioutil"
	"path"
	"time"
	"gopkg.in/yaml.v2"

	"github.com/imdario/mergo"
)
type Global struct {
	Retry				time.Duration	`yaml:"retry,omitempty"`
	DataAge				uint			`yaml:"data-age,omitempty"`
}
type Remed struct {
	Parallel			uint			`yaml:"parallel,omitempty"`
	CfgScriptsDir		string			`yaml:"scripts-dir,omitempty"`
	CfgPlaybookDir		string			`yaml:"playbook-dir,omitempty"`
	// derived
	ScriptsDir			string
	PlaybookDir			string
}
type Email struct {
	Smtp	string	`yaml:"smtp,omitempty"`
	From	string	`yaml:"from,omitempty"`
	To		string	`yaml:"to,omitempty"`
}

type NSysGitlab struct {
	Url			string	`yaml:"url,omitempty"`
	Group		string	`yaml:"repo"`
	Token		string	`yaml:"token"`
}
type NSysHpsm struct {
	Url			string	`yaml:"url"`
	User		string	`yaml:"user"`
	Pass		string	`yaml:"pass"`
	CreateEp	string	`yaml:"create-ep"`
	UpdateEp	string	`yaml:"update-ep"`
	CloseEp		string	`yaml:"close-ep"`
	Group		string	`yaml:"workgroup"`
	Defaults	map[string]string `yaml:"defaults,omitempty"`
}
type NSysMock struct {
	Url				string	`yaml:"url"`
}

type NotifySys struct {
	Gitlab	NSysGitlab	`yaml:"gitlab,omitempty"`
	Hpsm	NSysHpsm	`yaml:"hpsm,omitempty"`
	Mock	NSysMock	`yaml:"mock,omitempty"`
}
type Notify struct {
	Default		string		`yaml:"default,omitempty"`
	Resolved	bool		`yaml:"close-resolved,omitempty"`
	Sys			NotifySys	`yaml:"systems,omitempty"`
}
type Config struct {
	Global		Global		`yaml:"global,omitempty"`
	Remed		Remed		`yaml:"remed,omitempty"`
	Email		Email		`yaml:"email,omitempty"`
	Notify		Notify		`yaml:"notify,omitempty"`
}

func New() (*Config) {
	return &Config {
		// defaults
		Global: Global{
			Retry: time.Duration(10 * time.Minute),
			DataAge: 15,
		},
		Remed: Remed{
			Parallel: 24,
			CfgScriptsDir: "scripts",
			CfgPlaybookDir: "playbook",
		},
		Notify: Notify{
			Default: "mock",
			Resolved: true,
			Sys: NotifySys{
				Gitlab: NSysGitlab{
					Url: "https://gitlab.com",
				},
				Mock: NSysMock{
					Url: "http://localhost:6102/ticket",
				},
			},
		},
	}
}

func Load(fn string) (*Config, error) {
	cfg := New()
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

	if len(cfg.Remed.CfgPlaybookDir) > 0 {
		cfg.Remed.PlaybookDir = cfg.Remed.CfgPlaybookDir
	} else {
		cfg.Remed.PlaybookDir = "playbook"
	}
	if ! path.IsAbs(cfg.Remed.PlaybookDir) {
		cfg.Remed.PlaybookDir = path.Join(baseDir,cfg.Remed.PlaybookDir)
	}
	if len(cfg.Remed.CfgScriptsDir) > 0 {
		cfg.Remed.ScriptsDir = cfg.Remed.CfgScriptsDir
	} else {
		cfg.Remed.ScriptsDir = "scripts"
	}
	if ! path.IsAbs(cfg.Remed.ScriptsDir) {
		cfg.Remed.ScriptsDir = path.Join(baseDir,cfg.Remed.ScriptsDir)
	}

	return cfg, nil
}
