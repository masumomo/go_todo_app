package store

import (
	"errors"
	"github.com/masumomo/go_todo_app/entity"
)

var (
	Tasks       = &TasksStore{Tasks: map[entity.TaskID]*entity.Task{}}
	ErrNotFOund = errors.New("not found")
)

type TasksStore struct {
	lastID entity.TaskID
	Tasks  map[entity.TaskID]*entity.Task
}

func (ts *TasksStore) Add(t *entity.Task) (entity.TaskID, error) {
	ts.lastID++
	t.ID = ts.lastID
	ts.Tasks[t.ID] = t
	return t.ID, nil
}

func (ts *TasksStore) All() entity.Tasks {
	tasks := make([]*entity.Task, len(ts.Tasks))
	for i, t := range ts.Tasks {
		tasks[i-1] = t
	}
	return tasks
}
