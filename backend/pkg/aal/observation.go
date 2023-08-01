package aal

import (
	"fmt"
	"strings"
	"time"
)

type Observation struct {
	ID        string
	SensorID  string
	Value     string
	Timestamp time.Time
}

func (o *Observation) InsertQuery() Query {
	obsID := fmt.Sprintf("obs_%s", o.SensorID)
	builder := strings.Builder{}

	builder.WriteString(`PREFIX xsd: <http://www.w3.org/2001/XMLSchema#>` + "\n")
	builder.WriteString(`PREFIX sosa: <http://www.w3.org/ns/sosa/>` + "\n")
	builder.WriteString(`PREFIX rdf: <http://www.w3.org/1999/02/22-rdf-syntax-ns#>` + "\n")
	builder.WriteString(`PREFIX : <http://www.semanticweb.org/matheusdbampi/ontologies/2023/6/aal-ontology-lite/>` + "\n")

	builder.WriteString("INSERT DATA {" + "\n")
	builder.WriteString(`	:` + obsID + ` rdf:type sosa:Observation .` + "\n")
	builder.WriteString(`	:` + obsID + ` sosa:hasSimpleResult ` + o.Value + ` .` + "\n")
	builder.WriteString(`	:` + obsID + ` sosa:madeBySensor :` + o.SensorID + ` .` + "\n")
	builder.WriteString(`	:` + obsID + ` sosa:resultTime "` + o.Timestamp.Format("2006-01-02T15:04:05") + `"^^xsd:dateTime .` + "\n")
	builder.WriteString("}")

	return Query(builder.String())
}
