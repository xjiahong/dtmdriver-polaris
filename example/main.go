package main

import (
	"fmt"
	"time"

	"github.com/dtm-labs/dtmdriver"
	pd "github.com/dtm-labs/dtmdriver-polaris"
)

func main() {
	err := dtmdriver.Use(pd.Name)
	fmt.Println("got err", err)
	err = dtmdriver.GetDriver().RegisterGrpcService("polaris://0.0.0.0:8080/dtm?namespace=ssv-szzj",
		"")
	fmt.Println("got err", err)
	time.Sleep(30 * time.Second)
}
