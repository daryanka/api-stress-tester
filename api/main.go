package main

import (
	"fmt"
	"github.com/daryanka/api-stress-tester/api/clients"
	"github.com/daryanka/api-stress-tester/api/controllers"
	"github.com/daryanka/api-stress-tester/api/utils"
	"github.com/daryanka/api-stress-tester/api/websocket_conn"
	"github.com/joho/godotenv"
	"os"
)

func main() {
	err := godotenv.Load()
	if err != nil {
		panic(err.Error())
	}

	utils.InitLogger()
	clients.InitDatabase()
	engine := controllers.StartRouter()

	websocket_conn.NewHub()

	port := os.Getenv("PORT")
	if port == "" {
		port = "8081"
	}

	//go printConns()

	if err = engine.Run(fmt.Sprintf(":%v", port)); err != nil {
		panic(err.Error())
	}
}

//func printConns() {
//	for {
//		time.Sleep(2 * time.Second)
//
//		fmt.Println("conns", websocket_conn.Hub.Connections)
//		fmt.Println("writes", websocket_conn.Hub.WriteChannels)
		//if websocket_conn.Hub.Connections[7] != nil {
		//	for i := 0; i < 10; i++ {
		//		websocket_conn.WriteMessage(7, []byte("hello world"))
		//	}
		//} else {
		//	fmt.Println("NIL")
		//}
	//}
//}
