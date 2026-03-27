import torch
import numpy as np
from typing import Tuple, Optional
from transformers import AutoTokenizer, AutoModelForCausalLM
from session_manager import SESSION_STORE, SessionState

class InferenceEngine:
    """推理引擎"""
    def __init__(self, model_name: str = "TinyLlama/TinyLlama-1.1B-Chat-v1.0"):
        self.model_name = model_name
        self.tokenizer = AutoTokenizer.from_pretrained(model_name)
        self.model = AutoModelForCausalLM.from_pretrained(
            model_name,
            torch_dtype=torch.float32,
            device_map="auto"
        )
        self.model.eval()

    def prefill(self, session_id: str, prompt: str, max_new_tokens: int = 128,
                temperature: float = 0.7, top_p: float = 0.9) -> dict:
        """Prefill 阶段：首次推理"""
        import time
        start_time = time.time()

        # 创建或获取会话
        session = SESSION_STORE.get_session(session_id)
        if not session:
            session = SESSION_STORE.create_session(session_id)

        # Tokenize
        inputs = self.tokenizer(prompt, return_tensors="pt")
        input_ids = inputs["input_ids"].to(self.model.device)
        attention_mask = inputs.get("attention_mask", torch.ones_like(input_ids))

        # Forward pass with cache
        with torch.no_grad():
            outputs = self.model(
                input_ids=input_ids,
                attention_mask=attention_mask,
                use_cache=True,
                return_dict=True
            )

        # 保存 KV Cache
        session.past_key_values = outputs.past_key_values
        session.input_ids = input_ids
        session.attention_mask = attention_mask
        session.token_count = input_ids.shape[1]

        # 生成第一个 token
        logits = outputs.logits[:, -1, :]
        next_token = self._sample(logits, temperature, top_p)

        ttft_ms = (time.time() - start_time) * 1000

        return {
            "session_id": session_id,
            "generated_text": self.tokenizer.decode(next_token),
            "prompt_tokens": input_ids.shape[1],
            "generated_tokens": 1,
            "cache_hit": False,
            "mode": "prefill",
            "ttft_ms": ttft_ms,
            "total_latency_ms": ttft_ms
        }

    def decode(self, session_id: str, new_prompt: str, max_new_tokens: int = 128,
               temperature: float = 0.7, top_p: float = 0.9) -> dict:
        """Decode 阶段：增量推理"""
        import time
        start_time = time.time()

        session = SESSION_STORE.get_session(session_id)
        if not session or session.past_key_values is None:
            # Fallback to prefill
            return self.prefill(session_id, new_prompt, max_new_tokens, temperature, top_p)

        # Tokenize 新输入
        new_inputs = self.tokenizer(new_prompt, return_tensors="pt")
        new_input_ids = new_inputs["input_ids"].to(self.model.device)

        generated_tokens = []
        current_input_ids = new_input_ids

        for _ in range(max_new_tokens):
            with torch.no_grad():
                outputs = self.model(
                    input_ids=current_input_ids,
                    past_key_values=session.past_key_values,
                    use_cache=True,
                    return_dict=True
                )

            # 更新 KV Cache
            session.past_key_values = outputs.past_key_values

            # Sample next token
            logits = outputs.logits[:, -1, :]
            next_token = self._sample(logits, temperature, top_p)
            generated_tokens.append(next_token.item())

            # 检查 EOS
            if next_token.item() == self.tokenizer.eos_token_id:
                break

            current_input_ids = next_token.unsqueeze(0)

        session.token_count += len(generated_tokens)

        ttft_ms = (time.time() - start_time) * 1000

        return {
            "session_id": session_id,
            "generated_text": self.tokenizer.decode(generated_tokens),
            "prompt_tokens": new_input_ids.shape[1],
            "generated_tokens": len(generated_tokens),
            "cache_hit": True,
            "mode": "decode",
            "ttft_ms": ttft_ms,
            "total_latency_ms": ttft_ms
        }

    def _sample(self, logits: torch.Tensor, temperature: float = 0.7,
                top_p: float = 0.9) -> torch.Tensor:
        """采样下一个 token"""
        logits = logits / temperature

        # Top-p 采样
        sorted_logits, sorted_indices = torch.sort(logits, descending=True)
        cumsum = torch.cumsum(torch.softmax(sorted_logits, dim=-1), dim=-1)
        sorted_indices_to_remove = cumsum > top_p
        sorted_indices_to_remove[..., 0] = False

        indices_to_remove = sorted_indices[sorted_indices_to_remove]
        logits[:, indices_to_remove] = float('-inf')

        probs = torch.softmax(logits, dim=-1)
        next_token = torch.multinomial(probs, num_samples=1)

        return next_token
