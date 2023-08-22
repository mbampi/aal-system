import docker
import time

# Parameters
container_names = ['aal-system-fuseki-1', 'aal-system-backend-1']
sampling_interval = 10  # seconds
sampling_time = 60 * 60 * 4 # 4 hours
total_samples = int(60 / sampling_interval * sampling_time) # 4 hours
samples = 0

client = docker.DockerClient(base_url='unix://var/run/docker.sock')

# get_container_stats 
def get_container_stats(container) -> (float, float):
    stats = container.stats(stream=False)

    mem_usage: float = stats['memory_stats'].get('usage', 0)
    mem_cache: float = stats['memory_stats']['stats'].get('cache', 0)
    memory_usage = mem_usage - mem_cache
    memory_usage = memory_usage / (1024 * 1024 * 1024) # convert to GB

    cpu_total = float(stats["cpu_stats"]["cpu_usage"]["total_usage"])
    system_cpu = float(stats["cpu_stats"]["system_cpu_usage"])
    online_cpus = stats["cpu_stats"]["online_cpus"]
    cpu_delta = cpu_total - float(stats["precpu_stats"]["cpu_usage"]["total_usage"])
    system_delta = system_cpu - float(stats["precpu_stats"]["system_cpu_usage"])
    if system_delta > 0 and cpu_delta > 0:
        cpu_percent = (cpu_delta / system_delta) * online_cpus * 100.0
    else:
        cpu_percent = 0.0
    print(f"{container.name} = Memory: {memory_usage:.2f} GB  |  CPU: {cpu_percent:.2f}%")

    return memory_usage, cpu_percent

def print_average():
    print("\n")
    for idx, name in enumerate(container_names):
        average_mem = total_mems[idx] / samples
        average_cpu = total_cpus[idx] / samples
        print(f"For {name} ({samples}/{total_samples}):")
        print(f"  Average Memory Usage: {(average_mem):.2f} GB")
        print(f"  Average CPU Usage: {average_cpu:.2f}%")

containers = []
for name in container_names:
    try:
        container = client.containers.get(name)
        containers.append(container)
    except docker.errors.NotFound:
        print(f"Container '{name}' not found!")
        exit(1)

total_mems = [0.0 for _ in container_names]
total_cpus = [0.0 for _ in container_names]

for _ in range(total_samples):
    for idx, container in enumerate(containers):
        mem, cpu = get_container_stats(container)
        total_mems[idx] += mem
        total_cpus[idx] += cpu
    samples += 1

    if samples % 6 == 0: # every minute
        print_average()
    
    time.sleep(sampling_interval)

print("\n")
print("------- Final Results -------")
print_average()

client.close()
