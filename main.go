package main

import (
	"encoding/json"
	"fmt"
	"net"
	"net/http"
	"os"
	"os/signal"
	"syscall"
)

type Request struct {
	Name string
}

type Response struct {
	Err *string
}

type PluginActivateResponse struct {
	Implements []string
}

type CreateRequest struct {
	Request
	Opts map[string]interface{}
}

type MountPathResponse struct {
	Mountpoint string
	Response
}

func main() {
	// As per https://docs.docker.com/engine/extend/plugin_api/
	http.HandleFunc("/Plugin.Activate", func(w http.ResponseWriter, r *http.Request) {
		response, _ := json.Marshal(PluginActivateResponse{
			Implements: []string{"VolumeDriver"},
		})

		w.Header().Set("Content-Type", "application/vnd.docker.plugins.v1+json")
		w.Write(response)
	})

	// As per https://docs.docker.com/engine/extend/plugins_volume/
	http.HandleFunc("/VolumeDriver.Create", func(w http.ResponseWriter, r *http.Request) {
		var request CreateRequest
		var response Response

		decoder := json.NewDecoder(r.Body)
		err := decoder.Decode(&request)
		if err != nil {
			message := "Error Creating: " + err.Error()
			response.Err = &message
		}

		err = os.MkdirAll("/tmp/"+request.Name, 0755)
		if err != nil {
			message := "Error Removing: " + err.Error()
			response.Err = &message
		}

		rawResponse, _ := json.Marshal(response)
		w.Header().Set("Content-Type", "application/vnd.docker.plugins.v1+json")
		w.Write(rawResponse)
	})

	http.HandleFunc("/VolumeDriver.Remove", func(w http.ResponseWriter, r *http.Request) {
		var request Request
		var response Response

		decoder := json.NewDecoder(r.Body)
		err := decoder.Decode(&request)
		if err != nil {
			message := "Error Removing: " + err.Error()
			response.Err = &message
		}

		err = os.RemoveAll("/tmp/" + request.Name)
		if err != nil {
			message := "Error Removing: " + err.Error()
			response.Err = &message
		}

		rawResponse, _ := json.Marshal(response)
		w.Header().Set("Content-Type", "application/vnd.docker.plugins.v1+json")
		w.Write(rawResponse)
	})

	http.HandleFunc("/VolumeDriver.Mount", func(w http.ResponseWriter, r *http.Request) {
		var request Request
		var response MountPathResponse

		decoder := json.NewDecoder(r.Body)
		err := decoder.Decode(&request)
		if err != nil {
			message := "Error Creating: " + err.Error()
			response.Err = &message
		}

		response.Mountpoint = "/tmp/" + request.Name

		rawResponse, _ := json.Marshal(response)
		w.Header().Set("Content-Type", "application/vnd.docker.plugins.v1+json")
		w.Write(rawResponse)
	})

	http.HandleFunc("/VolumeDriver.Path", func(w http.ResponseWriter, r *http.Request) {
		var request Request
		var response MountPathResponse

		decoder := json.NewDecoder(r.Body)
		err := decoder.Decode(&request)
		if err != nil {
			message := "Error Creating: " + err.Error()
			response.Err = &message
		}

		response.Mountpoint = "/tmp/" + request.Name

		rawResponse, _ := json.Marshal(response)
		w.Header().Set("Content-Type", "application/vnd.docker.plugins.v1+json")
		w.Write(rawResponse)
	})

	http.HandleFunc("/VolumeDriver.Unmount", func(w http.ResponseWriter, r *http.Request) {
		var request CreateRequest
		var response Response

		decoder := json.NewDecoder(r.Body)
		err := decoder.Decode(&request)
		if err != nil {
			message := "Error Creating: " + err.Error()
			response.Err = &message
		}

		rawResponse, _ := json.Marshal(response)
		w.Header().Set("Content-Type", "application/vnd.docker.plugins.v1+json")
		w.Write(rawResponse)
	})

	err := os.MkdirAll("/run/docker/plugins", 0755)
	if err != nil {
		fmt.Println("error creating socket directory")
	}

	unixListener, err := net.Listen("unix", "/run/docker/plugins/barebones.sock")
	if err != nil {
		fmt.Println("listen error", err.Error())
		return
	}

	// Handle common process-killing signals so we can gracefully shut down:
	sigc := make(chan os.Signal, 1)
	signal.Notify(sigc, os.Interrupt, os.Kill, syscall.SIGTERM)
	go func(c chan os.Signal) {
		// Wait for a SIGINT or SIGKILL:
		sig := <-c
		fmt.Printf("Caught signal %s: shutting down.\n", sig)
		// Stop listening (and unlink the socket if unix type):
		unixListener.Close()
		// And we're done:
		os.Exit(0)
	}(sigc)

	http.Serve(unixListener, nil)
}
