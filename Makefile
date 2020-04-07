GOOGLE_CLOUD_PROJECT=vitalsign-2bc48

test:
	go test ./... -v

run:
	GOOGLE_CLOUD_PROJECT=${GOOGLE_CLOUD_PROJECT} go run app.go