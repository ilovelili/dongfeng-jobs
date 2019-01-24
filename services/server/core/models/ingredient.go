package models

// Ingredient ingredient entry
type Ingredient struct {
	ID   int64  `dapper:"id,primarykey,autoincrement,table=ingredients"`
	Name string `dapper:"name"`
}
