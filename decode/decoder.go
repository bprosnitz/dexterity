package decode

import (
  "io"
  "fmt"
  "reflect"
)

type Uleb uint32
type Ulebp1 uint32
type Sleb int32

type CustomReader interface {
  Read(r io.Reader, h Helper) error
}

type Helper interface {
  LookupSize(name string) (uint32, bool)
}

func Decode(r io.ReadSeeker, x interface{}) error {
  d := decoder{
    r: r,
  }
  rv := reflect.ValueOf(x)
  for rv.Kind() == reflect.Ptr {
    rv = rv.Elem()
  }
  return d.Decode(rv, "")
}

type decoderStackEntry struct {
  sizes map[string]uint32
}

type decoder struct {
  r io.ReadSeeker
  stack []decoderStackEntry
}

func (d *decoder) pop() {
  d.stack = d.stack[:len(d.stack) - 1]
}

func (d *decoder) push() {
  d.stack = append(d.stack, decoderStackEntry{
    sizes: map[string]uint32{},
  })
}

func (d *decoder) LookupSize(name string) (uint32, bool) {
  for i := len(d.stack) - 1; i >= 0; i-- {
    if v, ok := d.stack[i].sizes[name]; ok {
      return v, true
    }
  }
  return 0, false
}

func (d *decoder) setSize(name string, v uint32) {
  d.stack[len(d.stack) - 1].sizes[name] = v
}

func (d *decoder) Decode(rv reflect.Value, tag reflect.StructTag) error {
  if cr, ok := rv.Addr().Interface().(CustomReader); ok {
    return cr.Read(d.r, d)
  }

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
    listsize := tag.Get("listsize")
    if listsize != "" {
      d.setSize(listsize, v)
    }
    rv.SetUint(uint64(v))
    return nil
  case Uleb:
    v, err := readUleb(d.r)
    if err != nil {
      return err
    }
    listsize := tag.Get("listsize")
    if listsize != "" {
      d.setSize(listsize, v)
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
      listtag := tag.Get("listtag")
      if listtag == "" {
        return fmt.Errorf("missing listtag on slice")
      }
      size, ok := d.LookupSize(listtag)
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
    d.push()
    for i := 0; i < rv.NumField(); i++ {
      if err := d.Decode(rv.Field(i), rt.Field(i).Tag); err != nil {
        return err
      }
    }
    d.pop()
    return nil
  case reflect.Ptr:
    x, err := readUint32(d.r)
    if err != nil {
      return err
    }

    // If no index, this is an offset.
    if x == 0 {
      // nil pointer
      return nil
    }
    origOffset, _ := d.r.Seek(0, 1)
    if _, err := d.r.Seek(int64(x), 0); err != nil {
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
