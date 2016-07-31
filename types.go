package dexterity

import (
  "github.com/bprosnitz/dexterity/decode"
)

type Dex struct {
  Header DexHeader
  StringIds []DexStringIdItem `listtag:"StringIds"`
  TypeIds []DexTypeIdItem `listtag:"TypeIds"`
  ProtoIds []DexProtoIdItem `listtag:"ProtoIds"`
  FieldIds []DexFieldIdItem `listtag:"FieldIds"`
  MethodIds []DexMethodIdItem `listtag:"MethodIds"`
  ClassDefs []DexClassDefItem `listtag:"ClassDefs"`
}

type DataOffset struct {
  Offset uint32
  StringItem *DexStringIdItem
  ClassDefItem *DexClassDefItem
}

type DataOffsets []DataOffset

func (do DataOffsets) Len() int {
  return len(do)
}
func (do DataOffsets) Less(i, j int) bool {
  return do[i].Offset < do[j].Offset
}
func (do DataOffsets) Swap(i, j int) {
  do[i], do[j] = do[j], do[i]
}

type DexHeader struct {
  Magic [8]byte
  Checksum uint32
  Signature [20]byte
  FileSize uint32
  HeaderSize uint32
  EndianTag uint32
  LinkSize uint32
  LinkOff uint32
  MapOff uint32
  StringIdsSize decode.Size `listsize:"StringIds"`
  StringIdsOff uint32
  TypeIdsSize decode.Size `listsize:"TypeIds"`
  TypeIdsOff uint32
  ProtoIdsSize decode.Size `listsize:"ProtoIds"`
  ProtoIdsOff uint32
  FieldIdsSize decode.Size `listsize:"FieldIds"`
  FieldIdsOff uint32
  MethodIdsSize decode.Size `listsize:"MethodIds"`
  MethodIdsOff uint32
  ClassDefsSize decode.Size `listsize:"ClassDefs"`
  ClassDefsOff uint32
  DataSize uint32
  DataOff uint32
}

type DexStringIdItem struct {
  StringData *DexStringData
}

type DexStringData struct {
  Utf16Size decode.Uleb
  Value string
}

type DexTypeIdItem struct {
  DescriptorIdx uint32
}

type DexProtoIdItem struct {
  ShortyIdx uint32
  ReturnTypeIdx uint32
  ParametersOff uint32
}

type DexFieldIdItem struct {
  ClassIdx uint16
  TypeIdx uint16
  NameIdx uint32
}

type DexMethodIdItem struct {
  ClassIdx uint16
  ProtoIdx uint16
  NameIdx uint32
}

type DexClassDefItem struct {
  ClassIdx uint32
  AccessFlags uint32
  SuperclassIdx uint32
  InterfacesOff uint32
  SourceFileIdx uint32
  AnnotationsOff uint32
  ClassDataOff uint32
  StaticValuesOff uint32
}

type DexClassDefData struct {
  StaticFields []DexEncodedField
  InstanceFields []DexEncodedField
  DirectMethods []DexEncodedMethod
  VirtualMethods []DexEncodedMethod
}

type DexEncodedField struct {
  FieldIdxDiff uint32
  AccessFlags uint32
}

type DexEncodedMethod struct {
  MethodIdxDiff uint32
  AccessFlags uint32
  CodeOff uint32
}
