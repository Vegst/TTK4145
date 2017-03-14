package orders

import (
	. "../def"
	"encoding/json"
	"io/ioutil"
	"os"
	"log"
	"fmt"
)



func WriteToBackup(o Orders){
    ordersJson, errEncode := json.Marshal(o)

    if errEncode != nil {
		fmt.Println("error encoding json: ", errEncode)
	}

	filename := "orders/backups/data"

	_, errOpen := os.Open(filename)
	if errOpen != nil {
		fmt.Println("No file to write to, creating file...")
		_, _ = os.Create(filename)
	}

	errWrite := ioutil.WriteFile(filename, ordersJson, 0644)
	if errWrite != nil {
		fmt.Println("Error writing to file")
		log.Fatal(errWrite)
	}
}

func ReadFromBackup(file string) Orders{

	filename, errOpen := os.Open(file)
	if errOpen != nil {
		fmt.Println("No file to read from, creating file...")
		_, _ = os.Create(file)
		return true, decoded_client
	}

	data := make([]byte, 1024)
	n, errRead := filename.Read(data)
	if errREad != nil {
		fmt.Println("Error reading from file")
		fmt.Println(errRead)
	}

	errDecode := json.Unmarshal(data[:n], &o)
	if errDecode != nil {
		fmt.Println("Error decoding orders from backup")
	}
	return o
}