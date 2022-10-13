# BUILDTIME REQUIREMENTS

* [Go](https://golang.org/) 1.19+
* [accio](https://github.com/mcandre/accio) v0.0.3

## Recommended

* [ASDF](https://asdf-vm.com/) 0.10
* [snyk](https://www.npmjs.com/package/snyk) 1.996.0 (`npm install -g snyk@1.996.0`)

# SECURITY AUDIT

```console
$ snyk test
```

# INSTALL

```console
$ go install ./...
```

# UNINSTALL

```console
$ rm "$GOPATH/bin/rouge"
```
