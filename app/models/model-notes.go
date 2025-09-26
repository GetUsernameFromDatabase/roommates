package models

type Note struct {
	ModelBase
	Name    string `form:"name"`
	HouseID string `form:"house_id"`
	Content string `form:"content"`
}
