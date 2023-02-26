VERSION=0.3.3
BIN_NAME=adcl

all: build windows mac

bin_dir:
	mkdir -p ./bin/

_versions_equal:
	./check-versions.sh

build: bin_dir
	go build -o ./bin/${BIN_NAME} ./cmd/cli/

windows: bin_dir
	GOOS=windows GOARCH=386 go build -o ./bin/${BIN_NAME}.exe ./cmd/cli/

mac: bin_dir
	GOOS=darwin GOARCH=amd64 go build -o ./bin/${BIN_NAME}_darwin ./cmd/cli/

_build-alpine: bin_dir
	CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo -o ./bin/${BIN_NAME} ./cmd/cli/

docker-build: bin_dir
	docker build --target builder -t ${BIN_NAME}:${VERSION}-build . -f Dockerfile.build
	docker container create --name extract ${BIN_NAME}:${VERSION}-build
	docker container cp extract:/app/bin/${BIN_NAME} ./bin/
	docker container cp extract:/app/bin/${BIN_NAME}.exe ./bin/
	docker container cp extract:/app/bin/${BIN_NAME}_darwin ./bin/
	docker container rm -f extract

docker-cli:
	docker build --no-cache -t ${BIN_NAME} .
	echo "Try: docker run adcl"

releases: _versions_equal all
	tar -czvf ./bin/${BIN_NAME}v${VERSION}_linux.tar.gz ./bin/${BIN_NAME}
	tar -czvf ./bin/${BIN_NAME}v${VERSION}_win.tar.gz ./bin/${BIN_NAME}.exe
	tar -czvf ./bin/${BIN_NAME}v${VERSION}_darwin.tar.gz ./bin/${BIN_NAME}_darwin

install:
	go build -o ${GOPATH}/bin/${BIN_NAME} ./cmd/cli/

uninstall:
	rm -f ${GOPATH}/bin/${BIN_NAME}

clean_bin:
	rm -f ./bin/${BIN_NAME}
	rm -f ./bin/${BIN_NAME}_darwin
	rm -f ./bin/${BIN_NAME}.exe

clean_releases:
	rm -f ./bin/${BIN_NAME}v${VERSION}_linux.tar.gz
	rm -f ./bin/${BIN_NAME}v${VERSION}_win.tar.gz
	rm -f ./bin/${BIN_NAME}v${VERSION}_darwin.tar.gz

clean: clean_bin clean_releases
	rmdir ./bin/
