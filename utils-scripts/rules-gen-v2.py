def generate_rules(n):
    """Generate rules based on the given number."""
    rules = []
    for i in range(1, n+1):
        rule = f"""[rule_{i}:
    (?o rdf:type sosa:Observation) 
    (?o sosa:madeBySensor ?s) 
    (?s sosa:observes ?op)
    (?o sosa:hasSimpleResult ?r)
    greaterThan(?r, {i})
    lessThan(?r, {i+2})
    ->
    (:s :inferred ?op)
]"""
        rules.append(rule)
    return "\n\n".join(rules)

def main():
    """Main function to execute the script."""
    n = int(input("Enter the value of n: "))
    rules = generate_rules(n)

    # Write to a file
    with open("rules.txt", "w") as f:
        f.write(rules)

    print(f"Generated {n} rules and written to rules.txt")

if __name__ == "__main__":
    main()
