// Errorf formats its arguments like fmt.Sprintf,
// logs the result, and marks the test failed.
void testing_T_Errorf(void* self, const char* format, ...);

// Fatalf is like Errorf but marks the test as fatally failed.
// The test function must return right after calling it.
void testing_T_Fatalf(void* self, const char* format, ...);
