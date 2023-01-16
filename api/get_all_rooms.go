package api

import (
	databasePackage "hostel_hopper/infrastructure"
	tokenMechanism "hostel_hopper/token"
	"net/http"

	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type roomResponse struct {
	UserID     primitive.ObjectID `json:"user_id" bson:"_id"`
	RoomNumber string             `json:"current_room_number" bson:"current_room_number" `
	BlockName  string             `json:"current_block_name" bson:"current_block_name"`
	NoOfBeds   int                `json:"current_no_of_beds" bson:"current_no_of_beds"`
}

type getRoomsResponse struct {
	Rooms []roomResponse `json:"rooms"`
}

func GetAllRooms(ctx *gin.Context, database *databasePackage.Database) {
	payload := ctx.MustGet("authorization_payload").(*tokenMechanism.Payload)
	roomsCollection := database.MongoClient.Database("hostel_hopper").Collection("rooms")
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
	var response getRoomsResponse

	cursor, err := roomsCollection.Find(ctx, bson.M{})

	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "Unable to get rooms"})
		return
	}
	defer cursor.Close(ctx)

	for cursor.Next(ctx) {
		var roomResponse roomResponse
		cursor.Decode(&roomResponse)
		response.Rooms = append(response.Rooms, roomResponse)
	}

	ctx.JSON(http.StatusOK, response)

}
