from fastapi import FastAPI, HTTPException
from pydantic import BaseModel
import logging
from inference_engine import InferenceEngine
from session_manager import SESSION_STORE
import uuid

# 配置日志
logging.basicConfig(level=logging.INFO)
logger = logging.getLogger(__name__)

app = FastAPI()
engine = InferenceEngine()

class GenerateRequest(BaseModel):
    session_id: str = None
    prompt: str
    max_new_tokens: int = 128
    temperature: float = 0.7
    top_p: float = 0.9
    top_k: int = 50
    stream: bool = False

class GenerateResponse(BaseModel):
    success: bool
    data: dict

@app.post("/api/generate")
async def generate(request: GenerateRequest):
    """生成文本接口"""
    try:
        # 生成或使用提供的 session_id
        session_id = request.session_id or str(uuid.uuid4())

        # 检查是否是首次请求
        session = SESSION_STORE.get_session(session_id)
        is_first_request = session is None or session.past_key_values is None

        if is_first_request:
            result = engine.prefill(
                session_id=session_id,
                prompt=request.prompt,
                max_new_tokens=request.max_new_tokens,
                temperature=request.temperature,
                top_p=request.top_p
            )
        else:
            result = engine.decode(
                session_id=session_id,
                new_prompt=request.prompt,
                max_new_tokens=request.max_new_tokens,
                temperature=request.temperature,
                top_p=request.top_p
            )

        logger.info(f"Generate: session_id={session_id}, mode={result['mode']}, "
                   f"cache_hit={result['cache_hit']}, ttft_ms={result['ttft_ms']:.2f}")

        return GenerateResponse(success=True, data=result)

    except Exception as e:
        logger.error(f"Generate error: {str(e)}")
        raise HTTPException(status_code=500, detail=str(e))

@app.post("/api/sessions/{session_id}/clear")
async def clear_session(session_id: str):
    """清空会话"""
    try:
        SESSION_STORE.delete_session(session_id)
        logger.info(f"Session cleared: {session_id}")
        return {"success": True, "message": "session cleared"}
    except Exception as e:
        logger.error(f"Clear session error: {str(e)}")
        raise HTTPException(status_code=500, detail=str(e))

@app.get("/healthz")
async def health():
    """健康检查"""
    return {"status": "ok"}

@app.get("/metrics")
async def metrics():
    """指标接口"""
    return {
        "total_sessions": len(SESSION_STORE.sessions),
        "active_sessions": sum(1 for s in SESSION_STORE.sessions.values()
                              if not s.is_expired())
    }

if __name__ == "__main__":
    import uvicorn
    uvicorn.run(app, host="0.0.0.0", port=8000)
