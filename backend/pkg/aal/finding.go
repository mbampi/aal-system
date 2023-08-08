package aal

import (
	"aalsystem/pkg/fuseki"
	"log"
	"strings"
	"time"
)

type Finding struct {
	Name      string    `json:"name"`
	Patient   string    `json:"patient"`
	Sensor    string    `json:"sensor"`
	Value     string    `json:"value"`
	Timestamp time.Time `json:"timestamp"`
}

func findingsQuery() string {
	query := `
PREFIX sosa: <http://www.w3.org/ns/sosa/>
PREFIX rdf: <http://www.w3.org/1999/02/22-rdf-syntax-ns#>
PREFIX foaf: <http://xmlns.com/foaf/0.1/>
PREFIX skos: <http://www.w3.org/2004/02/skos/core#>
PREFIX : <http://www.semanticweb.org/matheusdbampi/ontologies/2023/7/aal-ontology#>

SELECT ?patientName ?sensor ?findingName ?value ?time
WHERE {
	?patient :hasFinding ?finding .
	?patient :hasName ?patientName .
	?finding :inferredBy ?observation .
	?finding skos:prefLabel ?findingName .
	?observation sosa:hasSimpleResult ?value .
	?observation sosa:madeBySensor ?sensor .
	?observation sosa:resultTime ?time 
}`
	return query
}

func resultToFindings(results fuseki.QueryResult) []Finding {
	bindings := results.Results.Bindings
	if len(bindings) == 0 {
		return nil
	}

	findings := make([]Finding, 0, len(bindings))
	for _, binding := range results.Results.Bindings {
		finding := binding["findingName"].Value
		finding = finding[strings.LastIndex(finding, "/")+1:]

		patient := binding["patientName"].Value
		value := binding["value"].Value

		sensor := binding["sensor"].Value
		sensor = sensor[strings.LastIndex(sensor, "#")+1:]

		layout := "2006-01-02T15:04:05"
		time, err := time.Parse(layout, binding["time"].Value)
		if err != nil {
			log.Println("failed to parse time:", err)
		}

		findings = append(findings, Finding{
			Name:      finding,
			Patient:   patient,
			Value:     value,
			Sensor:    sensor,
			Timestamp: time,
		})
	}

	return findings
}

// IsDuplicate returns true if the finding is already in the list.
func (f *Finding) IsDuplicate(list []Finding) bool {
	for _, finding := range list {
		if f.Name == finding.Name &&
			f.Patient == finding.Patient &&
			f.Sensor == finding.Sensor {
			return true
		}
	}
	return false
}
