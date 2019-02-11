/* 2019-01-07 (cc) <paul4hough@gmail.com>
   config validation
*/
package config

import (
	"strings"
	"testing"

	"gopkg.in/yaml.v2"
)

func TestLoadFull(t *testing.T) {

	var cfgExp = Config{
		ListenAddr:			":9201",
		TicketDefaultSys:	"gitlab",
		TicketDefaultGrp:	"user/project",
		CfgScriptsDir:		"/etc/agate-scripts",
		CfgPlaybookDir:		"/etc/agate-playbook",
		MaxDays:			15,
		CloseResolved:		true,
		// EmailSmtp:			"localhost:25",
		// EmailFrom:			"no-reply-agate@nowhere.none",
		GitlabURL:			"https://gitlab.com/api/v4",
		GitlabToken:		"secret-token",
		HpsmURL:			"https://hpsm/api/v3",
		HpsmCreateEp:		"incident2",
		HpsmUpdateEp:		"incident2",
		HpsmCloseEp:		"incident3",
		HpsmUser:			"hpsm",
		HpsmPass:			"pass",
		MockURL:			"http://localhost:9202/ticket",
		ScriptsDir:			"/etc/agate-scripts",
		PlaybookDir:		"/etc/agate-playbook",
	}

	cfgfn := "testdata/config.good.full.yml"

	cfgGot, err := LoadFile(cfgfn)

	if err != nil {
		t.Errorf("LoadFile %s: %s",cfgfn,err)
	}

	gotYml, err := yaml.Marshal(cfgGot)
	if err != nil {
		t.Fatalf("yaml.Marshal: %s",err)
	}

	expYml, err := yaml.Marshal(cfgExp)
	if err != nil {
		t.Fatalf("yaml.Marshal: %s",err)
	}
	gotLines := strings.Split(string(gotYml), "\n")
	expLines := strings.Split(string(expYml),"\n")

	for i, gv := range gotLines {
		if gv != expLines[i] {
			t.Fatalf("\n%s !=\n%s\nGOT:\n%s\nEXP:\n%s\n",
				gv,expLines[i],gotYml,expYml)
		}
	}

}

func TestLoadMin(t *testing.T) {

	var cfgExp = Config{
		ListenAddr:			":9201",
		TicketDefaultSys:	"gitlab",
		TicketDefaultGrp:	"user/project",
		ScriptsDir:			"testdata/scripts",
		PlaybookDir:		"testdata/playbook",
	}

	cfgfn := "testdata/config.good.min.yml"

	cfgGot, err := LoadFile(cfgfn)

	if err != nil {
		t.Errorf("LoadFile %s: %s",cfgfn,err)
	}

	gotYml, err := yaml.Marshal(cfgGot)
	if err != nil {
		t.Fatalf("yaml.Marshal: %s",err)
	}

	expYml, err := yaml.Marshal(cfgExp)
	if err != nil {
		t.Fatalf("yaml.Marshal: %s",err)
	}
	gotLines := strings.Split(string(gotYml), "\n")
	expLines := strings.Split(string(expYml),"\n")

	for i, gv := range gotLines {
		if gv != expLines[i] {
			t.Fatalf("\n%s !=\n%s\nGOT:\n%s\nEXP:\n%s\n",
				gv,expLines[i],gotYml,expYml)
		}
	}

}
