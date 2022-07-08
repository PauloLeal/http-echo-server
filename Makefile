BUILD_DIR=.build
COVER_PROFILE_FILE="${BUILD_DIR}/go-cover.tmp"

.PHONY: target clean mk-build-dir update-deps build-deps clean-test test cover-html badge

target: test

clean:
	rm -rf $(TARGET_FILE) $(BUILD_DIR)

############## build tasks

mk-build-dir:
	@mkdir -p ${BUILD_DIR}

update-deps:
	@go get -u -d -v ./...
	
build-deps:
	@go get -d -v ./...

############## test tasks

clean-test:
	@go fmt ./...
	@go clean -testcache

test: clean-test
	go test -p 1 ./...

cover-html: mk-build-dir clean-test
	@go test -p 1 -coverprofile=${COVER_PROFILE_FILE} ./... ; echo
	@go tool cover -html=${COVER_PROFILE_FILE}
	$(MAKE) badge

badge:
	@go install github.com/jpoles1/gopherbadger@latest
	gopherbadger -md="README.md" -png=false 1>&2 2> /dev/null
	@if [ -f coverage.out ]; then rm coverage.out ; fi; 
