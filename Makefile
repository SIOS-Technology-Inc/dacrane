install:
	cd src && go build -o /usr/local/bin/dacrane main.go
uninstall:
	rm /usr/local/bin/dacrane
