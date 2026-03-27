import time
from typing import Optional, Dict, List, Tuple
import torch
from transformers import AutoTokenizer, AutoModelForCausalLM

class SessionState:
    """会话状态管理"""
    def __init__(self, session_id: str):
        self.session_id = session_id
        self.past_key_values = None
        self.input_ids = None
        self.attention_mask = None
        self.token_count = 0
        self.last_access = time.time()
        self.created_at = time.time()

    def update_access_time(self):
        self.last_access = time.time()

    def is_expired(self, ttl_seconds: int = 3600) -> bool:
        return time.time() - self.last_access > ttl_seconds


class SessionStore:
    """会话存储管理"""
    def __init__(self, max_sessions: int = 100):
        self.sessions: Dict[str, SessionState] = {}
        self.max_sessions = max_sessions

    def create_session(self, session_id: str) -> SessionState:
        if len(self.sessions) >= self.max_sessions:
            # LRU 淘汰
            oldest_id = min(self.sessions.keys(),
                          key=lambda k: self.sessions[k].last_access)
            del self.sessions[oldest_id]

        session = SessionState(session_id)
        self.sessions[session_id] = session
        return session

    def get_session(self, session_id: str) -> Optional[SessionState]:
        session = self.sessions.get(session_id)
        if session:
            session.update_access_time()
        return session

    def delete_session(self, session_id: str):
        if session_id in self.sessions:
            del self.sessions[session_id]

    def cleanup_expired(self, ttl_seconds: int = 3600):
        expired = [sid for sid, s in self.sessions.items()
                  if s.is_expired(ttl_seconds)]
        for sid in expired:
            del self.sessions[sid]


# 全局会话存储
SESSION_STORE = SessionStore()
