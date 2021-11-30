package pg_dao

import (
	"database/sql"
	"gitlab.com/distributed_lab/kit/pgdb"
	"gitlab.com/distributed_lab/logan/v3/errors"
)

const (
	IdColumn        = "id"
	CreatedAtColumn = "created_at"
	OrderAscending  = "asc"
	OrderDescending = "desc"
)

var ErrNotFound = errors.New("record not found")

// A DAO describes main methods for common data access object.
// Notice that you should use Clone() for every new request.
type DAO interface {
	Clone() DAO
	New() DAO

	Create(dto interface{}) (int64, error)

	FilterByID(id int64) DAO
	FilterGreater(col string, val interface{}) DAO
	FilterLess(col string, val interface{}) DAO
	FilterByColumn(col string, val interface{}) DAO

	Get(dto interface{}) (bool, error)
	Select(list interface{}) error

	Limit(limit uint64) DAO
	OrderByDesc(col string) DAO
	OrderByAsc(col string) DAO

	UpdateWhereID(id int64) DAO
	UpdateColumn(col string, val interface{}) DAO
	Update() error

	DeleteWhereVal(col string, val interface{}) DAO
	DeleteWhereID(id int64) DAO
	Delete() error

	Page(params pgdb.OffsetPageParams) DAO

	Transaction(fn func(q DAO) error) error
	TransactionSerializable(fn func(q DAO) error) error
	TransactionWithLevel(level sql.IsolationLevel, fn func(q DAO) error) error
}
