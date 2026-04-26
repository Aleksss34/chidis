package app

import (
	"bufio"
	"chidis/internal/router"
	"chidis/internal/server"
	"chidis/internal/storage"
	"log"
	"net"
)

const (
	network = "tcp"
	port    = ":6721"
)

func RunServer() {
	if err := storage.InitStorage(); err != nil {
		log.Fatal("не удалось прочитать базу данных, ошибка:", err)
	}
	if err := Recovery(); err != nil {
		log.Println(err)
	}

	tcpAddr, err := net.ResolveTCPAddr(network, port)
	if err != nil {
		log.Fatal("не получилось сгенерировать tcp адрес, ошибка:", err)
	}

	listener, err := net.ListenTCP(network, tcpAddr)
	if err != nil {
		log.Fatal("не получилось запустить слушателя, ошибка:", err)
	}
	log.Println("сервер успешно запущен")
	defer listener.Close()
	go storage.TickExpired()
	for {
		conn, err := listener.Accept()
		if err != nil {
			log.Println("не получилось подключится к клиенту, ошибка:", err)
		}
		go server.HandleClient(conn)
	}
}
func Recovery() error {
	storage.IsLoading = true
	defer func() { storage.IsLoading = false }()
	if _, err := storage.DBFile.Seek(0, 0); err != nil {
		return err
	}
	scannerDB := bufio.NewScanner(storage.DBFile)
	log.Println("читаю предыдущие записи:")

	for scannerDB.Scan() {
		answer, err := router.Routing(scannerDB.Text())
		if err != nil {
			return err
		}
		log.Println(answer)
	}
	return nil
}
