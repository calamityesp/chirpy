package main

import (
	"fmt"
	"net/http"
	"sort"
	"strconv"

	"github.com/calamityesp/chirpy/common"
)

func (cfg *apiConfig) handlerChirpsRetrieve(w http.ResponseWriter, r *http.Request) {
	dbChirps, err := cfg.DB.GetChirps()
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "Couldn't retrieve chirps")
		return
	}
	// get the authorId
	param := r.URL.Query().Get("author_id")
	if param != "" {
		cfg.handlerChirpsRetrieveById(w, r, param)
		return
	}

	chirps := []common.Chirp{}
	for _, dbChirp := range dbChirps {
		chirps = append(chirps, common.Chirp{
			ID:   dbChirp.ID,
			Body: dbChirp.Body,
		})
	}

	sort.Slice(chirps, func(i, j int) bool {
		return chirps[i].ID < chirps[j].ID
	})

	respondWithJSON(w, http.StatusOK, chirps)
}

func (cfg *apiConfig) handlerChirpsRetrieveById(w http.ResponseWriter, r *http.Request, author_id string) {
	dbChirps, err := cfg.DB.GetChirps()
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "Couldn't retrieve chirps")
		return
	}

	// get the authorId
	authorId, err := strconv.Atoi(author_id)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, err.Error())
	}

	chirps := []common.Chirp{}
	for _, dbChirp := range dbChirps {
		if dbChirp.Author_Id == authorId {
			chirps = append(chirps, common.Chirp{
				ID:        dbChirp.ID,
				Body:      dbChirp.Body,
				Author_Id: dbChirp.Author_Id,
			})
		}

	}

	sort.Slice(chirps, func(i, j int) bool {
		return chirps[i].ID < chirps[j].ID
	})

	respondWithJSON(w, http.StatusOK, chirps)
}

func (cfg *apiConfig) handlerChirpRetrieveById(w http.ResponseWriter, r *http.Request) {
	emptyChirp := common.Chirp{}
	param := r.PathValue("chirpId")
	fmt.Println(param)
	chirpId, err := strconv.Atoi(param)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "Failed to Convert Chirp Id")
	}

	fChirp, err := cfg.DB.GetChirpById(chirpId)
	if err != nil {
		respondWithJSON(w, http.StatusNotFound, fChirp)
		return
	} else if fChirp == emptyChirp {
		respondWithJSON(w, http.StatusNotFound, fChirp)
		return
	}

	respondWithJSON(w, http.StatusOK, fChirp)
}
