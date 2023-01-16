package api

import (
	"context"
	"net/http"
	"strings"

	databasePackage "hostel_hopper/infrastructure"
	tokenMechanism "hostel_hopper/token"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"golang.org/x/crypto/bcrypt"
)

type registrationResponse struct {
	SessionID          uuid.UUID                `json:"session_id"`
	User               registrationUserResponse `json:"user"`
	AccessToken        string                   `json:"access_token"`
	AccessTokenExpire  time.Time                `json:"access_token_expires_at"`
	RefreshToken       string                   `json:"refresh_token"`
	RefreshTokenExpire time.Time                `json:"refresh_token_expires_at"`
}

type registrationUserResponse struct {
	FirstName string `json:"first_name" binding:"required"`
	LastName  string `json:"last_name" binding:"required"`
}

type registrationUserRequest struct {
	Email      string `json:"email" binding:"required"`
	Password   string `json:"password" binding:"required"`
	RegNumber  string `json:"reg_number" binding:"required"`
	FirstName  string `json:"first_name" binding:"required"`
	LastName   string `json:"last_name" binding:"required"`
	RoomNumber string `json:"room_number" binding:"required"`
	BlockName  string `json:"block_name" binding:"required"`
	NoOfBeds   int    `json:"no_of_beds" binding:"required"`
	Gender     int    `json:"gender" binding:"required"`
}

func registerationNewUserResponse(user databasePackage.User) registrationUserResponse {
	return registrationUserResponse{
		FirstName: user.FirstName,
		LastName:  user.LastName,
	}
}

func Register(ctx *gin.Context, database *databasePackage.Database, tokenMaker tokenMechanism.Maker) {
	userCollection := database.MongoClient.Database("hostel_hopper").Collection("users")
	roomsCollection := database.MongoClient.Database("hostel_hopper").Collection("rooms")
	var input registrationUserRequest
	if err := ctx.ShouldBindJSON(&input); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	var tempUser databasePackage.User
	err := userCollection.FindOne(ctx, bson.M{"email": input.Email}).Decode(&tempUser)
	if err == nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "User Already Exists"})
		return
	}

	input.Email = strings.Trim(input.Email, " ")
	if input.Email == "" {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "Email empty"})
		return
	}

	input.Password = strings.Trim(input.Password, " ")
	if input.Password == "" {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "Password Empty"})
		return
	}

	input.FirstName = strings.Trim(input.FirstName, " ")
	if input.FirstName == "" {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "First Name empty"})
		return
	}

	input.LastName = strings.Trim(input.LastName, " ")
	if input.LastName == "" {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "Last Name empty"})
		return
	}
	input.BlockName = strings.Trim(input.BlockName, " ")
	if input.LastName == "" {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "Blockname empty"})
		return
	}
	input.RoomNumber = strings.Trim(input.RoomNumber, " ")
	if input.RoomNumber == "" {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "RoomNumber empty"})
		return
	}

	if input.NoOfBeds <= 0 {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "no fo beds invalid"})
		return
	}
	if input.Gender > 1 {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "asa kuch nh hota"})
		return
	}

	hash, err := bcrypt.GenerateFromPassword([]byte(input.Password), bcrypt.MinCost)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "Password error"})
		return
	}
	var created_at, _ = time.Parse(time.RFC3339, time.Now().Format(time.RFC3339))
	userID := primitive.NewObjectID()
	user := databasePackage.User{
		ID:          userID,
		RegNumber:   input.RegNumber,
		Email:       input.Email,
		Gender:      input.Gender,
		Password:    string(hash),
		FirstName:   input.FirstName,
		LastName:    input.LastName,
		RoomNumber:  input.RoomNumber,
		BlockName:   input.BlockName,
		NoOfBeds:    input.NoOfBeds,
		CreatedAt:   created_at,
		Preferences: []databasePackage.Preferences{},
		IsAdmin:     false,
	}
	_, errorInsertion := userCollection.InsertOne(context.TODO(), user)
	if errorInsertion != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "Unable to add user to DB"})
		return
	}

	_, errorInsertion = roomsCollection.InsertOne(context.TODO(), databasePackage.Room{
		UserID:            user.ID,
		InitialRoomNumber: input.RoomNumber,
		InitialBlockName:  input.BlockName,
		InitialNoOfBeds:   input.NoOfBeds,
		CurrentRoomNumber: input.RoomNumber,
		CurrentBlockName:  input.BlockName,
		CurrentNoOfBeds:   input.NoOfBeds,
	})
	if errorInsertion != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "Unable to add room to DB"})
		return
	}

	registerUserResponseProfile := registerationNewUserResponse(user)
	accessToken, accessTokenPayload, err := tokenMaker.GenerateToken(user.Email, user.Gender, userID.Hex(), time.Hour*48)

	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "Access Token Creation Failed"})
		return
	}
	refreshToken, refreshTokenPayload, err := tokenMaker.GenerateToken(user.Email, user.Gender, userID.Hex(), time.Hour*36)

	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "Refresh Token Creation Failed"})
		return
	}

	rsp := registrationResponse{
		SessionID:          refreshTokenPayload.ID,
		AccessToken:        accessToken,
		AccessTokenExpire:  accessTokenPayload.ExpireAt,
		RefreshToken:       refreshToken,
		RefreshTokenExpire: refreshTokenPayload.ExpireAt,
		User:               registerUserResponseProfile,
	}
	ctx.JSON(http.StatusOK, rsp)

}
