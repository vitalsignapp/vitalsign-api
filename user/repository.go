package user

import (
	"context"
	"fmt"

	"cloud.google.com/go/firestore"
)

// User User
type User struct {
	ID       string `json:"id"`
	Password string `json:"password"`
}

// .collection("userData")
// .where("userId", "==", "999")

// NewChangePassword NewChangePassword
func NewChangePassword(fs *firestore.Client) func(context.Context, string, string) error {
	return func(ctx context.Context, userID string, password string) error {
		fmt.Println(userID)
		_, err := fs.Collection("userData").Doc(userID).Get(ctx)
		if err != nil {
			return err
		}
		_, err = fs.Collection("userData").Doc(userID).Set(ctx, map[string]interface{}{"password": password}, firestore.MergeAll)
		if err != nil {
			return err
		}
		return nil
	}
}
