build:
	go build ./...

win-build:
	env GOOS=windows GOARCH=amd64 CGO_ENABLED=1 CC=x86_64-w64-mingw32-gcc make build

win-run-tests-1:
	cd tests && go test -o ../tests.test -c -covermode=atomic -coverprofile=c.out -coverpkg=../...
	wine ./tests.test -test.v=true -test.coverprofile=c.out

win-run-tests:
	env GOOS=windows GOARCH=amd64 CGO_ENABLED=1 CC=x86_64-w64-mingw32-gcc make win-run-tests-1
	make run-tests-2

run-tests:
	cd tests && go test -o ../tests.test -c -covermode=atomic -coverprofile=c.out -coverpkg=../...
	./tests.test -test.v=true -test.coverprofile=c.out
	make run-tests-2

run-tests-2:
	sed -i -e '\:^github.com/sudachen/go-fp/:d' c.out
	cp c.out gocov.txt
	sed -i -e 's:github.com/sudachen/go-dl/::g' c.out

run-cover:
	go tool cover -html=gocov.txt

run-cover-tests: run-tests run-cover

