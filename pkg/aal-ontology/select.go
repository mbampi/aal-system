package aalontology

// PatientDiseases returns a SPARQL query that returns the diseases of the given patient.
// The patient is identified by its id.
// The query returns the diseases.
func PatientDiseases(patientID string) Query {
	return Query(`PREFIX rdfs: <http://www.w3.org/2000/01/rdf-schema#>
	PREFIX rdf: <http://www.w3.org/1999/02/22-rdf-syntax-ns#>
	PREFIX : <http://www.semanticweb.org/matheusdbampi/ontologies/2023/6/aal-ontology-lite/>
	
	SELECT ?medicalCondition
	WHERE {
	  ?patient rdf:type/rdfs:subClassOf* :Patient .
	  ?patient :hasID "` + patientID + `" .
	  ?patient :hasMedicalCondition ?medicalCondition .
	}`)
}

// ValueByPropertyQuery returns a SPARQL query that returns the value of the given property of the given patient.
// The patient is identified by its id.
// The query returns the value of the property.
func ValueByPropertyQuery(patientID, property string) Query {
	return Query(`PREFIX sosa: <http://www.w3.org/ns/sosa/>
	PREFIX rdfs: <http://www.w3.org/2000/01/rdf-schema#>
	PREFIX rdf: <http://www.w3.org/1999/02/22-rdf-syntax-ns#>
	PREFIX foaf: <http://xmlns.com/foaf/0.1/>
	PREFIX : <http://www.semanticweb.org/matheusdbampi/ontologies/2023/6/aal-ontology-lite/>

	SELECT ?value
	WHERE {
		?patient rdf:type/rdfs:subClassOf* :Patient .
		?patient :hasID "` + patientID + `" .
		?patient :wears ?sensor .
		?sensor sosa:madeObservation ?obs .
		?obs sosa:observes :` + property + ` .
		?obs sosa:hasSimpleResult ?value .
	}`)
}

// PatientObservationsQuery returns a SPARQL query that returns the observations of the given patient.
// The patient is identified by its id.
// The query returns the patient, the property and the value of the observation.
func PatientObservationsQuery(patientID string) Query {
	return Query(`PREFIX sosa: <http://www.w3.org/ns/sosa/>
	PREFIX rdfs: <http://www.w3.org/2000/01/rdf-schema#>
	PREFIX rdf: <http://www.w3.org/1999/02/22-rdf-syntax-ns#>
	PREFIX foaf: <http://xmlns.com/foaf/0.1/>
	PREFIX : <http://www.semanticweb.org/matheusdbampi/ontologies/2023/6/aal-ontology-lite/>
	
	SELECT ?patient ?prop ?value
	WHERE {
		?patient rdf:type/rdfs:subClassOf* :Patient .
		?patient :hasID "` + patientID + `" .
		?patient :wears ?sensor .
		?sensor sosa:madeObservation ?obs .
		?obs sosa:observes ?prop .
		?obs sosa:hasSimpleResult ?value
	}`)
}

// ValuesInPatientRoomQuery returns a SPARQL query that returns
// the values of the given property of the given patient in the given room.
// The patient is identified by its id.
// The query returns the value of the property.
func ValuesInPatientRoomQuery(patientID string) Query {
	return Query(`PREFIX sosa: <http://www.w3.org/ns/sosa/>
	PREFIX rdfs: <http://www.w3.org/2000/01/rdf-schema#>
	PREFIX rdf: <http://www.w3.org/1999/02/22-rdf-syntax-ns#>
	PREFIX foaf: <http://xmlns.com/foaf/0.1/>
	PREFIX : <http://www.semanticweb.org/matheusdbampi/ontologies/2023/6/aal-ontology-lite/>

	SELECT ?room ?sensor ?prop ?value
	WHERE {
		?patient rdf:type/rdfs:subClassOf* :Patient .
		?patient :hasID "` + patientID + `" .
		?patient :isIn ?room .
		?room :hasSensor ?sensor .
		?sensor sosa:madeObservation ?obs .
		?obs sosa:observes ?prop .
		?obs sosa:hasSimpleResult ?value
	}`)
}
