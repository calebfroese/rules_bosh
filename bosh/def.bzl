def _bosh_release_impl(ctx):
    inputs = [f for f in ctx.files.packages] + [f for f in ctx.files.jobs]
    outputs = [ctx.outputs.out]

    args = ctx.actions.args()
    args.add(["-output", ctx.outputs.out.path])
    args.add(["-name", ctx.label.name])
    args.add(["-stemcellDistro", ctx.attr.stemcell_distro])
    args.add(["-stemcellVersion", ctx.attr.stemcell_version])
    for package in ctx.files.packages:
        args.add("-package")
        args.add(package.path)
    for job in ctx.files.jobs:
        args.add("-job")
        args.add(job.path)

    ctx.actions.run(
        inputs=inputs,
        outputs=outputs,
        arguments=[args],
        mnemonic="BoshRelease",
        progress_message="Building %s BOSH release" % ctx.label.name,
        executable=ctx.executable._builder,
    )

bosh_release = rule(
    _bosh_release_impl,
    attrs = {
        "jobs": attr.label_list(
            mandatory = True,
        ),
        "packages": attr.label_list(
            mandatory = True,
        ),
        "stemcell_distro": attr.string(
            mandatory = True,
        ),
        "stemcell_version": attr.string(
            mandatory = True,
        ),
        "_builder": attr.label(
            default = Label("//bosh:buildrel"),
            allow_single_file = True,
            executable = True,
            cfg = "host",
        )
    },
    outputs = {
        "out": "%{name}.tgz",
    },
)

def _bosh_job_impl(ctx):
    inputs = [f for f in ctx.files.templates] + [ctx.file.monit, ctx.file.spec]
    outputs = [ctx.outputs.out]

    args = ctx.actions.args()
    args.add(["-output", ctx.outputs.out.path])
    args.add(["-manifest", ctx.file.spec.path])
    args.add(["-monit", ctx.file.monit.path])
    for template in ctx.files.templates:
        args.add("-template")
        args.add(template.path)

    ctx.actions.run(
        inputs=inputs,
        outputs=outputs,
        arguments=[args],
        mnemonic="BoshJob",
        progress_message="Building %s BOSH job" % ctx.label.name,
        executable=ctx.executable._builder,
    )

bosh_job = rule(
    _bosh_job_impl,
    attrs = {
        "templates": attr.label_list(
            allow_files = True,
        ),
        "monit": attr.label(
            allow_single_file = True,
            mandatory = True,
        ),
        "spec": attr.label(
            allow_single_file = True,
            mandatory = True,
        ),
        "_builder": attr.label(
            default = Label("//bosh:buildjob"),
            allow_single_file = True,
            executable = True,
            cfg = "host",
        )
    },
    outputs = {
        "out": "%{name}.tgz",
    },
)

def _bosh_package_impl(ctx):
    inputs = [f for f in ctx.files.srcs]
    outputs = [ctx.outputs.out]

    args = ctx.actions.args()
    args.add(["-output", ctx.outputs.out.path])
    for pkg in ctx.files.srcs:
        args.add("-file")
        args.add(pkg.path)

    ctx.actions.run(
        inputs=inputs,
        outputs=outputs,
        arguments=[args],
        mnemonic="BoshPackage",
        progress_message="Building %s BOSH package" % ctx.label.name,
        executable=ctx.executable._builder,
    )

bosh_package = rule(
    _bosh_package_impl,
    attrs = {
        "srcs": attr.label_list(
            allow_files = True,
            mandatory = True,
        ),
        "_builder": attr.label(
            default = Label("//bosh:buildpkg"),
            allow_single_file = True,
            executable = True,
            cfg = "host",
        )
    },
    outputs = {
        "out": "%{name}.tgz",
    },
)