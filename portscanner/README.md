# Port Scanner

This is a lightweight and efficient TCP port scanner written in Go. The tool allows scanning a range of ports on a specified host, leveraging concurrency for high performance. It also includes features like progress tracking, timeout configuration, and result output to a file.

---

## Features

- **Concurrent Scanning:** Utilizes multiple workers to scan ports in parallel.
- **Progress Bar:** Displays a real-time progress bar during the scanning process.
- **Timeout Control:** Customizable timeout to avoid hanging on unresponsive ports.
- **Output to File:** Option to save results to a file.

---

## Prerequisites

- Go 1.19 or later installed on your system.

---

## Installation

Clone the repository and navigate to the project directory:

```bash
git clone https://github.com/GMELUM/bhg
cd portscanner
```

Build the project:

```bash
go build -o portscanner
```

---

## Usage

Run the program with the following command-line arguments:

```bash
./portscanner [flags]
```

### Flags

| Flag | Default Value | Description                                                              |
| ---- | ------------- | ------------------------------------------------------------------------ |
| `-h` | `127.0.0.1`   | Host to scan.                                                            |
| `-s` | `1`           | Start port for scanning.                                                 |
| `-e` | `65535`       | End port for scanning.                                                   |
| `-w` | `10`          | Number of parallel workers for concurrent scanning.                      |
| `-t` | `0`           | Timeout for a response from each port.                                   |
| `-o` | `""`          | File to save the output. If omitted, results are printed to the console. |

---

## Examples

### Scan ports on a host

```bash
./portscanner -h 192.168.1.1 -s 1 -e 1000
```

This will scan ports 1 to 1000 on `192.168.1.1` and display the results in the console.

### Save output to a file

```bash
./portscanner -h 192.168.1.1 -s 1 -e 1000 -o results.txt
```

This will scan ports 1 to 1000 on `192.168.1.1` and save the results to `results.txt`.

### Set a timeout

```bash
./portscanner -h 192.168.1.1 -t 2s
```

This sets a timeout of 2 seconds for each port scan attempt.

---

## Output

- **Console Output:** Open ports will be printed line by line.
- **File Output:** If the `-o` flag is used, results will be written to the specified file.

---

## License

This project is licensed under the MIT License. See the LICENSE file for details.