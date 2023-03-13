# Secret app milestone 1

## config

- Golang

## setup

Copy `.env-example` to `.env`

## usage

Run
```
go run main.go
```

Now you can add secrets with
```
curl --request POST \
  --url http://localhost:8080/ \
  --header 'Content-Type: application/json' \
  --data '{
	"plain_text": "secret"
}'
```

And retrieve the precendent with
```
curl --request GET \
  --url http://localhost:8080/ \
  --header 'Content-Type: application/json' \
  --data '{
	"id": "5ebe2294ecd0e0f08eab7690d2a6ee69"
}'
```
Change the id based on what you input.