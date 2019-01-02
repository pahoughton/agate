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

	promp "github.com/prometheus/client_golang/prometheus"
)


func procAnsible(a *AmgrAlert, tid string) error {

	// create inventory file for ansible
	node := strings.Split(a.Labels["instance"],":")[0]
	invfile, err := ioutil.TempFile("/tmp", "inventory")
	if err != nil {
		return fmt.Errorf("ioutil.TempFile: %s",err.Error())
	}
	defer os.Remove(invfile.Name())
	if _, err := invfile.WriteString(node + "\n"); err != nil {
		return fmt.Errorf("WriteString: %s",err.Error())
	}
	if err := invfile.Close(); err != nil {
		return fmt.Errorf("Close: %s",err.Error())
	}

	cmdargs := []string{"-i", invfile.Name(),"-e"}

	avars := "agate_role=" + a.Labels["ansible"]

	if _, ok := a.Labels["ansible_vars"]; ok {
		avars += " " + a.Labels["ansible_vars"]
	}
	cmdargs = append(cmdargs,avars,*args.Playbook)

	cmdout, err := exec.Command("ansible-playbook",cmdargs...).CombinedOutput()

	var cmdstatus string

	if err != nil {
		cmdstatus = "error"
	} else {
		cmdstatus = "success"
	}
	if len(tid) > 0 {
		tcom := fmt.Sprintf("command: anisble-playbook %v",cmdargs)
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

	prom.AnsiblePlays.With(
		promp.Labels{
			"role": a.Labels["ansible"],
			"status": cmdstatus,
		})
	return nil
}
