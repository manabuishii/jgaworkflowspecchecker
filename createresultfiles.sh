#!/bin/bash

#for FILE in `cat output/HG001740.9.0.result.txt`
for FILE in `cat XX00000.txt`
do
  echo $FILE > test/resultfile/success/$FILE
done
