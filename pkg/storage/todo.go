package storage

import (
	"errors"
	"strconv"
	"time"
)

// Todo represents a todo entry
type Todo struct {
	ID      string `json:"id"`
	Name    string `json:"name"`
	Content string `json:"content"`
}

// TodoRepo describes a repository of todo items
type TodoRepo interface {
	Find(id string) (*Todo, error)
	Create(name, content string) (*Todo, error)
}

type ephemeralTodoRepo struct {
	todos []*Todo
}

// NewEphemeralTodoRepo returns a TodoRepo that stores records in memory
func NewEphemeralTodoRepo() TodoRepo {
	return &ephemeralTodoRepo{todos: make([]*Todo, 0, 10)}
}

func (r *ephemeralTodoRepo) Find(id string) (*Todo, error) {
	for _, t := range r.todos {
		if t.ID == id {
			return t, nil
		}
	}

	return nil, errors.New("No todo found with the specified id")
}

func (r *ephemeralTodoRepo) Create(name, content string) (*Todo, error) {
	t := &Todo{
		ID:      strconv.FormatInt(time.Now().Unix(), 10),
		Name:    name,
		Content: content,
	}

	r.todos = append(r.todos, t)
	return t, nil
}
