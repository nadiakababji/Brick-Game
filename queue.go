// ho scelto di rappresentare le queue come liste concatenate con tail perchè ho continua necessita di inserire elementi alla fine della coda e eliminare il primo elemento della coda e
// con la lista concatenata ( a differenza delle slice ) la complessità temporale di tutte le funzioni che mi servono è  = O(1)
package main

import "fmt"

// Struttura dati per un nodo della lista
type Node struct {
	data     string
	distanza int
	next     *Node
}

// Struttura dati per la coda
type Queue struct {
	head *Node
	tail *Node
}

// Inserisci un elemento alla fine della coda
func (q *Queue) Enqueue(item string, dis int) {
	newNode := &Node{data: item, distanza: dis, next: nil}

	if q.head == nil {
		q.head = newNode
		q.tail = newNode
	} else {
		q.tail.next = newNode
		q.tail = newNode
	}
}

// Rimuovi e restituisci l'elemento dalla testa della coda
func (q *Queue) Dequeue() (string, int) {
	if q.head == nil {
		panic("La coda è vuota")
	}

	data := q.head.data
	distanza := q.head.distanza
	q.head = q.head.next

	if q.head == nil {
		q.tail = nil
	}

	return data, distanza
}

// Restituisci true se la coda è vuota
func (q *Queue) IsEmpty() bool {
	return q.head == nil
}

func (q *Queue) PrintQueue() {
	fmt.Print("Contenuto della coda: ")

	p := q.head
	for p != nil {
		fmt.Printf("%s %d ", p.data, p.distanza)
		p = p.next
	}

	fmt.Println()
}
