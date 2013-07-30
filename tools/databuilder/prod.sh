#!/bin/bash

./build.sh
mv -i bayesian.data ../../data/bayesian.data
mv -i dict.data ../../data/dict.data
chmod 666 ../../data/bayesian.data
