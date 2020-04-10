package patient

import (
	"context"
	"encoding/json"
	"net/http"

	"cloud.google.com/go/firestore"
	"github.com/gorilla/mux"
	"google.golang.org/api/iterator"
)

func NewScheduler(fs *firestore.Client) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		vars := mux.Vars(r)
		ctx := context.Background()
		iter := fs.Collection("patientData").Where("HN", "==", vars["patientID"]).Documents(ctx)
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

		json.NewEncoder(w).Encode(&data)
	}
}
