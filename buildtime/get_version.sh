#!/bin/sh

# Get the latest tag name.
version="latest"

if command -v git >/dev/null
then
  version=$(git describe --tags --always --abbrev=0)
fi

# APP_VERSION can be passed as ARG (environment variable) for docker builds.
if [ -n "$APP_VERSION" ]; then
  version="$APP_VERSION"
fi

echo "$version" > version.txt
