package database

import (
	"encoding/json"
	"errors"
	"log"
	"os"
	"sync"
)

type Chirp struct {
    ID int      `json:"id"`
    Body string `json:"body"`
}

type DB struct {
    path string
    mux *sync.RWMutex
}

type DBStructure struct {
    Chirps map[int]Chirp `json:"chirps"`
}

// NewDB creates a new database connection
// and creates the database file if it doesn't exist
func NewDB(path string) (*DB, error) {
    log.Println("Creating new database connection.")
    db := &DB{ path: path, mux: &sync.RWMutex{} }

    err := db.ensureDB()
    if err != nil {
        log.Printf("Something went horribly wrong while ensuring database.\nError: %v\n", err)
        return nil, err
    }

    return db, nil
}

// CreateChirp creates a new chirp and saves it to disk
func (db *DB) CreateChirp(body string) (Chirp, error) {
    dbStructure, err := db.loadDB()
    if err != nil {
        return Chirp{}, err
    }

    id := len(dbStructure.Chirps) + 1

    chirp := Chirp{ ID: id, Body: body }
    dbStructure.Chirps[id] = chirp

    err = db.writeDB(dbStructure)
    if err != nil {
        return Chirp{}, err
    }

    return chirp, nil
}

// GetChirps returns all chirps in the database
func (db *DB) GetChirps() ([]Chirp, error) {
    dbStructure, err := db.loadDB()
    chirps := make([]Chirp, 0, len(dbStructure.Chirps))
    if err != nil {
        return nil, err
    }

    for _, chirp := range dbStructure.Chirps {
        chirps = append(chirps, chirp)
    }
    
    return chirps, nil
}

// ensureDB creates a new database file if it doesn't exist
func (db *DB) ensureDB() error {
    log.Println("Ensuring database...")

    _, err := os.ReadFile(db.path)
    if errors.Is(err, os.ErrNotExist) {
        log.Println("Database does not exist, creating database...")
        err := db.createDB()
        if err != nil {
            return err
        }
        log.Println("Database created.")
    } else if err != nil {
        log.Println("Error while ensuring database.")
        log.Printf("Error: %v\n", err)
        return err
    } else {
        log.Println("No errors while ensuring database.")
    }

    log.Println("Database ensured.")

    return nil
}

// loadDB reads the database file into memory
func (db *DB) loadDB() (DBStructure, error) {
    db.mux.RLock()
    defer db.mux.RUnlock()

    data, err := os.ReadFile(db.path)
    if err != nil {
        return DBStructure{}, err
    }

    dbStructure := DBStructure{}
    err = json.Unmarshal(data, &dbStructure)
    if err != nil {
        log.Println("An error occurred while unmarshalling the json from database.")
        return DBStructure{}, err
    }

    return dbStructure, nil
}

// writeDB writes the database file to disk
func (db *DB) writeDB(dbStructure DBStructure) error {
    db.mux.Lock()
    defer db.mux.Unlock()

    data, err := json.Marshal(dbStructure)
    if err != nil {
        log.Println("An error occurred while marshalling json for write.")
        return err
    }

    err = os.WriteFile(db.path, data, 0600)
    if err != nil {
        log.Println("An error occurred while writing the database file.")
        return err
    }

    return nil
}

func (db *DB) createDB() error {
    dbStructure := DBStructure{ Chirps: map[int]Chirp{} }
    return db.writeDB(dbStructure)
}
