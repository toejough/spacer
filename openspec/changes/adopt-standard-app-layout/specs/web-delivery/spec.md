## ADDED Requirements

### Requirement: The app publishes named assets, not a directory
The app SHALL serve only assets it has explicitly published, from a designated published directory. Content that merely sits beside those assets — test files, fixtures, notes, source — SHALL NOT be reachable. No directory listing SHALL be returned for any path.

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

### Requirement: The server does not depend on where it was compiled
The app SHALL locate its templates and assets by a path known at run time, or carry them inside the binary. It SHALL NOT derive them from the location the source occupied at compile time.

#### Scenario: The repository moves and the app still serves
- **WHEN** the built binary is run with the repository at a different path from the one it was compiled at
- **THEN** pages render normally, rather than failing in a way that looks like a routing fault

#### Scenario: A rebuild under the service manager actually rebuilds
- **WHEN** the app is started by the service manager, whose environment does not include the variables a toolchain needs
- **THEN** the build either succeeds or fails loudly — it does not fall through to a previously built binary while reporting the app as started
