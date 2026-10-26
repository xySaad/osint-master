# OSINT Master

**OSINT Master** (`osintmaster`) is a modular Open Source Intelligence (OSINT) command-line tool written in Go. It is designed to assist security analysts, penetration testers, and investigators in gathering publicly available information across platforms, including social media profile reconnaissance and activity tracking.

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
- [Disclaimer & Ethical Use](#disclaimer--ethical-use)

---

## Features

- **GitHub Intelligence**:
  - Retrieves profile details: display name, email, bio, company, location, blog/website, Twitter handle, follower/following counts, public repository count, and account creation date.
  - Identifies recent public activity (Push events, Pull Requests, Comments, Releases, Issues, etc.) including direct links and timestamps.
- **Instagram Intelligence**:
  - Retrieves profile metadata: full name, biography, verified badge status, account privacy status (`is_private`), post count, and follower/following metrics.
  - Extracts latest activity (post/video links and approximate timestamp resolved from Instagram snowflake IDs).
- **LinkedIn Intelligence**:
  - Gathers public intelligence.
  - Extracts profile information: full name, headline, current company (`worksFor`), location, education (`alumniOf`), follower count, connection count, and public profile URL.
  - Retrieves latest public activity (posts or published articles, publication time, and article/post URL) via structured Schema.org JSON-LD parsing.
- **Facebook Intelligence**:
  - Gathers public intelligence from public Facebook profiles and pages.
  - Extracts full name, bio/about section, follower count, "talking about" metrics, and canonical URL via link-preview metadata.
- **X (Twitter) Intelligence**:
  - Gathers public intelligence from X (formerly Twitter) profiles.
  - Extracts display name, bio/description, location, follower/following counts, post count, account creation date, verification badge status, and canonical profile URL.
- **Dual Output Streaming**:
  - Mirror scan results simultaneously to stdout and to an output file using the `-o` argument.
- **Clean CLI Interface**:
  - Intuitive subcommand and argument parsing.

---

## Project Architecture

```
osint-master/
├── go.mod
├── go.sum
├── main.go                       # CLI entrypoint and subcommand router
└── src/
    ├── lib/
    │   ├── clone_output.go       # Multi-writer stdout redirection utility
    │   └── print_if_set.go       # Conditional formatted output helper
    ├── options/
    │   └── options.go            # Argument extraction and parsing logic
    └── tools/
        └── socialmedia/
            ├── facebook.go       # Public Facebook intelligence collector
            ├── github.go         # GitHub profile and event collector
            ├── instagram.go      # Instagram profile scraper and parser
            ├── linkedin.go       # Public LinkedIn intelligence collector
            ├── socialmedia.go    # Subcommand handler for 'username'
            └── x.go              # Public X (Twitter) intelligence collector
```

---

## Prerequisites

- [Go](https://go.dev/) (v1.22+ recommended, project targets Go 1.25)
- Internet connection for fetching profile information and querying public endpoints

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

Search for account profiles and recent activity across supported platforms for a given username:

```bash
# Using go run:
go run main.go username <target_username>

# Or with compiled binary:
./osintmaster username <target_username>
```

#### Example Output:

```bash
$ go run . username satyanadella
  - Github -
Username: satyanadella
Name: Satya Nadella
Bio: CEO of Microsoft
Company: @Microsoft
Location: Redmond, Washington
Followers: 4665
Following: 2
Public Repos: 0
Created At: May 2018
Last Activity: N/A

  - Instagram -
Username: satyanadella
Full Name: Satya Nadella
Bio: Chairman and CEO of Microsoft Corporation
Verified false
Private: false
Followers: 6
Following: 0
Posts: 0
Last Activity: N/A

  - Linkedin -
Username: satyanadella
Full Name: Satya Nadella
Headline: As chairman and CEO of Microsoft, I define my mission and that of my company as…
Company: Microsoft
Location: Redmond, Washington, United States
Education: The University of Chicago Booth School of Business
Followers: 12254400
Connections: 500+
URL: https://www.linkedin.com/in/satyanadella
Last Activity:
 Type: Post
 Time: 2026-10-23 16:29:08.068 +0000 UTC
 URL: https://www.linkedin.com/posts/satyanadella_super-to-see-volkswagen-group-putting-copilot-activity-7519434747178049536-cvN-

  - Facebook -
Username: satyanadella
Name: Satya Nadella
URL: https://www.facebook.com/satya.nadella/

  - X -
Username: satyanadella
Name: Satya Nadella
Bio: Chairman and CEO at Microsoft
Verified: true
Followers: 9644511
Following: 320
Posts: 2413
Joined: February 2009
URL: https://x.com/satyanadella
```

---

### File Output (`-o`)

To save the scan results directly to a file while retaining terminal visibility:

```bash
./osintmaster username <target_username> -o report.txt
```

---

## Roadmap & Planned Modules

- [x] **`username`**: Social media profiling
  - [x] GitHub
  - [x] Instagram
  - [x] LinkedIn
  - [x] Facebook
  - [x] Twitter / X
  - [ ] Reddit
  - [ ] GitLab
- [ ] **`ip`**: IP geolocation, ASN lookup, reverse DNS, and threat intelligence.
- [ ] **`domain`**: WHOIS lookups, DNS enumeration, subdomain discovery, and SSL certificate inspection.
- [ ] **JSON Export**: Structured output format (`--json`) for automated pipeline ingestion.

---

## Disclaimer & Ethical Use

> [!WARNING]
> This tool is developed strictly for educational purposes, authorized security testing, and lawful digital investigation. Users are solely responsible for ensuring compliance with all applicable local, national, and international laws, as well as the terms of service of any third-party platforms queried. The authors assume no liability for misuse or damage caused by this software.
