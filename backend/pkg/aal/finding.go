package aal

import (
	"aalsystem/pkg/fuseki"
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

SELECT ?patientName ?findingName ?value
WHERE {
  ?patient :hasFinding ?finding .
  ?patient :hasName ?patientName .
  ?finding :inferredBy ?observation .
  ?finding skos:prefLabel ?findingName .
  ?obs sosa:hasSimpleResult ?value .
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
		// TODO: add time

		findings = append(findings, Finding{
			Name:    finding,
			Patient: patient,
			Value:   value,
		})
	}

	return findings
}
