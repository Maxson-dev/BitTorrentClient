package main

import (
	"fmt"
	"math/rand"
	"testing"
	"time"
)

func TestBitField_SetPiece(t *testing.T) {
	rand.Seed(time.Now().UnixMicro())

	const sample = 10000

	var bitField BitField = make([]byte, sample+400)
	for i := 0; i < sample; i++ {
		bitField[i] = byte(rand.Uint32())
	}

	var strRepr1 string
	for i := 0; i < len(bitField); i++ {
		strRepr1 += fmt.Sprintf("%08b-", bitField[i])
	}
	fmt.Println(strRepr1)

	var bits = make([]int, sample)
	for i := 0; i < sample; i++ {
		bits[i] = rand.Intn(sample)
	}

	for i := 0; i < len(bits); i++ {
		bitField.SetPiece(bits[i])
		if !bitField.HasPiece(bits[i]) {
			t.Errorf("SetPiece(%d) didn`t set bit", bits[i])
		}
	}

	var strRepr2 string
	for i := 0; i < len(bitField); i++ {
		strRepr2 += fmt.Sprintf("%08b-", bitField[i])
	}
	fmt.Println(strRepr2)

}

func TestBitField_HasPiece(t *testing.T) {
	var bitField BitField = []byte{8, 9, 4}
	var strRepr1 string
	for i := 0; i < len(bitField); i++ {
		strRepr1 += fmt.Sprintf("%08b-", bitField[i])
	}
	fmt.Println(strRepr1)
	tests := []int{4, 21}
	wants := []bool{true, true}
	for i, test := range tests {
		result := bitField.HasPiece(test)
		if result != wants[i] {
			t.Errorf("HasPiece(%d) == %t", test, result)
		}
	}
}
