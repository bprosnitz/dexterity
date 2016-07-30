package dexterity

import (
  "fmt"
  "sort"
  "io"
)

const dexMagic = "dex\n035\000"

func Process(r io.ReadSeeker, offset uint32, value interface{}) error {
//  if offset != 0xffffffff {
    origOffset, _ := r.Seek(0, 1)
    defer r.Seek(origOffset, 0)
    if _, err := r.Seek(int64(offset), 0); err != nil {
      return err
    }

  switch x := value.(type) {
  case *Dex:
    if err := readDexHeader(r, &x.Header); err != nil {
      return err
    }
    x.StringIds = make([]DexStringIdItem, x.Header.StringIdsSize)
    for i := 0; i < int(x.Header.StringIdsSize); i++ {
    //  offsetHere, _ := r.Seek(0, 1)
      //fmt.Printf("offsetHere: %v\n", offsetHere)
    //  x.StringIds[i]
  /*    if err := Process(r, 0xffffffff, &x.StringIds[i]); err != nil {
        return err
      }
    //  r.Seek(offsetHere + 4, 0)
    }
  case *DexStringIdItem:*/
    if err := read(r, &x.StringIds[i].StringDataOff); err != nil {
      return err
    }
    if err := Process(r, offset, &x.StringIds[i].StringData); err != nil {
      return err
    }
  }
  case *DexStringData:
    if err := read(r, &x.Utf16Size); err != nil {
      return err
    }
    str, err := readMutf8(r)
    if err != nil {
      return err
    }
    x.Value = str
  }
  return nil
}

func ReadDex(r io.Reader, dex *Dex) error {
  err := readDexHeader(r, &dex.Header)
  if err != nil {
    return err
  }

  dex.StringIds = make([]DexStringIdItem, dex.Header.StringIdsSize)
  for i := 0; i < int(dex.Header.StringIdsSize); i++ {
    if err := read(r, &dex.StringIds[i].StringDataOff); err != nil {
      return err
    }
    dex.DataOffsets = append(dex.DataOffsets, DataOffset{
        Offset: dex.StringIds[i].StringDataOff,
        StringItem: &dex.StringIds[i],
      })
  }

  dex.TypeIds = make([]DexTypeIdItem, dex.Header.TypeIdsSize)
  for i := 0; i < int(dex.Header.TypeIdsSize); i++ {
    if err := read(r, &dex.TypeIds[i].DescriptorIdx); err != nil {
      return err
    }
  }

  dex.ProtoIds = make([]DexProtoIdItem, dex.Header.ProtoIdsSize)
  for i := 0; i < int(dex.Header.ProtoIdsSize); i++ {
    if err := read(r, &dex.ProtoIds[i].ShortyIdx); err != nil {
      return err
    }
    if err := read(r, &dex.ProtoIds[i].ReturnTypeIdx); err != nil {
      return err
    }
    if err := read(r, &dex.ProtoIds[i].ParametersOff); err != nil {
      return err
    }
  }

  dex.FieldIds = make([]DexFieldIdItem, dex.Header.FieldIdsSize)
  for i := 0; i < int(dex.Header.FieldIdsSize); i++ {
    if err := read(r, &dex.FieldIds[i].ClassIdx); err != nil {
      return err
    }
    if err := read(r, &dex.FieldIds[i].TypeIdx); err != nil {
      return err
    }
    if err := read(r, &dex.FieldIds[i].NameIdx); err != nil {
      return err
    }
  }

  dex.MethodIds = make([]DexMethodIdItem, dex.Header.MethodIdsSize)
  for i := 0; i < int(dex.Header.MethodIdsSize); i++ {
    if err := read(r, &dex.MethodIds[i].ClassIdx); err != nil {
      return err
    }
    if err := read(r, &dex.MethodIds[i].ProtoIdx); err != nil {
      return err
    }
    if err := read(r, &dex.MethodIds[i].NameIdx); err != nil {
      return err
    }
  }

  dex.ClassDefs = make([]DexClassDefItem, dex.Header.ClassDefsSize)
  for i := 0; i < int(dex.Header.ClassDefsSize); i++ {
    if err := read(r, &dex.ClassDefs[i].ClassIdx); err != nil {
      return err
    }
    if err := read(r, &dex.ClassDefs[i].AccessFlags); err != nil {
      return err
    }
    if err := read(r, &dex.ClassDefs[i].SuperclassIdx); err != nil {
      return err
    }
    if err := read(r, &dex.ClassDefs[i].InterfacesOff); err != nil {
      return err
    }
    if err := read(r, &dex.ClassDefs[i].SourceFileIdx); err != nil {
      return err
    }
    if err := read(r, &dex.ClassDefs[i].AnnotationsOff); err != nil {
      return err
    }
    if err := read(r, &dex.ClassDefs[i].ClassDataOff); err != nil {
      return err
    }
    dex.DataOffsets = append(dex.DataOffsets, DataOffset{
        Offset: dex.ClassDefs[i].ClassDataOff,
        ClassDefItem: &dex.ClassDefs[i],
      })
    if err := read(r, &dex.ClassDefs[i].StaticValuesOff); err != nil {
      return err
    }
  }

  sort.Sort(dex.DataOffsets)
  return nil
}

