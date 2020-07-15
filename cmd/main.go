package main

func main() {

	client, closeConnection := ConnectMongo()
	defer func() { closeConnection() }()

	c := Init(client)
}
