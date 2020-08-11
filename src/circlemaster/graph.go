/*
  Copyright (c) 2017 The Mode Group

  Licensed under the Apache License, Version 2.0 (the "License");
  you may not use this file except in compliance with the License.
  You may obtain a copy of the License at

      https://www.apache.org/licenses/LICENSE-2.0

  Unless required by applicable law or agreed to in writing, software
  distributed under the License is distributed on an "AS IS" BASIS,
  WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
  See the License for the specific language governing permissions and
  limitations under the License.
*/

package circlemaster

import (
	"sync"
	"time"
)

const MaxNodes = 50
const EdgeTimeout = time.Duration(30 * time.Second)

type Graph struct {
	ObjectiveFunction string
	Destination       int

	Edges [MaxNodes][MaxNodes]time.Time

	mutex sync.Mutex
}

type NetworkManager struct {
	GraphManagers sync.Map // networkId -> graphManager
}

type GraphManager struct {
	Graphs map[string]map[int]*Graph // ObjectiveFunction -> Destination -> Graph
}

func NewNetworkManager() *NetworkManager {
	nm := new(NetworkManager)
	return nm
}

func NewGraphManager() *GraphManager {
	gm := new(GraphManager)
	gm.Graphs = make(map[string]map[int]*Graph)
	return gm
}

func (nm *NetworkManager) GetGraphManager(networkId string) *GraphManager {
	graphmanager, _ := nm.GraphManagers.LoadOrStore(networkId, NewGraphManager())
	return graphmanager.(*GraphManager)
}

func (gm *GraphManager) GetGraph(objectiveFunction string, destination int) *Graph {
	if _, contains := gm.Graphs[objectiveFunction]; !contains {
		gm.Graphs[objectiveFunction] = make(map[int]*Graph)
	}
	destinationMap := gm.Graphs[objectiveFunction]

	if _, contains := destinationMap[destination]; !contains {
		destinationMap[destination] = &Graph{ObjectiveFunction: objectiveFunction, Destination: destination}
	}

	return gm.Graphs[objectiveFunction][destination]
}


func (graph *Graph) RequestEdge(from int, to int) time.Time {
	graph.mutex.Lock()
	now := time.Now().UTC()

	// If edge is free and it forms loops, we don't give permission
	if !graph.edgeIsGranted(from, to, now) && graph.FormsLoop(from, to, now) {
		graph.mutex.Unlock()
		return time.Time{}
	}

	graph.Edges[from][to] = now
	graph.mutex.Unlock()
	return now.Add(EdgeTimeout)
}

func (graph *Graph) ReleaseEdge(from int, to int) {
	reset := time.Time{}
	graph.mutex.Lock()
	graph.Edges[from][to] = reset
	graph.mutex.Unlock()
}

/*
 * Works by trying to find a path from 'to' to 'from', if that exists we would have a loop by adding 'from' to 'to'
 * edge.
 */
func (graph *Graph) FormsLoop(from int, to int, now time.Time) bool {
	nodesToVisit := make([]int, 0, MaxNodes)
	nodesToVisit = append(nodesToVisit, to)
	for len(nodesToVisit) > 0 {
		currentNode := nodesToVisit[0]

		if graph.edgeIsGranted(currentNode, from, now) {
			// we have a loop
			return true
		}

		// Adds the nodes we can reach from currentNode to the visit list.
		for index, _ := range graph.Edges[currentNode] {
			if graph.edgeIsGranted(currentNode, index, now) {
				nodesToVisit = append(nodesToVisit, index)
			}
		}

		nodesToVisit = nodesToVisit[1:]
	}

	return false
}

// Returns true if the given edge was granted to someone.
func (graph *Graph) edgeIsGranted(from int, to int, now time.Time) bool {
	return now.Sub(graph.Edges[from][to]) < EdgeTimeout
}

func (gm *GraphManager) RequestEdge(objectiveFunction string, destination int, from int, to int) time.Time {
	graph := gm.GetGraph(objectiveFunction, destination)
	return graph.RequestEdge(from, to)
}

func (gm *GraphManager) ReleaseEdge(objectiveFunction string, destination int, from int, to int) {
	graph := gm.GetGraph(objectiveFunction, destination)
	graph.ReleaseEdge(from, to)
}

