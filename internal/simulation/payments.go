package simulation

import (
	"fmt"
	"lightning-network/internal/analyzer"
	"math/rand"
	"time"

	"gonum.org/v1/gonum/graph/topo"
)

// la funzione Paga prende in input un nodo u, il nodo v di arrivo e la quantità di satoshi da inviare
// restituisce True se il pagamento va andato a buon fine, False altrimenti, il numero di hop effettuati per arrivare a destinazione
// true se esiste un percorso tra u e v (mancanza liquidità), false se non esiste un percorso tra u e v (path non esistente)
func Paga(nodoMittente, nodoDestinatario int64, pagamento float64, lng *analyzer.LNGraph) (bool, int, bool) {
	// se il mittente e il destinatario sono uguali il pagamento fallisce
	if nodoMittente == nodoDestinatario {
		return false, 0, false
	}

	visitati := make(map[int64]bool) //mappa per i nodi visitati
	//mappe per ricostruire il percorso da destinatario a sorgente
	parentNode := make(map[int64]int64)
	parentEdge := make(map[int64]analyzer.LightningEdge)

	coda := []int64{nodoMittente} //coda inizializzata con il primo nodo
	visitati[nodoMittente] = true
	path := false

	for len(coda) > 0 {
		nodo := coda[0] //estrae il nodo della coda
		coda = coda[1:] //"rimuove" il nodo appena estratto dalla coda

		if nodo == nodoDestinatario {
			path = true
			break
		}

		adiacenti := lng.Graph.From(nodo)
		for adiacenti.Next() {
			next := adiacenti.Node().ID()

			if !visitati[next] { //se il nodo adiacente non è stato ancora visitato
				//recupera l'arco di lng.Graph tra nodo e next e rendilo di tipo LightningEdge
				arco := lng.Graph.EdgeBetween(nodo, next).(analyzer.LightningEdge)

				if arco.SatoshiDisponibili(nodo) >= pagamento {
					visitati[next] = true
					parentNode[next] = nodo
					parentEdge[next] = arco
					coda = append(coda, next)
				}
			}
		}
	}

	if !path {
		if topo.PathExistsIn(lng.Graph, lng.Graph.Node(nodoMittente), lng.Graph.Node(nodoDestinatario)) {
			return false, 0, true
		} else {
			return false, 0, false
		}
	}

	nodo := nodoDestinatario
	hop := 0

	for nodo != nodoMittente {
		//ricostruiamo il percorso da Destinatario a Mittente
		prev := parentNode[nodo]
		channel := parentEdge[nodo]

		//aggiorniamo il bilancio sull'arco utilizzato
		channel.AggiornaBilancio(prev, pagamento)
		nodo = prev
		//incrementiamo numero di hop effettuati
		hop++
	}

	return true, hop, true
}

func RandomProcess(numPagamenti int, pagamento float64, lng *analyzer.LNGraph) {
	fmt.Printf("\n--- AVVIO SIMULAZIONE DI %d PAGAMENTI ---\n", numPagamenti)
	//estraggo tutti gli id dei nodi per la scelta casuale dei nodi
	nodi := lng.Graph.Nodes()
	IDArray := []int64{}

	for nodi.Next() {
		IDArray = append(IDArray, nodi.Node().ID())
	}

	numNodi := len(IDArray)
	if numNodi < 2 {
		fmt.Printf("Numero di nodi nel grafo insufficienti per la simulazione.")
		return
	}
	rand := rand.New(rand.NewSource(time.Now().UnixNano()))

	pagamentiRiusciti := 0
	pagamentiFalliti := 0
	fallimentoLiquidita := 0
	fallimentoPath := 0
	totaleHops := 0

	for i := 0; i < numPagamenti; i++ {
		startNodeIndex := rand.Intn(numNodi)
		destNodeIndex := rand.Intn(numNodi)
		//fintanto che il nodo destinazione corrisponde a quello di partenza
		for startNodeIndex == destNodeIndex {
			destNodeIndex = rand.Intn(numNodi)
		}

		startNodeID := IDArray[startNodeIndex]
		destNodeID := IDArray[destNodeIndex]

		esito, hops, pathExists := Paga(startNodeID, destNodeID, pagamento, lng)
		//se il pagamento va a buon fine
		if esito {
			pagamentiRiusciti++
			totaleHops += hops
		} else {
			pagamentiFalliti++
			if pathExists {
				fallimentoLiquidita++
			} else {
				fallimentoPath++
			}

		}
		if (i+1)%1000 == 0 {
			fmt.Printf("Progresso: %d di %d pagamenti simulati...\n", i+1, numPagamenti)
		}
	}

	fmt.Println("\n--- RISULTATI DELLA SIMULAZIONE ---")
	fmt.Printf("Pagamenti tentati: %d\n", numPagamenti)
	fmt.Printf("Pagamenti Riusciti: %d (%.2f%%)\n", pagamentiRiusciti, float64(pagamentiRiusciti)/float64(numPagamenti)*100)
	fmt.Printf("Pagamenti Falliti: %d (%.2f%%)\n", pagamentiFalliti, float64(pagamentiFalliti)/float64(numPagamenti)*100)
	fmt.Printf("Fallimenti per liquidità insufficiente: %d (%.2f%%)\n", fallimentoLiquidita, float64(fallimentoLiquidita)/float64(numPagamenti)*100)
	fmt.Printf("Fallimenti per percorso non esistente: %d (%.2f%%)\n", fallimentoPath, float64(fallimentoPath)/float64(numPagamenti)*100)

	if pagamentiRiusciti > 0 {
		fmt.Printf("Lunghezza media percorso: %.2f salti\n", float64(totaleHops)/float64(pagamentiRiusciti))
	}
}
