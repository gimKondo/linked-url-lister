# Feature Specification: URL Crawler for NotebookLM Import

**Feature Branch**: `001-url-crawler`
**Created**: 2025-12-10
**Status**: Draft
**Input**: Web crawler to list URLs under a specified website for NotebookLM import

## User Scenarios & Testing *(mandatory)*

### User Story 1 - Basic URL Crawling (Priority: P1)

A user wants to extract all relevant page URLs from a documentation website (e.g., GitLab Handbook) to import into NotebookLM for analysis. The user provides a starting URL and receives a list of all discoverable pages that contain meaningful text content.

**Why this priority**: This is the core functionality of the tool. Without basic crawling, the tool has no value.

**Independent Test**: Can be fully tested by running the tool against a small test website with known link structure and verifying the output contains expected URLs.

**Acceptance Scenarios**:

1. **Given** a valid URL to a documentation site, **When** the user runs the tool with that URL, **Then** the tool outputs a list of URLs that are linked from the starting page
2. **Given** a URL with nested pages (e.g., /docs/guide/chapter1), **When** the tool crawls, **Then** it discovers and lists pages in subdirectories of the starting URL
3. **Given** pages with varying amounts of text, **When** the tool filters results, **Then** only pages with substantial text content (default: 500+ characters) are included

---

### User Story 2 - Output for NotebookLM (Priority: P2)

A user wants to copy the URL list directly into NotebookLM's import interface. The output format must be compatible with NotebookLM's URL input requirements.

**Why this priority**: Output formatting is essential for the stated purpose but depends on basic crawling working first.

**Independent Test**: Can be tested by verifying output format matches NotebookLM's expected input (one URL per line, valid absolute URLs).

**Acceptance Scenarios**:

1. **Given** a completed crawl, **When** the user views the output, **Then** each URL appears on its own line
2. **Given** a completed crawl, **When** the output is pasted into NotebookLM, **Then** NotebookLM accepts all URLs without format errors
3. **Given** a large site with many pages, **When** the user requests JSON output format, **Then** the tool outputs URLs in a structured JSON array

---

### User Story 3 - Crawl Control (Priority: P3)

A user wants to control crawl behavior to limit scope, avoid overloading servers, or focus on specific sections of a site.

**Why this priority**: Control features improve usability but are not essential for basic operation.

**Independent Test**: Can be tested by running crawls with different options and verifying behavior changes accordingly.

**Acceptance Scenarios**:

1. **Given** a maximum depth option, **When** the user specifies --max-depth=2, **Then** the tool only follows links up to 2 levels deep from the starting URL
2. **Given** a minimum text length option, **When** the user specifies --min-text=1000, **Then** only pages with 1000+ characters of text are included
3. **Given** a rate limit concern, **When** the tool crawls, **Then** it respects a configurable delay between requests (default: 100ms)

---

### Edge Cases

- What happens when the starting URL returns a 404 or other error? → Tool reports error and exits with non-zero status
- What happens when a page requires authentication? → Tool skips the page and continues crawling accessible pages
- What happens when the site has circular links? → Tool tracks visited URLs and avoids revisiting
- What happens when the site has thousands of pages? → Tool provides progress output to stderr; user can set max-pages limit
- What happens when robots.txt disallows crawling? → Tool respects robots.txt by default (can be overridden with flag)

## Requirements *(mandatory)*

### Functional Requirements

- **FR-001**: System MUST accept a starting URL as a command-line argument
- **FR-002**: System MUST discover URLs by following links from the starting page
- **FR-003**: System MUST recursively crawl pages that are under the same URL path prefix as the starting URL
- **FR-004**: System MUST filter out pages that require authentication (non-200 responses, login redirects)
- **FR-005**: System MUST filter out pages with insufficient text content (configurable threshold, default 500 characters)
- **FR-006**: System MUST output discovered URLs to stdout (one per line by default)
- **FR-007**: System MUST track visited URLs to avoid infinite loops from circular links
- **FR-008**: System MUST respect robots.txt directives by default
- **FR-009**: System MUST support JSON output format via --json flag
- **FR-010**: System MUST support configurable maximum crawl depth via --max-depth flag
- **FR-011**: System MUST support configurable minimum text length via --min-text flag
- **FR-012**: System MUST support configurable request delay via --delay flag
- **FR-013**: System MUST output progress information to stderr during crawl
- **FR-014**: System MUST handle network errors gracefully and continue crawling other pages

### Key Entities

- **URL**: A web page address to be crawled or output; attributes include absolute URL string, crawl depth, and discovery status
- **Page**: A fetched web resource; attributes include URL, HTTP status, text content length, and list of extracted links
- **CrawlResult**: The final output; a collection of URLs that passed all filters (path prefix, authentication, text content)

## Success Criteria *(mandatory)*

### Measurable Outcomes

- **SC-001**: Users can obtain a complete URL list from a 100-page documentation site in under 5 minutes
- **SC-002**: 100% of output URLs are accessible without authentication when manually verified
- **SC-003**: 95% of output URLs contain the minimum specified text content when manually verified
- **SC-004**: Output can be directly pasted into NotebookLM without manual reformatting
- **SC-005**: Tool completes crawl of GitLab Handbook (https://handbook.gitlab.com/) top-level sections successfully

## Assumptions

- Users have network access to the target websites
- Target websites serve HTML content (not JavaScript-only SPAs that require browser rendering)
- NotebookLM accepts plain text URL lists (one URL per line) as input
- A reasonable default for "substantial text content" is 500 characters of visible text
- Default request delay of 100ms provides reasonable balance between speed and server courtesy
