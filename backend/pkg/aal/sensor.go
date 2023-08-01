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
	sensorID := "sensor_" + s.ID
	return Query(`PREFIX : <http://www.semanticweb.org/matheusdbampi/ontologies/2023/6/aal-ontology-lite/>
	INSERT DATA {
		:` + sensorID + ` rdf:type :Sensor .
		:` + sensorID + ` :hasName "` + s.Name + `" .
		:` + sensorID + ` :observes :` + string(s.ObservableProperty) + ` .
		:` + sensorID + ` :hasFeatureOfInterest :` + string(s.FeatureOfInterest) + ` .
		:` + sensorID + ` :locatedAt :` + string(s.InstalledAt) + ` .
	}`)
}

func (s *Sensor) RemoveQuery() Query {
	sensorID := "sensor_" + s.ID
	return Query(`PREFIX : <http://www.semanticweb.org/matheusdbampi/ontologies/2023/6/aal-ontology-lite/>
	DELETE {
		:` + sensorID + ` ?p ?o .
	}
	WHERE {
		:` + sensorID + ` ?p ?o .
	}`)
}
