package main

func main() {
	server := InitializeApp()
	server.Run(":9999")
}

// go run .
// curl -i -X GET http://localhost:9999/post/1
