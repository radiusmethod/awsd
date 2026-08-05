## v0.3.0 (August 5, 2026)
* Added `awsd init <shell>` — one line in your rc file (`eval "$(awsd init zsh)"`) now replaces the `awsd` alias, the completion `source`, and the persistence snippet from the README. [#53]
* Added fish support, and PowerShell tab completion.
* Added `awsd shellenv [shell]`, which prints the export/unset statements for the active profile and region. This is what the generated function evals.
* `~/.awsd` parsing now lives only in the Go code. `scripts/_awsd` and `scripts/powershell/awsd.ps1` are thin shims over `awsd shellenv`.
* Values are now shell-quoted, so profile names containing spaces or quotes work.
* Deprecated `alias awsd="source _awsd"` and `source _awsd_autocomplete`. Both still work; use `awsd init` instead.

## v0.2.0 (April 27, 2026)
* Added region switching: `awsd set region [name]` (interactive picker if no name given), `awsd unset region`, `awsd list regions`.
* Added `awsd set profile [name]` and `awsd unset profile` as explicit forms — bare `awsd <profile>` still works.
* `~/.awsd` now uses a `key=value` format (`profile=...` / `region=...`). Legacy single-line files are still readable; the next write upgrades them.
* Wrapper now also exports `AWS_DEFAULT_REGION` alongside `AWS_REGION`.

## v0.1.3 (March 9, 2025)
* Adds circular scrolling

## v0.1.2 (July 24, 2024)
* Added PowerShell support. [#19] thanks, @jinxiao
* Ensure the ~/.awsd file exists before attempting to read its contents to prevent errors.

## v0.1.1 (June 28, 2024)
* Replaced regular expressions with the `ini` package for extracting profiles from AWS config files.
* Fixed an issue where extra spaces between the "profile" keyword and the profile name could prevent the profile from being set.

## v0.1.0 (May 2, 2024)
* Increase zsh autocompletion compatibility.

## v0.0.9 (April 4, 2024)
* Fixes issue with help command shorthand flag `-h`. [#20] thanks, @Masamerc

## v0.0.8 (December 23, 2023)
* Added autocomplete script to install.

## v0.0.7  (October 20, 2023)
* Update for new organization.

## v0.0.6  (October 6, 2023)
* Refactored codebase.

## v0.0.5  (October 6, 2023)
* Added support for passing arbitrary profile names as arguments. [#4] thanks, @withakay
* Added `awsd list` command to simply list all profiles.

## v0.0.4  (October 4, 2023)
* Don't append default profile if it already exists.

## v0.0.3  (September 22, 2023)
* Added additional error checking.

## v0.0.2  (May 27, 2022)
* Added ability to search AWS profiles. [#4] thanks, @M1kep

## v0.0.1 (November 30, 2021)
* Initial Release
