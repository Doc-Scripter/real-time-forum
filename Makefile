run :
	clear
	cd backend && go run .
test:
	cd backend && go test -v ./...
format:
	cd backend && go fmt -w -s .
clean:
	rm -rf application.log
reset:
	cd backend && cd data
	migrate down
	migrate up
	make clean