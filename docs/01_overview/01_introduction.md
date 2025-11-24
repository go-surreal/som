# Introduction

SOM (SurrealDB Object Mapper) is an ORM and query builder for [SurrealDB](https://surrealdb.com/) with built-in model mapping and type-safe query operation generator. It provides an easy and sophisticated database access layer.

## What is SurrealDB?

SurrealDB is a cutting-edge database system that offers a SQL-style query language with real-time queries and efficient related data retrieval. It supports both schema-full and schema-less data handling.

With its full graph database functionality, SurrealDB enables advanced querying and analysis by allowing records (or vertices) to be connected with edges, each with its own properties and metadata. This facilitates multi-table, multi-depth document retrieval without complex JOINs, all within the database.

*(Information extracted from the [official homepage](https://surrealdb.com))*.

## Why SOM?

SOM provides:

- **Type-safe queries**: Compile-time checked queries prevent runtime errors
- **Code generation**: Automatically generates repository and query builder code from your Go models
- **Native Go experience**: Work with your domain models directly, no manual mapping required
- **Graph support**: First-class support for SurrealDB's graph capabilities via edges
