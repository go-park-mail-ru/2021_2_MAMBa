.PHONY: test
test:
	go test -v -coverpkg=./... -coverprofile=profile.cov ./...
	cat profile.cov | grep -v ".pb.go:|mock|easyjson.go:|handlers.go|app" > profile1.cov
	go tool cover -func profile1.cov