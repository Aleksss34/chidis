package app

import (
	"chidis/internal/server"
	"log"
	"net"
)

const (
	network = "tcp"
	port    = ":6721"
)

func RunServer() {
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
	for {
		conn, err := listener.Accept()
		if err != nil {
			log.Println("не получилось подключится к клиенту, ошибка:", err)
		}
		go server.HandleClient(conn)
	}
}
