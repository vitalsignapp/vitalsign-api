GOOGLE_CLOUD_PROJECT=vitalsigns-426ee

test:
	go test ./... -v

run:
	GOOGLE_CLOUD_PROJECT=${GOOGLE_CLOUD_PROJECT} go run main.go