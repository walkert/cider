.PHONY: build, clean, install test
install:
	@cd cmd/cider && CGO_ENABLED=0 go install
	@echo cider has been installed

build:
	@cd cmd/cider && CGO_ENABLED=0 go build
	@echo cider has been built

test:
	@go test -cover

clean:
	@rm -f cmd/cider/cider
