# mws-api

[![Build Status](https://travis-ci.org/MathWebSearch/mwsapi.svg?branch=master)](https://travis-ci.org/MathWebSearch/mwsapi)

A [golang](https://golang.org) library and set of tools to setup, query and maintain a MathWebSearch and Temasearch instance.

## Overview

- `cmd`: Implementation of commands
  - `cmd/temaquery`: Queries A joined MathWebSearch + ElasticSearch Setup
  - `cmd/mwsapid`: HTTP Daemon serving temasearch queries
  - `cmd/mwsquery`: Queries a (plain) MathWebSearch instance for MathWebSeach Queries
  - `cmd/elasticquery`: Queries a (plain) Elasticsearch instance for Tema Queries
  - `cmd/elasticsync`: Creates and maintains an Elasticsearch instance for use with Temasearch
- `connection`: Contains Connection Code to MathWebSearch and ElasticSearch
- `engine`: Underlying code used by commands above
  - `engine/elasticsync`: Creates and maintains an Elasticsearch instance for use with Temasearch
  - `engine/elasticengine`: ElasticSearch only queries
  - `engine/mwsengine`: MathWebSearch queries
  - `engine/temaengine`: TemaSearch Queries
- `query`: Implements query parsing and serializing
- `result`: Implements result parsing and serializing
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
In the following we call each such document a _Harvest Element_.
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
We call each file a _segment_.
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

## API Daemon

The program in `cmd/mwsapid` implements an HTTP Daemon that can answer all the queries above.
It can be configured using the command line arguments:

```
Usage of mwsapid:
  -elastic-host string
        Host to use for elasticsearch
  -elastic-port int
        Port to use for elasticsearch (default 9200)
  -mws-host string
        Host to use for mathwebsearch. If omitted, disable mathwebsearch support
  -mws-port int
        Port to use for mathwebsearch (default 8080)
  -host string
        Host to listen on (default "localhost")
  -port int
        Port to listen on for queries (default 3000)
```

### General Structure

The server supports three kinds of requests -- which are described in more detail below.
When using `POST` requests, all parameters should be encoded using JSON in the body.

For all requests, the server will respond with a JSON object in the body.
By default, this corresponds to a simple `application/json` response.
However, when the URL-parameter `callback` is provided, a [`JSONP`](https://web.archive.org/web/20160304044218/http://www.json-p.org/) response is sent instead.

Furthermore, the server makes use of the following status codes:

| Code | Description           | Meaning                                                                        |
| ---- | --------------------- | ------------------------------------------------------------------------------ |
| 200  | OK                    | Request suceeded and the body will contain the desired response.               |
| 400  | Bad Request           | Malformed request, this occurs if some parameters are out of range or missing. |
| 404  | Not Found             | The given request was not found or is not supported by the server.             |
| 405  | Method Not Allowed    | The request method (i.e. POST or GET) is not allowed for the given request.    |
| 500  | Internal Server Error | Something went wrong when trying to answer the query.                          |

When responding with a non-200 status code, the body will always contain a JSON string with a detailed error message.
This message is not intended for end users, instead it should be used by developers to debug the issue at hand.

### Search Result Serialization

#### Result

A search results are represented using [the Result struct](result/result.go) as follows:

|   Field   |          Type          | Optional |                                                                Description                                                                |
| --------- | ---------------------- | -------- | ----------------------------------------------------------------------------------------------------------------------------------------- |
| kind      | `string`               | no       | Type of response returned. One of `"mwsd"`, `"tema"`, `"elastic"`, `"elastic-document"` or `""`. The `elastic` ones are currently unused. |
| took      | `number`               | yes      | Time it took to run the query in Nanoseconds. Optional.                                                                                   |
| stats     | `Dict<string, number>` | yes      | Time in nanoseconds spent in specific phases. Component names may change in the future.                                                   |
| total     | `number`               | no       | Total number of results for the query, regardless how many are returned.                                                                  |
| from      | `number`               | no       | 0-based number this set of results starts at.                                                                                             |
| size      | `number`               | no       | Number of results returned.                                                                                                               |
| variables | `Array<QueryVariable>` | yes      | Query Variables found within the original query. See detailed description below.                                                          |
| ids       | `Array<number>`        | yes      | Internal result ids, when requested.                                                                                                      |
| hits      | `Array<Hit>`           | yes      | The list of matching hits. See detailed description below.                                                                                |

#### Hit

A hit is represented using [the Hit struct](result/hit.go) as following:

|  Field   |         Type         | Optional |                                              Description                                              |
| -------- | -------------------- | -------- | ----------------------------------------------------------------------------------------------------- |
| id       | `string`             | yes      | (Possibly internal) id of this hit.                                                                   |
| url      | `string`             | yes      | Url of this hit.                                                                                      |
| xpath    | `string`             | yes      | Xpath of the query term to the formulae referred to by this hit.                                      |
| element  | `HarvestElement`     | yes      | Harvest element (`aka <mws:data> element`) belonging to this hit. See below for detailed description. |
| score    | `number`             | yes      | Score of this Hit as determined by ElasticSearch.                                                     |
| snippets | `Array<string>`      | yes      | Snipets that caused this hit to gain the score. TemaSearch only.                                      |
| math_ids | `Array<MathFormula>` | no       | Formulae found within this hit. See detailed description below.                                       |

#### HarvestElement

An `<mws:data>` element is represented using [the HarvestElement struct](result/harvest.go) as follows:

|  Field   |                   Type                   | Optional |                                                                                   Description                                                                                    |
| -------- | ---------------------------------------- | -------- | -------------------------------------------------------------------------------------------------------------------------------------------------------------------------------- |
| metadata | `any`                                    | no       | Metadata belonging to this element. When representing a valid JSON object, returns that JSON object. When empty, an empty JSON obeject. Otherwise a string representing the XML. |
| segment  | `string`                                 | yes      | Name of the segment (i.e. `.harvest` file) this element comes from.                                                                                                              |
| text     | `string`                                 | no       | Text contained in this document. Might contained replaced math ids.                                                                                                              |
| mws_id   | `Dict<number, Dict<string, MathFormula>` | yes      | Information about each math identifier within this document, and a map from the internal mws ids to the appropriate formula. See below for detailed description.                 |
| mws_ids  | `Array<number>`                          | yes      | List of math identifiers in this documents. Corresponds to the keys of `mws_id`.                                                                                                 |
| math     | `Dict<string, string>`                   | no       | Source code of replaced math elements within this document.                                                                                                                      |

#### MathFormula

A math formula is represented using [the MathFormula struct](result/formula.go) as following:

|  Field  |          Type          | Optional |                          Description                          |
| ------- | ---------------------- | -------- | ------------------------------------------------------------- |
| source  | `string`               | yes      | MathML Element Representing entire formula.                   |
| durl    | `string`               | yes      | Document URL this formula is contained in.                    |
| url     | `string`               | yes      | Local Formula Id                                              |
| xpath   | `string`               | yes      | XPath from the formula to the matching subterm.               |
| subterm | `string`               | yes      | MathML Element representing the matching subterm.             |
| subst   | `Dict<string, string>` | yes      | MathML Elements representing values for the subsituted terms. |

#### QueryVariable

A query variable is represented using [the QueryVariable struct](result/variable.go) as following:

| Field |   Type   | Optional |                       Description                        |
| ----- | -------- | -------- | -------------------------------------------------------- |
| name  | `string` | no       | Name of this Query Variable                              |
| xpath | `string` | no       | XPath from the root of the query to the variable itself. |

### Status Request

The Status Handler is called running a GET on `/`.
It takes no parameters and returns a [StatusResponse](engine/server.go) with the following structure:

|  Field  |         Type         | Optional |              Description               |
| ------- | -------------------- | -------- | -------------------------------------- |
| name    | `string`             | no       | Name of this server. Always `mwsapid`. |
| tagline | `string`             | no       | Server Tagline.                        |
| engines | `Dict<string, bool>` | no       | Supported "engines" or routes.         |

Example Response:

```json
{
  "name": "mwsapid",
  "tagline": "You know, for math",
  "engines": {
    "mws": true,
    "tema": false
  }
}
```

### MathWebSearch Request

The MWS Handler is called running a POST on `/mws/`.
It takes parameters of type [MWSAPIRequest](engine/mwsengine/handler.go).

|    Field    |      Type       | Optional |                                                                  Description                                                                  |
| ----------- | --------------- | -------- | --------------------------------------------------------------------------------------------------------------------------------------------- |
| expressions | `Array<string>` | yes      | List of MathWebSearch expressions. Each should be a the body of a single "mws:expression" tag, using the "mws" and "m" predefined namespaces. |
| mwsids      | `boolean`       | yes      | When true, do not return MathWebSearch results, but only their IDs                                                                            |
| count       | `boolean`       | yes      | When true, return only count of results, not results themselves.                                                                              |
| from        | `number`        | yes      | Used for pagination. 0-based index to start result set at. Defaults to 0, must be >= 0.                                                       |
| size        | `number`        | yes      | Used for pagination. Maximum number of results to returns. Defaults to 10, must be between 0 and 100 inclusive.                               |

For example, when the server is running on localhost at port 3000, the following curl command could be used to make a simple request:

`curl -d '{"expressions":["<mws:qvar>x</mws:qvar>"]}' -H "Content-Type: application/json" -X POST http://localhost:3000/mws/`

If the count parameter is true, the response will be a single json number.
Otherwise, the server will return a [Result](result/result.go) struct, see the above section for details on how this looks.
This behaviour is identical to the `mwsquery` executable.
Example responses can be found in the [cmd/mwsquery/cmd/testdata](cmd/mwsquery/cmd/testdata) folder.

### TmeaSearch Request

Not yet documented.

## Docker

For convenience, a Dockerfile serving the `API` daemon is provided.
It can be found at the automated build [mathwebsearch/mwsapi](https://hub.docker.com/r/mathwebsearch/mwsapi) on DockerHub.
It can be run as follows:

```
docker run mathwebsearch/mwsapi
```

It serves the API Daemon (see above) on port 3000 by default and can be customized using the following environment variables:

- MWSAPID_HOST: Host to listen for requests. Defaults to "0.0.0.0".
- MWSAPID_PORT: Port to listen for requests. Defaults to 3000.

- MWSAPID_MWS_HOST: Host to expect MathWebSearch Daemon on. If omitted, MathWebSearch support is disabled.
- MWSAPID_MWS_PORT: Port to expect MathWwebSearch Daemon on. Defaults to 8080.

- MWSAPID_ELASTIC_HOST: Host to expect Elasticsearch Daemon on. If omitted, TemaSearch support is disabled.
- MWSAPID_ELASTIC_PORT: Port to expected Elasticsearch on. Default to 9200.

Furthermore, a Docker Image for elasticsync also exists.
See [MathWebSearch/tema-elasticsync](https://github.com/MathWebSearch/tema-elasticsync) for details.

## License

GPL3, see [LICENSE](LICENSE).
