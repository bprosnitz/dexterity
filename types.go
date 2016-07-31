package dexterity

import (
  "io"
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
  AnnotationsOff *DexAnnotationsDirectoryItem
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

type DexAnnotationsDirectoryItem struct {
  ClassAnnotationsOff uint32
  FieldSize uint32 `listsize:"AnnotationFields"`
  AnnotatedMethodsSize uint32 `listsize:"AnnotationMethods"`
  AnnotationParametersSize uint32 `listsize:"AnnotationParameters"`
  FieldAnnotations []DexFieldAnnotation `listtag:"AnnotationFields"`
  MethodAnnotations []DexMethodAnnotation `listtag:"AnnotationMethods"`
  ParameterAnnotations []DexParameterAnnotations `listtag:"AnnotationParameters"`
}

type DexFieldAnnotation struct {
  FieldIdx uint32
  AnnotationsOff uint32
}

type DexMethodAnnotation struct {
  MethodIdx uint32
  AnnotationsOff uint32
}

type DexParameterAnnotations struct {
  MethodIdx uint32
  AnnotationsOff uint32
}

type DexAnnotationItem struct {
  Visibility uint8
  Annotation DexEncodedAnnotation
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

type DexEncodedAnnotation struct {
  TypeIdx decode.Uleb
  Size decode.Uleb `listsize:"AnnotationElements"`
  Elements []DexAnnotationElement `listtag:"AnnotationElements"`
}

type DexAnnotationElement struct {
  NameIdx decode.Uleb
  Value DexEncodedValue
}

type DexEncodedValue struct {
  T uint8
  V []byte
}

func (ev *DexEncodedValue) Read(r io.Reader) error {
panic("S")
}
