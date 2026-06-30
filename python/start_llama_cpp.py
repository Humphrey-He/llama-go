#!/usr/bin/env python3
import sys
import subprocess
import logging
from pathlib import Path

logging.basicConfig(level=logging.INFO)
logger = logging.getLogger(__name__)

def install_deps():
    """安装依赖"""
    try:
        import llama_cpp
    except ImportError:
        logger.info("安装 llama-cpp-python...")
        subprocess.check_call([sys.executable, "-m", "pip", "install", "llama-cpp-python[server]", "-q"])

def main():
    install_deps()

    from llama_cpp.server.app import create_app
    import uvicorn

    model_path = "../models/Dolphin3.0-Llama3.1-8B-Q4_K_S.gguf"

    if not Path(model_path).exists():
        logger.error(f"模型文件不存在: {model_path}")
        return

    logger.info(f"加载模型: {model_path}")
    logger.info(f"大小: {Path(model_path).stat().st_size / 1024 / 1024:.2f} MB")

    app = create_app(
        model=model_path,
        n_gpu_layers=0,
        n_ctx=2048,
        verbose=True
    )

    logger.info("启动服务器: http://localhost:8000")
    uvicorn.run(app, host="0.0.0.0", port=8000)

if __name__ == "__main__":
    main()
