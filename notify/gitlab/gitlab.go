/* 2019-01-07 (cc) <paul4hough@gmail.com>
   gitlab issue interface

FIXME handle:
bad repo
bad issue (missing/deleted)
close closed

*/
package gitlab

import (
	"fmt"
	"strconv"
	"strings"

	gl "github.com/xanzy/go-gitlab"

	"github.com/pahoughton/agate/config"
	"github.com/pahoughton/agate/notify/nid"
)

type Gitlab struct {
	tsys	uint8
	grp		string
	debug	bool
	c		*gl.Client
}

func New(cfg config.NSysGitlab, tsys int,dbg bool) *Gitlab {
	g := &Gitlab{
		tsys:	uint8(tsys),
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

func (g *Gitlab)Create(prj, title, desc string, ) (nid.Nid, error) {

	i, resp, err := g.c.Issues.CreateIssue(prj,&gl.CreateIssueOptions{
		Title: gl.String(title),
		Description: gl.String("```\n"+desc+"\n```\n"),
	})

	if err != nil {
		if g.debug {
			fmt.Printf("dbg ERROR gitlab.CreateIssue: %T %v\n",err,err)
		}
		if resp != nil {
			glresp := string(err.(*gl.ErrorResponse).Body)
			if prj != g.grp &&
				strings.Contains(glresp, "Project Not Found") {
				return g.Create(
					g.grp, title,
					"prj default override: "+prj + " to "+g.grp+"\n"+desc)
			} else {
				// got unknown error response, non-recoverable
				panic(err.Error()+"\n\n"+prj+"\n"+title+"\n"+desc)
			}
		}
		return nil, err
	}
	if g.debug {
		fmt.Printf("gitlab.CreateIssue: ret issue: %v\n",i)
	}
	return nid.NewString(g.tsys,fmt.Sprintf("%s:%d",prj,i.IID)), nil
}

func (g *Gitlab)Update(id nid.Nid, cmt string) error {

	nida := strings.Split(id.Id(),":")
	prj := nida[0]
	issue, err := strconv.Atoi(nida[1])
	if err != nil {
		return fmt.Errorf("atoi: %s - %s",nida[1],err)
	}
	if g.debug {
		fmt.Printf("gitlab.AddComment: nid '%s' nida '%v' nida0 '%s' nida1 '%s' prj '%s' issue '%d'\n",
			id.Id(),
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

func (g *Gitlab)Close(id nid.Nid, cmt string) error {

	if len(cmt) > 0 {
		g.Update(id,cmt)
	}

	nida := strings.Split(id.Id(),":")
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
