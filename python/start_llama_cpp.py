#!/usr/bin/env python3
import logging
from pathlib import Path
from llama_cpp import Llama
from llama_cpp.server.app import create_app
import uvicorn

logging.basicConfig(level=logging.INFO)
logger = logging.getLogger(__name__)

def main():
    model_path = "internal/model/Dolphin3.0-Llama3.1-8B-Q4_K_S.gguf"

    if not Path(model_path).exists():
        logger.error(f"模型文件不存在: {model_path}")
        return

    logger.info(f"加载模型: {model_path}")
    logger.info(f"大小: {Path(model_path).stat().st_size / 1024 / 1024:.2f} MB")

    # 创建 OpenAI 兼容服务器
    app = create_app(
        model=model_path,
        n_gpu_layers=-1,  # 使用所有 GPU 层
        n_ctx=2048,
        verbose=True
    )

    logger.info("启动服务器: http://localhost:8000")
    uvicorn.run(app, host="0.0.0.0", port=8000)

if __name__ == "__main__":
    main()
