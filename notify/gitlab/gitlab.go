/* 2019-01-07 (cc) <paul4hough@gmail.com>
   gitlab issue interface
*/
package gitlab

import (
	"fmt"
	"strconv"
	"strings"

	gl "github.com/xanzy/go-gitlab"

	"github.com/pahoughton/agate/config"
)

type Gitlab struct {
	name	string
	grp		string
	debug	bool
	c		*gl.Client
}

func New(cfg config.NSysGitlab,name string,dbg bool) *Gitlab {
	g := &Gitlab{
		name:	name,
		grp:	cfg.Group,
		debug:	dbg,
		c:		gl.NewClient(nil, cfg.Token),
	}
	g.c.SetBaseURL(cfg.Url)
	return g
}

func (g *Gitlab)Group() string {
	return g.grp
}

func (g *Gitlab)Create(prj, title, desc string, ) ([]byte, error) {

	i, resp, err := g.c.Issues.CreateIssue(prj,&gl.CreateIssueOptions{
		Title: gl.String(title),
		Description: gl.String("```\n"+desc+"\n```\n"),
	})
	if resp != nil && resp.Body != nil {
		defer resp.Body.Close()
	}

	if err != nil {
		return nil, err
	}
	if g.debug {
		fmt.Printf("gitlab.CreateIssue: ret issue: %v\n",i)
	}
	return []byte(fmt.Sprintf("%s:%d",prj,i.IID)), nil
}

func (g *Gitlab)Update(id []byte, cmt string) error {

	nida := strings.Split(string(id),":")
	prj := nida[0]
	issue, err := strconv.Atoi(nida[1])
	if err != nil {
		return fmt.Errorf("atoi: %s - %s",nida[1],err)
	}
	if g.debug {
		fmt.Printf("gitlab.AddComment: nid '%s' nida '%v' nida0 '%s' nida1 '%s' prj '%s' issue '%d'\n",
			string(id),
			nida,
			nida[0],
			nida[1],
			prj,
			issue)
	}
	_, resp, err := g.c.Notes.CreateIssueNote(
		prj,
		issue,
		&gl.CreateIssueNoteOptions{
			Body: gl.String("```\n"+cmt+"\n```\n"),
		})

	if resp != nil && resp.Body != nil {
		defer resp.Body.Close()
	}

	if err != nil {
		return err
	}
	return nil
}

func (g *Gitlab)Close(id []byte, cmt string) error {

	if len(cmt) > 0 {
		g.Update(id,cmt)
	}

	nida := strings.Split(string(id),":")
	prj := nida[0]
	issue, err := strconv.Atoi(nida[1])
	if err != nil {
		return fmt.Errorf("atoi: %s - %s",nida[1],err)
	}

	_, resp, err := g.c.Issues.UpdateIssue(
		prj,
		issue,
		&gl.UpdateIssueOptions{
			StateEvent: gl.String("close"),
		})

	if resp != nil && resp.Body != nil {
		defer resp.Body.Close()
	}

	if err != nil {
		return err
	}
	return nil
}
