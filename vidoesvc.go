package main

import (
	"net/http"
	"fmt"
	"mime"
	"io"
	"encoding/json"
	"mime/multipart"
	"github.com/gorilla/mux"
	"github.com/urfave/negroni"
	"github.com/mitchellh/goamz/aws"
  "github.com/mitchellh/goamz/s3"
)

type VideoMetadata struct{
	Id string
}

type ContentRangeHeader struct{
	Start int
	End int
	Total int
}

func indexVideo(w http.ResponseWriter, req *http.Request) {

}

func createVideo(w http.ResponseWriter, req *http.Request) {
	contentType := req.Header.Get("Content-Type")

	if contentType != "" {
		fmt.Println(contentType)

		mediaType, params, err := mime.ParseMediaType(contentType)

		if err == nil {
			if mediaType == "application/json" {
				postVideoMetadata(req.Body)
				// TODO: reply to client
			} else if mediaType == "multipart/related" {
				handleMultipartRelated(req, params)
			} else {
				w.Header().Set("Content-Type", "application/json")
				w.WriteHeader(http.StatusBadRequest)
				resp, _ := json.Marshal("{\"code\": 400}")
				w.Write(resp)
			}
		} // TODO: else send 415
	} // TODO: else send 415
}

func handleMultipartRelated(req *http.Request, params map[string]string) (error, error) {
	boundary, ok := params["boundary"]
	if !ok {

	}
	mpReader := multipart.NewReader(req.Body, boundary)
	jsonPart, err := mpReader.NextPart()
	if err != nil {
		return nil, err
	}

	contentType := jsonPart.Header.Get("Content-Type")
	mediaType, params, err := mime.ParseMediaType(contentType)
	charset, ok := params["charset"]
	if (err != nil) || (mediaType != "application/json") || (ok && charset != "utf-8") {
		return nil, err // FIXME: create a new error
	}

	videoMetadata, err := postVideoMetadata(jsonPart)
	if err != nil {
		return nil, err
	}
	videoPart, err := mpReader.NextPart()
	if err != nil {
		return nil, err
	}

	_, err = postVideoData(videoPart, videoMetadata)
	if err != nil {
		return nil, err
	}
	return nil, nil
}

func postVideoMetadata(r io.Reader) (*VideoMetadata, error) {
	fmt.Println("Sending Video Metadata")
	return nil, nil
}

func postVideoData(r io.Reader, v *VideoMetadata) (*VideoMetadata, error){
	fmt.Println("Sending Video Data")
	basicPost(r)
	return nil, nil
}

func basicPost(r io.Reader) {
	auth, err := aws.EnvAuth()
  if err != nil {
    fmt.Println(err)
  }
  client := s3.New(auth, aws.USEast)
	bucket := client.Bucket("mg4-video-staging")
	err = bucket.PutReader("whatever", r, 355856563, "video/mp4", "")

	if err != nil {
		fmt.Println(err)
	}
}

func showVideo(w http.ResponseWriter, req *http.Request) {

}

func updateVideo(w http.ResponseWriter, req *http.Request) {

}

func getVideoStatus(w http.ResponseWriter, req *http.Request) {

}

func main() {
	router := mux.NewRouter()

	router.HandleFunc("/", indexVideo).
		Methods("GET")
	router.HandleFunc("/", createVideo).
		Methods("POST")
	router.HandleFunc("/{videoId}", showVideo).
		Methods("GET")
	router.HandleFunc("/{videoId}", updateVideo).
		Methods("PATCH", "PUT")
	router.HandleFunc("/{videoId}", getVideoStatus).
		Methods("HEAD")


	n := negroni.Classic()
	n.UseHandler(router)
	n.Run(":3000")
}
