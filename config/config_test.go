/* 2019-01-07 (cc) <paul4hough@gmail.com>
   config validation
*/
package config

import (
	"encoding/json"
	"net/url"
	"reflect"
	"regexp"
	"strings"
	"testing"
	"time"

	"github.com/stretchr/testify/require"
	"gopkg.in/yaml.v2"
)

func TestLoadFull(t *testing.T) {

	var cfgExp = Config{
		ListenAddr:			":9201",
		BaseDir:			"/var/lib/agate",
		MaxDays:			15,
		TicketDefaultSys:	"gitlab",
		TicketDefaultGrp:	"user/project",
		CloseResolved:		true,
		EmailSmtp:			"localhost:25",
		EmailFrom:			"no-reply-agate@nowhere.none",
		GitlabUrl:			"https://gitlab.com/api/v4",
		GitlabToken:		"secret-token",
		HpsmUrl:			"https://hpsm/api/v3",
		HpsmUser:			"hpsm",
		HpsmPass:			"pass",
		MockURL:			"http://localhost:9202/ticket",
		DataDir:			"/var/lib/agate/data"
		PlaybookDir:		"/var/lib/agate/playbook"
		ScriptsDir:			"/var/lib/agate/scripts"
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
