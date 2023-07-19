package sqliteparserutils

import (
	"github.com/antlr/antlr4/runtime/Go/antlr/v4"
	"github.com/libsql/sqlite-antlr4-parser/sqliteparser"
)

var openToCloseTokensThatCanContainSemiColonInside = map[int]int{
	sqliteparser.SQLiteLexerBEGIN_: sqliteparser.SQLiteLexerEND_,
}

type stack[T any] struct {
	data []T
}

func newStack[T any]() *stack[T] {
	return &stack[T]{}
}

func (s *stack[T]) push(elem T) {
	s.data = append(s.data, elem)
}

func (s *stack[T]) pop() T {
	elem := s.data[len(s.data)-1]
	s.data = s.data[:len(s.data)-1]
	return elem
}

func (s *stack[T]) peek() T {
	return s.data[len(s.data)-1]
}

func (s *stack[T]) len() int {
	return len(s.data)
}

func SplitStatement(statement string) []string {
	statementStream := antlr.NewInputStream(statement)

	lexer := sqliteparser.NewSQLiteLexer(statementStream)
	tokenStream := antlr.NewCommonTokenStream(lexer, 0)

	stmtIntervals := make([]*antlr.Interval, 0)
	currentIntervalStart := -1
	openTokensStack := newStack[int]()

	for currentToken := tokenStream.LT(1); currentToken.GetTokenType() != antlr.TokenEOF; currentToken = tokenStream.LT(1) {
		tokenStream.Consume()

		if currentIntervalStart == -1 {
			if currentToken.GetTokenType() == sqliteparser.SQLiteLexerSCOL {
				continue
			}
			currentIntervalStart = currentToken.GetTokenIndex()
		}

		if _, isOpenToken := openToCloseTokensThatCanContainSemiColonInside[currentToken.GetTokenType()]; isOpenToken {
			openTokensStack.push(currentToken.GetTokenType())
		}

		if openTokensStack.len() > 0 {
			if closeToken, ok := openToCloseTokensThatCanContainSemiColonInside[openTokensStack.peek()]; ok && closeToken == currentToken.GetTokenType() {
				openTokensStack.pop()
			}
		} else if currentToken.GetTokenType() == sqliteparser.SQLiteLexerSCOL {
			stmtIntervals = append(stmtIntervals, antlr.NewInterval(currentIntervalStart, currentToken.GetTokenIndex()-1))
			currentIntervalStart = -1
		}
	}

	if currentIntervalStart != -1 {
		stmtIntervals = append(stmtIntervals, antlr.NewInterval(currentIntervalStart, tokenStream.LT(1).GetTokenIndex()))
	}

	stmts := make([]string, 0)
	for _, stmtInterval := range stmtIntervals {
		stmts = append(stmts, tokenStream.GetTextFromInterval(stmtInterval))
	}

	return stmts
}
