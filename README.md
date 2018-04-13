# BOSH rules for Bazel

## About

These rules provide support for building BOSH releases using the [Bazel][bazel]
build system. This may be useful to you if you're already using Bazel for
building your project and would like to start using BOSH for deployment.

This project is still in early stages but it can already build valid compiled
releases for BOSH directors.

[bazel]: https://bazel.build/

## Usage

You should familiarize yourself with the [components of a BOSH
release][release] if you haven't already done so. The rules do not use the
standard toolchain to build the release but the basic components are still the
same.

[release]: https://bosh.io/docs/create-release.html

The core API has 3 rules at the moment. Here are some examples below from the
[BPM branch which uses Bazel][bpm-branch] to build its release. I'm still
working on this documentation: expect something more proper in the future.

[bpm-branch]: https://github.com/cloudfoundry-incubator/bpm-release/tree/bazel

``` python
load("@com_github_xoebus_rules_bosh//bosh:def.bzl", "bosh_package")

bosh_package(
    name = "bpm",
    srcs = [
        "//bpm/cmd/bpm",
    ],
)
```

``` python
load("@com_github_xoebus_rules_bosh//bosh:def.bzl", "bosh_release")

bosh_release(
    name = "bpmrelease",
    jobs = [
        "//bosh/jobs/bpm:bpm",
        "//bosh/jobs/test-server:test-server",
    ],
    packages = [
        ":bpm",
        ":bpm-runc",
        ":test-server",
    ],
    stemcell_distro = "ubuntu-trusty",
    stemcell_version = "1234",
)
```

``` python
load("@com_github_xoebus_rules_bosh//bosh:def.bzl", "bosh_job")

bosh_job(
    name = "bpm",
    monit = ":monit",
    spec = ":spec",
    templates = [
        "templates/bpm",
        "templates/pre-start.erb",
        "templates/setup.erb",
    ],
    visibility = ["//bosh:__pkg__"],
)
```
