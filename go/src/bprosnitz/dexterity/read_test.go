package dexterity

import (
  "testing"
  "bytes"
)

func TestReadUleb(t *testing.T) {
    testCases := []struct{
      input []byte
      output uint32
      outputErr bool
    }{
      {
        input: []byte{},
        outputErr: true,
      },
      {
        input: []byte{0x00},
        output: 0,
      },
      {
        input: []byte{0x7f},
        output: 127,
      },
      {
        input: []byte{0x80},
        outputErr: true,
      },
      {
        input: []byte{0x80, 0x00},
        output: 0,
      },
      {
        input: []byte{0x81, 0x00},
        output: 1,
      },
      {
        input: []byte{0x80, 0x01},
        output: 128,
      },
      {
        input: []byte{0x80, 0x7f},
        output: 16256,
      },
      {
        input: []byte{0xa9, 0xc4, 0x84, 0x98, 0x01},
        output: 318841385,
      },
    }

    for _, testCase := range testCases {
      out, err := readUleb(bytes.NewReader(testCase.input))
      if testCase.outputErr {
        if err == nil {
          t.Errorf("%#v: expected error, got none", testCase.input)
        }
      } else {
        if err != nil {
          t.Errorf("%#v: unexpected error: %v", testCase.input, err)
          continue
        }
        if out != testCase.output {
          t.Errorf("%#v: got %v, want %v", testCase.input, out, testCase.output)
        }
      }
    }
}

func TestReadUlebP1(t *testing.T) {
    testCases := []struct{
      input []byte
      output uint32
      outputErr bool
    }{
      {
        input: []byte{},
        outputErr: true,
      },
      {
        input: []byte{0x00},
        output: 0xffffffff,
      },
      {
        input: []byte{0x7f},
        output: 126,
      },
      {
        input: []byte{0x80},
        outputErr: true,
      },
      {
        input: []byte{0x80, 0x00},
        output: 0xffffffff,
      },
      {
        input: []byte{0x81, 0x00},
        output: 0,
      },
      {
        input: []byte{0x80, 0x01},
        output: 127,
      },
      {
        input: []byte{0x80, 0x7f},
        output: 16255,
      },
      {
        input: []byte{0xa9, 0xc4, 0x84, 0x98, 0x01},
        output: 318841384,
      },
    }

    for _, testCase := range testCases {
      out, err := readUlebP1(bytes.NewReader(testCase.input))
      if testCase.outputErr {
        if err == nil {
          t.Errorf("%#v: expected error, got none", testCase.input)
        }
      } else {
        if err != nil {
          t.Errorf("%#v: unexpected error: %v", testCase.input, err)
          continue
        }
        if out != testCase.output {
          t.Errorf("%#v: got %v, want %v", testCase.input, out, testCase.output)
        }
      }
    }
}

func TestReadSleb(t *testing.T) {
    testCases := []struct{
      input []byte
      output int32
      outputErr bool
    }{
      {
        input: []byte{},
        outputErr: true,
      },
      {
        input: []byte{0x00},
        output: 0,
      },
      {
        input: []byte{0x7f},
        output: -1,
      },
      {
        input: []byte{0x80},
        outputErr: true,
      },
      {
        input: []byte{0x80, 0x00},
        output: 0,
      },
      {
        input: []byte{0x81, 0x00},
        output: 1,
      },
      {
        input: []byte{0x80, 0x01},
        output: 128,
      },
      {
        input: []byte{0x80, 0x7f},
        output: -128,
      },
      {
        input: []byte{0xa9, 0xc4, 0x84, 0x98, 0x01},
        output: 318841385,
      },
    }

    for _, testCase := range testCases {
      out, err := readSleb(bytes.NewReader(testCase.input))
      if testCase.outputErr {
        if err == nil {
          t.Errorf("%#v: expected error, got none", testCase.input)
        }
      } else {
        if err != nil {
          t.Errorf("%#v: unexpected error: %v", testCase.input, err)
          continue
        }
        if out != testCase.output {
          t.Errorf("%#v: got %v, want %v", testCase.input, out, testCase.output)
        }
      }
    }
}
