# OSINT Master

**OSINT Master** (`osintmaster`) is a modular Open Source Intelligence (OSINT) command-line tool written in Go for gathering public information across social media platforms and domain infrastructure.

---

## Prerequisites

- [Go](https://go.dev/) (v1.22+)

---

## Setup

```bash
# Clone the repository
git clone https://github.com/xySaad/osint-master
cd osint-master

# Download dependencies
go mod download

# Build binary
go build -o osintmaster main.go
```

---

## Usage

```bash
./osintmaster [command] [arguments] [options]
```

### Global Options

- `-o <FileName>`: Save output to a file while displaying it in the terminal.
- `--help`: Show help and usage information.

---

### `username`

Searches public profiles and activity across GitHub, Instagram, LinkedIn, Facebook, and X (Twitter).

#### Syntax

```bash
./osintmaster username <target_username> [-o <FileName>]
```

#### Arguments & Options

| Argument / Option   | Required | Description                                    |
| ------------------- | -------- | ---------------------------------------------- |
| `<target_username>` | Yes      | Username to search across supported platforms. |

#### Example

```bash
./osintmaster username satyanadella
```

```text
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
 URL: https://www.linkedin.com/posts/satyanadella_...

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

### `domain`

Performs subdomain enumeration, DNS resolution, SSL/TLS certificate inspection, and subdomain takeover detection.

#### Syntax

```bash
./osintmaster domain <domain> [wordlist.txt] [fingerprints.json] [-o <FileName>]
```

#### Arguments & Options

| Argument / Option     | Required | Description                                                                       |
| --------------------- | -------- | --------------------------------------------------------------------------------- |
| `<domain>`            | Yes      | Target domain to scan (e.g. `example.com`).                                       |
| `[wordlist.txt]`      | No       | Custom wordlist for subdomain discovery _(default: `src/assets/subdomains.txt`)_. |
| `[fingerprints.json]` | No       | Custom takeover fingerprints file _(default: `src/assets/fingerprints.json`)_.    |

#### Examples

```bash
# Scan a domain using default wordlist and fingerprints
./osintmaster domain example.com

# Scan with a custom wordlist
./osintmaster domain example.com custom_wordlist.txt

# Scan with custom wordlist and custom fingerprints
./osintmaster domain example.com custom_wordlist.txt custom_fingerprints.json

# Save output to a file
./osintmaster domain example.com -o domain_scan.txt
```

#### Example Output

```text
$ ./osintmaster domain example.com
example.com [172.66.147.243 104.20.23.154 2606:4700:10::ac42:93f3 2606:4700:10::6814:179a]
2026-09-26 -> 2026-12-25 / 50 days left
Subject: CN=example.com
Issuer: CN=Cloudflare TLS Issuing ECC CA 3,O=SSL Corporation,C=US
DNS Names: [example.com *.example.com]

[POTENTIAL TAKEOVER] blog.example.com -> custom.github.io (Service: github)
blog.example.com [185.199.108.153 185.199.109.153]
2026-01-10 -> 2026-04-10 / 65 days left
Subject: CN=custom.github.io
Issuer: CN=DigiCert Global G2 TLS RSA SHA256 2020 CA1,O=DigiCert Inc,C=US
DNS Names: [custom.github.io]
```

---

### `ip`

Retrieves geolocation data and ISP details using `ip-api.com` and checks for historical data associated with the IP using `otx.alienvault.com/api`.

#### Limitations

`ip-api.com` is limited to 45 requests per minute.
`otx.alienvault.com/api` is limited to 10–20 requests per minute.

#### Syntax

```bash
./osintmaster ip <ip_address> [-o <FileName>]
```

#### Arguments & Options

| Argument / Option | Required | Description                       |
| ----------------- | -------- | --------------------------------- |
| `<ip_address>`    | Yes      | Target IP address (IPv4 or IPv6). |

#### Examples

```bash
# Lookup an IP address
./osintmaster ip 8.8.8.8

# Save output to a file
./osintmaster ip 8.8.8.8 -o ip_scan.txt
```

#### Example Output

```text
$ ./osintmaster ip 8.8.8.8
  - Geolocation -
Country: United States
Region: Virginia
City: Ashburn
Latitude: 39.03
Longitude: -77.5
Timezone: America/New_York

  - ISP Details -
ISP: Google LLC
Organization: Google Public DNS
AS: AS15169 Google LLC

  - Historical Abuse Data -
Historical Reports: 0
```

---

## Disclaimer & Ethical Use

> [WARNING]
> This tool is developed strictly for educational purposes, authorized security testing, and lawful digital investigation. Users are solely responsible for ensuring compliance with all applicable local, national, and international laws. The authors assume no liability for misuse or damage caused by this software.
