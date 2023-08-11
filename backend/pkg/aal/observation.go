package aal

import (
	"fmt"
	"strings"
	"time"
)

type Observation struct {
	Name      string    `json:"name"`
	Sensor    string    `json:"sensor"`
	Value     string    `json:"value"`
	Timestamp time.Time `json:"timestamp"`
}

func (o *Observation) InsertQuery() Query {
	obsID := fmt.Sprintf("obs_%s", o.Sensor)
	builder := strings.Builder{}

	builder.WriteString(`PREFIX xsd: <http://www.w3.org/2001/XMLSchema#>` + "\n")
	builder.WriteString(`PREFIX sosa: <http://www.w3.org/ns/sosa/>` + "\n")
	builder.WriteString(`PREFIX rdf: <http://www.w3.org/1999/02/22-rdf-syntax-ns#>` + "\n")
	builder.WriteString(`PREFIX : <http://www.semanticweb.org/matheusdbampi/ontologies/2023/7/aal-ontology#>` + "\n")

	builder.WriteString("DELETE {" + "\n")
	builder.WriteString(`	:` + obsID + ` rdf:type sosa:Observation .` + "\n")
	builder.WriteString(`	:` + obsID + ` ?a ?b .` + "\n")
	builder.WriteString("}" + "\n")

	builder.WriteString("INSERT {" + "\n")
	builder.WriteString(`	:` + obsID + ` rdf:type sosa:Observation .` + "\n")
	if o.Value == "on" {
		builder.WriteString(`	:` + obsID + ` sosa:hasSimpleResult "true"^^xsd:boolean .` + "\n")
	} else if o.Value == "off" {
		builder.WriteString(`	:` + obsID + ` sosa:hasSimpleResult "false"^^xsd:boolean .` + "\n")
	} else {
		builder.WriteString(`	:` + obsID + ` sosa:hasSimpleResult ` + o.Value + ` .` + "\n")
	}
	builder.WriteString(`	:` + obsID + ` sosa:madeBySensor :` + o.Sensor + ` .` + "\n")
	builder.WriteString(`	:` + obsID + ` sosa:resultTime "` + o.Timestamp.Format("2006-01-02T15:04:05") + `"^^xsd:dateTime .` + "\n")
	builder.WriteString(`   :` + obsID + ` sosa:hasFeatureOfInterest :patient1.` + "\n")
	builder.WriteString("}")

	builder.WriteString("WHERE {" + "\n")
	builder.WriteString(` OPTIONAL {` + "\n")
	builder.WriteString(`	:` + obsID + ` rdf:type sosa:Observation .` + "\n")
	builder.WriteString(`	:` + obsID + ` ?a ?b .` + "\n")
	builder.WriteString("}" + "\n")
	// builder.WriteString(`	:` + obsID + ` rdf:type sosa:Observation .` + "\n")
	// builder.WriteString(`	:` + obsID + ` sosa:hasSimpleResult ?old_value .` + "\n")
	// builder.WriteString(`	:` + obsID + ` sosa:madeBySensor :` + o.Sensor + ` .` + "\n")
	// builder.WriteString(`	:` + obsID + ` sosa:resultTime ?old_timestamp .` + "\n")
	// builder.WriteString(" }" + "\n")
	builder.WriteString("}")

	return Query(builder.String())
}
