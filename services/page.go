package services

import (
	"context"
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

// func GetPages(c *gin.Context) {
// 	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
// 	defer cancel()

// 	mode := c.Query("mode") // Get the mode parameter from the query

// 	var pages []entities.Pages
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

// 	for results.Next(ctx) {
// 		var pgs entities.Pages
// 		err := results.Decode(&pgs)
// 		if err != nil {
// 			c.JSON(http.StatusInternalServerError, gin.H{"message": err.Error()})
// 			return
// 		}
// 		var rowInfo []entities.Row // Modify "entities.Row" to match your actual row structure
// 		for _, rowNum := range pgs.Rows {
// 			// Fetch the row information from the "rows" collection based on rowNum
// 			rowFilter := bson.M{"_id": rowNum} // Use _id instead of "row"
// 			rowResult := rowCollection.FindOne(ctx, rowFilter)
// 			if rowResult.Err() != nil {
// 				c.JSON(http.StatusInternalServerError, gin.H{"message": "Failed to retrieve row information"})
// 				return
// 			}
// 			var row entities.Row
// 			if err := rowResult.Decode(&row); err != nil {
// 				c.JSON(http.StatusInternalServerError, gin.H{"message": err.Error()})
// 				return
// 			}
// 			rowInfo = append(rowInfo, row)
// 		}

// 		pgs.RowsInfo = rowInfo // Assign the row information to the new RowsInfo field

// 		if pgs.Mode == "desktop" {
// 			pgs.Meta.Title = ""
// 			pgs.Meta.Description = ""
// 		}
// 		pages = append(pages, pgs)
// 	}

// 	c.JSON(http.StatusOK, gin.H{"success": true, "massage": "page", "body": pages})
// }
// func GetPages(c *gin.Context) {
// 	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
// 	defer cancel()

// 	mode := c.Query("mode") // Get the mode parameter from the query

// 	var pages []entities.Pages
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

// 	for results.Next(ctx) {
// 		var pgs entities.Pages
// 		err := results.Decode(&pgs)
// 		if err != nil {
// 			c.JSON(http.StatusInternalServerError, gin.H{"message": err.Error()})
// 			return
// 		}

// 		if pgs.Mode == "desktop" {
// 			pgs.Meta.Title = ""
// 			pgs.Meta.Description = ""
// 		}

// 		// Fetch and populate row information for each row ID
// 		var rows []entities.Row
// 		for _, rowID := range pgs.Rows {
// 			var row entities.Row
// 			err := rowCollection.FindOne(ctx, bson.M{"_id": rowID}).Decode(&row)
// 			if err != nil {
// 				c.JSON(http.StatusInternalServerError, gin.H{"message": "Failed to fetch row information"})
// 				return
// 			}
// 			rows = append(rows, row)
// 		}

// 		// Replace Rows field with the fetched row information
// 		pgs.Rows = rows

// 		pages = append(pages, pgs)
// 	}

// 	c.JSON(http.StatusOK, gin.H{"success": true, "message": "page", "body": pages})
// }
// func GetPages(c *gin.Context) {
// 	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
// 	defer cancel()

// 	mode := c.Query("mode") // Get the mode parameter from the query

// 	var pages []entities.Pages
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

// 	for results.Next(ctx) {
// 		var pgs entities.Pages
// 		err := results.Decode(&pgs)
// 		if err != nil {
// 			c.JSON(http.StatusInternalServerError, gin.H{"message": err.Error()})
// 			return
// 		}

// 		if pgs.Mode == "desktop" {
// 			pgs.Meta.Title = ""
// 			pgs.Meta.Description = ""
// 		}

// 		// Fetch and populate row information for each row ID
// 		var rows []entities.Row
// 		for _, rowID := range pgs.Rows {
// 			var row entities.Row
// 			err := rowCollection.FindOne(ctx, bson.M{"_id": rowID}).Decode(&row)
// 			if err != nil {
// 				c.JSON(http.StatusInternalServerError, gin.H{"message": "Failed to fetch row information"})
// 				return
// 			}
// 			rows = append(rows, row)
// 		}

// 		// Replace Rows field with the fetched row IDs
// 		rowIDs := make([]int, len(rows))
// 		for i, row := range rows {
// 			rowIDs[i] = row.ID // Assuming Row.ID is of type int
// 		}
// 		pgs.Rows = rowIDs

// 		pages = append(pages, pgs)
// 	}

// 	c.JSON(http.StatusOK, gin.H{"success": true, "message": "page", "body": pages})
// }

func GetPages(c *gin.Context) {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	mode := c.Query("mode") // Get the mode parameter from the query

	//var pgs []entities.Pages
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

		if pgs.Mode == "desktop" {
			pgs.Meta.Title = ""
			pgs.Meta.Description = ""
		}

		// Fetch and populate row information for each row ID
		var rows []entities.Row
		for _, rowID := range pgs.Rows {
			var row entities.Row
			err := rowCollection.FindOne(ctx, bson.M{"_id": rowID}).Decode(&row)
			if err != nil {
				c.JSON(http.StatusInternalServerError, gin.H{"message": "Failed to fetch row information"})
				return
			}
			rows = append(rows, row)
		}

		// Convert the rows to the desired format
		var formattedRows []map[string]interface{}
		for _, row := range rows {
			formattedRow := map[string]interface{}{
				"_id":             row.ID,
				"fluid":           row.Fluid,
				"backgroundColor": row.BackGroundColor,
				"cols":            row.Cols,
				"pageId":          row.PageId,
				"createdAt":       row.CreatedAt,
				"updatedAt":       row.UpdatedAt,
				"__v":             row.V,
			}
			formattedRows = append(formattedRows, formattedRow)
		}

		// Create a map for the current page
		pageMap := map[string]interface{}{
			"_id":       pgs.Id,
			"meta":      pgs.Meta,
			"mode":      pgs.Mode,
			"rows":      formattedRows, // Replace pgs.Rows with formattedRows
			"url":       pgs.Url,
			"createdAt": pgs.CreatedAt,
			"updatedAt": pgs.UpdatedAt,
			"__v":       pgs.V,
		}

		responseBody = append(responseBody, pageMap)
	}

	// Send the response with the formatted JSON structure
	c.JSON(http.StatusOK, gin.H{"success": true, "message": "page", "body": responseBody})
}
