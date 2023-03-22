package lib

import (
	"fmt"
	"time"

	restyLib "github.com/go-resty/resty/v2"
)

var Resty *restyLib.Client

func init() {
	if Resty == nil {
		fmt.Println("resty init...")
		Resty = restyLib.New()
		Resty.SetTimeout(10 * time.Second)
	}
}
