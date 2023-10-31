package start

import "github.com/don764372409/go-im-sdk/client"

func Start(url string) {
	client := client.New(url)
	client.Startup()
}
