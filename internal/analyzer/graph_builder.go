package analyzer

import (
	"fmt"
	"lightning-network/internal/parser"
	"math/rand"
	"os"
	"time"

	"gonum.org/v1/gonum/graph"
	"gonum.org/v1/gonum/graph/graphs/gen"
	"gonum.org/v1/gonum/graph/simple"
)

type BalanceState struct {
	BalanceA float64 // bilancio A->B
	BalanceB float64 // bilancio B->A
}

type LightningEdge struct {
	F, T     graph.Node
	Capacity float64
	Balance  *BalanceState //utilizzo un puntatore ad una struct per aggiornamento in tempo costante del bilancio
}

func (e LightningEdge) From() graph.Node {
	return e.F
}
func (e LightningEdge) To() graph.Node {
	return e.T
}
func (e LightningEdge) ReversedEdge() graph.Edge {
	return LightningEdge{
		F:        e.T,
		T:        e.F,
		Capacity: e.Capacity,
		Balance:  e.Balance,
	}
}

func (e LightningEdge) SatoshiDisponibili(nodeID int64) float64 {
	if nodeID == e.F.ID() {
		return e.Balance.BalanceA
	}
	return e.Balance.BalanceB

}

func (e LightningEdge) AggiornaBilancio(nodeID int64, pagamento float64) {
	if nodeID == e.F.ID() {
		e.Balance.BalanceA -= pagamento
		e.Balance.BalanceB += pagamento
	} else {
		e.Balance.BalanceB -= pagamento
		e.Balance.BalanceA += pagamento
	}
}

type LNGraph struct {
	Graph      *simple.UndirectedGraph
	PubKeyToID map[string]int64
	IDToPubKey map[int64]string
}

func BuildGraph(nodi []parser.NodeData, archi []parser.EdgeData) *LNGraph {
	//dati i nodi e gli archi parsati dal file JSON costruisci un grafo di tipo LNGraph
	fmt.Println("Costruzione del grafo con GoNum...")

	graph := simple.NewUndirectedGraph()

	pubKeyToId := make(map[string]int64)
	idToPubKey := make(map[int64]string)

	for _, nData := range nodi {
		node := graph.NewNode()
		graph.AddNode(node)

		pubKeyToId[nData.PubKey] = node.ID()
		idToPubKey[node.ID()] = nData.PubKey
	}

	numArchi := 0
	for _, eData := range archi {
		idNode1, existNode1 := pubKeyToId[eData.Node1pub]
		idNode2, existNode2 := pubKeyToId[eData.Node2pub]

		if existNode1 && existNode2 {
			if idNode1 != idNode2 {
				edge := LightningEdge{
					F:        graph.Node(idNode1), // From
					T:        graph.Node(idNode2), // To
					Capacity: float64(eData.Capacity),
					Balance: &BalanceState{
						BalanceA: float64(eData.Capacity) / 2.0,
						BalanceB: float64(eData.Capacity) / 2.0,
					},
				}

				graph.SetEdge(edge)
				numArchi++
			}
		}
	}

	fmt.Printf("Grafo costruito: %d nodi, %d archi validi elaborati\n", graph.Nodes().Len(), numArchi)

	return &LNGraph{
		Graph:      graph,
		PubKeyToID: pubKeyToId,
		IDToPubKey: idToPubKey,
	}
}

func ErdosRenyiGraph(n int, p float64, sampleGraph *LNGraph) *LNGraph {
	rand := rand.New(rand.NewSource(time.Now().UnixNano()))
	archi := sampleGraph.Graph.Edges()
	EdgeArray := []LightningEdge{}

	for archi.Next() {
		EdgeArray = append(EdgeArray, archi.Edge().(LightningEdge))
	}
	numArchi := len(EdgeArray)

	graph := simple.NewUndirectedGraph()
	gen.Gnp(graph, n, p, nil)
	graphEdges := graph.Edges()
	for graphEdges.Next() {
		edge := graphEdges.Edge()

		randEdgeIndex := rand.Intn(numArchi)
		randEdge := EdgeArray[randEdgeIndex]

		lightningEdge := LightningEdge{
			F:        edge.From(),
			T:        edge.To(),
			Capacity: randEdge.Capacity,
			Balance: &BalanceState{
				BalanceA: randEdge.Capacity / 2.0,
				BalanceB: randEdge.Capacity / 2.0,
			},
		}
		graph.RemoveEdge(edge.From().ID(), edge.To().ID())
		graph.SetEdge(lightningEdge)
	}

	fmt.Printf("Grafo Erdos Renyi generato: %d nodi, %d archi validi elaborati\n", graph.Nodes().Len(), graph.Edges().Len())

	return &LNGraph{
		Graph:      graph,
		PubKeyToID: nil,
		IDToPubKey: nil,
	}
}

func RipristinaBilanci(lng *LNGraph) {
	archi := lng.Graph.Edges()
	for archi.Next() {
		e := archi.Edge().(LightningEdge)
		e.Balance.BalanceA = e.Capacity / 2.0
		e.Balance.BalanceB = e.Capacity / 2.0
	}
}

func CopiaGrafo(grafoOriginale *LNGraph) *LNGraph {
	grafoCopia := simple.NewUndirectedGraph()
	nodi := grafoOriginale.Graph.Nodes()
	for nodi.Next() {
		nodo := nodi.Node()
		grafoCopia.AddNode(nodo)
	}

	archi := grafoOriginale.Graph.Edges()
	for archi.Next() {
		arco := archi.Edge().(LightningEdge)
		grafoCopia.SetEdge(arco)
	}

	pubKeyToId := make(map[string]int64)
	idToPubKey := make(map[int64]string)

	for pubKey, id := range grafoOriginale.PubKeyToID {
		pubKeyToId[pubKey] = id
	}

	for id, pubKey := range grafoOriginale.IDToPubKey {
		idToPubKey[id] = pubKey
	}

	return &LNGraph{
		Graph:      grafoCopia,
		PubKeyToID: pubKeyToId,
		IDToPubKey: idToPubKey,
	}
}

func ExportGraphDot(lng *LNGraph, filename string, maxEdges int) error {
	fmt.Printf("Avvio esportazione di %d archi del grafo in %s...\n", maxEdges, filename)

	file, err := os.Create(filename)
	if err != nil {
		return err
	}
	defer file.Close()

	_, err = file.WriteString("graph LightningNetwork {\n")
	if err != nil {
		return err
	}

	archi := lng.Graph.Edges()
	count := 0

	for archi.Next() {
		if count >= maxEdges {
			break
		}

		e := archi.Edge()

		pub1 := lng.IDToPubKey[e.From().ID()][:6]
		pub2 := lng.IDToPubKey[e.To().ID()][:6]

		line := fmt.Sprintf("    \"%s\" -- \"%s\";\n", pub1, pub2)
		file.WriteString(line)

		count++
	}

	_, err = file.WriteString("}\n")
	return err
}
