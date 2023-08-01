package aal

import (
	"encoding/json"
	"fmt"
	"io"
	"os"
	"time"
)

// Patient represents a patient in the AAL system.
// It is a wrapper around the patient entity in the ontology.
// It is responsible for creating the patient entity in the ontology.
type Patient struct {
	ID                              string    `json:"id"`
	Name                            string    `json:"name"`
	BirthDate                       time.Time `json:"birth_date"`
	HasDiseases                     []URI     `json:"has_disease"`
	TakesMedications                []URI     `json:"takes_mdication"`
	HasPropensityToAdverseReactions []URI     `json:"has_propensity_to_adverse_reaction"`
}

// NewPatientFromFile creates a new patient from a JSON file.
func NewPatientFromFile(id, filename string) (*Patient, error) {
	var patient *Patient
	file, err := os.Open(filename)
	if err != nil {
		return nil, fmt.Errorf("failed to open patient file: %w", err)
	}
	data, err := io.ReadAll(file)
	if err != nil {
		return nil, fmt.Errorf("failed to read patient file: %w", err)
	}
	err = json.Unmarshal(data, &patient)
	if err != nil {
		return nil, fmt.Errorf("failed to parse patient file: %w", err)
	}
	patient.ID = id
	return patient, nil
}

func (p *Patient) InsertQuery() Query {
	return Query(`PREFIX : <http://www.semanticweb.org/matheusdbampi/ontologies/2023/6/aal-ontology-lite/>
	PREFIX rdf: <http://www.w3.org/1999/02/22-rdf-syntax-ns#>

	INSERT DATA {
		:patient_` + p.ID + ` rdf:type :Patient .
		:patient_` + p.ID + ` :hasName "` + p.Name + `" .
		:patient_` + p.ID + ` :hasBirtDate "` + fmt.Sprint(p.BirthDate) + `" .
		:patient_` + p.ID + ` :hasID "` + p.ID + `" .
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
