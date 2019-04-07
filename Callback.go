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
typedef void (*AddBreakpointCallbackType)(unsigned int address);
void callAddBreakpointCallback(AddBreakpointCallbackType callback, unsigned int address) {
    return (callback)(address);
}
*/
import "C"
