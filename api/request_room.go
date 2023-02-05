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

type requestRoomRequest struct {
	RequestedBy string `json:"requested_by" binding:"required"`
	RequestedTo string `json:"requested_to" binding:"required"`
}

func RequestRoom(ctx *gin.Context, database *databasePackage.Database) {
	payload := ctx.MustGet("authorization_payload").(*tokenMechanism.Payload)
	var requestRoomRequest requestRoomRequest
	if err := ctx.ShouldBindJSON(&requestRoomRequest); err != nil {
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
	requestRoomRequest.RequestedBy = strings.Trim(requestRoomRequest.RequestedBy, " ")
	requestRoomRequest.RequestedTo = strings.Trim(requestRoomRequest.RequestedTo, " ")

	if requestRoomRequest.RequestedBy == "" || requestRoomRequest.RequestedTo == "" {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "invalid input"})
		return
	}

	// match := bson.M{"_id": userID}

	// _, updateErr := collection.UpdateOne(ctx, match, )
	// if updateErr != nil {
	// 	ctx.JSON(http.StatusInternalServerError, gin.H{"error": "Unable to add db"})
	// 	return
	// }
	ctx.JSON(http.StatusOK, gin.H{"message": "success"})

}
