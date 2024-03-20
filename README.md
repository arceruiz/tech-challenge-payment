# FIAP - TechChallenge - Payment Service

# Description

This service is responsible to receive the payment request from a SQS queue, send it to the payment provider (not implemented) and receive it's callback with the payment state.
After receiving the state from payment provider, it will update the payment status and send a new message in a SQS queue about this state.
We have a diagram about a flow of this service here: [Create Flow ](./docs/diagrams/flow-diagram.png), [Cancelled Flow](./docs/diagrams/cancelled-flow-diagram.png) 

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

## Manually testing the API

On directory ```/api``` there's a collection that can be imported on Insomnia or similar so you can test manually the application's API.

## Running the unit tests

Simply run ```make run-tests``` and let the magic happens. At the end it will automatically open an html with the coverage % for every package.
We also have the most recently applied unit tests file in this [folder](./docs/unit-tests-results/payment-unit.png) too. And there is a html file about the last unit tests [execution](./docs/unit-tests-results/coverage.html).

## Test + Build + Bake Image

Simply run ```make test-build-bake``` and let the magic happens. The docker file will run the unit-tests, build the application and bake the docker image for the application.

## Infrastructure

This application runs in a k8s cluster. The manifests about the configuration of this application are in this [repository](https://github.com/mauriciodm1998/payment-service-gitops).