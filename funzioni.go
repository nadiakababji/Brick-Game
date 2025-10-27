package main

import (
	"fmt"
	"strings"
)

// restituisce la stringa di un mattoncino nel formato σ : α, β
func (m mattoncino) String() string {
	return fmt.Sprintf("%s: %s, %s", m.nome, m.formaBordoSx, m.formaBordoDx)
}

// verifica se un mattoncino è in scatola
func inScatola(g gioco, sigma string) bool {
	if _, inGioco := g.mattoncinoPerNome[sigma]; !inGioco {
		return false
	}
	if _, inTavolo := g.FilaPerMattoncino[sigma]; inTavolo {
		return false
	}
	return true
}

// *** Questa funzione viene chiamata solo quando si è sicuri che la lista di nomi
// *** contiene solo mattoncini presenti nella scatola e che ogni coppia di mattoncini si incastrano tra di loro
func creaFila(g gioco, listaNomi string) {

	var l filaConCoda

	input := strings.Split(listaNomi, " ")

	sigma := ""
	for i := 0; i < len(input); i++ {

		sigma = input[i][1:]
		AggiungiNodoInCoda(&l, sigma, input[i][0] == '-')
		g.FilaPerMattoncino[sigma] = &l

	}

	return

}

func bfs(grafo map[string][]string, bordiToMattoncino map[string][]string, alpha string, beta string, exlude bool, limiteProfondità int) string { //restituisce la sequenza e la sua lunghezza

	numArchiTraNodi := map[string]int{} //dato il nome del mattoncino mi dice se è stato usato, serve per escludere dei mattoncini senza cambiare la mappa bordiToMattoncini
	for k, v := range bordiToMattoncino {
		numArchiTraNodi[k] = len(v)
	}

	pred := map[string]string{} //dato un bordo mi dice il bordo precedente
	queue := &Queue{}           //*visited appena si aggiunge alla queue

	key := ""
	if exlude {
		key = alpha + "+" + beta
		if beta < alpha {
			key = beta + "+" + alpha
		}
		if _, ok := numArchiTraNodi[key]; ok {
			if numArchiTraNodi[key] == 1 {
				delete(numArchiTraNodi, key)
			} else {
				numArchiTraNodi[key]--
			}
		}
		exlude = false
	}

	pred[alpha] = ""
	queue.Enqueue(alpha, 0)

	nodoCurr := ""
	profondità := 0

	for !queue.IsEmpty() {

		nodoCurr, profondità = queue.Dequeue()
		if profondità >= limiteProfondità {

			return ""
		}

		//trovato nodo finale
		if nodoCurr == beta {
			break

		}

		for _, vicino := range grafo[nodoCurr] { //aggiungo tutti i vicini ,se escludi=true e beta è presente tra i vicini non lo metto

			key = nodoCurr + "+" + vicino
			if vicino < nodoCurr {
				key = vicino + "+" + nodoCurr
			}

			//aggiungo vicino alla coda se c'è ancora una arco che collega nodo corrente con il vicino
			if _, ok := numArchiTraNodi[key]; ok {

				if numArchiTraNodi[key] == 1 {
					delete(numArchiTraNodi, key)
				} else {
					numArchiTraNodi[key]--
				}
				queue.Enqueue(vicino, profondità+1)
				if _, ok := pred[vicino]; !ok { //aggiorno il pred del vicino solo se esso non era già stato considerato da un altro nodo
					pred[vicino] = nodoCurr
				}
				continue
			}

		}

	}

	//caso in cui il loop è terminato perchè non ci sono più elementi nella queue e non si è arrivati al nodo beta
	if _, ok := pred[beta]; !ok {
		return ""

	}

	//costruire la sequenza di forme
	sequenza := beta
	curr := beta
	for curr != alpha {
		sequenza = pred[curr] + "+" + sequenza
		curr = pred[curr]
	}
	return sequenza
}

func calcolaSottosequenzaMassima(s1 string, s2 string) string {
	m := len(s1) + 1
	n := len(s2) + 1

	matrice := make([][]int, m)
	for i := range matrice {
		matrice[i] = make([]int, n)
	}

	for i := 1; i < m; i++ {
		for j := 1; j < n; j++ {
			if s1[i-1] == s2[j-1] {
				matrice[i][j] = matrice[i-1][j-1] + 1
			} else {
				matrice[i][j] = max(matrice[i-1][j], matrice[i][j-1])
			}
		}
	}

	// La matrice contiene la lunghezza della sottosequenza massima tra i prefissi di s1 e s2.
	seq := ""

	//ricostruzione sequenza
	i, j := m-1, n-1
	for i > 0 && j > 0 {
		if s1[i-1] == s2[j-1] {
			seq = string(s1[i-1]) + seq
			i--
			j--
		} else if matrice[i-1][j] > matrice[i][j-1] {
			i--
		} else {
			j--
		}
	}
	return seq

}

// restituisce il massimo tra due numeri.
func max(a, b int) int {
	if a > b {
		return a
	}
	return b
}

func sottoSequnzaMaxDiMattonciniInComune(a1 []string, a2 []string) int {

	m := len(a1) //escludo l'ultima parola di a1 e di a2
	n := len(a2)
	matrice := make([][]int, m)
	for i := range matrice {
		matrice[i] = make([]int, n)
	}

	for i := 1; i < m; i++ {

		for j := 1; j < n; j++ {
			if a1[i-1]+"+"+a1[i] == a2[j-1]+"+"+a2[j] || a1[i]+"+"+a1[i-1] == a2[j-1]+"+"+a2[j] { //aggiungo il + sempre perchè senza considerebbe a bc come ab c
				matrice[i][j] = matrice[i-1][j-1] + 1
			} else {
				matrice[i][j] = max(matrice[i-1][j], matrice[i][j-1])
			}
		}
	}

	return matrice[m-1][n-1]
}
