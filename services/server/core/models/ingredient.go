package models

// Ingredient ingredient entry
type Ingredient struct {
	ID       int64  `dapper:"id,primarykey,autoincrement,table=ingredients" csv:"-"`
	Name     string `dapper:"name" csv:"原料名称"`
	Material string `dapper:"material" csv:"物料名称"`
}
