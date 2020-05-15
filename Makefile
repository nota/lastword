all: app

app:
	go build

deploy: app
	gcloud functions deploy lastword --runtime go111 --trigger-http --allow-unauthenticated --entry-point LastWord --env-vars-file .env.yaml
