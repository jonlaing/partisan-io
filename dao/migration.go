package dao

type Migratable interface {
	Table() string
	Migrations() []Migration
}

type Migration func(stamp int64, db *Database) error

func (db *Database) CreateTable(columns map[string]string) error {
	return nil
}

func (db *Database) InsertColumn(table, column, props string) error {
	return nil
}

func (db *Database) AlterColumn(table, column, props string) error {
	return nil
}

func (db *Database) RemoveColumn(table, column string) error {
	return nil
}
