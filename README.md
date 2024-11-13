# Backstage AI related administration CLI - bac

A CLI that facilitates injecting AI model metadata from various sources into the Backstage Catalog

## Contributing

All contributions are welcome. The [Apache 2 license](http://www.apache.org/licenses/) is used and does not require any 
contributor agreement to submit patches. That said, the preference at this time for issue tracking is not GitHub issues
in this repository.  

Rather, visit the team's current [Jira project](https://issues.redhat.com/projects/DEVAI/issues).

See [the development guide](docs/DEVELOPMENT.md) for details on how to build and test any contributions you make.

## Usage

At a high level, the `bac` CLI

- Provides for the generation of YAML formatted definitions of Backstage `Components`, `Resources`, and `APIs` catalog entities by accessing external systems that provide AI model metadata
- Which external systems are supported is expected to grow over time, at least in the short term.
- Once that YAML information is stored in a HTTP accessible file, the `bac` CLI then provides commands to instructs a specific Backstage instance to import those entities into its catalog.  This will show up as a Backstage `Location` in the catalog, where it is a parent of the `Components`, `Resources` and `APIs`
- Those `Components`, `Resources`, and `APIs` will have specific AI related `types` which will allow for distinguishing from other `Components`, `Resources` and `APIs` in the catalog
- Later on, it also allows for the deletion of Backstage `Locations` and any `Components`, `Resources`, and `APIs` defined by that `Location`
- Lastly, the `bac` CLI allows for retrieving any AI related `Components`, `Resources` and `APIs`

To received detailed usage information and example invocations, after building the `bac` executable, you can run
```shell
bac help
```

This invocation will also provide the current list of subcommands.  Similarly, running 
```shell
bac help <subcommand>
bac help <subcommand> <subcommand>
```
will provide usage information, example invocations, optional flags, and additional subcommands for the current list of subcommands.

## Potential tl;dr

First, our [background document](docs/background.md) gets into the scenarios and personas we are targeting with this CLI,
as well as rationale for the syntax, language(s), and the like.

Then, our [roadmap document](docs/roadmap.md) provides a snapshot of the more immediate changes we have planned with, 
Jira references when ideas reach sufficient priority to warrant official tracking.
