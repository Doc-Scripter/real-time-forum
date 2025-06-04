run :
	clear
	cd backend && go run .
test:
	cd backend && go test -v ./...
format:
	cd backend && gofmt -w -s .
clean:
	cd backend/logs && rm -rf application.log
reset:
	cd backend && cd data
	migrate down
	migrate up
	make clean
push:
	git push github && git push gitea
merge:
	git merge master
git:
	@echo "If you get an error , remove all remote -> git remote remove <branch name>"
	@if git remote | grep -q "origin"; then \
		echo "Removing existing 'origin' remote..."; \
		git remote remove origin; \
	fi 
	@if git remote | grep -q "github"; then \
		echo "Removing existing 'origin' remote..."; \
		git remote remove github; \
	fi
	git remote add github https://github.com/Doc-Scripter/real-time-forum.git
	git remote add gitea https://learn.zone01kisumu.ke/git/cliffootieno/real-time-forum
	@echo "Remotes have been updated."
pull:
	git fetch github
	git pull --all