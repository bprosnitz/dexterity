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
    if err := read(r, &dex.StringIds[i]); err != nil {
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


func read(r io.Reader, i interface{}) error {
    return binary.Read(r, binary.LittleEndian, i)
}
