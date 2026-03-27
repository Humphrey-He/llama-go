from fastapi import FastAPI
from pydantic import BaseModel
import torch
from transformers import AutoTokenizer, AutoModelForCausalLM

app = FastAPI()

# 轻量模型（MVP 用小模型，快速运行）
MODEL_NAME = "TinyLlama/TinyLlama-1.1B-Chat-v1.0"
tokenizer = AutoTokenizer.from_pretrained(MODEL_NAME)
model = AutoModelForCausalLM.from_pretrained(MODEL_NAME, torch_dtype=torch.float32)

class InferRequest(BaseModel):
    prompt: str
    kv_cache: list = None

class InferResponse(BaseModel):
    text: str
    kv: dict
    token_num: int

@app.post("/infer", response_model=InferResponse)
async def infer(req: InferRequest):
    # 1. 分词
    inputs = tokenizer(req.prompt, return_tensors="pt")

    # 2. 推理 + KV 缓存输出
    with torch.no_grad():
        outputs = model(**inputs, use_cache=True)

    # 3. 提取 KV 缓存
    past_kv = outputs.past_key_values
    kv = {
        "keys": past_kv[0][0].squeeze(0).tolist(),
        "vals": past_kv[0][1].squeeze(0).tolist()
    }

    # 4. 生成文本
    text = tokenizer.decode(outputs.logits.argmax(dim=-1))

    return InferResponse(
        text=text,
        kv=kv,
        token_num=len(inputs.input_ids[0])
    )

if __name__ == "__main__":
    import uvicorn
    uvicorn.run(app, host="0.0.0.0", port=8000)