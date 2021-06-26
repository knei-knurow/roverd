all: roverd

roverd: ./roverd.go ./env.go
	go build -o roverd ./roverd.go ./env.go

install:
	cp ./roverd /usr/local/bin

setup:
	cp ./roverd.service /etc/systemd/system
	mkdir -p /etc/systemd/system/roverd.service.d
	cp ./env.conf /etc/systemd/system/roverd.service.d

uninstall:
	rm /usr/local/bin/roverd

clean:
	rm ./roverd
