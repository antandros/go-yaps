# Yaps Example

we need create interface for plugin (if we want)

```bash
cd example
mkdir plugins
go run cmd/create/create_interface.go
```

now we can create a plugin binary

```bash
mkdir bin
cd plugin_cmd
go build -o ../bin/test test.go
```

now start a example

```bash
go run cmd/run.go
```
