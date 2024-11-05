# gityup-go

Given a folder that contains git projects:

```text
~/code
├── project1
├── project2
├── etc
```

loop through them and run `git smart-pull` [^1] and `git remote update origin --prune`.

## Reason

A great way to start learning programming languages is to develop small programs, especially solving problems for yourself.

I decided to write a tool to sync projects I work on and ones I have cloned down just for reference.

Keeping them up-to-date helps me not branch from an old commit ref as well as keeping reference projects current.

## Development

```bash
make
make build
```

### Linting

The linting functionality comes from installing these tools:

- goimports

```bash
go install golang.org/x/tools/cmd/goimports@latest
```

- staticcheck

```bash
go install honnef.co/go/tools/cmd/staticcheck@latest
```

- golangci-lint

<https://golangci-lint.run/welcome/install/>

Then you can run the lint make target:

```bash
make lint
```

### Build for other platforms

```bash
export GOOS=linux
export GOARCH=amd64
make build
```

## Security

There's a nice blog post [here](https://jarosz.dev/article/writing-secure-go-code/) which I've tried to implement here.

Many of the recommended tools are included in the golangci-lint call in the Linting step.

### Vulnerabilities check

The makefile has targets for running security scans.

- govulncheck

```bash
go install golang.org/x/vuln/cmd/govulncheck@latest
```

Then you can run the vuln make target:

```bash
make vuln
```

## Other Languages

- Haskell: <https://github.com/kevwil/gityup-haskell>
- Lua: <https://github.com/kevwil/gityup-lua>
- Python: <https://github.com/kevwil/gityup-py>

## Footnotes

[^1]: git-smart Ruby gem <https://github.com/kevwil/git-smart> which is a fork from [here](https://github.com/geelen/git-smart) with some fixes.
