/* 2018-12-25 (cc) <paul4hough@gmail.com>
   process alert script remediation
*/
package main

import (
	"fmt"
	"os"
	"os/exec"
	"strings"
	"path/filepath"

	log "github.com/sirupsen/logrus"
)


func procScript(a *AmgrAlert) {

	node := strings.Split(a.Labels["instance"],":")[0]

	scriptfn := filepath.Join(*scriptDir,a.Labels["script"])

	cmdargs := []string{node}

	if _, ok := a.Labels["script_arg"]; ok {
		cmdargs = append(cmdargs, a.Labels["script_arg"])
	}

	cmdout, err := exec.Command(scriptfn,cmdargs...).CombinedOutput()

	a.RemedOut = string(cmdout)

	if err != nil {
		fmt.Fprintf(os.Stderr,"script %s %v\noutput:\n%s\n",
			scriptfn,
			cmdargs,
			cmdout)
		log.Error(err)
		a.Status = "script failed"
		createTicket(a)
	} else {
		a.Status = "remediated"
		createTicket(a)
	}

	// fixme debug
	fmt.Fprintf(os.Stderr,"script out\n%s\n",cmdout)
	log.Debug("script " + a.Labels["script"] + " complete")

	scriptProcd.Inc()
}
