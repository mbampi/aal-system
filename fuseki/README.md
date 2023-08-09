# Fuseki

## Description

This is the Fuseki server for the AAL project. It is a SPARQL server that uses Apache Jena Fuseki with Generic Reasoner and SWRL Rules.

## Files

- `dev.Dockerfile`: Dockerfile to build the Fuseki server image
- `config.ttl`: Fuseki configuration file
- `general.rules`: general SWRL rules file for sensors and observations inference
- `medical.rules`: medical SWRL rules file for medical inference, where the medical rules are stored

- `aal-ontology.ttl`: AAL Ontology file (not in git because it has 700MB)
