package main

import (
	"bufio"
	"fmt"
	"log"
	"net"
	"strings"
)

var clients = make(map[net.Conn]bool)

func broadcastMessage(message string) {
	for client := range clients {
		_, err := client.Write([]byte(message))

		if err != nil {
			log.Println("Mesaj gönderilemedi:", err)
			client.Close()
			delete(clients, client)
		}
	}
}

func main() {

	listener, err := net.Listen("tcp", ":8080")
	if err != nil {
		log.Fatal("Bağlantı dinlenemedi:", err)
	}
	defer listener.Close()

	fmt.Println("Sunucu başlatıldı. Bağlantı bekleniyor...")
	go broadcastMessage("Kullanıcı adı Girmeniz gerekiyor")

	for {
		conn, err := listener.Accept()
		if err != nil {
			log.Println("Bağlantı kabul edilemedi:", err)
			continue
		}
		go broadcastMessage("Kullanıcı adını belirle ve içeriye gir: ")
		clients[conn] = true

		go func(conn net.Conn) {
			defer conn.Close()

			username, err := bufio.NewReader(conn).ReadString('\n')
			if err != nil {
				log.Println("Kullanıcı adı alınamadı:", err)
				delete(clients, conn)
				return
			}
			log.Println("Kullanıcı adı yazmanız gerekiyor: ")
			username = strings.TrimSpace(username)

			fmt.Println("Yeni bir istemci bağlandı:", username)

			go broadcastMessage(username + " katıldı.\n")

			for {
				message, err := bufio.NewReader(conn).ReadString('\n')
				if err != nil {
					log.Println("Mesaj alınamadı:", err)
					break
				}

				message = strings.TrimSpace(message)

				if message == "/quit" {
					break
				}

				broadcastMessage(username + ": " + message + "\n")
			}

			delete(clients, conn)
			fmt.Println("Bir istemci ayrıldı:", username)
			conn.Close()
		}(conn)
	}
}
