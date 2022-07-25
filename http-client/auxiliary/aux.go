package aux

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
)

func printResponse(resp *http.Response) {
	if resp.StatusCode == 200 {
		b, err := io.ReadAll(resp.Body)
		if err != nil {
			panic(err)
		}
		prettyJSON, err := formatJSON(b)
		if err != nil {
			panic(err)
		}
		fmt.Printf("%s\n", prettyJSON)
	} else {
		fmt.Printf("Status is: %s\n", resp.Status)
	}
}

func formatJSON(data []byte) ([]byte, error) {
	var out bytes.Buffer
	err := json.Indent(&out, data, "", "    ")
	if err == nil {
		return out.Bytes(), err
	}
	return data, nil
}
