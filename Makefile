

build:
	CGO_ENABLED=0 go build -tags most no_sqlite3 moderncsqlite -ldflags '-s -w'

run:
	${HOME}/go/bin/reflex -r '(go|html)' -s go run main.go
