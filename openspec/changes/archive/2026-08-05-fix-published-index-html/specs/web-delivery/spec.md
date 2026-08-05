## MODIFIED Requirements

### Requirement: The app publishes named assets, not a directory

The app SHALL serve only assets it has explicitly published, from a designated published directory. Content that merely sits beside those assets — test files, fixtures, notes, source — SHALL NOT be reachable. No directory listing SHALL be returned for any path. The entry page SHALL be published the same way as every other asset: linked to its source, not maintained as a separate hand-authored copy.

#### Scenario: A test file beside the assets is not served
- **WHEN** a test file sits in the same source directory as the app's stylesheets and scripts, and is requested by exact path
- **THEN** it is refused, while the stylesheets and scripts are served

#### Scenario: The asset directory cannot be browsed
- **WHEN** the asset path is requested as a directory
- **THEN** no index is returned, so what is present cannot be enumerated

#### Scenario: An unpublished file stays unreachable whatever it is called
- **WHEN** a file with an unfamiliar name or extension is added beside the assets and not published
- **THEN** it is unreachable, because reachability follows from having been published

#### Scenario: Editing a published asset is live
- **WHEN** a published asset is edited
- **THEN** the next request returns the edited content, with no rebuild or restart

#### Scenario: The entry page reflects its source, not a stand-in
- **WHEN** the entry page's source template is edited
- **THEN** the next request for the entry page returns the edited content — the same guarantee already made for every other published asset

#### Scenario: The entry page is never served from an intermediate cache
- **WHEN** the entry page is requested
- **THEN** the response carries `Cache-Control: no-store`, the same guarantee already made for assets under `/static/`, so a deploy is never hidden behind a cached HTML document
