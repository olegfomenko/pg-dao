package testing

import (
	"testing"

	sq "github.com/Masterminds/squirrel"
)

type QueryUpdateChecker func(*testing.T, sq.UpdateBuilder)
type QueryDeleteChecker func(*testing.T, sq.DeleteBuilder)
type QuerySelectChecker func(*testing.T, sq.SelectBuilder)

func DefaultQueryUpdateChecker(*testing.T, sq.UpdateBuilder) {}
func DefaultQueryDeleteChecker(*testing.T, sq.DeleteBuilder) {}
func DefaultQuerySelectChecker(*testing.T, sq.SelectBuilder) {}

type MockData struct {
	selectChecker QuerySelectChecker
	updateChecker QueryUpdateChecker
	deleteChecker QueryDeleteChecker
	Entry         interface{}
	Error         error
	Ok            bool
}
