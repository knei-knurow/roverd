all: roverd

roverd: ./roverd.go ./lidar.go ./head.go
	go build ./roverd.go ./lidar.go ./head.go

install:
	cp ./roverd /usr/local/bin

clean:
	rm -rf ./roverd
