#!/bin/bash
_=$(grep Version README.md | awk '{print "git rev-parse refs/tags/"$3" > /dev/null 2>&1;"}' | sh)
if [ "$(echo $?)" == 0 ]
then
  TAG=$(grep Version README.md | awk '{print $3}')
  echo "Version $TAG already exists. Please update the string version on README.md"
	exit 1
fi