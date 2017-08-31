package main

import (
	"github.com/ghodss/yaml"
	"github.com/xanzy/go-gitlab"
	"io/ioutil"
	"os/user"
	"path"
)

func check(err error) {
	if err != nil {
		panic(err)
	}
}

func get_issues_path() string {
	current_user, err := user.Current()
	check(err)
	return path.Join(current_user.HomeDir, ".gitlab-issues.yaml")
}

func Save(issues []*gitlab.Issue) {
	yaml_issues, err := yaml.Marshal(issues)
	check(err)
	err = ioutil.WriteFile(get_issues_path(), yaml_issues, 0644)
	check(err)
}

func Load() (issues []*gitlab.Issue, err error) {
	saved_issues, err := ioutil.ReadFile(get_issues_path())
	if err != nil {
		return
	}
	err = yaml.Unmarshal(saved_issues, &issues)
	if err != nil {
		return
	}
	return
}
