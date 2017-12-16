#!/usr/bin/env python
# written for python 3!

import colorama
import gzip
import json
import sys
import os
from os import path
from urllib import (request, error as urllib_error)


class NonValidationError(Exception):
    pass


class Message:
    def __init__(self, d):
        self.type = d.get("type")
        self.message = d.get("message")
        self.extract = d.get("extract")
        self.line = d.get("firstLine", d.get("lastLine"))
        self.column = d.get("firstColumn")
        self.highlight_start = d.get("hiliteStart")
        self.highlight_length = d.get("hiliteLength")

        # if we have a subtype, it may be a warning, or a fatal error
        if "subtype" in d:
            self.type = d.get("subtype")

        # pick colors based off message type
        if self.type == "info":
            self.type_color = colorama.Fore.RESET
            self.highlight_color = colorama.Back.RESET
        elif self.type == "warning":
            self.type_color = colorama.Fore.YELLOW
            self.highlight_color = colorama.Back.YELLOW
        elif self.type == "error":
            self.type_color = colorama.Fore.RED + colorama.Style.BRIGHT
            self.highlight_color = colorama.Back.RED
        elif self.type == "fatal":
            self.type_color = colorama.Fore.RED + colorama.Style.DIM
            self.highlight_color = colorama.Back.RED
        else:
            # something unrelated to validating happened
            raise NonValidationError(self.message)

        self.extract = "{}{}{}{}{}".format(self.extract[:self.highlight_start], self.highlight_color,
                                           self.extract[
                                           self.highlight_start:self.highlight_start + self.highlight_length],
                                           colorama.Style.RESET_ALL,
                                           self.extract[self.highlight_start + self.highlight_length:])

    def __str__(self):
        return "{}{}({}:{}): {}{}\n{}\n".format(
            self.type_color, self.type, self.line, self.column,
            self.message, colorama.Style.RESET_ALL, self.extract)

    def __repr__(self):
        return self.__str__()


class Result:
    def __init__(self, filepath, json_str):
        self.filepath = filepath
        self.messages = [Message(m) for m in json_str["messages"]]

    def __str__(self):
        if len(self.messages) != 0:
            return "{}:\n\t{}".format(self.filepath, "\n\t".join(map(str, self.messages)))
        else:
            return ""

    def __repr__(self):
        return self.__str__()


def validate_html_file(filepath):
    try:
        req = request.Request("https://validator.w3.org/nu/?out=json", method="POST")
        req.add_header("Content-Type", "text/html; charset=utf-8")
        req.add_header("Accept-Encoding", "gzip")

        with open(filepath, 'rb') as f:
            req.data = f.read()

        with request.urlopen(req) as resp:
            if resp.info().get("Content-Encoding") == "gzip":
                resp_buf = gzip.decompress(resp.read())
            else:
                resp_buf = resp.read()

        return Result(filepath, json.loads(resp_buf.decode("utf-8")))

    except urllib_error.URLError as e:
        print("Error in network comms for cache '{}': {}", filepath, e)

    except NonValidationError as e:
        print("Error unrelated to validation for cache '{}': {}", filepath, e)


def html_tool(paths):
    results = []

    for p in paths:
        results.append(validate_html_file(p))

    return results


def help_tool():
    print("Usage: validate.py tool paths...")
    print("'tool' can be one of: html, help")
    print("'paths' is a list of files or directories to walk")


def find_children(dir_path, file_type):
    filepaths = []

    # for all files, if its extension (minus the period) matches file_type, add it to our files list
    # recurse on sub directories
    for root, subDirs, files in os.walk(top=dir_path, followlinks=True):
        for f in files:
            path_split = path.splitext(f)
            if len(path_split) > 1 and path_split[1] != "" and path_split[1][1:] == file_type:
                filepaths.append(path.join(root, f))

        for sub in subDirs:
            filepaths.extend(find_children(path.join(root, sub), file_type))

    return filepaths


def main(args):
    if len(args) < 2 or args[1] == "-h" or args[1] == "--help" or args[1] == "help":
        help_tool()
        return

    tool = args[1]
    paths = args[2:]
    files = []

    for p in paths:
        if path.isdir(p):
            files.extend(find_children(p, tool))
        else:
            files.append(p)

    if tool == "html":
        results = html_tool(files)
    else:
        print("Invalid tool name.")
        help_tool()
        return

    for r in results:
        print(r)


if __name__ == "__main__":
    main(sys.argv)
