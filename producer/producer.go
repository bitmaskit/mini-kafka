package main

import (
	"bufio"
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"os"
	"strings"
)

type producer struct {
	url string
}

func New(url string) *producer {
	return &producer{url: url}
}

func (p *producer) Publish(payload []byte) error {
	resp, err := http.Post(fmt.Sprintf("%s/publish", p.url), "application/json", bytes.NewBuffer(payload))
	if err != nil {
		return err
	}
	return resp.Body.Close()
}

func main() {
	prd := New("http://localhost:9092")
	scanner := bufio.NewScanner(os.Stdin)

	fmt.Println("=== Publisher CLI ===")
	fmt.Println("Type: <topic> <message>")
	fmt.Println("Example: foo hello world")
	fmt.Println("------------------------")

	for {
		fmt.Print("📝 > ")
		if !scanner.Scan() {
			break
		}
		line := scanner.Text()
		parts := strings.SplitN(line, " ", 2)
		if len(parts) != 2 {
			fmt.Println("⚠️  Invalid. Format: <topic> <message>")
			continue
		}

		payload := map[string]string{
			"topic": parts[0],
			"data":  parts[1],
		}
		body, err := json.Marshal(payload)
		if err != nil {
			fmt.Println("❌ Failed to marshal:", err)
			continue
		}

		if err = prd.Publish(body); err != nil {
			fmt.Println("❌ Failed to publish:", err)
			continue
		}
		fmt.Printf("✅ Sent to [%s]: %s\n", parts[0], parts[1])
	}
}
