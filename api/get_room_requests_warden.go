package api

import (
	databasePackage "hostel_hopper/infrastructure"
	tokenMechanism "hostel_hopper/token"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

// type room struct {
// 	UserID     primitive.ObjectID `json:"user_id"`
// 	RoomNumber string             `json:"current_room_number"`
// 	BlockName  string             `json:"current_block_name"`
// 	NoOfBeds   int                `json:"current_no_of_beds"`
// }
type roomRequestsResponse struct {
	RequestID   string    `json:"request_id"`
	RequestedBy string    `json:"requested_by"`
	RequestedTo string    `json:"requested_to"`
	CreatedAt   time.Time `bson:"created_at"`
}

type getRequestsResponse struct {
	Requests []roomRequestsResponse `json:"requests"`
}

func GetAllRoomRequests(ctx *gin.Context, database *databasePackage.Database) {
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
	var response getRequestsResponse
	warderRoomsCollection := database.MongoClient.Database("hostel_hopper").Collection("warden_room_requests")
	cursor, err := warderRoomsCollection.Find(ctx, bson.M{})

	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "Unable to get requests"})
		return
	}
	defer cursor.Close(ctx)

	for cursor.Next(ctx) {
		var roomRequest databasePackage.RoomRequest
		cursor.Decode(&roomRequest)

		response.Requests = append(response.Requests, roomRequestsResponse{
			RequestID:   roomRequest.ID.Hex(),
			RequestedBy: roomRequest.RequestedBy.Hex(),
			RequestedTo: roomRequest.RequestedTo.Hex(),
			CreatedAt:   roomRequest.CreatedAt,
		})
	}

	ctx.JSON(http.StatusOK, response)

}
