# API Key Rotator Makefile
# 轻量级和企业级构建

.PHONY: build-lightweight build-enterprise build-all test-lightweight test-enterprise

# 构建默认版本 (轻量级 - SQLite + 内存缓存)
build:
	docker build -t api-key-rotator .

# 构建企业级版本 (MySQL + Redis)
build-enterprise:
	docker build -f Dockerfile.enterprise -t api-key-rotator:enterprise .

# 构建所有版本
build-all: build build-enterprise

# 轻量级构建 (别名)
build-lightweight: build

# 测试轻量级版本
test-lightweight:
	@echo "Testing lightweight build..."
	docker run --rm -p 8001:8000 \
		-e ADMIN_PASSWORD="test123" \
		-e JWT_SECRET="test_jwt_secret" \
		-e GLOBAL_PROXY_KEYS="test_key" \
		-v $(PWD)/test-data:/app/data \
		api-key-rotator:lightweight &
	@echo "Lightweight version started on port 8001"
	@sleep 5
	@echo "Stopping test container..."
	@docker stop $$(docker ps -q --filter "ancestor=api-key-rotator:lightweight")

# 测试企业级版本 (需要MySQL和Redis服务)
test-enterprise:
	@echo "Testing enterprise build..."
	docker-compose -f docker-compose.test.yml up -d mysql redis
	@sleep 10
	docker run --rm -p 8002:8000 \
		--network apikeyrotator_default \
		-e ADMIN_PASSWORD="test123" \
		-e JWT_SECRET="test_jwt_secret" \
		-e GLOBAL_PROXY_KEYS="test_key" \
		-e DB_HOST="mysql" \
		-e DB_USER="testuser" \
		-e DB_PASSWORD="testpass" \
		-e DB_NAME="testdb" \
		-e REDIS_HOST="redis" \
		api-key-rotator:enterprise &
	@echo "Enterprise version started on port 8002"
	@sleep 5
	@echo "Stopping test containers..."
	@docker stop $$(docker ps -q --filter "ancestor=api-key-rotator:enterprise")
	@docker-compose -f docker-compose.test.yml down

# 发布到Docker Hub (需要先登录)
publish-lightweight: build-lightweight
	docker tag api-key-rotator:lightweight yourusername/api-key-rotator:latest
	docker tag api-key-rotator:lightweight yourusername/api-key-rotator:lightweight
	docker push yourusername/api-key-rotator:latest
	docker push yourusername/api-key-rotator:lightweight

publish-enterprise: build-enterprise
	docker tag api-key-rotator:enterprise yourusername/api-key-rotator:enterprise
	docker push yourusername/api-key-rotator:enterprise

publish-all: publish-lightweight publish-enterprise

# 清理构建缓存
clean:
	docker system prune -f
	docker volume prune -f

# 显示构建信息
info:
	@echo "API Key Rotator Build Information"
	@echo "================================"
	@echo "Lightweight: SQLite + Memory Cache"
	@echo "Enterprise:  MySQL + Redis (includes all features)"
	@echo ""
	@echo "Build targets:"
	@echo "  build-lightweight - Build lightweight version"
	@echo "  build-enterprise  - Build enterprise version"
	@echo "  build-all         - Build both versions"
	@echo "  test-lightweight  - Test lightweight build"
	@echo "  test-enterprise   - Test enterprise build"
	@echo "  publish-lightweight - Push lightweight to Docker Hub"
	@echo "  publish-enterprise   - Push enterprise to Docker Hub"
	@echo "  clean             - Clean Docker build cache"