package testing

import (
	"database/sql"
	"fmt"
	"reflect"
	"testing"

	sq "github.com/Masterminds/squirrel"
	pg "github.com/olegfomenko/pg-dao"
	"gitlab.com/distributed_lab/kit/pgdb"
	"gitlab.com/distributed_lab/logan/v3/errors"
)

type dao struct {
	tableName  string
	sql        sq.SelectBuilder
	upd        sq.UpdateBuilder
	dlt        sq.DeleteBuilder
	mocksOrder *[]MockData
	t          *testing.T
}

func NewDAO(tableName string, t *testing.T) pg.DAO {
	return &dao{
		tableName:  tableName,
		sql:        sq.Select(tableName + ".*").From(tableName),
		upd:        sq.Update(tableName),
		dlt:        sq.Delete(tableName),
		mocksOrder: nil,
		t:          t,
	}
}

func (d *dao) Clone() pg.DAO {
	d.t.Log("Cloning dao")
	return &dao{
		tableName:  d.tableName,
		sql:        sq.Select(d.tableName + ".*").From(d.tableName),
		upd:        sq.Update(d.tableName),
		dlt:        sq.Delete(d.tableName),
		mocksOrder: d.mocksOrder,
		t:          d.t,
	}
}

func (d *dao) New() pg.DAO {
	d.t.Log("Initializing new raws")
	return &dao{
		tableName:  d.tableName,
		sql:        sq.Select(d.tableName + ".*").From(d.tableName),
		upd:        sq.Update(d.tableName),
		dlt:        sq.Delete(d.tableName),
		mocksOrder: d.mocksOrder,
		t:          d.t,
	}
}

func (d *dao) Count() pg.DAO {
	d.t.Log("Counting records")
	return &dao{
		tableName:  d.tableName,
		sql:        sq.Select("count(*)").From(d.tableName),
		upd:        sq.Update(d.tableName),
		dlt:        sq.Delete(d.tableName),
		mocksOrder: d.mocksOrder,
		t:          d.t,
	}
}

func (d *dao) Create(dto interface{}) (int64, error) {
	d.t.Log("Saving new record ", dto)

	if len(*d.mocksOrder) == 0 {
		d.t.Fatal("empty mocks")
	}

	mock := (*d.mocksOrder)[0]
	next := (*d.mocksOrder)[1:]
	*d.mocksOrder = next

	return mock.Entry.(int64), mock.Error
}

// NOTE: Not defined in DAO interface
//
// func (d *dao) FilterOnlyAfter(time time.Time) pg.DAO {
// 	d.log.Debug("Filtering after time: ", time.String())
// 	d.sql = d.sql.Where(sq.Gt{pg.CreatedAtColumn: time})
// 	return d
// }

// NOTE: Not defined in DAO interface
//
// func (d *dao) FilterOnlyBefore(time time.Time) pg.DAO {
// 	d.log.Debug("Filtering before time: ", time.String())
// 	d.sql = d.sql.Where(sq.Lt{pg.CreatedAtColumn: time})
// 	return d
// }

func (d *dao) FilterByID(id int64) pg.DAO {
	d.t.Logf("Filtering by id: %d", id)
	d.sql = d.sql.Where(sq.Eq{pg.IdColumn: id})
	return d
}

func (d *dao) FilterGreater(col string, val interface{}) pg.DAO {
	d.t.Logf("Filtering greater column: %s value: %v", col, val)
	d.sql = d.sql.Where(sq.Gt{col: val})
	return d
}

func (d *dao) FilterLess(col string, val interface{}) pg.DAO {
	d.t.Logf("Filtering less column: %s value: %v", col, val)
	d.sql = d.sql.Where(sq.Lt{col: val})
	return d
}

func (d *dao) FilterByColumn(col string, val interface{}) pg.DAO {
	d.t.Logf("Filtering by column: %s value: %v", col, val)
	d.sql = d.sql.Where(sq.Eq{col: val})
	return d
}

