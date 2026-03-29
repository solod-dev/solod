package main

import "solod.dev/so/os"

func envTest() {
	{
		// Setenv, Getenv.
		err := os.Setenv("SO_TEST_KEY", "test_value")
		if err != nil {
			panic("Setenv failed")
		}
		val := os.Getenv("SO_TEST_KEY")
		if val != "test_value" {
			panic("Getenv: wrong value")
		}
	}
	{
		// LookupEnv - present.
		os.Setenv("SO_TEST_LOOKUP", "found")
		val, ok := os.LookupEnv("SO_TEST_LOOKUP")
		if !ok {
			panic("LookupEnv: should be present")
		}
		if val != "found" {
			panic("LookupEnv: wrong value")
		}
	}
	{
		// LookupEnv - absent.
		_, ok := os.LookupEnv("SO_TEST_NONEXISTENT_VAR_XYZ")
		if ok {
			panic("LookupEnv: should not be present")
		}
	}
	{
		// Unsetenv.
		os.Setenv("SO_TEST_UNSET", "bye")
		err := os.Unsetenv("SO_TEST_UNSET")
		if err != nil {
			panic("Unsetenv failed")
		}
		val := os.Getenv("SO_TEST_UNSET")
		if val != "" {
			panic("Unsetenv: should be empty")
		}
	}
	{
		// Getenv on PATH (should always be set).
		path := os.Getenv("PATH")
		if len(path) == 0 {
			panic("Getenv PATH: empty")
		}
	}
}
