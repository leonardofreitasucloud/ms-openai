package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"time"
)

func main() {
	log.Println("*** [INFO]: ListenAndServe starting on port 8090...")

	mux := http.NewServeMux()
	mux.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprint(w, "Welcome to ms-openai")
	})

	mux.HandleFunc("/api/v1", getApi)
	//mux.HandleFunc("/api/v1/chatgpt/completions", getCompletions)

	// Verifica se o servidor está ativo antes de iniciar
	go isUp()

	err := http.ListenAndServe(":8090", mux)
	if err != nil {
		fmt.Printf("server closed\n")
		log.Fatal(err)
	}

}

func requestOpenAiCompletions() (string, error) {
	// Dados da API Key
	apiKey := os.Getenv("API_KEY")

	// Dados da mensagem para o ChatGPT
	input := "Definição de Kubernetes"

	// Constrói a requisição HTTP Post
	requestBody, err := json.Marshal(map[string]interface{}{
		"model":       "ada-search-document",
		"prompt":      input,
		"max_tokens":  100,
		"temperature": 1.0,
	})
	if err != nil {
		fmt.Println("*** [ERROR]: Erro ao contruir o corpo da requisição: ", err.Error())
		return "", err
	}

	req, err := http.NewRequest("POST", "https://api.openai.com/v1/completions", bytes.NewBuffer(requestBody))
	if err != nil {
		fmt.Println("*** [ERROR]: Erro ao criar a requisição: ", err.Error())
		return "", err
	}

	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", "Bearer "+apiKey)

	//Executa a requisição
	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		fmt.Println("*** [ERROR]: Erro ao executar a requisição: ", err.Error())
		return "", err
	}
	defer resp.Body.Close()

	// Lê a resposta
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		fmt.Println("*** [ERROR]: Erro ao ler a resposta: ", err.Error())
		return "", err
	}

	return string(body), nil
}

func isUp() {
	// Aguarda 1 segundo antes de executar a verificação
	time.Sleep(3 * time.Second)

	req, err := http.NewRequest("POST", "http://localhost:8090/", nil)
	if err != nil {
		log.Fatal(err)
	}

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		log.Fatal(err)
	}
	defer resp.Body.Close()

	if resp.StatusCode == http.StatusOK {
		log.Println("*** [INFO]: Requisição http://localhost:8090/ bem-sucedida (200 OK)")
		log.Println("*** [INFO]: ListenAndServe started !")
	} else {
		log.Println("*** [ERROR]: Requisição retornou um status diferente de 200 OK:", resp.Status)
	}
}

func getApi(w http.ResponseWriter, r *http.Request) {
	log.Println("*** [INFO]: request: /api/v1")
	io.WriteString(w, "Request /api/v1")
}

func getCompletions(w http.ResponseWriter, r *http.Request) {
	responseBody, err := requestOpenAiCompletions()
	if err != nil {
		log.Println("*** [ERROR]: Erro na solicitação do OpenAI Completions:", err)
		http.Error(w, "Erro na solicitação do OpenAI Completions", http.StatusInternalServerError)
		return
	}

	log.Println("*** [INFO]: request: /api/v1/chatgpt/completions")
	io.WriteString(w, responseBody)
}
