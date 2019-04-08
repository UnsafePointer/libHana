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
typedef void (*ContinueCallbackType)();
void callContinueCallback(ContinueCallbackType callback) {
    return (callback)();
}
typedef void (*AddLoadWatchpointCallbackType)(unsigned int address);
void callAddLoadWatchpointCallback(AddLoadWatchpointCallbackType callback, unsigned int address) {
    return (callback)(address);
}
typedef void (*AddStoreWatchpointCallbackType)(unsigned int address);
void callAddStoreWatchpointCallback(AddStoreWatchpointCallbackType callback, unsigned int address) {
    return (callback)(address);
}
typedef void (*RemoveBreakpointCallbackType)(unsigned int address);
void callRemoveBreakpointCallback(RemoveBreakpointCallbackType callback, unsigned int address) {
    return (callback)(address);
}
typedef void (*RemoveLoadWatchpointCallbackType)(unsigned int address);
void callRemoveLoadWatchpointCallback(RemoveLoadWatchpointCallbackType callback, unsigned int address) {
    return (callback)(address);
}
typedef void (*RemoveStoreWatchpointCallbackType)(unsigned int address);
void callRemoveStoreWatchpointCallback(RemoveStoreWatchpointCallbackType callback, unsigned int address) {
    return (callback)(address);
}
*/
import "C"
