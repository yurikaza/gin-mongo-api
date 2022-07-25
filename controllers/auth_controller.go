package controllers

import (
	"context"
	"gin-mongo-api/configs"
	"gin-mongo-api/models"
	"gin-mongo-api/utils"
	"gin-mongo-api/responses"
	"net/http"
	"time"
	"log"

	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
    "github.com/gin-contrib/sessions"
    "golang.org/x/crypto/bcrypt"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
)

func HashPassword(password string) (string, error) {
    bytes, err := bcrypt.GenerateFromPassword([]byte(password), 14)
    return string(bytes), err
}

func CheckPasswordHash(password, hash string) bool {
    err := bcrypt.CompareHashAndPassword([]byte(hash), []byte(password))
    return err == nil
}

var authCollection *mongo.Collection = configs.GetCollection(configs.DB, "users")
var authValidate = validator.New()

func Register() gin.HandlerFunc {
	return func(c *gin.Context) {
		ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
		var user models.User
		defer cancel()

		//validate the request body
		if err := c.BindJSON(&user); err != nil {
			c.JSON(http.StatusBadRequest, responses.UserResponse{Status: http.StatusBadRequest, Message: "error", Data: map[string]interface{}{"data": err.Error()}})
			return
		}

		//use the validator library to validate required fields
		if validationErr := authValidate.Struct(&user); validationErr != nil {
			c.JSON(http.StatusBadRequest, responses.UserResponse{Status: http.StatusBadRequest, Message: "error", Data: map[string]interface{}{"data": validationErr.Error()}})
			return
		}

        hashedPassword, _ := HashPassword(user.Password)

		newUser := models.User{
			Id:       primitive.NewObjectID(),
			Name:     user.Name,
			Email:    user.Email,
			Password: hashedPassword,
		}

		result, err := authCollection.InsertOne(ctx, newUser)
		if err != nil {
			c.JSON(http.StatusInternalServerError, responses.UserResponse{Status: http.StatusInternalServerError, Message: "error", Data: map[string]interface{}{"data": err.Error()}})
			return
		}

		c.JSON(http.StatusCreated, responses.UserResponse{Status: http.StatusCreated, Message: "success", Data: map[string]interface{}{"data": result}})
	}
}

func Login() gin.HandlerFunc {
	return func(c *gin.Context) {
		ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
		var user models.User
		var dbUser models.User

		//validate the request body
		if err := c.BindJSON(&user); err != nil {
			c.JSON(http.StatusBadRequest, responses.UserResponse{Status: http.StatusBadRequest, Message: "error", Data: map[string]interface{}{"data": err.Error()}})
			return
		}

        err := authCollection.FindOne(ctx, bson.M{"email":user.Email}).Decode(&dbUser)
		defer cancel()

        if err != nil {
            c.JSON(http.StatusInternalServerError, gin.H{"error": "login or passowrd is incorrect"})
            return
        }		

		userPass := []byte(user.Password)
        dbPass := []byte(dbUser.Password)
        passErr := bcrypt.CompareHashAndPassword(dbPass, userPass)

        if passErr != nil{
	       log.Println(passErr)
           return
        }
        
        mongoId := dbUser.Id
        stringObjectID := mongoId.Hex()
        token, err := utils.CreateJwtToken(c, stringObjectID)

        if err != nil{
	       log.Println(err)
           return
        }

        session := sessions.Default(c)
        v := session.Get("jwt")

		c.JSON(http.StatusCreated, responses.UserResponse{Status: http.StatusCreated, Message: "success", Data: map[string]interface{}{"data": dbUser, "token": token}})
	}
}

