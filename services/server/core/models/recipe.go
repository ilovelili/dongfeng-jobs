package models

// Recipe recipe entry
type Recipe struct {
	ID             int64  `dapper:"id,primarykey,autoincrement,table=recipes"`
	Name           string `dapper:"name"`
	Ingredient     int64  `dapper:"ingredient_id"`
	IngredientName string `dapper:"-"`
	CreatedBy      string `dapper:"created_by"`
}
