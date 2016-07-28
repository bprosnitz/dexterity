package dexterity

import (
  "fmt"
  "sort"
  "io"
)

const dexMagic = "dex\n035\000"

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
