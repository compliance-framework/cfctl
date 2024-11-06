# cfctl

A command line interface to Compliance Framework


## Installations

### Mac
```bash
brew install coreutils go jq
```

```bash
go build -o cfctl
sudo mv cfctl /usr/local/bin/cfctl
mkdir ~/.cfctl
touch ~/.cfctl/config
```

Add to ~/.cfctl/config:
```nano
default: dev
contexts:
  dev:
    url: "http://localhost:8080/api"
```


