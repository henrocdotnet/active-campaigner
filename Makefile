
GO=go
GOBUILD=$(GO) build


ifeq ($(OS), Windows_NT)
	EXE_EXT='.exe'
else
	EXE_EXT=''
endif

all: build-cli

build-cli:
	$(GOBUILD) -o bin/cli$(EXE_EXT) ./campaigner-cli/

test:
	echo "Not implemented yet"
