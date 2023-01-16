package infrastructure

import (
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type User struct {
	ID          primitive.ObjectID `bson:"_id" binding:"required"`
	RegNumber   string             `bson:"reg_number" binding:"required"`
	Gender      int                `bson:"gender" binding:"required"`
	Email       string             `bson:"email" binding:"required"`
	Password    string             `bson:"password" binding:"required"`
	FirstName   string             `bson:"first_name" binding:"required"`
	LastName    string             `bson:"last_name" binding:"required"`
	CreatedAt   time.Time          `bson:"created_at" binding:"required"`
	RoomNumber  string             `bson:"room_number" binding:"required"`
	BlockName   string             `bson:"block_name" binding:"required"`
	NoOfBeds    int                `bson:"no_of_beds" binding:"required"`
	IsAdmin     bool               `bson:"is_admin" binding:"required"`
	Preferences []Preferences      `bson:"preferences" binding:"required"`
}
type Preferences struct {
	BlockName    string `json:"block_name" bson:"block_name" binding:"required"`
	NumberOfBeds int    `json:"no_of_beds" bson:"no_of_beds" binding:"required"`
}
type Room struct {
	UserID            primitive.ObjectID `bson:"_id"  binding:"required"`
	InitialRoomNumber string             `json:"initial_room_number" bson:"initial_room_number"  binding:"required"`
	InitialBlockName  string             `json:"initial_block_name" bson:"initial_block_name" binding:"required"`
	InitialNoOfBeds   int                `json:"initial_no_of_beds" bson:"initial_no_of_beds" binding:"required"`
	CurrentRoomNumber string             `json:"current_room_number" bson:"current_room_number"  binding:"required"`
	CurrentBlockName  string             `json:"current_block_name" bson:"current_block_name" binding:"required"`
	CurrentNoOfBeds   int                `json:"current_no_of_beds" bson:"current_no_of_beds" binding:"required"`
}
