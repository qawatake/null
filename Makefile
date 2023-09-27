test:
	go mod tidy -modfile=go_test.mod
	go test ./... -modfile go_test.mod -shuffle=on -race

test.cover:
	go mod tidy -modfile=go_test.mod
	go test -modfile=go_test.mod -race -shuffle=on -coverprofile=coverage.txt -covermode=atomic ./...

mod.clean:
	rm -f go.mod go.sum
	cat go.mod.bk > go.mod
