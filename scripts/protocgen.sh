#!/usr/bin/env bash
set -e
home=$PWD

proto_dirs=$(find ./ -name 'buf.yaml' -print0 | xargs -0 -n1 dirname | sort | uniq)
for dir in $proto_dirs; do
  echo "Generating proto code for $dir"

  cd $dir

  if [ -f "buf.gen.pulsar.yaml" ]; then
    buf generate --template buf.gen.pulsar.yaml
    if [ -d "./api" ]; then
      cp -r ./api $home
      rm -rf ./api
    fi
  fi

  if [ -f "buf.gen.gogo.yaml" ]; then
      for file in $(find . -maxdepth 5 -name '*.proto'); do
        if grep -q "option go_package" "$file"; then
          buf generate --template buf.gen.gogo.yaml $file
        fi
    done

    if [ -d "./mandu/x" ]; then
      cp -r ./mandu/x $home
      rm -rf ./mandu/x
    fi
  fi

  cd $home
done

go mod tidy
