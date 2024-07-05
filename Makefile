.PHONY: build run clean

build:
	docker build -t codepulse-img .

run:
	docker run --name codepulse-container -d -p 8080:8080 codepulse-img 

clean:
	docker rm -f codepulse-container
	docker rmi -f codepulse-img