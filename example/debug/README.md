This document describes an example of how print debugging is performed in Dacrane.

## Prepare

```bash
$ dacrane init
```

## Debugging Code

The `debug` plugin has `dummy` and `print` resources.
The `dummy` resource creates a virtual resource that does not exist.
The `print` resource prints the data given as arguments to the console.

For example, the following module call can create a dummy resource

```yaml
- name: dummy
  module: debug/resource/dummy
  arguments:
    a: "dummy-resource"
    b: 1
```

The dummy resource will have the following output

```json
{
  "a": "dummy-resource",
  "b": 1
}
```

For example, the following module call will print arguments

```yaml
- name: print-calculate
  module: debug/resource/print
  arguments:
    expression: 1 + 2 * 3
    result: ${{ 1 + 2 * 3 }}
```

The console output will look like this

```json
{
  "expression": "1 + 2 * 3",
  "result": 7
}
```

Check dacrane.yaml for the complete sample code.

## Apply Debug

The sample code can be executed as follows (Note the console output.)

```bash
$ dacrane apply test test -a arg="'"foo"'"

[test.print-arg (debug/resource/print)] Creating...
{
  "arg": "foo"
}
[test.print-arg (debug/resource/print)] Created.

[test.print-calculate (debug/resource/print)] Creating...
{
  "expression": "1 + 2 * 3",
  "result": 7
}
[test.print-calculate (debug/resource/print)] Created.

[test.print-dummy (debug/resource/print)] Creating...
{
  "resource": {
    "a": "dummy-resource",
    "b": 1
  }
}
[test.print-dummy (debug/resource/print)] Created.
```

### Destroy

The created instance can be deleted as follows

```bash
$ dacrane dacrane destroy test

[test.print-dummy (debug/resource/print)] Deleting...
{
  "resource": {
    "a": "dummy-resource",
    "b": 1
  }
}
[test.print-dummy (debug/resource/print)] Deleted.
[test.print-calculate (debug/resource/print)] Deleting...
{
  "expression": "1 + 2 * 3",
  "result": 7
}
[test.print-calculate (debug/resource/print)] Deleted.
[test.print-arg (debug/resource/print)] Deleting...
{
  "arg": "foo"
}
[test.print-arg (debug/resource/print)] Deleted.
[test.dummy (debug/resource/dummy)] Deleting...
[test.dummy (debug/resource/dummy)] Deleted.
```
