all: roverd

roverd: ./roverd.go ./env.go
	go build -o roverd ./roverd.go ./env.go

install:
	cp ./roverd /usr/local/bin

clean:
	rm -rf ./roverd *.log *.pid
