package migrations

var DataMap = make(map[string]MigrationInte)

type MigrationInte interface {
	Up()
	Down()
}