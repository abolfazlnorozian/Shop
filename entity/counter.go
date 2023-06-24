package entity

type Counter struct {
	ID      string `bson:"_id"`
	OrderID int    `bson:"order_id"`
}
