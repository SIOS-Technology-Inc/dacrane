## Prepare

```bash
$ terraform --version
$ az --version
```

```bash
$ az login
```

## Quick Start (Local Docker)

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

## Quick Start (App Service)

```bash
$ dacrane apply quick-start-as qs-as
```

```bash
$ curl https://qs-dacrane-sample-as.azurewebsites.net

hello world
```

```bash
$ dacrane destroy qs-as
```

## Deploy Practically

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
$ curl http://localhost:3000

hello world
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
$ curl https://dev-dacrane-sample-as.azurewebsites.net/

hello world
```

```bash
$ dacrane destroy dev
```

```bash
$ dacrane destroy api-v1
```

```bash
$ dacrane destroy base
```
