all: roverd

roverd: ./roverd.go ./lidar.go
	go build ./roverd.go ./lidar.go

install:
	cp ./roverd /usr/local/bin

clean:
	rm -rf ./roverd
