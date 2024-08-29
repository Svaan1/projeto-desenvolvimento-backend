air:
	air -c ./configs/.air.toml

dbuild:
	docker build -f build/Dockerfile -t server:latest .

drun:
	docker run -d -p 8080:8080 --name backend-container server:latest

dclean:
	docker rmi server:latest || true
	docker stop backend-container || true
	docker rm backend-container || true

docs:
	godoc -http=:6060