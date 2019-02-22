/* 2019-02-14 (cc) <paul4hough@gmail.com>
   run alert remediation
*/
package amgr

func (am *Amgr)Fix(a *alert.Alert,tid string,multi bool) {

	aname = a.Name()
	node = a.Node()

	ardir := path.Join(am.proc.PlaybookDir,"roles",aname)
	finfo, err := os.Stat(ardir)
	if err == nil && finfo.IsDir() {
		emsg := ""
		out, err := am.proc.Ansible(node,a.Labels)
		if err != nil {
			emsg = "ERROR: " + err.Error() + "\n"
			am.Error("ansible - " + err.Error() + "\n")
		}
		tcom := "ansible remediation results\n" + emsg + out

		if err = am.ticket.Comment(a,tid,tcom,multi); err != nil {
			am.Error(fmt.Sprintf("ticket add comment: %s\n%s",err,tcom))
		}
	}

	sfn := path.Join(am.proc.ScriptsDir,aname)
	finfo, err = os.Stat(sfn)
	if err == nil && (finfo.Mode() & 0111) != 0 {
		emsg := ""
		out, err := am.proc.Script(node,a.Labels)
		if err != nil {
			emsg = "ERROR: " + err.Error() + "\n"
			am.Error("script - " + err.Error() + "\n")
		}
		tcom := "script remediation results\n" + emsg + out

		if err = am.ticket.Comment(a,tid,tcom,multi); err != nil {
			am.Error(fmt.Sprintf("ticket add comment: %s\n%s",err,tcom))
		}
	}
}
