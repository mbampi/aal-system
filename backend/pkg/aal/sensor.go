package aal

type System struct {
	ID        string
	Name      string
	Sensors   []Sensor
	Actuators []Actuator
}

type Sensor struct {
	ID                 string
	Name               string
	FeatureOfInterest  URI    // Patient, Room, etc
	ObservableProperty string // HeartRate, Temperature, etc
	InstalledAt        URI    // Patient, Room, Bed, Window, etc
	Unit               string // Unit of the observable property
}

type Actuator struct {
	ID                 string
	Name               string
	FeatureOfInterest  URI // Patient, Room, etc
	ActuatableProperty URI // Temperature, etc
	LocatedAt          URI // Patient, Room, Bed, Window, etc
}

// GetSensorQuery returns the query to look for a specific sensor in the ontology
// and return its attributes.
func GetSensorQuery(sensorID string) Query {
	return Query(`PREFIX : <http://www.semanticweb.org/matheusdbampi/ontologies/2023/7/aal-ontology#>
	PREFIX sosa: <http://www.w3.org/ns/sosa/>
	PREFIX dogont: <http://elite.polito.it/ontologies/dogont.owl#>
	PREFIX rdf: <http://www.w3.org/1999/02/22-rdf-syntax-ns#>
	PREFIX skos: <http://www.w3.org/2004/02/skos/core#>

	SELECT ?observableProperty ?installedAt
	WHERE {
		:` + sensorID + ` rdf:type sosa:Sensor .
		OPTIONAL { 
			:` + sensorID + ` sosa:observes ?prop . 
			?prop skos:prefLabel ?observableProperty .
		}
		OPTIONAL { :` + sensorID + ` dogont:isIn ?installedAt . }
	}`)
}
