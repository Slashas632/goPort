package output

import (
	"encoding/json"
	"os"
	"sync"
)

type Result struct {
	IP     string `json:"ip"`
	Port   int    `json:"port"`
	Banner string `json:"banner"`
}

var (
	results []Result
	mu      sync.Mutex
)

var JSONPath string

func Init(path string) {
	JSONPath = path
}

func AddResult(ip string, port int, banner string) {
	mu.Lock()
	defer mu.Unlock()
	results = append(results, Result{IP: ip, Port: port, Banner: banner})

}

func SaveJson(path string) error {
	mu.Lock()
	defer mu.Unlock()

	f, err := os.Create(path)
	if err != nil {
		return err
	}
	defer f.Close()

	encoder := json.NewEncoder(f)
	encoder.SetIndent("", "  ")
	return encoder.Encode(results)
}
