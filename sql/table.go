package sql

type PostgresTable[t any] interface {
	Init() error
}
