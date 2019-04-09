/* 2018-12-25 (cc) <paul4hough@gmail.com>

process alert ansible remediation
- create inventory file w/ node
- create playbook with variables from labels
- run ansible role
- return output
*/
package remed

import (
	"fmt"
	"io/ioutil"
	"os"
	"os/exec"
	"path"

	promp "github.com/prometheus/client_golang/prometheus"
	"github.com/pahoughton/agate/amgr/alert"
)

func (r *Remed) AnsibleAvail(labels alert.LabelSet) bool {
	aname, ok := labels["alertname"]
	if ok {
		ardir := path.Join(r.playbookDir,"roles",string(aname))
		finfo, err := os.Stat(ardir)
		return err == nil && finfo.IsDir()
	} else {
		return ok
	}
}

func (r *Remed)Ansible( node string, labels alert.LabelSet) (string, error) {
	taname, ok := labels["alertname"]
	if ! ok {
		return "", r.errorf("no alertname label: Ansible(%s,%v)",node,labels)
	}
	aname := string(taname)

	// create inventory file for ansible
	invfile, err := ioutil.TempFile("/tmp", "inventory")
	if err != nil {
		return "", r.errorf("ioutil.TempFile: %s",err.Error())
	}
	defer os.Remove(invfile.Name())
	if _, err := invfile.WriteString(node + "\n"); err != nil {
		return "", r.errorf("WriteString: %s",err.Error())
	}
	if err := invfile.Close(); err != nil {
		return "", r.errorf("Close: %s",err.Error())
	}

	// create playbook
	pbfile, err := ioutil.TempFile(r.playbookDir,node)
	if err != nil {
		return "", r.errorf("ioutil.TempFile: %s",err.Error())
	}
	defer os.Remove(pbfile.Name())
	pbvars := "  vars:\n"
	for k, v := range labels {
		pbvars += "    " +string(k)+": "+string(v)+"\n"
	}
	pbcont := `---
- name: agate {{ agate_role }} remediation
  hosts: all
` + pbvars + `
  roles:
    - "{{ agate_role }}"
`
	if r.debug {fmt.Printf("proc.Ansible-playbook:\n%s\n",pbcont)}

	if _, err := pbfile.WriteString(pbcont); err != nil {
		return "", r.errorf("WriteString: %s",err.Error())
	}
	if err := pbfile.Close(); err != nil {
		return "", r.errorf("Close: %s",err.Error())
	}

	arole := "agate_role=" + aname

	cmdargs := []string{"-i",invfile.Name(),"-e",arole,pbfile.Name()}

	cmdout, err := exec.Command("ansible-playbook",cmdargs...).CombinedOutput()

	var cmdstatus string

	if err != nil {
		cmdstatus = "error"
	} else {
		cmdstatus = "success"
	}

	out := fmt.Sprintf("command: anisble-playbook %v",cmdargs)
	out += "results: " + cmdstatus + "\n"
	if err != nil {
		out += "cmd error: " + err.Error() + "\n"
	}
	out += "output:\n" + string(cmdout)

	if r.debug {
		fmt.Printf("DEBUG: ansible-playbook %v\noutput: %s\n",cmdargs,cmdout)
	}

	r.metrics.ansible.With(
		promp.Labels{"role": aname,"status": cmdstatus}).Inc()

	return out, err
}
