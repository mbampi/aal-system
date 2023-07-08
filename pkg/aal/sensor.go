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
	FeatureOfInterest  URI // Patient, Room, etc
	ObservableProperty URI // HeartRate, Temperature, etc
	InstalledAt        URI // Patient, Room, Bed, Window, etc

	ValueField string // Field of the Home Assistant state object that contains the value of the observable property
	Unit       string // Unit of the observable property
}

type Actuator struct {
	ID                 string
	Name               string
	FeatureOfInterest  URI // Patient, Room, etc
	ActuatableProperty URI // Temperature, etc
	LocatedAt          URI // Patient, Room, Bed, Window, etc
}

func (s *Sensor) InsertQuery() Query {
	return Query(`PREFIX : <http://www.semanticweb.org/matheusdbampi/ontologies/2023/6/aal-ontology-lite/>
	INSERT DATA {
		:sensor_` + s.ID + ` rdf:type :Sensor .
		:sensor_` + s.ID + ` :hasName "` + s.Name + `" .
		:sensor_` + s.ID + ` :observes :` + string(s.ObservableProperty) + ` .
		:sensor_` + s.ID + ` :hasFeatureOfInterest :` + string(s.FeatureOfInterest) + ` .
		:sensor_` + s.ID + ` :locatedAt :` + string(s.InstalledAt) + ` .
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
