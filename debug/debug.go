package debug

import (
	"fmt"
	"os"
)

func Info() {
	fmt.Println("=== debug.Info <begin> =========")
	dir, err := os.Getwd()
	if err != nil {
		fmt.Println("debug.Info | error:", err)
	}
	fmt.Println("Current Dir : ", dir)
	fmt.Println("=== debug.Info <end> ===========")
}
