
# AAL System
Ontology-Based Health Monitoring System For Ambient Assisted Living

## Description

The goal of this system is to be a semantic connector between Smart Home IoT data and a Health ontology.
Enabling SPARQL queries to be smartly performed with the reasoning engine, providing information for a better AAL management.

![system architecture](sys-architecture.png "System Architecture")

## System capabilities:

- [X] Connects to Home Assistant (Smart Home)
- [X] Request entities state from Home Assistant
- [X] Connects to Apache Jena Fuseki (SPARQL Server)
- [X] Queries SPARQL Server
- [X] Dockerize Apache Jena Fuseki Server
- [ ] Dockerize AAL System
- [ ] CRUD of instances/data in RDF graph database
- [ ] Add/Remove SPARQL queries
- [ ] Run SPARQL queries regularly
- [ ] Run queries only when sensor data changed
- [ ] Build a docker compose to run the Apache Jena Server and the AAL System
- [ ] Serve REST api to allow client connection

## Using

### Pre-requisites
- [Docker](https://docs.docker.com/engine/install/)
- [Task](https://taskfile.dev/#/installation)
- [Go](https://golang.org/doc/install)
- [Home Assistant](https://www.home-assistant.io/docs/installation/)
  - Should be running (usually in a Raspberry Pi) in the same network as this machine 
  - Should have the [RESTful API](https://www.home-assistant.io/integrations/rest/) enabled

### Environment variables

```bash
HASSIO_TOKEN=<your home assistant token>
```

##

### Build

```bash
task build
```

### Run

```bash
task up
```

### Stop

```bash
task down
```