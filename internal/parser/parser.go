package parser

import (
	"encoding/json"
	"os"
)

type NodeData struct {
	PubKey string `json:"pub_key"`
	Alias  string `json:"alias"`
}

type EdgeData struct {
	Node1pub string  `json:"node1_pub"`
	Node2pub string  `json:"node2_pub"`
	Capacity float64 `json:"capacity,string"` // ,string rende possibile la conversione da string a float64
}

func LoadJSON(filename string) ([]NodeData, []EdgeData, error) {
	var nodes []NodeData //slice che contiene i nodi del grafo JSON
	var edges []EdgeData //slice che contiene gli archi del grafo JSON

	file, err := os.Open(filename) // apertura del file
	if err != nil {
		return nil, nil, err //se avviene errore ritornalo
	}
	defer file.Close() // posticipa la chiusura del file dopo un qualsiasi return

	decoder := json.NewDecoder(file) // inizializza il decoder

	_, err = decoder.Token() // consuma il primo carattere del JSON, ovvero '{'
	if err != nil {
		return nil, nil, err
	}

	for decoder.More() { // finchè esiste qualcosa da leggere
		t, err := decoder.Token() // consuma il prossimo elemento
		if err != nil {
			break
		}

		key := t.(string) // associa l'elemento consumato a key e rendilo stringa

		if key == "nodes" {
			decoder.Token()      // consuma l'elemento '['
			for decoder.More() { // se esiste qualcosa da leggere
				var node NodeData
				// leggi l'oggetto JSON e mettilo nei campi di node.
				// se non genera errore (quindi esiste un nodo), inseriscilo nello slice
				if err := decoder.Decode(&node); err == nil {
					nodes = append(nodes, node)
				}
			}
			decoder.Token() // consuma l'elemento ']'
		} else if key == "edges" {
			decoder.Token()      // consuma l'elemento '['
			for decoder.More() { // se esiste qualcosa da leggere
				var edge EdgeData
				// leggi l'oggetto JSON e mettilo nei campi di node.
				// se non genera errore (quindi esiste un arco), inseriscilo nello slice
				if err := decoder.Decode(&edge); err == nil {
					edges = append(edges, edge)
				}
			}
			decoder.Token() // consuma l'elemento ']'
		}
	}
	return nodes, edges, nil
}
