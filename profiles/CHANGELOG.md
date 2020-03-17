# Changelog

All notable changes to this project will be documented in this file.

The format is based on [Keep a Changelog](https://keepachangelog.com/en/1.0.0/)

## [Unreleased]

## 2020-03-17

### Added

- README.md
- CHANGELOG.md
- Parser to parse traffic signals when process node for unidb in lib-unidb/relations.lua 
- type of traffic_signals in relation_types in car-unidb.lua to load relations when it contains traffic signals. 

### Removed
- type of route in relations_types in car-unidb.lua since there is no type of route in unidb

## 2020-03-15

### Changed

- Revert PR-67 to keep car.lua and lib/* are the same as before

### Removed

- test_speed_unit.lua which is unneccessary

## 2020-03-11

### Added

- lib-unidb to hold all scripts related to unidb
- car-unidb.lua to support unidb for car style

