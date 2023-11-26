package main

import (
	"context"
	"database/sql"
	"fmt"
	"math/rand"
	"time"

	"github.com/gin-gonic/gin" //gin
	"github.com/jackc/pgx/v5"  // pgx PostgreSQL driver
	"github.com/jackc/pgx/v5/pgtype"
	"golangTest/golangTest"
	"net/http"
)

// User structure
type User struct {
	ID            int    `json:"id"`
	Name          string `json:"name"`
	PhoneNumber   string `json:"phone_number"`
	OTP           string `json:"otp"`
	OTPExpiration string `json:"otp_expiration"`
}

var ctx = context.Background()
var db golangTest.DBTX

func main() {
	var err error
	// Connecting to the database
	db, err := pgx.Connect(ctx, "user=faraz dbname=test sslmode=verify-full")
	if err != nil {

	}
	defer db.Close(ctx)

	router := gin.Default()

	// Creating routes using gin
	router.POST("/api/users", createUser)
	router.POST("/api/users/generateotp", generateOTP)
	router.POST("/api/users/verifyotp", verifyOTP)

	router.Run(":8080")
}

func createUser(c *gin.Context) {
	var newUser User
	// Parsing the JSON body into the newUser struct
	if err := c.BindJSON(&newUser); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Checking if the phone number already exists in the database
	var exists bool
	queries := golangTest.New(db)
	exists, err := queries.CheckPhoneExistence(ctx, newUser.PhoneNumber)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Database error"})
		return
	}

	if exists {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Phone number already in use"})
		return
	}

	// Inserting the new user into the database
	err = queries.CreateUser(ctx, golangTest.CreateUserParams{
		Name:        newUser.Name,
		PhoneNumber: newUser.PhoneNumber,
	})
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Unable to create user"})
		return
	}

	c.JSON(http.StatusCreated, gin.H{"message": "User created successfully"})
}

func generateOTP(c *gin.Context) {
	var User User
	// Parsing the JSON body to get the user's phone number
	if err := c.BindJSON(&User); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Checking if the phone number exists in the database
	var exists bool
	queries := golangTest.New(db)
	exists, err := queries.CheckPhoneExistence(ctx, User.PhoneNumber)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Database error"})
		return
	}

	if !exists {
		c.JSON(http.StatusNotFound, gin.H{"error": "User not found"})
		return
	}

	// Generating a random 4-digit OTP
	otp := fmt.Sprintf("%04d", rand.Intn(10000))

	// Setting OTP expiration time (1 minute from now)
	expirationTime := time.Now().Add(time.Minute)

	// Updating user's OTP in the database
	// _, err = queries.UpdateUserOTP("UPDATE users SET otp = $1, otp_expiration_time = $2 WHERE phone_number = $3", otp, expirationTime, User.PhoneNumber)
	err = queries.UpdateUserOTP(ctx, golangTest.UpdateUserOTPParams{
		PhoneNumber:       User.PhoneNumber,
		Otp:               pgtype.Text{String: otp},
		OtpExpirationTime: pgtype.Timestamp{Time: expirationTime},
	})
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Unable to generate OTP"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "OTP generated successfully"})
}

func verifyOTP(c *gin.Context) {
	var User User

	// Parsing the JSON body to get the phone number and OTP
	if err := c.BindJSON(&User); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Retrieving the OTP and its expiration time from the database
	var storedOTP, otpExpiration string
	queries := golangTest.New(db)
	otpRow, err := queries.GetOTP(ctx, User.PhoneNumber)
	if err != nil {
		if err == sql.ErrNoRows {
			// User with the given phone number not found
			c.JSON(http.StatusNotFound, gin.H{"error": "User not found"})
		} else {
			// Database error
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Database error"})
		}
		return
	}
	otpRow.Otp.Scan(&storedOTP)
	otpRow.OtpExpirationTime.Scan(&otpExpiration)

	// Checking if the OTP is correct
	if storedOTP != User.OTP {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Incorrect OTP"})
		return
	}

	// Checking if the OTP has expired
	expirationTime, _ := time.Parse(time.RFC3339, otpExpiration)
	if time.Now().After(expirationTime) {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "OTP has expired"})
		return
	}

	// OTP is correct and not expired
	c.JSON(http.StatusOK, gin.H{"message": "OTP verified successfully"})
}
