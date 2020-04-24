package ward

import (
	"context"

	"cloud.google.com/go/firestore"
	"google.golang.org/api/iterator"
)

type Ward struct {
	ID          string `json:"id"`
	Name        string `json:"name"`
	CreatedDate int    `json:"createdDate"`
}
type PatientRoomDate struct {
	Date      string `json:"date"`
	Microtime int    `json:"microtime"`
	Week      int    `json:"week"`
}
type PatientRoom struct {
	AddTime     int             `json:"addTime"`
	Date        PatientRoomDate `json:"date"`
	HospitalKey string          `json:"hospitalKey"`
	Name        string          `json:"name"`
}

func NewRepository(fs *firestore.Client) func(context.Context, string) []Ward {
	return func(ctx context.Context, hospitalKey string) []Ward {
		iter := fs.Collection("patientRoom").Where("hospitalKey", "==", hospitalKey).Documents(ctx)
		defer iter.Stop()
		wards := []Ward{}
		for {
			doc, err := iter.Next()
			if err == iterator.Done {
				break
			}
			if err != nil {
				continue
			}

			p := PatientRoom{}
			err = doc.DataTo(&p)
			if err != nil {
				continue
			}

			wards = append(wards, Ward{
				ID:          doc.Ref.ID,
				Name:        p.Name,
				CreatedDate: p.Date.Microtime,
			})
		}

		return wards
	}
}

func AddNewRepository(fs *firestore.Client) func(context.Context, RoomRequest) error {
	return func(ctx context.Context, r RoomRequest) error {

		_, _, err := fs.Collection("patientRoom").Add(ctx, map[string]interface{}{
			"name":        r.Name,
			"hospitalKey": r.HospitalKey,
			"addTime":     r.AddTime,
			"date": map[string]interface{}{
				"date":      r.Date.Date,
				"microtime": r.Date.Microtime,
				"week":      r.Date.Week,
			},
		})
		if err != nil {
			return err
		}
		return nil
	}
}
