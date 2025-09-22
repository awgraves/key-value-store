package main

func main() {
	kvAPIv1BaseURL := getKVServiceAPIv1BaseURL()
	client := NewAPIv1Client(kvAPIv1BaseURL)
	r := setupRouter(client)
	r.Run(":8081")
}
