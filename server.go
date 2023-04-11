package main

import ("encoding/gob"
		"fmt"
		"net"
		"net/http" ; "io")

//HTTP server
	func hello(res http.ResponseWriter, req *http.Request){
		res.Header().Set(
			"Content-Type", "Text/html",
		)
		io.WriteString(
			res, 
			`<doctype html>
				<html>
					<head>
						<title>Hello World</title>
					</head>
				<body>
					Hello World!
				</body>
				</html>`,
		)
	}


//normal server
		func server(){
			//listenning on a certain port
			ln, err := net.Listen("tcp", ":9999")
			if err != nil {
				fmt.Println(err)
				return
			}

			for {
				//accept a connection
				c, err := ln.Accept()
				if err != nil {
					fmt.Println(err)
					continue
				}
				//handling the server connection
				go handleServerConnection(c)


			}
		}

		func handleServerConnection(c net.Conn ) {
			//recieving messages
			var message string
			err := gob.NewDecoder(c).Decode(&message)
			if err != nil {
				fmt.Println(err)
			}else{
				fmt.Println("Message '", message, "' has been recieved,")
			}
			c.Close()
		}

		func client( message string) {
			//client connecting to server
			c, err := net.Dial("tcp", "localhost")
			if err != nil {
				fmt.Println(err)
				return
			}

			//sending message to the server
			fmt.Println("Sending message ", message)
			err = gob.NewEncoder(c).Encode(message)
			if err != nil{
				fmt.Println(err)
			}
			c.Close()
		}

		func main() {
			http.HandleFunc("/hello", hello)
			http.ListenAndServe(":9000", nil)
		}