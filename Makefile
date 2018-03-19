test:
	go test -v -race

coverprofile:
	go test -coverprofile=coverage.out

cover: coverprofile
	go tool cover -html=coverage.out

deps:
	dep ensure