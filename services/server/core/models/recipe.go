package models

// Recipe recipe entry
type Recipe struct {
	ID             int64   `dapper:"id,primarykey,autoincrement,table=recipes" csv:"-"`
	Name           string  `dapper:"name" csv:"菜品名称"`
	Ingredient     int64   `dapper:"ingredient_id" csv:"-"`
	IngredientName string  `dapper:"-" csv:"原料名称"`
	UnitAmount     float64 `dapper:"unit_amount" csv:"-"`
	CreatedBy      string  `dapper:"created_by" csv:"-"`
}
