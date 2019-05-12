package models

// IngredientNutrition ingredient nutrition entry
type IngredientNutrition struct {
	ID               int64   `dapper:"id,primarykey,autoincrement,table=ingredient_nutritions" csv:"-"`
	Ingredient       string  `dapper:"ingredient" csv:"ingredient"`
	Protein100g      float64 `dapper:"protein_100g" csv:"protein_100g"`
	Fat100g          float64 `dapper:"fat_100g" csv:"fat_100g"`
	Carbohydrate100g float64 `dapper:"carbohydrate_100g" csv:"carbohydrate_100g"`
	Heat100g         float64 `dapper:"heat_100g" csv:"heat_100g"`
	Calcium100g      float64 `dapper:"calcium_100g" csv:"calcium_100g"`
	Iron100g         float64 `dapper:"iron_100g" csv:"iron_100g"`
	Zinc100g         float64 `dapper:"zinc_100g" csv:"zinc_100g"`
	VA100g           float64 `dapper:"va_100g" csv:"va_100g"`
	VB1100g          float64 `dapper:"vb1_100g" csv:"vb1_100g"`
	VB2100g          float64 `dapper:"vb2_100g" csv:"vb2_100g"`
	VC100g           float64 `dapper:"vc_100g" csv:"vc_100g"`
	CategoryID       int64   `dapper:"category_id" csv:"category_id"`
}
