package main

import "fmt"

// elemento nella fila
type filaNode struct {
	nomeMattoncino string
	flipped        bool
	next           *filaNode
}

type filaConCoda struct {
	testa, coda *filaNode
}

// inserisce nuovo elemento nella fila
func nuovoNodo(sigma string, flipped bool) *filaNode {
	return &filaNode{sigma, flipped, nil}
}

// aggiunge un elemento alla fine della fila
func AggiungiNodoInCoda(l *filaConCoda, sigma string, flipped bool) {
	if l == nil {
		l = &filaConCoda{testa: nuovoNodo(sigma, flipped), coda: nuovoNodo(sigma, flipped)}
		return
	}
	if l.coda == nil {
		l.coda = nuovoNodo(sigma, flipped)
		l.testa = l.coda
	} else {
		l.coda.next = nuovoNodo(sigma, flipped)
		l.coda = l.coda.next
	}
}

// stampa la fila nel formato richiesto
func stampaFilaConCoda(g gioco, l filaConCoda) {
	fmt.Println("(")
	p := l.testa
	for p != nil {
		if p.flipped {
			m := g.mattoncinoPerNome[p.nomeMattoncino]
			fmt.Println(mattoncino{m.formaBordoDx, m.formaBordoSx, m.nome})
		} else {
			fmt.Println(g.mattoncinoPerNome[p.nomeMattoncino])
		}
		p = p.next
	}
	fmt.Println(")")
}
