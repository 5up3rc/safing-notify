#/bin/bash

which 2goarray >/dev/null
if [ $? -ne 0 ]; then
  echo "2goarray not installed."
  echo "run to install: go get github.com/cratonica/2goarray"
  exit 1
fi

if [ ! -f "$GOPATH/src/github.com/mattn/go-gtk/tools/make_inline_pixbuf/make_inline_pixbuf.go" ]; then
  echo "2goarray not installed."
  echo "run to install: go get github.com/cratonica/2goarray"
  exit 1
fi

for entry in $(ls *.png); do

  name=${entry%.png}
  name=${name#icon}
  filename="${name}.go"

  # png byte array
  echo "//+build darwin" > png$filename
  cat "$entry" | 2goarray $name icons >> png$filename
  if [ $? -ne 0 ]; then
    echo "Error processing PNG $entry"
    exit
  fi

  # GTK pixbuf
  echo -e "//+build linux\n" > pb$filename
  echo "package icons" >> pb$filename
  go run $GOPATH/src/github.com/mattn/go-gtk/tools/make_inline_pixbuf/make_inline_pixbuf.go ${name} ${entry} | tail -n +2 >> pb$filename
  if [ $? -ne 0 ]; then
    echo "Error processing Pixbuf $entry"
    exit
  fi

done
