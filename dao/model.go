package dao

type Model interface {
	Table() string
	Columns() map[string]interface{}
}

func (db *Database) GetByID(m Model, id uint64) error {
	colNames, scanVals := dissectColumns(m)

	row := db.Conn.QueryRow(`SELECT ? FROM ? WHERE id = ?`, colNames, m.Table(), id)
	row.Scan(scanVals...)

	return nil
}

func (db *Database) GetCollection(m Model, sel string, inter ...[]interface{}) error {
	// colNames, scanVals := dissectColumns(m)

	return nil
}
