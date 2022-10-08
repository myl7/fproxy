# fproxy

HTTP reverse proxy forwarding file access with local file persistence (a.k.a. cache but separated by file)

## Motivation

The situation is that, there are some remote files served by HTTP, and we want to not only reversely proxy them, but save them in our servers when access comes to avoid, e.g., the provider stopping serving the files

nginx can be configured to cache the responses permanently.
But the cache is integrated and it is not easy to extract the original files from them.

The program instead trys to cache/persist the files the same as their remote hierarchy

## Features

## License

Copyright (C) 2022 myl7

SPDX-License-Identifier: Apache-2.0
