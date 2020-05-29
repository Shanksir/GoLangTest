package main

import (
	"bufio"
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
	"os"
	"strconv"
	"strings"
	"time"
)

const monitoramentos = 4
const delay = 5

func exibeIntro() {
	var versao float32 = 1.1
	fmt.Println("Esse programa está na versão:", versao)
	fmt.Println("Selecione a opção que desejar:")
}

func leComando() int {
	var resposta int
	fmt.Scan(&resposta)

	return resposta
}

func leSitesArquivo() []string {
	var sites []string
	arquivo, err := os.Open("sites.txt")
	if err != nil {
		fmt.Println("Ocorreu um erro:", err)
	}
	leitor := bufio.NewReader(arquivo)

	for {
		linha, err := leitor.ReadString('\n')
		linha = strings.TrimSpace(linha)
		sites = append(sites, linha)
		if err == io.EOF {
			break
		}
	}
	return sites
}

func imprimeLogs() {
	arquivo, err := ioutil.ReadFile("log.txt")
	if err != nil {
		fmt.Println("Ocorreu um erro:", err)
	}
	fmt.Println(string(arquivo))
}

func registraLog(site string, status bool) {
	arquivo, err := os.OpenFile("log.txt", os.O_RDWR|os.O_CREATE|os.O_APPEND, 0666)
	if err != nil {
		fmt.Println("Ocorreu um erro:", err)
	}

	arquivo.WriteString(time.Now().Format("02/01/2006 15:04:05") + " - " + site + " Online: " + strconv.FormatBool(status) + "\n")

	arquivo.Close()
}

func iniciarMonitoramento() {
	fmt.Println("Monitorando...")

	sites := leSitesArquivo()

	for i := 0; i < monitoramentos; i++ {
		for i, site := range sites {
			resposta, err := http.Get(site)
			if err != nil {
				fmt.Println("Ocorreu um erro:", err)
			}
			if resposta.StatusCode == 200 {
				fmt.Println("Site ", i, ":", site, " OK!")
				registraLog(site, true)
			} else {
				fmt.Println("Site ", i, ":", site, " OFF!")
				registraLog(site, false)
			}
		}
		time.Sleep(delay * time.Second)
		fmt.Println("")
	}
}

func main() {

	for {

		fmt.Println("1 - Iniciar Monitoramento.")
		fmt.Println("2 - Exibir Log.")
		fmt.Println("3 - Sair.")

		switch leComando() {
		case 1:
			iniciarMonitoramento()
		case 2:
			fmt.Println("Exibindo log...")
			imprimeLogs()
		case 3:
			fmt.Println("Saindo...")
			os.Exit(0)
		default:
			fmt.Println("Opção inválida.")
			os.Exit(-1)
		}
	}

}
