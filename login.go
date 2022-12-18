package main

// Thanks to Depado for the code example
// https://github.com/Depado/gin-auth-example/

import (
	"math/rand"
	"net/http"
	"strings"
	"time"
	"github.com/gin-gonic/gin"
)

type student_login_data struct {
	Roll_No string `json:"roll_no"`
	Password string `json:"password"`
	Token string `json:"token"`
}

type admin_login_data struct {
	Username string `json:"username"`
	Password string `json:"password"`
	Token string `json:"token"`
}

// This map will store the username/token key value pairs
var users = make(map[string]string)

// AdminAuthRequired is a simple middleware to check the session
func AdminAuthRequired(c *gin.Context) {
	var login_data admin_login_data

	if c.BindJSON(&login_data) != nil {
		c.AbortWithStatus(http.StatusUnauthorized)
        return
    }

	if v, ok := users[login_data.Username]; ok && v != login_data.Token {
        c.AbortWithStatus(http.StatusUnauthorized)
        return
	}
	// Continue down the chain to handler etc
	c.Next()
}

// AuthRequired is a simple middleware to check the session
func AuthRequired(c *gin.Context) {
	var login_data student_login_data

	if c.BindJSON(&login_data) != nil {
        return
    }

	if v, ok := users[login_data.Roll_No]; ok && v != login_data.Token {
        c.AbortWithStatus(http.StatusUnauthorized)
        return
	}
	// Continue down the chain to handler etc
	c.Next()
}

// Secure enough for a prototype
func GenerateUserToken() string {
	var letterRunes = []rune("abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ1234567890!@#$%^&")
	b := make([]rune, 255)

    rand.Seed(time.Now().UnixNano())
    for i := range b {
        b[i] = letterRunes[rand.Intn(len(letterRunes))]
    }
    return string(b)
}

// login is a handler that parses a form and checks for specific data

func AdminLoginFunc(c *gin.Context) {
	var login_data admin_login_data
	
	if c.BindJSON(&login_data) != nil {
        return
    }

	// Validate form input
	if strings.Trim(login_data.Username, " ") == "" || strings.Trim(login_data.Password, " ") == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Parameters can't be empty"})
		return
	}

	// Check for roll_no and password match, usually from a database
	if !CheckUsernameAndPassword(login_data.Username, login_data.Password) {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Authentication failure"})
		return
	}

	// Generate token and make it a key pair with roll_no
	token := GenerateUserToken()
	users[login_data.Username] = token
	c.JSON(http.StatusOK, gin.H{"token": token})
}

func LoginFunc(c *gin.Context) {
	var login_data student_login_data

	if c.BindJSON(&login_data) != nil {
        return
    }

	// Validate form input
	if strings.Trim(login_data.Roll_No, " ") == "" || strings.Trim(login_data.Password, " ") == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Parameters can't be empty"})
		return
	}

	// Check for roll_no and password match, usually from a database
	if !CheckRollNoAndPassword(login_data.Roll_No, login_data.Password) {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Authentication failure"})
		return
	}

	// Generate token and make it a key pair with roll_no
	token := GenerateUserToken()
	users[login_data.Roll_No] = token
	c.JSON(http.StatusOK, gin.H{"token": token})
}

func LogoutFunc(c *gin.Context) {
	var login_data student_login_data

	if c.BindJSON(&login_data) != nil {
        return
    }

	if v, ok := users[login_data.Roll_No]; ok && v == login_data.Token {
        // If the username/token pair is found in the users map,
        // remove this username from the users map
        // and respond with an HTTP success status
        delete(users, login_data.Roll_No)
        c.JSON(http.StatusOK, gin.H{"message": "Successfully logged out"})
    } else {
        // If the username/token pair is not found in the users map,
        // respond with an HTTP error
        c.AbortWithStatus(http.StatusUnauthorized)
    }
}

func AdminLogoutFunc(c *gin.Context) {
	var login_data admin_login_data

	if c.BindJSON(&login_data) != nil {
		c.AbortWithStatus(http.StatusUnauthorized)
        return
    }

	if v, ok := users[login_data.Username]; ok && v == login_data.Token {
        // If the username/token pair is found in the users map,
        // remove this username from the users map
        // and respond with an HTTP success status
        delete(users, login_data.Username)
        c.JSON(http.StatusOK, gin.H{"message": "Successfully logged out"})
    } else {
        // If the username/token pair is not found in the users map,
        // respond with an HTTP error
        c.AbortWithStatus(http.StatusUnauthorized)
    }
}

func SigninUser(c *gin.Context) {
	var student_data students_database

	if c.BindJSON(&student_data) != nil {
        return
    }
    
    AddDataToDatabase(student_data)
}
