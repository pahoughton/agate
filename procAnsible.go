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

	promp "github.com/prometheus/client_golang/prometheus"
)


func procAnsible(a *AmgrAlert) {

	// create inventory file for ansible
	node := strings.Split(a.Labels["instance"],":")[0]
	invfile, err := ioutil.TempFile("/tmp", "inventory")
	if err != nil {
		fmt.Println("FATAL-ioutil.TempFile: %s",err.Error())
		os.Exit(2)
	}
	defer os.Remove(invfile.Name())
	if _, err := invfile.WriteString(node + "\n"); err != nil {
		fmt.Println("FATAL-WriteString: %s",err.Error())
		os.Exit(2)
	}
	if err := invfile.Close(); err != nil {
		fmt.Println("FATAL-Close: %s",err.Error())
		os.Exit(2)
	}

	pbookfn := filepath.Join(*args.PlaybookDir,a.Labels["ansible"] + ".yml")

	cmdargs := []string{"-i", invfile.Name()}

	if _, ok := a.Labels["ansible_vars"]; ok {
		cmdargs = append(cmdargs, "-e", a.Labels["ansible_vars"])
	}
	cmdargs = append(cmdargs, pbookfn)

	cmdout, err := exec.Command("ansible-playbook",cmdargs...).CombinedOutput()

	a.RemedOut = string(cmdout)

	if err != nil {
		a.Status = "ansible failed"
		fmt.Printf("ERROR-ansible-%s:%s\n%s\n",
			a.Labels["ansible"],err.Error(),cmdout)
		createTicket(a)
	} else {
		a.Status = "remediated"
		createTicket(a)
	}
	prom.AnsiblePlays.With(
		promp.Labels{
			"playbook": a.Labels["ansible"],
			"status": a.Status,
		})
	if *args.Debug {
		fmt.Printf("ansible out\n%s\n",cmdout)
	}
}
