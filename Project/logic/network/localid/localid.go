package localid

import (
	"./localip"
	"os"
	"fmt"
)

func ID() string {
	localIP, err := localip.LocalIP()
	if err != nil {
		fmt.Println(err)
		localIP = "DISCONNECTED"
	}
	return fmt.Sprintf("peer-%s-%d", localIP, os.Getpid())
}
