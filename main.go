package main

import (
	"database/sql"
	"fmt"
	"log"
	"net/http"
	"strconv"

	_ "github.com/denisenkom/go-mssqldb"
	"github.com/gin-gonic/gin"
)

type Person struct {
	Id       int    `json:"id" form:"id"`
	Username string `json:"username" form:"username"`
	Email    string `json:"email" form:"email"`
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
		rows.Close()
		c.JSON(http.StatusOK, gin.H{
			"data": listPerson,
		})
	})
	r.GET("/:id", func(c *gin.Context) {
		id := c.Param("id")
		var person Person
		row, err := db.Query("select * from Persons where id=@id", sql.Named("id", id))
		// row, err := db.QueryContext(c, `select * from Persons where id = @id;`, sql.Named("id", id))
		if err != nil {
			log.Fatal(err)
		}
		for row.Next() {
			row.Scan(&person.Id, &person.Username, &person.Email, &person.Age)
		}
		defer row.Close()
		if person.Id > 0 {
			c.JSON(http.StatusOK, gin.H{
				"status": "success",
				"data":   person,
			})
			return
		} else {
			c.JSON(http.StatusNotFound, gin.H{
				"status": "failed",
				"msg":    "Person not found",
			})
			return
		}

	})

	r.DELETE("/:id", func(c *gin.Context) {
		id := c.Param("id")
		stmt, err := db.Prepare("DELETE FROM Persons WHERE id=@id")
		if err != nil {
			log.Fatal(err)
		}
		rs, err := stmt.Exec(sql.Named("id", id))
		if err != nil {
			log.Fatal(err)
		}
		row, err := rs.RowsAffected()
		if err != nil {
			log.Fatal(err)
		}
		defer stmt.Close()
		rows := row
		if row == 0 {
			c.JSON(http.StatusNotFound, gin.H{
				"status": "failed",
				"msg":    "Person not found or deleted failed",
			})

		} else {
			c.JSON(http.StatusOK, gin.H{
				"status": "success",
				"count":  rows,
			})
		}
	})

	r.POST("/", func(c *gin.Context) {
		var person Person
		if err := c.BindJSON(&person); err != nil {
			return
		}
		var pr Person
		r, err := db.Query("select top 1 * from Persons where email=@email", sql.Named("email", person.Email))
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{
				"status": "failed",
				"error":  err,
			})
			return
		}
		for r.Next() {
			r.Scan(&pr.Id, &pr.Username, &pr.Email, &pr.Age)

		}
		r.Close()
		if pr.Email == person.Email {
			c.JSON(http.StatusBadRequest, gin.H{
				"status":  "failed",
				"message": "Email already exists",
			})
			return
		} else {
			row, err := db.Query("INSERT INTO Persons(username, email, age) VALUES (@username, @email, @age);select ID = convert(bigint, SCOPE_IDENTITY())", sql.Named("username", person.Username), sql.Named("email", person.Email), sql.Named("age", person.Age))
			if err != nil {
				c.JSON(http.StatusBadRequest, gin.H{
					"status": "failed",
					"error":  err,
				})
				return
			}
			var lastInsertId int64
			for row.Next() {
				row.Scan(&lastInsertId)

			}
			person.Id = int(lastInsertId)
			row.Close()
			c.JSON(http.StatusCreated, gin.H{
				"success": true,
				"person":  person,
			})
		}

	})

	r.PUT("/:id", func(c *gin.Context) {
		id := c.Param("id")
		var person Person
		var pr Person
		if err := c.BindJSON(&pr); err != nil {
			return
		}
		row, err := db.Query("select * from Persons where id=@id", sql.Named("id", id))
		if err != nil {
			log.Fatal(err)
		}
		for row.Next() {
			row.Scan(&person.Id, &person.Username, &person.Email, &person.Age)
		}
		row.Close()
		i, _ := strconv.Atoi(id)
		if person.Id == i {
			r, err := db.Query("UPDATE Persons SET username = @username, email = @email, age = @age WHERE id = @id;select ID = convert(bigint, SCOPE_IDENTITY())", sql.Named("username", pr.Username), sql.Named("email", pr.Email), sql.Named("age", pr.Age), sql.Named("id", id))
			if err != nil {
				log.Fatal(err)
			}
			pr.Id = person.Id
			r.Close()
			c.JSON(http.StatusOK, gin.H{
				"success": true,
				"person":  pr,
			})
		} else {
			c.JSON(http.StatusNotFound, gin.H{
				"status":  "failed",
				"message": "Person not found",
			})
		}

	})
	fmt.Println("Server running http://localhost:4000")
	r.Run(":4000")
}
