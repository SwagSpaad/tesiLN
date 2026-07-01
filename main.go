package main

import (
	"fmt"
	"lightning-network/internal/analyzer"
	"lightning-network/internal/parser"
	"lightning-network/internal/simulation"
	"log"
	"math/rand"
	"time"
)

func main() {
	start := time.Now()
	fmt.Println("Analizzo grafo lightning network...")
	fmt.Println("Lettura file JSON...")
	nodi, archi, err := parser.LoadJSON("data/graph.json")

	if err != nil {
		log.Fatalf("Errore durante la lettura del file: %v", err)
	}
	grafo := analyzer.BuildGraph(nodi, archi)
	erdosRenyiGraph := analyzer.ErdosRenyiGraph(grafo.Graph.Nodes().Len(), 0.0003, grafo)
	erdosRenyiGraph.ComponentiConnesse()
	/*
		err = analyzer.ExportGraphDot(grafo, "graph_view.dot", 50418)
		if err != nil {
			log.Fatalf("Errore durante l'esportazione: %v", err)
		} else {
			fmt.Println("File graph_view.dot generato correttamente.")
		}
	*/

	fmt.Println("Calcolo il grado medio del grafo...")
	avgDegree := analyzer.AvgDegree(grafo)
	fmt.Printf("Il grado medio del grafo è: %.2v\n", avgDegree)

	Simulazione(1000, 1000, grafo, 50, "LN")
	Simulazione(10000, 1000, erdosRenyiGraph, 50, "ER")

	fmt.Printf("Termino.\n")
	tempoTotale := time.Since(start)
	fmt.Printf("Tempo trascorso: %v\n", tempoTotale)
}

func Simulazione(numPagamenti int, pagamento float64, lng *analyzer.LNGraph, nodiRimossiPerSimulazione int, targetGraph string) {
	rand := rand.New(rand.NewSource(time.Now().UnixNano()))
	lng.ComponentiConnesse()

	centralityMap := simulation.CalcolaBetweennessCentrality(lng)
	hubOrdinati := simulation.OrdinaPerCentralita(centralityMap)
	numNodi := len(hubOrdinati)

	fmt.Printf("\n--- Simulazione n.1: grafo %s intatto ---\n", targetGraph)
	simulation.RandomProcess(numPagamenti, pagamento, lng)
	fmt.Printf("\n--- Numero nodi: %d ---\n", lng.Graph.Nodes().Len())
	analyzer.RipristinaBilanci(lng)

	lngCopy := analyzer.CopiaGrafo(lng)
	fmt.Printf("\n--- Simulazione n.2: %d guasti casuali su %s ---\n", nodiRimossiPerSimulazione, targetGraph)
	for i := 0; i < nodiRimossiPerSimulazione; i++ {
		randIndex := rand.Intn(numNodi)
		nodeID := hubOrdinati[randIndex]
		lngCopy.Graph.RemoveNode(nodeID)
	}
	simulation.RandomProcess(numPagamenti, pagamento, lngCopy)
	lngCopy.ComponentiConnesse()
	fmt.Printf("\n--- Numero nodi: %d ---\n", lngCopy.Graph.Nodes().Len())

	analyzer.RipristinaBilanci(lng)
	lngCopy = analyzer.CopiaGrafo(lng)
	fmt.Printf("\n--- Simulazione n.3: %d attacchi mirati su %s ---\n", nodiRimossiPerSimulazione, targetGraph)
	for _, nodeID := range hubOrdinati[0:nodiRimossiPerSimulazione] {
		lngCopy.Graph.RemoveNode(nodeID)
	}
	simulation.RandomProcess(numPagamenti, pagamento, lngCopy)
	lngCopy.ComponentiConnesse()
	fmt.Printf("\n--- Numero nodi: %d ---\n", lngCopy.Graph.Nodes().Len())
}
