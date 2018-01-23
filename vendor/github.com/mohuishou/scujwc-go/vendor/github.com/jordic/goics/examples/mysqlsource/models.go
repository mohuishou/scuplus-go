package main

import (
	"database/sql"
	"fmt"
	"log"
	"time"

	_ "github.com/go-sql-driver/mysql"
	"github.com/jordic/goics"
)

// We define our struct as mysql row
type Reserva struct {
	Clau       string         `db:"clau"`
	Data       time.Time      `db:"data"`
	Apartament string         `db:"apartament"`
	Nits       int            `db:"nits"`
	Extra      sql.NullString `db:"limpieza"`
	Llegada    sql.NullString `db:"llegada"`
	Client     string         `db:"client"`
}

// A collection of rous
type ReservasCollection []*Reserva

// We implement ICalEmiter interface that will return a goics.Componenter.
func (rc ReservasCollection) EmitICal() goics.Componenter {

	c := goics.NewComponent()
	c.SetType("VCALENDAR")
	c.AddProperty("CALSCAL", "GREGORIAN")
	c.AddProperty("PRODID;X-RICAL-TZSOURCE=TZINFO", "-//tmpo.io")

	for _, ev := range rc {
		s := goics.NewComponent()
		s.SetType("VEVENT")
		dtend := ev.Data.AddDate(0, 0, ev.Nits)
		k, v := goics.FormatDateField("DTEND", dtend)
		s.AddProperty(k, v)
		k, v = goics.FormatDateField("DTSTART", ev.Data)
		s.AddProperty(k, v)
		s.AddProperty("UID", ev.Clau+"@whotells.com")
		des := fmt.Sprintf("%s Llegada: %s", ev.Extra.String, ev.Llegada.String)
		s.AddProperty("DESCRIPTION", des)
		s.AddProperty("SUMMARY", fmt.Sprintf("Reserva de %s", ev.Client))
		s.AddProperty("LOCATION", ev.Apartament)

		c.AddComponent(s)
	}

	return c

}

// Get data from database populating ReservasCollection
func GetReservas() ReservasCollection {

	t := time.Now()
	q := `SELECT clau, data, nits, limpieza, llegada, 
			CONCAT(b.tipus, " ", b.title) as apartament,
			CONCAT(c.nom, " ", c.cognom) as client 
		FROM reserves_reserva as a LEFT join crm_client as c on a.client_id = c.id
		LEFT JOIN reserves_apartament as b on a.apartament_id = b.id
		WHERE data>=? and a.status<=3`

	reservas := ReservasCollection{}
	err := Db.Select(&reservas, q, t.Format("2006-01-02"))
	_ = Db.Unsafe()
	if err != nil {
		log.Println(err)
	}
	return reservas
}
