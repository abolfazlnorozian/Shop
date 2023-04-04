package dbconnect

import (
	"shop/db"
	"shop/entity"
)

func MgoFindAllProducts() []entity.Products {
	var pro []entity.Products
	db.MgoConnect.Find(nil).All(&pro)

	return pro
}
