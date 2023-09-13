from google.protobuf import descriptor as _descriptor
from google.protobuf import message as _message
from typing import ClassVar as _ClassVar, Mapping as _Mapping, Optional as _Optional, Union as _Union

DESCRIPTOR: _descriptor.FileDescriptor

class Empty(_message.Message):
    __slots__ = []
    def __init__(self) -> None: ...

class ConfigResponse(_message.Message):
    __slots__ = ["success", "data"]
    SUCCESS_FIELD_NUMBER: _ClassVar[int]
    DATA_FIELD_NUMBER: _ClassVar[int]
    success: bool
    data: bytes
    def __init__(self, success: bool = ..., data: _Optional[bytes] = ...) -> None: ...

class InTypes(_message.Message):
    __slots__ = ["index", "type"]
    INDEX_FIELD_NUMBER: _ClassVar[int]
    IN_FIELD_NUMBER: _ClassVar[int]
    TYPE_FIELD_NUMBER: _ClassVar[int]
    index: int
    type: str
    def __init__(self, index: _Optional[int] = ..., type: _Optional[str] = ..., **kwargs) -> None: ...

class OutType(_message.Message):
    __slots__ = ["index", "out", "type"]
    INDEX_FIELD_NUMBER: _ClassVar[int]
    OUT_FIELD_NUMBER: _ClassVar[int]
    TYPE_FIELD_NUMBER: _ClassVar[int]
    index: int
    out: bytes
    type: str
    def __init__(self, index: _Optional[int] = ..., out: _Optional[bytes] = ..., type: _Optional[str] = ...) -> None: ...

class FunctionRequest(_message.Message):
    __slots__ = ["function", "struct"]
    IN_FIELD_NUMBER: _ClassVar[int]
    FUNCTION_FIELD_NUMBER: _ClassVar[int]
    STRUCT_FIELD_NUMBER: _ClassVar[int]
    function: str
    struct: str
    def __init__(self, function: _Optional[str] = ..., struct: _Optional[str] = ..., **kwargs) -> None: ...

class ErrorMessage(_message.Message):
    __slots__ = ["code", "message", "data"]
    CODE_FIELD_NUMBER: _ClassVar[int]
    MESSAGE_FIELD_NUMBER: _ClassVar[int]
    DATA_FIELD_NUMBER: _ClassVar[int]
    code: int
    message: str
    data: bytes
    def __init__(self, code: _Optional[int] = ..., message: _Optional[str] = ..., data: _Optional[bytes] = ...) -> None: ...

class StatResponse(_message.Message):
    __slots__ = ["avgresponse"]
    AVGRESPONSE_FIELD_NUMBER: _ClassVar[int]
    avgresponse: int
    def __init__(self, avgresponse: _Optional[int] = ...) -> None: ...

class FunctionResponse(_message.Message):
    __slots__ = ["data", "success", "client", "error"]
    DATA_FIELD_NUMBER: _ClassVar[int]
    SUCCESS_FIELD_NUMBER: _ClassVar[int]
    CLIENT_FIELD_NUMBER: _ClassVar[int]
    ERROR_FIELD_NUMBER: _ClassVar[int]
    data: bytes
    success: bool
    client: str
    error: ErrorMessage
    def __init__(self, data: _Optional[bytes] = ..., success: bool = ..., client: _Optional[str] = ..., error: _Optional[_Union[ErrorMessage, _Mapping]] = ...) -> None: ...
