install:
	go install ./cmd/tfplugingen-kubernetes

test:
	go test ./...

.PHONY: install test
