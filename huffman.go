package main

import (
    "fmt"
    "sort"
    huffman "github.com/icza/huffman"
    bitstream "github.com/dgryski/go-bitstream"
)

// I use uint16 throughout all the tree even though permutations needs only uint8. Ease of reading
type ValueType uint16

// Codes are the huffman codes store in a uint64, starting with MSB
type Codes []uint64
// Bits are paired with Codes and specify the number of bits to read out of Codes, starting with MSB
type Bits []uint8

// Only used for encoding process
var pCodes Codes
var cCodes [22] Codes
var iCodes [22] Codes
var pBits Bits
var cBits [22] Bits
var iBits [22] Bits

// Used in decoding process to walk the trees one bit at a time (also for generating the intial codes)
var pTree *huffman.Node
var cTrees [22]*huffman.Node
var iTrees [22]*huffman.Node

func Encode(codes Codes, bits Bits, bw *bitstream.BitWriter, v uint16) {
    // Write out |bits[v]| number of bits from |codes[v]| ; starting with MSB
    bw.WriteBits(codes[v], int(bits[v]))
}

func Decode(node *huffman.Node, br *bitstream.BitReader) uint16 {
    for {
        bit, _ := br.ReadBit()
        if bit {
            node = node.Right
        } else {
            node = node.Left
        }
        if node != nil {
            return uint16(node.Value)
        }
    }
    return 0xFFFF
}

func newHuffTree(counts []uint64) *huffman.Node {
    leaves := []*huffman.Node{}
    for i, v := range counts {
        leaves = append(leaves, &huffman.Node{Value: huffman.ValueType(i), Count: int(v)})
    }
    return huffman.Build(leaves)
}

func populateCodes(root *huffman.Node, codes Codes, bits Bits) {
    // Used for building a list of all huffman codes for easy lookup when encoding
    var traverse func(n *huffman.Node, code uint64, bits byte)
    traverse = func(n *huffman.Node, code uint64, count byte) {
        if n.Left == nil {
            v := uint16(n.Value)
            codes[v] = code
            bits[v] = count
            return
        }
        count++
        traverse(n.Left, code<<1, count)
        traverse(n.Right, code<<1+1, count)
    }
    traverse(root, 0, 0)
}

func convertCountToLeaves(count map[uint16]uint64) []*huffman.Node {
    leaves := make([]*huffman.Node, len(count))
    i := 0
    for k, v := range count {
        leaves[i] = &huffman.Node{Value: huffman.ValueType(k), Count: int(v)}
        i++
    }
    sort.Stable(huffman.SortNodes(leaves))
    return leaves
}

func generateCodesFromCount(pcount map[uint16]uint64, ccount map[uint16]uint64, icount map[uint16]uint64) {
    // Return leaves from counts, sorted in ascending order
    pleaves := convertCountToLeaves(pcount)
    cleaves := convertCountToLeaves(ccount)
    ileaves := convertCountToLeaves(icount)

    // Build the huffman trees for the permutations and populate the codes for lookup during encode
    pTree = huffman.BuildSorted(pleaves)
    populateCodes(pTree, pCodes, pBits)
    for i := 0; i < 22; i += 1 {
        if i != 21 { // |permutations[21]| has a single combination; no information to encode!
            newComb := []*huffman.Node{}
            for _, v := range cleaves {
                if uint16(v.Value) < combinations[i] {
                    newComb = append(newComb, v)
                }
            }
            if (len(newComb) < 2) == false {
                cTrees[i] = huffman.BuildSorted(newComb)
                populateCodes(cTrees[i], cCodes[i], cBits[i])
            }
        }
        if i != 0 { // |permutations[0]| has a single iteration per combination; no information to encode!
            newIter := []*huffman.Node{}
            for _, v := range ileaves {
                if uint16(v.Value) < iterations[i] {
                    newIter = append(newIter, v)
                }
            }
            if (len(newIter) < 2) == false {
                iTrees[i] = huffman.BuildSorted(newIter)
                populateCodes(iTrees[i], iCodes[i], iBits[i])
            }
        }
    }
}

func init() {
    fmt.Println("Say something!")

    /*
      Allocate the space for all the trees and codes. Some models will not use any of these, this
      fact does not change the convience of having a statically sized list.

      The encoding and decoding processes should not be requesting indexes which were never
      populated. I would prefer to *enforce* this in code, but there are enough panic()'s spread
      around to warm my heart.
    */
    pCodes = make(Codes, 22)
    pBits = make(Bits, 22)
    for i := 0; i < 22; i += 1 {
        cCodes[i] = make(Codes, combinations[i])
        cBits[i] = make(Bits, combinations[i])
        iCodes[i] = make(Codes, iterations[i])
        iBits[i] = make(Bits, iterations[i])
    }
}
