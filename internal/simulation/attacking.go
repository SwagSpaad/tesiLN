package simulation

import (
	"lightning-network/internal/analyzer"
	"sort"

	"gonum.org/v1/gonum/graph/network"
)

type NodoCentralita struct {
	NodeID     int64
	Centralita float64
}

func CalcolaDegreeCentrality(lng *analyzer.LNGraph) map[int64]float64 {
	centralityMap := make(map[int64]float64)
	nodi := lng.Graph.Nodes()
	for nodi.Next() {
		node := nodi.Node()
		centralityMap[node.ID()] = float64(lng.Graph.From(node.ID()).Len())
	}
	return centralityMap
}

func CalcolaBetweennessCentrality(lng *analyzer.LNGraph) map[int64]float64 {
	centralityMap := network.Betweenness(lng.Graph)
	return centralityMap
}

func OrdinaPerCentralita(centralityMap map[int64]float64) []int64 {
	nodeSlice := make([]NodoCentralita, 0, len(centralityMap))

	for id, centr := range centralityMap {
		nodeSlice = append(nodeSlice, NodoCentralita{
			NodeID:     id,
			Centralita: centr,
		})
	}

	sort.Slice(nodeSlice, func(i, j int) bool {
		return nodeSlice[i].Centralita > nodeSlice[j].Centralita
	})

	idOrdinati := make([]int64, len(nodeSlice))
	for i, nodo := range nodeSlice {
		idOrdinati[i] = nodo.NodeID
	}

	return idOrdinati
}
