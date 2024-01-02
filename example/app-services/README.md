This document explains that Dacrane deploys a Node.js API service into Docker and App Service.

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

Set a password for any MySQL, for example:

```bash
export MYSQL_PASSWORD="your_password"
```

## Quick Start (Local Docker)

This section explains the way deploys API service into Docker with one command.

```bash
$ dacrane apply quick-start qs
```

```bash
$ curl http://localhost:3000/status

{"db":{"reachable":true}}
```

```bash
$ curl http://localhost:3000/users

[{"id":1,"name":"alice"},{"id":2,"name":"bob"}]
```

```bash
$ dacrane destroy qs
```

## Quick Start (App Service)

This section explains the way deploys API service into App Service with one command.

```bash
$ dacrane apply quick-start-as qs-as
```

```bash
$ curl https://qs-dacrane-sample-as.azurewebsites.net/status

{"db":{"reachable":true}}
```

```bash
$ curl https://qs-dacrane-sample-as.azurewebsites.net/users

[{"id":1,"name":"alice"},{"id":2,"name":"bob"}]
```

```bash
$ dacrane destroy qs-as
```

## Deploy Practically

This section explains the more practical way to deploy.

```bash
$ dacrane apply base base -a '{ prefix: dacrane }'
```

```bash
$ dacrane apply api-image api-v1 \
  -a '{ tag: v1, acr: "${{ instances.base.modules.acr }}" }'
```

```bash
$ dacrane apply local-docker local \
  -a '{ image: "${{ instances.api-v1.modules.api-local-image.modules.image }}" }'
```

```bash
$ dacrane apply db-migration schema-v1-local -a '
version: v1
network: ${{ instances.local.modules.net.name }}
mysql:
  username: root
  password: my-secret-pw
  host: ${{ instances.local.modules.db.name }}
  database: api
'
```

```bash
$ curl http://localhost:3000/status

{"db":{"reachable":true}}
```

```bash
$ curl http://localhost:3000/users

[{"id":1,"name":"alice"},{"id":2,"name":"bob"}]
```

```bash
$ dacrane apply app-service dev -a '
env: dev
spec: low
base: ${{ instances.base }}
api: ${{ instances.api-v1 }}
'
```

```bash
$ dacrane apply db-migration schema-v1-dev -a '
version: v1
mysql:
  username: ${{ instances.dev.modules.mysql.administrator_login }}@${{ instances.dev.modules.mysql.name }}
  password: ${{ instances.dev.modules.mysql.administrator_login_password }}
  host: ${{ instances.dev.modules.mysql.fqdn }}
  database: ${{ instances.dev.modules.mysql-database.parameter.database }}
'
```

```bash
$ curl https://dev-dacrane-sample-as.azurewebsites.net/status

{"db":{"reachable":true}}
```

```bash
$ curl https://dev-dacrane-sample-as.azurewebsites.net/users

[{"id":1,"name":"alice"},{"id":2,"name":"bob"}]
```

```bash
$ dacrane destroy schema-v1-dev
```

```bash
$ dacrane destroy dev
```

```bash
$ dacrane destroy schema-v1-local
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
