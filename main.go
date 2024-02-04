package pg_dao

import (
	"context"
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
// Notice that you should use Clone() to create new session and New() to use the same.
type DAO interface {
	Clone() DAO
	New() DAO
	Count() DAO

	Create(dto interface{}) (int64, error)
	CreateCtx(ctx context.Context, dto interface{}) (int64, error)

	FilterByID(id int64) DAO
	FilterGreater(col string, val interface{}) DAO
	FilterLess(col string, val interface{}) DAO
	FilterByColumn(col string, val interface{}) DAO

	Get(dto interface{}) (bool, error)
	GetCtx(ctx context.Context, dto interface{}) (bool, error)

	Select(list interface{}) error
	SelectCtx(ctx context.Context, list interface{}) error

	Limit(limit uint64) DAO
	OrderByDesc(col string) DAO
	OrderByAsc(col string) DAO

	UpdateWhereID(id int64) DAO
	UpdateColumn(col string, val interface{}) DAO

	Update() error
	UpdateCtx(ctx context.Context) error

	DeleteWhereVal(col string, val interface{}) DAO
	DeleteWhereID(id int64) DAO
	Delete() error
	DeleteCtx(ctx context.Context) error

	Page(params pgdb.OffsetPageParams, column string) DAO
	Cursor(params pgdb.CursorPageParams, column string) DAO

	Transaction(fn func(q DAO) error) error
	TransactionSerializable(fn func(q DAO) error) error
	TransactionWithLevel(level sql.IsolationLevel, fn func(q DAO) error) error

	ExecRaw(func(raw *pgdb.DB) error) error
	ExecRawCtx(ctx context.Context, fn func(ctx context.Context, raw *pgdb.DB) error) error
}
