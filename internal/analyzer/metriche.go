package analyzer

import (
	"fmt"
	"os"

	"gonum.org/v1/gonum/graph"
	"gonum.org/v1/gonum/graph/network"
	"gonum.org/v1/gonum/graph/path"
	"gonum.org/v1/gonum/graph/simple"
	"gonum.org/v1/gonum/graph/topo"
)

func AvgDegree(lng *LNGraph) float64 {
	archi := lng.Graph.Edges().Len()
	nodi := lng.Graph.Nodes().Len()

	avgDegree := (2 * archi) / nodi
	return float64(avgDegree)
}

func (lng *LNGraph) ComponentiConnesse() (error, *LNGraph) {
	fmt.Println("Calcolo il numero di componenti connesse...")
	componentiConnesse := topo.ConnectedComponents(lng.Graph)
	numeroComponenti := len(componentiConnesse)

	fmt.Printf("Ci sono %v componenti connesse nel grafo\n", numeroComponenti)

	// genero il file contenente il numero di componenti, il numero di nodi per ogni componente e il diametro
	file, err := os.Create("NodiPerComponente.md")
	if err != nil {
		return err, nil
	}
	defer file.Close()

	giantComp := []graph.Node{}
	count := 1
	for _, c := range componentiConnesse {
		if len(c) > len(giantComp) {
			giantComp = c
		}
		diametro := lng.CalcoloDiametro(c)
		line := fmt.Sprintf("Componente %v -- nodi: %v -- diametro: %v\n", count, len(c), diametro)
		if _, err = file.WriteString(line); err != nil {
			return err, nil
		}
		count++
	}
	giantCompGraph := lng.GeneraSottografo(giantComp)
	return err, giantCompGraph
}

func (lng *LNGraph) TotCapacita() int {
	totCapacita := 0
	archi := lng.Graph.Edges()
	for archi.Next() {
		arco := archi.Edge().(LightningEdge)
		totCapacita += int(arco.Capacity)
	}
	return totCapacita
}

func (lng *LNGraph) GeneraSottografo(nodiComp []graph.Node) *LNGraph {
	sottoGrafo := simple.NewUndirectedGraph()
	pubKeyToId := make(map[string]int64)
	idToPubKey := make(map[int64]string)

	for _, nodo := range nodiComp {
		sottoGrafo.AddNode(nodo)
		pubKeyNodo := lng.IDToPubKey[nodo.ID()]
		pubKeyToId[pubKeyNodo] = nodo.ID()
		idToPubKey[nodo.ID()] = pubKeyNodo
	}

	for _, nodo := range nodiComp {
		adiacenti := lng.Graph.From(nodo.ID())
		for adiacenti.Next() {
			nodoDestinazione := adiacenti.Node()
			// se il nodo esiste nel sottografo e l'id del nodo A è minore del nodo B
			// (serve per controllare una sola volta in tempo costante se esiste un arco tra i due nodi)
			if sottoGrafo.Node(nodoDestinazione.ID()) != nil && nodo.ID() < nodoDestinazione.ID() {
				eData := lng.Graph.EdgeBetween(nodo.ID(), nodoDestinazione.ID())
				sottoGrafo.SetEdge(eData)
			}
		}
	}

	return &LNGraph{
		Graph:      sottoGrafo,
		PubKeyToID: pubKeyToId,
		IDToPubKey: idToPubKey,
	}
}

func (lng *LNGraph) CalcoloDiametro(componente []graph.Node) float64 {
	numNodi := len(componente)
	diametro := 0.0

	if numNodi <= 5000 {
		subGrafo := lng.GeneraSottografo(componente)
		AllPaths := path.DijkstraAllPaths(subGrafo.Graph)
		ecc := network.Eccentricity(subGrafo.Graph, AllPaths)

		for _, dist := range ecc {
			diametro = max(dist, diametro)
		}
	} else {
		//TODO: inserire 2sweepBFS per componenti grandi e utilizzare calcolo preciso per componenti piccole (<500 nodi)
	}
	return diametro
}
