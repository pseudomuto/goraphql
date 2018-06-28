package graph

import (
	"github.com/graphql-go/graphql"

	"github.com/pseudomuto/goraphql/pkg/storage"
)

var todoType = graphql.NewObject(graphql.ObjectConfig{
	Name: "Todo",
	Fields: graphql.Fields{
		"id":      &graphql.Field{Type: graphql.String},
		"name":    &graphql.Field{Type: graphql.NewNonNull(graphql.String)},
		"content": &graphql.Field{Type: graphql.NewNonNull(graphql.String)},
	},
})

type todoGraph struct {
	repo storage.TodoRepo
}

// NewTodoNode returns a graph node for Todo objects
func NewTodoNode(repo storage.TodoRepo) Node {
	return &todoGraph{repo: repo}
}

func (g *todoGraph) Queries() graphql.Fields {
	return graphql.Fields{
		"todo": &graphql.Field{
			Type:        todoType,
			Description: "Find a particular todo by id",
			Args: graphql.FieldConfigArgument{
				"id": &graphql.ArgumentConfig{Type: graphql.NewNonNull(graphql.String)},
			},
			Resolve: func(params graphql.ResolveParams) (interface{}, error) {
				id, _ := params.Args["id"].(string)
				return g.repo.Find(id)
			},
		},
	}
}

func (g *todoGraph) Mutations() graphql.Fields {
	return graphql.Fields{
		"createTodo": &graphql.Field{
			Type:        todoType,
			Description: "Create a new todo item",
			Args: graphql.FieldConfigArgument{
				"name":    &graphql.ArgumentConfig{Type: graphql.NewNonNull(graphql.String)},
				"content": &graphql.ArgumentConfig{Type: graphql.NewNonNull(graphql.String)},
			},
			Resolve: func(params graphql.ResolveParams) (interface{}, error) {
				name, _ := params.Args["name"].(string)
				content, _ := params.Args["content"].(string)
				return g.repo.Create(name, content)
			},
		},
	}
}
