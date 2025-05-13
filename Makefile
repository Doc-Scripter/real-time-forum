run :
	clear
	cd backend && go run .
test:
	cd backend && go test -v ./...
format:
	cd backend && go fmt -w -s .