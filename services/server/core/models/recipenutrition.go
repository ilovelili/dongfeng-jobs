package models

// RecipeNutrition recipe nutrition entry
type RecipeNutrition struct {
	ID           int64   `dapper:"id,primarykey,autoincrement,table=recipe_nutritions" csv:"-"`
	Recipe       string  `dapper:"recipe" csv:"recipe"`
	Carbohydrate float64 `dapper:"carbohydrate" csv:"carbohydrate"`
	Dietaryfiber float64 `dapper:"dietaryfiber" csv:"dietaryfiber"`
	Protein      float64 `dapper:"protein" csv:"protein"`
	Fat          float64 `dapper:"fat" csv:"fat"`
	Heat         float64 `dapper:"heat" csv:"heat"`
}
