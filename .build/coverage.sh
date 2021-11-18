#!/bin/bash

set -e

TOTAL=$(go tool cover -func=coverage.out | grep total | grep -Eo '[0-9]+\.[0-9]+')

if (( $(bc <<< "$TOTAL <= 50") )); then
  COLOR=red
elif (( $(bc <<< "$TOTAL > 80") )); then
  COLOR=green
else
  COLOR=orange
fi

BADGE="https:\/\/img.shields.io\/badge\/coverage-$TOTAL%25-$COLOR"

sed "s/!\[CoverageBadge\].*/![CoverageBadge](${BADGE})/g" < README.md > _README.md
mv _README.md README.md
