package dao

func dissectColumns(m Model) (ks []string, vs []interface{}) {
	for k, v := range m.Columns() {
		ks = append(ks, k)
		vs = append(vs, v)
	}
	return
}
