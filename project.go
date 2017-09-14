package dashboard

import (
	"log"
	"sync"
	"time"
)

var (
	defaultProjectMap map[string]*Project
	defaultProjects   = []*Project{
		newProject("sensu-plugin", "sensu-plugins/sensu-plugin", "master", "sensu-plugin"),
		newProject("sensu-plugins-slack", "sensu-plugins/sensu-plugins-slack", "master", "sensu-plugins-slack"),
		newProject("sensu-extension", "sensu-extensions/sensu-extension", "master", "sensu-extension"),
		newProject("sensu-extensions-influxdb", "sensu-extensions/sensu-extensions-influxdb", "master", "sensu-extensions-influxdb"),
	}
)

func init() {
	go resetProjectsPeriodically()
}

func resetProjectsPeriodically() {
	for range time.Tick(time.Hour / 2) {
		log.Println("resetting projects' cache")
		resetProjects()
	}
}

func resetProjects() {
	for _, p := range defaultProjects {
		p.reset()
	}
}

type Project struct {
	Name    string `json:"name"`
	Nwo     string `json:"nwo"`
	Branch  string `json:"branch"`
	GemName string `json:"gem_name"`

	Gem     *RubyGem      `json:"gem"`
	Travis  *TravisReport `json:"travis"`
	GitHub  *GitHub       `json:"github"`
	fetched bool
}

func (p *Project) fetch() {
	rubyGemChan := rubygem(p.GemName)
	travisChan := travis(p.Nwo, p.Branch)
	githubChan := github(p.Nwo)

	if p.Gem == nil {
		p.Gem = <-rubyGemChan
	}

	if p.Travis == nil {
		p.Travis = <-travisChan
	}

	if p.GitHub == nil {
		p.GitHub = <-githubChan
	}

	p.fetched = true
}

func (p *Project) reset() {
	p.fetched = false
	p.Gem = nil
	p.Travis = nil
	p.GitHub = nil
}

func buildProjectMap() {
	defaultProjectMap = map[string]*Project{}
	for _, p := range defaultProjects {
		defaultProjectMap[p.Name] = p
	}
}

func newProject(name, nwo, branch, rubygem string) *Project {
	return &Project{
		Name:    name,
		Nwo:     nwo,
		Branch:  branch,
		GemName: rubygem,
	}
}

func getProject(name string) *Project {
	if defaultProjectMap == nil {
		buildProjectMap()
	}

	if p, ok := defaultProjectMap[name]; ok {
		if !p.fetched {
			p.fetch()
		}
		return p
	}

	return nil
}

func getAllProjects() []*Project {
	var wg sync.WaitGroup
	for _, p := range defaultProjects {
		wg.Add(1)
		go func(project *Project) {
			project.fetch()
			wg.Done()
		}(p)
	}
	wg.Wait()
	return defaultProjects
}

func getProjects() []*Project {
	return defaultProjects
}
