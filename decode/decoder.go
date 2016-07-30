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
  rv := reflect.ValueOf(x)
  for rv.Kind() == reflect.Ptr {
    rv = rv.Elem()
  }
  return d.Decode(rv, "")
}

type decoder struct {
  r io.ReadSeeker
  sizes map[string]uint32
}

func (d *decoder) Decode(rv reflect.Value, tag reflect.StructTag) error {
  rt := rv.Type()
  switch rv.Interface().(type) {
  case uint8:
    var b [1]byte
    if _, err := d.r.Read(b[:]); err != nil {
      return err
    }
    rv.SetUint(uint64(b[0]))
    return nil
  case uint16:
    v, err := readUint16(d.r)
    if err != nil {
      return err
    }
    rv.SetUint(uint64(v))
    return nil
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
  case string:
    v, err := readMutf8(d.r)
    if err != nil {
      return err
    }
    rv.SetString(v)
    return nil
  }
  switch rt.Kind() {
  case reflect.Array:
    for i := 0; i < rt.Len(); i++ {
      if err := d.Decode(rv.Index(i), ""); err != nil {
        return err
      }
    }
    return nil
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
  case reflect.Ptr:
    offset, err := readUint32(d.r)
    if err != nil {
      return err
    }
    origOffset, _ := d.r.Seek(0, 1)
    if _, err := d.r.Seek(int64(offset), 0); err != nil {
      return err
    }
    childRv := reflect.New(rt.Elem())
    rv.Set(childRv)
    if err := d.Decode(childRv.Elem(), tag); err != nil {
      return err
    }
    d.r.Seek(origOffset, 0)
    return nil
  }
  return fmt.Errorf("unhandled type: %v", rt)
}
