# Super Todo

Try it out now at https://super-todo.app!

This project is a dramatically overengineered todo app I created to practice React, Go, gRPC, microservice design, containerization, and Kubernetes deployments.

## Deployment

This project was tested locally using [k3s](https://k3s.io/) but is currently deployed on GCP Cloud Run due to high cost of GCP Compute and GKE.

## Component Overview

**Client:** Bulit in [Next.js](https://nextjs.org/), using [TailwindCSS](https://tailwindcss.com/) with [shadcn](https://ui.shadcn.com/) UI components. Consumes EventSources from the Gateway API to reduce unnecessary fetching and client side re-rendering.

**Gateway API:** REST API written in Go using [gin](https://github.com/gin-gonic/gin) that acts as a gRPC client to the other microservices and serves HTTP requests to the client.

**User Service:** [gRPC](https://grpc.io/) server written in Go responsible for handling anything related to "User" objects.

**Todo Service:** [gRPC](https://grpc.io/) server written in Go responsible for handling anything related to "Todo" objects.

**Combine Service:** [gRPC](https://grpc.io/) server written in Go responsible for handling anything related to "Combine" objects, which keep track of which users and todos are related.

**PostgreSQL:** Locally run using a container during development, currently hosted on [Supabase](https://supabase.com/).

**Redis:** A docker image of Redis deployed as a container. Used by the Go microservices to cache database requests and as a message broker to horizontally scale Gateway API replicas.

**nginx:** Serves as an ingress container on GCP, with the other containers running as sidecars inside a single Cloud Run instance.
