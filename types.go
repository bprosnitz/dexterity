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
  StringIdsSize uint32 `listsize:"StringIds"`
  StringIdsOff uint32
  TypeIdsSize uint32 `listsize:"TypeIds"`
  TypeIdsOff uint32
  ProtoIdsSize uint32 `listsize:"ProtoIds"`
  ProtoIdsOff uint32
  FieldIdsSize uint32 `listsize:"FieldIds"`
  FieldIdsOff uint32
  MethodIdsSize uint32 `listsize:"MethodIds"`
  MethodIdsOff uint32
  ClassDefsSize uint32 `listsize:"ClassDefs"`
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
  ClassDataOff *DexClassDefData
  StaticValuesOff uint32
}

type DexClassDefData struct {
  StaticFieldsSize decode.Uleb `listsize:"StaticFields"`
  InstanceFieldSize decode.Uleb `listsize:"InstanceFields"`
  DirectMethodsSize decode.Uleb `listsize:"DirectMethods"`
  VirtualMethodsSize decode.Uleb `listsize:"VirtualMethods"`
  StaticFields []DexEncodedField `listtag:"StaticFields"`
  InstanceFields []DexEncodedField `listtag:"InstanceFields"`
  DirectMethods []DexEncodedMethod `listtag:"DirectMethods"`
  VirtualMethods []DexEncodedMethod `listtag:"VirtualMethods"`
}

type DexEncodedField struct {
  FieldIdxDiff decode.Uleb
  AccessFlags decode.Uleb
}

type DexEncodedMethod struct {
  MethodIdxDiff decode.Uleb
  AccessFlags decode.Uleb
  CodeOff decode.Uleb
}
