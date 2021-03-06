package postgres

import (
	"reflect"
	"testing"
	"time"

	pb "github.com/azizbaxodirov/todo-service/genproto"
)

func TestTodoRepo_Create(t *testing.T) {
	tests := []struct {
		name  string
		input pb.Todo
		want  pb.Todo
	}{
		{
			name: "Fist test",
			input: pb.Todo{
				Id:       "0d512776-60ed-4980-b8a3-6904a2234fd9",
				Assignee: "Test First Assignee",
				Title:    "Some Title",
				Summary:  "Summary",
				Deadline: "2021-10-15",
				Status:   "inactive",
			},
			want: pb.Todo{
				Assignee: "Test First Assignee",
				Title:    "Some Title",
				Summary:  "Summary",
				Deadline: "2021-10-15T00:00:00Z",
				Status:   "inactive",
			},
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			got, err := pgRepo.Create(tc.input)
			if err != nil {
				t.Fatalf("%s: expected: %v, got: %v", tc.name, tc.want, got)
			}
			got.Id = ""
			if !reflect.DeepEqual(tc.want, got) {
				t.Fatalf("%s: expected: %v, got: %v", tc.name, tc.want, got)
			}
		})
	}
}

func TestTodoRepo_Get(t *testing.T) {
	tests := []struct {
		name  string
		input string
		want  pb.Todo
	}{
		{
			name:  "successful",
			input: "0d512776-60ed-4980-b8a3-6904a2234fd7",
			want: pb.Todo{
				Assignee: "Test First Assignee",
				Title:    "Some Title",
				Summary:  "Summary",
				Deadline: "2021-10-15T00:00:00Z",
				Status:   "inactive",
			},
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			got, err := pgRepo.Get(tc.input)
			if err != nil {
				t.Fatalf("%s: expected: %v, got: %v", tc.name, tc.want, got)
			}
			got.Id = ""
			if !reflect.DeepEqual(tc.want, got) {
				t.Fatalf("%s: expected: %v, got: %v", tc.name, tc.want, got)
			}
		})
	}
}

type TestInput struct {
	time  time.Time
	page  int64
	limit int64
}

type TestWant struct {
	todos []pb.Todo
	count int64
}

func TestTodoRepo_List(t *testing.T) {
	tests := []struct {
		name  string
		input TestInput
		want  TestWant
	}{
		{
			name: "successful",
			input: TestInput{
				page:  1,
				limit: 2,
			},
			want: TestWant{
				todos: []pb.Todo{
					{
						Assignee: "Test First Assignee",
						Title:    "Some Title",
						Summary:  "Summary",
						Deadline: "2021-10-15T00:00:00Z",
						Status:   "inactive",
					},
					{
						Assignee: "Test First Assignee",
						Title:    "Some Title",
						Summary:  "Summary",
						Deadline: "2021-10-15T00:00:00Z",
						Status:   "inactive",
					},
				},
				count: 4,
			},
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			got, count, err := pgRepo.List(tc.input.page, tc.input.limit)
			if err != nil {
				t.Fatalf("%s: expected: %v, got: %v", tc.name, tc.want, got)
			}

			if count == tc.want.count {
				for i, j := range tc.want.todos {
					if j.Assignee != got[i].Assignee || j.Title != got[i].Title || j.Summary != got[i].Summary || j.Deadline != got[i].Deadline || j.Status != got[i].Status {
						t.Fatalf("%s: expected: %v, got: %v", tc.name, tc.want.todos, got)
					}
				}
			} else {
				t.Fatalf("%s: expected: %v, got: %v", tc.name, tc.want.todos, got)
			}
		})
	}
}

func TestTodoRepo_Update(t *testing.T) {
	tests := []struct {
		name  string
		input pb.Todo
		want  pb.Todo
	}{
		{
			name: "successful",
			input: pb.Todo{
				Id:       "0d512776-60ed-4980-b8a3-6904a2234fd9",
				Assignee: "Updated",
				Title:    "Updated",
				Summary:  "Updated",
				Deadline: "2021-12-14",
				Status:   "Updated",
			},
			want: pb.Todo{
				Assignee: "Updated",
				Title:    "Updated",
				Summary:  "Updated",
				Deadline: "2021-12-14T00:00:00Z",
				Status:   "Updated",
			},
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			got, err := pgRepo.Update(tc.input)
			if err != nil {
				t.Fatalf("%s: expected: %v, got: %v", tc.name, tc.want, got)
			}
			got.Id = ""
			if !reflect.DeepEqual(tc.want, got) {
				t.Fatalf("%s: expected: %v, got: %v", tc.name, tc.want, got)
			}
		})
	}
}

func TestTodoRepo_Delete(t *testing.T) {
	tests := []struct {
		name  string
		input string
		want  error
	}{
		{
			name:  "successful",
			input: "0d512776-60ed-4980-b8a3-6904a2234fd8",
			want:  nil,
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			err := pgRepo.Delete(tc.input)
			if err != nil {
				t.Fatalf("%s: expected: %v, got: %v", tc.name, tc.want, err)
			}
		})
	}
}

func TestTodoRepo_ListOverdue(t *testing.T) {
	layout := "2006-01-02"
	ParsedTime, _ := time.Parse(layout, "2012-10-10")

	tests := []struct {
		name  string
		input TestInput
		want  TestWant
	}{
		{
			name: "successful",
			input: TestInput{
				time:  ParsedTime,
				page:  1,
				limit: 5,
			},
			want: TestWant{
				todos: []pb.Todo{
					{
						Assignee: "Test First Assignee",
						Title:    "Some Title",
						Summary:  "Summary",
						Deadline: "2021-10-15T00:00:00Z",
						Status:   "inactive",
					},
					{
						Assignee: "Test First Assignee",
						Title:    "Some Title",
						Summary:  "Summary",
						Deadline: "2021-10-15T00:00:00Z",
						Status:   "inactive",
					},
					{
						Assignee: "Updated",
						Title:    "Updated",
						Summary:  "Updated",
						Deadline: "2021-12-14T00:00:00Z",
						Status:   "Updated",
					},
				},
				count: 3,
			},
		},
	}
	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			got, count, err := pgRepo.ListOverdue(tc.input.time, tc.input.page, tc.input.limit)
			if err != nil {
				t.Fatalf("%s: expected: %v, got: %v", tc.name, tc.want, got)
			}

			if tc.want.count == int64(count) {
				for i, j := range tc.want.todos {
					if j.Assignee != got[i].Assignee || j.Title != got[i].Title || j.Summary != got[i].Summary || j.Deadline != got[i].Deadline || j.Status != got[i].Status {
						t.Fatalf("%s: expected: %v, got: %v", tc.name, tc.want.todos, got)
					}
				}
			} else {
				t.Fatalf("%s: expected count: %v, got count: %v", tc.name, tc.want.count, count)
			}
		})
	}
}
