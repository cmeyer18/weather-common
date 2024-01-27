package common_tables

type IIdTable[T any] interface {
	Init() error

	Insert(item T) error

	Select(id string) (*T, error)

	Delete(id string) error
}
