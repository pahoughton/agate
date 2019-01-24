/* 2018-12-25 (cc) <paul4hough@gmail.com>
   process alert script remediation
*/
package proc

import (
	"fmt"
	"io/ioutil"
	"os"
	"os/exec"
	"path/filepath"

	"gopkg.in/yaml.v2"

	pmod "github.com/prometheus/common/model"
	promp "github.com/prometheus/client_golang/prometheus"
)


func (p *Proc)Script(node string, labels pmod.LabelSet) (string, error) {

	lfile, err := ioutil.TempFile("/tmp",node)
	if err != nil {
		return "", fmt.Errorf("ioutil.TempFile: %s",err.Error())
	}
	defer os.Remove(lfile.Name())

	lyml, err := yaml.Marshal(labels)
	if err != nil {
		return "", fmt.Errorf("yaml.Marshal - %s\n%v",err,labels)
	}
	if _, err := lfile.Write(lyml); err != nil {
		return "", fmt.Errorf("Write: %s",err.Error())
	}
	if err := lfile.Close(); err != nil {
		return "", fmt.Errorf("Close: %s",err.Error())
	}

	scriptfn := filepath.Join(p.ScriptsDir,string(labels["alertname"]))

	cmdargs := []string{node,lfile.Name()}

	cmdout, err := exec.Command(scriptfn,cmdargs...).CombinedOutput()

	var cmdstatus string

	if err != nil {
		cmdstatus = "error"
	} else {
		cmdstatus = "success"
	}

	tcom := fmt.Sprintf("command: %s %v",scriptfn,cmdargs)
	tcom += "results: " + cmdstatus + "\n"
	if err != nil {
		tcom += "cmd error: " + err.Error() + "\n"
	}
	tcom += "output:\n" + string(cmdout)

	if p.Debug {
		fmt.Printf("DEBUG: script %v\noutput: %s\n",cmdargs,cmdout)
	}

	p.ScriptsRun.With(
		promp.Labels{
			"script": string(labels["alertname"]),
			"status": cmdstatus,
		}).Inc()

	return tcom, err
}
