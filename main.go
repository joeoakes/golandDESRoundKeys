package main

import (
	"fmt"
)

const (
	keySize     = 64
	halfKeySize = keySize / 2
	numRounds   = 16
)

// Initial Permutation Table
var initialPerm = [56]int{
	57, 49, 41, 33, 25, 17, 9,
	1, 58, 50, 42, 34, 26, 18,
	10, 2, 59, 51, 43, 35, 27,
	19, 11, 3, 60, 52, 44, 36,
	63, 55, 47, 39, 31, 23, 15,
	7, 62, 54, 46, 38, 30, 22,
	14, 6, 61, 53, 45, 37, 29,
	21, 13, 5, 28, 20, 12, 4,
}

// Number of bit shifts
var shiftTable = [16]int{
	1, 1, 2, 2,
	2, 2, 2, 2,
	1, 2, 2, 2,
	2, 2, 2, 1,
}

// Permutation applied on the shifted key to get the round key
var keyPerm = [48]int{
	14, 17, 11, 24, 1, 5,
	3, 28, 15, 6, 21, 10,
	23, 19, 12, 4, 26, 8,
	16, 7, 27, 20, 13, 2,
	41, 52, 31, 37, 47, 55,
	30, 40, 51, 45, 33, 48,
	44, 49, 39, 56, 34, 53,
	46, 42, 50, 36, 29, 32,
}

func main() {
	// Example 64-bit key (8 bytes)
	key := []byte("mysecret")

	// Generate round keys
	roundKeys, err := generateDESKeys(key)
	if err != nil {
		fmt.Println("Error generating keys:", err)
		return
	}

	// Print round keys
	for i, roundKey := range roundKeys {
		fmt.Printf("Round %d Key: %x\n", i+1, roundKey)
	}
}

func generateDESKeys(key []byte) ([][]byte, error) {
	if len(key) != 8 {
		return nil, fmt.Errorf("key must be 8 bytes long")
	}

	// Convert key to a 64-bit binary representation
	keyBits := bytesToBits(key)

	// Apply the initial key permutation
	keyPlus := permute(keyBits, initialPerm[:])

	// Split the key into two halves
	left := keyPlus[:halfKeySize/2]
	right := keyPlus[halfKeySize/2:]

	var roundKeys [][]byte
	for i := 0; i < numRounds; i++ {
		// Perform left shifts
		left = leftShift(left, shiftTable[i])
		right = leftShift(right, shiftTable[i])

		// Combine left and right halves and apply key permutation
		combined := append(left, right...)
		roundKey := permute(combined, keyPerm[:])
		roundKeys = append(roundKeys, bitsToBytes(roundKey))
	}

	return roundKeys, nil
}

func bytesToBits(data []byte) []int {
	bits := make([]int, 0, len(data)*8)
	for _, b := range data {
		for i := 7; i >= 0; i-- {
			bits = append(bits, int((b>>i)&1))
		}
	}
	return bits
}

func bitsToBytes(bits []int) []byte {
	bytes := make([]byte, 0, (len(bits)+7)/8)
	for i := 0; i < len(bits); i += 8 {
		var b byte
		for j := 0; j < 8; j++ {
			b |= byte(bits[i+j]) << (7 - j)
		}
		bytes = append(bytes, b)
	}
	return bytes
}

func permute(original []int, table []int) []int {
	permuted := make([]int, len(table))
	for i, pos := range table {
		permuted[i] = original[pos-1]
	}
	return permuted
}

func leftShift(data []int, shifts int) []int {
	return append(data[shifts:], data[:shifts]...)
}
