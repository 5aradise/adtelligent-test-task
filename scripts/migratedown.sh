#!/bin/bash

cd sql/schema
goose postgres postgres://paradise:7889@127.0.0.1:5432/adtelligent down
