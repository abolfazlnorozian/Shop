package dbconnect

// func MgoFindAllProducts() ([]*entity.Products, error) {

// 	var products []*entity.Products
// 	// db.MongoConnect.Find(context.Background(), pro)
// 	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
// 	var DB = db.MD()
// 	var Collection =db.GetCollection(DB, "products")
// 	defer cancel()

// 	// return pro
// 	cursor, err := Collection.Find(ctx, bson.D{{}})
// 	if err != nil {
// 		return nil, err
// 	}
// 	for cursor.Next(ctx) {
// 		var pro entity.Products
// 		err := cursor.Decode(&pro)
// 		if err != nil {
// 			return nil, err
// 		}
// 		products = append(products, &pro)
// 	}
// 	if err := cursor.Err(); err != nil {
// 		if err != nil {
// 			return nil, err
// 		}
// 		cursor.Close(ctx)
// 		if len(products) == 0 {
// 			return nil, errors.New("document not found")
// 		}

// 	}

// 	return products, nil
// }
