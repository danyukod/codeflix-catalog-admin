package main

import "testing"

func TestMainFunction(t *testing.T) {
	// Since main() calls fmt.Println and doesn't return anything,
	// we test that it runs without panicking.
	defer func() {
		if r := recover(); r != nil {
			t.Errorf("main() panicked: %v", r)
		}
	}()

	main()
}
