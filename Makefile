APP = mydocker
GO  = go

.PHONY: build

build:
	$(GO) build -ldflags "$(LDFLAGS)" -o $(APP) 

clean:
	$(GO) clean


run: build
	./$(APP)

test:
	$(GO) test -v ./...

