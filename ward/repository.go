package ward

import (
	"context"

	"cloud.google.com/go/firestore"
	"google.golang.org/api/iterator"
)

func NewRepository(fs *firestore.Client) func(context.Context, string) []map[string]interface{} {
	return func(ctx context.Context, hospitalKey string) []map[string]interface{} {
		iter := fs.Collection("patientRoom").Where("hospitalKey", "==", hospitalKey).Documents(ctx)
		defer iter.Stop()

		data := []map[string]interface{}{}
		for {
			doc, err := iter.Next()
			if err == iterator.Done {
				break
			}
			if err != nil {
				// TODO: Handle error.
			}
			data = append(data, doc.Data())
		}

		return data
	}
}
