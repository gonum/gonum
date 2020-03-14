#!/bin/bash

if [ "${TRAVIS_CPU_ARCH}" != "amd64" ];  then
	# TODO(kortschak): Remove this when git fetch --shallow-exclude is supported by git on arm64.
	# See https://travis-ci.community/t/git-fetch-shallow-exclude-does-not-work-on-arm64-architecture/7583
	echo "Skip lint check on arm64"
	exit 0
fi

if [ -n "${TAGS}" ];  then
	echo "Skip redundant lint check"
	exit 0
fi

BRANCH=${TRAVIS_BRANCH}
if [ "${BRANCH}" == "master" ] && [ -n "${TRAVIS_PULL_REQUEST_BRANCH}" ]; then
	BRANCH=${TRAVIS_PULL_REQUEST_BRANCH}
fi

if [ "${BRANCH}" == "master" ]; then
	# Don't run linter on master; it's too late by then.
	exit 0
fi

if [ "${TRAVIS_COMMIT}" == "${TRAVIS_TAG}" ]; then
	# Don't run linter on tag pushes.
	exit 0
fi

set -xe

# Get all the commits on the branch, and the base commit.
git fetch --shallow-exclude=master origin ${BRANCH}
git log --oneline
if [ -z "${TRAVIS_PULL_REQUEST_BRANCH}" ]; then
	# Get the previous commit as well if we are not in a pull request build.
	# Otherwise we already have the base commit.
	COMMITS=$(git log --oneline | wc -l)
	git fetch --depth=$((COMMITS+1)) origin ${BRANCH}
fi

# Travis does not correctly report the branch commit range in ${TRAVIS_COMMIT_RANGE}
# if there has been an amended commit, so just get it from the git log.
SINCE=$(git log --oneline --reverse | head -n1 | cut -f1 -d' ')

# Lint changes since we were on master.
golangci-lint run --new-from-rev=${SINCE}
