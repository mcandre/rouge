# BUILDTIME REQUIREMENTS

* [Go](https://golang.org/) 1.17+
* [accio](https://github.com/mcandre/accio) v0.0.2

## Recommended

* [snyk](https://www.npmjs.com/package/snyk) 1.893.0 (`npm install -g snyk@1.893.0`)

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