package main

import (
	"fmt"
	ui "github.com/gizak/termui"
	"github.com/urfave/cli"
	"github.com/xanzy/go-gitlab"
	"os"
	"sort"
	"time"
)

type Client struct {
	git *gitlab.Client
}

func max(a int, b int) int {
	if a < b {
		return b
	}
	return a
}

var git Client
var offline = false

func (client Client) get_open_issues(project string) ([]*gitlab.Issue, bool) {
	var issues []*gitlab.Issue
	var err error = nil
	used_local := false
	if !offline {
		options := &gitlab.ListProjectIssuesOptions{State: gitlab.String("opened")}
		issues, _, err = client.git.Issues.ListProjectIssues(project, options)
		if err == nil {
			Save(issues)
		}
	}
	
	if offline || err != nil {
		issues, err = Load()
		used_local = true
		check(err)
	} 

	sort.Slice(issues, func(i, j int) bool {
		return issues[i].IID < issues[j].IID
	})
	return issues, used_local
}

func Update() {
	ui.Body.Align()
	ui.Render(ui.Body)
}

func listIssues(issues []*gitlab.Issue, selected int) (r []string) {
	r = make([]string, len(issues))
	for i, issue := range issues {
		var labels string
		if len(issue.Labels) > 0 {
			labels = fmt.Sprintf("[%v](fg-yellow,fg-bold)", issue.Labels)
		}
		if i == selected {
			r[i] = fmt.Sprintf("[➤ #%v %v](fg-cyan,fg-bold) %v", issue.IID, issue.Title, labels)
		} else {
			r[i] = fmt.Sprintf("  [#%v](fg-bold,fg-green) [%v](fg-bold) %v", issue.IID, issue.Title, labels)
		}
	}
	return
}

func display_issues() {
	loading_icons := []rune{'⡿','⣟','⣯','⣷','⣾','⣽','⣻','⢿',}
	var issues []*gitlab.Issue
	loaded := false

	i := 0

	issues_panel := ui.NewList()
	issues_panel.BorderLabel = "Issues"
	issues_panel.Height = ui.TermHeight()
	issues_panel.Items = listIssues(issues, i)
	issues_panel.Overflow = "wrap"
	issues_details := ui.NewPar("")
	issues_details.Height = ui.TermHeight()
	issues_details.BorderLabel = "Details"

	default_layout := []*ui.Row {ui.NewRow(
		ui.NewCol(6, 0, issues_panel),
		ui.NewCol(6, 0, issues_details),
	)}
	details_layout := []*ui.Row {
		ui.NewCol(12, 0, issues_details),
	}

	ui.Body.Rows = default_layout

	change_selected := func() {
		if len(issues) == 0 {
			return
		}
		if len(issues) > issues_panel.Height {
			issues_panel.Items = listIssues(issues, i)[max(0, i-issues_panel.Height/2):]
		} else {
			issues_panel.Items = listIssues(issues, i)
		}
		issues_panel.Y = i
		issues_details.Text = fmt.Sprintf("[%v](fg-bold,fg-green)\n[Opened by: %v on %v\nAssigned To: %v](fg-bold,fg-blue)\n\n%v",
			issues[i].Title,
			issues[i].Author.Name,
			issues[i].CreatedAt.Format("2/01/2006"),
			issues[i].Assignee.Name,
			issues[i].Description)
		Update()
	}
	change_selected()

	ui.Handle("/sys/kbd/q", func(ui.Event) {
		ui.StopLoop()
	})

	ui.Handle("/sys/kbd/<up>", func(ui.Event) {
		i = (i - 1 + len(issues)) % len(issues)
		change_selected()
	})

	ui.Handle("/sys/kbd/<down>", func(ui.Event) {
		i = (i + 1) % len(issues)
		change_selected()
	})
	ui.Handle("/sys/kbd/<right>", func(ui.Event) {
		ui.Body.Rows = details_layout
		Update()
	})
	ui.Handle("/sys/kbd/<left>", func(ui.Event) {
		ui.Body.Rows = default_layout
		Update()
	})

	loader_ticker := time.NewTicker(time.Millisecond * 100)

	go func() {
		loader_index := 0
		for _ = range loader_ticker.C {
			loading_rune := fmt.Sprintf("%c", loading_icons[loader_index%len(loading_icons)])
			issues_details.Text = loading_rune
			issues_panel.Items = []string{loading_rune}
			loader_index++
			Update()
		}
	}()

	go func() {
		var used_local bool
		issues, used_local = git.get_open_issues("3528291")
		if used_local {
			issues_panel.BorderLabel += " (Offline)"
		}
		loaded = true
		loader_ticker.Stop()
		change_selected()
	}()

	Update()
	ui.Loop()
}

func main() {

	app := cli.NewApp()

	app.Name = "gitlab"
	app.Description = "A terminal dashboard for gitlab"

	app.Flags = []cli.Flag{
		cli.BoolFlag{
			Destination: &offline,
			Name:        "offline, o",
			Usage:       "Use saved data instead of live data",
		},
	}
	app.Action = func(c *cli.Context) error {
		git = Client{git: gitlab.NewClient(nil, "XqedyiDAFJ7zyV3h-tVp")}
		err := ui.Init()
		if err != nil {
			return err
		}
		defer ui.Close()

		display_issues()
		return nil
	}

	app.Run(os.Args)

}
