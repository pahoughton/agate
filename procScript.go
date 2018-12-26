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

	cleansfn := strings.Replace(a.Labels["script"],"/","-",-1)
	scriptfn := filepath.Join(*scriptDir,cleansfn)

	aout, err := exec.Command(
		"echo",
		scriptfn,
		node).
			CombinedOutput()

	if err != nil {
		fmt.Fprintf(os.Stderr,"script out\n%s\n",aout)
		log.Fatal(err)
	}

	// fixme debug
	fmt.Fprintf(os.Stderr,"script out\n%s\n",aout)
	log.Debug("script " + a.Labels["script"] + " complete")

	scriptProcd.Inc()
}
