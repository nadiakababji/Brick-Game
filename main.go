package main

import (
	"bufio"
	"fmt"
	"os"
	"strings"
)

func main() {

	g := gioco{
		mattoncinoPerNome: make(map[string]mattoncino),
		FilaPerMattoncino: make(map[string]*filaConCoda),
	}

	var s string
	fmt.Scan(&s)
	file, err := os.Open(s)
	if err != nil {
		fmt.Println("errore apertura file di input")
		return
	}

	myScanner := bufio.NewScanner(file)
	for myScanner.Scan() {

		line := myScanner.Text()
		command := strings.Split(line, " ")

		switch command[0] {
		case "m":

			inserisciMattoncino(g, command[1], command[2], command[3])

		case "s":

			stampaMattoncino(g, command[1])

		case "d":
			listaNomi := strings.Join(command[1:], " ")
			disponiFila(g, listaNomi)

		case "S":

			stampaFila(g, command[1])

		case "e":

			eliminaFila(g, command[1])

		case "f":

			disponiFilaMinima(g, command[1], command[2])

		case "M":

			sottostringaMassima(command[1], command[2])

		case "i":
			indiceCacofonia(g, command[1])

		case "c":
			costo(g, command[1], command[2])

		case "q":
			return

		default:

			fmt.Println("comando sconosciuto")

		}

	}
}
