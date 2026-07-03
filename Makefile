PROTO_DIR=internal/controller/grpc/v1/inventory/pb

proto:
	cd $(PROTO_DIR) && PATH="$$PATH:$$(go env GOPATH)/bin" buf dep update && PATH="$$PATH:$$(go env GOPATH)/bin" buf generate

docker_build:
	docker build --no-cache -t inventory:latest .

docker_up:
	docker compose up -d

docker_down:
	docker compose down --volumes
	docker container prune -f
	docker image prune -f
	docker volume prune -f

docker_restart: docker_down docker_build docker_up
docker_start: docker_down docker_up

REST_HOST=localhost:8081
GRPC_HOST=localhost:50051
TRACE_ID=trace-$(shell date +%s)

# REST tests
rest-create:
	curl -s -X POST http://$(REST_HOST)/api/v1/inventories \
		-H "Content-Type: application/json" \
		-H "X-Request-ID: $(TRACE_ID)" \
		-d '{"product_id":"ABC-001","name":"Mouse","description":"Wireless mouse","price":299.99,"quantity":50}' | jq .

rest-create-2:
	curl -s -X POST http://$(REST_HOST)/api/v1/inventories \
		-H "Content-Type: application/json" \
		-H "X-Request-ID: $(TRACE_ID)" \
		-d '{"product_id":"ABC-002","name":"Keyboard","description":"Mechanical keyboard","price":599.99,"quantity":30}' | jq .

rest-get-all:
	curl -s http://$(REST_HOST)/api/v1/inventories \
		-H "X-Request-ID: $(TRACE_ID)" | jq .

rest-get-by-id:
	curl -s http://$(REST_HOST)/api/v1/inventories/ABC-001 \
		-H "X-Request-ID: $(TRACE_ID)" | jq .

rest-test: rest-create rest-create-2 rest-get-all rest-get-by-id

# gRPC tests
grpc-list:
	grpcurl -plaintext $(GRPC_HOST) list

grpc-create:
	grpcurl -plaintext \
		-H "x-request-id: $(TRACE_ID)" \
		-d '{"product_id":"XYZ-001","name":"Monitor","description":"4K display","price":9999.99,"quantity":10}' \
		$(GRPC_HOST) inventory.InventoryService/CreateProduct

grpc-create-2:
	grpcurl -plaintext \
		-H "x-request-id: $(TRACE_ID)" \
		-d '{"product_id":"XYZ-002","name":"Headset","description":"Noise cancelling","price":2499.99,"quantity":20}' \
		$(GRPC_HOST) inventory.InventoryService/CreateProduct

grpc-get-all:
	grpcurl -plaintext \
		-H "x-request-id: $(TRACE_ID)" \
		$(GRPC_HOST) inventory.InventoryService/GetProducts

grpc-get-by-id:
	grpcurl -plaintext \
		-H "x-request-id: $(TRACE_ID)" \
		-d '{"product_id":"XYZ-001"}' \
		$(GRPC_HOST) inventory.InventoryService/GetProductByID

grpc-test: grpc-create grpc-create-2 grpc-get-all grpc-get-by-id

db-audit:
	docker exec postgres psql -U postgres -d inventory \
		-c "SELECT product_id, created_by_event, created_by_trace_id, updated_by_event, updated_by_trace_id FROM products;"

test: rest-test grpc-test db-audit

k8s-up:
	kubectl apply -f k8s/namespace.yaml
	kubectl apply -f k8s/configmap.yaml
	kubectl apply -f k8s/secret.yaml
	kubectl apply -f k8s/postgres.yaml
	kubectl apply -f k8s/redis.yaml
	kubectl apply -f k8s/api.yaml
	kubectl apply -f k8s/grpc.yaml

k8s-down:
	kubectl delete -f k8s/

k8s-status:
	kubectl get all -n inventory