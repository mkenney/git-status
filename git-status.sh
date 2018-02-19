#!/usr/bin/env bash

# This is a simple, clean `git` status line for your shell prompt. The
# `git-status.sh` script defines a function called, creatively,
# `__git_status()` that returns a string indicating the current state of
# your local repository. The function returns a string describing:
#
# * `origin`/`local`/`tag`/`detached` origin indicator
# * branch name/tag name/commit hash position indicator
# * number of commits behind origin: `>n`
# * number of commits ahead of origin: `<n`
# * number of untracked files: `?n`
# * number of deleted files: `Dn`
# * number of added files: `+n`
# * number of modified files: `𝚫n`
# * number of renamed files: `↪n`
# * number of files with unstaged changes: `∴n`
# * total number of files: `#n`
#
# MIT License
# Copyright (c) 2018 Michael Kenney
#
# Permission is hereby granted, free of charge, to any person obtaining a copy
# of this software and associated documentation files (the "Software"), to deal
# in the Software without restriction, including without limitation the rights
# to use, copy, modify, merge, publish, distribute, sublicense, and/or sell
# copies of the Software, and to permit persons to whom the Software is
# furnished to do so, subject to the following conditions:
#
# The above copyright notice and this permission notice shall be included in all
# copies or substantial portions of the Software.
#
# THE SOFTWARE IS PROVIDED "AS IS", WITHOUT WARRANTY OF ANY KIND, EXPRESS OR
# IMPLIED, INCLUDING BUT NOT LIMITED TO THE WARRANTIES OF MERCHANTABILITY,
# FITNESS FOR A PARTICULAR PURPOSE AND NONINFRINGEMENT. IN NO EVENT SHALL THE
# AUTHORS OR COPYRIGHT HOLDERS BE LIABLE FOR ANY CLAIM, DAMAGES OR OTHER
# LIABILITY, WHETHER IN AN ACTION OF CONTRACT, TORT OR OTHERWISE, ARISING FROM,
# OUT OF OR IN CONNECTION WITH THE SOFTWARE OR THE USE OR OTHER DEALINGS IN THE
# SOFTWARE.

__git_status() {
    set -o pipefail
    local added=0
    local added_str=
    local ahead_str=
    local behind_str=
    local deleted=0
    local deleted_str=
    local detached=0
    local modified=0
    local modified_str=
    local output=0
    local ref_name=$(basename $(git symbolic-ref HEAD 2> /dev/null) 2> /dev/null)
    local renamed=0
    local renamed_str=
    local total=0
    local total_str=
    local tree_position=
    local unstaged_str=
    local untracked=0
    local untracked_str=

    git rev-parse HEAD &> /dev/null; exit_code=$?
    if [ "0" = "$exit_code" ]; then
        hash=$(git rev-parse --short=10 HEAD 2> /dev/null); exit_code=$?
    fi
    if [ "0" != "$exit_code" ]; then
        git rev-parse --abbrev-ref HEAD &> /dev/null; exit_code=$?
        if [ "0" = "$exit_code" ]; then
            hash=$(git rev-parse --abbrev-ref HEAD)
        else
            hash="master"
        fi
    fi

    ref_name=$(git rev-parse --abbrev-ref HEAD 2> /dev/null); exit_code=$?
    ref_source="origin"
    if [ "0" != "$exit_code" ]; then
        # See if on a tag
        ref_name="$(git describe --tags 2> /dev/null)"; exit_code=$?
        ref_source="tag"
        if [ "0" != "$exit_code" ]; then

            # See if referencing a name
            ref_name=$(basename $(git symbolic-ref HEAD 2> /dev/null) 2> /dev/null); exit_code=$?
            ref_source="local"
            if [ "0" != "$exit_code" ]; then
                # Assume detached state
                ref_name=$hash
                ref_source="detached"
            fi
        fi
    else
        tree_position=$(git for-each-ref --format='%(upstream:short)' $(git symbolic-ref -q HEAD 2> /dev/null) | head -1) 2> /dev/null
        if [ "" = "$tree_position" ]; then
            ref=$(git symbolic-ref -q HEAD 2> /dev/null)
            if [ "" = "$ref" ]; then
                ref_name="$hash"
                ref_source="detached"
            else
                ref_name=$(basename $ref)
                ref_source="local"
            fi
        elif grep -q '\->' <<< $(git show -s --pretty=%d HEAD); then
            ref_name=$(basename $tree_position)
            ref_source=$(dirname $tree_position)
        else
            ref_name="$hash"
            ref_source="detached"
        fi
    fi

    # Branch information
    compare_ref="HEAD"
    if [ "origin" = "$ref_source" ]; then
        compare_ref="${ref_source}/${ref_name}"
    fi

echo git rev-list --left-right --count $hash...$compare_ref
    rev_list=$(git rev-list --left-right --count $hash...$compare_ref &> /dev/null); exit_code=$?
    if [ "0" = "$exit_code" ] && [ "" != "$rev_list" ]; then

        ahead_str="<$rev_list | awk '{print $1}') "
        behind_str=">$rev_list | awk '{print $2}') "
        output=1
    fi

    #git rev-list $compare_ref...$hash &> /dev/null; exit_code=$?
    #if [ "0" = "$behind_str" ] && [ "" != "$(git rev-list $compare_ref...$hash)" ]; then
    #    behind_str=">$(echo $(git rev-list $compare_ref...$hash) | wc | awk '{print $1}') "
    #    output=1
    #fi

    # Files with unstaged changes
    if [ "" != "$(git diff --name-only)" ]; then
        # ♐ ± ~ ∵ ∴
        unstaged_str="∴$(git diff --name-only | wc | awk '{print $1}') "
        output=1
    fi

    # Tabulate all change states
    while read line; do
        flag1=${line:0:1}
        flag2=${line:1:1}
        if [ "" != "$line" ]; then
            total=$((total + 1))

            # Added files
            if [ "A" = "$flag1" ] || [ "A" = "$flag2" ]; then
                added=$((added + 1))
                # +
                added_str="+$added "
                output=1
            fi

            # Deleted files
            if [ "D" = "$flag1" ] || [ "D" = "$flag2" ]; then
                deleted=$((deleted + 1))
                # × ␥ ␡
                deleted_str="D$deleted "
                output=1
            fi

            # Modified files
            if [ "M" = "$flag1" ] || [ "M" = "$flag2" ]; then
                modified=$((modified + 1))
                # ≠ ≢ 𝚫
                modified_str="𝚫$modified "
                output=1
            fi

            # Renamed files
            if [ "R" = "$flag1" ] || [ "R" = "$flag2" ]; then
                renamed=$((renamed + 1))
                # ⤿  ↪
                renamed_str="↪$renamed "
                output=1
            fi

            # Untracked files
            if [ "?" = "$flag1" ] || [ "?" = "$flag2" ]; then
                untracked=$((untracked + 1))
                # ∑ ?
                untracked_str="?$untracked "
                output=1
            fi
        fi
    done << EOF
$(git status --porcelain)
EOF

    # Total files
    if [ "0" != "$total" ]; then
        total_str="#$total "
        output=1
    fi

    if [ "0" != "$output" ]; then
        echo "$(echo -e "${ref_source}/${ref_name} ${behind_str}${ahead_str}${untracked_str}${deleted_str}${added_str}${renamed_str}${modified_str}${unstaged_str}${total_str}" | sed -e 's/[[:space:]]*$//')"
    else
        echo "${ref_source}/${ref_name}"
    fi
}

export -f __git_status
