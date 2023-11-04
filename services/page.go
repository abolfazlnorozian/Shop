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

		// Modify pgs.Meta if needed
		if pgs.Mode == "desktop" {
			pgs.Meta.Title = ""
			pgs.Meta.Description = ""
		}

		pageMap := map[string]interface{}{
			"_id":  pgs.Id,
			"meta": pgs.Meta,
			"mode": pgs.Mode,
			"url":  pgs.Url,
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
			}

			// ...

			cols := make([]map[string]interface{}, 0)

			for _, colID := range row.Cols {
				var col entities.Column
				err := colCollection.FindOne(ctx, bson.M{"_id": colID}).Decode(&col)
				if err != nil {
					log.Printf("Failed to fetch column information for column ID %d: %v", colID, err)
					c.JSON(http.StatusInternalServerError, gin.H{"message": "Failed to fetch column information"})
					return
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
					"moreUrl":         col.MoreUrl,
					"name":            col.Name,
					"rowId":           col.RowId,
					"createdAt":       col.CreatedAt,
					"updatedAt":       col.UpdatedAt,
					"__v":             col.V,
				}

				// Handle different content structures
				var contentData []interface{}

				if col.Content != nil {
					// Check if the "content" is an array of objects
					if contentArray, ok := col.Content.([]entities.Content); ok {
						for _, item := range contentArray {
							// Convert each item into an object without keys
							contentItem := map[string]interface{}{
								"alt":   item.Alt,
								"link":  item.Link,
								"image": item.Image,
							}
							contentData = append(contentData, contentItem)
						}
					} else if contentObject, ok := col.Content.(entities.Content); ok {
						// Convert the content object into an object without keys
						contentItem := map[string]interface{}{
							"alt":   contentObject.Alt,
							"link":  contentObject.Link,
							"image": contentObject.Image,
						}
						contentData = append(contentData, contentItem)
					} else {
						// Handle other cases where "content" is not an array or object
						contentData = append(contentData, col.Content)
					}
				} else {
					// If "content" is null in the database, set it to null in the response
					contentData = nil
				}

				colMap["content"] = contentData
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

// 	mode := c.Query("mode")

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

// 	var responseBody = make(map[string]interface{})

// 	responseBody["success"] = true
// 	responseBody["message"] = "page"

// 	var pages []interface{}

// 	for results.Next(ctx) {
// 		var pgs entities.Pages
// 		err := results.Decode(&pgs)
// 		if err != nil {
// 			c.JSON(http.StatusInternalServerError, gin.H{"message": err.Error()})
// 			return
// 		}

// 		// Modify pgs.Meta if needed
// 		if pgs.Mode == "desktop" {
// 			pgs.Meta.Title = ""
// 			pgs.Meta.Description = ""
// 		}

// 		page := map[string]interface{}{
// 			"meta": pgs.Meta,
// 			"mode": pgs.Mode,
// 			"rows": []interface{}{}, // Initialize rows array
// 		}

// 		for _, rowID := range pgs.Rows {
// 			var row entities.Row
// 			err := rowCollection.FindOne(ctx, bson.M{"_id": rowID}).Decode(&row)
// 			if err != nil {
// 				c.JSON(http.StatusInternalServerError, gin.H{"message": "Failed to fetch row information"})
// 				return
// 			}

// 			rowMap := map[string]interface{}{
// 				"fluid":           row.Fluid,
// 				"backgroundColor": row.BackgroundColor,
// 				"cols":            []interface{}{}, // Initialize cols array
// 			}

// 			for _, colID := range row.Cols {
// 				log.Printf("Fetching data for column ID: %d", colID) // Add this line
// 				var col entities.Column
// 				err := colCollection.FindOne(ctx, bson.M{"_id": colID}).Decode(&col)
// 				if err != nil {
// 					log.Printf("Failed to fetch column information for column ID %d: %v", colID, err)
// 					c.JSON(http.StatusInternalServerError, gin.H{"message": "Failed to fetch column information"})
// 					return
// 				}

// 				colMap := map[string]interface{}{
// 					"size":            col.Size,
// 					"elevation":       col.Elevation,
// 					"padding":         col.Padding,
// 					"radius":          col.Radius,
// 					"margin":          col.Margin,
// 					"backgroundColor": col.BackgroundColor,
// 					"dataUrl":         col.DataUrl,
// 					"isMore":          col.IsMore,
// 					"dataType":        col.DataType,
// 					"layoutType":      col.LayoutType,
// 					"moreUrl":         col.MoreUrl,
// 					"name":            col.Name,
// 					"rowId":           col.RowId,
// 					"createdAt":       col.CreatedAt,
// 					"updatedAt":       col.UpdatedAt,
// 					"__v":             col.V,
// 				}

// 				switch col.Content.(type) {
// 				case []interface{}:
// 					contentArray := col.Content.([]interface{})
// 					// Handle the case where contentArray is an array
// 					var contentData []interface{}
// 					for _, item := range contentArray {
// 						contentItem := item.(map[string]interface{})
// 						contentData = append(contentData, contentItem)
// 					}
// 					colMap["content"] = contentData

// 				case map[string]interface{}:
// 					// Handle the case where col.Content is an object
// 					contentData := col.Content.(map[string]interface{})
// 					colMap["content"] = contentData
// 				}

// 				rowMap["cols"] = append(rowMap["cols"].([]interface{}), colMap)
// 			}

// 			page["rows"] = append(page["rows"].([]interface{}), rowMap)
// 		}

// 		pages = append(pages, page)
// 	}

// 	responseBody["body"] = map[string]interface{}{
// 		"pages": pages,
// 	}

// 	c.JSON(http.StatusOK, responseBody)
// }

//**************************************************************************************************************************

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

// 	var responseBody []map[string]interface{} // Create a slice of maps for building the response body

// 	for results.Next(ctx) {
// 		var pgs entities.Pages
// 		err := results.Decode(&pgs)
// 		if err != nil {
// 			c.JSON(http.StatusInternalServerError, gin.H{"message": err.Error()})
// 			return
// 		}

// 		// Modify pgs.Meta if needed
// 		if pgs.Mode == "desktop" {
// 			pgs.Meta.Title = ""
// 			pgs.Meta.Description = ""
// 		}

// 		var formattedRows []map[string]interface{}
// 		for _, rowID := range pgs.Rows {
// 			var row entities.Row
// 			err := rowCollection.FindOne(ctx, bson.M{"_id": rowID}).Decode(&row)
// 			if err != nil {
// 				c.JSON(http.StatusInternalServerError, gin.H{"message": "Failed to fetch row information"})
// 				return
// 			}

// 			var formattedCols []map[string]interface{}
// 			for _, colID := range row.Cols {
// 				var col entities.Column
// 				err := colCollection.FindOne(ctx, bson.M{"_id": colID}).Decode(&col)
// 				if err != nil {
// 					log.Printf("Failed to fetch column information for column ID %d: %v", colID, err)
// 					c.JSON(http.StatusInternalServerError, gin.H{"message": "Failed to fetch column information"})
// 					return
// 				}

// 				formattedCol := map[string]interface{}{
// 					"_id":             col.ID,
// 					"size":            col.Size,
// 					"elevation":       col.Elevation,
// 					"padding":         col.Padding,
// 					"radius":          col.Radius,
// 					"margin":          col.Margin,
// 					"backgroundColor": col.BackgroundColor,
// 					"dataUrl":         col.DataUrl,
// 					"isMore":          col.IsMore,
// 					"dataType":        col.DataType,
// 					"layoutType":      col.LayoutType,
// 					"moreUrl":         col.MoreUrl,
// 					"content":         col.Content,
// 					"name":            col.Name,
// 					"rowId":           col.RowId,
// 					"createdAt":       col.CreatedAt,
// 					"updatedAt":       col.UpdatedAt,
// 					"__v":             col.V,
// 				}

// 				formattedCols = append(formattedCols, formattedCol)
// 			}

// 			formattedRow := map[string]interface{}{
// 				"_id":             row.ID,
// 				"fluid":           row.Fluid,
// 				"backgroundColor": row.BackgroundColor,
// 				"cols":            formattedCols,
// 				"pageId":          row.PageId,
// 				"createdAt":       row.CreatedAt,
// 				"updatedAt":       row.UpdatedAt,
// 				"__v":             row.V,
// 			}

// 			formattedRows = append(formattedRows, formattedRow)
// 		}

// 		pageMap := map[string]interface{}{
// 			"_id":       pgs.Id,
// 			"meta":      pgs.Meta,
// 			"mode":      pgs.Mode,
// 			"rows":      formattedRows,
// 			"url":       pgs.Url,
// 			"createdAt": pgs.CreatedAt,
// 			"updatedAt": pgs.UpdatedAt,
// 			"__v":       pgs.V,
// 		}

// 		responseBody = append(responseBody, pageMap)
// 	}

// 	c.JSON(http.StatusOK, gin.H{"success": true, "message": "page", "body": responseBody})
// }

//**********************************************************************************************

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

// 	var responseBody gin.H // Create a map for building the response body

// 	var rows []gin.H // Create a slice of maps for rows

// 	for results.Next(ctx) {
// 		var pgs entities.Pages
// 		err := results.Decode(&pgs)
// 		if err != nil {
// 			c.JSON(http.StatusInternalServerError, gin.H{"message": err.Error()})
// 			return
// 		}

// 		var formattedRows []gin.H
// 		for _, rowID := range pgs.Rows {
// 			var row entities.Row
// 			err := rowCollection.FindOne(ctx, bson.M{"_id": rowID}).Decode(&row)
// 			if err != nil {
// 				c.JSON(http.StatusInternalServerError, gin.H{"message": "Failed to fetch row information"})
// 				return
// 			}

// 			var formattedCols []gin.H
// 			for _, colID := range row.Cols {
// 				var col entities.Column
// 				err := colCollection.FindOne(ctx, bson.M{"_id": colID}).Decode(&col)
// 				if err != nil {
// 					log.Printf("Failed to fetch column information for column ID %d: %v", colID, err)
// 					c.JSON(http.StatusInternalServerError, gin.H{"message": "Failed to fetch column information"})
// 					return
// 				}
// 				fmt.Printf("Column Content: %+v\n", col.Content)

// 				content, err := handleColumnContent(col.Content)
// 				if err != nil {
// 					c.JSON(http.StatusInternalServerError, gin.H{"message": "Failed to handle column content"})
// 					return
// 				}

// 				colMap := gin.H{
// 					"_id":             col.ID,
// 					"size":            col.Size,
// 					"elevation":       col.Elevation,
// 					"padding":         col.Padding,
// 					"radius":          col.Radius,
// 					"margin":          col.Margin,
// 					"backgroundColor": col.BackgroundColor,
// 					"dataUrl":         col.DataUrl,
// 					"isMore":          col.IsMore,
// 					"dataType":        col.DataType,
// 					"layoutType":      col.LayoutType,
// 					"moreUrl":         col.MoreUrl,
// 					"content":         content,
// 					"name":            col.Name,
// 					"rowId":           col.RowId,
// 					"createdAt":       col.CreatedAt,
// 					"updatedAt":       col.UpdatedAt,
// 					"__v":             col.V,
// 				}

// 				formattedCols = append(formattedCols, colMap)
// 			}

// 			rowMap := gin.H{
// 				"_id":             row.ID,
// 				"fluid":           row.Fluid,
// 				"backgroundColor": row.BackgroundColor,
// 				"cols":            formattedCols,
// 				"pageId":          row.PageId,
// 				"createdAt":       row.CreatedAt,
// 				"updatedAt":       row.UpdatedAt,
// 				"__v":             row.V,
// 			}

// 			formattedRows = append(formattedRows, rowMap)
// 		}

// 		pgsMap := gin.H{
// 			"_id":       pgs.Id,
// 			"meta":      pgs.Meta,
// 			"mode":      pgs.Mode,
// 			"rows":      formattedRows,
// 			"url":       pgs.Url,
// 			"createdAt": pgs.CreatedAt,
// 			"updatedAt": pgs.UpdatedAt,
// 			"__v":       pgs.V,
// 		}

// 		rows = append(rows, pgsMap)
// 	}

// 	responseBody = gin.H{
// 		"success": true,
// 		"message": "page",
// 		"body": gin.H{
// 			"meta": gin.H{
// 				"keywords": []interface{}{}, // You can update this with actual keywords
// 			},
// 			"mode": mode, // Update this with the actual mode
// 			"rows": rows,
// 		},
// 	}

// 	c.JSON(http.StatusOK, responseBody)
// }

//************************************************************************************

// func handleColumnContent(content interface{}) ([]entities.Content, error) {
// 	var contentData []entities.Content

// 	// Check if the content is an array
// 	if contentArray, isArray := content.([]interface{}); isArray {
// 		// Handle array content
// 		for _, item := range contentArray {
// 			if itemMap, isMap := item.(map[string]interface{}); isMap {
// 				// Process itemMap as needed
// 				// Assuming the itemMap structure is similar to your Content struct
// 				contentItem, err := processContentItem(itemMap)
// 				if err != nil {
// 					return nil, err
// 				}
// 				contentData = append(contentData, contentItem)
// 			}
// 		}
// 	} else if contentMap, isMap := content.(map[string]interface{}); isMap {
// 		// Handle object content
// 		contentItem, err := processContentItem(contentMap)
// 		if err != nil {
// 			return nil, err
// 		}
// 		contentData = append(contentData, contentItem)
// 	}

// 	// Return the extracted content data
// 	return contentData, nil
// }

// func processContentItem(contentMap map[string]interface{}) (entities.Content, error) {
// 	var contentItem entities.Content

// 	if alt, ok := contentMap["alt"].(string); ok {
// 		contentItem.Alt = alt
// 	}

// 	if link, ok := contentMap["link"].(string); ok {
// 		contentItem.Link = link
// 	}

// 	if imageObj, ok := contentMap["image"].(map[string]interface{}); ok {
// 		if url, ok := imageObj["url"].(string); ok {
// 			contentItem.Image.URL = url
// 		}

// 		if id, ok := imageObj["_id"].(string); ok {
// 			contentItem.Image.Id = id
// 		}
// 	}

// 	return contentItem, nil
// }

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

// 	var responseBody []entities.Pages // Create a slice of entities.Pages for building the response body

// 	for results.Next(ctx) {
// 		var pgs entities.Pages
// 		err := results.Decode(&pgs)
// 		if err != nil {
// 			c.JSON(http.StatusInternalServerError, gin.H{"message": err.Error()})
// 			return
// 		}

// 		// Modify pgs.Meta if needed
// 		if pgs.Mode == "desktop" {
// 			pgs.Meta.Title = ""
// 			pgs.Meta.Description = ""
// 		}

// 		var formattedRows []entities.Row
// 		for _, rowID := range pgs.Rows {
// 			var row entities.Row
// 			err := rowCollection.FindOne(ctx, bson.M{"_id": rowID}).Decode(&row)
// 			if err != nil {
// 				c.JSON(http.StatusInternalServerError, gin.H{"message": "Failed to fetch row information"})
// 				return
// 			}

// 			var formattedCols []entities.Column
// 			for _, colID := range row.Cols {
// 				var col entities.Column
// 				err := colCollection.FindOne(ctx, bson.M{"_id": colID}).Decode(&col)
// 				if err != nil {
// 					log.Printf("Failed to fetch column information for column ID %d: %v", colID, err)
// 					c.JSON(http.StatusInternalServerError, gin.H{"message": "Failed to fetch column information"})
// 					return
// 				}

// 				content, err := handleColumnContent(col.Content)
// 				if err != nil {
// 					c.JSON(http.StatusInternalServerError, gin.H{"message": "Failed to handle column content"})
// 					return
// 				}

// 				col.Content = content
// 				formattedCols = append(formattedCols, col)
// 			}

// 			// Create a slice of int to hold the column IDs
// 			var colIDs []int

// 			// Populate the colIDs slice with the IDs from formattedCols
// 			for _, col := range formattedCols {
// 				colIDs = append(colIDs, col.ID)
// 			}

// 			// Assign colIDs to row.Cols
// 			row.Cols = colIDs

// 			formattedRows = append(formattedRows, row)
// 		}

// 		// Create a slice of int to hold the row IDs
// 		var rowIDs []int

// 		// Populate the rowIDs slice with the IDs from formattedRows
// 		for _, row := range formattedRows {
// 			rowIDs = append(rowIDs, row.ID)
// 		}

// 		// Assign rowIDs to pgs.Rows
// 		pgs.Rows = rowIDs

// 		responseBody = append(responseBody, pgs)
// 	}

// 	c.JSON(http.StatusOK, gin.H{"success": true, "message": "page", "body": responseBody})
// }
