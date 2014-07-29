package features

import "github.com/goxmpp/goxmpp/stream"

type depGraph map[string][]string

var DependancyGraph stream.DependancyManageable = NewDependancyGraph()

func NewDependancyGraph() depGraph {
	return depGraph(make(map[string][]string))
}

func (dg depGraph) Add(name string, depends ...string) {
	for _, dep := range depends {
		if _, ok := dg[dep]; !ok {
			dg[dep] = []string{name}
		} else {
			dg[dep] = append(dg[dep], name)
		}
	}
}

func (dg depGraph) Get(name string) []string {
	active := []string{}
	for _, dep := range dg[name] {
		active = append(active, dep)
	}

	return active
}
