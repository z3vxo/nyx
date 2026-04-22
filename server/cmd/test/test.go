package main

import (
	"bytes"
	"crypto/tls"
	"encoding/binary"
	"flag"
	"fmt"
	"io"
	"net/http"
	"os"

	"github.com/google/uuid"
)

func writeField(buf *bytes.Buffer, s string) {
	b := []byte(s)
	length := make([]byte, 4)
	binary.LittleEndian.PutUint32(length, uint32(len(b)))
	buf.Write(length)
	buf.Write(b)
}

func buildRegisterPayload() []byte {
	var buf bytes.Buffer

	// CMD TYPE - 4 bytes LE (1 = register)
	cmdType := make([]byte, 4)
	binary.LittleEndian.PutUint32(cmdType, 1)
	buf.Write(cmdType)

	writeField(&buf, uuid.NewString())        // GUID
	writeField(&buf, "TESTLAB\\operator")     // USERNAME
	writeField(&buf, "DESKTOP-NYX01")         // HOSTNAME
	writeField(&buf, "192.168.1.100")         // INTERNAL IP
	writeField(&buf, `C:\Windows\System32\svchost.exe`) // PROCESS PATH

	// PID - 4 bytes LE
	pid := make([]byte, 4)
	binary.LittleEndian.PutUint32(pid, 1337)
	buf.Write(pid)

	// PPID - 4 bytes LE
	ppid := make([]byte, 4)
	binary.LittleEndian.PutUint32(ppid, 512)
	buf.Write(ppid)

	// IS_ELEVATED - 1 byte
	buf.WriteByte(1)

	// ARCH - 1 byte (1 = x64, 0 = x86)
	buf.WriteByte(1)

	// MINOR VERSION - 2 bytes LE
	minor := make([]byte, 2)
	binary.LittleEndian.PutUint16(minor, 0)
	buf.Write(minor)

	// MAJOR VERSION - 2 bytes LE
	major := make([]byte, 2)
	binary.LittleEndian.PutUint16(major, 10)
	buf.Write(major)

	// BUILD VERSION - 2 bytes LE
	build := make([]byte, 2)
	binary.LittleEndian.PutUint16(build, 19045)
	buf.Write(build)

	return buf.Bytes()
}

func main() {
	host := flag.String("host", "127.0.0.1", "listener host")
	port := flag.Int("port", 8080, "listener port")
	endpoint := flag.String("endpoint", "/ms/upload", "POST endpoint")
	proto := flag.String("proto", "https", "protocol: http or https")
	flag.Parse()

	payload := buildRegisterPayload()
	url := fmt.Sprintf("%s://%s:%d%s", *proto, *host, *port, *endpoint)

	fmt.Printf("[*] Sending register payload (%d bytes) to %s\n", len(payload), url)

	client := &http.Client{}
	if *proto == "https" {
		client.Transport = &http.Transport{
			TLSClientConfig: &tls.Config{InsecureSkipVerify: true},
		}
	}
	resp, err := client.Post(url, "application/octet-stream", bytes.NewReader(payload))
	if err != nil {
		fmt.Fprintf(os.Stderr, "[-] Request failed: %v\n", err)
		os.Exit(1)
	}
	defer resp.Body.Close()

	body, _ := io.ReadAll(resp.Body)
	fmt.Printf("[*] Status: %s\n", resp.Status)
	fmt.Printf("[*] Body:   %s\n", string(body))
}
