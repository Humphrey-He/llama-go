#!/usr/bin/env python3
"""
Ollama Backend Adapter - 替代 llama-cpp-python
使用 Ollama 提供 OpenAI 兼容 API
"""
import logging
import requests
import json

logging.basicConfig(level=logging.INFO)
logger = logging.getLogger(__name__)

class OllamaServer:
    def __init__(self, model_name="llama3.1:8b", host="http://localhost:11434"):
        self.model_name = model_name
        self.host = host

    def check_health(self):
        """检查 Ollama 服务状态"""
        try:
            resp = requests.get(f"{self.host}/api/tags", timeout=5)
            if resp.status_code == 200:
                models = resp.json().get("models", [])
                logger.info(f"Ollama 可用，已加载模型: {[m['name'] for m in models]}")
                return True
        except Exception as e:
            logger.error(f"Ollama 不可用: {e}")
            return False

    def chat_completion(self, messages, stream=False):
        """OpenAI 兼容的聊天接口"""
        url = f"{self.host}/api/chat"
        payload = {
            "model": self.model_name,
            "messages": messages,
            "stream": stream
        }

        try:
            resp = requests.post(url, json=payload, stream=stream, timeout=60)
            resp.raise_for_status()

            if stream:
                for line in resp.iter_lines():
                    if line:
                        yield json.loads(line)
            else:
                return resp.json()
        except Exception as e:
            logger.error(f"请求失败: {e}")
            raise

def main():
    logger.info("=" * 60)
    logger.info("启动 Ollama 后端服务")
    logger.info("=" * 60)

    server = OllamaServer()

    if not server.check_health():
        logger.error("\n请先安装并启动 Ollama:")
        logger.error("1. 下载: https://ollama.com/download")
        logger.error("2. 安装后自动启动")
        logger.error("3. 运行: ollama pull llama3.1:8b")
        return

    logger.info("\nOllama 后端已就绪")
    logger.info("Go 服务可以连接到: http://localhost:11434")

if __name__ == "__main__":
    main()
