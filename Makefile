install:
	go build -o /usr/local/bin/dacrane src/main.go
uninstall:
	rm /usr/local/bin/dacrane
build-parser:
	goyacc -o src/code/parser/parser.go src/code/parser/parser.go.y
