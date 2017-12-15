#!/bin/bash

export TMP_DIR=/tmp/esalert-temp
mkdir $TMP_DIR

cp -rf sample static $TMP_DIR

go build esalert.go
go build check_rule.go

cp esalert check_rule $TMP_DIR

cd $TMP_DIR
tar -czvf esalert.tar.gz ./

cd -
cp $TMP_DIR/esalert.tar.gz ./
rm -rf $TMP_DIR check_rule esalert

