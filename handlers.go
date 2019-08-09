package main

import (
	"bytes"
	"encoding/json"
	"github.com/gorilla/handlers"
	"github.com/gorilla/mux"
	"github.com/prometheus/common/log"
	"github.com/sirupsen/logrus"
	"image/jpeg"
	"io"
	"net/http"
	"strconv"
)

type Handler struct {
	le *logrus.Entry
	s Service
}

func NewHandler(le *logrus.Entry, s Service) Handler {
	return Handler{
		le: le,
		s: s,
	}
}

type ReqBody struct {
	Text string `json:"text"`
}

func (h *Handler) CreateImage(w http.ResponseWriter, r *http.Request) {
	var body io.Reader = r.Body
	dec := json.NewDecoder(body)

	var reqBody ReqBody
	if err := dec.Decode(&reqBody); err != nil {
		h.le.WithError(err).Error("could not decode request body")
		w.WriteHeader(http.StatusBadRequest)
		//fmt.Fprintln(w, fmt.Sprintf(msgTemplate, "could not decode request body"))
		return
	}

	rawImage, err := h.s.CreateImage(reqBody.Text)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	buffer := new(bytes.Buffer)
	if err := jpeg.Encode(buffer, rawImage, nil); err != nil {
		h.le.WithError(err).Error("could not encode image")
	}

	w.Header().Set("Content-Type", "image/jpeg")
	w.Header().Set("Content-Length", strconv.Itoa(len(buffer.Bytes())))
	w.WriteHeader(http.StatusCreated)
	if _, err := w.Write(buffer.Bytes()); err != nil {
		h.le.WithError(err).Error("could not write image")
	}
}

func (h *Handler) RegisterHandlers(router *mux.Router) {
	log.Info("registering handlers")
	createImageHandler := handlers.MethodHandler{
		"POST": http.HandlerFunc(h.CreateImage),
	}

	router.Handle("/create", createImageHandler)
}
