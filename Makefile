build:
	go build ./...

win-build-cross-all:
	go build ./...
	cd tests && go test -o ../tests.exe -c -covermode=atomic -coverprofile=c.out -coverpkg=../...

win-build:
	env GOOS=windows GOARCH=amd64 CGO_ENABLED=1 CC=x86_64-w64-mingw32-gcc make win-build-cross-all

win-run-tests:
	env GOOS=windows GOARCH=amd64 CGO_ENABLED=1 CC=x86_64-w64-mingw32-gcc make win-build-cross-all
	wine ./tests.exe -test.v=true -test.coverprofile=c.out
	make run-tests-2

run-tests-1:
	cd tests && go test -o ../tests.test -c -covermode=atomic -coverprofile=c.out -coverpkg=../...
	./tests.test -test.v=true -test.coverprofile=c.out

run-tests-2:
	sed -i -e '\:^github.com/sudachen/go-foo/:d' c.out
	sed -i -e 's:github.com/sudachen/go-dl/::g' c.out
	awk '/\.go/{print "github.com/sudachen/go-dl/"$$0}/^mode/{print $$0}' < c.out > gocov.txt

run-tests: run-tests-1 run-tests-2

run-cover:
	go tool cover -html=gocov.txt

run-cover-tests: run-tests run-cover

run-cover-all:
	make run-tests-1
	cp c.out c1.out
	make win-run-tests
	mv c.out c2.out
	cp c1.out c3.out
	tail -n +2 c2.out >> c3.out
	head -n 1 c3.out > c.out
	tail -n +2 c3.out | sort >> c.out
	make run-tests-2
	make run-cover
