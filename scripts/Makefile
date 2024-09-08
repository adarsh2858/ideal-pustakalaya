CMD_DIR = ../cmd
BIN_DIR= ../bin

default:
	@for f in $(shell ls -d ${CMD_DIR}/*/); do \
		echo $${f} | sed -E 's/\.\.\/cmd\/([a-zA-Z]+)\/$$/building \1.../'; \
		pushd $${f} > /dev/null; \
		go build -o ../${BIN_DIR}/; \
		popd > /dev/null; \
	done

install:
	@pushd ../cmd/cli > /dev/null && \
	echo "installing {{.repoName}} cli at ${GOPATH}/bin" && \
	go build -o ${GOPATH}/bin/{{.repoName}} && \
	popd > /dev/null
