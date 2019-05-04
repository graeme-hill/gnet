To build everything:

```
bazel build //...
```

To add a new go dependency:

```
govendor fetch URL
bazel run //:gazelle
```
