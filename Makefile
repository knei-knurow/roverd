all: roverd

roverd: ./roverd.go ./movement.go
	go build -o roverd ./roverd.go ./movement.go

install:
	cp ./roverd /usr/local/bin

clean:
	rm -rf ./roverd *.log *.pid
