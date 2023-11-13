package services

import (
	"context"
	"log"
	"net/http"
	"shop/database"
	"shop/entities"
	"time"

	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

var pagesCollection *mongo.Collection = database.GetCollection(database.DB, "pages")
var rowCollection *mongo.Collection = database.GetCollection(database.DB, "rows")
var colCollection *mongo.Collection = database.GetCollection(database.DB, "columns")

// func GetPages(c *gin.Context) {

// 	mode := c.DefaultQuery("mode", "")
// 	filter := bson.M{}

// 	if mode != "" {
// 		filter["mode"] = mode

// 	}
// 	cur, err := pagesCollection.Find(c, filter)
// 	if err != nil {
// 		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
// 		return
// 	}
// 	fmt.Println("cur:", cur)
// 	defer cur.Close(c)
// 	var pages []entities.Pages
// 	// responseBody["success"]=true
// 	// responseBody["message"]="page"

// 	for cur.Next(c) {
// 		var pgs entities.Pages
// 		err := cur.Decode(&pgs)
// 		if err != nil {
// 			c.JSON(http.StatusInternalServerError, gin.H{"message": err.Error()})
// 			return
// 		}

// 		pageMap := map[string]interface{}{
// 			"_id": pgs.Id,
// 		}

// 	}

// }

//***********************************************************************************************

func GetPages(c *gin.Context) {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	mode := c.Query("mode") // Get the mode parameter from the query

	filter := bson.M{}

	if mode != "" {
		filter["mode"] = mode
	}

	results, err := pagesCollection.Find(ctx, filter)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"message": "Failed to retrieve collection"})
		return
	}
	defer results.Close(ctx)

	var responseBody = make(map[string]interface{})

	responseBody["success"] = true
	responseBody["message"] = "page"

	pages := make(map[string]interface{}, 0)

	for results.Next(ctx) {
		var pgs entities.Pages
		err := results.Decode(&pgs)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"message": err.Error()})
			return
		}

		// // Modify pgs.Meta if needed
		// if pgs.Mode == "desktop" {
		// 	pgs.Meta.Title = ""
		// 	pgs.Meta.Description = ""
		// }

		pageMap := map[string]interface{}{
			"_id":       pgs.Id,
			"meta":      pgs.Meta,
			"mode":      pgs.Mode,
			"url":       pgs.Url,
			"createdAt": pgs.CreatedAt,
			"updatedAt": pgs.UpdatedAt,
			"__v":       pgs.V,
		}

		rows := make([]map[string]interface{}, 0)

		for _, rowID := range pgs.Rows {
			var row entities.Row
			err := rowCollection.FindOne(ctx, bson.M{"_id": rowID}).Decode(&row)
			if err != nil {
				c.JSON(http.StatusInternalServerError, gin.H{"message": "Failed to fetch row information"})
				return
			}

			rowMap := map[string]interface{}{
				"fluid":           row.Fluid,
				"backgroundColor": row.BackgroundColor,
				"createdAt":       row.CreatedAt,
				"updatedAt":       row.UpdatedAt,
				"pageId":          row.PageId,
				"__v":             row.V,
				"_id":             row.ID,
			}

			// ...

			cols := make([]map[string]interface{}, 0)

			for _, colID := range row.Cols {
				var col entities.Column

				err := colCollection.FindOne(c, bson.M{"_id": colID}).Decode(&col)
				if err != nil {
					log.Printf("Failed to fetch column information for column ID %d: %v", colID, err)
					c.JSON(http.StatusInternalServerError, gin.H{"message": err.Error()})
					return
				}

				//**************************************************************************8

				simplifiedContent := make(map[string]interface{})

				switch colContent := col.Content.(type) {

				case primitive.A:

					for _, item := range colContent {

						rowMap := make(map[string]interface{})

						switch item := item.(type) {
						case primitive.E:

							key := item.Key
							value := item.Value

							switch key {
							case "alt":
								if alt, ok := value.(string); ok {
									rowMap["alt"] = alt
								}
							case "link":
								if link, ok := value.(string); ok {
									rowMap["link"] = link
								}
							case "image":
								if image, ok := value.(map[string]interface{}); ok {
									imageUrl, _ := image["url"].(string)
									imageMap := map[string]interface{}{
										"url": imageUrl,
									}
									rowMap["image"] = imageMap
								}
							default:

							}
						default:

						}

						simplifiedContent = rowMap
					}

				case primitive.D:
					// Handle document case
					rowMap := make(map[string]interface{})

					for _, item := range colContent {

						key := item.Key
						value := item.Value

						switch key {
						case "alt":
							if alt, ok := value.(string); ok {
								rowMap["alt"] = alt
							}
						case "link":
							if link, ok := value.(string); ok {
								rowMap["link"] = link
							}
						case "image":
							switch value := value.(type) {
							case primitive.D:
								imageMap := make(map[string]interface{})
								for _, subItem := range value {
									subKey := subItem.Key
									subValue := subItem.Value

									switch subKey {
									case "url":
										if url, ok := subValue.(string); ok {
											imageMap["url"] = url
										}

									}
								}
								rowMap["image"] = imageMap
							}
						default:

						}
					}
					simplifiedContent = rowMap

					// simplifiedContent = append(simplifiedContent, map[string]interface{}{"content": rowMap})
				}

				colMap := map[string]interface{}{
					"_id":             col.ID,
					"size":            col.Size,
					"elevation":       col.Elevation,
					"padding":         col.Padding,
					"radius":          col.Radius,
					"margin":          col.Margin,
					"backgroundColor": col.BackgroundColor,
					"dataUrl":         col.DataUrl,
					"isMore":          col.IsMore,
					"dataType":        col.DataType,
					"layoutType":      col.LayoutType,
					"content":         simplifiedContent,
					"moreUrl":         col.MoreUrl,
					"name":            col.Name,
					"rowId":           col.RowId,
					"createdAt":       col.CreatedAt,
					"updatedAt":       col.UpdatedAt,
					"__v":             col.V,
				}

				// // Check the type of the content field
				// switch v := col.Content.(type) {
				// case []interface{}:
				// 	// If it's an array, assume it's the first type of content
				// 	colMap["content"] = v
				// case map[string]interface{}:
				// 	// If it's a map, assume it's the second type of content
				// 	colMap["content"] = []interface{}{v}
				// }

				// colMap := entities.Column{
				// 	ID:              col.ID,
				// 	Size:            col.Size,
				// 	Padding:         col.Padding,
				// 	Radius:          col.Radius,
				// 	Margin:          col.Margin,
				// 	BackgroundColor: col.BackgroundColor,
				// 	DataUrl:         col.DataUrl,
				// 	IsMore:          col.IsMore,
				// 	DataType:        col.DataType,
				// 	LayoutType:      col.LayoutType,
				// 	Content:         col.Content,
				// 	MoreUrl:         col.MoreUrl,
				// 	Name:            col.Name,
				// 	RowId:           col.RowId,
				// 	CreatedAt:       col.CreatedAt,
				// 	UpdatedAt:       col.UpdatedAt,
				// 	V:               col.V,
				// }

				cols = append(cols, colMap)

			}

			// ...

			rowMap["cols"] = cols
			rows = append(rows, rowMap)

		}

		pageMap["rows"] = rows
		pages = pageMap
	}

	responseBody["body"] = pages

	c.JSON(http.StatusCreated, responseBody)
}

