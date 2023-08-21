
initial_number = int(input("Enter the initial number: "))
final_number = int(input("Enter the final number: "))

for i in range(initial_number, final_number + 1):
    sensor_name = f"sensor.sensor_{i}"
    alias = f"sensor_{i}"
    print(f'm.AddSensor("{sensor_name}", "{alias}")')
