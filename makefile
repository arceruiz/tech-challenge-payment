run-db:
	docker-compose -f deployments/db-docker-compose.yml up -d
run-tests:
	go test $$(go list ./... | grep -v /data/) -coverprofile=cover.out.tmp && cat ./cover.out.tmp | grep -v "mock.go" > ./cover.out && go tool cover -html=cover.out 
test:
	go test $$(go list ./... | grep -v /data/) -coverprofile=cover.out.tmp
test-build-bake:
	docker build -t mauricio1998/payment-service . -f build/Dockerfile