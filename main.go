package lastword

import (
	"context"
	"fmt"
	"net/http"
	"os"

	"cloud.google.com/go/datastore"
)

var project = os.Getenv("GCP_PROJECT")
var authToken = os.Getenv("AUTH_TOKEN")
var kind = os.Getenv("KIND")

type lastword struct {
	Status string
}

func writeStatus(ctx context.Context, client *datastore.Client, key, status string) error {
	k := datastore.NameKey(kind, key, nil)
	_, err := client.Put(ctx, k, &lastword{Status: status})
	return err
}
func getStatus(ctx context.Context, client *datastore.Client, key string) (string, error) {
	k := datastore.NameKey(kind, key, nil)
	var lw lastword
	err := client.Get(ctx, k, &lw)
	return lw.Status, err
}

func showError(w http.ResponseWriter, err error) {
	w.WriteHeader(http.StatusInternalServerError)
	fmt.Fprintf(w, err.Error())
}

// LastWord is entrypoint
func LastWord(w http.ResponseWriter, r *http.Request) {
	if r.FormValue("token") != authToken {
		w.WriteHeader(http.StatusForbidden)
		fmt.Fprintf(w, "need token param!")
		return
	}
	ctx := context.Background()
	client, err := datastore.NewClient(ctx, project)
	if err != nil {
		showError(w, err)
		return
	}
	key := r.FormValue("key")
	if key == "" {
		w.WriteHeader(http.StatusBadRequest)
		fmt.Fprintf(w, "need key param!")
		return
	}
	if r.Method == http.MethodPost {
		status := r.FormValue("status")
		err := writeStatus(ctx, client, key, status)
		if err != nil {
			showError(w, err)
			return
		}

		fmt.Fprintf(w, "wrote { %v: %v }", key, status)
		return
	}

	status, err := getStatus(ctx, client, key)
	if err != nil {
		showError(w, err)
		return
	}
	fmt.Fprintf(w, status)
}
