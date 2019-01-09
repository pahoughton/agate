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

	promp "github.com/prometheus/client_golang/prometheus"
)


func (p *Proc)Script(
	node	string,
	labels	map[string]string,
	tsys	string,
	tid		string) error {

	lfile, err := ioutil.TempFile("/tmp",node)
	if err != nil {
		return fmt.Errorf("ioutil.TempFile: %s",err.Error())
	}
	defer os.Remove(lfile.Name())

	lyml, err := yaml.Marshal(labels)
	if err != nil {
		return fmt.Errorf("yaml.Marshal - %s\n%v",err,labels)
	}
	if _, err := lfile.Write(lyml); err != nil {
		return fmt.Errorf("Write: %s",err.Error())
	}
	if err := lfile.Close(); err != nil {
		return fmt.Errorf("Close: %s",err.Error())
	}

	scriptfn := filepath.Join(p.ScriptsDir,labels["script"])

	cmdargs := []string{node,lfile.Name()}

	cmdout, err := exec.Command(scriptfn,cmdargs...).CombinedOutput()

	var cmdstatus string

	if err != nil {
		cmdstatus = "error"
	} else {
		cmdstatus = "success"
	}
	if len(tid) > 0 {
		tcom := fmt.Sprintf("command: %s %v",scriptfn,cmdargs)
		tcom += "results: " + cmdstatus + "\n"
		if err != nil {
			tcom += "cmd error: " + err.Error() + "\n"
		}
		tcom += "output:\n" + string(cmdout)
		if err = p.Ticket.AddTidComment(tsys,tid,tcom); err != nil {
			return fmt.Errorf("ticket comment - %s",err.Error())
		}
	}
	if p.Debug {
		fmt.Printf("DEBUG: script %v\noutput: %s\n",cmdargs,cmdout)
	}

	p.ScriptsRun.With(
		promp.Labels{
			"script": labels["script"],
			"status": cmdstatus,
		}).Inc()

	return nil
}
