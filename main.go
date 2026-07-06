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
	//	rand := rand.New(rand.NewSource(time.Now().UnixNano()))

	nodi, archi, err := parser.LoadJSON("data/graph.json")
	if err != nil {
		log.Fatalf("Errore durante la lettura del file: %v", err)
	}

	const maxRimozioni = 150 //1% nodi giant comp
	const pagamentiPerStep = 10000
	const importoPagamento = 50000.0

	grafoLN := analyzer.BuildGraph(nodi, archi)
	grafoER_GC := analyzer.ErdosRenyiGraph(len(nodi), 0.000336, grafoLN)
	//	grafoER_AM := analyzer.CopiaGrafo(grafoER_GC) // copio il grafo ER per la futura simulazione degli attacchi mirati

	// analyzer.EsportaGradoCSV(grafoLN, "GradoLN.csv")
	// analyzer.EsportaGradoCSV(grafoER_GC, "GradoER.csv")
	grafoER_GC.ComponentiConnesse()
	// fmt.Println("=== AVVIO SIMULAZIONE COMPLETA NO GIANT===...")

	// var wg sync.WaitGroup
	// wg.Add(4)

	// go func() {
	// 	defer wg.Done()
	// 	fmt.Println("=== SCENARIO 1: GRAFO LN - GUASTI CASUALI ===")

	// 	// Simulazione guasti casuali su LN
	// 	grafoLN_GC := analyzer.BuildGraph(nodi, archi)
	// 	//_, giantComp_GC := grafoLN_GC.ComponentiConnesse()
	// 	centralityMap := simulation.CalcolaDegreeCentrality(grafoLN_GC)
	// 	hubOrdinati := simulation.OrdinaPerCentralita(centralityMap)
	// 	nodiDaRimuovere := make([]int64, 0, maxRimozioni)
	// 	numNodi := len(hubOrdinati) //lng.Graph.Nodes().Len()
	// 	for i := 0; i < maxRimozioni; i++ {
	// 		randIndex := rand.Intn(numNodi)
	// 		for hubOrdinati[randIndex] == -1 {
	// 			randIndex = rand.Intn(numNodi)
	// 		}
	// 		nodiDaRimuovere = append(nodiDaRimuovere, hubOrdinati[randIndex])
	// 		hubOrdinati[randIndex] = -1
	// 	}

	// 	err = simulation.EsportaCSV(grafoLN_GC, nodiDaRimuovere, maxRimozioni, pagamentiPerStep, importoPagamento, "LN_GuastiCasuali_noGiant.csv")
	// 	if err != nil {
	// 		log.Fatalf("Errore CSV Guasti Casuali LN no Giant: %v", err)
	// 	}
	// 	fmt.Println("=== FINE SCENARIO 1: CSV GENERATO ===")
	// }()
	// go func() {
	// 	defer wg.Done()
	// 	//Simulazione attacchi mirati su LN
	// 	fmt.Println("=== SCENARIO 2: GRAFO LN - ATTACCHI MIRATI ===")
	// 	grafoLN_AM := analyzer.BuildGraph(nodi, archi)
	// 	//_, giantComp_AM := grafoLN_AM.ComponentiConnesse()
	// 	centralityMap := simulation.CalcolaDegreeCentrality(grafoLN_AM)
	// 	hubOrdinati := simulation.OrdinaPerCentralita(centralityMap)

	// 	err = simulation.EsportaCSV(grafoLN_AM, hubOrdinati, maxRimozioni, pagamentiPerStep, importoPagamento, "LN_AttacchiMirati_noGiant.csv")
	// 	if err != nil {
	// 		log.Fatalf("Errore CSV Attacchi Mirati LN: %v", err)
	// 	}
	// 	fmt.Println("=== FINE SCENARIO 2: CSV GENERATO ===")
	// }()

	// go func() {
	// 	defer wg.Done()
	// 	fmt.Println("=== SCENARIO 3: GRAFO ER - GUASTI CASUALI ===")
	// 	//Simulazione guasti casuali su ER

	// 	//_, giantComp_ER_GC := grafoER_GC.ComponentiConnesse()
	// 	centralityMap := simulation.CalcolaDegreeCentrality(grafoER_GC)
	// 	hubOrdinati := simulation.OrdinaPerCentralita(centralityMap)
	// 	nodiDaRimuovere := make([]int64, 0, maxRimozioni)
	// 	numNodi := len(hubOrdinati) //lng.Graph.Nodes().Len()
	// 	for i := 0; i < maxRimozioni; i++ {
	// 		randIndex := rand.Intn(numNodi)
	// 		for hubOrdinati[randIndex] == -1 {
	// 			randIndex = rand.Intn(numNodi)
	// 		}
	// 		nodiDaRimuovere = append(nodiDaRimuovere, hubOrdinati[randIndex])
	// 		hubOrdinati[randIndex] = -1
	// 	}

	// 	err = simulation.EsportaCSV(grafoER_GC, nodiDaRimuovere, maxRimozioni, pagamentiPerStep, importoPagamento, "ER_GuastiCasuali_noGiant.csv")
	// 	if err != nil {
	// 		log.Fatalf("Errore CSV Guasti Casuali ER: %v", err)
	// 	}
	// 	fmt.Println("=== FINE SCENARIO 3: CSV GENERATO ===")
	// }()
	// go func() {
	// 	defer wg.Done()
	// 	fmt.Println("=== SCENARIO 4: GRAFO ER - ATTACCHI MIRATI ===")
	// 	//Simulazione guasti casuali su ER
	// 	//_, giantComp_ER_AM := grafoER_AM.ComponentiConnesse()
	// 	centralityMap := simulation.CalcolaDegreeCentrality(grafoER_AM)
	// 	hubOrdinati := simulation.OrdinaPerCentralita(centralityMap)
	// 	err = simulation.EsportaCSV(grafoER_AM, hubOrdinati, maxRimozioni, pagamentiPerStep, importoPagamento, "ER_AttacchiMirati_noGiant.csv")
	// 	if err != nil {
	// 		log.Fatalf("Errore CSV Attacchi Mirati ER: %v", err)
	// 	}
	// 	fmt.Println("=== FINE SCENARIO 4: CSV GENERATO ===")
	// }()
	// wg.Wait()
	// fmt.Printf("=== FINE SIMULAZIONE COMPLETA ===\n")
	tempoTotale := time.Since(start)
	fmt.Printf("Tempo trascorso: %v\n", tempoTotale)
	/*
		fmt.Println("Analizzo grafo lightning network...")
		fmt.Println("Lettura file JSON...")
		nodi, archi, err := parser.LoadJSON("data/graph.json")

		if err != nil {
			log.Fatalf("Errore durante la lettura del file: %v", err)
		}
		grafo := analyzer.BuildGraph(nodi, archi)
		grafoErdosRenyi := analyzer.ErdosRenyiGraph(grafo.Graph.Nodes().Len(), 0.0003, grafo)
		grafoErdosRenyi.ComponentiConnesse()

			err = analyzer.ExportGraphDot(grafo, "graph_view.dot", 50418)
			if err != nil {
				log.Fatalf("Errore durante l'esportazione: %v", err)
			} else {
				fmt.Println("File graph_view.dot generato correttamente.")
			}


		fmt.Println("Calcolo il grado medio del grafo...")
		avgDegree := analyzer.AvgDegree(grafo)
		fmt.Printf("Il grado medio del grafo è: %.2v\n", avgDegree)

		Simulazione(100000, 1000, grafo, 50, "LN")
		Simulazione(100000, 1000, grafoErdosRenyi, 50, "ER")

		fmt.Printf("Termino.\n")
		tempoTotale := time.Since(start)
		fmt.Printf("Tempo trascorso: %v\n", tempoTotale)
	*/
}

