package dexterity

import (
  "io"
  "encoding/binary"
  "fmt"
)

func readUleb(r io.Reader) (uint32, error) {
  b := make([]byte, 1)
  hasNext := true
  shiftAmt := uint(0)
  var value uint32
  for hasNext {
    if n, err := r.Read(b); err != nil {
      return 0, err
    } else if n != 1 {
      return 0, fmt.Errorf("expected to read 1 byte, but read none")
    }
    hasNext = b[0] & 0x80 != 0
    value |= uint32(b[0] & 0x7f) << shiftAmt
    shiftAmt += 7
  }
  return value, nil
}

func readUlebP1(r io.Reader) (uint32, error) {
  u, err := readUleb(r)
  return u - 1, err
}

func readSleb(r io.Reader) (int32, error) {
  b := make([]byte, 1)
  hasNext := true
  shiftAmt := uint(0)
  var signBit bool
  var value uint32
  for hasNext {
    if n, err := r.Read(b); err != nil {
      return 0, err
    } else if n != 1 {
      return 0, fmt.Errorf("expected to read 1 byte, but read none")
    }
    hasNext = b[0] & 0x80 != 0
    signBit = b[0] & 0x40 != 0
    value |= uint32(b[0] & 0x7f) << shiftAmt
    shiftAmt += 7
  }
  if signBit {
    value |= 0xffffffff << shiftAmt
  }
  return int32(value), nil
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
