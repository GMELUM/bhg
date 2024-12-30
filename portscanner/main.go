package main

import (

	_ "embed"

	"bufio"
	"flag"
	"fmt"
	"net"
	"os"
	"sync"
	"time"

	"encoding/json"

	"github.com/schollz/progressbar/v3"

)

//go:embed ports.json
var jsonData []byte

type Info struct {
	Description string `json:"description"`
	UDP         bool   `json:"udp"`
	Status      string `json:"status"`
	Port        string `json:"port"`
	TCP         bool   `json:"tcp"`
}

type Ports map[string][]Info

var (
	workers = flag.Int("w", 10, "Number of parallel workers")
	host    = flag.String("h", "127.0.0.1", "Host to scan")
	start   = flag.Int("s", 1, "Start port for scanning")
	end     = flag.Int("e", 65535, "End port for scanning")
	timeout = flag.Duration("t", 0, "Timeout to wait for a response from the server.")
	output  = flag.String("o", "", "File to save output")
)

func main() {

	defer func() {
		if r := recover(); r != nil {
			fmt.Printf("Program crashed with panic: %v\n", r)
		}

		fmt.Println("Press any key to exit...")
		reader := bufio.NewReader(os.Stdin)
		_, _ = reader.ReadString('\n')
	}() // Ensure graceful exit handling

	flag.Parse()

	// Validate port range
	if *start < 1 || *end > 65535 || *start > *end {
		fmt.Println("Invalid port range. Ports must be between 1 and 65535, and start <= end.")
		return
	}

	var portsList Ports
	err := json.Unmarshal(jsonData, &portsList)
	if err != nil {
		panic(err)
	}

	// Channel to distribute ports to workers
	ports := make(chan int, *workers)

	// Create progress bar
	bar := progressbar.NewOptions((*end)-(*start)+1,
		progressbar.OptionSetDescription("Scanning Ports"),
		progressbar.OptionShowCount(),
		progressbar.OptionEnableColorCodes(true),
		progressbar.OptionThrottle(time.Second/4), // Throttled updates
		progressbar.OptionSetTheme(progressbar.Theme{
			Saucer:        "█",
			SaucerHead:    "█",
			SaucerPadding: " ",
			BarStart:      "|",
			BarEnd:        "|",
		}),
	)

	// Thread-safe map to store open ports
	var openPorts sync.Map

	// File for real-time output if specified
	var file *os.File
	var writer *bufio.Writer
	if *output != "" {
		var err error
		file, err = os.Create(*output)
		if err != nil {
			fmt.Printf("\nError creating file %s: %v\n", *output, err)
			return
		}
		defer file.Close()
		writer = bufio.NewWriter(file)
	}

	// WaitGroup to wait for all workers to finish
	var wg sync.WaitGroup

	// Launch workers
	for i := 0; i < *workers; i++ {
		wg.Add(1)
		go func(
			host string,
			ports chan int,
			openPorts *sync.Map,
			wg *sync.WaitGroup,
			bar *progressbar.ProgressBar,
			writer *bufio.Writer,
		) {
			defer wg.Done()

			for port := range ports {

				bar.Add(1) // Update progress bar

				address := fmt.Sprintf("%s:%d", host, port)

				var (
					conn net.Conn
					err  error
				)

				if (*timeout) == 0 {
					// Use Dial to not avoid hanging indefinitely
					conn, err = net.Dial("tcp", address)
				} else {
					// Use DialTimeout to avoid hanging indefinitely
					conn, err = net.DialTimeout("tcp", address, *timeout)
				}

				if err == nil {
					conn.Close()
					openPorts.Store(port, true) // Store open port in a thread-safe map

					// Write open port to file
					if writer != nil {

						if info, ok := portsList[fmt.Sprint(port)]; ok {
							for _, elem := range info {
								fmt.Fprintf(writer, "%-9d | %v \n", port, elem.Description)
								writer.Flush()
							}

						} else {
							fmt.Fprintf(writer, "%-9d | %v \n", port, "Not specified")
							writer.Flush()
						}

					}
				}

			}
		}(*host, ports, &openPorts, &wg, bar, writer)
	}

	// Send port numbers to the channel
	for i := *start; i <= *end; i++ {
		ports <- i
	}
	close(ports) // Close the channel after sending all ports

	// Wait for all workers to finish
	wg.Wait()

	// Output results
	if *output == "" {
		// Print to console if no output file is specified
		fmt.Println("\nOpen Ports:")
		openPorts.Range(func(key, value interface{}) bool {
			fmt.Println(key) // Keys are the open ports
			return true
		})
	} else {
		fmt.Printf("\nResults saved to %s\n", *output)
	}
}
