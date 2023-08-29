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
	"go.mongodb.org/mongo-driver/mongo"
)

var pagesCollection *mongo.Collection = database.GetCollection(database.DB, "pages")
var rowCollection *mongo.Collection = database.GetCollection(database.DB, "rows")
var colCollection *mongo.Collection = database.GetCollection(database.DB, "columns")

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

	var responseBody []map[string]interface{} // Create a slice of maps for building the response body

	for results.Next(ctx) {
		var pgs entities.Pages
		err := results.Decode(&pgs)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"message": err.Error()})
			return
		}

		// Modify pgs.Meta if needed
		if pgs.Mode == "desktop" {
			pgs.Meta.Title = ""
			pgs.Meta.Description = ""
		}

		var formattedRows []map[string]interface{}
		for _, rowID := range pgs.Rows {
			var row entities.Row
			err := rowCollection.FindOne(ctx, bson.M{"_id": rowID}).Decode(&row)
			if err != nil {
				c.JSON(http.StatusInternalServerError, gin.H{"message": "Failed to fetch row information"})
				return
			}

			var formattedCols []map[string]interface{}
			for _, colID := range row.Cols {
				var col entities.Column
				err := colCollection.FindOne(ctx, bson.M{"_id": colID}).Decode(&col)
				if err != nil {
					log.Printf("Failed to fetch column information for column ID %d: %v", colID, err)
					c.JSON(http.StatusInternalServerError, gin.H{"message": "Failed to fetch column information"})
					return
				}

				formattedCol := map[string]interface{}{
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
					"moreUrl":         col.MoreUrl,
					"content":         col.Content,
					"name":            col.Name,
					"rowId":           col.RowId,
					"createdAt":       col.CreatedAt,
					"updatedAt":       col.UpdatedAt,
					"__v":             col.V,
				}

				formattedCols = append(formattedCols, formattedCol)
			}

			formattedRow := map[string]interface{}{
				"_id":             row.ID,
				"fluid":           row.Fluid,
				"backgroundColor": row.BackgroundColor,
				"cols":            formattedCols,
				"pageId":          row.PageId,
				"createdAt":       row.CreatedAt,
				"updatedAt":       row.UpdatedAt,
				"__v":             row.V,
			}

			formattedRows = append(formattedRows, formattedRow)
		}

		pageMap := map[string]interface{}{
			"_id":       pgs.Id,
			"meta":      pgs.Meta,
			"mode":      pgs.Mode,
			"rows":      formattedRows,
			"url":       pgs.Url,
			"createdAt": pgs.CreatedAt,
			"updatedAt": pgs.UpdatedAt,
			"__v":       pgs.V,
		}

		responseBody = append(responseBody, pageMap)
	}

	c.JSON(http.StatusOK, gin.H{"success": true, "message": "page", "body": responseBody})
}
