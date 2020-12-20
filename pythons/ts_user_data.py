import time
from typing import Any, Dict, Union, Tuple


StringMap = Dict[str, str]
StringAnyValueMap = Dict[str, Any]
MaybeManyValueStringMap = Dict[str, Union[Tuple[str], str]]


class Request:
    user_data: StringMap
    query: MaybeManyValueStringMap
    headers: MaybeManyValueStringMap
    body: str

    def __init__(self, data):
        self.query = data["query"]
        self.headers = data["headers"]
        self.user_data = data["user_data"]
        self.body = data.get("body", "")

    def resp(self, response: StringAnyValueMap):
        return {"response": response, "user_data": self.user_data}


def handler(data: StringAnyValueMap):
    r = Request(data)
    if r.query.get("time"):
        r.user_data["time"] = str(time.time())
    return r.resp({
        "message": "hello",
        "echo": {
            "body": r.body,
            "user_data": r.user_data,
            "headers": r.headers,
            "query": r.query,
        }
    })
