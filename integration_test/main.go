package main

import (
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
	"os"
	"os/exec"
	"strconv"
	"strings"
	"time"
)

func main() {
	const port = "8771"
	const exporterBinaryName = "process_exporter"
	const dummyBinaryName = "prcexpintdum"
	const dummyDescripiveName = "iamdummy"
	run("go", "build", "-o", exporterBinaryName, "github.com/setlog/process_exporter/cmd/exporter")
	defer os.Remove(exporterBinaryName)
	run("go", "build", "-o", dummyBinaryName, "github.com/setlog/process_exporter/integration_test/dummy")
	defer os.Remove(dummyBinaryName)
	cmd := exec.Command("./"+exporterBinaryName, "-port", port, "-binary", dummyBinaryName)
	err := cmd.Start()
	if err != nil {
		panic(err)
	}
	defer cmd.Process.Signal(os.Interrupt)

	cmdDummy := exec.Command("./"+dummyBinaryName, "-name", dummyDescripiveName)
	dummyReader, err := cmdDummy.StdoutPipe()
	if err != nil {
		panic(err)
	}
	dummyWriter, err := cmdDummy.StdinPipe()
	if err != nil {
		panic(err)
	}
	defer dummyWriter.Close()
	err = cmdDummy.Start()
	if err != nil {
		panic(err)
	}
	b := []byte{0}
	_, err = io.ReadAtLeast(dummyReader, b, 1)
	if err != nil {
		panic(err)
	}
	if b[0] != 0x2A {
		panic(fmt.Sprintf("Byte was %x. Expected 0x2A", b[0]))
	}
	time.Sleep(time.Second)

	resp, err := http.Get("http://localhost:" + port + "/metrics")
	if err != nil {
		panic(err)
	}
	defer resp.Body.Close()
	dummyWriter.Close()
	if resp.StatusCode != http.StatusOK {
		panic(fmt.Sprintf("got code %d when %d expected", resp.StatusCode, http.StatusOK))
	}
	contents, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		panic(err)
	}
	contentsString := string(contents)
	lines := strings.Split(contentsString, "\n")
	gotLine := false
	for _, line := range lines {
		prefix := fmt.Sprintf("mine_disk_write_bytes{bin=\"%s\",name=\"%s\",pid=\"%d\"} ", dummyBinaryName, dummyDescripiveName, cmdDummy.Process.Pid)
		if strings.HasPrefix(line, prefix) {
			writeBytes, err := strconv.Atoi(line[len(prefix):])
			if err != nil {
				panic(fmt.Sprintf("failed to parse write_bytes from line %s: %v", line, err))
			}
			if writeBytes < 1024 || writeBytes > 8192 {
				panic(fmt.Sprintf("exporter reported %d bytes written. Expected something in the range [1024;8192].", writeBytes))
			}
			gotLine = true
		}
	}
	if !gotLine {
		panic(fmt.Sprintf("exporter did not report io byte write count"))
	}
	err = cmdDummy.Wait()
	if err != nil {
		panic(err)
	}
}

func run(command string, args ...string) {
	cmd := exec.Command(command, args...)
	out, err := cmd.CombinedOutput()
	if err != nil {
		fmt.Println(string(out))
		panic(err)
	}
}
