```bash
$ terraform --version
$ az --version
```

```bash
$ az login
```

```bash
$ dacrane apply base base -a '{ prefix: dacrane }'
```

```bash
$ dacrane apply api-image api-v1 \
  -a '{ tag: v1, acr: "${{ instances.base.modules.acr }}" }'
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
