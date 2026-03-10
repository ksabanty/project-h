# Reddit Sports Highlights Aggregator

A Go application that fetches and displays recent sports highlight videos from multiple subreddits using the Reddit API. The application automatically filters for video posts from the past 24 hours and sorts them by popularity.

## Table of Contents

- [Features](#features)
- [Prerequisites](#prerequisites)
- [Reddit API Setup](#reddit-api-setup)
- [Installation](#installation)
- [Configuration](#configuration)
  - [Environment Variables](#environment-variables)
  - [Subreddit Queries](#subreddit-queries)
- [Building and Running](#building-and-running)
- [Customization](#customization)
- [Troubleshooting](#troubleshooting)

## Features

- Fetches recent video posts from multiple subreddits
- Filters posts by flair (e.g., "Highlight", "Media")
- Only displays posts from the last 24 hours
- Filters for video content from popular platforms (v.redd.it, streamable.com, YouTube, Twitter, etc.)
- Sorts results by upvotes
- Outputs clickable terminal links for easy access
- Caches Reddit OAuth tokens to minimize API calls

## Prerequisites

Before building and running this application, ensure you have the following installed:

- **Go** version 1.23.0 or higher ([download here](https://golang.org/dl/))
- A **Reddit account** to create an API application
- Terminal with clickable link support (most modern terminals)

## Reddit API Setup

To use the Reddit API, you must create a Reddit application and obtain authentication credentials.

### Step 1: Create a Reddit Application

1. Log in to your Reddit account
2. Navigate to [https://www.reddit.com/prefs/apps](https://www.reddit.com/prefs/apps)
3. Scroll to the bottom and click **"create another app..."** or **"are you a developer? create an app..."**
4. Fill out the form:
   - **name**: Choose any name (e.g., "Sports Highlights Fetcher")
   - **App type**: Select **"script"**
   - **description**: Optional
   - **about url**: Optional
   - **redirect uri**: Enter `http://localhost:8080` (required but not used)
5. Click **"create app"**

### Step 2: Retrieve Your Credentials

After creating the app, you will see your application details:

- **CLIENT_ID**: The string directly under "personal use script" (14 characters)
- **CLIENT_SECRET**: The string labeled "secret" (27 characters)

Keep these credentials secure and do not share them publicly.

## Installation

### Step 1: Clone or Download the Repository

```bash
git clone <repository-url>
cd project-h
```

Or download and extract the source code to a directory of your choice.

### Step 2: Install Dependencies

The application uses Go modules for dependency management. Install all required dependencies:

```bash
go mod download
```

## Configuration

### Environment Variables

Create a `.env` file in the project root directory with your Reddit API credentials:

```bash
CLIENT_ID=your_client_id_here
CLIENT_SECRET=your_client_secret_here
USER_AGENT=YourApp/1.0 by YourRedditUsername
```

**Field Descriptions:**

- **CLIENT_ID**: Your Reddit app's client ID (from Reddit API Setup)
- **CLIENT_SECRET**: Your Reddit app's secret key (from Reddit API Setup)
- **USER_AGENT**: A unique identifier for your app. Reddit requires this to follow their API rules. Format: `AppName/Version by RedditUsername`

**Example:**

```bash
CLIENT_ID=ABc123XYz45678
CLIENT_SECRET=aBcDeFgHiJkLmNoPqRsTuVwXyZ1
USER_AGENT=SportsHighlightApp/1.0 by u/your_username
```

### Subreddit Queries

The application reads subreddit and flair configurations from `assets/subreddit_queries.json`. This file determines which subreddits to query and which post flairs to filter.

**File Location:** `assets/subreddit_queries.json`

**Format:**

```json
[
    {
        "subreddit": "subreddit_name",
        "search_query": "flair_name:FlairName"
    }
]
```

**Field Descriptions:**

- **subreddit**: The name of the subreddit (without the "r/" prefix)
- **search_query**: The flair filter in the format `flair_name:FlairName`. Leave empty (`""`) to search all flairs in that subreddit

**Default Configuration:**

```json
[
    {
        "subreddit": "soccer",
        "search_query": "flair_name:Media"
    },
    {
        "subreddit": "nfl",
        "search_query": "flair_name:Highlight"
    },
    {
        "subreddit": "nba",
        "search_query": "flair_name:Highlight"
    },
    {
        "subreddit": "baseball",
        "search_query": ""
    }
]
```

**To Add More Subreddits:**

1. Identify the subreddit you want to add
2. Visit the subreddit and identify the flair used for highlight posts (usually visible on post listings)
3. Add a new entry to the JSON array:

```json
{
    "subreddit": "hockey",
    "search_query": "flair_name:Streamable"
}
```

**Special Behavior:**

- For **r/baseball**, the application additionally filters for posts containing `[Highlight]` in the title
- Empty `search_query` will search all posts regardless of flair

## Building and Running

### Option 1: Run Directly with Go

Execute the application without building a binary:

```bash
go run main.go
```

### Option 2: Build and Run Executable

Build a standalone executable:

```bash
go build -o highlights main.go
```

Run the compiled executable:

```bash
./highlights
```

### Expected Output

When run successfully, the application will:

1. Load authentication credentials from `.env`
2. Obtain or use cached Reddit API access token
3. Fetch posts from each configured subreddit
4. Display progress: `Fetching posts from r/soccer`
5. Output clickable links to highlight posts, sorted by popularity

**Example Output:**

```
Fetching posts from r/soccer
Fetching posts from r/nfl
Fetching posts from r/nba
Fetching posts from r/baseball
Orange County SC's Vuk Latinovich with a fantastic own goal
Patrick Mahomes throws 45-yard TD to Travis Kelce
LeBron James with the monster dunk
Aaron Judge crushes 450-foot home run
```

Links are clickable in supported terminals (click to open in browser).

## Customization

### Change Time Window

By default, the application fetches posts from the last 24 hours. To modify this:

1. Open `main.go`
2. Locate the line: `oneDayAgo := now - 86400`
3. Change `86400` to your desired number of seconds:
   - 12 hours: `43200`
   - 48 hours: `172800`
   - 1 week: `604800`

### Modify Video Domain Filters

The application filters for posts from specific video domains. To add or remove domains:

1. Open `main.go`
2. Locate the `isVideoPost` condition (around line 127)
3. Add or remove domain checks following the existing pattern:

```go
strings.Contains(child.Data.Domain, "yourdomain.com") ||
```

### Change Number of Posts Fetched

To fetch more or fewer posts per subreddit:

1. Open `main.go`
2. Locate: `limit=50` in the API URL
3. Change `50` to your desired limit (maximum: 100 per Reddit API)

## Troubleshooting

### Error: "Error loading .env file"

**Solution:** Ensure a `.env` file exists in the project root directory with the required variables.

### Error: "Error getting response" or HTTP 401

**Solution:** Verify your `CLIENT_ID`, `CLIENT_SECRET`, and `USER_AGENT` are correct in the `.env` file.

### No Results Displayed

**Possible Causes:**
- No posts matching criteria in the last 24 hours
- Incorrect flair names in `subreddit_queries.json`
- Subreddit has no posts with the specified flair

**Solution:** Check the subreddit manually to verify flair names and recent activity.

### Error: "Error opening queries file"

**Solution:** Ensure `assets/subreddit_queries.json` exists and is properly formatted JSON.

### Links Not Clickable

**Solution:** Use a modern terminal that supports OSC 8 hyperlinks (iTerm2, GNOME Terminal 3.35+, Windows Terminal, etc.).

### Token Cache Issues

The application caches tokens in `token_cache.json`. If experiencing authentication issues:

```bash
rm token_cache.json
```

Then run the application again to fetch a fresh token.

## License

This project is provided as-is for personal use.
