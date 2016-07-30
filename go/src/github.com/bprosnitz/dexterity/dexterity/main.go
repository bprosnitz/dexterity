package main

import (
  "flag"
  "log"
  "os"
  "fmt"
  "bprosnitz/dexterity"
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
    log.Fatal(err)
  }
  fmt.Printf("%#v\n", dex)
  f.Close()
}
