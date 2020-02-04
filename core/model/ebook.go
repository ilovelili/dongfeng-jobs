package model

// Ebook Ebook entity
type Ebook struct {
	BaseModel
	Pupil     Pupil
	PupilID   uint
	Date      string
	Hash      string
	Converted bool
}

// Pupil pupil entity
type Pupil struct {
	BaseModel
	Name      string
	Class     Class
	ClassID   uint
	CreatedBy string
}

// Class class entity
type Class struct {
	BaseModel
	Year      string `gorm:"unique_index:idx_year_name" json:"year" csv:"学年"`
	Name      string `gorm:"unique_index:idx_year_name" json:"name" csv:"班级"`
	CreatedBy string `json:"created_by" csv:"-"`
}
