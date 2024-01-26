# FIAP - TechChallenge - Payment Service

## Features

- Create Payments
- Search Payments By ID
- Get All Payments
- Receive Callbacks from payment providers

## How To Run Locally

First of all we need the DataBase. To set it up you have 2 options:

Option 1: $```docker-compose -f deployments/db-docker-compose.yml up -d```

Option 2: $```make run-db```

Both are going to have the same result.

Then you can run the application:

### VSCode - Debug
The launch.json file is already configured for debuging. Just hit F5 and be happy.

### Running directly from go

Option 1: $```go run cmd/client/main.go```

Option 2: $```make run-app```

## Manually testing the API

On directory ```/api``` there's a collection that can be imported on Insomnia or similar so you can test manually the application's API.

## Running the unit tests

Simply run ```make run-tests``` and let the magic happens. At the end it will automatically open an html with the coverage % for every package.