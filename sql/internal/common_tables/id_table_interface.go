package common_tables

type IIdTable[T any] interface {
	Insert(item T) error

	Select(id string) (*T, error)

	Delete(id string) error
}
