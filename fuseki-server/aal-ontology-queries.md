
# AAL Ontology Lite Queries

### Query all patients
```sparql
PREFIX rdfs: <http://www.w3.org/2000/01/rdf-schema#>
PREFIX rdf: <http://www.w3.org/1999/02/22-rdf-syntax-ns#>
PREFIX : <http://www.semanticweb.org/matheusdbampi/ontologies/2023/6/aal-ontology-lite/>

SELECT ?patient
WHERE {
  ?patient rdf:type/rdfs:subClassOf* :Patient .
}
```

### Query sensors worn by patients
```sparql
PREFIX rdfs: <http://www.w3.org/2000/01/rdf-schema#>
PREFIX rdf: <http://www.w3.org/1999/02/22-rdf-syntax-ns#>
PREFIX : <http://www.semanticweb.org/matheusdbampi/ontologies/2023/6/aal-ontology-lite/>

SELECT ?sensor
WHERE {
  ?patient rdf:type/rdfs:subClassOf* :Patient .
  ?patient :wears ?sensor .
}
```

### Query sensors worn by an specific patient (Matheus)
```sparql
PREFIX rdfs: <http://www.w3.org/2000/01/rdf-schema#>
PREFIX rdf: <http://www.w3.org/1999/02/22-rdf-syntax-ns#>
PREFIX foaf: <http://xmlns.com/foaf/0.1/>
PREFIX : <http://www.semanticweb.org/matheusdbampi/ontologies/2023/6/aal-ontology-lite/>

SELECT ?sensor
WHERE {
    ?patient rdf:type/rdfs:subClassOf* :Patient .
  	?patient foaf:name "Matheus" .
  	?patient :wears ?sensor .
}
```

### Query all patients with a specific disease
```sparql
PREFIX rdfs: <http://www.w3.org/2000/01/rdf-schema#>
PREFIX rdf: <http://www.w3.org/1999/02/22-rdf-syntax-ns#>
PREFIX : <http://www.semanticweb.org/matheusdbampi/ontologies/2023/6/aal-ontology-lite/>

SELECT ?patient ?d
WHERE {
  ?patient rdf:type/rdfs:subClassOf* :Patient .
  ?patient :hasMedicalCondition ?d .
  ?d rdf:type :Diabetes .
}
```

### Query all patients and their diseases
```sparql
PREFIX rdfs: <http://www.w3.org/2000/01/rdf-schema#>
PREFIX rdf: <http://www.w3.org/1999/02/22-rdf-syntax-ns#>
PREFIX : <http://www.semanticweb.org/matheusdbampi/ontologies/2023/6/aal-ontology-lite/>

SELECT ?patient ?d
WHERE {
  ?patient rdf:type/rdfs:subClassOf* :Patient .
  ?patient :hasMedicalCondition ?d .
}
```

### Query all patients properties and values
```sparql
PREFIX sosa: <http://www.w3.org/ns/sosa/>
PREFIX rdfs: <http://www.w3.org/2000/01/rdf-schema#>
PREFIX rdf: <http://www.w3.org/1999/02/22-rdf-syntax-ns#>
PREFIX : <http://www.semanticweb.org/matheusdbampi/ontologies/2023/6/aal-ontology-lite/>

SELECT DISTINCT ?patient ?sensor ?prop ?value
WHERE {
  ?patient rdf:type/rdfs:subClassOf* :Patient .
  ?patient :wears ?sensor .
  ?sensor sosa:madeObservation ?obs .
  ?obs sosa:observes ?prop .
  ?obs sosa:hasSimpleResult ?value
}
```

### Query to add a new observation
```sparql
PREFIX xsd: <http://www.w3.org/2001/XMLSchema#>
PREFIX sosa: <http://www.w3.org/ns/sosa/>
PREFIX rdf: <http://www.w3.org/1999/02/22-rdf-syntax-ns#>
PREFIX : <http://www.semanticweb.org/matheusdbampi/ontologies/2023/6/aal-ontology-lite/>

INSERT DATA {
  :obs1 rdf:type sosa:Observation .
  :obs1 sosa:observes :HeartRate .
  :obs1 sosa:hasSimpleResult 80 .
  :obs1 sosa:madeBySensor :emfit .
  :obs1 sosa:resultTime "2021-07-01T00:00:00"^^xsd:dateTime .
}
```

### Query to list all observations
```sparql
PREFIX sosa: <http://www.w3.org/ns/sosa/>
PREFIX rdf: <http://www.w3.org/1999/02/22-rdf-syntax-ns#>
PREFIX : <http://www.semanticweb.org/matheusdbampi/ontologies/2023/6/aal-ontology-lite/>

SELECT ?obs ?sensor ?prop ?value
WHERE {
  ?obs rdf:type sosa:Observation .
  ?obs sosa:observes ?prop .
  ?obs sosa:hasSimpleResult ?value .
  ?obs sosa:madeBySensor ?sensor
}
```

### Query to get the last observation of a patient wearable sensor
```sparql
PREFIX xsd: <http://www.w3.org/2001/XMLSchema#>
PREFIX sosa: <http://www.w3.org/ns/sosa/>
PREFIX rdfs: <http://www.w3.org/2000/01/rdf-schema#>
PREFIX rdf: <http://www.w3.org/1999/02/22-rdf-syntax-ns#>
PREFIX : <http://www.semanticweb.org/matheusdbampi/ontologies/2023/6/aal-ontology-lite/>

SELECT ?obs ?prop ?value
WHERE {
  ?patient rdf:type/rdfs:subClassOf* :Patient .
  ?patient :wears ?sensor .
  ?sensor sosa:madeObservation ?obs .
  ?obs sosa:observes ?prop .
  ?obs sosa:hasSimpleResult ?value
}
ORDER BY DESC(?obs)
LIMIT 1
```