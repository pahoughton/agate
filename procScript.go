/* 2018-12-25 (cc) <paul4hough@gmail.com>
   process alert script remediation
*/
package main

import (
	"fmt"
	"os/exec"
	"strings"
	"path/filepath"

	promp "github.com/prometheus/client_golang/prometheus"
)


func procScript(a *AmgrAlert, tid string) error {

	node := strings.Split(a.Labels["instance"],":")[0]

	scriptfn := filepath.Join(*args.ScriptDir,a.Labels["script"])

	cmdargs := []string{node}

	if _, ok := a.Labels["script_arg"]; ok {
		cmdargs = append(cmdargs, a.Labels["script_arg"])
	}

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
		if err = addTicketComment(tid,tcom); err != nil {
			fmt.Println("ERROR: ticket comment - ",err.Error())
		}
	}
	if *args.Debug {
		fmt.Printf("DEBUG: ansible-playbook %v\noutput: %s\n",cmdargs,cmdout)
	}

	prom.ScriptsRun.With(
		promp.Labels{
			"script": a.Labels["script"],
			"status": a.Status,
		})
	return nil
}