func (d *dao) Get(dto interface{}) (bool, error) {
	d.t.Log("Getting record")
	if reflect.ValueOf(dto).Type().Kind() != reflect.Ptr {
		return false, errors.New("argument is not a pointer")
	}

	if len(*d.mocksOrder) == 0 {
		d.t.Fatal("empty mocks")
	}

	mock := (*d.mocksOrder)[0]
	next := (*d.mocksOrder)[1:]
	*d.mocksOrder = next

	mock.CheckSelectBuilder(d.sql)

	reflect.Indirect(reflect.ValueOf(dto)).Set(reflect.Indirect(reflect.ValueOf(mock.Entry)))
	return mock.Ok, mock.Error
}

func (d *dao) Select(list interface{}) error {
	d.t.Log("Selecting record")
	if reflect.ValueOf(list).Type().Kind() != reflect.Ptr {
		return errors.New("argument is not a slice pointer")
	}

	if len(*d.mocksOrder) == 0 {
		d.t.Fatal("empty mocks")
	}

	mock := (*d.mocksOrder)[0]
	next := (*d.mocksOrder)[1:]
	*d.mocksOrder = next

	mock.CheckSelectBuilder(d.sql)

	reflect.Indirect(reflect.ValueOf(list)).Set(reflect.Indirect(reflect.ValueOf(mock.Entry)))
	return mock.Error
}

func (d *dao) Limit(limit uint64) pg.DAO {
	d.t.Logf("Limiting rows: %d", limit)
	d.sql = d.sql.Limit(limit)
	return d
}

func (d *dao) OrderByDesc(col string) pg.DAO {
	d.t.Logf("Ordering descending column: %s", col)
	d.sql = d.sql.OrderBy(fmt.Sprintf("%s %s", col, pg.OrderDescending))
	return d
}

func (d *dao) OrderByAsc(col string) pg.DAO {
	d.t.Logf("Ordering ascending column: %s", col)
	d.sql = d.sql.OrderBy(fmt.Sprintf("%s %s", col, pg.OrderAscending))
	return d
}

func (d *dao) UpdateWhereID(id int64) pg.DAO {
	d.t.Logf("Updating by id: %d", id)
	d.upd = d.upd.Where(sq.Eq{pg.IdColumn: id})
	return d
}

func (d *dao) UpdateColumn(col string, val interface{}) pg.DAO {
	d.t.Logf("Updating column: %s value: %v", col, val)
	d.upd = d.upd.Set(col, val)
	return d
}

func (d *dao) Update() error {
	d.t.Log("Updating record")
	if len(*d.mocksOrder) == 0 {
		d.t.Fatal("empty mocks")
	}

	mock := (*d.mocksOrder)[0]
	next := (*d.mocksOrder)[1:]
	*d.mocksOrder = next

	mock.CheckUpdateBuilder(d.upd)

	return mock.Error
}

func (d *dao) DeleteWhereVal(col string, val interface{}) pg.DAO {
	d.t.Log("Deleting where column:", col, " value: ", val)
	d.dlt = d.dlt.Where(sq.Eq{col: val})
	return d
}

func (d *dao) DeleteWhereID(id int64) pg.DAO {
	d.t.Logf("Deleting by id: %d", id)
	d.dlt = d.dlt.Where(sq.Eq{pg.IdColumn: id})
	return d
}

func (d *dao) Delete() error {
	d.t.Log("Deleting record")
	if len(*d.mocksOrder) == 0 {
		d.t.Fatal("empty mocks")
	}

	mock := (*d.mocksOrder)[0]
	next := (*d.mocksOrder)[1:]
	*d.mocksOrder = next

	mock.CheckDeleteBuilder(d.dlt)

	return mock.Error
}

func (d *dao) Page(params pgdb.OffsetPageParams) pg.DAO {
	d.t.Log("Applying page parms")
	d.sql = params.ApplyTo(d.sql, "id")
	return d
}

func (d *dao) Transaction(fn func(q pg.DAO) error) (err error) {
	d.t.Log("Starting db transaction")
	defer d.t.Log("Finishing db transaction")
	return fn(d)
}

func (d *dao) TransactionSerializable(fn func(q pg.DAO) error) error {
	d.t.Log("Starting db serializable transaction")
	defer d.t.Log("Finishing db serializable transaction")
	return fn(d)
}

func (d *dao) TransactionWithLevel(level sql.IsolationLevel, fn func(q pg.DAO) error) error {
	d.t.Log("Starting db transaction with level: ", level)
	defer d.t.Logf("Finishing db transaction with level: %d", level)
	return fn(d)
}
