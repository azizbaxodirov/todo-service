package postgres

import (
	"database/sql"
	"fmt"
	"time"

	"github.com/jmoiron/sqlx"

	pb "github.com/azizbaxodirov/todo-service/genproto"
)

type todoRepo struct {
	db *sqlx.DB
}

// NewTodoRepo ...
func NewTodoRepo(db *sqlx.DB) *todoRepo {
	return &todoRepo{db: db}
}

func (r *todoRepo) Create(todo pb.Todo) (pb.Todo, error) {
	var id string
	err := r.db.QueryRow(`
        INSERT INTO todos(assignee, title, summary, deadline, status)
        VALUES ($1, $2, $3, $4, $5) returning id`, todo.Assignee, todo.Title, todo.Summary, todo.Deadline, todo.Status).Scan(&id)
	if err != nil {
		return pb.Todo{}, err
	}

	var NewTodo pb.Todo

	NewTodo, err = r.Get(id)

	if err != nil {
		return pb.Todo{}, err
	}

	return NewTodo, nil
}

func (r *todoRepo) Get(id string) (pb.Todo, error) {
	var todo pb.Todo
	err := r.db.QueryRow(`
        SELECT id, assignee, title, summary, deadline, status FROM todos WHERE id = $1`, id).Scan(&todo.Id, &todo.Assignee, &todo.Title, &todo.Summary, &todo.Deadline, &todo.Status)
	if err != nil {
		return pb.Todo{}, err
	}

	return todo, nil
}

func (r *todoRepo) ListOverdue(req time.Time, page, limit int64) ([]*pb.Todo, int64, error) {
	offset := (page - 1) * limit
	rows, err := r.db.Queryx(`
				SELECT id, assignee, title, summary, deadline, status
				FROM todos WHERE deadline >= $1 LIMIT $2 OFFSET $3`, req, limit, offset)
	if err != nil {
		return nil, 0, err
	}
	if err = rows.Err(); err != nil {
		return nil, 0, err
	}
	defer rows.Close()

	var (
		todos []*pb.Todo
		count int64
	)
	for rows.Next() {
		var todo pb.Todo
		err = rows.Scan(
			&todo.Id,
			&todo.Assignee,
			&todo.Title,
			&todo.Summary,
			&todo.Deadline,
			&todo.Status,
		)
		if err != nil {
			return nil, 0, err
		}
		todos = append(todos, &todo)
	}

	err = r.db.QueryRow(`SELECT count(*) FROM todos WHERE deadline >= $1`, req).Scan(&count)
	if err != nil {
		return nil, 0, err
	}

	return todos, count, nil
}

func (r *todoRepo) List(page, limit int64) ([]*pb.Todo, int64, error) {
	offset := (page - 1) * limit
	rows, err := r.db.Queryx(
		`SELECT id, assignee, title, summary, deadline, status FROM todos LIMIT $1 OFFSET $2`,
		limit, offset)
	if err != nil {
		return nil, 0, err
	}
	if err = rows.Err(); err != nil {
		return nil, 0, err
	}
	defer rows.Close() // nolint:err check

	var (
		todos []*pb.Todo
		count int64
	)
	for rows.Next() {
		var todo pb.Todo
		err = rows.Scan(&todo.Id, &todo.Assignee, &todo.Title, &todo.Summary, &todo.Deadline, &todo.Status)
		if err != nil {
			return nil, 0, err
		}
		todos = append(todos, &todo)
	}

	err = r.db.QueryRow(`SELECT count(*) FROM todos`).Scan(&count)
	if err != nil {
		return nil, 0, err
	}

	return todos, count, nil
}

func (r *todoRepo) Update(todo pb.Todo) (pb.Todo, error) {
	result, err := r.db.Exec(`UPDATE todos SET assignee=$1, title=$2, summary=$3, deadline=$4, status=$5 WHERE id=$6`,
		todo.Assignee, todo.Title, todo.Summary, todo.Deadline, todo.Status, todo.Id,
	)
	if err != nil {
		return pb.Todo{}, err
	}

	if i, _ := result.RowsAffected(); i == 0 {
		return pb.Todo{}, sql.ErrNoRows
	}

	var NewTodo pb.Todo

	NewTodo, err = r.Get(todo.Id)

	fmt.Println(result, NewTodo)

	if err != nil {
		return pb.Todo{}, err
	}

	return NewTodo, nil
}

func (r *todoRepo) Delete(id string) error {
	result, err := r.db.Exec(`DELETE FROM todos WHERE id = $1`, id)
	if err != nil {
		return err
	}

	if i, _ := result.RowsAffected(); i == 0 {
		return sql.ErrNoRows
	}

	return nil
}
