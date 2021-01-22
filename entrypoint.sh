#!/bin/sh
set -e

/pulumi/bin/pulumi login file:///pulumi/ccmaas/state
/pulumi/bin/ccmaas server
