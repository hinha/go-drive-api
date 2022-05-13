package main

import (
	"encoding/json"
	"fmt"
	"github.com/gorilla/mux"
	"log"
	"net/http"
	"time"
)

func main() {
	oauth := NewOAuthClient()
	r := mux.NewRouter()
	r.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		fmt.Fprintf(w, "Auth URL: %v", oauth.CodeURL())
	})
	r.HandleFunc("/authorization/google/callback", func(w http.ResponseWriter, r *http.Request) {
		q := r.URL.Query()
		code := q.Get("code")
		token, err := oauth.GetToken(code)
		if err != nil {
			w.WriteHeader(http.StatusBadRequest)
			fmt.Fprintf(w, "%v", err)
		}
		oauth.SaveToken(token)
		http.Redirect(w, r, fmt.Sprintf("%s/drive", r.Host), http.StatusPermanentRedirect)
		return
		// authorization/google/callback?state=state-token&code=4/0AX4XfWi5h5nU8mpllMi-vnXvieQ7vBij3_vDNmPtYSceyXoVy37w0ZpfQXa_Gzr2AadRiw&scope=https://www.googleapis.com/auth/drive.metadata.readonly
	})
	r.HandleFunc("/drive/file-list", func(w http.ResponseWriter, r *http.Request) {

		service, err := oauth.DriveService(r.Context())
		if err != nil {
			writeError(w, err)
			return
		}

		listFile, err := service.FilesList(20)
		if err != nil {
			writeError(w, err)
			return
		}

		response := Response{}
		response.Data.Files = listFile.Files
		files, _ := json.Marshal(response)
		w.WriteHeader(http.StatusOK)
		w.Header().Set("Content-Type", "application/json")
		w.Write(files)
	})

	srv := &http.Server{
		Handler: r,
		Addr:    "127.0.0.1:9099",
		// Good practice: enforce timeouts for servers you create!
		WriteTimeout: 15 * time.Second,
		ReadTimeout:  15 * time.Second,
	}
	log.Fatal(srv.ListenAndServe())
}

func writeError(w http.ResponseWriter, err error) {
	w.WriteHeader(http.StatusOK)
	fmt.Fprintf(w, "%v", err)
}
