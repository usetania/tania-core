# Changelog
All notable changes to this project will be documented in this file.

The format is based on [Keep a Changelog](http://keepachangelog.com/en/1.0.0/)
and this project currently doesn't strictly adhere to [Semantic Versioning](http://semver.org/spec/v2.0.0.html).

## [Unreleased]
### Added 
- Add `CHANGELOG.md` based on the [Keep a Changelog](http://keepachangelog.com/en/1.0.0/)
- Add `app_port` config for configurable backend port

### Changed
- Change [paked/configure](https://github.com/paked/configure) package with [spf13/viper](https://github.com/spf13/viper) because [paked/configure](https://github.com/paked/configure) doesn't support config of slice
- Change `redirect_uri` config to use array of string instead of single string value to handle multiple host

## [1.5.1] - 2018-04-14
### Fixed
- https://github.com/Tanibox/tania-core/issues/9

## [1.5.0] - 2018-04-03
### Added
- Rewrite from PHP to Go
- Feature: Crop batch management. Now you can make a group of plants to grow in the same area or even different area.
- Feature: Inventories management. You can record not only seeds but also your other farm tools and utilities.