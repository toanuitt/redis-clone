# Redis Clone

A lightweight Redis clone implemented in Go that supports core Redis functionality including:

- Basic key-value operations
- Key expiration
- Sorted sets
- Sets
- Count-Min Sketch (CMS)
- Bloom filters
- Memory eviction policies (LRU, LFU, Random)

## Features

### Data Structures
- Key-value store with TTL support
- Sorted Sets (using B+ Tree)
- Simple Sets
- Count-Min Sketch
- Bloom Filter 

### Memory Management
- Key eviction policies:
  - allkeys-random
  - allkeys-lru (Least Recently Used)
  - allkeys-lfu (Least Frequently Used)
- Configurable eviction pool size
- TTL-based key expiration

### Server Features
- TCP server using I/O multiplexing (epoll on Linux, kqueue on macOS)
- RESP (Redis Serialization Protocol) protocol support
- Graceful shutdown handling
- Configurable connection limits

## Getting Started

### Prerequisites
- Go 1.22 or higher

### Installation
```bash
git clone https://github.com/yourusername/redis-clone.git
cd redis-clone
go build [main.go](http://_vscodecontentref_/1)