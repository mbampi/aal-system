package fuseki

import (
	"bytes"
	"fmt"
)

// Atom represents a SWRL atom.
type Atom struct {
	Type  string
	Class string
	Args  []string
}

// SWRLRule represents a SWRL rule.
type SWRLRule struct {
	Body []Atom
	Head []Atom
}

// ConvertSWRLToTTL converts a SWRL rule to a Turtle (TTL) string.
// The rule is represented as a SWRL Imp (Implication) in Turtle.
func ConvertSWRLToTTL(ruleName string, rule SWRLRule) string {
	var buffer bytes.Buffer

	prefixes := `
@prefix swrl: <http://www.w3.org/2003/11/swrl#> .
@prefix xsd: <http://www.w3.org/2001/XMLSchema#> .
@prefix rdf: <http://www.w3.org/1999/02/22-rdf-syntax-ns#> .
@prefix rdfs: <http://www.w3.org/2000/01/rdf-schema#> .
@prefix owl: <http://www.w3.org/2002/07/owl#> .
@prefix sosa: <http://www.w3.org/ns/sosa/> .
@prefix ssn: <http://www.w3.org/ns/ssn/> .
@prefix : <http://www.semanticweb.org/matheusdbampi/ontologies/2023/6/aal-ontology-lite/> .
`
	buffer.WriteString(prefixes)

	buffer.WriteString(":" + ruleName + " a swrl:Imp ;\n")
	buffer.WriteString("    swrl:body [\n        a swrl:AtomList ;\n")

	// Process body
	for _, atom := range rule.Body {
		buffer.WriteString(createAtomRDF(atom))
	}
	for range rule.Body {
		buffer.WriteString("        ] \n")
	}

	buffer.WriteString("    ] ;\n")
	buffer.WriteString("    swrl:head [\n        a swrl:AtomList ;\n")

	// Process head
	for _, atom := range rule.Head {
		buffer.WriteString(createAtomRDF(atom))
	}
	for range rule.Head {
		buffer.WriteString("        ] \n")
	}

	buffer.WriteString("    ] .\n")

	return buffer.String()
}

// createAtomRDF creates the RDF for a SWRL atom.
func createAtomRDF(atom Atom) string {
	var rdf string

	switch atom.Type {
	case "class":
		rdf = fmt.Sprintf(`        rdf:first [
            a swrl:ClassAtom ;
            swrl:classPredicate :%s ;
            swrl:argument1 %s
        ];`, atom.Class, atom.Args[0])
	case "property":
		rdf = fmt.Sprintf(`        rdf:first [
            a swrl:IndividualPropertyAtom ;
            swrl:propertyPredicate :%s ;
            swrl:argument1 %s ;
            swrl:argument2 %s
        ];`, atom.Class, atom.Args[0], atom.Args[1])
	}

	rdf += "\n        rdf:rest [\n            a swrl:AtomList ;\n"

	return rdf
}

// Example:
// sosa:Observation(?o) ^ sosa:observes(?o, HeartRate) ^ sosa:hasSimpleResult(?o, ?r) ^ swrlb:greaterThan(?r, 85) -> Trigger(?heartRateAlarm) ^ event(?heartRateAlarm, "heart_rate_alarm")
// heartRateRule := fuseki.SWRLRule{
// 	Body: []fuseki.Atom{
// 		{Type: "class", Class: "sosa:Observation", Args: []string{":o"}},
// 		{Type: "property", Class: "sosa:observes", Args: []string{":o", ":HeartRate"}},
// 	},
// 	Head: []fuseki.Atom{
// 		{Type: "class", Class: "Trigger", Args: []string{":heartRateAlarm"}},
// 	},
// }

// ttl := fuseki.ConvertSWRLToTTL("heartRateRule", heartRateRule)
// fmt.Println(ttl)
