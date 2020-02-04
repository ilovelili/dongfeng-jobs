package model

// BreakfastOrLunch breakfast or lunch emum
type BreakfastOrLunch uint

const (
	// AllMenuType all
	AllMenuType BreakfastOrLunch = 0
	// Breakfast breakfast
	Breakfast BreakfastOrLunch = 1
	// Lunch lunch
	Lunch BreakfastOrLunch = 2
	// Snack snack
	Snack BreakfastOrLunch = 3
)

// JuniorOrSenior juinor class menu or senior class menu
type JuniorOrSenior uint

const (
	// AllClass all
	AllClass JuniorOrSenior = 0
	// Junior junior class
	Junior JuniorOrSenior = 1
	// Senior senior class
	Senior JuniorOrSenior = 2
)

// Menu entity
type Menu struct {
	BaseModel
	Date               string
	Recipe             *Recipe
	RecipeID           uint
	BreakfastOrLunch   BreakfastOrLunch
	JuniorOrSenior     JuniorOrSenior
	JuniorBreakfastRaw string `gorm:"-"`
	JuniorLunchRaw     string `gorm:"-"`
	JuniorSnackRaw     string `gorm:"-"`
	SeniorBreakfastRaw string `gorm:"-"`
	SeniorLunchRaw     string `gorm:"-"`
	SeniorSnackRaw     string `gorm:"-"`
}

// Recipe entity
type Recipe struct {
	BaseModel
	Name              string           `csv:"菜品名称"`
	Ingredients       []*Ingredient    `gorm:"many2many:recipe_ingredients" csv:"-"`
	CSVIngredient     string           `gorm:"-" csv:"原料名称"`
	RecipeNutrition   *RecipeNutrition `csv:"-"`
	RecipeNutritionID *uint            `csv:"-"`
}

// RecipeNutrition entity
type RecipeNutrition struct {
	BaseModel
	Recipe       string  `csv:"recipe"`
	Carbohydrate float64 `csv:"carbohydrate"`
	Dietaryfiber float64 `csv:"dietaryfiber"`
	Protein      float64 `csv:"protein"`
	Fat          float64 `csv:"fat"`
	Heat         float64 `csv:"heat"`
}

// Ingredient ingredient entity
type Ingredient struct {
	BaseModel
	Recipes              []*Recipe `gorm:"many2many:recipe_ingredients" csv:"-"`
	Ingredient           string    `gorm:"unique_index" json:"ingredient"`
	Alias                string
	IngredientCategory   *IngredientCategory
	IngredientCategoryID uint
	Protein100g          float64 `gorm:"column:protein_100g"`
	ProteinDaily         float64
	Fat100g              float64 `gorm:"column:fat_100g"`
	FatDaily             float64
	Carbohydrate100g     float64 `gorm:"column:carbohydrate_100g"`
	CarbohydrateDaily    float64
	Heat100g             float64 `gorm:"column:heat_100g"`
	HeatDaily            float64
	Calcium100g          float64 `gorm:"column:calcium_100g"`
	CalciumDaily         float64
	Iron100g             float64 `gorm:"column:iron_100g"`
	IronDaily            float64
	Zinc100g             float64 `gorm:"column:zinc_100g"`
	ZincDaily            float64
	VA100g               float64 `gorm:"column:va_100g"`
	VADaily              float64
	VB1100g              float64 `gorm:"column:vb1_100g"`
	VB1Daily             float64
	VB2100g              float64 `gorm:"column:vb2_100g"`
	VB2Daily             float64
	VC100g               float64 `gorm:"column:vc_100g"`
	VCDaily              float64
	CreatedBy            string
}

// IngredientCategory ingredient category
type IngredientCategory struct {
	BaseModel
	Category  string
	CreatedBy string
}
