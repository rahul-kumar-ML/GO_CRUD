package main

import (
	"fmt"
	"net/http"
	"os"

	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

func main() {
	fmt.Println("Gin framework")

	//define router for gin
	router := gin.Default()

	//Load .env variables for credentials
	err_env := godotenv.Load()
	if err_env != nil {
		println("Error in fetching the credentials", err_env)
	}

	//connect to postgressDB
	dsn := os.Getenv("Credentials")
	// dsn := "host=localhost user=postgres password=root dbname=rahul port=5432 sslmode=disable TimeZone=Asia/Shanghai"
	fmt.Println(dsn)
	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
	fmt.Print(err, db)

	//Check connected? with time
	var time string = ""
	db.Raw("SELECT now()").Scan(&time)
	if time != "" {
		fmt.Println("DB Connected at", time)
	} else {
		fmt.Println("DB NOT Connected")
	}

	//CREATE USER PROFILE Operation
	router.POST("/adduser", func(c *gin.Context) {

		var requestBody struct {
			UserName string `json:"name"`
			Id       string `json:"id"`
		}
		err := c.BindJSON(&requestBody)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}
		c.JSON(200, gin.H{"data": requestBody.UserName + "User Added"})

		fmt.Println(requestBody)
		db.Exec("INSERT INTO users values(?,?);", requestBody.UserName, requestBody.Id)
	})

	//READ USERNAME FROM DB
	router.GET("/getuser/:id", func(c *gin.Context) {
		id := c.Param("id")

		var name string = ""
		db.Raw("SELECT name from users where id=?;", id).Scan(&name)
		c.JSON(200, gin.H{
			"name": name,
			"id":   id,
		})
	})

	//UPDATE USERNAME
	router.POST("/changename/:id", func(c *gin.Context) {
		id := c.Param("id")

		var requestBody struct {
			UserName string `json:"name"`
			Id       string `json:"id"`
		}
		err := c.BindJSON(&requestBody)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}
		c.JSON(200, gin.H{"data": "Successfully Updated Name",
			"id":   id,
			"name": requestBody.UserName,
		})
		fmt.Println(requestBody)
		db.Exec("Update users set name=? where id=?;", requestBody.UserName, id)
	})

	//DELETE USERNAME FROM DB
	router.DELETE("/deluser/:id", func(c *gin.Context) {
		id := c.Param("id")
		db.Exec("DELETE from users where id=?;", id)
		c.JSON(200, gin.H{"success": id + "User Deleted Successfully"})
	})
	router.Run(":8080")
}
