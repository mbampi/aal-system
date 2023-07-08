
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
- [X] Add Openllet Reasoner to Apache Jena Fuseki Server
- [X] Websocket connection to Home Assistant
- [ ] Load phyisical environment from config file
- [ ] Load patient data from config file
- [ ] Add/Remove/Update SPARQL queries to Fuseki via REST API
- [ ] Run queries only when sensor data changed
- [ ] Insert SWRL rules in Fuseki via REST API
- [ ] Serve REST API to allow client connection
- [ ] Dockerize AAL System
- [ ] Build a docker compose to run the Apache Jena Server and the AAL System

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

### Update

```bash
task update
```

## High-level AAL System Overview

Pre-requisites:
- There is a patient
- The patient lives in a house
- The house has sensors
- The sensors are connected to a smart home system (e.g. home assistant)
- The smart home system has an API

Installation:
  - Gather information about the sensors in the house (may install new sensors if needed)
  - Install the AAL System
  - Configure the AAL System with the environment and sensors information
  - Connect the AAL System to the smart home system API (to get sensor data)
  - Connect the AAL System to the EHR system (to get patient data)
  - Based on the patient data, the AAL System will activate the monitoring rules to be used
  
Usage:
  - The doctor will use the EHR system to update the patient data
  - The AAL System will update automatically 
    - with the new patient data (e.g. new medication)
    - new sensor data (e.g. heart rate)
  - New rules can be added to the AAL System
  - AAL System will
    - provide real-time information about the patient health status (which can be used by the doctor/caregiver to make decisions)
    - provide historical information about the patient health status
    - send alerts based on the rules
    - take actions on the smart home actuators based on the rules
