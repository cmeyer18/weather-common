module github.com/cmeyer18/weather-common/v5/sql/db_migrations/v5

go 1.22

require (
	github.com/cmeyer18/weather-common/v5 v5.0.5
	github.com/golang-migrate/migrate/v4 v4.17.0
	github.com/lib/pq v1.10.9
)

replace github.com/cmeyer18/weather-common/v5 => ../../../weather-common

require (
	github.com/hashicorp/errwrap v1.1.0 // indirect
	github.com/hashicorp/go-multierror v1.1.1 // indirect
	go.uber.org/atomic v1.7.0 // indirect
)
