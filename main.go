package main

import (
    "fmt"
    "io"
    "os"
    //"bufio"
    bitstream "github.com/dgryski/go-bitstream"
)

func check(e error) bool {
    if e == io.EOF { return true }
    if e != nil {
        panic(e)
    }
    return false
}

func readData(r io.Reader, cb func(uint32)) {
    var err error
    read24 := make([]uint8, 3)
    for {
        var bytesread int
        bytesread, err = r.Read(read24)
        if check(err) { break }
        data := uint32(0)
        switch bytesread {
           case 3:
               data = uint32(read24[2]) | uint32(read24[1]) << 8 | uint32(read24[0]) << 16
           case 2:
               data = uint32(read24[1]) | uint32(read24[0]) << 8
           case 1:
               data = uint32(read24[0])
           default:
               panic("I dont think this should happen")
        }
        cb(data)
    }
}

/*
func main() {
    fmt.Println("Hello there")
    r, err := os.Open("/tmp/2")
    defer r.Close()
    check(err)
    br := bitstream.NewReader(r)

    w, err := os.OpenFile("/tmp/3", os.O_RDWR|os.O_CREATE, 0755)
    defer w.Close()
    check(err)

    bytes := make([]uint8, 3)
    for size := 0; size < 2; size += 1 {
        data := uint32(0)
        p := uint8(Decode(br))
        c := uint16(0)
        i := uint16(0)
        if p == 21 {
            bits, err := br.ReadBits(16)
            check(err)
            i = uint16(bits)
        } else {
            c = uint16(Decode(CTrees[p], br))
            if p != 0 {
                i = uint16(Decode(ITrees[p], br))
            }
        }
        data = PCI2Data(p, c, i)
        pos := uint32(16)
        for i := uint32(0); i < 3; i += 1 {
            bytes[i] = uint8(data >> pos)
            pos -= 8
        }
        w.Write(bytes)
    }
}
*/

func main() {
    fmt.Println("Hello there")
    r, err := os.Open("/tmp/1")
    defer r.Close()
    check(err)
    pSymbols := map[uint16]uint64{}
    cSymbols := map[uint16]uint64{}
    iSymbols := map[uint16]uint64{}
    readData(r, func(data uint32) {
        p, c, i := Data2PCI(data)
        pSymbols[uint16(p)] += 1
        if p != 21 {
            cSymbols[c] += 1
        }
        if p != 0 {
            iSymbols[i] += 1
        }
    })
    generateCodesFromCount(pSymbols, cSymbols, iSymbols)
    fmt.Println(len(pSymbols), len(cSymbols), len(iSymbols))

    r.Seek(0, 0)
    w, err := os.OpenFile("/tmp/2", os.O_RDWR|os.O_CREATE, 0755)
    defer w.Close()
    check(err)
    bw := bitstream.NewWriter(w)

    readData(r, func(data uint32) {
        p, c, i := Data2PCI(data)
        Encode(pCodes, pBits, bw, uint16(p))
        if p != 21 && len(cCodes[p]) > 1 {
            Encode(cCodes[p], cBits[p], bw, c)
        }
        if p != 0 && len(iCodes[p]) > 1 {
            Encode(iCodes[p], iBits[p], bw, i)
        }
    })
    bw.Flush(bitstream.Bit(false))
}
