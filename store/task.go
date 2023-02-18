package store

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"time"

	"github.com/masumomo/go_todo_app/config"
	"github.com/masumomo/go_todo_app/entity"

	_ "github.com/go-sql-driver/mysql"
	"github.com/jmoiron/sqlx"
)

var (
	Tasks       = &TasksStore{Tasks: map[entity.TaskID]*entity.Task{}}
	ErrNotFOund = errors.New("not found")
)

func New(ctx context.Context, cfg *config.Config) (*sqlx.DB, func(), error) {
	db, err := sql.Open("mysql",
		fmt.Sprintf(
			"%s:%s@tcp(%s:%d)/%s?parseTime=true",
			cfg.DBUser, cfg.DBPassword,
			cfg.DBHost, cfg.DBPort,
			cfg.DBName,
		),
	)
	if err != nil {
		return nil, nil, err
	}
	ctx, cancel := context.WithTimeout(ctx, 2*time.Second)
	defer cancel()
	if err := db.PingContext(ctx); err != nil {
		return nil, func() { _ = db.Close() }, err
	}
	xdb := sqlx.NewDb(db, "mysql")
	return xdb, func() { _ = db.Close() }, nil
}

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

func (r *Repository) ListTasks(
	ctx context.Context, db Queryer,
) (entity.Tasks, error) {
	tasks := entity.Tasks{}
	sql := `SELECT
				id,
				title,
				status,
				created,
				modified
			FROM task;`
	if err := db.SelectContext(ctx, &tasks, sql); err != nil {
		return nil, err
	}
	return tasks, nil
}

func (r *Repository) AddTask(
	ctx context.Context, db Execer, t *entity.Task,
) error {
	t.Created = r.Clocker.Now()
	t.Modified = r.Clocker.Now()
	sql := `INSERT INTO task
	(title, status, created, modified)
	VALUES (?, ?, ?, ?)`
	result, err := db.ExecContext(ctx, sql, t.Title, t.Status, t.Created, t.Modified)
	if err != nil {
		return err
	}
	id, err := result.LastInsertId()

	if err != nil {
		return err
	}
	t.ID = entity.TaskID(id)
	return nil
}
