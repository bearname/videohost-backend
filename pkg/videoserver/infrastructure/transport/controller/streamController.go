package controller

import (
	"fmt"
	"github.com/bearname/videohost/pkg/common/infrarstructure/transport/controller"
	"github.com/bearname/videohost/pkg/videoserver/domain/repository"
	"github.com/gorilla/mux"
	log "github.com/sirupsen/logrus"
	"net/http"
)

type StreamController struct {
	controller.BaseController
	videoRepository repository.VideoRepository
}

func NewStreamController(videoRepository repository.VideoRepository) *StreamController {
	v := new(StreamController)
	v.videoRepository = videoRepository
	return v
}

func (c *StreamController) StreamHandler(writer http.ResponseWriter, request *http.Request) {
	writer.Header().Set("Access-Control-Allow-Origin", "*")
	writer.Header().Set("Access-Control-Allow-Methods", "POST, GET, OPTIONS, PUT, DELETE")
	writer.Header().Set("Access-Control-Allow-Headers", "Accept, Content-Type, Content-Length, Accept-Encoding, X-CSRF-Token, Authorization")
	if (*request).Method == "OPTIONS" {
		writer.WriteHeader(http.StatusNoContent)
		return
	}
	fmt.Println("readust", request.RequestURI)
	vars := mux.Vars(request)
	var videoId string
	var ok bool

	if videoId, ok = vars["videoId"]; !ok {
		c.BaseController.WriteResponse(writer, http.StatusBadRequest, false, "Failed get videoId")
		return
	}

	video, err := c.videoRepository.Find(videoId)
	if err != nil {
		c.BaseController.WriteResponse(writer, http.StatusBadRequest, false, "Video not found")
		return
	}

	mediaBase := c.getMediaBase(video.Id)
	segName, ok := vars["segName"]
	log.Info("videoId: " + videoId + " segName " + segName)
	if !ok {
		m3u8Name := "index.m3u8"
		log.Info("serveHlsM3u8")
		c.serveHlsM3u8(writer, request, mediaBase, m3u8Name)
	} else {
		log.Info("serveHlsTs")
		c.serveHlsTs(writer, request, mediaBase, segName)
	}
}

func (_ *StreamController) getMediaBase(id string) string {
	return "content\\" + id
}

func (_ *StreamController) serveHlsM3u8(w http.ResponseWriter, r *http.Request, mediaBase, m3u8Name string) {
	mediaFile := fmt.Sprintf("%s\\%s", mediaBase, m3u8Name)
	http.ServeFile(w, r, mediaFile)
	w.Header().Set("Content-Type", "application/x-mpegURL")
}

func (_ *StreamController) serveHlsTs(w http.ResponseWriter, r *http.Request, mediaBase, segName string) {
	mediaFile := fmt.Sprintf("%s\\%s", mediaBase, segName)
	http.ServeFile(w, r, mediaFile)
	w.Header().Set("Content-Type", "video/MP2T")
}