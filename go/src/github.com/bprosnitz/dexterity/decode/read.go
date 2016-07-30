package decode

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

func readUint32(r io.Reader) (uint32, error) {
  var u uint32
  err := read(r, &u)
  return u, err
}

func readUint16(r io.Reader) (uint16, error) {
  var u uint16
  err := read(r, &u)
  return u, err
}

func readMutf8(r io.Reader) (string, error) {
  // TODO(bprosnitz) Fully support MUTF-8 modifications
  b := make([]byte, 1)
  var bytes []byte
  for {
    if _, err := r.Read(b); err != nil {
        return "", err
    }
    if b[0] == 0 {
      return string(bytes), nil
    }
    bytes = append(bytes, b[0])
  }
}
