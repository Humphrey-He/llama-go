# 启动 Go 服务
run-go:
	go run main.go

# 启动 Python 服务
run-py:
	cd python && pip install -r requirements.txt && python kv_server.py

# 容器启动
up:
	docker-compose up -d