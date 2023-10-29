package services

import (
	"encoding/json"
	"io/ioutil"
	"net/http"
	"shop/entities"

	"github.com/gin-gonic/gin"
)

// func State(c *gin.Context) {
// 	// Manually adjust the requested path to include the double slash
// 	requestedPath := "//json/state.json"

// 	// Construct the actual file path
// 	filePath := "." + requestedPath

// 	// Check if the file exists
// 	_, err := os.Stat(filePath)
// 	if err != nil {
// 		c.JSON(http.StatusNotFound, gin.H{"error": "File not found"})
// 		return
// 	}

// 	// Open and read the file
// 	file, err := os.Open(filePath)
// 	if err != nil {
// 		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
// 		return
// 	}
// 	defer file.Close()

// 	// Serve the file content
// 	c.Header("Content-Type", "application/json")
// 	c.Status(http.StatusOK)

// 	_, err = io.Copy(c.Writer, file)
// 	if err != nil {
// 		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
// 	}
// }

// ***********************************************************************8

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

// func State(c *gin.Context) {
// 	path := c.Param("path")
// 	// Process the path as needed
// 	if strings.HasPrefix(path, "json/state.json") {
// 		fileContent, err := ioutil.ReadFile("state.json")
// 		if err != nil {
// 			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
// 			return
// 		}
// 		var states []entities.State
// 		err = json.Unmarshal(fileContent, &states)
// 		if err != nil {
// 			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
// 			return
// 		}

// 		c.JSON(http.StatusOK, states)
// 	} else {
// 		c.JSON(http.StatusNotFound, gin.H{"error": "Not Found"})
// 	}
// }

//******************************************************************************8

// func State(c *gin.Context) {
// 	path := c.Param("path")
// 	if path != "json/state.json" {
// 		c.JSON(http.StatusNotFound, gin.H{"error": "Not Found"})
// 		return
// 	}

// 	fileContent, err := ioutil.ReadFile("state.json")
// 	if err != nil {
// 		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
// 		return
// 	}

// 	var states []entities.State
// 	err = json.Unmarshal(fileContent, &states)
// 	if err != nil {
// 		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
// 		return
// 	}

// 	c.JSON(http.StatusOK, states)
// }
//***************************************************************
// func State(c *gin.Context) {
// 	path := c.Param("path")

// 	path = strings.TrimRight(path, "/")

// 	if path != "//json/state.json" {
// 		c.JSON(http.StatusNotFound, gin.H{"error": "Not Found"})
// 		return
// 	}

// 	fileContent, err := ioutil.ReadFile("state.json")
// 	if err != nil {
// 		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
// 		return
// 	}

// 	var states []entities.State
// 	err = json.Unmarshal(fileContent, &states)
// 	if err != nil {
// 		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
// 		return
// 	}

// 	c.JSON(http.StatusOK, states)
// }

//********************************************************8
// func State(c *gin.Context) {
// 	path := c.Param("path")

// 	// Check if the path ends with "json/state.json"
// 	if strings.HasSuffix(path, "json/state.json") {
// 		fileContent, err := ioutil.ReadFile("state.json")
// 		if err != nil {
// 			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
// 			return
// 		}

// 		var states []entities.State
// 		err = json.Unmarshal(fileContent, &states)
// 		if err != nil {
// 			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
// 			return
// 		}

// 		c.JSON(http.StatusOK, states)
// 	} else {
// 		c.JSON(http.StatusNotFound, gin.H{"error": "Not Found"})
// 	}
// }
