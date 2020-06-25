.PHONY: build clean

build:
	go build -o _bin/lambda-sample ./lambda

clean:
	rm _bin/*
