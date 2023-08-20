def generate_instances(n):
    """Generate instances based on the given number."""
    instances_list = []
    for i in range(1, n + 1):
        instance = f"""###  http://www.semanticweb.org/matheusdbampi/ontologies/2023/7/aal-ontology#sensor_{i}
:sensor_{i} rdf:type owl:NamedIndividual ,
                          <http://www.w3.org/ns/sosa/Sensor> ;
                 <http://www.w3.org/ns/sosa/observes> <http://snomed.info/id/86290005> ."""
        instances_list.append(instance)
    return "\n\n".join(instances_list)

def main():
    """Main function to execute the script."""
    n = int(input("Enter the value of n: "))
    instances = generate_instances(n)

    # Write to a file
    with open("instances.txt", "w") as f:
        f.write(instances)

    print(f"Generated ontology instances for {n} sensors and written to instances.txt")

if __name__ == "__main__":
    main()
