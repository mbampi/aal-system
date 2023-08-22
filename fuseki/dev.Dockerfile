# Using Ubuntu as a base image
FROM ubuntu:latest

# Set the Working Directory inside the container
WORKDIR /fuseki

# Install necessary utilities
RUN apt-get update && \
    apt-get install -y wget tar openjdk-11-jre-headless && \
    apt-get clean && \
    rm -rf /var/lib/apt/lists/*

# Download and install Fuseki
RUN wget https://downloads.apache.org/jena/binaries/apache-jena-fuseki-4.9.0.tar.gz && \
    tar xzf apache-jena-fuseki-4.9.0.tar.gz && \
    rm apache-jena-fuseki-4.9.0.tar.gz

# Copy the ontology file
COPY /aal-ontology.ttl /fuseki/aal-ontology.ttl

# Copy the configuration file
COPY /config.ttl /fuseki/config.ttl

# Copy the swrl rules files
COPY /general.rules /fuseki/general.rules
COPY /medical.rules /fuseki/medical.rules

# Expose port for Fuseki
EXPOSE 3030

# Optimize Fuseki Memory
ENV JVM_ARGS="-Xmx2g -Xms512m"

# Run Fuseki server with provided configuration on startup
CMD ["./apache-jena-fuseki-4.9.0/fuseki-server", "--config=/fuseki/config.ttl"]
