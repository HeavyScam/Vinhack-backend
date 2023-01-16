package api

import (
	databasePackage "hostel_hopper/infrastructure"
	tokenMechanism "hostel_hopper/token"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"go.mongodb.org/mongo-driver/bson"
	"golang.org/x/crypto/bcrypt"
)

type loginUserResponse struct {
	SessionID          uuid.UUID    `json:"session_id"`
	User               userResponse `json:"user"`
	AccessToken        string       `json:"access_token"`
	AccessTokenExpire  time.Time    `json:"access_token_expires_at"`
	RefreshToken       string       `json:"refresh_token"`
	RefreshTokenExpire time.Time    `json:"refresh_token_expires_at"`
}

type loginUserRequest struct {
	Email    string `json:"email" binding:"required"`
	Password string `json:"password" binding:"required"`
}

type userResponse struct {
	FirstName string `json:"first_name" binding:"required"`
	LastName  string `json:"last_name" binding:"required"`
}

func newUserResponse(user databasePackage.User) userResponse {
	return userResponse{
		FirstName: user.FirstName,
		LastName:  user.LastName,
	}
}

func Login(ctx *gin.Context, database *databasePackage.Database, tokenMaker tokenMechanism.Maker) {
	collection := database.MongoClient.Database("hostel_hopper").Collection("users")
	var input loginUserRequest
	if err := ctx.ShouldBindJSON(&input); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	var user databasePackage.User
	err := collection.FindOne(ctx, bson.M{"email": input.Email}).Decode(&user)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "User does not exists"})
		return
	}
	err = bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(input.Password))
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "Password did not matched"})
		return
	}

	loginUserResponseProfile := newUserResponse(user)
	accessToken, accessTokenPayload, err := tokenMaker.GenerateToken(user.Email, user.Gender, user.ID.Hex(), time.Hour*48)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "Access Token Creation Failed"})
		return
	}

	refreshToken, refreshTokenPayload, err := tokenMaker.GenerateToken(user.Email, user.Gender, user.ID.Hex(), time.Hour*36)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "Access Token Creation Failed"})
		return
	}

	rsp := loginUserResponse{
		SessionID:          refreshTokenPayload.ID,
		AccessToken:        accessToken,
		AccessTokenExpire:  accessTokenPayload.ExpireAt,
		RefreshToken:       refreshToken,
		RefreshTokenExpire: refreshTokenPayload.ExpireAt,
		User:               loginUserResponseProfile,
	}
	ctx.JSON(http.StatusOK, rsp)
}
