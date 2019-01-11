package main

import "sort"

/*
  These are the 22 permutation matches that 3-bytes (24-bits) can produce. Looking at them, they
  make sense. I came about these numbers through brute-force, but visually reviewing them shows
  that we can generate these matchings (in order).

  An example:
    Take a random number from the range 0 to 2^24 as "data". Represent it in 8 octal digits, pad
    leading 0's when needed.

    data = 0o63674163
    If we count the occurance of each octal digit, we get a table that looks like this:

      +-------+---+---+---+---+---+---+---+---+
      | Data  | 0 | 1 | 2 | 3 | 4 | 5 | 6 | 7 |
      +-------+---+---+---+---+---+---+---+---+
      | Count | 0 | 1 | 0 | 2 | 1 | 0 | 3 | 1 |
      +-------+---+---+---+---+---+---+---+---+

    Sorting that count map will produce a list:
      []uint8{3, 2, 1, 1, 1, 0, 0, 0}

    That list is equal to |permutations[15]|

    The 22 symbols we encode for the |p| in PCI abstractly represent these groupings.
*/
var permutations = [22][8]uint8 {
    [8]uint8{8, 0, 0, 0, 0, 0, 0, 0},
    [8]uint8{7, 1, 0, 0, 0, 0, 0, 0},
    [8]uint8{6, 2, 0, 0, 0, 0, 0, 0},
    [8]uint8{6, 1, 1, 0, 0, 0, 0, 0},
    [8]uint8{5, 3, 0, 0, 0, 0, 0, 0},
    [8]uint8{5, 2, 1, 0, 0, 0, 0, 0},
    [8]uint8{5, 1, 1, 1, 0, 0, 0, 0},
    [8]uint8{4, 4, 0, 0, 0, 0, 0, 0},
    [8]uint8{4, 3, 1, 0, 0, 0, 0, 0},
    [8]uint8{4, 2, 2, 0, 0, 0, 0, 0},
    [8]uint8{4, 2, 1, 1, 0, 0, 0, 0},
    [8]uint8{4, 1, 1, 1, 1, 0, 0, 0},
    [8]uint8{3, 3, 2, 0, 0, 0, 0, 0},
    [8]uint8{3, 3, 1, 1, 0, 0, 0, 0},
    [8]uint8{3, 2, 2, 1, 0, 0, 0, 0},
    [8]uint8{3, 2, 1, 1, 1, 0, 0, 0},
    [8]uint8{3, 1, 1, 1, 1, 1, 0, 0},
    [8]uint8{2, 2, 2, 2, 0, 0, 0, 0},
    [8]uint8{2, 2, 2, 1, 1, 0, 0, 0},
    [8]uint8{2, 2, 1, 1, 1, 1, 0, 0},
    [8]uint8{2, 1, 1, 1, 1, 1, 1, 0},
    [8]uint8{1, 1, 1, 1, 1, 1, 1, 1},
}

/*
  We only have 22 permutation groups. Each permutation group has a specific number of combinations
  associated with it. Visual means are most useful for me to explain again.

  Taking |permutations[0]| as a start:
    [8]uint8{8, 0, 0, 0, 0, 0, 0, 0}

  That string says that all 8 octal digits are duplicates. Valid matching 24 bit values are:
    0o00000000
    0o11111111
    0o22222222
    0o33333333
    0o44444444
    0o55555555
    0o66666666
    0o77777777

  There are only 8 combinations since there are only 8 possible octal values.
  The first permutation group is also special in that it has a single iteration. (It doesn't mater
  the order we return the octal digits in, all bits will match up)

  Taking |permutations[1]|:
    [8]uint8{7, 1, 0, 0, 0, 0, 0, 0}

  This reads as 7 octal digits are the same, and one is different from those. Some matching values:
    0o00000001
    0o00000002
    0o00000003
    0o00200000
    0o07000000
    0o11111110
    0o11111112
    0o77707777

  This group has lots of combinations and lots of iterations. The math used to calculate these
  numbers for each group is based on |permutations| and combinitorial math
*/
var combinations = [22]uint16 {
    8,      // 8!/(1!*(8-1)!)
    56,     // 7!/(1!*(7-1)!) * 8!/(1!*(8-1)!)
    56,     // 8!/(2!*(8-2)!) * 2
    168,    // 7!/(2!*(7-2)!) * 8!/(1!*(8-1)!)
    56,     // 8!/(2!*(8-2)!) * 2
    336,    // 6!/(1!*(6-1)!) * 8!/(2!*(8-2)!) * 2
    280,    // 7!/(3!*(7-3)!) * 8!/(1!*(8-1)!)
    28,     // 8!/(2!*(8-2)!)
    336,    // 6!/(1!*(6-1)!) * 8!/(2!*(8-2)!) * 2
    168,    // 8!/(3!*(8-3)!) * 3
    840,    // 6!/(2!*(6-2)!) * 8!/(2!*(8-2)!) * 2
    280,    // 7!/(4!*(7-4)!) * 8!/(1!*(8-1)!)
    168,    // 8!/(3!*(8-3)!) * 3
    420,    // 6!/(2!*(6-2)!) * 8!/(2!*(8-2)!)
    840,    // 5!/(1!*(5-1)!) * 8!/(3!*(8-3)!) * 3
    1120,   // 6!/(3!*(6-3)!) * 8!/(2!*(8-2)!) * 2
    168,    // 7!/(5!*(7-5)!) * 8!/(1!*(8-1)!)
    70,     // 8!/(4!*(8-4)!)
    560,    // 5!/(2!*(5-2)!) * 8!/(3!*(8-3)!)
    420,    // 6!/(4!*(6-4)!) * 8!/(2!*(8-2)!)
    56,     // 7!/(6!*(7-6)!) * 8!/(1!*(8-1)!)
    1,      // 8!/(8!*(8-8)!)
}

