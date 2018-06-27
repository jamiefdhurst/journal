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
brackets. URLs are parameterised to include post slugs, as opposed to IDs. 

## Available Endpoints

### Retrieve all posts

**Method/URL:** `GET /api/v1/post`

**Successful Response:** `200`

Contains all current post reources in reverse date order.

```json
[
    {
        "id": 1,
        "slug": "example-post",
        "title": "An Example Post",
        "date": "2018-05-18T12:53:22Z",
        "content": "<p>TEST</p><p>:gif:id:cE1qRt8nl6Neo:</p>"
    }
]
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
    "id": 1,
    "slug": "example-post",
    "title": "An Example Post",
    "date": "2018-05-18T12:53:22Z",
    "content": "<p>TEST</p><p>:gif:id:cE1qRt8nl6Neo:</p>"
}
```
**Error Responses:** 

`404` - Post with provided slug could not be found.

--

### Create a post

**Method/URL:** `PUT /api/v1/post`

Post is provided as JSON, ommitting the ID and slug:

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
    "id": 2,
    "slug": "a-brand-new-post",
    "title": "A Brand New Post",
    "date": "2018-06-28T00:42:12Z",
    "content": "<p>This is a brand new post, completely.</p>"
}
```

**Error Responses:**

* `400` - Incorrect parameters supplied - the date, title and content must be
provided.

--

### Update a post

**Method/URL:** `POST /api/v1/post/{slug}`

The post's slug is used as the individual URI for the resource.

E.g.: `/api/v1/post/example-post`

Keys to update within the post can be one or more of `date`, `title` and
`content`:

```json
{
    "content": "<p>I'm only changing the content this time.</p>"
}
```

Or:

```json
{
    "date": "2018-06-21T09:12:00Z",
    "title": "Even Braver New World",
    "content": "<p>I changed a bit more on this attempt.</p>"
}
```

When updating the post, the slug remains constant, even when the title changes.

**Successful Response:** `200`

```json
{
    "id": 2,
    "slug": "a-brand-new-post",
    "title": "Even Braver New World",
    "date": "2018-06-21T09:12:00Z",
    "content": "<p>I changed a bit more on this attempt.</p>"
}
```

**Error Responses:**

* `400` - Incorrect parameters supplied - at least one or more of the date,
title and content must be provided.
* `404` - Post with provided slug could not be found.
