package dexterity

import (
  "io"
  "github.com/bprosnitz/dexterity/decode"
)

type Dex struct {
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

  StringIds []DexStringIdItem `listtag:"StringIds"`
  TypeIds []DexTypeIdItem `listtag:"TypeIds"`
  ProtoIds []DexProtoIdItem `listtag:"ProtoIds"`
  FieldIds []DexFieldIdItem `listtag:"FieldIds"`
  MethodIds []DexMethodIdItem `listtag:"MethodIds"`
  ClassDefs []DexClassDefItem `listtag:"ClassDefs"`
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
  Parameters *DexTypeList
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
  Interfaces *DexTypeList
  SourceFileIdx uint32
  Annotations *DexAnnotationsDirectoryItem
  ClassData *DexClassDefData
  StaticValues *DexEncodedArray
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
  ClassAnnotations *DexAnnotationSet
  FieldSize uint32 `listsize:"AnnotationFields"`
  AnnotatedMethodsSize uint32 `listsize:"AnnotationMethods"`
  AnnotationParametersSize uint32 `listsize:"AnnotationParameters"`
  FieldAnnotations []DexFieldAnnotation `listtag:"AnnotationFields"`
  MethodAnnotations []DexMethodAnnotation `listtag:"AnnotationMethods"`
  ParameterAnnotations []DexParameterAnnotations `listtag:"AnnotationParameters"`
}

type DexAnnotationSet struct {
  Size uint32 `listsize:"AnnotationSet"`
  Items []*DexAnnotationItem `listtag:"AnnotationSet"`
}

type DexFieldAnnotation struct {
  FieldIdx uint32
  Annotations *DexAnnotationSet
}

type DexMethodAnnotation struct {
  MethodIdx uint32
  Annotations *DexAnnotationSet
}

type DexParameterAnnotations struct {
  MethodIdx uint32
  Annotations *DexAnnotationSet
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

type DexTypeList struct {
  Size uint32 `listsize:"TypeItems"`
  Items []DexTypeItem `listtag:"TypeItems"`
}

type DexEncodedArray struct {
  Size decode.Uleb `listsize:"ArrayItems"`
  Values []DexEncodedValue `listtag:"ArrayItems"`
}

type DexTypeItem struct {
  TypeIdx uint16
}

type DexCodeItem struct {
  
}

type DexEncodedValue struct {
  T uint8
  V []byte
}

func (ev *DexEncodedValue) Read(r io.Reader) error {
  panic("not yet implemented")
}
