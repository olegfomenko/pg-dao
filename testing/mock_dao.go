package testing

import (
	"fmt"
	sq "github.com/Masterminds/squirrel"
	pg "github.com/olegfomenko/pg-dao"
	"gitlab.com/distributed_lab/kit/pgdb"
	"gitlab.com/distributed_lab/logan/v3/errors"
	"reflect"
	"time"
)

type dao struct {
	tableName  string
	sql        sq.SelectBuilder
	upd        sq.UpdateBuilder
	dlt        sq.DeleteBuilder
	mocksOrder []MockData
}

func NewDAO(tableName string, mocks ...MockData) pg.DAO {
	return &dao{
		tableName:  tableName,
		sql:        sq.Select(tableName + ".*").From(tableName),
		upd:        sq.Update(tableName),
		dlt:        sq.Delete(tableName),
		mocksOrder: mocks,
	}
}

func (d *dao) Clone() pg.DAO {
	return &dao{
		tableName: d.tableName,
		sql:       sq.Select(d.tableName + ".*").From(d.tableName),
		upd:       sq.Update(d.tableName),
		dlt:       sq.Delete(d.tableName),
	}
}

func (d *dao) New() pg.DAO {
	return &dao{
		tableName: d.tableName,
		sql:       sq.Select(d.tableName + ".*").From(d.tableName),
		upd:       sq.Update(d.tableName),
		dlt:       sq.Delete(d.tableName),
	}
}

func (d *dao) Create(dto interface{}) (int64, error) {
	if len(d.mocksOrder) == 0 {
		panic("empty mocks")
	}

	mock := d.mocksOrder[0]
	d.mocksOrder = d.mocksOrder[1:]

	return mock.Entry.(int64), mock.Error
}

func (d *dao) Get(dto interface{}) (bool, error) {
	if reflect.ValueOf(dto).Type().Kind() != reflect.Ptr {
		return false, errors.New("argument is not a pointer")
	}

	if len(d.mocksOrder) == 0 {
		panic("empty mocks")
	}

	mock := d.mocksOrder[0]
	d.mocksOrder = d.mocksOrder[1:]

	mock.CheckSelectBuilder(d.sql)

	reflect.Indirect(reflect.ValueOf(dto)).Set(reflect.Indirect(reflect.ValueOf(mock.Entry)))
	return mock.Ok, mock.Error
}

func (d *dao) Select(list interface{}) error {
	if reflect.ValueOf(list).Type().Kind() != reflect.Ptr {
		return errors.New("argument is not a slice pointer")
	}

	if len(d.mocksOrder) == 0 {
		panic("empty mocks")
	}

	mock := d.mocksOrder[0]
	d.mocksOrder = d.mocksOrder[1:]

	mock.CheckSelectBuilder(d.sql)

	reflect.Indirect(reflect.ValueOf(list)).Set(reflect.Indirect(reflect.ValueOf(mock.Entry)))
	return mock.Error
}

func (d *dao) Delete() error {
	if len(d.mocksOrder) == 0 {
		panic("empty mocks")
	}

	mock := d.mocksOrder[0]
	d.mocksOrder = d.mocksOrder[1:]

	mock.CheckDeleteBuilder(d.dlt)

	return mock.Error
}

func (d *dao) Update() error {
	if len(d.mocksOrder) == 0 {
		panic("empty mocks")
	}

	mock := d.mocksOrder[0]
	d.mocksOrder = d.mocksOrder[1:]

	mock.CheckUpdateBuilder(d.upd)

	return mock.Error
}

func (d *dao) FilterByID(id int64) pg.DAO {
	d.sql = d.sql.Where(sq.Eq{pg.IdColumn: id})
	return d
}

func (d *dao) FilterOnlyAfter(time time.Time) pg.DAO {
	d.sql = d.sql.Where(sq.Gt{pg.CreatedAtColumn: time})
	return d
}

func (d *dao) FilterOnlyBefore(time time.Time) pg.DAO {
	d.sql = d.sql.Where(sq.Lt{pg.CreatedAtColumn: time})
	return d
}

func (d *dao) FilterGreater(col string, val interface{}) pg.DAO {
	d.sql = d.sql.Where(sq.Gt{col: val})
	return d
}

func (d *dao) FilterLess(col string, val interface{}) pg.DAO {
	d.sql = d.sql.Where(sq.Lt{col: val})
	return d
}

func (d *dao) FilterByColumn(col string, val interface{}) pg.DAO {
	d.sql = d.sql.Where(sq.Eq{col: val})
	return d
}

func (d *dao) Limit(limit uint64) pg.DAO {
	d.sql = d.sql.Limit(limit)
	return d
}

func (d *dao) OrderByDesc(col string) pg.DAO {
	d.sql = d.sql.OrderBy(fmt.Sprintf("%s %s", col, pg.OrderDescending))
	return d
}

func (d *dao) OrderByAsc(col string) pg.DAO {
	d.sql = d.sql.OrderBy(fmt.Sprintf("%s %s", col, pg.OrderAscending))
	return d
}

func (d *dao) UpdateWhereID(id int64) pg.DAO {
	d.upd = d.upd.Where(sq.Eq{pg.IdColumn: id})
	return d
}

func (d *dao) UpdateColumn(col string, val interface{}) pg.DAO {
	d.upd = d.upd.Set(col, val)
	return d
}

func (d *dao) DeleteWhereVal(col string, val interface{}) pg.DAO {
	d.dlt = d.dlt.Where(sq.Eq{col: val})
	return d
}

func (d *dao) DeleteWhereID(id int64) pg.DAO {
	d.dlt = d.dlt.Where(sq.Eq{pg.IdColumn: id})
	return d
}
func (d *dao) Page(params pgdb.OffsetPageParams) pg.DAO {
	panic("implement me")
}

func (d *dao) Transaction(fn func(q pg.DAO) error) (err error) {
	return fn(d)
}

func (d *dao) TransactionSerializable(fn func(q pg.DAO) error) error {
	return fn(d)
}
