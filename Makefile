.PHONY: build-backend
build-backend:
	go build -o ./bin/circadian

.PHONY: build-frontend
build-frontend:
	cd ui && \
	npm install && \
	npm run build

.PHONY: run-frontend
run-frontend:
	cd ui && \
	npm run dev

# .PHONY: run-backend
# run-backend: build-backend
#	./bin/circadian

.PHONY: run-backend
run-backend:
	go run main.go ./sample/config.yml

.PHONY: build-all
build-all: build-backend build-frontend
	
