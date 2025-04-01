APP = mydocker
GO  = go

.PHONY: build

build:
	@$(GO) build -gcflags "all=-N -l" -ldflags "$(LDFLAGS)" -o $(APP)

clean:
	@$(GO) clean

test:
	@$(GO) test -v ./...

