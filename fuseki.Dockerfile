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
RUN wget https://downloads.apache.org/jena/binaries/apache-jena-fuseki-4.8.0.tar.gz && \
    tar xzf apache-jena-fuseki-4.8.0.tar.gz && \
    rm apache-jena-fuseki-4.8.0.tar.gz

# Copy Openllet reasoner jar files
COPY ./openllet-jars /fuseki/run/extra

# Copy the configuration file
COPY ./config.ttl /fuseki/config.ttl
COPY ./snomedct.ttl /fuseki/snomedct.ttl

# Expose port for Fuseki
EXPOSE 3030

# Run Fuseki server with provided configuration on startup
CMD ["./apache-jena-fuseki-4.8.0/fuseki-server", "--config=/fuseki/config.ttl"]
