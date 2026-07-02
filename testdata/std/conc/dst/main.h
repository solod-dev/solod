#pragma once
#include "so/builtin/builtin.h"
#include "so/conc/conc.h"
#include "so/errors/errors.h"
#include "so/mem/mem.h"
#include "so/sync/sync.h"
#include "so/time/time.h"

// -- Types --

typedef struct main_Task main_Task;

// Task carries one task's input, output and error through a *Task.
typedef struct main_Task {
    so_int in;
    so_int out;
    so_Error err;
} main_Task;
