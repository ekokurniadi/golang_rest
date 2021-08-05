package main

import (
	"bytes"
	"database/sql"
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
	_ "github.com/go-sql-driver/mysql"
)

func main() {
	db, err := sql.Open("mysql", "root:@tcp(127.0.0.1:3306)/belajar_golang")
	err = db.Ping()
	if err != nil {
		panic("Gagal terkoneksi ke database")
	}
	defer db.Close()
	router := gin.Default()

	type Album struct {
		Id        int    `json:"id"`
		Foto      string `json:"foto"`
		Caption   string `json:"caption"`
		CreatedAt string `json:"createdAt"`
	}
	type User struct {
		Id       int     `json:"id"`
		Nama     string  `json:"nama"`
		Email    string  `json:"email"`
		Password string  `json:"password"`
		Album    []Album `json:"albums"`
	}

	// mendapatkan data by id user
	router.GET("/user/:id", func(c *gin.Context) {
		var (
			user   User
			result gin.H
		)
		id := c.Param("id")
		row := db.QueryRow("select id,nama,email,password from user where id=?", id)
		err = row.Scan(&user.Id, &user.Nama, &user.Email, &user.Password)

		if err != nil {
			result = gin.H{
				"message": "Data tidak ditemukan",
				"count":   0,
			}
		} else {
			result = gin.H{
				"message": "sucess",
				"values":  user,
			}
		}
		c.JSON(http.StatusOK, result)
	})

	router.GET("/user/album/:id", func(c *gin.Context) {
		var (
			albums []Album
			album  Album
			user   User
			result gin.H
		)
		id := c.Param("id")
		row := db.QueryRow("select id,nama,email,password from user where id=?", id)
		err = row.Scan(&user.Id, &user.Nama, &user.Email, &user.Password)
		if err != nil {
			fmt.Print(err.Error())
		}
		rows, errs := db.Query("select id,foto,caption,created_at from user_album where id_user=?", id)
		if errs != nil {
			fmt.Print(errs.Error())
		}
		for rows.Next() {
			errs = rows.Scan(&album.Id, &album.Foto, &album.Caption, &album.CreatedAt)
			albums = append(albums, album)
			if errs != nil {
				fmt.Print(errs.Error())
			}
		}
		final := User{Id: user.Id, Nama: user.Nama, Email: user.Email, Password: user.Password, Album: albums}
		if err != nil {
			result = gin.H{
				"status":  "Error",
				"message": "Data tidak ditemukan",
			}
		} else {
			result = gin.H{
				"status": "Success",
				"values": final,
			}
		}
		defer rows.Close()
		c.JSON(http.StatusOK, result)

	})

	// mendapatkan seluruh data user
	router.GET("/user", func(c *gin.Context) {
		var (
			user  User
			users []User
		)
		rows, err := db.Query("SELECT id,nama,email,password from user")
		if err != nil {
			fmt.Print(err.Error())
		}
		for rows.Next() {
			err = rows.Scan(&user.Id, &user.Nama, &user.Email, &user.Password)
			users = append(users, user)
			if err != nil {
				fmt.Print(err.Error())
			}

		}
		defer rows.Close()
		c.JSON(http.StatusOK, gin.H{
			"message": "Success",
			"values":  users,
		})
	})

	// create user baru
	router.POST("/user", func(c *gin.Context) {
		var buffer bytes.Buffer

		nama := c.PostForm("nama")
		email := c.PostForm("email")
		password := c.PostForm("password")

		smt, err := db.Prepare("INSERT INTO user (nama,email,password) values (?,?,?)")
		if err != nil {
			fmt.Print(err.Error())
		}

		_, err = smt.Exec(nama, email, password)
		if err != nil {
			fmt.Print(err.Error())
		}

		buffer.WriteString(nama)
		buffer.WriteString(" ")
		buffer.WriteString(email)
		defer smt.Close()

		data := buffer.String()
		c.JSON(http.StatusOK, gin.H{
			"message": fmt.Sprintf("Berhasil menambahkan data %s", data),
		})
	})

	// update data user
	router.PUT("/user", func(c *gin.Context) {
		var buffer bytes.Buffer
		id := c.PostForm("id")
		nama := c.PostForm("nama")
		email := c.PostForm("email")
		password := c.PostForm("password")

		smt, err := db.Prepare("UPDATE user set nama=?,email=?,password=? where id=?")

		if err != nil {
			fmt.Print(err.Error())
		}

		_, err = smt.Exec(nama, email, password, id)

		if err != nil {
			fmt.Print(err.Error())
		}

		buffer.WriteString(nama)
		buffer.WriteString(" ")
		buffer.WriteString(email)
		defer smt.Close()

		data := buffer.String()
		c.JSON(http.StatusOK, gin.H{
			"message": fmt.Sprintf("Berhasil merubah data %s", data),
		})
	})

	router.DELETE("/user", func(c *gin.Context) {
		id := c.PostForm("id")
		smt, err := db.Prepare("DELETE FROM user where id=?")
		if err != nil {
			fmt.Print(err.Error())
		}
		_, err = smt.Exec(id)

		if err != nil {
			fmt.Print(err.Error())
		}

		c.JSON(http.StatusOK, gin.H{
			"message": fmt.Sprintf("Berhasil menghapus data %s", id),
		})
	})

	router.Run(":8080")

}
