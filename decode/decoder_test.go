package decode_test

import (
  "testing"
  "reflect"
  "bytes"
  "github.com/bprosnitz/dexterity/decode"
)

type uint32s struct {
  A uint32
  B uint32
  C uint32
}

type lebs struct {
  A decode.Uleb
  B decode.Ulebp1
  C decode.Sleb
}

type list struct {
  A uint32 `listsize:"A"`
  B uint32 `listsize:"B"`
  AList []uint32 `listtag:"A"`
  BList []uint32 `listtag:"B"`
}

type mutf8 struct {
  A string
  B string
}

type ptr struct {
  A uint32
  P *ptrElem
  B uint32
  NilPtr *ptrElem
}

type ptrElem struct {
  A uint32
}

type array struct {
  A [2]uint32
  B [3]uint8
}

func TestDecode(t *testing.T) {
  tests := []struct{
    input []byte
    empty interface{}
    final interface{}
  }{
    {
      input: []byte{1, 0, 0, 0, 2, 0, 0, 0, 3, 0, 0, 0},
      empty: &uint32s{},
      final: &uint32s{1,2,3},
    },
    {
      input: []byte{0x80, 0x7f, 0x80, 0x7f, 0x80, 0x7f},
      empty: &lebs{},
      final: &lebs{16256, 16255, -128},
    },
    {
      input: []byte{0x61, 0x62, 0x63, 0x00, 0x61, 0x62, 0x00},
      empty: &mutf8{},
      final: &mutf8{"abc", "ab"},
    },
    {
      input: []byte{2, 0, 0, 0, 1, 0, 0, 0, 3, 0, 0, 0, 4, 0, 0, 0, 5, 0, 0, 0},
      empty: &list{},
      final: &list{2, 1, []uint32{3, 4}, []uint32{5}},
    },
    {
      input: []byte{1, 0, 0, 0, 17, 0, 0, 0, 3, 0, 0, 0, 0, 0, 0, 0, 0, 2, 0, 0, 0},
      empty: &ptr{},
      final: &ptr{1, &ptrElem{2}, 3, nil},
    },
    {
      input: []byte{1, 0, 0, 0, 2, 0, 0, 0, 3, 4, 5},
      empty: &array{},
      final: &array{[2]uint32{1,2},[3]uint8{3,4,5}},
    },
  }
  for _, test := range tests {
    if err := decode.Decode(bytes.NewReader(test.input), test.empty); err != nil {
      t.Errorf("%#v: %v", test.input, err)
      continue
    }
    if !reflect.DeepEqual(test.empty, test.final) {
      t.Errorf("%#v: got %#v, want %#v", test.input, test.empty, test.final)
    }
  }
}
