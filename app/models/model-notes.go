package models

type Note struct {
	ModelBase
	Name    string `form:"name"`
	Content string `form:"content"`
}
