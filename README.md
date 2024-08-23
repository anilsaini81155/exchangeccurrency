#steps to run the application

1)go mod tidy
2)go build ./cmd/server
3)go run ./cmd/server/main.go

Username  HashedPassword                                               UserType      
abc@g.com $2a$14$xNVmXN6dvywGoHaS/tYqg.OtJ/PvzTNzHcygO0akY/7VmEb.Re6dW exchange_admin
xyz@g.com $2a$14$MLgc2bSwHZvZIXMOCqmb/e9PPNWVDLLSYI6hbR5B3ItvFqNibSwIS exchange_user

A)APIS:

1)GET API: 

This api requires exchange_pair as the query params.

CURL :

curl --location 'http://localhost:8000/rate?exchange_pair=USD%2FEUR'

curl --location 'http://localhost:8000/login' \
--header 'Content-Type: application/json' \
--data-raw '{"username":"abc@g.com",
"password":"$2a$14$xNVmXN6dvywGoHaS/tYqg.OtJ/PvzTNzHcygO0akY/7VmEb.Re6dW"
}'

Response:
{"exchange_pair":"USD/EUR","exchange_rate":151,"message":"Exchange rates and pair fetched successfully."}

2)POST Login API:

This api requires username and password from the config file to get the token , only exchange admin can get the token.Token validity is of 24 hrs.

CURL :
curl --location 'http://localhost:8000/login' \
--header 'Content-Type: application/json' \
--data-raw '{"username":"abc@g.com",
"password":"$2a$14$xNVmXN6dvywGoHaS/tYqg.OtJ/PvzTNzHcygO0akY/7VmEb.Re6dW"
}'

Response:

{"token":"eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJ1c2VybmFtZSI6ImFiY0BnLmNvbSIsImV4cCI6MTcyNDQ4MzIzNn0.mctqdD6u11K9HoSJSictf9KtWZmb1MEeVFouf1HfKgE"}

3)POST Rate API:

This api requires exchange_pair,exchange_rate,username in the input request to store the data , without jwt token this api will not work.
Once successfully updated the websocket client will get the msg.

curl --location 'http://localhost:8000/rate' \
--header 'Content-Type: application/json' \
--header 'Authorization: Bearer eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJ1c2VybmFtZSI6ImFiY0BnLmNvbSIsImV4cCI6MTcyNDQ4MzIzNn0.mctqdD6u11K9HoSJSictf9KtWZmb1MEeVFouf1HfKgE' \
--data '{"username":"testuser",
"exchange_pair":"USD/EUR",
"exchange_rate": 151.00
}'

B)WebSocket:

1) For client connection use below url to connect to the server.
ws://localhost:8000/ws

##passsing msg should be in json only.
