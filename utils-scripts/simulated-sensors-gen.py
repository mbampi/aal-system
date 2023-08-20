import yaml

def generate_sensor_sequence(n):
    sensors = []
    for i in range(1, n+1):
        sensor = {
            "platform": "simulated",
            "name": f"sensor_{i}",
            "unit": "Q",
            "amplitude": 0,
            "mean": 50,
            "spread": 10,
            "seed": i,
            "relative_to_epoch": False
        }
        sensors.append(sensor)
    # Dumps each sensor separately and then joins them with two newlines for the required separation
    return "\n".join([yaml.dump([sensor], default_flow_style=False) for sensor in sensors])

n = int(input("Enter the number of sensors: "))
yaml_output = generate_sensor_sequence(n)

with open('sensors.yaml', 'w') as f:
    f.write(yaml_output)

print("YAML sequence has been saved to 'sensors.yaml'")