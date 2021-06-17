all: roverd

roverd: ./roverd.go
	go build -o roverd ./roverd.go

install:
	cp ./roverd /usr/local/bin

clean:
	rm -rf ./roverd *.log *.pid
