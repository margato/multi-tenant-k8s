# Multi-tenant Application

   - [Containers](#containers)
   - [Routing and Authorization](#routing-and-authorization)
   - [Run locally](#run-locally)
     - [Requirements](#requirements)
     - [1. Build application image](#1-build-application-image)
     - [2. Create Kubernetes namespaces](#2-create-kubernetes-namespaces)
     - [3. Install application Helm chart](#3-install-application-helm-chart)
     - [4. Start tunnel](#4-start-tunnel)
     - [5. Make a request](#5-make-a-request)



In a multi-tenant architecture, depending on business maturity, tenants isolation and budget, different approaches can be taken:
- **Pool:** Share-everything — Tenants share resources, but are logically isolated, for example, at the database schema level.
- **Silo:** Share-nothing — Each tenant have dedicated resources, providing resource-level isolation, avoiding the [Noisy Neighbors problem](https://en.wikipedia.org/wiki/Cloud_computing_issues#Performance_interference_and_noisy_neighbors).
- **Bridge:** Hybrid — Share common services, while critical workloads are isolated.

---

In this proof of concept, the application follows a bridge architecture with a single entry point, where requests are routed to each tenant's application plane, which are isolated from one another.

![Doc](./docs/Architecture.jpg)

## Containers
Each tenant has its own namespace, ensuring logical isolation within a single Kubernetes cluster. In this proof of concept (POC), I do not dive deep into node placement, however, in a production environment, tenants could be isolated into different node groups using EKS, for example, based on their criticality.

There is a single source-of-truth to the [application codebase](./apps/products-api/), which is managed as an artifact by [Helm using charts](./infra/helm). This allows the application to be deployed across multiple application planes (namespaces) while only changing values such as memory and CPU limits and HPA configuration.

## Routing and Authorization

The routing layer is powered by Envoy Proxy as a reverse proxy, and Open Policy Agent (OPA), providing authorization to the requested application plane. Communication between Envoy and OPA is done by GRPC and locally inside the Pod, as a sidecar, providing low latency.

Each request is routed based on a HTTP Header named as `x-tenant-id`. To provide security, OPA avaliates the JWT provided to check for a match between the header value and token claims, as well as to validate tenant existence. Note that the Rego implementation does not validate JWT signature nor token expiration, because this is only a proof of concept, in a production environment, this validation is essential.

## Run locally

### Requirements
- Docker
- Minikube
- Helm

### 1. Build application image
Inside `apps/products-api` folder, execute:
```bash
make build
```

### 2. Create Kubernetes namespaces
Inside `apps/products-api` folder, with minikube running, execute:
```bash
kubectl create ns ecommerce-shared
kubectl create ns ecommerce-tenant-1
kubectl create ns ecommerce-tenant-2
kubectl create ns ecommerce-tenant-3
```

### 3. Install application Helm chart
Inside `infra/helm/products-api` folder, with minikube running, execute:

```bash
helm install products-api-tenant-1 . -f values/tenant-1.yaml
helm install products-api-tenant-2 . -f values/tenant-2.yaml
helm install products-api-tenant-3 . -f values/tenant-3.yaml
```

Inside `infra/helm/request-router` folder, with minikube running, execute:
```bash
helm install request-router . -f values/default.yaml
```

### 4. Start tunnel
With minikube running, execute:

```bash
minikube tunnel
```

### 5. Make a request
```bash
curl --request GET \
  --url http://localhost:8080/debug \
  --header 'Authorization: Bearer eyJ0eXAiOiJKV1QiLCJhbGciOiJIUzI1NiJ9.eyJpc3MiOiJPbmxpbmUgSldUIEJ1aWxkZXIiLCJpYXQiOjE3Mjc3NTA1NTcsImV4cCI6MTcyNzc1MDc5OSwiYXVkIjoid3d3LmV4YW1wbGUuY29tIiwic3ViIjoib3N2YWxkb0BlbWFpbC5jb20iLCJlbWFpbCI6Im9zdmFsZG9AZXhhbXBsZS5jb20iLCJ0ZW5hbnQiOiJ0ZW5hbnQtMSJ9.-lBOl5jpRAnQ8YOVCYp2DygldDbWT8avMzL6m4sSkzI' \
  --header 'x-tenant-id: tenant-1'
```

Response 200 OK:
```json
{
  "host_name": "products-api-7df5fb5b94-xh8wx",
  "shop_name":"Tenant 1 Shop"
}
```

This endpoint is from the **products-api**, and shows the hostname and shop name; you should see the name of the tenant in the shop name attribute.

This JWT token was generated using [JSON Web Token Builder](http://jwtbuilder.jamiekurtz.com/) and includes a custom claim called `tenant`.
This particular token is for tenant `tenant-1`, if you want to explore further, you can generate another token for tenant `tenant-2` and tenant `tenant-3`.
