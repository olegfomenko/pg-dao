# pg-dao

Data Access Object for PostgreSQL for Distributed Lab projects.

## Usage

```go
package main

import pg "github.com/olegfomenko/pg-dao"

type Entry struct {
	Id   int64  `db:"id" structs:"-"`
	Name string `db:"name" structs:"name"`
}

func main() {
	// Loading config
	cfg := config.New(kv.MustFromEnv())

	// Creating DAO instance for table "entries"
	dao := pg.NewDAO(cfg.DB(), "entries")

	// Saving record
	// id - saved entry index
	id, err := dao.Create(Entry{
		Name: "First Entry",
	})

	// Cloning DAO for new session
	dao = dao.Clone()

	// Cleaning queries in DAO in current session
	dao = dao.New()

	// Updating entry
	err = dao.UpdateWhereID(id).UpdateColumn("name", "New First Entry").Update()

	// Getting entry by id
	var entry Entry
	// if ok is false - there is no entry with provided id
	ok, err := dao.New().FilterByID(id).Get(&entry)

	// Getting entry by field
	ok, err = dao.New().FilterByColumn("name", "New First Entry").Get(&entry)

	// Deleting entry
	err = dao.New().DeleteWhereID(id).Delete()
	
	
	// Creating transaction
	err = dao.Clone().TransactionSerializable(func(q pg.DAO) error {
		ok, err = q.FilterByID(id).Get(&entry)
		if err != nil {
			// rollback transaction
			return err
        	}
		
        
        	err = q.New().UpdateWhereID(id).UpdateColumn("name", "Updated First Entry").Update()
		if err != nil {
			// rollback transaction
			return err
		}

		// commit transaction
		return nil
	})
}
```

## Testing
Creating mocked dao

```go
package tests

import (
    "log"
	"testing"
	
	mockdb "github.com/olegfomenko/pg-dao/testing"
	"github.com/Masterminds/squirrel"
)

type Entry struct {
	Id   int64  `db:"id" structs:"-"`
	Name string `db:"name" structs:"name"`
}

func TestSimple(t *testing.T) {
	// Creating mock db with sql responses order
	mockDB := mockdb.New(t, "test").
		Add(Entry{}, false, nil).
		Add(Entry{
			Id: 1,
			Name: "First Entry",
		}, true, nil).
		Add(nil, true, nil).
		DAO()
	
	var entry Entry
	
	// returns false, nil due to first mock data
	ok, err := mockDB.FilterByID(1).Get(&entry)
	
	
	// returns true, nil and fills entry like Entry{1, "First Entry"} due second mock data
	ok, err = mockDB.New().FilterByID(1).Get(&entry)
	
	log.Printf("Entry: %+v", entry)
}

func Checker(t *testing.T, query squirrel.SelectBuilder) {
	str, _, err := query.ToSql()
	if err != nil {
		t.Error(err)
	}
	if str != "SELECT * FROM entries WHERE id = ?" {
		t.Errorf("Invalid query: %s", str)
	}
}

func TestWithCheckers(t *testing.T) {
	// Create mocked db with sql responses order, but
	// this time with function that will check each sql query

	mockDB := mockdb.New(t, "test").
		Add(Entry{Name: "First Entry"}, true, nil).
		SelectChecker(Checker).
		DAO()

	// returns true, nil and fills entry like Entry{1, "First Entry"} due first mock data
	// and checks sql query with Checker function
	err = mockDB.New().FilterByID(1).Get(&entry)
}
```
