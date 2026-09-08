# Audit Log Service

A small HTTP service for recording audit events... Who did what, to what, when, and from where. Then queries them back out again. Written in Go.

## The Idea

I wanted to learn Go, but my learning style is very much the "hands-on" approach. I learn best when I'm *doing*, so I brainstormed a small project that would let me lean on the skills I already have while picking up a new language from scratch.

And so, BEHOLD! My audit log service. It's being built in stages, and this README grows with it.

## Status

Still the early days! Right now the service accepts audit events over HTTP, validates them,
fills in what the server is responsible for, and prints them. Nothing is stored yet —
events are logged to the terminal and dropped. Storage is next.

**Built so far**

- HTTP server with its own router and explicit server configuration
- `POST /events` accepting structured audit events as JSON
- Separate wire and stored types, so clients can't set server-owned fields
- Server-stamped `received_at`, with `occurred_at` falling back to it and being flagged when absent
- Validation of required fields, rejecting malformed or incomplete events with a 400

**Planned**

- In-memory storage behind a swappable interface
- Asynchronous ingestion with graceful shutdown
- Query endpoints: filter by actor, resource, and time range
- Persistent storage with tamper-evident hash chaining
  
## Running It

Requires Go 1.27 or newer.

```
go run ./cmd/auditd
```

The server listens on `:8080`. In another terminal, send an event:

```
curl.exe -X POST localhost:8080/events -H "Content-Type: application/json" -d "@testdata/event_valid.json"
```

You should get back `OK!`, and the server terminal prints the event with its server-stamped `received_at`.

To see the fallback in action, send one with no `occurred_at`:

```
curl.exe -X POST localhost:8080/events -H "Content-Type: application/json" -d "@testdata/event_missing_occurred_at.json"
```

The event comes back with `occurred_at` equal to `received_at` and `occurred_at_inferred` set to true.

A malformed body, or one missing `user`, `action`, or `resource`, gets a 400.

On Windows, use `curl.exe` rather than `curl` — PowerShell aliases the latter to a different tool. On macOS and Linux, plain `curl` is correct.

## Design Notes
**Custom ServeMux VS DefaultServeMux**

I decided to go with making my own ServeMux rather than using DefaultServeMux.

Why?

A package-level global is shared mutable state that any package in my dependency tree can write to at init time, without my knowledge.

I.e. if I or any package I use imports pprof and I use DefaultServeMux for my server, pprof registers its routes on DefaultServeMux... Which means those routes are on a live, publicly available server. And one of those routes is dumping my program's memory where anyone can see. Not hip, not cool.

THAT BEING SAID... DefaultServeMux isn't bad, it's just a global. And globals are risky business in any kind of engineering.

-----------

**Explicit http.Server struct, not bare ListenAndServe** 

Constructing the struct is what makes timeouts and graceful shutdown possible in the future. Taking the tutorial shortcut costs a rewrite later.

-----------

**Two timestamps** 

`occurred_at` is the client's claim, recorded but never trusted. `received_at` is stamped on arrival and is the one you can vouch for. They diverge legitimately (batching, queue lag) and illegitimately (bad clocks, backdating), so both are needed.

Missing `occurred_at` falls back to `received_at`, and gets flagged. Losing evidence is an audit log's worst failure, so a less precise record beats a dropped one. `OccurredAtInferred` keeps it honest rather than laundering a guess into a client claim. 

Cost: consumers now have to know about the flag, or a time-range query silently mixes reported and inferred times.

-----------

**Separate wire type and stored type** 

`EventRequest` holds only what a client may set; `Event` adds server-owned fields. If those were decodable, a client could send an invented timestamp and mark it verified. 

Related rule: never mutate the request in place, since that destroys the record of what was actually claimed.

-----------

**Event embeds EventRequest rather than duplicating fields** 

Chosen partly to exercise Go's composition mechanism.

Tradeoff: it couples the stored schema to the wire schema. 

Exit condition: switch to duplication once they need to diverge, most likely when records get hash-chained.

## AI Usage

The implementation is mine. Every line of Go in this repository was written by me. I also made the design decisions and worked out the reasoning behind them.

AI assisted with everything around the code: generating the module file, creating empty package directories, running git commands, keeping track of decisions I'd made across sessions so I didn't lose them, and helping draft and tidy the wording of this README. It also acted as a reviewer and a sounding board while I learned the language (i.e. it argued with my choices, which is why the Design Notes exist).

I think that's a reasonable use of a good tool without defeating the purpose of making this project.

Thank you for reading!

— Charley ヾ( ˃ᴗ˂ )◞ • *✰
