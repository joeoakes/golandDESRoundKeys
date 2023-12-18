package main

import (
	"encoding/hex"
	"fmt"
)

// PC1 permutation table
var pc1 = [56]byte{
	56, 48, 40, 32, 24, 16, 8,
	0, 57, 49, 41, 33, 25, 17,
	9, 1, 58, 50, 42, 34, 26,
	18, 10, 2, 59, 51, 43, 35,
	62, 54, 46, 38, 30, 22, 14,
	6, 61, 53, 45, 37, 29, 21,
	13, 5, 60, 52, 44, 36, 28,
	20, 12, 4, 27, 19, 11, 3,
}

// PC2 permutation table
var pc2 = [48]byte{
	13, 16, 10, 23, 0, 4,
	2, 27, 14, 5, 20, 9,
	22, 18, 11, 3, 25, 7,
	15, 6, 26, 19, 12, 1,
	40, 51, 30, 36, 46, 54,
	29, 39, 50, 44, 32, 47,
	43, 48, 38, 55, 33, 52,
	45, 41, 49, 35, 28, 31,
}

// Left shift schedule
var shifts = [16]byte{
	1, 1, 2, 2, 2, 2, 2, 2,
	1, 2, 2, 2, 2, 2, 2, 1,
}

func leftCircularShift(slice []byte, shiftCount int) {
	length := len(slice)
	temp := make([]byte, length)
	copy(temp, slice)

	for i := 0; i < length; i++ {
		newIndex := (i + shiftCount) % length
		slice[i] = temp[newIndex]
	}
}

// Generate DES round keys from the original 64-bit key
func generateRoundKeys(key []byte) [16][48]byte {
	roundKeys := [16][48]byte{}

	// Apply the PC1 permutation to the original key
	var pc1Key [56]byte
	for i := 0; i < 56; i++ {
		pc1Key[i] = key[pc1[i]]
	}

	// Split the key into left and right halves
	left := pc1Key[:28]
	right := pc1Key[28:]

	// Generate round keys
	for i := 0; i < 16; i++ {
		// Perform left circular shift on left and right halves
		// Perform a left circular shift on the left and right halves
		leftCircularShift(left, i)
		leftCircularShift(right, i)
		//left = left[shifts[i]:] + left[:shifts[i]]
		//right = right[shifts[i]:] + right[:shifts[i]]

		// Combine left and right halves
		combined := append(left, right...)

		// Apply the PC2 permutation to get the round key
		for j := 0; j < 48; j++ {
			roundKeys[i][j] = combined[pc2[j]]
		}
	}

	return roundKeys
}

func main() {
	// 64-bit original key in hexadecimal format
	originalKeyHex := "133457799BBCDFF1"

	// Convert the original key from hexadecimal to bytes
	originalKey, _ := hex.DecodeString(originalKeyHex)

	// Generate the round keys
	roundKeys := generateRoundKeys(originalKey)

	// Print the round keys in hexadecimal format
	for i, key := range roundKeys {
		fmt.Printf("Round %2d: %s\n", i+1, hex.EncodeToString(key[:]))
	}
}
