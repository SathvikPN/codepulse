.phony: build

build:
	docker build -t leetcode-proxy-server .

run:
	docker run -d -p 8080:8080 --name leetcode-proxy-server leetcode-proxy-server

