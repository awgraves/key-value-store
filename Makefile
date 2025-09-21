.PHONY: dev-kvs-up dev-kvs-down

dev-kvs-up:
	cd kv_service && docker-compose -f docker-compose.dev.yaml up --build -d

dev-kvs-down:
	cd kv_service && docker-compose -f docker-compose.dev.yaml down
