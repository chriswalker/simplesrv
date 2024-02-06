# simplesrv
`simplesrv` is a [Go project template](https://go.dev/blog/gonew) for a bare-bones HTTP server, with a basic layered architecture, a backing database ([SQLite](https://sqlite.com/index.html)) and basic unit tests in a few places - not comprehensive, but enough to get things going when starting up a new templated project.

## Building and running `simplesrv`
The `run.sh` Bash script takes the place of excessive Makefiles. To build:

```sh
$ ./run.sh build
```

which results in the `simplesrv` binary being built into the `bin/` directory.

```sh
$ ./run.sh run
```

will run the server, by default on port 8080 and using an SQL database file called `app.db`.

For a list of tasks in the Run script, run:

```sh
$ ./run.sh help
```

Note the only third-party dependency is [golangci-lint](https://golangci-lint.run/) for the `lint` task.

## Using as a project template
1. Install `gonew`:

    ```sh
    $ go install golang.org/x/tools/cmd/gonew@latest
    ```

    (see [here](https://go.dev/blog/gonew) for more information)

2. Then checkout `simplesrv` as a new templated project by:

    ```sh
    $ gonew github.com/chriswalker/simplesrv example/mysrv
    $ cd ./mysrv
    # aaaand work away
    ```
