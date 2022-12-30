package clause_test

import (
	"geeorm/clause"
	c "geeorm/clause"
	"reflect"
	"testing"
)

func testSelect(t *testing.T) {
	var clause clause.Clause
	clause.Set(c.LIMIT, 3)
	clause.Set(c.SELECT, "User", []string{"*"})
	clause.Set(c.WHERE, "Name = ?", "Tom")
	clause.Set(c.ORDERBY, "Age ASC")
	sql, vars := clause.Build(c.SELECT, c.WHERE, c.ORDERBY, c.LIMIT)
	t.Log(sql, vars)
	if sql != "SELECT * FROM User WHERE Name = ? ORDER BY Age ASC LIMIT ?" {
		t.Fatal("failed to build SQL")
	}
	if !reflect.DeepEqual(vars, []interface{}{"Tom", 3}) {
		t.Fatal("failed to build SQLVars")
	}
}

func TestClause_Build(t *testing.T) {
	t.Run("select", func(t *testing.T) {
		testSelect(t)
	})
}
