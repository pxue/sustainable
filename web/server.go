package web

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"

	"github.com/pxue/sustainable/data"
	"upper.io/db.v3"
)

func enableCors(w http.ResponseWriter) {
	w.Header().Set("Access-Control-Allow-Origin", "*")
}

func handleVisual(w http.ResponseWriter, r *http.Request) {
	enableCors(w)

	type count struct {
		*data.Factory
		Count int64 `db:"count"`
	}

	query := data.DB.Select("f.name", "f.lat", "f.lon", db.Raw("count(ps.*)")).
		From("factories f").
		LeftJoin("product_suppliers ps").On("ps.factory_id = f.id").
		Where(db.Cond{"lat": db.NotEq(0)}).
		GroupBy("f.id")

	var counts []*count
	if err := query.All(&counts); err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		fmt.Fprintf(w, err.Error())
		return
	}

	json.NewEncoder(w).Encode(counts)
}

func New() error {
	log.Println("web launching on 8082")

	http.HandleFunc("/visual", handleVisual)
	return http.ListenAndServe(":8082", nil)
}
