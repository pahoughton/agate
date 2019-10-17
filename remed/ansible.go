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
	pmod "github.com/prometheus/common/model"
)

func (r *Remed) AnsibleAvail(task string) bool {
	ardir := path.Join(r.playbookDir,"roles",task)
	finfo, err := os.Stat(ardir)
	return err == nil && finfo.IsDir()
}

func (r *Remed)Ansible(task string, node string, labels pmod.LabelSet) (string, error) {
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

	arole := "agate_role=" + task

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
		promp.Labels{"role": task,"status": cmdstatus}).Inc()

	return out, err
}
