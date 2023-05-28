package snomed

import "fmt"

// SuperClassesOf returns a SPARQL query that returns the superclasses of the given class.
func SuperClassesOf(class string) string {
	return fmt.Sprintf(`PREFIX rdfs: <http://www.w3.org/2000/01/rdf-schema#>
		PREFIX skos: <http://www.w3.org/2004/02/skos/core#>
		PREFIX : <http://snomed.info/id/>

		SELECT DISTINCT ?subclass ?label
		WHERE
		{
			:%s rdfs:subClassOf* ?subclass .
			?subclass skos:prefLabel ?label .
		}`, class)
}

// SubclassesOf returns a SPARQL query that returns the subclasses of the given class.
func SubclassesOf(class string) string {
	return fmt.Sprintf(`PREFIX rdfs: <http://www.w3.org/2000/01/rdf-schema#>
		PREFIX skos: <http://www.w3.org/2004/02/skos/core#>
		PREFIX : <http://snomed.info/id/>

		SELECT DISTINCT ?x ?label
		WHERE
		{
			?x rdfs:subClassOf :%s .
			?x skos:prefLabel  ?label .
		}`, class)
}

// ClassLabelAndDescription returns a SPARQL query that returns the label and description of all classes.
func ClassLabelAndDescriptionWithLimit(limit int) string {
	return fmt.Sprintf(`PREFIX owl: <http://www.w3.org/2002/07/owl#>
			PREFIX rdfs: <http://www.w3.org/2000/01/rdf-schema#>
			
			SELECT DISTINCT ?class ?label ?description
			WHERE {
				?class a owl:Class.
				OPTIONAL { ?class rdfs:label ?label}
				OPTIONAL { ?class rdfs:comment ?description}
			}
			LIMIT %d`, limit)
}

// AllSubclasses returns a SPARQL query that returns all subclasses, recursively, of the given class.
func AllSubclassesWithLimit(id string, limit int) string {
	return fmt.Sprintf(`PREFIX owl: <http://www.w3.org/2002/07/owl#>
			PREFIX rdfs: <http://www.w3.org/2000/01/rdf-schema#>
			PREFIX : <http://snomed.info/id/>
			
			SELECT DISTINCT ?class ?label ?description
			WHERE {
				?class a owl:Class.
				?class rdfs:subClassOf+ :%s.
				OPTIONAL { ?class rdfs:label ?label}
				OPTIONAL { ?class rdfs:comment ?description}
			}
			LIMIT %d`, id, limit)
}

// RelatedBodyPart returns a SPARQL query that returns the related body part of the given class.
func RelatedBodyPart(class string) string {
	const AssociatedMorphology = "363698007"
	return fmt.Sprintf(`PREFIX rdfs: <http://www.w3.org/2000/01/rdf-schema#>
		PREFIX skos: <http://www.w3.org/2004/02/skos/core#>
		PREFIX : <http://snomed.info/id/>

		SELECT DISTINCT ?x ?label
		WHERE
		{
			:%s :%s ?x .
			?x skos:prefLabel ?label .
		}`, class, AssociatedMorphology)
}

// RelatedFinding returns a SPARQL query that returns the related finding of the given class.
func RelatedFinding(class string) string {
	const FindingSite = "363698007"
	return fmt.Sprintf(`PREFIX rdfs: <http://www.w3.org/2000/01/rdf-schema#>
		PREFIX skos: <http://www.w3.org/2004/02/skos/core#>
		PREFIX : <http://snomed.info/id/>

		SELECT DISTINCT ?x ?label
		WHERE
		{
			:%s :%s ?x .
			?x skos:prefLabel ?label .
		}`, class, FindingSite)
}
