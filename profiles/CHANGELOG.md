# Changelog

All notable changes to profiles will be documented in this file.All sentences start with ADDED, REMOVED, CHANGED, or FIXED

The format is based on [Keep a Changelog](https://keepachangelog.com/en/1.0.0/)

## 2020-03-17

- ADDED README.md [#232](https://github.com/Telenav/osrm-backend/pull/232)
- ADDED CHANGELOG.md [#232](https://github.com/Telenav/osrm-backend/pull/232)
- ADDED a parser to parse traffic signals when process node for unidb in lib-unidb/relations.lua [#232](https://github.com/Telenav/osrm-backend/pull/232)
- ADDED traffic_signals in relation_types in car-unidb.lua to load relations when it contains traffic signals [#232](https://github.com/Telenav/osrm-backend/pull/232)
- REMOVED route in relations_types in car-unidb.lua since there is no type of route in unidb [#232](https://github.com/Telenav/osrm-backend/pull/232)

## 2020-03-15

- CHANGED to revert PR-67 to keep car.lua and lib/* are the same as before [#228](https://github.com/Telenav/osrm-backend/pull/228)
- REMOVED test_speed_unit.lua which is unneccessary [#228](https://github.com/Telenav/osrm-backend/pull/228)

## 2020-03-11

- ADDED lib-unidb to hold all scripts related to unidb [#214](https://github.com/Telenav/osrm-backend/pull/214)
- ADDED car-unidb.lua to support unidb for car style [#214](https://github.com/Telenav/osrm-backend/pull/214)

