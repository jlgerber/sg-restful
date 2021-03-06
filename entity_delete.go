package main

import (
	"bytes"
	"encoding/json"
	"io/ioutil"
	"net/http"
	"strconv"
	"strings"

	log "github.com/Sirupsen/logrus"
	"github.com/gorilla/context"
	"github.com/gorilla/mux"
)

type deleteResponse struct {
	Results   bool   `json:"results"`
	Exception bool   `json:"exception,omitempty"`
	Message   string `json:"message,omitempty"`
	ErrorCode int    `json:"error_code,omitempty"`
}

func entityDeleteHandler(config clientConfig) func(rw http.ResponseWriter, req *http.Request) {
	return func(rw http.ResponseWriter, req *http.Request) {
		vars := mux.Vars(req)
		entityType, ok := vars["entity_type"]
		if !ok {
			log.Errorf("Missing Entity Type")
			rw.WriteHeader(http.StatusBadRequest)
			return
		}

		var entityID int
		var err error
		entityIDStr, ok := vars["id"]
		if ok {
			entityID, err = strconv.Atoi(entityIDStr)
			log.Debugf("Id: %v  Error: %s", entityID, err)
			if err != nil {
				rw.WriteHeader(http.StatusBadRequest)
				return
			}
		} else {
			rw.WriteHeader(http.StatusBadRequest)
			return
		}
		log.Debugf("Entity: %s - %d", entityType, entityID)

		query := map[string]interface{}{
			"type": entityType,
			"id":   entityID,
		}

		sgConn, ok := context.GetOk(req, "sgConn")
		if !ok {
			rw.WriteHeader(http.StatusInternalServerError)
			return
		}
		sg := sgConn.(Shotgun)
		sgReq, err := sg.Request("delete", query)
		if err != nil {
			log.Error("Request Error: ", err)
			rw.WriteHeader(http.StatusInternalServerError)
			return
		}

		var deleteResp deleteResponse
		respBody, err := ioutil.ReadAll(sgReq.Body)
		if err != nil {
			log.Error(err)
			rw.WriteHeader(http.StatusBadGateway)
			return
		}

		log.Debugf("Json Response: %s", respBody)

		err = json.Unmarshal(respBody, &deleteResp)
		if err != nil {
			log.Error(err)
			rw.WriteHeader(http.StatusBadGateway)
			return
		}

		log.Debugf("Response: %v", deleteResp)

		if deleteResp.Exception {
			if strings.Contains(deleteResp.Message, "Permission") {
				rw.WriteHeader(http.StatusForbidden)
			} else if strings.Contains(deleteResp.Message, "does not exist") {
				rw.WriteHeader(http.StatusNotFound)
			} else {
				rw.WriteHeader(http.StatusBadRequest)
			}
			rw.Write(bytes.NewBufferString(deleteResp.Message).Bytes())
			return
		}

		// I'm not sure this can even happen
		if !deleteResp.Results {
			rw.WriteHeader(http.StatusNotFound)
			return
		}

		rw.WriteHeader(http.StatusOK)
	}
}
