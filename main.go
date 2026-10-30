package main

import (
	"bytes"
	"crypto/rand"
	"encoding/binary"
	"fmt"
	"io"
	"log"
	"net"
	"time"
)

type FileServer struct {
	network string
	address string
}

func NewFileServer(network string, address string) *FileServer {
	return &FileServer{network: network, address: address}
}

func (fs *FileServer) start() {
	line, err := net.Listen(fs.network, fs.address)
	if err != nil {
		log.Fatal(err)
	}

	for {
		conn, err := line.Accept()
		if err != nil {
			log.Fatal(err)
		}
		fs.readLoop(conn)
	}
}
func (fs *FileServer) readLoop(conn net.Conn) {
	buf := new(bytes.Buffer)
	for {
		var size int64
		binary.Read(conn, binary.LittleEndian, &size)
		n, err := io.CopyN(buf, conn, size)
		if err != nil {
			log.Fatal(err)
		}

		fmt.Println("\n\n", buf.Bytes())
		fmt.Printf("recieved %d byte over tcp", n)
	}
}
func (fs *FileServer) sendFile(size int) error {
	file := make([]byte, size)
	_, err := io.ReadFull(rand.Reader, file)
	if err != nil {
		return err
	}

	conn, err := net.Dial(fs.network, fs.address)
	if err != nil {
		return err
	}

	binary.Write(conn, binary.LittleEndian, int64(size))
	n, err := io.CopyN(conn, bytes.NewReader(file), int64(size)) //copyn beacause we need end of file in reciever reader
	if err != nil {
		return err
	}

	fmt.Printf("written %d byte over tcp\n", n)
	return nil
}

func main() {
	fileSrever := NewFileServer("tcp", ":3000")
	go func() {
		time.Sleep(time.Second * 5)
		fileSrever.sendFile(40000000)
	}()

	fileSrever.start()

}
