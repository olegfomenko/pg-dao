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
	mockdb "github.com/olegfomenko/pg-dao/testing"
	"testing"
)

type Entry struct {
	Id   int64  `db:"id" structs:"-"`
	Name string `db:"name" structs:"name"`
}

func TestMain(t *testing.T) {
	cfg := config.New(kv.MustFromEnv())
	
	// Creating mock db with sql responses order
	// You can provide CheckSelectBuilder, CheckUpdateBuilder 
	// or CheckDeleteBuilder functions for checking sql query.
	// Or use mockdb.DefaultSelect for skipping sql checks
	mockDB := mockdb.NewDAO("entries", cfg.Log(),
		mockdb.MockData{
			CheckSelectBuilder: mockdb.DefaultSelect,
			CheckUpdateBuilder: mockdb.DefaultUpdate,
			CheckDeleteBuilder: mockdb.DefaultDelete,
			Entry:              Entry{},
			Error:              nil,
			Ok:                 false,
			T:                  t,
		},
		mockdb.MockData{
			CheckSelectBuilder: mockdb.DefaultSelect,
			CheckUpdateBuilder: mockdb.DefaultUpdate,
			CheckDeleteBuilder: mockdb.DefaultDelete,
			Entry: Entry{
				Id:   1,
				Name: "First Entry",
			},
			Error: nil,
			Ok:    true,
			T:     t,
		},
		mockdb.MockData{
			CheckSelectBuilder: mockdb.DefaultSelect,
			CheckUpdateBuilder: mockdb.DefaultUpdate,
			CheckDeleteBuilder: mockdb.DefaultDelete,
			Entry:              nil,
			Error:              nil,
			Ok:                 true,
			T:                  t,
		},
	)
	
	var entry Entry
	
	// returns false, nil due to first mock data
	ok, err := mockDB.FilterByID(1).Get(&entry)
	
	
	// returns true, nil and fills entry like Entry{1, "First Entry"} due second mock data
	ok, err = mockDB.New().FilterByID(1).Get(&entry)
}
```
