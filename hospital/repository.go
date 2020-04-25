package hospital

import (
	"context"

	"cloud.google.com/go/firestore"
)

// Hospital Hospital
type Hospital struct {
	DomainPrefix     string           `json:"domainPrefix"`
	Name             string           `json:"name"`
	VitalSignsConfig []HospitalConfig `json:"vitalsignsConfig"`
}

// HospitalConfig HospitalConfig
type HospitalConfig struct {
	Status bool   `json:"status"`
	Sym    string `json:"sym"`
}

// NewUpdateConfigPatient NewUpdateConfigPatient
func NewUpdateConfigPatient(fs *firestore.Client) func(context.Context, string, []HospitalConfig) error {
	return func(ctx context.Context, hospitalID string, hospitalConfigs []HospitalConfig) error {
		_, err := fs.Collection("hospital").Doc(hospitalID).Get(ctx)
		if err != nil {
			return err
		}

		data := []map[string]interface{}{}

		for _, element := range hospitalConfigs {
			data = append(data, map[string]interface{}{
				"status": element.Status,
				"sym":    element.Sym,
			})
		}

		_, err = fs.Collection("hospital").Doc(hospitalID).Set(ctx, map[string]interface{}{"vitalSignsConfig": data}, firestore.MergeAll)
		if err != nil {
			return err
		}
		return nil
	}
}