/*
  The iterations are calulated by some simple psuedo code based on |permutations|:

    total = 0
    for i in permutations[p] {
        total *= permutations[i]!
    }
    iterations[p] = 8! / total
*/
var iterations = [22]uint16 {
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

/*
  Offsets are only used when converting from the PCI rank to the original data.
  When we populate |pci|, we just stack all the data from the sorted PCI ranks into one list and
  use this table as a multiplier to grab the location of our original data.
  It is essentially a big brute-forced cheatsheet of answers for the program.

  I have been using the list visually to help with the code, but the list could be generated like:
    offsets := [22]uint32{}
    offsets[0] := 0
    for i := 1; i < 22; i += 1 {
      offsets[i] = combinations[i-1] * iterations[i-1]
    }
*/
var offsets = [22]uint32 {
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
/*
  NOTE:
  I believe it should be possible to get from the PCI rank to the original data *without* this
  brute-forced cheat sheet by plugging in the appropriate symbols from the combinations against the
  permutations table given it is an ordered list. The iterations might require actual iteration
  through a shuffle. Both of those trigger my "this will probably run slower" alarms.

  Nevertheless, I keep exploring this in my head. It would would allow all of the memory we use to,
  optionally, be dropped. Thats pretty interesting as a mechanism to me.
*/

/*
  These count tables are not currently used, but are exported. Ideally, I will be able to implement
  an arithmetic coder for these codes instead of the huffman trees. That is the dream.

  Each time a specific permutation is found, the count for that permutation is increased. The same
  for the combination and iterations.

  If you plan to use them, these do not account for the two outer edges of the |permutations|.
  The two cases are:
    |permutations[0]| has a single iteration. It has 8 combinations
    |permutations[21]| has a single combination. It has 40320 iterations

  With that in mind, a model based on these counts might want to remove those counters from the
  totals.
*/
var PermCount = [22]uint32{}
var CombCount = [1120]uint32{}
var IterCount = [40320]uint32{}

/*
  These lists have a few things going on. They are used when converting to and from our PCI ranks.
  When we read in our data, we are reading 3 bytes -- 24 bits of data. We shift that to a uint32 to
  make it easier to lookup in memory. The value of the data is the index in these lists.

  Helper functions are exported for usability as PCI2Data() and Data2PCI() which shows the usage of
  these lists.
*/
var plist = [16777216]uint8{}
var clist = [16777216]uint16{}
var ilist = [16777216]uint16{}
var pci = [16777216]uint32{}

func Data2PCI(d uint32) (uint8, uint16, uint16){
    return plist[d], clist[d], ilist[d]
}

func PCI2Data(p uint8, c uint16, i uint16) uint32 {
    return pci[offsets[p] + (uint32(c) * uint32(iterations[p])) + uint32(i)]
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

func comparePermutation(a, b [8]uint8) bool {
    // return true if a == b else false
    for i := range a {
        if a[i] != b[i] {
            return false
        }
    }
    return true
}

func unpack(s uint32) [8]uint8 {
    /*
      Take an uint32 and grab 24 bits starting from LSB. The example bitshifting that goes on.
      Ex: Data == uint32(1382749)
      s = 00010101 00011001 01011101
      for {
        // Last byte of |s|
        e = 01011101
        AND 00000111
        e = 00000101
        // using |e| as the index, iterate the occurance of the value
        // s is now bitshifted 3 positions to the right to remove the value and the loop repeats
        // s = 000101010001100101011101 -- in first iteration
        s = 000000101010001100101011
      }
    */
    ret := [8]uint8{}
    for i := uint8(0); i < 8; i += 1 {
        e := uint8(s)
        e &= 0x7
        ret[e]++
        s >>= 3
    }
    // We sort the values we found so that they are ordered and can be comared against |permutations|
    sort.Slice(ret[:], func(i, j int) bool { return ret[i] > ret[j] })
    return ret
}

func whichPermutationGroup(data uint32) uint8 {
    dups := unpack(data)
    // All 24 bit numbers fall into one of the 22 permutation descriptors
    // Anything outside of that is un-possible
    for i, v := range permutations {
        if comparePermutation(dups, v) {
            return uint8(i)
        }
    }
    panic("We did not find a match for our permutation group, something is terribly wrong")
    return 0 // I am a terrible programmer and go is not my... go-to language. Sam make compile. Sam sleep
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
        size := uint32(iterations[permutation])

        combination := elements / size
        iteration := elements % size

        // This counts as a safety checks right?
        // permutations groups are always 22, but this shows whats happening clearer
        if permutation >= uint8(len(permutations)) {
            panic("Permutation is borked")
        }
        if combination >= uint32(combinations[permutation]) {
            panic("Combination is borked")
        }
        if iteration >= uint32(iterations[permutation]) {
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
