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
	"testing"
	"time"
)

const ObjFunc1 = "ObjFunc1"
const ObjFunc2 = "ObjFunc2"
const Destination1 = 1
const Destination2 = 2
const Destination3 = 3
const Destination4 = 4

/* All those tests assumes the following 4-node topology
 * 1 <---> 2
 * 1 <---> 4
 * 2 <---> 3
 * 2 <---> 4
 * 3 <---> 4
 */

func validTime(expTime time.Time) bool {
	return !expTime.Equal(time.Time{})
}

func TestRequestOneEdge(t *testing.T) {
	gm := NewGraphManager()
	res := gm.RequestEdge(ObjFunc1, Destination3, 1, 2)
	if !validTime(res) {
		t.Errorf("Requesting first edge failed")
	}
}

func TestRequestSameLinkTwice(t *testing.T) {
	gm := NewGraphManager()
	res := gm.RequestEdge(ObjFunc1, Destination3, 1, 2)
	if !validTime(res) {
		t.Errorf("Requesting 1->2 edge failed")
	}

	resSecond := gm.RequestEdge(ObjFunc1, Destination3, 1, 2)
	if !validTime(resSecond) {
		t.Errorf("Requesting 1->2 edge again failed")
	}
}

func TestRequestSameLinkInBothDirections(t *testing.T) {
	gm := NewGraphManager()
	res := gm.RequestEdge(ObjFunc1, Destination3, 1, 2)
	if !validTime(res) {
		t.Errorf("Requesting 1->2 edge failed")
	}

	resSecond := gm.RequestEdge(ObjFunc1, Destination3, 2, 1)
	if validTime(resSecond) {
		t.Errorf("Requesting 2->1 edge succeeded after requesting 1->2")
	}
}

func TestTriangle124LoopWhenSendingTo3(t *testing.T) {
	gm := NewGraphManager()
	res := gm.RequestEdge(ObjFunc1, Destination3, 1, 2)
	if !validTime(res) {
		t.Errorf("Requesting 1->2 edge failed")
	}
	res = gm.RequestEdge(ObjFunc1, Destination3, 2, 4)
	if !validTime(res) {
		t.Errorf("Requesting 2->4 edge failed")
	}
	res = gm.RequestEdge(ObjFunc1, Destination3, 4, 1)
	if validTime(res) {
		t.Errorf("Requesting 2->4 edge succeeded. It would form a loop")
	}
}

func TestDifferentObjectiveFunctions(t *testing.T) {
	gm := NewGraphManager()
	res := gm.RequestEdge(ObjFunc1, Destination3, 1, 2)
	if !validTime(res) {
		t.Errorf("Requesting 1->2 edge failed")
	}

	resSecond := gm.RequestEdge(ObjFunc2, Destination3, 2, 1)
	if !validTime(resSecond) {
		t.Errorf("Requesting 2->1 edge failed on a different objective function")
	}
}

func TestDifferentDestinations(t *testing.T) {
	gm := NewGraphManager()
	res := gm.RequestEdge(ObjFunc1, Destination3, 1, 2)
	if !validTime(res) {
		t.Errorf("Requesting 1->2 edge failed")
	}

	resSecond := gm.RequestEdge(ObjFunc2, Destination4, 2, 1)
	if !validTime(resSecond) {
		t.Errorf("Requesting 2->1 edge failed on a different destination")
	}
}

func TestTriangle124LoopWhenSendingTo3WithRelease(t *testing.T) {
	gm := NewGraphManager()
	res := gm.RequestEdge(ObjFunc1, Destination3, 1, 2)
	if !validTime(res) {
		t.Errorf("Requesting 1->2 edge failed")
	}
	res = gm.RequestEdge(ObjFunc1, Destination3, 2, 4)
	if !validTime(res) {
		t.Errorf("Requesting 2->4 edge failed")
	}
	res = gm.RequestEdge(ObjFunc1, Destination3, 4, 1)
	if validTime(res) {
		t.Errorf("Requesting 2->4 edge succeeded. It would form a loop")
	}

	gm.ReleaseEdge(ObjFunc1, Destination3, 2, 4)
	res = gm.RequestEdge(ObjFunc1, Destination3, 4, 1)
	if !validTime(res) {
		t.Errorf("Requesting 2->4 edge did not succeeded even after releasing 2->4")
	}
}

func TestNetworkManager(t *testing.T) {
	nm := NewNetworkManager()
	gm1 := nm.GetGraphManager("net1")
	gm2 := nm.GetGraphManager("net2")
	if gm1 == gm2 {
		t.Errorf("Should be 2 different GraphManager")
	}

	gm3 := nm.GetGraphManager("net1")
	if gm1 != gm3 {
		t.Errorf("Should be the same GraphManager")
	}
}
