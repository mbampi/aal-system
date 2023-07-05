package aalontology

type Sensor struct {
	ID   string
	Name string
}

func (s *Sensor) InsertQuery() Query {
	return Query(`PREFIX : <http://www.semanticweb.org/matheusdbampi/ontologies/2023/6/aal-ontology-lite/>
	INSERT DATA {
		:sensor_` + s.ID + ` rdf:type :Sensor .
		:sensor_` + s.ID + ` :hasName "` + s.Name + `" .
	}`)
}

func (s *Sensor) RemoveQuery() Query {
	return Query(`PREFIX : <http://www.semanticweb.org/matheusdbampi/ontologies/2023/6/aal-ontology-lite/>
	DELETE {
		:sensor_` + s.ID + ` ?p ?o .
	}
	WHERE {
		:sensor_` + s.ID + ` ?p ?o .
	}`)
}

func (s *Sensor) Observes(property URI) Query {
	return Query(`PREFIX : <http://www.semanticweb.org/matheusdbampi/ontologies/2023/6/aal-ontology-lite/>
	INSERT DATA {
		:sensor_` + s.ID + ` :observes :` + string(property) + ` .
	}`)
}

func (s *Sensor) HasFeatureOfInterest(foi URI) Query {
	return Query(`PREFIX : <http://www.semanticweb.org/matheusdbampi/ontologies/2023/6/aal-ontology-lite/>
	INSERT DATA {
		:sensor_` + s.ID + ` :hasFeatureOfInterest :` + string(foi) + ` .
	}`)
}
