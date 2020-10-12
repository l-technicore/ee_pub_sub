package utils

import(
	"testing"
	"os/exec"
	"bytes"
	"log"
)

func TestInitiateLogger_and_CloseLogger(t *testing.T) {
	log.Printf("Creating a Test File for logging.")
	if err := InitiateLogger("test_file_xyz.out"); err != nil {
		t.Error(err.Error())
	}
	log.Printf("Closing the opened logging Test File.")
	if err := CloseLogger(); err != nil {
		t.Error("File not closed after calling CloseLogger. Error: "+err.Error())
	}
	cmd := exec.Command("rm", "test_file_xyz.out")
	var out bytes.Buffer
	cmd.Stdout = &out
	if err := cmd.Run();err == nil{
		log.Printf("Deleting the Test File: %v", out.String())
	}else{
		t.Error("Command finished with error:" + err.Error())
	}
}