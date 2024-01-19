# PS Live デモ

## Dacraneのプロジェクトを初期化する

```
$ dacrane init
```

## Azureへの接続情報などの環境変数を読み込む

```
$ source .env
```

## demo環境を構築する

12分くらい待つ。

```bash
$ dacrane apply demo demo \
  -a prefix=demo
```

待っている間にコードの説明

* index.js
* schemas
* dacrane.yaml
* .dacrane
* コマンドの説明

## 動作確認

2分くらい待つ（App ServiceのコンテナのPull&Runに時間がかかってる）

```bash
$ curl https://demo-sample-as.azurewebsites.net/status

{"db":{"reachable":true}}
```

初回のリクエスに時間がかかるのでその間にAzureポータルでできたリソースを確認します。

```bash
$ curl https://demo-sample-as.azurewebsites.net/users

[{"id":1,"name":"alice"},{"id":2,"name":"bob"}]
```

デモは以上。

以下は参考。

次のようにモジュールになっていればなんでもapplyできる。

```
$ dacrane apply resource-group demo-rg \
  -a prefix=demo
```

```
$ dacrane apply container-registry demo-acr \
  -a prefix=demo \
  -a 'resource_group=${{ demo-rg.rg }}'
```

```
$ dacrane apply container-image demo-image-v1 \
  -a tag=v1 \
  -a 'acr=${{ demo-acr.acr }}'
```

Azure Database For MySQLを作る。

```
$ dacrane apply azure-mysql demo-db \
  -a prefix=demo \
  -a 'resource_group=${{ demo-rg.rg }}'
```

```
$ dacrane apply mysql-database demo-mysql-db \
  -a database=demo \
  -a 'mysql=${{ demo-db.mysql }}'
```

```
$ dacrane apply mysql-migration demo-mysql-mig-v1 \
  -a database=demo \
  -a 'mysql=${{ demo-db.mysql }}' \
  -a up_sql=schemas/v1-up.sql \
  -a down_sql=schemas/v1-down.sql
```

```
$ dacrane apply app-service demo-as \
  -a prefix=demo \
  -a 'resource_group=${{ demo-rg.rg }}' \
  -a 'acr=${{ demo-acr.acr }}' \
  -a 'mysql=${{ demo-db.mysql }}' \
  -a 'api=${{ demo-image-v1.remote }}' \
  -a 'database=demo'
```
