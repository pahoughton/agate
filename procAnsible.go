/* 2018-12-25 (cc) <paul4hough@gmail.com>
   process alert ansible remediation
*/
package main

import (
	"fmt"
	"io/ioutil"
	"os"
	"os/exec"
	"strings"
	"path/filepath"

	log "github.com/sirupsen/logrus"
)


func procAnsible(a *AmgrAlert) {

	// create inventory file for ansible
	node := strings.Split(a.Labels["instance"],":")[0]
	invfile, err := ioutil.TempFile("/tmp", "inventory")
	if err != nil {
		log.Fatal(err)
	}
	defer os.Remove(invfile.Name())
	if _, err := invfile.WriteString(node + "\n"); err != nil {
		log.Fatal(err)
	}
	if err := invfile.Close(); err != nil {
		log.Fatal(err)
	}

	pbookfn := filepath.Join(*pbookDir,a.Labels["ansible"] + ".yml")

	aout, err := exec.Command(
		"echo",
		"ansible-playbook",
		"-i",invfile.Name(),
		pbookfn).
			CombinedOutput()

	if err != nil {
		fmt.Fprintf(os.Stderr,"ansible out\n%s\n",aout)
		log.Fatal(err)
	}

	// fixme debug
	fmt.Fprintf(os.Stderr,"ansible out\n%s\n",aout)
	log.Debug("ansible " + a.Labels["ansible"] + " complete")

	ansibleProcd.Inc()
}
