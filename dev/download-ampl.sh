#!/bin/bash
cd `dirname $0`
BASEDIR=`pwd`
if [ "$#" -ne 1 ]; then
  	echo "Usage: $0 <URL>"
    exit 1
fi
set -ex

URL=$1
PACKAGE=`basename $URL`
curl -k -O $URL
if [[ $PACKAGE == *.zip ]]; then
    unzip $PACKAGE
else
    tar xzvf $PACKAGE
fi
rm $PACKAGE
mv ampl.lic.base ampl_lic
mv ampl.* ampl
cp ampl_lic ampl/ampl.lic
cd ampl
pwd
./ampl -vvq
