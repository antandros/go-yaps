# -*- coding: utf-8 -*-
# Generated by the protocol buffer compiler.  DO NOT EDIT!
# source: protocol.proto
"""Generated protocol buffer code."""
from google.protobuf import descriptor as _descriptor
from google.protobuf import descriptor_pool as _descriptor_pool
from google.protobuf import symbol_database as _symbol_database
from google.protobuf.internal import builder as _builder
# @@protoc_insertion_point(imports)

_sym_db = _symbol_database.Default()




DESCRIPTOR = _descriptor_pool.Default().AddSerializedFile(b'\n\x0eprotocol.proto\x12\x08protocol\"\x07\n\x05\x45mpty\"/\n\x0e\x43onfigResponse\x12\x0f\n\x07success\x18\x01 \x01(\x08\x12\x0c\n\x04\x64\x61ta\x18\x02 \x01(\x0c\"2\n\x07InTypes\x12\r\n\x05index\x18\x01 \x01(\x05\x12\n\n\x02in\x18\x02 \x01(\x0c\x12\x0c\n\x04type\x18\x03 \x01(\t\"3\n\x07OutType\x12\r\n\x05index\x18\x01 \x01(\x05\x12\x0b\n\x03out\x18\x02 \x01(\x0c\x12\x0c\n\x04type\x18\x03 \x01(\t\"R\n\x0f\x46unctionRequest\x12\x1d\n\x02in\x18\x01 \x03(\x0b\x32\x11.protocol.InTypes\x12\x10\n\x08\x66unction\x18\x02 \x01(\t\x12\x0e\n\x06struct\x18\x03 \x01(\t\";\n\x0c\x45rrorMessage\x12\x0c\n\x04\x63ode\x18\x01 \x01(\x05\x12\x0f\n\x07message\x18\x02 \x01(\t\x12\x0c\n\x04\x64\x61ta\x18\x03 \x01(\x0c\"#\n\x0cStatResponse\x12\x13\n\x0b\x61vgresponse\x18\x01 \x01(\x05\"w\n\x10\x46unctionResponse\x12\x0c\n\x04\x64\x61ta\x18\x01 \x01(\x0c\x12\x0f\n\x07success\x18\x02 \x01(\x08\x12\x0e\n\x06\x63lient\x18\x03 \x01(\t\x12*\n\x05\x65rror\x18\x04 \x01(\x0b\x32\x16.protocol.ErrorMessageH\x00\x88\x01\x01\x42\x08\n\x06_error2\xfb\x01\n\x0ePluginProtocol\x12<\n\rRequestConfig\x12\x0f.protocol.Empty\x1a\x18.protocol.ConfigResponse\"\x00\x12G\n\x0c\x43\x61llFunction\x12\x19.protocol.FunctionRequest\x1a\x1a.protocol.FunctionResponse\"\x00\x12/\n\tHeartBeat\x12\x0f.protocol.Empty\x1a\x0f.protocol.Empty\"\x00\x12\x31\n\x04Stat\x12\x0f.protocol.Empty\x1a\x16.protocol.StatResponse\"\x00\x42\x0cZ\n./protocolb\x06proto3')

_globals = globals()
_builder.BuildMessageAndEnumDescriptors(DESCRIPTOR, _globals)
_builder.BuildTopDescriptorsAndMessages(DESCRIPTOR, 'protocol_pb2', _globals)
if _descriptor._USE_C_DESCRIPTORS == False:

  DESCRIPTOR._options = None
  DESCRIPTOR._serialized_options = b'Z\n./protocol'
  _globals['_EMPTY']._serialized_start=28
  _globals['_EMPTY']._serialized_end=35
  _globals['_CONFIGRESPONSE']._serialized_start=37
  _globals['_CONFIGRESPONSE']._serialized_end=84
  _globals['_INTYPES']._serialized_start=86
  _globals['_INTYPES']._serialized_end=136
  _globals['_OUTTYPE']._serialized_start=138
  _globals['_OUTTYPE']._serialized_end=189
  _globals['_FUNCTIONREQUEST']._serialized_start=191
  _globals['_FUNCTIONREQUEST']._serialized_end=273
  _globals['_ERRORMESSAGE']._serialized_start=275
  _globals['_ERRORMESSAGE']._serialized_end=334
  _globals['_STATRESPONSE']._serialized_start=336
  _globals['_STATRESPONSE']._serialized_end=371
  _globals['_FUNCTIONRESPONSE']._serialized_start=373
  _globals['_FUNCTIONRESPONSE']._serialized_end=492
  _globals['_PLUGINPROTOCOL']._serialized_start=495
  _globals['_PLUGINPROTOCOL']._serialized_end=746
# @@protoc_insertion_point(module_scope)