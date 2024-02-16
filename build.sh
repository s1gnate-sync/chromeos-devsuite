#!/bin/sh

set -eu

make -C sucrosh
make -C devcoo
make -C sethosts
make -C skipass
make -C runsu

mv sucrosh/sucrosh bin	
mv devcoo/devcoo bin	
mv sethosts/sethosts bin	
mv skipass/skipass bin	
mv -f runsu/runsu bin	
