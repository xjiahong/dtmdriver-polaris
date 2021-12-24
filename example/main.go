package main

import (
	"fmt"
	pd "github.com/ychensha/dtmdriver-polaris"
	"github.com/yedf/dtmdriver"
	"time"
)

func main() {
	err := dtmdriver.Use(pd.Name)
	fmt.Println("got err", err)
	err = dtmdriver.GetDriver().RegisterGrpcService("polaris://0.0.0.0:8080/your.service?namespace=Test",
		"YOUR_TOKEN")
	fmt.Println("got err", err)
	time.Sleep(30 * time.Second)
}
