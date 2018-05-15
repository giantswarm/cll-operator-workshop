#!/usr/bin/env bash

repo="github.com/giantswarm/cll-operator-workshop"
group_versions="GROUP:VERSION"

dir="$( cd "$( dirname "${BASH_SOURCE[0]}" )" && pwd )"

cd ${dir}/../vendor/k8s.io/code-generator && ./generate-groups.sh \
    "deepcopy,client" \
    $repo/pkg \
    $repo/pkg/apis \
    $group_versions \
    --go-header-file ${dir}/boilerplate.go.txt
