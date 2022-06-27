package util

import (
	"fmt"
	"math/rand"
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

//GenerateID Generate 12 byte UUID. It will return 24 hexadimal characters
func GenerateID() string {
	return primitive.NewObjectID().Hex()
}

//RandStringBytes Generate 6 byte UID Alphanumeric
func RandStringBytes(length int) string {
	rand.Seed(time.Now().UnixNano())
	b := make([]byte, length)
	rand.Read(b)
	return fmt.Sprintf("%x", b)[:length]
}
