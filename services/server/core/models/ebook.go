package models

// Ebook Ebook entity
type Ebook struct {
	ID        int64  `dapper:"id,primarykey,autoincrement,table=ebooks"`
	Year      string `dapper:"year"`
	Class     string `dapper:"class"`
	Name      string `dapper:"name"`
	Date      string `dapper:"date"`
	Hash      string `dapper:"hash"`
	Converted bool   `dapper:"converted"`
}
