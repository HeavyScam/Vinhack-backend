package api

import (
	databasePackage "hostel_hopper/infrastructure"
	tokenMechanism "hostel_hopper/token"
	"net/http"

	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type acceptRoomRequest struct {
	RequestID        string `json:"request_id" binding:"required"`
	AcceptedResponse int    `json:"accepted_response" binding:"required"`
}

func AcceptRoom(ctx *gin.Context, database *databasePackage.Database) {
	payload := ctx.MustGet("authorization_payload").(*tokenMechanism.Payload)
	var acceptRoomRequest acceptRoomRequest
	if err := ctx.ShouldBindJSON(&acceptRoomRequest); err != nil {
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
	requestId, err := primitive.ObjectIDFromHex(acceptRoomRequest.RequestID)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "invalid request id"})
		return
	}

	if acceptRoomRequest.AcceptedResponse < 0 || acceptRoomRequest.AcceptedResponse > 1 {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "invalid input"})
		return
	}
	requestCollection := database.MongoClient.Database("hostel_hopper").Collection("room_requests")
	var roomRequest databasePackage.RoomRequest
	err = requestCollection.FindOneAndDelete(ctx, bson.M{"_id": requestId}).Decode(&roomRequest)

	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "unable to remove request"})
		return
	}
	warderRoomsCollection := database.MongoClient.Database("hostel_hopper").Collection("warden_room_requests")
	if acceptRoomRequest.AcceptedResponse == 1 {
		_, err = warderRoomsCollection.InsertOne(ctx, roomRequest)
		if err != nil {
			ctx.JSON(http.StatusBadRequest, gin.H{"error": "unable to add to warden db"})
			return
		}
	}

	ctx.JSON(http.StatusOK, gin.H{"message": "success"})

}
