#!/usr/bin/env python3
import os
import sys
import logging
from pathlib import Path

# 配置日志
logging.basicConfig(
    level=logging.INFO,
    format='%(asctime)s - %(name)s - %(levelname)s - %(message)s',
    stream=sys.stdout
)
logger = logging.getLogger(__name__)

def main():
    model_path = "/models/Dolphin3.0-Llama3.1-8B-Q4_K_S.gguf"

    logger.info("=" * 60)
    logger.info("Starting vLLM server")
    logger.info("=" * 60)

    # 检查模型文件
    if not Path(model_path).exists():
        logger.error(f"Model file not found: {model_path}")
        logger.info("Available files in /models:")
        if Path("/models").exists():
            for f in Path("/models").iterdir():
                logger.info(f"  - {f.name} ({f.stat().st_size / 1024 / 1024:.2f} MB)")
        sys.exit(1)

    logger.info(f"Model file found: {model_path}")
    logger.info(f"Model size: {Path(model_path).stat().st_size / 1024 / 1024:.2f} MB")

    # 检查 vLLM 安装
    try:
        import vllm
        logger.info(f"vLLM version: {vllm.__version__}")
    except ImportError as e:
        logger.error(f"vLLM not installed: {e}")
        logger.info("Installing vLLM...")
        os.system("pip install vllm -q")
        import vllm
        logger.info(f"vLLM installed: {vllm.__version__}")

    # 启动 vLLM
    logger.info("Starting vLLM API server...")
    logger.info(f"Device: CPU")
    logger.info(f"Port: 8000")

    os.system(
        f"python -m vllm.entrypoints.openai.api_server "
        f"--model {model_path} "
        f"--port 8000 "
        f"--device cpu "
        f"--dtype float32"
    )

if __name__ == "__main__":
    main()
