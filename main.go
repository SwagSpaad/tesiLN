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
	rand := rand.New(rand.NewSource(time.Now().UnixNano()))
	fmt.Println("Analizzo grafo lightning network...")
	fmt.Println("Lettura file JSON...")
	nodi, archi, err := parser.LoadJSON("data/graph.json")

	if err != nil {
		log.Fatalf("Errore durante la lettura del file: %v", err)
	}
	grafo := analyzer.BuildGraph(nodi, archi)

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

	_, giantComp := grafo.ComponentiConnesse()
	centralityMap := simulation.CalcolaBetweennessCentrality(giantComp)

	hubOrdinati := simulation.OrdinaPerCentralita(centralityMap)
	numNodi := len(hubOrdinati)

	nodiRimossiPerSimulazione := 50

	fmt.Printf("\n--- Simulazione n.1: grafo intatto ---\n")
	simulation.RandomProcess(10000, 1000, giantComp)
	giantComp.ComponentiConnesse()

	giantCompCopy := giantComp
	fmt.Printf("\n--- Simulazione n.2: %d guasti casuali ---\n", nodiRimossiPerSimulazione)
	for i := 0; i < nodiRimossiPerSimulazione; i++ {
		randIndex := rand.Intn(numNodi)
		nodeID := hubOrdinati[randIndex]
		giantCompCopy.Graph.RemoveNode(nodeID)
	}
	simulation.RandomProcess(10000, 1000, giantCompCopy)
	giantCompCopy.ComponentiConnesse()

	giantCompCopy = giantComp
	fmt.Printf("\n--- Simulazione n.2: %d attacchi mirati ---\n", nodiRimossiPerSimulazione)
	for _, nodeID := range hubOrdinati[0:nodiRimossiPerSimulazione] {
		giantCompCopy.Graph.RemoveNode(nodeID)
	}
	simulation.RandomProcess(10000, 1000, giantCompCopy)
	giantCompCopy.ComponentiConnesse()

	fmt.Printf("Termino.\n")
	tempoTotale := time.Since(start)
	fmt.Printf("Tempo trascorso: %v\n", tempoTotale)
}
