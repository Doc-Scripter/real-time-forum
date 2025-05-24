run :
	clear
	cd backend && go run .
test:
	cd backend && go test -v ./...
format:
	cd backend && go fmt -w -s .
clean:
	cd backend/logs && rm -rf application.log
reset:
	cd backend && cd data
	migrate down
	migrate up
	make clean
push:
	git push github && git push origin