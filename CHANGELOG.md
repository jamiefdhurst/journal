# Changelog

All notable changes to this project will be documented in this file.

The format is based on [Keep a Changelog](https://keepachangelog.com/en/1.0.0/),
and this project adheres to [Semantic Versioning](https://semver.org/spec/v2.0.0.html).

## [1.0.1] - 2026-05-01

### Fixed


- fix: migration order was incorrect, and added screenshot to README (a6e7c92)


- fix: migration order was incorrect, and added screenshot to README (#153) (605b297)


## [1.0.0] - 2026-04-29

### Added


- Add re-brand milestone via settings (82dbddd)


- Add branch protection to release branches (b5ce137)


### Breaking Changes


- breaking: complete v1.0.0 overhaul with new features and removed legacy support (244426b)

**Breaking changes:**
  - Remove Giphy integration — API key environment variable no longer supported
  - Remove AWS Lambda deployment support — returns to standard HTTP server model only
  - Rename J_DATA_PATH environment variable to J_WEB_PATH
  - Restructure entry point from journal.go to cmd/journal/main.go
  - Replace frontend build toolchain (gulp/sass/webpack) with pre-built CSS theme system
  - Move static assets from web/static/ to web/themes/default/

  **New features:**
  - Calendar view for browsing articles by month and year
  - Stats page and API endpoint showing article counts, word counts, and visit data
  - Random article endpoint available via web UI and API
  - Full Markdown rendering support for article content
  - Database migration system for zero-downtime schema changes
  - Visit and request tracking stored in the database

  **API and documentation:**
  - OpenAPI 3.0 specification added at api/openapi.yml
  - Installation and user guide documentation added under docs/

  **Build and release:**
  - Multi-platform binary builds for linux/darwin on amd64/arm64
  - apt/deb and yum/rpm packaging via nfpm
  - Multi-platform Docker images pushed to GHCR and Docker Hub
  - Automated semantic versioning from conventional commits
  - Homebrew tap support for macOS installs
  - Remove CGO and SQLite dependency from Docker build

  Test coverage raised to ~100% across all application packages.


### Other


- Rename milestone (a757097)


- Release wildcard verification will not work (3c68a13)


## [0.9.5] - 2024-09-02

### Changed


- Update deploy workflow to use DO server (b0bb69d)


## [0.9.4] - 2024-08-28

### Added


- Add create for database where it doesn't exist (47e2660)


- Add create for database where it doesn't exist (#97) (cd787f6)


### Changed


- Update JSON test to use a free alternative (18a9302)


## [0.9.3] - 2024-05-20

### Other


- Disallow everything (a148b6f)


## [0.9.2] - 2024-05-14

### Other


- Remove environment URL from environment (f1e021c)


## [0.9.1] - 2024-05-14

### Added


- Add deployment ability for journal to Lambda (3e7cf91)


## [0.9.0] - 2024-05-14

### Changed


- Update build files and app to support Lambda execution (61373aa)


- Update build files and app to support Lambda execution (#95) (1e6acb1)


## [0.8.5] - 2023-12-22

### Other


- Correct package.json - was full of journal.go (a3bfbb0)


## [0.8.4] - 2023-12-21

### Changed


- Update Dockerfile to link package to repository (4df6923)


### Other


- LABEL must be after FROM (64d19fe)


## [0.8.3] - 2023-12-21

### Changed


- Update README to show Actions and to remove Jenkins badge (502525a)


## [0.8.2] - 2023-12-21

### Fixed


- Fix version in journal.go (6485de5)


### Other


- Still not quite working... Trying a split (c808251)


## [0.8.1] - 2023-12-21

### Other


- Make sure the version changes are saving to all files (1006150)


## [0.8.0] - 2023-12-21

### Testing


- Test a vrsion bump (2e3e95f)


## [0.7.4] - 2023-12-21

### Changed


- Update to fix the major/minor/patch detection (10e88fd)


## [0.7.3] - 2023-12-21

### Changed


- Update to remove release from release name (f42e858)


## [0.7.2] - 2023-12-21

### Added


- Add initial test action (a0254f9)


-  (8a99c39)

Add Actions for Testing and Building


### Changed


- Update to use merge commits, remove milestones from settings (59c007f)


- Update paths and cache location (22dedb0)


- Switch out coverage tool (425b3a3)


- Update to use PAT to get around protected branches issue (8953664)


### Other


- Skip CI: updated version number (9954fb4)


- Correct permissions (bd632f4)


- GOPATH support to try and fix tests (bdb0745)


- Run commands from correct locations (b12c103)


- Publish code coverage report (5bec03a)


- Accidentally removed go2xunit (494ff1c)


- Missed version in replace action (56a532d)


- Build is only going to work on main, going to need to add it there (245ada5)


- Surround main in quotes (2be4197)


- Forgot to remove if condition for workflow: (36f4893)


- Try a different calculation action (bff4fbb)


- Debug what's happening (52d58fa)


- More attempts (dc88982)


- Skip CI properly and checkout with PAT (8dbc74a)


- No force admins (0421f97)


### Testing


- Testing build action under normal branch (b0fb0f9)


## [0.7.1.1] - 2023-10-28

### Other


- Skip CI: updated version number (40bf852)


## [0.7.1] - 2023-09-21

### Fixed


- Fix edit button placement (3d08856)


### Other


- Skip CI: updated version number (a24eb64)


## [0.7] - 2023-06-27

### Added


- Add correct HTML titles and descriptions (c99c4a0)


- Add Google blog structure to articles (b544342)


- Add robots.txt (db62a75)


- Add robots.txt and sitemap.xml (d4ce947)


- Add GA support (a5fa7b7)


- Added caching to static assets and versions to JS/CSS (fcd3064)


- Add favicons (d7b8911)


- Add humans.txt (f17eb89)


### Changed


- Switch up H-tags correctly (0e6fad0)


### Fixed


- Fix a bleed into button (fc53db5)


### Other


- Remove header H1 (9212765)


- Replace JSON placeholder test to use own Journal API (4863fa7)


- Skip CI: updated version number (f85d15b)


## [0.6.1.1] - 2023-04-20

### Other


- Coverage check improvements (545a620)


- Skip CI: updated version number (f129625)


## [0.6.1] - 2023-04-20

### Other


- FIx mobile styling to be more consistent (af93f6c)


- Skip CI: updated version number (bedb7d1)


## [0.6] - 2023-04-20

### Added


- Add session support and integrate (05d6aab)


### Fixed


- Fix mock response call (644f73d)


- Fix tests for introducing banners, and so old issues (9580045)


### Other


- Skip CI: updated version number (b4d4a0f)


## [0.5.0.1] - 2023-03-16

### Other


- README link to CI was wrong (1a5069c)


- Skip CI: updated version number (ea1fdcb)


## [0.5] - 2023-03-16

### Added


- Add GitHub link (a1bae70)


### Other


- Skip CI: updated version number (0a51e90)


## [0.4.3.2] - 2023-03-16

### Other


- Skip CI: updated version number (8b1599c)


## [0.4.3.1] - 2023-03-15

### Other


- Unnecessary scm check (1024b0f)


- Skip CI: updated version number (5724886)


## [0.4.3] - 2023-03-14

### Fixed


- Fix skip CI for now (dd43246)


### Other


- Skip CI: updated version number (7fac6c1)


## [0.4.2.2] - 2023-03-14

### Other


- Remove re-build of image, this is handled in release process (138e946)


- Skip CI: updated version number (0cc5230)


## [0.4.2.1] - 2023-03-14

### Other


- Reconfigure build to add docker image and deployment: (28a70d1)


- Skip CI: updated version number (c182d6a)


## [0.4.2] - 2023-03-12

### Fixed


- Fix to default to major go version (b76d21d)


### Other


- Skip CI: updated version number (f0fe7f4)


## [0.4.1] - 2023-03-05

### Other


- Rework build to support automatic versioning (d09bb58)


- Skip CI: updated version number (71963f9)


## [0.4] - 2023-03-05

### Added


- Add GitHub settings (20f96a7)


### Changed


- Update tests to download tooling correctly (125242f)


### Fixed


- Fix readme url to go-sqlite3 (#43) (d5fd6c3)


- Fix for glob-parent vulnerabilities with overrides (07f179d)


### Other


- Rework pipeline to use newer syntax (5de6e5e)


- Move projects and add re-brand milestone (d5f241b)


- Ensure sass still works (09994a6)


- Ensure version numbers are aligned (e69ef57)


- Skip CI: updated version number (e36d316)


## [0.3.0.3] - 2021-10-24

### Changed


- Update test Dockerfile too (5910820)


### Fixed


- Fix Dockerfile to use sqlite3 (a639889)


## [0.3.0.2] - 2021-03-27

### Added


- Add purpose to readme (7a2f6e6)


- Add missing test dockerfile (19d8db4)


### Changed


- Update build to provide unit tests and coverage (3ddee77)


### Other


- Not bothering with older go versions (99fa993)


## [0.3.0.1] - 2021-03-24

### Fixed


- Fix style issue, increase version (08936c1)


## [0.3] - 2021-03-21

### Added


- Add Jenkinsfile (8024c27)


- Add go.mod now this is a requirement (948ad6c)


- Add in new template, styles and fixes for showing better dates (4e632bb)


- Add quiet pagination, 20 articles per page (31d361c)


- Adding previous and next links to view pages (920a099)


- Add title and articles to config (ef30116)


- Add DB path parameter (6f32563)


- Add missing test coverage in edit/create API calls (1929f7e)


- Add excerpt functionality and better design work (ad35381)


- Added some documentation (b19cffa)


### Changed


- Update build badge (7de87af)


- Update Jenkins URLs (ffcd34a)


- Tidy Jenkinsfile (02d8958)


- Update README for new env variables (8bdf5ea)


### Fixed


- Fix some readme derp, and add in some more badges (d99e36b)


- Fix build status badge in README (28a20ca)


- Fix update test (0ee7375)


### Other


- Simplify build pipeline (2914d1b)


- Setup for configuration (b95ca57)


- Support port variable (f891202)


- Enable blocking editing and creating articles (f554bce)


## [0.2.2] - 2019-05-02

### Added


- Add Dockerfile and travis updates (a05a128)


-  (9e39f27)

Add Dockerfile and travis updates


- Add test to ensure duplicate titles are handled correctly (f99c1b1)


### Fixed


- Fix build issue, and remove go dependency to speed up build (bf83972)


- Fix status code (2fef2b0)


### Other


- Run all tests on build, not just API tests (8a04ac6)


- Increase version to v0.2.2 (b0d0d02)


- Automatically create database, and support missing GIPHY functionality (2f573ba)


- One test failing (a23ce86)


-  (4777ee4)

Better Handling of Missing Dependencies


- Merge branch 'master' of github.com:jamiefdhurst/journal (d5ab696)


## [0.2.1] - 2018-08-03

### Added


- Add travis CI (3126b84)


- Add build status (b30cb11)


### Fixed


- Fix title slugifying incorrectly (5cfd6b3)


- Fix API tests (046ec9c)


- Fix saved banner consistently showing (9e8f4d8)


- Fix create tests (5c5983e)


### Other


- Remove .DS_Store (7d29f6a)


- Increase version to 0.2.1 (9931f68)


- Prevent multiple slugs (baa92bf)


-  (607bc20)

v0.2.1


## [0.2] - 2018-06-27

### Added


- Add list and single endpoints (9968952)


- Add POST request for create (f382c5b)


- Add controller test (d462e33)


- Add all web controller tests (b8fea6e)


- Add final unit tests (439b52a)


- Add documentation (31a71ab)


- Add link to API docs (8e13f2e)


### Changed


- Update viewws and simplify router (9b46eb7)


-  (583e024)

Refactor


### Documentation


-  (24ad5a7)

Documentation


### Fixed


- Fix for #1 - textarea height (338ae81)


-  (5124053)

Fix for #1 - textarea height


- Fix views and paths again (e17b64b)


- Fix package issues identified in nsp (853d6a5)


-  (676e47d)

Fix package issues identified in nsp


- Fix (e645dd8)


### Other


- Better new post location (fda109b)


-  (b34dc81)

Better new post location


- Remove stale plugin files to tidy up CSS (2aded73)


-  (147802c)

Remove stale plugin files to tidy up CSS


-  (a4e1692)

First stage of RESTful API


- Editing in layout all works (df69a89)


- API edit endpoint ready (c5f38c0)


-  (3c93d22)

Edit functionality


- API key can be stored in DB (600a001)


- GIF support working (ade4897)


-  (31723a2)

GIF integration


- Move folders and rename server to app, to better understand purpose (a65509c)


- Finished refactor (2efe134)


- Small app improvements (36ab321)


- Router tests in place (5c37eb4)


- Starting database tests (1ca5832)


- Model tests complete, adapters extracted (4cdb057)


- New adapter tests (3322422)


- Corrected duplication (0acd701)


- Finish unit testing (ddf5fec)


- Restructure project (211598d)


- Restricture tests (6a698b8)


- Restructured to use Giphy API key env variable (1ef5235)


- Remove test data, use fixtures instead (198a9d0)


- API tests ready (80466c3)


- Merge branch 'master' into testing (03f8714)


- Slight format changes (f30c158)


### Testing


-  (816f65c)

Testing


## [0.1] - 2016-03-20

### Added


- Adding design aspects (113800d)


- Add a lovely simple editor (edc4bd4)


### Other


- Initial commit (0792f95)


- Splitting templates into separate files, tested working (02bf6c5)


- Form template and dumping of values (2f102b1)


- Entries now saving from form (a446de9)


- Working to display ALL posts in list (20af2a0)


- Individual entry working (041229e)


- Reorganised into separate files, much much cleaner (f609639)


- 404 working (90ef3f1)


- Abstracted the Journal-based interactions into the Journal file, controllers no longer handle database (b90e582)


- Built in validation and saving, think version 0.1 is done (1b45f72)


- Options (d76e890)


- Bullets on options (437ddff)


- Much better SRP handling - now have a server and database handled purely through model (7d85638)



