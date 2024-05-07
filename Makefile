install:
	go build -o /usr/local/bin/dacrane src/main.go
uninstall:
	rm /usr/local/bin/dacrane
build-parser:
	goyacc -o src/parser/parser.go src/parser/parser.go.y
