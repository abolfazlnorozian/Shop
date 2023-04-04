package dbconnect

import (
	"shop/db"
	"shop/entity"
)

func MgoFindAllCategories() []entity.Category {
	var pro []entity.Category
	db.MgoConnect.Find(nil).All(&pro)

	return pro
}
