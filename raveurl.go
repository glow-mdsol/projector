package main

import "sort"

// Rave URL
type RaveURL struct {
	URL       string
	URLID     int
	URLPrefix string
	Projects  []*Project
}

func (raveUrl *RaveURL) getProject(projectName string) *Project {
	for _, prj := range raveUrl.Projects {
		if prj.ProjectName == projectName {
			return prj
		}
	}
	return nil
}

// Sorter
type ByP func(v1, v2 *Project) bool

func (by ByP) Sort(projects []*Project) {
	vs := &projectSorter{
		projects: projects,
		by:       by,
	}
	sort.Sort(vs)
}

type projectSorter struct {
	projects []*Project
	by       func(v1, v2 *Project) bool
}

// Len is part of sort.Interface.
func (s *projectSorter) Len() int {
	return len(s.projects)
}

// Swap is part of sort.Interface.
func (s *projectSorter) Swap(i, j int) {
	s.projects[i], s.projects[j] = s.projects[j], s.projects[i]
}

// Less is part of sort.Interface. It is implemented by calling the "by" closure in the sorter.
func (s *projectSorter) Less(i, j int) bool {
	return s.by(s.projects[i], s.projects[j])
}

// get the ordered set of projects
func (raveUrl *RaveURL) getProjects() []*Project {
	name := func(p1, p2 *Project) bool {
		return p1.ProjectName < p2.ProjectName
	}
	ByP(name).Sort(raveUrl.Projects)
	return raveUrl.Projects
}

func (raveURL *RaveURL) fixupURL() {
	for _, project := range raveURL.Projects {
		for _, projectVersion := range project.Versions {
			projectVersion.fixUpNullValues()
		}
	}
}
