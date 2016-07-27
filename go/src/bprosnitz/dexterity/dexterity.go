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
  if err := processDex(f); err != nil {
    log.Fatal(err)
  }
  f.Close()
}

func processDex(r io.Reader) error {
  header, err := readDexHeader(r)
  if err != nil {
    return err
  }
  fmt.Printf("header: %#v\n", header)

  return nil
}

func readDexHeader(r io.Reader) (*DexHeader, error) {
  readMagic := make([]byte, 8)
  if _, err := r.Read(readMagic); err != nil {
    return nil, err
  }
  if string(readMagic) != dexMagic {
    return nil, fmt.Errorf("invalid magic string %q", string(readMagic))
  }

  header := &DexHeader{}
  if err := read(r, &header.Checksum); err != nil {
    return nil, err
  }
  if n, err := r.Read(header.Signature[:]); err != nil {
    return nil, err
  } else if n != 20 {
    return nil, fmt.Errorf("partial signature bytes read")
  }
  if err := read(r, &header.FileSize); err != nil {
    return nil, err
  }
  if err := read(r, &header.HeaderSize); err != nil {
    return nil, err
  }
  if err := read(r, &header.EndianTag); err != nil {
    return nil, err
  }
  if err := read(r, &header.LinkSize); err != nil {
    return nil, err
  }
  if err := read(r, &header.LinkOff); err != nil {
    return nil, err
  }
  if err := read(r, &header.MapOff); err != nil {
    return nil, err
  }
  if err := read(r, &header.StringIdsSize); err != nil {
    return nil, err
  }
  if err := read(r, &header.StringIdsOff); err != nil {
    return nil, err
  }
  if err := read(r, &header.TypeIdsSize); err != nil {
    return nil, err
  }
  if err := read(r, &header.TypeIdsOff); err != nil {
    return nil, err
  }
  if err := read(r, &header.ProtoIdsSize); err != nil {
    return nil, err
  }
  if err := read(r, &header.ProtoIdsOff); err != nil {
    return nil, err
  }
  if err := read(r, &header.FieldIdsSize); err != nil {
    return nil, err
  }
  if err := read(r, &header.FieldIdsOff); err != nil {
    return nil, err
  }
  if err := read(r, &header.MethodIdsSize); err != nil {
    return nil, err
  }
  if err := read(r, &header.MethodIdsOff); err != nil {
    return nil, err
  }
  if err := read(r, &header.ClassDefsSize); err != nil {
    return nil, err
  }
  if err := read(r, &header.ClassDefsOff); err != nil {
    return nil, err
  }
  if err := read(r, &header.DataSize); err != nil {
    return nil, err
  }
  if err := read(r, &header.DataOff); err != nil {
    return nil, err
  }

  return header, nil
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


func read(r io.Reader, i interface{}) error {
    return binary.Read(r, binary.LittleEndian, i)
}
