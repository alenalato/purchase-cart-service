#!/bin/bash

docker run --rm --network host -v "$(pwd)":/src grafana/k6 run /src/e2e_tests/order_post.js
