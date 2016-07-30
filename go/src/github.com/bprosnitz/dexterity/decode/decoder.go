package decode

import (
  "io"
  "fmt"
  "reflect"
)

type Uleb uint32
type Ulebp1 uint32
type Sleb int32

func Decode(r io.ReadSeeker, x interface{}) error {
  d := decoder{
    r: r,
  }
  return d.Decode(reflect.ValueOf(x))
}

type decoderStackEntry struct {
  pos int
}

type decoder struct {
  r io.ReadSeeker
  stack []decoderStackEntry
}

func (d *decoder) Decode(rv reflect.Value) error {
  for rv.Kind() == reflect.Ptr {
    rv = rv.Elem()
  }
  rt := rv.Type()
  for i := 0; i < rv.NumField(); i++ {
    switch x := rv.Field(i).Interface().(type) {
    case uint32:
      if err := read(d.r, x); err != nil {
        return err
      }
    case Uleb:
      v, err := readUleb(d.r)
      if err != nil {
        return err
      }
      rv.Field(i).SetUint(uint64(v))
    case Ulebp1:
      v, err := readUlebP1(d.r)
      if err != nil {
        return err
      }
      rv.Field(i).SetUint(uint64(v))
    case Sleb:
      v, err := readSleb(d.r)
      if err != nil {
        return err
      }
      rv.Field(i).SetInt(int64(v))
    default:
      return fmt.Errorf("unhandled type: %v", rt.Field(i))
    }
  }
  return nil
}
