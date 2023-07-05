package aalontology

import "fmt"

type Patient struct {
	ID   string
	Name string
	Age  int
}

func (p *Patient) InsertQuery() Query {
	return Query(`PREFIX : <http://www.semanticweb.org/matheusdbampi/ontologies/2023/6/aal-ontology-lite/>
	PREFIX rdf: <http://www.w3.org/1999/02/22-rdf-syntax-ns#>

	INSERT DATA {
		:patient_` + p.ID + ` rdf:type :Patient .
		:patient_` + p.ID + ` :hasName "` + p.Name + `" .
		:patient_` + p.ID + ` :hasAge "` + fmt.Sprint(p.Age) + `" .
	}`)
}

func (p *Patient) RemoveQuery() Query {
	return Query(`PREFIX : <http://www.semanticweb.org/matheusdbampi/ontologies/2023/6/aal-ontology-lite/>
	DELETE {
		:patient_` + p.ID + ` ?p ?o .
	}
	WHERE {
		:patient_` + p.ID + ` ?p ?o .
	}`)
}

func (p *Patient) WearsSensor(sensorID string) Query {
	return Query(`PREFIX : <http://www.semanticweb.org/matheusdbampi/ontologies/2023/6/aal-ontology-lite/>
	INSERT DATA {
		:patient_` + p.ID + ` :wears :sensor_` + sensorID + ` .
	}`)
}

func (p *Patient) HasDisease(disease URI) Query {
	return Query(`PREFIX : <http://www.semanticweb.org/matheusdbampi/ontologies/2023/6/aal-ontology-lite/>
	INSERT DATA {
		:patient_` + p.ID + ` :hasDisease :disease_` + string(disease) + ` .
	}`)
}

func (p *Patient) HasPropensityToAdverseReaction(reaction URI) Query {
	return Query(`PREFIX : <http://www.semanticweb.org/matheusdbampi/ontologies/2023/6/aal-ontology-lite/>
	INSERT DATA {
		:patient_` + p.ID + ` :hasPropensityToAdverseReaction :reaction_` + string(reaction) + ` .
	}`)
}

func (p *Patient) TakesMedication(medication URI) Query {
	return Query(`PREFIX : <http://www.semanticweb.org/matheusdbampi/ontologies/2023/6/aal-ontology-lite/>
	INSERT DATA {
		:patient_` + p.ID + ` :takesMedication :medication_` + string(medication) + ` .
	}`)
}

func (p *Patient) GetClinicalFindings() Query {
	return Query(`SELECT ?finding
	WHERE {
		:patient_` + p.ID + ` :hasClinicalFinding ?finding .
	}`)
}
