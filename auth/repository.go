package auth

import (
	"context"

	"cloud.google.com/go/firestore"
	"google.golang.org/api/iterator"
)

func CheckAuthen(fs *firestore.Client) func(context.Context, string, string, string) *LoginResponse {
	return func(ctx context.Context, email, password, hospitalKey string) *LoginResponse {
		iter := fs.Collection("userData").
			Where("email", "==", email).
			Where("password", "==", password).
			Where("hospitalKey", "==", hospitalKey).
			Documents(ctx)
		defer iter.Stop()

		p := LoginResponse{}
		for {
			doc, err := iter.Next()
			if err == iterator.Done {
				break
			}

			if err != nil {
				break
			}

			p.ID = doc.Ref.ID
			err = doc.DataTo(&p)
			if err != nil {
				return nil
			}
		}
		return &p
	}
}
