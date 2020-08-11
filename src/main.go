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

package src

import (
	"fmt"
	"net/http"
	"strconv"

	log "github.com/sirupsen/logrus"
	"github.com/gorilla/mux"
	"time"
	"circlemaster"
)


var NetworkManager = circlemaster.NewNetworkManager()
var Manager = circlemaster.NewGraphManager()

func main() {

	router := mux.NewRouter().StrictSlash(true)
	router.HandleFunc(
		"/network/{networkId}/circlemaster/{objFunction}/{dest:[0-9]+}/edges/{from:[0-9]+}/{to:[0-9]+}/request",
			RequestEdge)
	router.HandleFunc(
		"/network/{networkId}/circlemaster/{objFunction}/{dest:[0-9]+}/edges/{from:[0-9]+}/{to:[0-9]+}/release",
			ReleaseEdge)

	log.Fatal(http.ListenAndServe(":8080", router))
}

func sendResult(writer http.ResponseWriter, expTime time.Time) {
	fmt.Fprintf(writer, "%d", expTime.Unix())
}

func RequestEdge(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	networkId := vars["networkId"]
	objFunction := vars["objFunction"]
	dest, _ := strconv.Atoi(vars["dest"])
	from, _ := strconv.Atoi(vars["from"])
	to, _ := strconv.Atoi(vars["to"])

	log.Infof("Got a Request Edge (NetworkId=%s, ObjFunction=%s). Destination = %d | Edge(%d -> %d)\n",
		networkId, objFunction, from, to, dest)

	gm := NetworkManager.GetGraphManager(networkId)
	result := gm.RequestEdge(objFunction, dest, from, to)

	log.Infof("Granted Edge Permission with expiration to: %v. " +
		"(NetworkId=%s, ObjFunction=%s). Destination = %d | Edge(%d -> %d)\n", result)

	sendResult(w, result)
}

func ReleaseEdge(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	networkId := vars["networkId"]
	objFunction := vars["objFunction"]
	dest, _ := strconv.Atoi(vars["dest"])
	from, _ := strconv.Atoi(vars["from"])
	to, _ := strconv.Atoi(vars["to"])

	log.Infof("Got a Release Edge (NetworkId=%s, ObjFunction=%s). Destination = %d | Edge(%d -> %d)\n",
		networkId, objFunction, from, to, dest)

	gm := NetworkManager.GetGraphManager(networkId)
	gm.ReleaseEdge(objFunction, dest, from, to)

	log.Infof("Released Edge (NetworkId=%s, ObjFunction=%s). Destination = %d | Edge(%d -> %d)\n",
		networkId, objFunction, from, to, dest)

	w.WriteHeader(http.StatusOK)
}
