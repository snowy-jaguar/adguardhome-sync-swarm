# NOTICE

This project is a derivative of [adguardhome-sync](https://github.com/bakito/adguardhome-sync),
which is licensed under the Apache License, Version 2.0.

This fork adds Docker Swarm support, healthcheck integration, and secret-based
environment variable support using custom entrypoint logic.

---

## Original Copyright

Copyright 2021 bakito

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

> http://www.apache.org/licenses/LICENSE-2.0

---

## Changes and Additions in This Fork

Unless otherwise noted, all changes and additions in this repository are:

Copyright 2025 snowy-jaguar  
Licensed under the Apache License, Version 2.0

The following components have been added or modified:

- `Dockerfile` — Updated for Swarm support and multi-stage build logic
- `entrypoint.sh` — Reads secrets and maps them to environment variables
- `healthcheck.sh` — Implements a flexible healthcheck with optional authentication
- `.github/workflows/build.yml` — GitHub Actions for building multi-platform images with version injection
- `build.sh` — Local build helper with dynamic versioning
- `NOTICE.md` — This file