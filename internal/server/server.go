package server

import (
	"bufio"
	"chidis/internal/handlers"
	"chidis/internal/router"
	"log"
	"net"
)

func HandleClient(conn net.Conn) {
	defer conn.Close()
	readerConn := bufio.NewReader(conn)
	for {
		msg, err := readerConn.ReadString('\n')
		if err != nil {
			log.Println("не удалось прочитать сообщение, ошибка:", err)
		}
		answer, err := router.Routing(msg)
		if err == handlers.ErrQuit {
			conn.Write([]byte(answer))
			log.Println("клиент ушел")
			return
		}
		if err != nil {

			conn.Write([]byte("-dont get answer, error:" + err.Error() + "\n"))

			log.Println("ошибка для клиента: ", err)

		}
		conn.Write([]byte(answer + "\n"))
	}
}
