.PHONY: build run clean

build:
	docker build -t codepulse_codepulse .

run:
	docker run --name codepulse-container -d -p 8080:8080 codepulse_codepulse 

clean:
	docker rm -f codepulse-container
	docker rmi -f codepulse_codepulse