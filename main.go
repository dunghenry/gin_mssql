package main

import (
	"database/sql"
	"fmt"
	"net/http"

	_ "github.com/denisenkom/go-mssqldb"
	"github.com/gin-gonic/gin"
)

type Person struct {
	Id       int    `json:"id" form:"id"`
	Username string `json:"username" form:"username"`
	Email    string    `json:"email" form:"email"`
	Age      int    `json:"age" form:"age"`
}

func main() {
	gin.SetMode(gin.ReleaseMode)
	r := gin.New()
	r.Use(gin.Logger())
	r.Use(gin.Recovery())
	r.SetTrustedProxies(nil)
	db, _ := sql.Open("sqlserver", "sqlserver://sa:123456789@DESKTOP-929CP1R:49693?database=GINMSSQL&connection+timeout=30")
	r.GET("/", func(c *gin.Context) {
		var listPerson []Person
		rows, _ := db.Query("select id, username, email, age from Persons")
		for rows.Next() {
			var person Person
			rows.Scan(&person.Id, &person.Username, &person.Email, &person.Age)
			listPerson = append(listPerson, person)
		}
		c.JSON(http.StatusOK, gin.H{
			"data": listPerson,
		})
	})

	fmt.Println("Server running http://localhost:4000")
	r.Run(":4000")
}
