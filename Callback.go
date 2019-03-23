// package name: libHana
package main

/*
typedef unsigned int* (*GlobalRegistersCallbackType)();
unsigned int* callGlobalRegistersCallback(GlobalRegistersCallbackType callback) {
    return (callback)();
}
typedef unsigned char* (*MemoryReadCallbackType)(unsigned int address, unsigned int length);
unsigned char* callMemoryReadCallback(MemoryReadCallbackType callback, unsigned int address, unsigned int length) {
    return (callback)(address, length);
}
*/
import "C"
