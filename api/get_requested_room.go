package api

import (
	databasePackage "hostel_hopper/infrastructure"
	tokenMechanism "hostel_hopper/token"
	"net/http"

	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type getRequestedRoomResponse struct {
	RequestID string               `json:"request_id"`
	Room      databasePackage.Room `json:"room"`
}

func GetRequestedRoom(ctx *gin.Context, database *databasePackage.Database) {
	payload := ctx.MustGet("authorization_payload").(*tokenMechanism.Payload)
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
	var roomRequest databasePackage.RoomRequest
	roomsRequestsCollection := database.MongoClient.Database("hostel_hopper").Collection("room_requests")
	err = roomsRequestsCollection.FindOne(ctx, bson.M{"requested_to": userID}).Decode(&roomRequest)

	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "Unable to get requests"})
		return
	}

	var room databasePackage.Room
	roomsCollection := database.MongoClient.Database("hostel_hopper").Collection("rooms")
	err = roomsCollection.FindOne(ctx, bson.M{"_id": roomRequest.RequestedBy}).Decode(&room)

	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "Unable to get room"})
		return
	}
	ctx.JSON(http.StatusOK, getRequestedRoomResponse{
		RequestID: roomRequest.ID.Hex(),
		Room:      room,
	})

}
