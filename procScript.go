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


func procScript(a *AmgrAlert) {

	node := strings.Split(a.Labels["instance"],":")[0]

	scriptfn := filepath.Join(*args.ScriptDir,a.Labels["script"])

	cmdargs := []string{node}

	if _, ok := a.Labels["script_arg"]; ok {
		cmdargs = append(cmdargs, a.Labels["script_arg"])
	}

	cmdout, err := exec.Command(scriptfn,cmdargs...).CombinedOutput()

	a.RemedOut = string(cmdout)

	if err != nil {
		a.Status = "script failed"
		fmt.Printf("ERROR-script: %s %s\n%s\n",
			cmdargs,
			err.Error(),
			cmdout)
		createTicket(a)
	} else {
		a.Status = "remediated"
		createTicket(a)
	}
	prom.ScriptsRun.With(
		promp.Labels{
			"script": a.Labels["script"],
			"status": a.Status,
		})

	if *args.Debug {
		fmt.Printf("script out\n%s\n",cmdout)
	}
}