//******************************************************************

//*********************************************************************

// func GetPages(c *gin.Context) {
// 	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
// 	defer cancel()

// 	mode := c.Query("mode") // Get the mode parameter from the query

// 	filter := bson.M{}

// 	if mode != "" {
// 		filter["mode"] = mode
// 	}

// 	results, err := pagesCollection.Find(ctx, filter)
// 	if err != nil {
// 		c.JSON(http.StatusInternalServerError, gin.H{"message": "Failed to retrieve collection"})
// 		return
// 	}
// 	defer results.Close(ctx)

// 	var pagesResponse []entities.Pages

// 	for results.Next(ctx) {
// 		var pgs entities.Pages
// 		err := results.Decode(&pgs)
// 		if err != nil {
// 			c.JSON(http.StatusInternalServerError, gin.H{"message": err.Error()})
// 			return
// 		}

// 		page := entities.Pages{
// 			Id:        pgs.Id,
// 			Meta:      pgs.Meta,
// 			Mode:      pgs.Mode,
// 			Url:       pgs.Url,
// 			CreatedAt: pgs.CreatedAt,
// 			UpdatedAt: pgs.UpdatedAt,
// 			V:         pgs.V,
// 			Rows:      pgs.Rows,
// 		}
// 		var rows []entities.Row

// 		for _, rowID := range pgs.Rows {
// 			var row entities.Row
// 			err := rowCollection.FindOne(ctx, bson.M{"_id": rowID}).Decode(&row)
// 			if err != nil {
// 				c.JSON(http.StatusInternalServerError, gin.H{"message": "Failed to fetch row information"})
// 				return
// 			}

// 			rowPresentation := entities.Row{
// 				Fluid:           row.Fluid,
// 				BackgroundColor: row.BackgroundColor,
// 				CreatedAt:       row.CreatedAt,
// 				UpdatedAt:       row.UpdatedAt,
// 				PageId:          row.PageId,
// 				V:               row.V,
// 				ID:              row.ID,
// 			}
// 			var cols []entities.Column

// 			for _, colID := range row.Cols {
// 				var col entities.Column
// 				err := colCollection.FindOne(ctx, bson.M{"_id": colID}).Decode(&col)
// 				if err != nil {
// 					log.Printf("Failed to fetch column information for column ID %+v: %v", colID, err)
// 					c.JSON(http.StatusInternalServerError, gin.H{"message": "Failed to fetch column information"})
// 					return
// 				}

// 				colPresentation := entities.Column{
// 					ID:              col.ID,
// 					Size:            col.Size,
// 					Elevation:       col.Elevation,
// 					Padding:         col.Padding,
// 					Radius:          col.Radius,
// 					Margin:          col.Margin,
// 					BackgroundColor: col.BackgroundColor,
// 					DataUrl:         col.DataUrl,
// 					IsMore:          col.IsMore,
// 					LayoutType:      col.LayoutType,
// 					DataType:        col.DataType,
// 					MoreUrl:         col.MoreUrl,
// 					Content:         col.Content,
// 					Name:            col.Name,
// 					RowId:           col.RowId,
// 					CreatedAt:       col.CreatedAt,
// 					UpdatedAt:       col.UpdatedAt,
// 					V:               col.V,
// 				}

// 				cols = append(cols, colPresentation)
// 			}
// 			//rowPresentation["cols"]=cols

// 			rows = append(rows, rowPresentation)
// 		}

// 		pagesResponse = append(pagesResponse, page)
// 	}

// 	c.JSON(http.StatusCreated, pagesResponse)
// }
