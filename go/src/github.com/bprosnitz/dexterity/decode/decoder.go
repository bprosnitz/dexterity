package decode

import (
  "io"
  "fmt"
  "reflect"
)

type Uleb uint32
type Ulebp1 uint32
type Sleb int32
type Size uint32

func Decode(r io.ReadSeeker, x interface{}) error {
  d := decoder{
    r: r,
    sizes: map[string]uint32{},
  }
  return d.Decode(reflect.ValueOf(x), "")
}

type decoderStackEntry struct {
  pos int
}

type decoder struct {
  r io.ReadSeeker
  sizes map[string]uint32
  stack []decoderStackEntry
}

func (d *decoder) Decode(rv reflect.Value, tag reflect.StructTag) error {
  for rv.Kind() == reflect.Ptr {
    rv = rv.Elem()
  }
  rt := rv.Type()
  switch rv.Interface().(type) {
  case uint32:
    v, err := readUint32(d.r)
    if err != nil {
      return err
    }
    rv.SetUint(uint64(v))
    return nil
  case Uleb:
    v, err := readUleb(d.r)
    if err != nil {
      return err
    }
    rv.SetUint(uint64(v))
    return nil
  case Ulebp1:
    v, err := readUlebP1(d.r)
    if err != nil {
      return err
    }
    rv.SetUint(uint64(v))
    return nil
  case Sleb:
    v, err := readSleb(d.r)
    if err != nil {
      return err
    }
    rv.SetInt(int64(v))
    return nil
  case Size:
    sizetag := tag.Get("sizetag")
    if sizetag == "" {
      return fmt.Errorf("missing sizetag on size")
    }
    v, err := readUint32(d.r)
    if err != nil {
      return err
    }
    rv.SetUint(uint64(v))
    d.sizes[sizetag] = v
    return nil
  }
  switch rt.Kind() {
  case reflect.Slice:
      sizetag := tag.Get("sizetag")
      if sizetag == "" {
        return fmt.Errorf("missing sizetag on slice")
      }
      size, ok := d.sizes[sizetag]
      if !ok {
        return fmt.Errorf("missing matching sizetag definition")
      }
      rv.Set(reflect.MakeSlice(rt, int(size), int(size)))
      for i := 0; i < int(size); i++ {
        if err := d.Decode(rv.Index(i), ""); err != nil {
          return err
        }
      }
      return nil
  case reflect.Struct:
    for i := 0; i < rv.NumField(); i++ {
      if err := d.Decode(rv.Field(i), rt.Field(i).Tag); err != nil {
        return err
      }
    }
    return nil
  }
  return fmt.Errorf("unhandled type: %v", rt)
}
