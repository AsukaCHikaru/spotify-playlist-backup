package write

import (
	"encoding/json"
	"fmt"
	"os"
)

func WriteToJson(result any) {
	fmt.Println("Writing to file...")
	jsonData, err := json.MarshalIndent(result, "", "  ")
	if err != nil {
		fmt.Println(err.Error())
		return
	}

	err = os.WriteFile("output.json", []byte(string(jsonData)), 0644)
	if err != nil {
		fmt.Println(err.Error())
		return
	}
}
