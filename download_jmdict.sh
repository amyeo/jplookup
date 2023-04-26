#!/bin/bash

echo "Downloading file ..."
wget "ftp://ftp.edrdg.org/pub/Nihongo//JMdict.gz" -O "JMdict.gz"
echo "Unpacking ..."
gzip -dk JMdict.gz
echo "Done. Run the --init command to build the index."