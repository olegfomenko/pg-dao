package pg_dao

import (
	"database/sql"
	"fmt"
	sq "github.com/Masterminds/squirrel"
	"github.com/fatih/structs"
	"gitlab.com/distributed_lab/kit/pgdb"
	"gitlab.com/distributed_lab/logan/v3/errors"
	"reflect"
	"time"
)

type dao struct {
	tableName string
	db        *pgdb.DB
	sql       sq.SelectBuilder
	upd       sq.UpdateBuilder
	dlt       sq.DeleteBuilder
}

func NewDAO(db *pgdb.DB, tableName string) DAO {
	return &dao{
		tableName: tableName,
		db:        db,
		sql:       sq.Select(tableName + ".*").From(tableName),
		upd:       sq.Update(tableName),
		dlt:       sq.Delete(tableName),
	}
}

func (d *dao) Clone() DAO {
	return &dao{
		tableName: d.tableName,
		db:        d.db.Clone(),
		sql:       sq.Select(d.tableName + ".*").From(d.tableName),
		upd:       sq.Update(d.tableName),
		dlt:       sq.Delete(d.tableName),
	}
}

func (d *dao) New() DAO {
	return &dao{
		tableName: d.tableName,
		db:        d.db,
		sql:       sq.Select(d.tableName + ".*").From(d.tableName),
		upd:       sq.Update(d.tableName),
		dlt:       sq.Delete(d.tableName),
	}
}

func (d *dao) Create(dto interface{}) (int64, error) {
	clauses := structs.Map(dto)

	var id int64
	stmt := sq.Insert(d.tableName).SetMap(clauses).Suffix("returning id")
	err := d.db.Get(&id, stmt)

	return id, err
}

func (d *dao) Get(dto interface{}) (bool, error) {
	if reflect.ValueOf(dto).Type().Kind() != reflect.Ptr {
		return false, errors.New("argument is not a pointer")
	}
	err := d.db.Get(dto, d.sql)
	if err == sql.ErrNoRows {
		return false, nil
	}

	return true, err
}

func (d *dao) Select(list interface{}) error {
	if reflect.ValueOf(list).Type().Kind() != reflect.Ptr {
		return errors.New("argument is not a slice pointer")
	}

	err := d.db.Select(list, d.sql)
	if err == sql.ErrNoRows {
		return nil
	}

	return err
}

func (d *dao) FilterByID(id int64) DAO {
	d.sql = d.sql.Where(sq.Eq{IdColumn: id})
	return d
}

func (d *dao) FilterOnlyAfter(time time.Time) DAO {
	d.sql = d.sql.Where(sq.Gt{CreatedAtColumn: time})
	return d
}

func (d *dao) FilterOnlyBefore(time time.Time) DAO {
	d.sql = d.sql.Where(sq.Lt{CreatedAtColumn: time})
	return d
}

func (d *dao) FilterGreater(col string, val interface{}) DAO {
	d.sql = d.sql.Where(sq.Gt{col: val})
	return d
}

func (d *dao) FilterLess(col string, val interface{}) DAO {
	d.sql = d.sql.Where(sq.Lt{col: val})
	return d
}

func (d *dao) FilterByColumn(col string, val interface{}) DAO {
	d.sql = d.sql.Where(sq.Eq{col: val})
	return d
}

func (d *dao) Limit(limit uint64) DAO {
	d.sql = d.sql.Limit(limit)
	return d
}

func (d *dao) OrderByDesc(col string) DAO {
	d.sql = d.sql.OrderBy(fmt.Sprintf("%s %s", col, OrderDescending))
	return d
}

func (d *dao) OrderByAsc(col string) DAO {
	d.sql = d.sql.OrderBy(fmt.Sprintf("%s %s", col, OrderAscending))
	return d
}

func (d *dao) UpdateWhereID(id int64) DAO {
	d.upd = d.upd.Where(sq.Eq{IdColumn: id})
	return d
}

func (d *dao) UpdateColumn(col string, val interface{}) DAO {
	d.upd = d.upd.Set(col, val)
	return d
}

func (d *dao) Update() error {
	res, err := d.db.ExecWithResult(d.upd)
	if err != nil {
		return errors.Wrap(err, "unable to update row")
	}
	rowsAffected, err := res.RowsAffected()
	if err != nil {
		return errors.Wrap(err, "unable to get affected rows")
	}
	if rowsAffected == 0 {
		return ErrNotFound
	}
	return nil
}

func (d *dao) DeleteWhereVal(col string, val interface{}) DAO {
	d.dlt = d.dlt.Where(sq.Eq{col: val})
	return d
}

func (d *dao) DeleteWhereID(id int64) DAO {
	d.dlt = d.dlt.Where(sq.Eq{IdColumn: id})
	return d
}

func (d *dao) Delete() error {
	err := d.db.Exec(d.dlt)
	if err != nil {
		return errors.Wrap(err, "unable to delete row")
	}
	return nil
}

func (d *dao) Page(params pgdb.OffsetPageParams) DAO {
	d.sql = params.ApplyTo(d.sql, "id")
	return d
}

func (d *dao) Transaction(fn func(q DAO) error) (err error) {
	return d.db.Transaction(func() error {
		return fn(d)
	})
}

func (d *dao) TransactionSerializable(fn func(q DAO) error) error {
	return d.db.TransactionWithOptions(&sql.TxOptions{Isolation: sql.LevelSerializable}, func() error {
		return fn(d)
	})
}

func (d *dao) TransactionWithLevel(level sql.IsolationLevel, fn func(q DAO) error) error {
	return d.db.TransactionWithOptions(&sql.TxOptions{Isolation: level}, func() error {
		return fn(d)
	})
}
