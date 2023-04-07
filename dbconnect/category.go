package dbconnect

// func MgoFindAllCategories() ([]*entity.Category, error) {

// 	var categ []*entity.Category
// 	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
// 	var DB = db.MD()
// 	var Collection = getcollection.GetCollection(DB, "categories")
// 	defer cancel()

// 	cursor, err := Collection.Find(ctx, bson.M{})
// 	if err != nil {
// 		return nil, err
// 	}
// 	for cursor.Next(ctx) {
// 		var category entity.Category
// 		err := cursor.Decode(&category)
// 		if err != nil {
// 			return nil, err
// 		}
// 		categ = append(categ, &category)
// 	}
// 	if err := cursor.Err(); err != nil {
// 		if err != nil {
// 			return nil, err
// 		}
// 		cursor.Close(ctx)
// 		if len(categ) == 0 {
// 			return nil, errors.New("document not found")
// 		}

// 	}

// 	return categ, nil
// }
