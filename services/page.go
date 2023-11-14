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

func GetPages(c *gin.Context) {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	mode := c.Query("mode")

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

			cols := make([]map[string]interface{}, 0)

			for _, colID := range row.Cols {
				var col entities.Column

				err := colCollection.FindOne(c, bson.M{"_id": colID}).Decode(&col)
				if err != nil {
					log.Printf("Failed to fetch column information for column ID %d: %v", colID, err)
					c.JSON(http.StatusInternalServerError, gin.H{"message": err.Error()})
					return
				}

				var simplifiedContent interface{}

				switch colContent := col.Content.(type) {

				case primitive.A:
					contentSlice := make([]map[string]interface{}, 0)

					for _, item := range colContent {

						rowMap := make(map[string]interface{})

						switch item := item.(type) {
						case primitive.D:
							for _, kv := range item {
								key := kv.Key
								value := kv.Value
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
											case "_id":
												if id, ok := subValue.(string); ok {
													imageMap["id"] = id
												}

											}
										}
										rowMap["image"] = imageMap
									}
								default:

								}

							}

						}

						contentSlice = append(contentSlice, rowMap)

					}

					simplifiedContent = contentSlice

				case primitive.D:

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

				cols = append(cols, colMap)

			}

			rowMap["cols"] = cols
			rows = append(rows, rowMap)

		}

		pageMap["rows"] = rows
		pages = pageMap
	}

	responseBody["body"] = pages

	c.JSON(http.StatusCreated, responseBody)
}
