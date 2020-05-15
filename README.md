# N.O.M.A.D.D.

Naïeve Omnidirectional Machine-Automated Discount Dicovery
A set of tools with the common purpose of finding great opportunities to travel. Intend to leverage Web Scraping, Data Analytics, Graph Theory, Concurrency, and more to find cheap flights, interesting travel routes, and delegate the act of searching to an automated system using minimal background resources and compact storage methods


### Motivation

Every moment, opportunities to live and experience new cultures fade away. Because it is impossible to delegate all personal time to searching for travel opportunities, cheap travel options go unnoticed. In their wake, loom expensive tickets that steal from tight budgets better spent far away than online. Instead of missing great deals, delegate the tedium of searching to an automaton

## Requirements for Local use
* go (1.13+)[https://golang.org/doc/install]
* PostgreSQL (11.7+)[https://www.postgresql.org/docs/11/install-short.html] If persistant data needed locally
* Clone this repository with `go get https://github.com/GPKyte/nomad`

## Usage
0) `cd $GO_PATH/src/github.com/GPKyte/nomad`
1) `go build`
2) `go test`
3) `go run [--OPTIONS]`

### Flags
Currently flags are not supported but would include the following
Passive and Active actions are being implemented and include:
  `--scrape`          Begins scraping default queries in the background. Compare to terminal operator & in `action &`
  `--discover`        
  `--config`          Modify settings that guide certain data collection decisions
  `--help`            Display the usage`

### Testing

Effort is taken to test 0, 1, n boundary cases, any clear expectations of transformations, and algorithmic correctness when feasible
Testing is initiated through the go standard library and CLI like so:
  `go test`
  `go test ./<packageName> --run <patternRegEx>`

## Contributions

Any contributions are welcome! Pull Request any changes, and consider writing a small test to demonstrate what you expect to change.


## DevOps:
* Configuration file (git ignore)
    - Collection of sites to scrape and their settings
* Default Configuration file (Version Control, VC)


## Display of information and use cases
* As a DEV
    - Exchangeable components with polymorphistic behavior for tasks
    - Flat repository access to custom tools
    - Sanitized data in storage
    - Persistent log files and verbosity levels
    - Low effort integration and minimal redundant with Tableau viz


* As a USER:
    - Tune threshold for tracking a listing
    - Tune filter for alerts over tracked listings
    - View listing trends over time
    - View listing trends across vendors
    - View listing trends across locations
    - View details for specific listing
    - View detail summaries for top-priority listings
    - Track listings meeting criteria for a time range and a set of destinations
    - Create routes using multiday layovers (optimal for traveling to a pit-stop city)

