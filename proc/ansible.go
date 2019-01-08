/* 2018-12-25 (cc) <paul4hough@gmail.com>
   process alert ansible remediation
*/
package proc

import (
	"fmt"
	"io/ioutil"
	"os"
	"os/exec"

	promp "github.com/prometheus/client_golang/prometheus"
)


func (p *Proc)Ansible(
	node	string,
	labels	map[string]string,
	tsys	string,
	tid		string) error {

	// create inventory file for ansible
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

	pbfile, err := ioutil.TempFile(p.PlaybookDir,node)
	if err != nil {
		return fmt.Errorf("ioutil.TempFile: %s",err.Error())
	}
	defer os.Remove(pbfile.Name())
	pbvars := "  vars:\n"
	for k, v := range labels {
		pbvars += "    " +k+": "+v+"\n"
	}
	pbcont := `---
- name: agate {{ agate_role }} remediation
  hosts: all
` + pbvars + `
  roles:
    - "{{ agate_role }}"
`
	if _, err := pbfile.WriteString(pbcont); err != nil {
		return fmt.Errorf("WriteString: %s",err.Error())
	}
	if err := pbfile.Close(); err != nil {
		return fmt.Errorf("Close: %s",err.Error())
	}

	cmdargs := []string{"-i", invfile.Name(),"-e"}

	arole := "agate_role=" + labels["ansible"]

	cmdargs = append(cmdargs,arole,pbfile.Name())

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
		if err = p.Ticket.AddTidComment(tsys,tid,tcom); err != nil {
			return fmt.Errorf("ticket comment - %s",err.Error())
		}
	}
	if p.Debug {
		fmt.Printf("DEBUG: ansible-playbook %v\noutput: %s\n",cmdargs,cmdout)
	}

	p.AnsiblePlays.With(
		promp.Labels{
			"role": labels["ansible"],
			"status": cmdstatus,
		}).Inc()

	return nil
}
