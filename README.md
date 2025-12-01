# zymkey.go

A Go wrapper for the Zymkey `libzk_app_utils` C library, allowing you to interact with Zymbit security modules (Zymkey 4i, HSM4, HSM6, etc.) directly from Go applications.

## Prerequisites

*   **Hardware**: A Zymbit security module (e.g., Zymkey 4i, HSM4, HSM6) installed on your device (e.g., Raspberry Pi).
*   **Software**: The Zymbit software drivers and C libraries must be installed.
    
    For standard installations:
    ```bash
    curl -G https://s3.amazonaws.com/zk-sw-repo/install_zk_sw.sh | sudo bash
    ```

    For newest Raspberry Pi OS or Debian 13+:
    ```bash
    curl -G https://gist.githubusercontent.com/snowmerak/4649cb44b092fa0c3d3dc53308766ce3/raw/1f00513a4341ca92d0ac0ba5a9eaa0d083d00629/install_zk_sw.sh | sudo bash
    ```

## Installation

```bash
go get github.com/snowmerak/zymkey.go
```

## Usage

### Initialization

Always initialize the Zymkey context before using other functions and close it when done.

```go
package main

import (
	"fmt"
	"log"
	"github.com/snowmerak/zymkey.go"
)

func main() {
	zk, err := zymkey.NewZymkey()
	if err != nil {
		log.Fatalf("Failed to initialize Zymkey: %v", err)
	}
	defer zk.Close()
    
    // Use zk instance...
}
```

### Random Number Generation

You can generate random bytes directly or use the `io.Reader` interface.

```go
// Direct generation
bytes, err := zk.GenerateRandomBytes(32)
if err != nil {
    log.Fatal(err)
}

// Using io.Reader interface
buf := make([]byte, 16)
n, err := zk.Read(buf)
```

### Data Locking (Encryption)

Encrypt and decrypt data using the Zymkey's hardware key.

```go
data := []byte("sensitive data")

// Lock (Encrypt)
locked, err := zk.LockBuffer(data)
if err != nil {
    log.Fatal(err)
}

// Unlock (Decrypt)
unlocked, err := zk.UnlockBuffer(locked)
if err != nil {
    log.Fatal(err)
}

fmt.Printf("Decrypted: %s\n", string(unlocked))
```

### Real-Time Clock (RTC)

Get the secure time from the Zymkey RTC.

```go
t, err := zk.Now()
if err != nil {
    log.Fatal(err)
}
fmt.Printf("Secure Time: %v\n", t)
```

### LED Control

Control the LED on the Zymkey module.

```go
zk.LEDOn()
time.Sleep(1 * time.Second)
zk.LEDOff()

// Flash LED: On for 100ms, Off for 100ms, 5 times
zk.LEDFlash(100*time.Millisecond, 100*time.Millisecond, 5)
```

### Device Information & Sensors

Retrieve device details and sensor readings.

```go
model, _ := zk.GetModelNumber()
serial, _ := zk.GetSerialNumber()
fw, _ := zk.GetFirmwareVersion()

fmt.Printf("Model: %s, Serial: %s, FW: %s\n", model, serial, fw)

temp, err := zk.GetCPUTemp()
if err == nil {
    fmt.Printf("CPU Temp: %.2f°C\n", temp)
}
```

## License

MIT
