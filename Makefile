.PHONY: run

W_DIR  = /home/setupuser/nmlt
GOPATH = /usr/local/go
TARGET = bin/nmlt
SOURCE = main.go

run: ${TARGET}
	@cd ${W_DIR} && ${TARGET} -luadir bin/listup

build: ${SOURCE}
	@cd ${W_DIR} && ${GOPATH}/bin/go build -o ${TARGET} ${SOURCE}

test:
	@cd ${W_DIR} && ${GOPATH}/bin/go test -v

${TARGET}: ${SOURCE}
	@make build
