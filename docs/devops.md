# ⚙️ DevOps Stack

The `gyanesh-help` DevOps stack is a combat-tested toolkit designed for SREs and Platform Engineers. It standardizes infrastructure-as-code and container orchestration pipelines.

## The Environment Store: `devops`

Running `gyanesh-help store install devops` directly provisions:
- **Terraform / OpenTofu**: Infrastructure provisioning.
- **Kubernetes (kubectl & helm)**: Cluster management.
- **AWS/GCP/Azure CLIs**: Bound to your current identity.
- **Docker**: Container runtime tooling.

## Dedicated DevOps Macros

### Kubernetes
- **`k8s-pod-shell`**: `kubectl exec -it <pod> -- bash || kubectl exec -it <pod> -- sh`
  - *Why:* Gracefully falls back to `sh` if `bash` isn't installed in the tiny alpine image.
- **`k8s-logs`**: `kubectl logs -f <pod> --tail=100`
  - *Why:* Trails logs instantly without overwhelming your terminal buffer.
- **`port-forward`**: `kubectl port-forward svc/<service> <local>:<remote>`
  - *Why:* Standardized secure tunnel to internal backend APIs.

### Terraform
- **`tf-plan-clean`**: `terraform init && terraform fmt -recursive && terraform validate && terraform plan`
  - *Why:* Prevents bloated PRs by formatting the HCL before structurally validating the syntax.

### Docker
- **`docker-compose-restart`**: `docker-compose down && docker-compose up -d`
- **`nuke-docker-volume`**: `docker volume rm $(docker volume ls -qf dangling=true)`
  - *Why:* Reclaims gigabytes of disk space safely by only targeting dangling orphanage volumes.

### Cloud Integration
- **`gcloud-project-switch`**: `gcloud config set project $1`
  - *Why:* Switches GC contexts without having to memorize the underlying configuration tree.

---
*Achieve infrastructure zen. Run `gyanesh-help store install devops`.*
