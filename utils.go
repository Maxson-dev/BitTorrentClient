package main

import "unsafe"

func b2s(buf []byte) string {
	return *(*string)(unsafe.Pointer(&buf))
}
