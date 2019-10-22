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
	"github.com/pahoughton/agate/notify/note"
)

type Gitlab struct {
	name	string
	grp		string
	debug	bool
	c		*gl.Client
}

func New(name string, cfg config.NSysGitlab,dbg bool) *Gitlab {
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
func (self *Gitlab) Name() string {
	return self.name
}

func (g *Gitlab) Create(grp string, note note.Note, remcnt int) ([]byte, error) {

	i, resp, err := g.c.Issues.CreateIssue(grp,&gl.CreateIssueOptions{
		Title: gl.String(note.Title()),
		Description: gl.String("```\n"+note.Desc()+"\n```\n"),
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
	return []byte(fmt.Sprintf("%s:%d",grp,i.IID)), nil
}

func (g *Gitlab)Update(note note.Note, cmt string) (bool,error) {

	nida := strings.Split(string(note.Nid),":")
	prj := nida[0]
	issue, err := strconv.Atoi(nida[1])
	if err != nil {
		return false, fmt.Errorf("atoi: %s - %s",nida[1],err)
	}
	if g.debug {
		fmt.Printf("gitlab.AddComment: nid '%s' nida '%v' nida0 '%s' nida1 '%s' prj '%s' issue '%d'\n",
			string(note.Nid),
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
		return false,err
	}
	return false,nil
}

func (g *Gitlab)Close(note note.Note, cmt string) error {

	if len(cmt) > 0 {
		g.Update(note,cmt)
	}

	nida := strings.Split(string(note.Nid),":")
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
