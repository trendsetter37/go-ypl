package main

import "sort"

// These numbers can be gotten with math, but it is easier to staticly encode them
// All ordering is vital
// Magic combinatorial group descriptors
var Groups = [22][8]uint8 {
    [8]uint8{0, 0, 0, 0, 0, 0, 0, 8},
    [8]uint8{0, 0, 0, 0, 0, 0, 1, 7},
    [8]uint8{0, 0, 0, 0, 0, 0, 2, 6},
    [8]uint8{0, 0, 0, 0, 0, 1, 1, 6},
    [8]uint8{0, 0, 0, 0, 0, 0, 3, 5},
    [8]uint8{0, 0, 0, 0, 0, 1, 2, 5},
    [8]uint8{0, 0, 0, 0, 1, 1, 1, 5},
    [8]uint8{0, 0, 0, 0, 0, 0, 4, 4},
    [8]uint8{0, 0, 0, 0, 0, 1, 3, 4},
    [8]uint8{0, 0, 0, 0, 0, 2, 2, 4},
    [8]uint8{0, 0, 0, 0, 1, 1, 2, 4},
    [8]uint8{0, 0, 0, 1, 1, 1, 1, 4},
    [8]uint8{0, 0, 0, 0, 0, 2, 3, 3},
    [8]uint8{0, 0, 0, 0, 1, 1, 3, 3},
    [8]uint8{0, 0, 0, 0, 1, 2, 2, 3},
    [8]uint8{0, 0, 0, 1, 1, 1, 2, 3},
    [8]uint8{0, 0, 1, 1, 1, 1, 1, 3},
    [8]uint8{0, 0, 0, 0, 2, 2, 2, 2},
    [8]uint8{0, 0, 0, 1, 1, 2, 2, 2},
    [8]uint8{0, 0, 1, 1, 1, 1, 2, 2},
    [8]uint8{0, 1, 1, 1, 1, 1, 1, 2},
    [8]uint8{1, 1, 1, 1, 1, 1, 1, 1},
}
var Permutations = [22]uint16 {
    1,
    8,
    28,
    56,
    56,
    168,
    336,
    70,
    280,
    420,
    840,
    1680,
    560,
    1120,
    1680,
    3360,
    6720,
    2520,
    5040,
    10080,
    20160,
    40320,
}
var Combinations = [22]uint16 {
    8,
    56,
    56,
    168,
    56,
    336,
    280,
    28,
    336,
    168,
    840,
    280,
    168,
    420,
    840,
    1120,
    168,
    70,
    560,
    420,
    56,
    1,
}
var Offsets = [22]uint32 {
    0,
    8,
    456,
    2024,
    11432,
    14568,
    71016,
    165096,
    167056,
    261136,
    331696,
    1037296,
    1507696,
    1601776,
    2072176,
    3483376,
    7246576,
    8375536,
    8551936,
    11374336,
    15607936,
    16736896,
}

var PermCount = [22]uint32{}
var CombCount = [1120]uint32{}
var IterCount = [40320]uint32{}

var plist = [16777216]uint8{}
var clist = [16777216]uint16{}
var ilist = [16777216]uint16{}
var pci = [16777216]uint32{}

func compareIntSlice(a, b [8]uint8) bool {
    for i := range a {
        if a[i] != b[i] {
            return false
        }
    }
    return true
}

func unpack(s uint32) [8]uint8 {
    ret := [8]uint8{}
    for i := uint8(0); i < 8; i += 1 {
        e := uint8(s)
        e &= 0x7
        ret[e]++
        s >>= 3
    }
    sort.Slice(ret[:], func(i, j int) bool { return ret[i] < ret[j] })
    return ret
}

func whichPermutationGroup(data uint32) uint8 {
    dups := unpack(data)
    // All 24 bit numbers fall into one of the 22 permutation descriptors
    // Anything outside of that is un-possible
    for i, v := range Groups {
        if compareIntSlice(dups, v) {
            return uint8(i)
        }
    }
    panic("We did not find a match for our permutation group, something is terribly wrong")
    return 0 // I am a terrible programmer and go is not my... go-to language. Sam make compile. Sam sleep
}

func Data2PCI(d uint32) (uint8, uint16, uint16){
    return plist[d], clist[d], ilist[d]
}

func PCI2Data(p uint8, c uint16, i uint16) uint32 {
    return pci[Offsets[p] + (uint32(c) * uint32(Permutations[p])) + uint32(i)]
}

func ValidateIntegrity() {
    /*
      The sole purpose of this function is to self-check our tables match up.
      I am hoping this code doesn't break on different endian'd systems, if it does
      this function should catch it immediately
    */
    for data := uint32(0); data < 16777216; data += 1 {
        p, c, i := Data2PCI(data)
        pciData := PCI2Data(p, c, i)
        if pciData != data {
            panic("Our plist, clist, ilist info doesn't match our pci info")
        }
    }
}

func cial() {
    // A temporary list of list of ints by permutation
    var pints [22][]uint32
    for i := 0; i < 22; i += 1 {
        pints[i] = []uint32{}
    }

    // The real work is done here. Looping through all possible 24-bit values (our "data") and
    // sorting them into permutation groups while recording thier combination and iteration
    for data := uint32(0); data < 16777216; data += 1 {
        permutation := whichPermutationGroup(data)
        elements := uint32(len(pints[permutation]))
        size := uint32(Permutations[permutation])

        combination := elements / size
        iteration := elements % size

        // This counts as a safety checks right?
        if permutation >= uint8(len(Permutations)) {
            panic("Permutation is borked")
        }
        if combination >= uint32(Combinations[permutation]) {
            panic("Combination is borked")
        }
        if iteration >= uint32(Permutations[permutation]) {
            panic("Iteration is borked")
        }

        PermCount[permutation] += 1
        CombCount[combination] += 1
        IterCount[iteration] += 1
        plist[data] = permutation
        clist[data] = uint16(combination)
        ilist[data] = uint16(iteration)
        pints[permutation] = append(pints[permutation], data)
    }

    // Populate all 24-bit values sorted by thier PCI ranks
    i := 0
    for _, p := range pints {
        for _, v := range p {
            pci[i] = v
            i += 1
        }
    }
    ValidateIntegrity()
}

func init() {
    cial()
}
