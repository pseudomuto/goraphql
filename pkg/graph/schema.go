package graph

import (
	"github.com/graphql-go/graphql"
)

// Node describes a node within the graph
type Node interface {
	Queries() graphql.Fields
	Mutations() graphql.Fields
}

// NewSchema returns a new schema with queries and mutations from the supplied nodes available.
func NewSchema(nodes ...Node) (graphql.Schema, error) {
	return graphql.NewSchema(graphql.SchemaConfig{
		Query:    buildRootQuery(nodes),
		Mutation: buildRootMutations(nodes),
	})
}

func buildRootQuery(nodes []Node) *graphql.Object {
	fields := graphql.Fields{}

	for _, n := range nodes {
		for name, field := range n.Queries() {
			fields[name] = field
		}
	}

	return graphql.NewObject(graphql.ObjectConfig{
		Name:   "RootQuery",
		Fields: fields,
	})
}

func buildRootMutations(nodes []Node) *graphql.Object {
	fields := graphql.Fields{}

	for _, n := range nodes {
		for name, field := range n.Mutations() {
			fields[name] = field
		}
	}

	return graphql.NewObject(graphql.ObjectConfig{
		Name:   "RootQuery",
		Fields: fields,
	})
}
