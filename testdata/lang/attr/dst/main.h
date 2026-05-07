#pragma once
#include "so/builtin/builtin.h"

// -- Types --

typedef struct main_Exported main_Exported;

// Exported struct with so:attr.
//
typedef struct __attribute__((packed)) main_Exported {
    so_int v;
} main_Exported;

// -- Variables and constants --

// Exported volatile variable.
//
extern volatile so_int main_Counter;

// Exported thread-local variable.
//
extern _Thread_local so_int main_PerThread;

// -- Functions and methods --

// Exported function with so:attr.
//
__attribute__((noinline)) void main_Work(void);
