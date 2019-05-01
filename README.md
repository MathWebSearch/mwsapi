# mws-api

[![Build Status](https://travis-ci.org/MathWebSearch/mwsapi.svg?branch=master)](https://travis-ci.org/MathWebSearch/mwsapi)

A [golang](https://golang.org) library and set of tools to setup, query and maintain a MathWebSearch and Temasearch instance. 

## Overview

- `cmd`: Implementation of commands
    - `cmd/temaquery`: Queries A joined MathWebSearch + ElasticSearch Setu
    - `cmd/mwsapid`: HTTP Daemon serving temasearch queries (not yet implemented)
    - `cmd/mwsquery`: Queries a (plain) MathWebSearch instance for MathWebSeach Queries
    - `cmd/elasticquery`: Queries a (plain) Elasticsearch instance for Tema Queries
    - `cmd/elasticsync`: Creates and maintains an Elasticsearch instance for use with Temasearch
- `connection`: Contains Connection Code to MathWebSearch and ElasticSearch

(Rest is not up-to-date)

- `temasearch`: Code implementing joined (mws, tema) queries
- `mws`: Code interacting with MathWebSearch
- `tema`: Code interacting with Elasticsearch / TemaSearch
    - `tema/sync`: Code to synchronize a TemaSearch index with Elasticsearch
    - `tema/query`: Code to query TemaSearch
- `elasticutils`: Utility wrappers around Elasticsearch objects
- `utils`: General utility functions

## Processes

In the following we describe the basic functionality between all programs within this repository. 
This documentation is intended to serve as an entry point, and thus does not describe all implementation details. 
The most detailed reference is always the source code. 

In principle, the source code for all commands is found in the appropriate `cmd/` subdirectory. 
Executables can be built by either using thhe standard `go build ./cmd/$CMD` or by simply using `make $CMD`. 
The binaries will be placed in the root directory. 

### MathWebSearch Query

The program in `cmd/mwsquery` can run plain MathWebSearch Queries. 
Queries are defined by the [Query Struct](mws/wrapper.go). 

Each Query consists of a list of MathWebSearch Expressions.
A MathWebSearch Expression is a ContentMathML expression in XML Syntax with additional support for Query variables. 
For this one can use `<mws:qvar>` tags to specify universal variables. 
The text content of a `<mws:qvar>` is considered its name and qvars with the same name will match to the same expressions. 

Each Expression can be given as an argument to the `mwsquery` executable.
For example:

`./mwsquery '<mws:qvar>x</mws:qvar>'`

Normal results are returned as JSON to STDOUT. 
The results are defined by the [Result Struct](mws/result.go). 

All queries are paginated -- by default they return the first 10 results.
The parameters `-from` and `-size` can be used to customize the result set. 

Sometimes it is only important how many results are returned, not the results themselves.
For this purpose the `-count` flag can be provided. 

Additionally, instead of returning the full results, sometimes it is also desired to only return the ids of each found formulae. 
This can be useful for debugging and use inside a full TemaSearch scenario.
To achieve this, the `-ids` flag can be provided. 

### Elasticsearch Syncronization

The program in `cmd/elasticsync` creates and maintains an Elasticsearch Index for use with Tema-Search. 

A Temasearch Index is a set of JSON objects conforming to the [Temasearch Harvest Element Schema](tema/Element.go).
In the following we call each such document a *Harvest Element*. 
A group of elements belonging to a single .harvest file (which in and of itself belongs to one source file) is usually contained in one line of an elasticsearch index file ending in .json. 
For backward compatibility, in between each lines of items in the index, an additional document containing legacy meta-information should be stored.
These .json files are stored within one folder on disk. 

In order to make this index queryable, it needs to be kept in sync with an appropriate Elasticsearch index. 
To achieve this one could in principle perform the following process to syncronize the disk with the index:

- Delete all existing indexed documents from Elasticsearch (if any)
- Read each `.json` file from disk and then
- add the documents contained inside of it to Elasticsearch

This approach does not scale well with large datasets. 
Having to delete the entire database, only to add the same content back is too slow.

Instead we split the Temasearch index into into the different files and treat each file seperatly. 
We call each file a *segment*. 
To syncronize an updated on-disk index into Elasticsearch, we roughly do the following:

- Mark all existing segments in the database as 'untouched'
- For each segment from the ElasticSearch index to be added:
    - compute a hash of the segment
    - check if this segment with the same name is already stored in the database by comparing the hash
        - if yes, we do not need to do anything as it has not changed
        - if no, we remove the old segment documents (if any) and add the new documents belonging to this hash
    - mark the segment as 'touched' within this syncronization process
- Delete the documents belonging to any segment still marked as 'untouched'

This process is far more efficient -- only updating documents in the database that have actually been changed.
However, this process requires that two seperate ElasticSearch indexes are maintained. 
The first index -- called tema by convention -- contains the TemaSearch Index Documents and is most obvious. 
The second index is called tema-segments and contains a list of known segments as well as their hashes. 
As a hash implementation we use SHA256.

### Elasticsearch Query

The program in `cmd/elasticquery` can run queries against the elasticsearch part of Temasearch. 
Queries are defined by the [Query Struct](tema/query/main.go) and consist of two parameters:

- Some text to search the index for
- A list of mathwebsearch ids that were found by normal MathWebSearch

A query may have both text and ids to search for, but it must not be empty. 
These can be provided to `elasticquery` using the `text` and `ids` parameters. 
For example:

`./elasticquery -text "Hello" -ids 1,2,3`

Normal results are returned as JSON to STDOUT. 
The results are defined by the [Result Struct](tema/query/main.go). 

All queries are paginated -- by default they return the first 10 results.
The parameters `-from` and `-size` can be used to customize the result set. 

Sometimes it is only important how many results are returned, not the results themselves.
For this purpose the `-count` flag can be provided. 

Internally, each query consists of two phases:

- The Document Phase. This intially queries elasticsearch to find all matching documents. 
- The Highlight Phase. For each returned document, elasticsearch is queried again to highlight matching segments. 

A normal query runs both phases. 
For debugging, it is possible to only run the document phase by running the `-document-phase-only` flag. 

## Docker

For convenience, a Dockerfile serving the `API` daemon is provided. 
It can be found at the automated build [mathwebsearch/mwsapi](https://hub.docker.com/r/mathwebsearch/mwsapi) on DockerHub. 
It can be run as follows:

```
docker run mathwebsearch/mwsapi
```

Furthermore, a Docker Image for elasticsync also exists. 
See [MathWebSearch/tema-elasticsync](https://github.com/MathWebSearch/tema-elasticsync) for details. 

## License

GPL3, see [LICENSE](LICENSE). 