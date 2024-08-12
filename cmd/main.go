package main

import (
	"github.com/AlexChe360/go_api/models"
	"github.com/gin-gonic/gin"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"log"
	"net/http"
	"strconv"
)

var db *gorm.DB

func main() {
	r := gin.Default()

	var err error
	dns := "host=localhost user=postgres password=root dbname=go_api port=5432 sslmode=disable"
	db, err = gorm.Open(postgres.Open(dns), &gorm.Config{})
	failOnError(err, "Failed to connect to database")

	db.AutoMigrate(&models.Project{})

	r.GET("/ping", func(c *gin.Context) {
		c.JSON(200, gin.H{
			"message": "pong",
		})
	})

	r.POST("/projects", Create)
	r.GET("/projects", GetAll)
	r.PUT("/projects/:id", Update)
	r.DELETE("/projects/:id", Delete)

	r.Run()
}

func failOnError(err error, msg string) {
	if err != nil {
		log.Fatalf("%s: %s", msg, err)
	}
}

func Create(c *gin.Context) {
	var data models.Project

	if err := c.ShouldBindJSON(&data); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	project := &models.Project{Name: data.Name}
	result := db.Create(&project)

	if result.Error != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": result.Error.Error()})
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "create new",
		"data":    project,
	})
}

func GetAll(c *gin.Context) {
	var projects []models.Project
	result := db.Find(&projects)

	if result.Error != nil {
		c.JSON(http.StatusBadRequest, gin.H{"message": result.Error.Error()})
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "get all",
		"data":    projects,
	})
}

func Update(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid ID"})
		return
	}

	var payload models.Project

	if err := c.ShouldBindJSON(&payload); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request payload"})
		return
	}

	var project models.Project
	result := db.First(&project, id)
	if result.Error != nil {
		c.JSON(http.StatusBadRequest, gin.H{"message": result.Error.Error()})
	}

	project.Name = payload.Name
	db.Save(&project)

	c.JSON(http.StatusOK, gin.H{
		"message": "update project",
		"data":    project,
	})
}

func Delete(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid ID"})
		return
	}

	result := db.Delete(&models.Project{}, id)
	if result.Error != nil {
		c.JSON(http.StatusBadRequest, gin.H{"message": result.Error.Error()})
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "delete project",
	})

}
