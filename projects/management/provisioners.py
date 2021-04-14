# Copyright 2020, Pulumi Corporation.  All rights reserved.

import abc
import json
import hashlib
import io
import paramiko
import pulumi
from pulumi import dynamic, Input, Output
from pulumi.resource import Resource
import socket
import time
from typing import Any, Optional
from typing_extensions import TypedDict
from uuid import uuid4


def sha256sum(filename):
    h = hashlib.sha256()
    with open(filename, "rb") as f:
        data = f.read()
        h.update(data)
    return h.hexdigest()


def compare(a, b):
    try:
        val_a = json.dumps(a, sort_keys=True, indent=2)
        val_b = json.dumps(b, sort_keys=True, indent=2)
        if val_a != val_b:
            return False
    except TypeError:
        return False
    return True


# ConnectionArgs tells a provisioner how to access a remote resource. It includes the hostname
# and optional port (default is 22), username, password, and private key information.
@pulumi.input_type
class ConnectionArgs:
    host: pulumi.Input[str]
    port: Optional[pulumi.Input[int]]
    username: pulumi.Input[str]
    private_key: pulumi.Input[str]

    def __init__(self, host, username, private_key, port=22):
        self.host = host
        self.port = port
        self.username = username
        self.private_key = private_key


def connect(conn: dict) -> paramiko.SSHClient:
    ssh = paramiko.SSHClient()
    ssh.set_missing_host_key_policy(paramiko.AutoAddPolicy())
    # pkey = paramiko.RSAKey.from_private_key_file(filename=conn['private_key_file'])
    fileIO = io.StringIO(conn["private_key"])
    pkey = paramiko.RSAKey.from_private_key(fileIO)
    # Retry the connection until the endpoint is available (up to 2 minutes).
    retries = 0
    while True:
        try:
            ssh.connect(
                hostname=conn["host"],
                port=int(conn["port"]),
                username=conn["username"],
                pkey=pkey,
            )
            print("connection done")
            return ssh
        except (paramiko.SSHException, socket.error):
            if retries == 24:
                raise
            time.sleep(5)
            retries = retries + 1
            pass
        except Exception as e:
            print(e)
            raise e


class ProvisionerProvider(dynamic.ResourceProvider):
    __metaclass__ = abc.ABCMeta

    @abc.abstractmethod
    def on_create(self, inputs: Any) -> Any:
        return

    def create(self, inputs):
        outputs = self.on_create(inputs)
        return dynamic.CreateResult(id_=uuid4().hex, outs=outputs)

    def diff(self, _id, olds, news):
        # If anything changed in the inputs, replace the resource.
        diffs = []
        for key in olds:
            if key not in news:
                diffs.append(key)
            else:
                try:
                    olds_value = json.dumps(olds[key], sort_keys=True, indent=2)
                    news_value = json.dumps(news[key], sort_keys=True, indent=2)
                    if olds_value != news_value:
                        diffs.append(key)
                except TypeError:
                    diffs.append(key)
        for key in news:
            if key not in olds:
                diffs.append(key)
        return dynamic.DiffResult(
            changes=len(diffs) > 0, replaces=diffs, delete_before_replace=True
        )


# CopyFileProvider implements the resource lifecycle for the CopyFile resource type below.
class CopyFileProvider(ProvisionerProvider):
    def on_create(self, inputs: Any) -> Any:
        ssh = connect(inputs["conn"])
        scp = ssh.open_sftp()
        try:
            scp.put(inputs["src"], inputs["dest"])
        finally:
            scp.close()
            ssh.close()
        return inputs


# CopyFile is a provisioner step that can copy a file over an SSH connection.
class CopyFile(dynamic.Resource):
    def __init__(
        self,
        name: str,
        host_id: str,
        conn: pulumi.Input[ConnectionArgs],
        src: str,
        dest: str,
        opts: Optional[pulumi.ResourceOptions] = None,
    ):
        super().__init__(
            CopyFileProvider(),
            name,
            {
                "host_id": host_id,
                "conn": conn,
                "src": src,
                "dest": dest,
                "fileHash": sha256sum(src),
            },
            opts,
        )


# CopyFileFromStringProvider implements the resource lifecycle for CopyFileFromString resource type
class CopyFileFromStringProvider(ProvisionerProvider):
    def on_create(self, props: Any) -> Any:
        ssh = connect(props["conn"])
        scp = ssh.open_sftp()
        try:
            b = io.BytesIO(bytes(props["from_str"], encoding="utf8"))
            scp.putfo(b, props["dest"])
        finally:
            scp.close()
            ssh.close()
        return props


# CopyFileFromString is a provisioner that copy string to a new file on remote host over SSH connection
class CopyFileFromString(dynamic.Resource):
    def __init__(
        self,
        name: str,
        host_id: str,
        conn: ConnectionArgs,
        from_str: str,
        dest: str,
        opts: pulumi.ResourceOptions = None,
    ):
        super().__init__(
            CopyFileFromStringProvider(),
            name,
            {
                "host_id": host_id,
                "conn": conn,
                "from_str": from_str,
                "dest": dest,
            },
            opts,
        )


# RunCommandResult is the result of running a command.
class RunCommandResult(TypedDict):
    stdout: str
    """The stdout of the command that was executed."""
    stderr: str
    """The stderr of the command that was executed."""


# RemoteExecProvider implements the resource lifecycle for the RemoteExec resource type below.
class RemoteExecProvider(ProvisionerProvider):
    def on_create(self, inputs: Any) -> Any:
        ssh = connect(inputs["conn"])
        try:
            for command in inputs["commands"]:
                stdin, stdout, stderr = ssh.exec_command(command)
                err = "".join(stderr.readlines())
                if err != "":
                    raise Exception(
                        'remote execution "{}" failed: {}'.format(command, err)
                    )
        finally:
            ssh.close()
        return inputs


# RemoteExec runs remote one or more commands over an SSH connection. It returns the resulting
# stdout and stderr from the commands in the results property.
class RemoteExec(dynamic.Resource):
    def __init__(
        self,
        name: str,
        host_id: str,
        conn: ConnectionArgs,
        commands: list,
        opts: Optional[pulumi.ResourceOptions] = None,
    ):
        self.conn = conn
        """conn contains information on how to connect to the destination, in addition to dependency information."""
        self.commands = commands
        """The commands to execute. Exactly one of 'command' and 'commands' is required."""
        super().__init__(
            RemoteExecProvider(),
            name,
            {
                "host_id": host_id,
                "conn": conn,
                "commands": commands,
            },
            opts,
        )
