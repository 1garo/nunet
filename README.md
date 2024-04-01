# Nunet Challenge

The goal is to design and implement a system that facilitates seamless communication between machines and efficiently manages container deployment.

## Getting Started

### Prerequisites

- [Go](https://go.dev/doc/install)
- [Make](https://www.gnu.org/software/make/#download)
- [Docker Compose](https://docs.docker.com/compose/install/)

### Installation

1. Clone the repository:

```bash
$ git clone git@github.com:1garo/nunet.git
$ cd nunet
```

2. Up server containers:

```bash
$ make up
```

- Remove all containers:

  ```bash
  $ make down
  ```

3. Compile the client:
```bash
$ make cli
```

4. Run tests:
```bash
$ make test
```

Checkout [Makefile](./Makefile) to see all possible commands.

## Usage

I don't recommend trying to run the `servers` locally, prefer to use `make up`, because we depend on two services running.

After `make up`, server will be running at `localhost:50051` and `localhost:50052`.

To use the client, you can run `make cli`:

#### Examples

To get a description on the existing commands
```bash
$ ./cli --help
```

You pass an address (it's the only optional argument)
```bash
$ ./cli -program=ls -args=-l -addr=localhost:50052
```

The -args arguments are comma separeted
```bash
$ ./cli -program=echo -args=nu,net 
```

To check if the commands are really being replicated, a good way is to check the logs of the containers, for example:

```bash
$ docker container logs nunet-api01-1
$ docker container logs nunet-api02-1
```

You are gonna be able to see which port the container is running on.

Run `./cli -program=echo -args=nu,net` and then check the logs to see that they were replicated on both containers.

If you want to test other way around (default is `localhost:50051` aka `api01`), just use the -addr argument `-addr=localhost:50052`
