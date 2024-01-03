This document explains that Dacrane deploys a Node.js API service into Docker and Azure Kubernetes Service.

## Prepare

(If you don’t have a service principal) Create a service principal for Dacrane, for example:

```bash
$ az ad sp create-for-rbac --name "dacrane-${your_name}" --role="Contributor" --scopes="/subscriptions/${your_subscription_id}"

{
  "appId": "00000000-0000-0000-0000-000000000000",
  "displayName": "dacrane-your-name",
  "password": "12345678-0000-0000-0000-000000000000",
  "tenant": "10000000-0000-0000-0000-000000000000"
}
```

Store the credentials as Environment Variables, for example:

```bash
export ARM_CLIENT_ID="00000000-0000-0000-0000-000000000000"
export ARM_CLIENT_SECRET="12345678-0000-0000-0000-000000000000"
export ARM_TENANT_ID="10000000-0000-0000-0000-000000000000"
export ARM_SUBSCRIPTION_ID="20000000-0000-0000-0000-000000000000"
```

## Quick Start (Local Docker)

This section explains the way deploys API service into Docker with one command.

```bash
$ dacrane apply quick-start qs
```

```bash
$ curl http://localhost:3000

hello world
```

```bash
$ dacrane destroy qs
```

## Quick Start (Azure Kubernetes Service)

This section explains the way deploys API service into Azure Kubernetes Service with one command.

```bash
$ dacrane apply quick-start-aks qs-aks
```

```bash
$ curl http://qs-dacrane.japaneast.cloudapp.azure.com/

hello world
```

```bash
$ dacrane destroy qs-aks
```

## Deploy Practically

This section explains the more practical way to deploy.

```bash
$ dacrane apply base base -a '{ prefix: dacrane }'
```

```bash
$ dacrane apply api-image api-v1 \
  -a '{ tag: v1, acr: "${{ instances.base.acr }}" }'
```

```bash
$ dacrane apply local-docker local \
  -a '{ image: "${{ instances.api-v1.api-local-image.image }}" }'
```

```bash
$ curl http://localhost:3000

hello world
```

```bash
$ dacrane apply aks dev -a '
env: dev
base: ${{ instances.base }}
api: ${{ instances.api-v1 }}
'
```

```bash
$ curl http://dev-dacrane.japaneast.cloudapp.azure.com/

hello world
```

```bash
$ dacrane destroy dev
```

```bash
$ dacrane destroy local
```

```bash
$ dacrane destroy api-v1
```

```bash
$ dacrane destroy base
```
