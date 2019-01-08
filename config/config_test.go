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

func TestLoadHpsm(t *testing.T) {

	var cfgExp = Config{
		ListenAddr:		":9201",
		BaseDir:		"/var/lib/agate",
		MaxDays:		15,
		EmailSmtp:		"localhost:25",
		EmailFrom:		"no-reply-agate@nowhere.none",
		GitlabUrl:		"https://gitlab.com/api/v4",
		GitlabToken:	"secret-token",
		HpsmUrl:		"https://hpsm/api/v3",
		HpsmUser:		"hpsm",
		HpsmPass:		"pass",
	}

	cfgfn := "testdata/config.good.yml"

	cfgGot, err := LoadFile(cfgfn)

	if err != nil {
		t.Errorf("LoadFile %s: %s",cfgfn,err)
	}

	ymlGot, err := yaml.Marshal(cfgGot)
	if err != nil {
		t.Fatalf("yaml.Marshal: %s",err)
	}

	ymlExp, err := yaml.Marshal(cfgExp)
	if err != nil {
		t.Fatalf("yaml.Marshal: %s",err)
	}
	if ! reflect.DeepEqual(ymlGot, ymlExp) {
		t.Fatalf("%s: unexpected diff:\n  got:\n%s\n  exp:\n%s\n"
			cfgFn,
			ymlGot,
			ymlExp)
	}
}
