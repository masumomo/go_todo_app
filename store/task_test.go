package store

import (
	"context"
	"github.com/jmoiron/sqlx"
	"testing"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/google/go-cmp/cmp"

	"github.com/masumomo/go_todo_app/clock"
	"github.com/masumomo/go_todo_app/entity"
	"github.com/masumomo/go_todo_app/testutil"
)

// func TestNew(t *testing.T) {
// 	type args struct {
// 		ctx context.Context
// 		cfg *config.Config
// 	}
// 	tests := []struct {
// 		name    string
// 		args    args
// 		want    *sqlx.DB
// 		want1   func()
// 		wantErr bool
// 	}{
// 		// TODO: Add test cases.
// 	}
// 	for _, tt := range tests {
// 		t.Run(tt.name, func(t *testing.T) {
// 			got, got1, err := New(tt.args.ctx, tt.args.cfg)
// 			if (err != nil) != tt.wantErr {
// 				t.Errorf("New() error = %v, wantErr %v", err, tt.wantErr)
// 				return
// 			}
// 			if !reflect.DeepEqual(got, tt.want) {
// 				t.Errorf("New() got = %v, want %v", got, tt.want)
// 			}
// 			if !reflect.DeepEqual(got1, tt.want1) {
// 				t.Errorf("New() got1 = %v, want %v", got1, tt.want1)
// 			}
// 		})
// 	}
// }

// func TestTasksStore_Add(t *testing.T) {
// 	type fields struct {
// 		lastID entity.TaskID
// 		Tasks  map[entity.TaskID]*entity.Task
// 	}
// 	type args struct {
// 		t *entity.Task
// 	}
// 	tests := []struct {
// 		name    string
// 		fields  fields
// 		args    args
// 		want    entity.TaskID
// 		wantErr bool
// 	}{
// 		// TODO: Add test cases.
// 	}
// 	for _, tt := range tests {
// 		t.Run(tt.name, func(t *testing.T) {
// 			ts := &TasksStore{
// 				lastID: tt.fields.lastID,
// 				Tasks:  tt.fields.Tasks,
// 			}
// 			got, err := ts.Add(tt.args.t)
// 			if (err != nil) != tt.wantErr {
// 				t.Errorf("TasksStore.Add() error = %v, wantErr %v", err, tt.wantErr)
// 				return
// 			}
// 			if !reflect.DeepEqual(got, tt.want) {
// 				t.Errorf("TasksStore.Add() = %v, want %v", got, tt.want)
// 			}
// 		})
// 	}
// }

// func TestTasksStore_All(t *testing.T) {
// 	type fields struct {
// 		lastID entity.TaskID
// 		Tasks  map[entity.TaskID]*entity.Task
// 	}
// 	tests := []struct {
// 		name   string
// 		fields fields
// 		want   entity.Tasks
// 	}{
// 		// TODO: Add test cases.
// 	}
// 	for _, tt := range tests {
// 		t.Run(tt.name, func(t *testing.T) {
// 			ts := &TasksStore{
// 				lastID: tt.fields.lastID,
// 				Tasks:  tt.fields.Tasks,
// 			}
// 			if got := ts.All(); !reflect.DeepEqual(got, tt.want) {
// 				t.Errorf("TasksStore.All() = %v, want %v", got, tt.want)
// 			}
// 		})
// 	}
// }

func TestRepository_ListTask(t *testing.T) {
	ctx := context.Background()
	tx, err := testutil.OpenDBForTest(t).BeginTxx(ctx, nil)

	t.Cleanup(func() { _ = tx.Rollback() })
	if err != nil {
		t.Fatal()
	}
	wants := prepareTasks(ctx, t, tx)
	r := &Repository{}
	got, err := r.ListTasks(ctx, tx)
	if err != nil {
		t.Errorf("Unexpected err: %v", err)
	}
	if d := cmp.Diff(got, wants); len(d) != 0 {
		t.Errorf("Diff -got +want: %v", d)
	}
}

func prepareTasks(ctx context.Context, t *testing.T, con Execer) entity.Tasks {
	t.Helper()
	if _, err := con.ExecContext(ctx, "DELETE FROM task;"); err != nil {
		t.Logf("Failed to initialize task: %v", err)
	}
	c := clock.FixedClocker{}
	wants := entity.Tasks{
		{
			Title: "want task 1", Status: "todo",
			Created: c.Now(), Modified: c.Now(),
		},
		{
			Title: "want task 2", Status: "todo",
			Created: c.Now(), Modified: c.Now(),
		},
		{
			Title: "want task 3", Status: "done",
			Created: c.Now(), Modified: c.Now(),
		},
	}
	result, err := con.ExecContext(ctx,
		`INSERT INTO task (title, status, created, modified)
		 VALUES
		   (?, ?, ?, ?),
		   (?, ?, ?, ?),
		   (?, ?, ?, ?); 
		`,
		wants[0].Title, wants[0].Status, wants[0].Created, wants[0].Modified,
		wants[1].Title, wants[1].Status, wants[1].Created, wants[1].Modified,
		wants[2].Title, wants[2].Status, wants[2].Created, wants[2].Modified,
	)
	if err != nil {
		t.Fatal(err)
	}
	id, err := result.LastInsertId()
	if err != nil {
		t.Fatal(err)
	}
	wants[0].ID = entity.TaskID(id)
	wants[1].ID = entity.TaskID(id + 1)
	wants[2].ID = entity.TaskID(id + 2)

	return wants
}

func TestRepository_AddTask(t *testing.T) {
	t.Parallel()
	ctx := context.Background()

	c := clock.FixedClocker{}
	var wantID int64 = 20

	okTask := &entity.Task{
		Title:    "ok task",
		Status:   "todo",
		Created:  c.Now(),
		Modified: c.Now(),
	}

	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatal(err)
	}
	t.Cleanup(func() { db.Close() })

	mock.ExpectExec(
		`INSERT INTO task \(title, status, created, modified\) VALUES \(\?, \?, \?, \?\)`,
	).WithArgs(okTask.Title, okTask.Status, c.Now(), c.Now()).WillReturnResult(sqlmock.NewResult(wantID, 1))

	xdb := sqlx.NewDb(db, "mysql")
	r := &Repository{Clocker: c}

	if err := r.AddTask(ctx, xdb, okTask); err != nil {
		t.Errorf("got error %v", err)
	}
}
