package integration

import (
	"fmt"
	"io"
	"net/http"
	"os"
	"os/exec"
	"testing"
)

func TestIntegration(t *testing.T) {
	const port = "8771"
	run("go", "build", "-o", "exporter", "github.com/setlog/process_exporter/cmd/exporter")
	defer os.Remove("exporter")
	run("go", "build", "-o", "proc_exporter_integration_dummy", "github.com/setlog/process_exporter/test/dummy")
	defer os.Remove("proc_exporter_integration_dummy")
	cmd := exec.Command("exporter", "-port", port, "-binary", "proc_exporter_integration_dummy")
	err := cmd.Start()
	if err != nil {
		panic(err)
	}
	cmdDummy := exec.Command("proc_exporter_integration_dummy", "-name", "iamdummy")
	dummyReader, err := cmdDummy.StdoutPipe()
	if err != nil {
		panic(err)
	}
	dummyWriter, err := cmdDummy.StdinPipe()
	if err != nil {
		panic(err)
	}
	err = cmdDummy.Start()
	if err != nil {
		panic(err)
	}
	b := make([]byte, 1, 1)
	_, err = io.ReadAtLeast(dummyReader, b, 1)
	if err != nil {
		panic(err)
	}
	resp, err := http.Get("http://localhost:" + port)
	if err != nil {
		panic(err)
	}
	if resp.StatusCode != http.StatusOK {
		panic(fmt.Sprintf("got code %d when %d expected", resp.StatusCode, http.StatusOK))
	}
	err = dummyWriter.Close()
	if err != nil {
		panic(err)
	}
	err = cmdDummy.Wait()
	if err != nil {
		panic(err)
	}

}

func run(command string, args ...string) {
	cmd := exec.Command(command, args...)
	err := cmd.Run()
	if err != nil {
		panic(err)
	}
}
