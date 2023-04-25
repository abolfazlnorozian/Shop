package services

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"shop/db"
	"shop/entity"
	"shop/middleware"
	"time"

	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
)

var categoryCollection *mongo.Collection = db.GetCollection(db.DB, "categories")

func iterateChild(resChild []entity.Category, obj *entity.Category) {
	var ctx, cancel = context.WithTimeout(context.Background(), 100*time.Second)
	// fmt.Println("len child: ", len(resChild))
	defer cancel()

	obj.Children = make([]entity.Category, len(resChild))
	if len(resChild) > 0 {
		jsonData, _ := json.Marshal(resChild)
		fmt.Println("resChild: ", string(jsonData))
		obj.Children = append([]entity.Category{}, resChild...)
		for _, sv := range resChild {
			var resSc []entity.Category
			fc := bson.M{"parent": sv.Parent.Hex()}
			cursor, err := categoryCollection.Find(ctx, fc)
			defer cancel()
			if err != nil {
				fmt.Println("halo")
			}
			if err = cursor.All(ctx, &resSc); err != nil {
				fmt.Println("loha")
			}

			//iterateChild(resSc, &sv)
		}
	}
}

func FindAllCategories(c *gin.Context) {

	var ctx, cancel = context.WithTimeout(context.Background(), 100*time.Second)
	var response []*entity.Category
	var resChild []entity.Category
	// q := c.Query("q")
	// perPage, err := strconv.Atoi(c.Query("per_page"))
	// if err != nil || perPage == 0 {
	// 	perPage = 10
	// }
	// page, _ := strconv.Atoi(c.Query("page"))
	// if page > 0 { // page start from index 0
	// 	page -= 1
	// }
	// filter := bson.M{
	// 	"name": primitive.Regex{
	// 		Pattern: q,
	// 		Options: "i", // i for case sensitive uppercase or lowercase
	// 	},
	// 	"is_parent": true,
	// 	"parent":    "",
	// }
	//options := options.Find()
	// options.SetSkip(int64(page * perPage))
	// options.SetLimit(int64(perPage))
	// options.SetSort(bson.D{{Key: "__v", Value: 1}})
	// cursor, err := categoryCollection.Find(ctx, filter, options)
	// defer cancel()
	// if err != nil {
	// 	c.JSON(http.StatusInternalServerError, gin.H{
	// 		"status":  false,
	// 		"message": http.StatusInternalServerError,
	// 		"error":   err.Error(),
	// 	})
	// 	return
	// }
	// if err = cursor.All(ctx, &response); err != nil {
	// 	log.Fatal(err)
	// 	c.JSON(http.StatusInternalServerError, gin.H{
	// 		"status":  false,
	// 		"message": err.Error(),
	// 	})
	// 	return
	// }

	for _, v := range response {
		// fmt.Println(v.ID.Hex())
		filterChild := bson.D{bson.E{Key: "parent", Value: v.Parent}}
		cursor, err := categoryCollection.Find(ctx, filterChild)
		defer cancel()
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{
				"status":  false,
				"message": http.StatusInternalServerError,
				"error":   err.Error(),
			})
			return
		}
		if err = cursor.All(ctx, &resChild); err != nil {
			log.Fatal(err)
			c.JSON(http.StatusInternalServerError, gin.H{
				"status":  false,
				"message": err.Error(),
			})
			return
		}
		//iterateChild(resChild, v)
	}
	if len(response) > 0 {
		c.JSON(http.StatusOK, gin.H{
			"status":  true,
			"message": "Success retrieve data",
			"data":    response,
		})
	} else {
		c.JSON(http.StatusOK, gin.H{
			"status":  false,
			"message": "Data Not Found",
			"data":    response,
		})
	}

	//************************************************************

	// var categories []*entity.Category
	// //var resChild []entity.Category

	// results, err := categoryCollection.Find(c, bson.D{{}})
	// if err != nil {
	// 	c.JSON(http.StatusInternalServerError, gin.H{"massage": "Not Find Collection"})
	// 	return
	// }

	// for results.Next(c) {
	// 	var title entity.Category

	// 	err := results.Decode(&title)
	// 	if err != nil {
	// 		c.JSON(http.StatusInternalServerError, gin.H{"message": err.Error()})
	// 		return

	// 	}

	// 	categories = append(categories, &title)

	// }
	// if err := results.Err(); err != nil {
	// 	if err != nil {
	// 		c.JSON(http.StatusInternalServerError, gin.H{"message": "results has error"})
	// 		return
	// 	}
	// 	results.Close(c)
	// 	if len(categories) == 0 {
	// 		c.JSON(http.StatusInternalServerError, gin.H{"message": "document not found"})
	// 		return
	// 	}
	// 	return

	// }

	// c.JSON(http.StatusOK, response.Response{Status: http.StatusOK, Message: "success", Data: map[string]interface{}{"data": categories}})
}

func AddCategories(c *gin.Context) {

	var title entity.Category

	if err := c.ShouldBindJSON(&title); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	if err := middleware.CheckUserType(c, "admin"); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	if validationErr := validate.Struct(&title); validationErr != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": validationErr.Error()})
		return
	}

	_, err := categoryCollection.InsertOne(c, title)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": title})

}
