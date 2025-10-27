package main

import (
	"fmt"
	"math"
	"strings"
)

type mattoncino struct {
	formaBordoSx string
	formaBordoDx string
	nome         string
}

type gioco struct {
	mattoncinoPerNome map[string]mattoncino
	FilaPerMattoncino map[string]*filaConCoda
}

// aggiunge il mattoncino al gioco
func inserisciMattoncino(g gioco, alpha, beta, sigma string) {

	if alpha == beta {
		return
	}
	if _, ok := g.mattoncinoPerNome[sigma]; !ok {
		g.mattoncinoPerNome[sigma] = mattoncino{alpha, beta, sigma}
	}

	return
}

// stampa il mattoncino
func stampaMattoncino(g gioco, sigma string) {

	if value, ok := g.mattoncinoPerNome[sigma]; ok {
		fmt.Println(value)
	}

	return
}

func disponiFila(g gioco, listaNomi string) {
	//controllli:

	//caso: primo mattoncino
	input := strings.Split(listaNomi, " ")
	flip := input[0][0] == '-'
	sigma := input[0][1:]

	//controllo che non ci sia un mattoncino con lo stesso nome in scatola
	if !inScatola(g, sigma) {
		return
	}

	// mattoncino corrente
	m := g.mattoncinoPerNome[sigma]

	// Determino la forma del bordo destro, tenendo in considerazione la disposizione del mattoncino
	formaBordoDxPrev := m.formaBordoDx
	if flip {
		formaBordoDxPrev = m.formaBordoSx
	}

	formaBordoDx := "" //contine la forma dx del mattoncino corrente, tiene in considerazione la disposizione del mattoncino
	formaBordoSx := "" //contine il bordo sx del mattoncino corrente, tiene in considerazione la disposizione del mattoncino

	// controlli sui restanti mattoncini della fila
	for i := 1; i < len(input); i++ {
		flip = input[i][0] == '-'
		sigma = input[i][1:]

		if !inScatola(g, sigma) {
			return
		}

		m = g.mattoncinoPerNome[sigma]

		// Determino le forme dei bordi destro e sinistro, considerando la disposizione del mattoncino, per determinare quale bordo si dovrà incastrare con il mattoncino precedente
		formaBordoSx = m.formaBordoSx
		formaBordoDx = m.formaBordoDx
		if flip {
			formaBordoSx, formaBordoDx = formaBordoDx, formaBordoSx
		}

		// Controllo che il mattoncino corrente si incastro con quello precedente
		if formaBordoSx != formaBordoDxPrev {
			return
		}

		formaBordoDxPrev = formaBordoDx

	}

	//Fine controlli, creo la fila:
	creaFila(g, listaNomi)
	return

}

func stampaFila(g gioco, sigma string) {

	// Verifica se il mattoncino con nome sigma è presente in una fila
	if f, inTavolo := g.FilaPerMattoncino[sigma]; inTavolo {
		stampaFilaConCoda(g, *f)
	}

	return
}

func eliminaFila(g gioco, sigma string) {

	// Verifica che il mattoncino con nome sigma sia in una fila del tavolo
	if f, inTavolo := g.FilaPerMattoncino[sigma]; inTavolo {

		// tutti i mattoncini che compongono la fila vengono rimessi nella scatola
		p := f.testa
		for p != nil {
			delete(g.FilaPerMattoncino, p.nomeMattoncino)
			p = p.next
		}

	}

}

func disponiFilaMinima(g gioco, alpha, beta string) {

	grafo := map[string][]string{}             //lista di adiacenza
	bordiToMattoncino := map[string][]string{} //chiave è nella forma: formaBordoSx+formaBordoDx, aggiungo un "+" tra i due lati perchè un mattoncino che ha lati "ab" e  "c" e un altro che ha "a" "bc" darebbero stessa chiave

	var formaSx, formaDx, nome string

	for k, v := range g.mattoncinoPerNome {
		//per i mattoncini nella scatola
		if inScatola(g, k) {
			formaSx = v.formaBordoSx
			formaDx = v.formaBordoDx
			nome = v.nome

			// Aggiornamento grafo
			grafo[formaSx] = append(grafo[formaSx], formaDx)
			grafo[formaDx] = append(grafo[formaDx], formaSx)

			if formaSx < formaDx {
				bordiToMattoncino[formaSx+"+"+formaDx] = append(bordiToMattoncino[formaSx+"+"+formaDx], nome)
			} else {
				bordiToMattoncino[formaDx+"+"+formaSx] = append(bordiToMattoncino[formaDx+"+"+formaSx], nome)
			}
		}

	}

	// Controllo che alpha sia presente nel grafo
	sequenzaMin := "" //di bordi
	lenSequenzaMin := math.MaxInt

	if _, ok := grafo[alpha]; ok {

		if alpha == beta {

			sequenzaCurr := ""

			viciniVisitati := map[string]bool{} //se eseguo bfs una volta su un specifica forma, non lo rieseguo se nei vicini ci sono ripetizioni della stessa forma

			for _, vicino := range grafo[alpha] { //se ci sono più vicini ugualli posso evitare di fare la bfs

				if viciniVisitati[vicino] {
					continue
				}
				sequenzaCurr = bfs(grafo, bordiToMattoncino, vicino, beta, true, lenSequenzaMin) //exlude mi serve solo se c'è un mattoncino che ha arco vicino+beta
				//la bfs restituisce una stringa "" se non ha trovato nessuna sequenza con lunghezza minore della sequenza minima corrente
				if sequenzaCurr != "" {
					sequenzaMin = sequenzaCurr
					lenSequenzaMin = len(strings.Split(sequenzaMin, "+"))
				}
				viciniVisitati[vicino] = true

			}
			//aggiungo al cammino minimo il nodo di partenza
			sequenzaMin = alpha + "+" + sequenzaMin

		} else { //alpha!=beta
			// Controllo che beta sia presente nel grafo
			if _, ok := grafo[beta]; ok {
				sequenzaMin = bfs(grafo, bordiToMattoncino, alpha, beta, false, math.MaxInt)

			}
		}
	}

	//verifica di avere trovato un cammino da alpha a beta.
	if sequenzaMin == "" || sequenzaMin == alpha+"+" {
		fmt.Printf("non esiste fila da %s a %s\n", alpha, beta)
		return
	}

	// Costruzione della lista di nomi di mattoncini a partire dalla sequenzaMin
	listaNomi := ""
	seqForme := strings.Split(sequenzaMin, "+")
	key := seqForme[0]
	segno := ""

	for i := 1; i < len(seqForme); i++ {
		//scrivo la chiave nell'ordine corretto
		if key > seqForme[i] {
			key = seqForme[i] + "+" + key
		} else {
			key += "+" + seqForme[i]
		}

		if v, ok := bordiToMattoncino[key]; ok {

			//oriento il mattoncino nel modo corretto
			if g.mattoncinoPerNome[v[len(v)-1]].formaBordoDx == seqForme[i] {
				segno = "+"
			} else {
				segno = "-"
			}

			listaNomi += segno + v[len(v)-1] + " "

			if len(v) == 1 {
				delete(bordiToMattoncino, key)
			} else {
				bordiToMattoncino[key] = v[:len(v)-1]
			}

		}
		key = seqForme[i]

	}

	// Pulizia della lista di nomi e creazione della fila
	listaNomi = strings.TrimSpace(listaNomi)
	creaFila(g, listaNomi)
}

