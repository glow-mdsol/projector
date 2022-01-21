package main

import (
	"strings"
)

// Rave URL
type RaveURL struct {
	PreferredURL string `db:"url"`
	URLID        int    `db:"id"`
	AlternateURL string `db:"alternate_url"`
	Projects     []*Project
}

// Get the URL by looking across the two candidates
func (r *RaveURL) URL() string {
	if r.AlternateURL != "" {
		return r.AlternateURL
	}
	return r.PreferredURL
}

// Get the Prefix URL (eg pharma.mdsol.com => pharma)
func (r *RaveURL) URLPrefix() string {
	return strings.Split(r.URL(), ".")[0]
}

//func createRaveURL(r Record) *RaveURL {
//	raveURL := new(RaveURL)
//	raveURL.URL = r.URL
//	raveURL.URLID = r.URLID
//	raveURL.URLPrefix = strings.Split(r.URL, ".")[0]
//	return raveURL
//}

func (r *RaveURL) getProject(projectName string) *Project {
	for _, prj := range r.Projects {
		if prj.ProjectName == projectName {
			return prj
		}
	}
	return nil
}

// get the ordered set of projects
func (r *RaveURL) getProjects() []*Project {
	name := func(p1, p2 *Project) bool {
		return p1.ProjectName < p2.ProjectName
	}
	ByP(name).Sort(r.Projects)
	return r.Projects
}
