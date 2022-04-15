#!/bin/bash
cp -r js/* static/js && cd static/js/ && find . -name '*.js' -exec minify {} -o {} \;