func readDexHeader(r io.Reader, header *DexHeader) error {
  readMagic := make([]byte, 8)
  if _, err := r.Read(readMagic); err != nil {
    return err
  }
  if string(readMagic) != dexMagic {
    return fmt.Errorf("invalid magic string %q", string(readMagic))
  }

  if err := read(r, &header.Checksum); err != nil {
    return err
  }
  if n, err := r.Read(header.Signature[:]); err != nil {
    return err
  } else if n != 20 {
    return fmt.Errorf("partial signature bytes read")
  }
  if err := read(r, &header.FileSize); err != nil {
    return err
  }
  if err := read(r, &header.HeaderSize); err != nil {
    return err
  }
  if err := read(r, &header.EndianTag); err != nil {
    return err
  }
  if err := read(r, &header.LinkSize); err != nil {
    return err
  }
  if err := read(r, &header.LinkOff); err != nil {
    return err
  }
  if err := read(r, &header.MapOff); err != nil {
    return err
  }
  if err := read(r, &header.StringIdsSize); err != nil {
    return err
  }
  if err := read(r, &header.StringIdsOff); err != nil {
    return err
  }
  if err := read(r, &header.TypeIdsSize); err != nil {
    return err
  }
  if err := read(r, &header.TypeIdsOff); err != nil {
    return err
  }
  if err := read(r, &header.ProtoIdsSize); err != nil {
    return err
  }
  if err := read(r, &header.ProtoIdsOff); err != nil {
    return err
  }
  if err := read(r, &header.FieldIdsSize); err != nil {
    return err
  }
  if err := read(r, &header.FieldIdsOff); err != nil {
    return err
  }
  if err := read(r, &header.MethodIdsSize); err != nil {
    return err
  }
  if err := read(r, &header.MethodIdsOff); err != nil {
    return err
  }
  if err := read(r, &header.ClassDefsSize); err != nil {
    return err
  }
  if err := read(r, &header.ClassDefsOff); err != nil {
    return err
  }
  if err := read(r, &header.DataSize); err != nil {
    return err
  }
  if err := read(r, &header.DataOff); err != nil {
    return err
  }

  return nil
}
/*
func readDexClassDefData(r io.ReadSeeker, offset uint32, cd *DexClassDefData) error {
  if _, err := r.Seek(int64(offset), 0); err != nil {
    return err
  }
  staticFieldsSize, err := readUleb(r)
  if err != nil {
    return err
  }
  instanceFieldSize, err := readUleb(r)
  if err != nil {
    return err
  }
  directMethodsSize, err := readUleb(r)
  if err != nil {
    return err
  }
  virtualMethodsSize, err := readUleb(r)
  if err != nil {
    return err
  }

  cd.StaticFields = make([]DexEncodedField, staticFieldsSize)
  for i := range cd.StaticFields {
    if err := readEncodedField(&cd.StaticFields[i]); err != nil {
      return err
    }
  }
  cd.InstanceFields = make([]DexEncodedField, instanceFieldSize)
  for i := range cd.InstanceFields {
    if err := readEncodedField(&cd.InstanceFields[i]); err != nil {
      return err
    }
  }
  cd.DirectMethods = make([]DexEncodedField, directMethodsSize)
  for i := range cd.DirectMethods {
    if err := readEncodedField(&cd.DirectMethods[i]); err != nil {
      return err
    }
  }
  cd.VirtualMethods = make([]DexEncodedField, virtualMethodsSize)
  for i := range cd.VirtualMethods {
    if err := readEncodedField(&cd.VirtualMethods[i]); err != nil {
      return err
    }
  }
  return nil
}

func readEncodedField(r io.Reader, ef *DexEncodedField) error {
  var err error
  ef.FieldIdxDiff, err = readUleb(r)
  if err != nil {
    return err
  }
  ef.AccessFlags, err = readUleb(r)
  return err
}

func readEncodedMethod(r io.Reader, ef *DexEncodedMethod) error {
  var err error
  ef.MethodIdxDiff, err = readUleb(r)
  if err != nil {
    return err
  }
  ef.AccessFlags, err = readUleb(r)
  if err != nil {
    return err
  }
  ef.CodeOff, err = readUleb(r)
  return err
}
*/