func Simulazione(numPagamenti int, pagamento float64, grafo *analyzer.LNGraph, nodiRimossiPerSimulazione int, targetGraph string) {
	rand := rand.New(rand.NewSource(time.Now().UnixNano()))
	_, giantComp := grafo.ComponentiConnesse()

	centralityMap := simulation.CalcolaDegreeCentrality(giantComp)
	/*centralityMap := make(map[int64]float64)
	for nodi := lng.Graph.Nodes(); nodi.Next(); {
		node := nodi.Node()
		centralityMap[node.ID()] = 1
	}*/
	hubOrdinati := simulation.OrdinaPerCentralita(centralityMap)
	numNodi := len(hubOrdinati) //lng.Graph.Nodes().Len()

	fmt.Printf("\n--- Simulazione n.1: grafo %s intatto ---\n", targetGraph)
	simulation.RandomProcess(numPagamenti, pagamento, giantComp)
	fmt.Printf("\n--- Numero nodi: %d ---\n", giantComp.Graph.Nodes().Len())
	analyzer.RipristinaBilanci(giantComp)

	lngCopy := analyzer.CopiaGrafo(giantComp)
	fmt.Printf("\n--- Simulazione n.2: %d guasti casuali su %s ---\n", nodiRimossiPerSimulazione, targetGraph)
	for i := 0; i < nodiRimossiPerSimulazione; i++ {
		randIndex := rand.Intn(numNodi)
		for hubOrdinati[randIndex] == -1 {
			randIndex = rand.Intn(numNodi)
		}
		nodeID := hubOrdinati[randIndex]
		lngCopy.Graph.RemoveNode(nodeID)
		hubOrdinati[randIndex] = -1
	}
	simulation.RandomProcess(numPagamenti, pagamento, lngCopy)
	lngCopy.ComponentiConnesse()
	fmt.Printf("\n--- Numero nodi: %d ---\n", lngCopy.Graph.Nodes().Len())

	analyzer.RipristinaBilanci(giantComp)
	lngCopy = analyzer.CopiaGrafo(giantComp)
	fmt.Printf("\n--- Simulazione n.3: %d attacchi mirati su %s ---\n", nodiRimossiPerSimulazione, targetGraph)
	for _, nodeID := range hubOrdinati[0:nodiRimossiPerSimulazione] {
		lngCopy.Graph.RemoveNode(nodeID)
	}
	simulation.RandomProcess(numPagamenti, pagamento, lngCopy)
	lngCopy.ComponentiConnesse()
	fmt.Printf("\n--- Numero nodi: %d ---\n", lngCopy.Graph.Nodes().Len())
}
