.PHONY: build clean

build:
	go build -o _bin/lambda-sample ./lambda
	go build -o _bin/lambda-sample ./dynamodb

clean:
	rm _bin/*
