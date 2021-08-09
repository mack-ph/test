package main

import (
	"bufio"
	//"bytes"
	"encoding/json"
	"fmt"
	"golang.org/x/text/encoding/simplifiedchinese"
	//"golang.org/x/text/transform"
	"io"
	//"io/ioutil"
	"log"
	"os"
	"os/exec"
	//"syscall"
	//	"github.com/axgle/mahonia"
	//"github.com/mitchellh/mapstructure"
	"rest-go-demo/tool"
	"time"
)

func ConvertByte2String(byte []byte) string {

	var str string

	decodeBytes, _ := simplifiedchinese.GB18030.NewDecoder().Bytes(byte)
	str = string(decodeBytes)

	return str
}

func RunCmd(cmd1 *exec.Cmd) {
	log.Println("Starting  ")
	ppReader, err := cmd1.StdoutPipe()
	defer ppReader.Close()
	var bufReader = bufio.NewReader(ppReader)
	if err != nil {
		fmt.Printf("create cmd stdoutpipe failed,error:%s\n", err)
		os.Exit(1)
	}
	err = cmd1.Start()
	if err != nil {
		fmt.Printf("cannot start cmd1,error:%s\n", err)
		os.Exit(1)
	}

	var buffer []byte = make([]byte, 4096)
	for {
		n, err := bufReader.Read(buffer)
		if err != nil {
			if err == io.EOF {

				log.Println("Ending  ")
				fmt.Printf("pipi has Closed\n")
				time.Sleep(1 * time.Second)
				err := stopProcess(cmd1)
				if err != nil {
					fmt.Printf("stop child process failed,error:%s", err)
					os.Exit(1)
				}
				cmd1.Wait()

				break
			} else {
				fmt.Println("Read content failed")
			}
		}
		garbledStr := ConvertByte2String(buffer[:n])
		fmt.Print(garbledStr)
	}
}

func stopProcess(cmd *exec.Cmd) error {

	fmt.Printf("结束子进程%s成功\n", cmd.Path)
	if err := cmd.Process.Kill(); err != nil {
		log.Fatal("failed to kill process: ", err)
		return err
	}

	return nil
}

func main2() {

	datas := string(`
    { 
      "type": "person",
      "name":"dj",
      "age":18,
      "job": "programmer"
    } `)

	var m map[string]interface{}
	err := json.Unmarshal([]byte(datas), &m)
	if err != nil {
		log.Fatal(err)
	}
	log.Println(m["type"])

	c := tool.NewCycleQueue(3)
	fmt.Println(c.Push(1))
	fmt.Println("realLength", c.QueueLength(), "queue", c.Display())
	fmt.Println(c.Pop())
	fmt.Println(c.Push(2))
	fmt.Println("realLength", c.QueueLength(), "queue", c.Display())

	cmd1 := exec.Command("ping", "-t", "www.google.com")

	go RunCmd(cmd1)
	time.Sleep(10 * time.Second)
	stopProcess(cmd1)
	time.Sleep(3 * time.Second)
	log.Printf("game over  command  %s /n", cmd1)

}
