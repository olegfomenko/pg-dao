package testing

import (
	sq "github.com/Masterminds/squirrel"
	"testing"
)

type MockData struct {
	CheckSelectBuilder func(builder sq.SelectBuilder)
	CheckUpdateBuilder func(builder sq.UpdateBuilder)
	CheckDeleteBuilder func(builder sq.DeleteBuilder)
	Entry              interface{}
	Error              error
	Ok                 bool
	T                  *testing.T
}

func DefaultSelect(builder sq.SelectBuilder) {}
func DefaultUpdate(builder sq.SelectBuilder) {}
func DefaultDelete(builder sq.DeleteBuilder) {}
