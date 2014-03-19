webhook-monitor
===============

**Features:**
* Written in Go
* Listens on a port for webhook posts
* Uses referrer info to know what to do
* Keeps track of the last time each endpoint was called to prevent DOS attack
* Minimum interval defaults to 5 minutes (configurable with JSON file?)