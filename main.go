package main

import (
	"database/sql"
	"fmt"
	"net/http"

	"github.com/labstack/echo"
	_ "github.com/lib/pq"
	"github.com/rs/xid"
)

const (
	host     = "localhost"
	port     = 5432
	user     = "postgres"
	password = "password"
	dbname   = "Penjualan"
)

type M map[string]interface{}

type barang struct {
	Id    string `json:"p_id_barang"`
	Nama  string `json:"p_nama_barang"`
	Harga string `json:"p_harga"`
	Stok  string `json:"p_stok"`
}

type barang_array struct {
	barangs []barang `json:"data_barang"`
}

func main() {
	r := echo.New()

	r.GET("/get/barang", func(ctx echo.Context) error {
		psqlconn := fmt.Sprintf("host=%s port= %d user=%s password=%s dbname=%s sslmode=disable", host, port, user, password, dbname)
		db, err := sql.Open("postgres", psqlconn)

		if err != nil {
			data := M{"Pesan": err.Error(), "Status": 500}
			return ctx.JSON(http.StatusInternalServerError, data)
		}
		defer db.Close()

		readsql := "SELECT * from read_barang()"
		rows, error := db.Query(readsql)

		if error != nil {
			panic(error)
		}
		defer rows.Close()

		hasil := barang_array{}
		for rows.Next() {
			data_barang := barang{}
			error := rows.Scan(&data_barang.Id, &data_barang.Nama, &data_barang.Harga, &data_barang.Stok)
			if error != nil {
				panic(error)
			}
			hasil.barangs = append(hasil.barangs, data_barang)

		}

		data := M{"Data": hasil.barangs, "Pesan": "Berhasil dibaca", "Status": 200}
		return ctx.JSON(http.StatusOK, data)

	})

	r.POST("/post/barang", func(ctx echo.Context) error {
		psqlconn := fmt.Sprintf("host=%s port= %d user=%s password=%s dbname=%s sslmode=disable", host, port, user, password, dbname)
		db, err := sql.Open("postgres", psqlconn)

		if err != nil {
			data := M{"Pesan": err.Error(), "Status": 500}
			return ctx.JSON(http.StatusInternalServerError, data)
		}
		defer db.Close()

		id_barang := xid.New().Counter()
		nama_barang := ctx.FormValue("nama_barang")
		harga := ctx.FormValue("harga")
		stok := ctx.FormValue("stok")

		insertsql := "SELECT create_barang($1, $2, $3, $4)"

		_, e := db.Exec(insertsql, id_barang, nama_barang, harga, stok)
		if e != nil {
			data := M{"Pesan": err.Error(), "Status": 500}
			return ctx.JSON(http.StatusInternalServerError, data)
		}

		data := M{"Pesan": "Data berhasil dimasukkan", "Status": 200}
		return ctx.JSON(http.StatusOK, data)

	})

	r.PUT("/put/barang/:id_barang", func(ctx echo.Context) error {
		psqlconn := fmt.Sprintf("host=%s port= %d user=%s password=%s dbname=%s sslmode=disable", host, port, user, password, dbname)
		db, err := sql.Open("postgres", psqlconn)

		if err != nil {
			data := M{"Pesan": err.Error(), "Status": 500}
			return ctx.JSON(http.StatusInternalServerError, data)
		}
		defer db.Close()

		id_barang := ctx.Param("id_barang")
		nama_barang := ctx.FormValue("nama_barang")
		harga := ctx.FormValue("harga")
		stok := ctx.FormValue("stok")

		updatesql := "SELECT update_barang($1, $2, $3, $4)"

		_, e := db.Exec(updatesql, id_barang, nama_barang, harga, stok)
		if e != nil {
			data := M{"Pesan": err.Error(), "Status": 500}
			return ctx.JSON(http.StatusInternalServerError, data)
		}

		data := M{"Pesan": "Data berhasil diperbarui", "Status": 200}
		return ctx.JSON(http.StatusOK, data)
	})

	r.DELETE("/delete/barang/:id_barang", func(ctx echo.Context) error {
		psqlconn := fmt.Sprintf("host=%s port= %d user=%s password=%s dbname=%s sslmode=disable", host, port, user, password, dbname)
		db, err := sql.Open("postgres", psqlconn)

		if err != nil {
			data := M{"Pesan": err.Error(), "Status": 500}
			return ctx.JSON(http.StatusInternalServerError, data)
		}
		defer db.Close()

		id_barang := ctx.Param("id_barang")

		deletesql := "SELECT delete_barang($1)"

		_, e := db.Exec(deletesql, id_barang)
		if e != nil {
			data := M{"Pesan": err.Error(), "Status": 500}
			return ctx.JSON(http.StatusInternalServerError, data)
		}

		data := M{"Pesan": "Data berhasil dihapus", "Status": 200}
		return ctx.JSON(http.StatusOK, data)
	})

	r.Start(":9000")
}
