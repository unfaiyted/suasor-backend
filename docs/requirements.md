# URL Shortener API Technical Specification

## 1. System Architecture

- RESTful API built with Go
- Database: PostgreSQL for persistent storage
  <!--- Load Balancer for high availability-->

## 2. API Endpoints

```http

POST /api/v1/shorten
Request Body:{
    "fullUrl": "https://example.com/very/long/url",
    "customCode": "optional-custom-code"  // optional
}

Response:{
    "shortUrl": "http://short.url/abc123",
    "longUrl": "https://example.com/very/long/url",
    "createdAt": "2023-01-01T00:00:00Z",
    "expiresAt": "2024-01-01T00:00:00Z"
}

Redirect to Original URL
GET /{short_code}


Get URL Statistics
GET /api/v1/urls/{short_code}/stats

Response:{
    "shortCode": "abc123",
    "clicks": 1000,
    "createdAt": "2023-01-01T00:00:00Z",
    "lastAccessed": "2023-01-01T00:00:00Z"
}


```

## 3. Technical Requirements

URL Generation
Generate 6-8 character alphanumeric codes
Check for collisions before assigning
Support custom codes (optional)
Performance Requirements

## 4. Performance Goals

Response time: < 100ms for redirects
Availability: 99.9%
Support for 1000+ requests per second

## 5. Security Requirements

Rate limiting: 100 requests per hour per IP
URL validation to prevent malicious links
HTTPS enforcement
Input sanitization

## Core Components

### HTTP Server

RESTful API implementation
Request routing
Error handling

### Middleware

Logging
Rate limiting
Request validation

## Database Operations

Connection pooling
Transaction management
Query optimization
Migration management

## Stretch Goal - Caching Layer

Redis implementation
Cache invalidation strategy
TTL management

## Development Stack

Go 1.16+
PostgreSQL 13+
Redis 6+
Docker
Nginx (Load Balancer)

## Project Timeline

### Phase 1: Core Implementation

Basic API setup
Database implementation
URL shortening logic
Basic redirect functionality

### Phase 2: Enhanced Features

Analytics tracking
Custom URLs
API documentation
Monitoring setup

### Phase 3: Production Readiness

Caching layer
Security hardening
Performance optimization
Production deployment
