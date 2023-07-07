package aal

type Observation struct {
	ID        string
	SensorID  string
	Value     string
	Unit      string
	Timestamp string
}

func (o *Observation) InsertQuery() Query {
	return Query(`PREFIX : <http://www.semanticweb.org/matheusdbampi/ontologies/2023/6/aal-ontology-lite/>
	INSERT DATA {
		:observation_` + o.ID + ` rdf:type :Observation .
		:observation_` + o.ID + ` :madeBy :` + string(o.SensorID) + ` .
		:observation_` + o.ID + ` :hasValue "` + o.Value + `" .
		:observation_` + o.ID + ` :hasUnit "` + o.Unit + `" .
		:observation_` + o.ID + ` :hasTimestamp "` + o.Timestamp + `" .
	}`)
}
