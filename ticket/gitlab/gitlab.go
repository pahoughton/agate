/* 2019-01-07 (cc) <paul4hough@gmail.com>
   gitlab issue interface
*/
package gitlab

import (
	"errors"
	"fmt"
	"strconv"
	"strings"

	gl "github.com/xanzy/go-gitlab"
)

type Gitlab struct {
	Debug		bool
	defaultPrj	string
	c			*gl.Client
}

func New(url string, token string, dprj string) *Gitlab {
	g := &Gitlab{
		defaultPrj: dprj,
		c:			gl.NewClient(nil, token),
	}
	g.c.SetBaseURL(url)
	return g
}

func (g *Gitlab)CreateIssue(
	prj		string,
	title	string,
	desc	string) (string, error) {

	if len(prj) < 1 {
		if len(g.defaultPrj) < 1 {
			return "", errors.New("no gitlab project")
		}
		prj = g.defaultPrj
	}
	i, resp, err := g.c.Issues.CreateIssue(prj,&gl.CreateIssueOptions{
		Title: gl.String(title),
		Description: gl.String("```\n"+desc+"\n```\n"),
	})
	if err != nil {
		return "", fmt.Errorf("gl.CreateIssue: %s\nresp:\n%v",err,resp)
	}
	if g.Debug {
		fmt.Printf("gitlab.CreateIssue: ret issue: %v\n",i)
	}
	return fmt.Sprintf("%s:%d",prj,i.IID), nil
}

func (g *Gitlab)AddComment(tid string, cmt string) error {

	tida := strings.Split(tid,":")
	prj := tida[0]
	issue, err := strconv.Atoi(tida[1])
	if err != nil {
		return fmt.Errorf("atoi: %s - %s",tida[1],err)
	}
	if g.Debug {
		fmt.Printf("gitlab.AddComment: tid '%s' tida '%v' tida0 '%s' tida1 '%s' prj '%s' issue '%d'\n",
			tid,
			tida,
			tida[0],
			tida[1],
			prj,
			issue)
	}
	_, resp, err := g.c.Notes.CreateIssueNote(
		prj,
		issue,
		&gl.CreateIssueNoteOptions{
			Body: gl.String("```\n"+cmt+"\n```\n"),
		})

	if err != nil {
		return fmt.Errorf("gl.CreateIssueNote: %s\nresp:\n%v",err,resp)
	}
	return nil
}

func (g *Gitlab)Close(tid string) error {

	tida := strings.Split(tid,":")
	prj := tida[0]
	issue, err := strconv.Atoi(tida[1])
	if err != nil {
		return fmt.Errorf("atoi: %s - %s",tida[1],err)
	}

	_, resp, err := g.c.Issues.UpdateIssue(
		prj,
		issue,
		&gl.UpdateIssueOptions{
			StateEvent: gl.String("close"),
		})

	if err != nil {
		return fmt.Errorf("gl.CreateIssueNote: %s\nresp:\n%v",err,resp)
	}
	return nil
}
