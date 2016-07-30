package main

import (
  "flag"
  "log"
  "os"
  "github.com/bprosnitz/dexterity"
   "github.com/davecgh/go-spew/spew"
)

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
  dex := &dexterity.Dex{}
  if err := dexterity.ReadDex(f, dex); err != nil {
//  if err := dexterity.Process(f, 0, dex); err != nil {
    log.Fatal(err)
  }
  spew.Printf("%#v\n", dex)
  f.Close()
}
