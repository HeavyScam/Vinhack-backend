package api

import (
	databasePackage "hostel_hopper/infrastructure"
	tokenMechanism "hostel_hopper/token"
	"net/http"

	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type acceptRoomWardenRequest struct {
	RequestID        string `json:"request_id" binding:"required"`
	AcceptedResponse int    `json:"accepted_response" binding:"required"`
}

func AcceptRoomWarden(ctx *gin.Context, database *databasePackage.Database) {
	payload := ctx.MustGet("authorization_payload").(*tokenMechanism.Payload)
	var acceptRoomWardenRequest acceptRoomWardenRequest
	if err := ctx.ShouldBindJSON(&acceptRoomWardenRequest); err != nil {
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
	requestId, err := primitive.ObjectIDFromHex(acceptRoomWardenRequest.RequestID)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "invalid request id"})
		return
	}

	if acceptRoomWardenRequest.AcceptedResponse < 0 || acceptRoomWardenRequest.AcceptedResponse > 1 {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "invalid input"})
		return
	}
	warderRoomsCollection := database.MongoClient.Database("hostel_hopper").Collection("warden_room_requests")
	var roomRequest databasePackage.RoomRequest
	err = warderRoomsCollection.FindOneAndDelete(ctx, bson.M{"_id": requestId}).Decode(&roomRequest)

	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "unable to remove request"})
		return
	}
	roomsCollection := database.MongoClient.Database("hostel_hopper").Collection("rooms")
	if acceptRoomWardenRequest.AcceptedResponse == 1 {
		_, err = roomsCollection.DeleteOne(ctx, bson.M{"_id": roomRequest.RequestedBy})
		if err != nil {
			ctx.JSON(http.StatusBadRequest, gin.H{"error": "unable to remove room"})
			return
		}
		_, err = roomsCollection.DeleteOne(ctx, bson.M{"_id": roomRequest.RequestedTo})
		if err != nil {
			ctx.JSON(http.StatusBadRequest, gin.H{"error": "unable to remove room"})
			return
		}
	}

	ctx.JSON(http.StatusOK, gin.H{"message": "success"})

}
