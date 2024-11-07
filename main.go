package main

import (
	fileSrv "NodeBox/file_server"
	"NodeBox/node"
	"bytes"
	"flag"
	"fmt"
	"gopkg.in/yaml.v3"
	"io"
	"log"
	"os"
	"time"
)

var nodes = []string{":3001", ":3002", "3000"}

func main() {
	option := flag.String("option", "send", "option for file server |send -file ... | / |get -key ... |")
	//file := flag.String("file", "no", "file path for sending")
	//confFile := flag.String("conf", "config.yaml", "config file for file server")
	listenAddr := flag.String("addr", ":1234", "listen adders")
	flag.Parse()

	fsMain := makeServer(*listenAddr, "lala", nodes)

	// Запускаємо сервер асинхронно

	go func() {
		fsMain.Start()
	}()
	time.Sleep(2 * time.Second)

	switch *option {
	case "store":
		fmt.Println("store")
		file, err := os.ReadFile("Tanenbaum_Networks.pdf")
		if err != nil {
			fmt.Println(err)
		}
		err = fsMain.StoreFile("book", bytes.NewReader(file))
		if err != nil {
			fmt.Println(err)
		}
	case "get":
		data, err := fsMain.GetFile("book")
		if err != nil {
			fmt.Println(err)
		}

		err = fsMain.StoreFile("book", data)
		if err != nil {
			fmt.Println(err)
		}
	}

	fmt.Println("lala")

	select {}

}

type FSConfig struct {
	ListenAddr string   `yaml:"listenAddr"`
	ID         string   `yaml:"ID"`
	Nodes      []string `yaml:"Nodes"`
}

func parseConf(path string) (*FSConfig, error) {
	conf := &FSConfig{}

	file, err := os.Open(path)
	if err != nil {
		return nil, err
	}

	fileData, er := io.ReadAll(file)
	if er != nil {
		return nil, er
	}

	err = yaml.Unmarshal(fileData, conf)
	if err != nil {
		return nil, err
	}

	return conf, nil
}

func makeServer(listenAddr string, id string, nodes []string) *fileSrv.FileServer {
	tcpOpts := node.TCPTransportOptions{
		ListenAddress: listenAddr,
		Decoder:       node.DefaultDecoder{},
		HandshakeFunc: node.NILHandshake,
	}

	tcpTrsp := node.NewTCPTransport(tcpOpts)

	fileServerOpts := fileSrv.FileServerOpts{
		ID:                id,
		EncKey:            fileSrv.NewEncKey(),
		StorageRoot:       listenAddr + "_store",
		PathTransformFunc: fileSrv.CASPath,
		Transport:         tcpTrsp,
		BootstrapNodes:    nodes,
	}

	fs, err := fileSrv.NewFileServer(fileServerOpts)
	if err != nil {
		log.Fatal(err)
	}

	tcpTrsp.OnNode = fs.OnNode

	return fs
}
