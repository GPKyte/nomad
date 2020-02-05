# Discount Nomad

A set of tools with the common purpose of finding great opportunities to travel. Intend to leverage Web Scraping, Data Analytics, Graph Theory, Concurrency, and more to find cheap flights, interesting travel routes, and delegate the act of searching to an automated system using minimal background resources and compact storage methods

### Motivation

Every moment, opportunities to live and experience new cultures fade away. Because it is impossible to delegate all personal time to searching for travel opportunities, cheap travel options go unnoticed. In their wake, loom expensive tickets that steal from tight budgets better spent far away than online. Instead of missing great deals, delegate the tedium of searching to an automaton

## Requirements

## Usage

## Contributions

Any contributions are welcome! Follow some of these guide lines and suggestions and Pull Request any changes you've made and tested

  Polymorphism
  Don’t Repeat Yourself (DRY)
	Build own Data Struct & Algorithms where practical and non-distracting

## Results

Possible Modules + API and structure of project
* SecureCredentials Implements CleanDataConstant
* RawListing Implements DirtyDataConstant, GenericListing
* CleanListing Implements CleanDataConstant, GenericListing
* ExceptionSiteListingAsAnExample Extends CleanListing

* TraversableLocation Implements GenericGraphNode
* TraversableRoute Implements GenericGraphWalk

* Parse.HTML
* Parse.Listing
* Parse.XYZ
* Logger (Should be abstracted even if logging with a good std libr)

DevOps:
* Ansible to setup and run tasks
* Configuration file (git ignore)
    - Collection of sites to scrape and their settings
* Default Configuration file (Version Control, VC)


Display of information and use cases
* As a DEV
    - I want global var value tracing
    - Backwards debugging
    - Exchangeable components with polymorphistic behavior for tasks
    - Flat repository access to custom tools
    - Sanitized data in storage
    - Security for credentials
    - Process monitoring
    - Persistent log files and verbosity levels
    - Encrypt logs as filtered (e.g. RSA with hidden keys)
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
    - Create complex routes with optimal layovers (optimal for traveling to a pit-stop city)


Development ideas:
    Pair programming?? Help other Sen Sem projects/ be soundboard
    Use AI + Structured NLP to identify key information in page?
    Structure of DB? Relational?


Norms:
"Data, POJO, immutable data" What to call this in naming consistently?

