package main

import (
	"errors"
	"net/http"
	"github.com/gin-gonic/gin"
)

type user struct {
	ID				string	`json:"id"`
	Username	string	`json:"username"`
	About			string	`json:"about"`
}

var users = []user {
	{ ID: "1", Username: "MarcDev", About: "Fullstack Developer & Product Designer" },
	{ ID: "2", Username: "Mezo", About: "Would You Bot Developer" },
}

func getUsers(c *gin.Context) {
	c.IndentedJSON(http.StatusOK, users)
}

func findUser(id string) (*user, error) {
	for i, b := range users {
		if b.ID == id {
			return &users[i], nil
		}
	}

	return nil, errors.New("user not found")
}

func getUser(c *gin.Context) {
	id := c.Param("id")
	user, err := findUser(id)

	if err != nil {
		c.IndentedJSON(http.StatusNotFound, gin.H{"message": "User not found.", "code": 404})
		return
	}

	c.IndentedJSON(http.StatusOK, user)
}

func changeUsername(c *gin.Context) {
	id, idOk := c.GetQuery("id")
	new, newOk := c.GetQuery("new")

	if !idOk {
		c.IndentedJSON(http.StatusBadRequest, gin.H{"message": "Missing the id query parameter.", "code": 400})
		return
	}

	if !newOk {
		c.IndentedJSON(http.StatusBadRequest, gin.H{"message": "Missing the new query parameter.", "code": 400})
		return
	}

	user, err := findUser(id)

	if err != nil {
		c.IndentedJSON(http.StatusNotFound, gin.H{"message": "User not found.", "code": 404})
		return
	}

	user.Username = new
	c.IndentedJSON(http.StatusOK, user)
}

func changeAbout(c *gin.Context) {
	id, idOk := c.GetQuery("id")
	new := c.Param("new")

	if !idOk {
		c.IndentedJSON(http.StatusBadRequest, gin.H{"message": "Missing the id query parameter", "code": 400})
		return
	}

	user, err := findUser(id)

	if err != nil {
		c.IndentedJSON(http.StatusNotFound, gin.H{"message": "User not found.", "code": 404})
		return
	}

	user.About = new
	c.IndentedJSON(http.StatusOK, user)
}

func createUser(c *gin.Context) {
	var newUser user

	if err := c.BindJSON(&newUser); err != nil {
		return
	}

	users = append(users, newUser)
	c.IndentedJSON(http.StatusCreated, newUser)
}

func findIndexOfUser(id string) (int) {
	for k, v := range users {
			if id == v.ID {
					return k
			}
	}
	return -1
}

func deleteUser(c *gin.Context) {
	id := c.Param("id")

	var userIndex = findIndexOfUser(id)

	if userIndex == -1 {
		c.IndentedJSON(http.StatusNotFound, gin.H{"message": "User not found.", "code": 404})
		return
	}

	users = append(users[:userIndex], users[userIndex+1:]...)
	c.IndentedJSON(http.StatusOK, gin.H{"message": "User deleted successfully.", "code": 200})
}

func main() {
	router := gin.Default()

	router.GET("/users", getUsers)
	router.GET("/user/:id", getUser)

	router.POST("/user", createUser)

	router.PATCH("/username", changeUsername)
	router.PATCH("/about/:new", changeAbout)

	router.DELETE("/user/:id", deleteUser)

	router.Run("localhost:8080")
}