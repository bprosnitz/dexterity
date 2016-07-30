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

type sizeList struct {
  A decode.Size `sizetag:"A"`
  B decode.Size `sizetag:"B"`
  AList []uint32 `sizetag:"A"`
  BList []uint32 `sizetag:"B"`
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
      final: &lebs{127, 126, -128},
    },
    {
      input: []byte{2, 0, 0, 0, 1, 0, 0, 0, 3, 0, 0, 0, 4, 0, 0, 0, 5, 0, 0, 0},
      empty: &sizeList{},
      final: &sizeList{2, 1, []uint32{3, 4}, []uint32{5}},
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
