# User Guide

This guide covers the day-to-day use of Journal: browsing entries, creating
new posts, editing existing ones, and navigating the built-in views.

---

## The index

When you open Journal in your browser you land on the index page. It lists
your most recent entries in reverse-chronological order, showing a short
excerpt from each one. Use the pagination links at the bottom of the page to
move between pages; the number of entries per page is controlled by the
`J_POSTS_PER_PAGE` configuration variable (default: 20).

Each entry on the index shows:

- **Title** — click it to open the full entry
- **Date** — the date the entry is written for
- **Excerpt** — the first few words of the content (default: 50 words)
- **Edit** button — only visible when editing is enabled

---

## Reading an entry

Click an entry title or its **Read More** link to open the full post. The
content is written in Markdown and rendered as HTML. At the bottom of the
entry you will find links to navigate to the **previous** and **next**
entries.

Each entry URL follows the pattern `/<slug>`, where the slug is derived
automatically from the title when the post is created (for example, a title
of "My First Post" becomes `/my-first-post`).

---

## Creating an entry

> Creating entries must be enabled. If the **New Entry** link is not
> visible, the `J_CREATE` variable has been set to `0`.

1. Click **New Entry** in the navigation (or go to `/new` directly).
2. Fill in the three fields:
   - **Title** — required, at least 3 characters.
   - **Date** — required, in `YYYY-MM-DD` format. The date picker pre-fills
     today's date.
   - **Content** — required, at least 3 characters. Write in
     [Markdown](https://www.markdownguide.org/basic-syntax/) — headings,
     bold, italic, lists, links, code blocks and blockquotes are all
     supported.
3. Click **Save**.

On success you are redirected to the index, which shows a confirmation
message. The new entry appears at the top of the list.

### How slugs are generated

The URL slug is derived from the title automatically:

- Characters that are not letters or numbers are replaced with `-`
- Everything is lowercased
- If another entry already has the same slug, a numeric suffix is appended
  (e.g. `-2`, `-3`)

You cannot set the slug manually through the UI.

### Validation rules

The form will reject a submission if:

- The title is missing or fewer than 3 characters
- The date is missing or not in `YYYY-MM-DD` format
- The content is missing or fewer than 3 characters
- The slug derived from the title conflicts with a reserved path (`new`,
  `random`, `stats`) or the API prefix (`api/`)

---

## Editing an entry

> Editing entries must be enabled. If the **Edit** button is not visible,
> the `J_EDIT` variable has been set to `0`.

1. From the index or the entry view, click **Edit**.
2. Modify the **Title**, **Date**, or **Content** as needed.
3. Click **Save**.

The same validation rules apply as when creating. Saving an edit updates the
entry's last-modified timestamp but does not change its URL slug.

---

## The calendar

Navigate to `/calendar` to see a month-by-month calendar view of your
entries. Days that have at least one entry show the entry titles as links.

Use the navigation arrows to move between months and years:

- `/calendar` — current month
- `/calendar/2024` — all months in 2024
- `/calendar/2024/03` — March 2024

---

## Random entry

Go to `/random` (or click **Random** in the navigation if your theme
includes it) to jump to a randomly selected entry. This is useful for
rediscovering older writing.

---

## Stats

The `/stats` page shows an overview of your journal:

- **Total posts** and the date of the first post
- Current **configuration** (title, description, theme, posts per page,
  whether create/edit are enabled)
- **Daily visits** for the last 14 days, broken down by web and API hits
- **Monthly visits** for the full history

---

## Markdown reference

Entry content is written in Markdown. A quick reference:

| Format | Markdown |
|---|---|
| Heading | `# H1`, `## H2`, `### H3` |
| Bold | `**bold**` |
| Italic | `*italic*` |
| Link | `[text](https://example.com)` |
| Unordered list | `- item` |
| Ordered list | `1. item` |
| Inline code | `` `code` `` |
| Code block | ```` ```language ```` … ```` ``` ```` |
| Blockquote | `> quote` |
| Horizontal rule | `---` |