func sottostringaMassima(s1 string, s2 string) {

	fmt.Println(calcolaSottosequenzaMassima(s1, s2))

}

func indiceCacofonia(g gioco, sigma string) {

	indice := 0
	if l, inFila := g.FilaPerMattoncino[sigma]; inFila {

		p := l.testa
		nomeMattoncinoPrec := p.nomeMattoncino
		nomeMattoncinoCur := ""
		p = p.next
		for p != nil {
			nomeMattoncinoCur = p.nomeMattoncino
			indice += len(calcolaSottosequenzaMassima(nomeMattoncinoPrec, nomeMattoncinoCur))
			nomeMattoncinoPrec = nomeMattoncinoCur
			p = p.next
		}
		fmt.Println(indice)
	}
	return
}

func costo(g gioco, sigma string, seqForme string) {

	mattonciniDisponibili := map[string][]string{}

	var formaSx, formaDx, nome string

	seqForme1 := []string{}
	seqForme2 := strings.Split(seqForme, ",")
	if len(seqForme2) == 1 {
		return
	}
	if f, inTavolo := g.FilaPerMattoncino[sigma]; inTavolo {
		p := f.testa
		seqForme1 = append(seqForme1, g.mattoncinoPerNome[p.nomeMattoncino].formaBordoSx)
		for p != nil {
			nome = p.nomeMattoncino
			formaSx = g.mattoncinoPerNome[nome].formaBordoSx
			formaDx = g.mattoncinoPerNome[nome].formaBordoDx

			if p.flipped {
				seqForme1 = append(seqForme1, formaSx)
			} else {
				seqForme1 = append(seqForme1, formaDx)
			}

			//aggiungo il nome del mattoncino alla mappa.
			if formaSx < formaDx {
				mattonciniDisponibili[formaSx+"+"+formaDx] = append(mattonciniDisponibili[formaSx+"+"+formaDx], nome)
			} else {
				mattonciniDisponibili[formaDx+"+"+formaSx] = append(mattonciniDisponibili[formaDx+"+"+formaSx], nome)

			}
			p = p.next

		}
	} else {
		return
	}

	//aggiungo i mattoncini della scatola a mattonciniDisponibili
	for k, v := range g.mattoncinoPerNome {

		if _, inTavolo := g.FilaPerMattoncino[k]; !inTavolo {
			formaSx = v.formaBordoSx
			formaDx = v.formaBordoDx
			if formaSx < formaDx {
				mattonciniDisponibili[formaSx+"+"+formaDx] = append(mattonciniDisponibili[formaSx+"+"+formaDx], k)
			} else {
				mattonciniDisponibili[formaDx+"+"+formaSx] = append(mattonciniDisponibili[formaDx+"+"+formaSx], k)

			}

		}
	}

	key := seqForme2[0]
	//controllo che data la sequenza di forme in input sia possibile costruire la fila utilizzando solo i mattoncini disponibili
	for i := 1; i < len(seqForme2); i++ {

		if key > seqForme2[i] {
			key = seqForme2[i] + "+" + key
		} else {
			key += "+" + seqForme2[i]

		}

		if v, disponibile := mattonciniDisponibili[key]; disponibile {

			if len(v) == 1 {
				delete(mattonciniDisponibili, key)
			} else {
				mattonciniDisponibili[key] = v[:len(v)-1]
			}

		} else { //caso in cui non c'è un mattonicno disponibile e quindi non si può costruire la fila
			fmt.Println("indefinito")
			return
		}

		key = seqForme2[i]

	}

	//arrivati qui, è sicuramente possibile calcolare un costo
	i := sottoSequnzaMaxDiMattonciniInComune(seqForme1, seqForme2)
	costo := (len(seqForme1) - 1 - i) + (len(seqForme2) - 1 - i) //sotraggo 1 perchè numero di mattoncini=numero forme-1
	fmt.Println(costo)

}
