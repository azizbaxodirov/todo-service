package repo

import (
	"time"

	pb "github.com/azizbaxodirov/todo-service/genproto"
)

// TodoStorageI ...
type TodoStorageI interface {
	Create(pb.Todo) (pb.Todo, error)
	Get(id string) (pb.Todo, error)
	List(page, limit int64) ([]*pb.Todo, int64, error)
	Update(pb.Todo) (pb.Todo, error)
	Delete(id string) error
	ListOverdue(time time.Time, page, limit int64) ([]*pb.Todo, int64, error)
}
