# THIS IS IN NEED OF A RE-WRITE

## Old and not really correct:

### gateway

(currently) GraphQL gateway that covers all the services and serves as BFF for the web.

### pkg

Helper packages that are used across the services.

### web

The web app that serves as the frontend for the vediagames.com website.

### worker

The worker that handles all the background tasks:

- Updates
- Inserts
- Event logs
- Increases

Deployed as serverless app and sits behind an MQ

### search

Search service that provides search functionality for the website.
This service connects all other services that support search and adds business logic on top of them.

### db
Schema and migrations for the database. Also includes some sample data for local development.

### deploy
Manifests and code for deploying the services to the cloud:

- K8s manifests
- Terraform code
- Pulumi code
- Helm charts

### build
Build scripts for the services:

- Dockerfiles
- CI pipelines