package common_tables

type IIdTable[T any] interface {
	// Deprecated: use the weather-db-migrator docker images for setup
	Init() error

	Insert(item T) error

	Select(id string) (*T, error)

	Delete(id string) error

	Update(id string, item T) error
}
