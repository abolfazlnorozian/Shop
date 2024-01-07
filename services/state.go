package services

import (
	"encoding/json"
	"io/ioutil"
	"net/http"
	"shop/entities"

	"github.com/gin-gonic/gin"
)

func State(c *gin.Context) {
	fileContent, err := ioutil.ReadFile("state.json")
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	var states []entities.State
	err = json.Unmarshal(fileContent, &states)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, states)
}
