package api

import (
	databasePackage "hostel_hopper/infrastructure"
	tokenMechanism "hostel_hopper/token"
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type AddPreferencesRequest struct {
	Preferences []databasePackage.Preferences `json:"preferences" bson:"preferences" binding:"required"`
}

func AddPreferences(ctx *gin.Context, database *databasePackage.Database) {
	payload := ctx.MustGet("authorization_payload").(*tokenMechanism.Payload)
	var addPreferencesRequest AddPreferencesRequest
	if err := ctx.ShouldBindJSON(&addPreferencesRequest); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	userID, err := primitive.ObjectIDFromHex(payload.UserID)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": gin.H{
			"message": "invalid user id",
			"id":      payload.UserID,
		}})
		return
	}
	var user databasePackage.User
	collection := database.MongoClient.Database("hostel_hopper").Collection("users")
	err = collection.FindOne(ctx, bson.M{"_id": userID}).Decode(&user)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "User does not exists"})
		return
	}

	for _, v := range addPreferencesRequest.Preferences {
		v.BlockName = strings.Trim(v.BlockName, " ")

		if v.BlockName == "" || v.NumberOfBeds <= 0 {
			ctx.JSON(http.StatusBadRequest, gin.H{"error": "invalid input"})
			return
		}
	}

	update := bson.M{
		"preferences": addPreferencesRequest.Preferences,
	}
	match := bson.M{"_id": userID}

	_, updateErr := collection.UpdateOne(ctx, match, bson.M{"$set": update})
	if updateErr != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "Unable to add db"})
		return
	}
	ctx.JSON(http.StatusOK, gin.H{"message": "success"})

}
