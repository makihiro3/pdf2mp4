DEPS=$(shell find -type f -name '*.go')
backend: $(DEPS)
	CGO_ENABLED=0 GOOS=linux GOARCH=amd64 GOAMD64=v3 go build -o $@ -trimpath -buildmode pie .
