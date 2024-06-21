package main

import (
	"fmt"
	"net/http"
	"sort"
	"strconv"

	"github.com/calamityesp/chirpy/common"
)

type GetStruct struct {
	AuthorId string
	Sort     string
}

func (cfg *apiConfig) handlerChirpsRetrieve(w http.ResponseWriter, r *http.Request) {
	dbChirps, err := cfg.DB.GetChirps()
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "Couldn't retrieve chirps")
		return
	}
	// get the authorId
	getChirps := GetStruct{
		AuthorId: r.URL.Query().Get("author_id"),
		Sort:     r.URL.Query().Get("sort"),
	}
	if getChirps.AuthorId != "" {
		cfg.handlerChirpsRetrieveById(w, getChirps)
		return
	}

	chirps := []common.Chirp{}
	for _, dbChirp := range dbChirps {
		chirps = append(chirps, common.Chirp{
			ID:   dbChirp.ID,
			Body: dbChirp.Body,
		})
	}

	if getChirps.Sort == "asc" {
		sort.Slice(chirps, func(i, j int) bool {
			return chirps[i].ID < chirps[j].ID
		})
	} else {
		sort.Slice(chirps, func(i, j int) bool {
			return chirps[i].ID > chirps[j].ID
		})
	}

	respondWithJSON(w, http.StatusOK, chirps)
}

func (cfg *apiConfig) handlerChirpsRetrieveById(w http.ResponseWriter, cstruct GetStruct) {
	dbChirps, err := cfg.DB.GetChirps()
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "Couldn't retrieve chirps")
		return
	}

	// get the authorId
	authorId, err := strconv.Atoi(cstruct.AuthorId)
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
