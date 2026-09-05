# Audit Log Service

A small HTTP service for recording audit events... Who did what, to what, when, and from where. Then queries them back out again. Written in Go.

## The Idea

I wanted to learn Go, but my learning style is very much the "hands-on" approach. I learn best when I'm *doing*, so I brainstormed a small project that would let me lean on the skills I already have while picking up a new language from scratch.

And so, BEHOLD! My audit log service. It's being built in stages, and this README grows with it.

## Status

Early days! Right now this is an HTTP server that starts up, routes a request, and answers it. Small beginnings, but it runs.

**Built so far**

- HTTP server with its own router and explicit server configuration

**Planned**

- Accepting structured audit events as JSON over POST
- In-memory storage behind a swappable interface
- Asynchronous ingestion with graceful shutdown
- Query endpoints — filter by actor, resource, and time range
- Persistent storage with tamper-evident hash chaining

## Running It

Requires Go 1.27 or newer.

```
go run ./cmd/auditd
```

Then, in another terminal:

```
curl localhost:8080
```

## Design Notes
**Custom ServeMux VS DefaultServeMux**

I decided to go with making my own ServeMux rather than using DefaultServeMux.

Why?

A package-level global is shared mutable state that any package in my dependency tree can write to at init time, without my knowledge.

I.e. if I or any package I use imports pprof and I use DefaultServeMux for my server, pprof registers its routes on DefaultServeMux... Which means those routes are on a live, publicly available server. And one of those routes is dumping my program's memory where anyone can see. Not hip, not cool.

THAT BEING SAID... DefaultServeMux isn't bad, it's just a global. And globals are risky business in any kind of engineering.

## AI Usage

The implementation is mine — every line of Go in this repository was written by me.

AI assisted with scaffolding, tooling, and review: generating the module file, creating empty package directories, running git commands, and acting as a sounding board while I learned the language. I think that's a reasonable use of a good tool without defeating the purpose of making this project.

Thank you for reading!

— Charley ヾ( ˃ᴗ˂ )◞ • *✰
