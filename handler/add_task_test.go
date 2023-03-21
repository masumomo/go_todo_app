package handler

import (
	"bytes"
	"context"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/go-playground/validator/v10"
	"github.com/masumomo/go_todo_app/entity"
	"github.com/masumomo/go_todo_app/testutil"
)

func TestAddTask(t *testing.T) {
	t.Parallel()
	type want struct {
		status  int
		rspFile string
	}
	tests := map[string]struct {
		reqFile string
		want    want
	}{
		"ok": {
			reqFile: "testdata/add_task/ok_req.json.golden",
			want: want{
				status:  http.StatusOK,
				rspFile: "testdata/add_task/ok_res.json.golden",
			},
		},
		"bad": {
			reqFile: "testdata/add_task/bad_req.json.golden",
			want: want{
				status:  http.StatusBadRequest,
				rspFile: "testdata/add_task/bad_res.json.golden",
			},
		},
	}
	for n, tt := range tests {
		tt := tt
		t.Run(n, func(t *testing.T) {
			t.Parallel()
		})

		w := httptest.NewRecorder()
		r := httptest.NewRequest(
			http.MethodPost,
			"/task",
			bytes.NewReader(testutil.LoadFile(t, tt.reqFile)),
		)

		moq := &AddTaskServiceMock{}
		moq.AddTaskFunc = func(ctx context.Context, title string) (*entity.Task, error) {
			return nil, nil
		}
		sut := AddTask{
			Service:   moq,
			Validator: validator.New(),
		}

		sut.ServeHTTP(w, r)

		resp := w.Result()
		testutil.AssertResponse(t,
			resp, tt.want.status, testutil.LoadFile(t, tt.want.rspFile),
		)
	}
}
