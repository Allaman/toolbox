#!/bin/bash

xelatex "$1"
rm -rf -- *.aux *.log *.out *.sync-conflict* .cache
