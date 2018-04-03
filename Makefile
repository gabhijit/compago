.PHONY: all

node: node/*.go
	echo "Making node"
	go build -v -o node-agent ./node/


controller:
	echo "Making Controller"

all: node controller

clean:
	$(RM) -rf node-agent
