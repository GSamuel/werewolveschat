package chat

import (
	"fmt"
	"os"
	"os/exec"
	"runtime"
)

func Clear() {

	switch runtime.GOOS {
	case "windows":
		cmd := exec.Command("cmd", "/c", "cls")
		cmd.Stdout = os.Stdout
		cmd.Run()
		return
	case "linux":
		cmd := exec.Command("clear") //Linux example, its tested
		cmd.Stdout = os.Stdout
		cmd.Run()
	default:
		panic("Your platform is unsupported! I can't clear terminal screen :(")
	}

	fmt.Println(runtime.GOOS)
}
