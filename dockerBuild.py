import subprocess
from decimal import Decimal, getcontext

getcontext().prec = 1
# 读取当前版本号
with open("version.txt", "r") as f:
    version = f.read().strip()

# 将版本号加1
new_version = Decimal(version) + Decimal(0.1)

# 将新版本号写入文件
with open("version.txt", "w") as f:
    f.write(str(new_version))
# 执行 Docker build 命令
imageId = subprocess.run(["docker", "image", "ls", f"--filter=reference=email:latest", "--format", "{{.ID}}"],
                         capture_output=True).stdout.decode('utf-8').strip("\n")
containerId = subprocess.run(["docker", "ps", "--filter", "ancestor=email", "--format", "{{.ID}}"],
                             capture_output=True).stdout.decode('utf-8').strip("\n")
print(imageId)
print(containerId)

subprocess.run(["docker", "build", "--no-cache", "-t", f"email:v{new_version}", "."])
subprocess.run(["docker", "rmi", f"email:latest"])
subprocess.run(["docker", "tag", f"email:v{new_version}", "email:latest"])
if containerId != "":
    subprocess.run(["docker", "stop", f"{containerId}"])
    subprocess.run(["docker", "rm", f"{containerId}"])
if imageId != "":
    subprocess.run(["docker", "image", "rm", f"{imageId}"])
subprocess.run(["docker", "run", "-d", "--name", "email", "-p","8011:8011","--restart", "always", "email:latest"])
