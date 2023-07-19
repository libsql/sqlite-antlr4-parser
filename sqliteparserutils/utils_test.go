package sqliteparserutils_test

import (
	"reflect"
	"testing"

	"github.com/libsql/sqlite-antlr4-parser/sqliteparserutils"
)

func TestSplitStatement(t *testing.T) {
	tests := []struct {
		name  string
		value string
		want  []string
	}{
		{
			name:  "EmptyStatement",
			value: "",
			want:  []string{},
		},
		{
			name:  "OnlySemicolon",
			value: ";;;;",
			want:  []string{},
		},
		{
			name:  "SingleStatementWithoutSemicolon",
			value: "select 1",
			want:  []string{"select 1"},
		},
		{
			name:  "SingleStatementWithSemicolon",
			value: "select 1;",
			want:  []string{"select 1"},
		},
		{
			name:  "MultipleCorrectStatements",
			value: "select 1; INSERT INTO counter(country, city, value) VALUES(?, ?, 1) ON CONFLICT DO UPDATE SET value = IFNULL(value, 0) + 1 WHERE country = ? AND city = ?; select 2",
			want:  []string{"select 1", "INSERT INTO counter(country, city, value) VALUES(?, ?, 1) ON CONFLICT DO UPDATE SET value = IFNULL(value, 0) + 1 WHERE country = ? AND city = ?", "select 2"},
		},
		{
			name:  "MultipleWrongStatements",
			value: "select from table; INSERT counter(country, city, value) VALUES(?, ?, 1) ON CONFLICT DO UPDATE SET value = IFNULL(value, 0) + 1 WHERE country = ? AND city = ?; create something",
			want:  []string{"select from table", "INSERT counter(country, city, value) VALUES(?, ?, 1) ON CONFLICT DO UPDATE SET value = IFNULL(value, 0) + 1 WHERE country = ? AND city = ?", "create something"},
		},
		{
			name:  "MultipleWrongTokens",
			value: "sdfasdfigosdfg sadfgsd ggsadgf; sdfasdfasd; 1230kfvcasd; 213 dsf s 0 fs229dt",
			want:  []string{"sdfasdfigosdfg sadfgsd ggsadgf", "sdfasdfasd", "1230kfvcasd", "213 dsf s 0 fs229dt"},
		},
		{
			name:  "MultipleSemicolonsBetweenStatements",
			value: "select 1;;;;;; ;;; ; ; ; ; select 2",
			want:  []string{"select 1", "select 2"},
		},
		{
			name:  "CreateTriggerStatement",
			value: "CREATE TRIGGER update_updated_at AFTER UPDATE ON users FOR EACH ROW BEGIN UPDATE users SET updated_at = 0 WHERE id = NEW.id; end;",
			want:  []string{"CREATE TRIGGER update_updated_at AFTER UPDATE ON users FOR EACH ROW BEGIN UPDATE users SET updated_at = 0 WHERE id = NEW.id; end"},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := sqliteparserutils.SplitStatement(tt.value)
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("got %#v, want %#v", got, tt.want)
			}
		})
	}
}
