install:
	go build -o /usr/local/bin/dacrane src/main.go
uninstall:
	rm /usr/local/bin/dacrane
build-parser:
	goyacc -o src/core/evaluator/parser.go src/core/evaluator/parser.go.y
