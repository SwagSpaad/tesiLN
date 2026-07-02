package simulation

import (
	"encoding/csv"
	"fmt"
	"lightning-network/internal/analyzer"
	"math/rand"
	"os"
	"strconv"
	"time"

	"gonum.org/v1/gonum/graph/topo"
)

// la funzione Paga prende in input un nodo u, il nodo v di arrivo e la quantità di satoshi da inviare
// restituisce True se il pagamento va andato a buon fine, False altrimenti, il numero di hop effettuati per arrivare a destinazione
// true se esiste un percorso tra u e v (mancanza liquidità), false se non esiste un percorso tra u e v (path non esistente)
func Paga(nodoMittente, nodoDestinatario int64, pagamento float64, lng *analyzer.LNGraph) (bool, int64, bool) {
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
	hop := int64(0)

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
		for startNodeIndex == destNodeIndex {
			destNodeIndex = rand.Intn(numNodi)
		}
		nodoSorg := IDArray[startNodeIndex]
		nodoDest := IDArray[destNodeIndex]

		esito, hops, pathExists := Paga(nodoSorg, nodoDest, pagamento, lng)

		if esito {
			pagamentiRiusciti++
			totaleHops += int(hops)
		} else {
			pagamentiFalliti++
			if pathExists {
				fallimentoLiquidita++
			} else {
				fallimentoPath++
			}
		}
		if (i+1)%(numPagamenti/10) == 0 {
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

func EsportaCSV(lng *analyzer.LNGraph, nodiDaRimuovere []int64, totRimozioni int, numPagamenti int, valPagamento float64, nomeFile string) error {
	file, err := os.Create(nomeFile)
	if err != nil {
		return fmt.Errorf("errore durante la creazione del file: %v", err)
	}

	defer file.Close()

	writer := csv.NewWriter(file)
	defer writer.Flush()

	// Intestazione CSV
	header := []string{"NodiRimossi", "Fallimento per liquidità", "Fallimento per path inesistente"}
	err = writer.Write(header)
	if err != nil {
		return fmt.Errorf("Errore durante la scrittura dell'intestazione: %v", err)
	}

	fmt.Printf("\n--- AVVIO EXPORT CSV (%s) ---\n", nomeFile)
	for stato := 0; stato <= totRimozioni; stato++ {
		analyzer.RipristinaBilanci(lng)
		percFailLiquidità, percFailPath := calcolaFallimento(lng, numPagamenti, valPagamento)

		riga := []string{
			strconv.Itoa(stato),
			fmt.Sprintf("%.2f", percFailLiquidità),
			fmt.Sprintf("%.2f", percFailPath),
		}

		writer.Write(riga)
		fmt.Printf("Step %d di %d. Nodi rimossi: %d | Perc fail liquidita %.2f%% | Perc fail path %.2f%%\n", stato, totRimozioni, stato, percFailLiquidità, percFailPath)

		if stato < totRimozioni && stato < len(nodiDaRimuovere) {
			nodoDaRimuovere := nodiDaRimuovere[stato]
			lng.Graph.RemoveNode(nodoDaRimuovere)
		}
	}

	return err
}

func calcolaFallimento(lng *analyzer.LNGraph, numPagamenti int, valPagamento float64) (float64, float64) {
	nodi := lng.Graph.Nodes()
	IDArray := []int64{}

	for nodi.Next() {
		IDArray = append(IDArray, nodi.Node().ID())
	}

	numNodi := len(IDArray)
	if numNodi < 2 {
		return 0, 100.00
	}
	rand := rand.New(rand.NewSource(time.Now().UnixNano()))

	connComp := topo.ConnectedComponents(lng.Graph)
	mappaComp := make(map[int64]int)

	for idComp, comp := range connComp {
		for _, nodo := range comp {
			mappaComp[nodo.ID()] = idComp
		}
	}
	fallimentoLiquidita := 0
	fallimentoPath := 0

	for i := 0; i < numPagamenti; i++ {
		startNodeIndex := rand.Intn(numNodi)
		destNodeIndex := rand.Intn(numNodi)
		for startNodeIndex == destNodeIndex {
			destNodeIndex = rand.Intn(numNodi)
		}
		nodoSorg := IDArray[startNodeIndex]
		nodoDest := IDArray[destNodeIndex]

		if mappaComp[nodoSorg] != mappaComp[nodoDest] {
			fallimentoPath++
			continue
		}

		esito, _, _ := Paga(nodoSorg, nodoDest, valPagamento, lng)

		if !esito {
			fallimentoLiquidita++
		}
	}

	percFailLiquidità := float64(fallimentoLiquidita) / float64(numPagamenti) * 100
	percFailPath := float64(fallimentoPath) / float64(numPagamenti) * 100

	return percFailLiquidità, percFailPath
}
