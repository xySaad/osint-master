# OSINT Master

**OSINT Master** (`osintmaster`) is a modular Open Source Intelligence (OSINT) command-line tool written in Go. It is designed to assist security analysts, penetration testers, and investigators in gathering publicly available information across platforms, starting with social media profile reconnaissance and activity tracking.

---

## Table of Contents

- [Features](#features)
- [Project Architecture](#project-architecture)
- [Prerequisites](#prerequisites)
- [Installation & Build](#installation--build)
- [Usage](#usage)
  - [Global Syntax](#global-syntax)
  - [Username Reconnaissance (`username`)](#username-reconnaissance-username)
  - [File Output (`-o`)](#file-output--o)
- [Roadmap & Planned Modules](#roadmap--planned-modules)

---

## Features

- **GitHub Intelligence**:
  - Retrieves profile details: display name, email, bio, company, location, blog/website, Twitter handle, follower/following counts, public repository count, and account creation date.
  - Identifies recent public activity (Push events, Pull Requests, Comments, Releases, Issues, etc.) including direct links and timestamps.
- **Instagram Intelligence**:
  - Retrieves profile metadata: full name, biography, verified badge status, account privacy status (`is_private`), post count, and follower/following metrics.
  - Extracts latest activity (post/video links and approximate timestamp resolved from Instagram snowflake IDs).
- **Dual Output Streaming**:
  - Mirror scan results simultaneously to stdout and to an output file using the `-o` argument.
- **Clean CLI Interface**:
  - Intuitive subcommand and argument parsing.

---

## Prerequisites

- [Go](https://go.dev/) (v1.22+ recommended, project targets Go 1.25)
- Internet connection for fetching profile information and querying APIs

---

## Installation & Build

### 1. Clone the repository

```bash
git clone https://github.com/zone01oujda/osint-master.git
cd osint-master
```

### 2. Download dependencies

```bash
go mod download
```

### 3. Build binary

```bash
go build -o osintmaster main.go
```

You can now run `./osintmaster` directly or install it globally to your `$GOPATH/bin`:

```bash
go install
```

---

## Usage

### Global Syntax

```bash
osintmaster [COMMAND] ...args
```

#### Global Arguments:

- `-o <FileName>`: Mirror console output to the specified file.
- `--help`: Display usage and help message.

---

### Username Reconnaissance (`username`)

Search for account profiles and recent activity associated with a given username:

```bash
# Using go run:
go run main.go username <target_username>

# Or with compiled binary:
./osintmaster username <target_username>
```

#### Example Output:

```text
  - Github -
Username: octocat
Name: The Octocat
Company: @github
Location: San Francisco
Blog: https://github.blog
Followers: 24706
Following: 9
Public Repos: 8
Created At: Jan 2011
Last Activity: N/A

  - Instagram -
Username: octocat
Verified false
Private: false
Followers: 5
Following: 0
Posts: 3  - Instagram -
Username: octocat
Full Name: The Octocat
Bio: Official Instagram profile of Octocat
Verified: true
Private: false
Followers: 45.2K
Following: 120
Posts: 85
Last Activity:
 Type: Post
 Time: 2026-02-10 18:30:15 +0000 UTC
 URL: https://www.instagram.com/p/C_abc123/
Last Activity:
 Type: Post
 Time: 2012-07-18 19:07:49.686 +0100 +01
 URL: https://www.instagram.com/p/NO1BdXHJ3a
```

---

### File Output (`-o`)

To save the scan results directly to a file while retaining terminal visibility:

```bash
./osintmaster username <target_username> -o report.txt
```

---

## Roadmap & Planned Modules

- [x] **`username`**: Social media profiling (GitHub, Instagram).
- [ ] **`username` expansion**: Twitter/X, Reddit, LinkedIn, GitLab.
- [ ] **`ip`**: IP geolocation, ASN lookup, reverse DNS, and threat intelligence.
- [ ] **`domain`**: WHOIS lookups, DNS enumeration, subdomain discovery, and SSL certificate inspection.

---
