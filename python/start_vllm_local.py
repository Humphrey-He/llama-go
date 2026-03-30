#!/usr/bin/env python3
import os
import sys
import logging
from pathlib import Path

logging.basicConfig(
    level=logging.INFO,
    format='%(asctime)s - %(levelname)s - %(message)s',
    stream=sys.stdout
)
logger = logging.getLogger(__name__)

def main():
    model_path = "internal/model/Dolphin3.0-Llama3.1-8B-Q4_K_S.gguf"

    logger.info("=" * 60)
    logger.info("启动 vLLM 服务器")
    logger.info("=" * 60)

    # 检查模型文件
    if not Path(model_path).exists():
        logger.error(f"模型文件不存在: {model_path}")
        sys.exit(1)

    logger.info(f"模型文件: {model_path}")
    logger.info(f"模型大小: {Path(model_path).stat().st_size / 1024 / 1024:.2f} MB")

    # 检查 GPU
    try:
        import torch
        logger.info(f"PyTorch: {torch.__version__}")
        logger.info(f"CUDA 可用: {torch.cuda.is_available()}")
        if torch.cuda.is_available():
            logger.info(f"GPU: {torch.cuda.get_device_name(0)}")
            device = "cuda"
        else:
            logger.warning("GPU 不可用，使用 CPU")
            device = "cpu"
    except ImportError:
        logger.warning("PyTorch 未安装，使用 CPU")
        device = "cpu"

    # 启动 vLLM
    logger.info(f"设备: {device.upper()}")
    logger.info("启动中...")

    cmd = f"python -m vllm.entrypoints.openai.api_server --model {model_path} --port 8000"
    if device == "cpu":
        cmd += " --device cpu --dtype float32"

    os.system(cmd)

if __name__ == "__main__":
    main()
