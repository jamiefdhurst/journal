# Journal API Documentation - v1

## Introduction

The API provides full access to posts within the journal - viewing, creating
and editing.

### HTTP Requests

* `GET` - Retrieve one or more resources
* `PUT` - Create a new resource.
* `POST` - Update a resource.

### Response Codes

* `200` - Request was successful.
* `400` - Request was not understood, or required parameters were missing.
* `404` - Resource was not found.

### URL Parameters

When specified within endpoints, URL parameters are shown within `{}` curly
brackets. URLs are parametrised to include post slugs, as opposed to IDs.

## Available Endpoints

### Retrieve all posts

**Method/URL:** `GET /api/v1/post`

**Successful Response:** `200`

Contains all current post resources in reverse date order, paginated. The
`links` property containers next and previous links, and `pagination` contains
information on the total posts, pages and posts per page.

```json
{
    "links": {
        "prev": "/api/v1/post?page=1",
        "next": "/api/v1/post?page=3"
    },
    "pagination": {
        "current_page": 2,
        "total_pages": 3,
        "posts_per_page": 1,
        "total_posts": 3
    },
    "posts": [
        {
            "url": "/api/v1/post/example-post",
            "title": "An Example Post",
            "date": "2018-05-18T00:00:00Z",
            "content": "TEST",
            "created_at": "2018-05-18T15:16:17Z",
            "updated_at": "2018-05-18T15:16:17Z"
        }
    ]
}
```

**Error Responses:** *None*

--

### Retrieve a single post

**Method/URL:** `GET /api/v1/post/{slug}`

The post's slug is used as the individual URI for the resource.

E.g.: `/api/v1/post/example-post`

**Successful Response:** `200`

Contains the single post.

```json
{
    "url": "/api/v1/post/example-post",
    "title": "An Example Post",
    "date": "2018-05-18T00:00:00Z",
    "content": "TEST",
    "created_at": "2018-05-18T15:16:17Z",
    "updated_at": "2018-05-18T15:16:17Z"
}
```

**Error Responses:**

`404` - Post with provided slug could not be found.

--

### Create a post

**Method/URL:** `PUT /api/v1/post`

Post is provided as JSON, omitting the ID and slug:

```json
{
    "title": "A Brand New Post",
    "date": "2018-06-28T00:42:12Z",
    "content": "<p>This is a brand new post, completely.</p>"
}
```

The date can be provided in the following formats:

* `2018-06-28`
* `2018-06-28T00:42:12Z`

**Successful Response:** `200`

```json
{
    "url": "/api/v1/post/a-brand-new-post",
    "title": "A Brand New Post",
    "date": "2018-06-28T00:42:12Z",
    "content": "This is a brand new post, completely."
}
```

**Error Responses:**

* `400` - Incorrect parameters supplied - the date, title and content must be
provided.

--

### Retrieve a random post

**Method/URL:** `GET /api/v1/post/random`

**Successful Response:** `200`

Contains a randomly selected post.

```json
{
    "url": "/api/v1/post/example-post",
    "title": "An Example Post",
    "date": "2018-05-18T12:53:22Z",
    "content": "TEST"
}
```

**Error Responses:**

`404` - No posts exist in the system.

--

### Update a post

**Method/URL:** `POST /api/v1/post/{slug}`

The post's slug is used as the individual URI for the resource.

E.g.: `/api/v1/post/example-post`

Keys to update within the post can be one or more of `date`, `title` and
`content`:

```json
{
    "content": "I'm only changing the content this time."
}
```

Or:

```json
{
    "date": "2018-06-21T09:12:00Z",
    "title": "Even Braver New World",
    "content": "I changed a bit more on this attempt."
}
```

When updating the post, the slug remains constant, even when the title changes.

**Successful Response:** `200`

```json
{
    "url": "/api/v1/post/a-brand-new-post",
    "title": "Even Braver New World",
    "date": "2018-06-21T09:12:00Z",
    "content": "I changed a bit more on this attempt."
}
```

**Error Responses:**

* `400` - Incorrect parameters supplied - at least one or more of the date,
title and content must be provided.
* `404` - Post with provided slug could not be found.

---

### Stats

**Method/URL:** `GET /api/v1/stats`

**Successful Response:** `200`

Retrieve statistics, configuration information and visit summaries for the
current installation.

```json
{
    "posts": {
        "count": 3,
        "first_post_date": "Monday January 1, 2018"
    },
    "configuration": {
        "title": "Jamie's Journal",
        "description": "A private journal containing Jamie's innermost thoughts",
        "theme": "default",
        "posts_per_page": 20,
        "google_analytics": false,
        "create_enabled": true,
        "edit_enabled": true
    },
    "visits": {
        "daily": [
            {
                "date": "2025-01-01",
                "api_hits": 20,
                "web_hits": 30,
                "total": 50
            }
        ],
        "monthly": [
            {
                "month": "2025-01",
                "api_hits": 200,
                "web_hits": 300,
                "total": 500
            }
        ]
    }
}
```

**Error Responses:** *None*
