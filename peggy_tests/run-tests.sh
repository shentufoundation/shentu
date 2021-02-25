#!/bin/bash
TEST_TYPE=$1
set -eu

# Run test entry point script
docker exec shentu_test_instance /bin/sh -c "pushd /shentu/ && peggy_tests/container-scripts/integration-tests.sh 1 $TEST_TYPE"