package models

type House struct {
	ModelBase
	Name         string   `form:"name"`
	RoommateKeys []string `form:"roommates[]"`
	// used only by htmx to get user suggestions when adding roomates to house
	SearchedUser string `form:"searched_user"`
	// make sure to match indices with RoommateKeys
	//  used just to display things for user
	RoommateLabels []string `form:"roommates_labels[]"`
}
