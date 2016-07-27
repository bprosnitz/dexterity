package main

import (
  "flag"
  "fmt"
  "log"
  "encoding/binary"
  "os"
  "io"
)

const dexMagic = "dex\n035\000"

func main() {
  flag.Parse()

  if flag.NArg() < 1 {
    log.Fatalf("please specify dex filename")
  }
  filename := flag.Arg(0)
  f, err := os.Open(filename)
  if err != nil {
    log.Fatal(err)
  }
  dex := &Dex{}
  if err := readDex(f, dex); err != nil {
    log.Fatal(err)
  }
  fmt.Printf("%#v\n", dex)
  f.Close()
}

func readDex(r io.Reader, dex *Dex) error {
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

type Dex struct {
  Header DexHeader
  StringIds []DexStringIdItem
  TypeIds []DexTypeIdItem
  ProtoIds []DexProtoIdItem
  FieldIds []DexFieldIdItem
  MethodIds []DexMethodIdItem
  ClassDefs []DexClassDefItem

  DataOffsets []DataOffset
}

type DataOffset struct {
  Offset uint32
  StringItem *DexStringIdItem
  ClassDefItem *DexClassDefItem
}

type DexHeader struct {
  Checksum uint32
  Signature [20]byte
  FileSize uint32
  HeaderSize uint32
  EndianTag uint32
  LinkSize uint32
  LinkOff uint32
  MapOff uint32
  StringIdsSize uint32
  StringIdsOff uint32
  TypeIdsSize uint32
  TypeIdsOff uint32
  ProtoIdsSize uint32
  ProtoIdsOff uint32
  FieldIdsSize uint32
  FieldIdsOff uint32
  MethodIdsSize uint32
  MethodIdsOff uint32
  ClassDefsSize uint32
  ClassDefsOff uint32
  DataSize uint32
  DataOff uint32
}

type DexStringIdItem struct {
  StringDataOff uint32
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

func read(r io.Reader, i interface{}) error {
    return binary.Read(r, binary.LittleEndian, i)
}

func consume(r io.Reader, n int) error {
  b := make([]byte, n)
  on, err := r.Read(b)
  if err != nil {
    return err
  }
  if on != n {
    return fmt.Errorf("invalid number of bytes read")
  }
  return nil
}
