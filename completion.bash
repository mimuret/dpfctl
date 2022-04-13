# bash completion for dpfctl                               -*- shell-script -*-

__dpfctl_debug()
{
    if [[ -n ${BASH_COMP_DEBUG_FILE} ]]; then
        echo "$*" >> "${BASH_COMP_DEBUG_FILE}"
    fi
}

# Homebrew on Macs have version 1.3 of bash-completion which doesn't include
# _init_completion. This is a very minimal version of that function.
__dpfctl_init_completion()
{
    COMPREPLY=()
    _get_comp_words_by_ref "$@" cur prev words cword
}

__dpfctl_index_of_word()
{
    local w word=$1
    shift
    index=0
    for w in "$@"; do
        [[ $w = "$word" ]] && return
        index=$((index+1))
    done
    index=-1
}

__dpfctl_contains_word()
{
    local w word=$1; shift
    for w in "$@"; do
        [[ $w = "$word" ]] && return
    done
    return 1
}

__dpfctl_handle_go_custom_completion()
{
    __dpfctl_debug "${FUNCNAME[0]}: cur is ${cur}, words[*] is ${words[*]}, #words[@] is ${#words[@]}"

    local shellCompDirectiveError=1
    local shellCompDirectiveNoSpace=2
    local shellCompDirectiveNoFileComp=4
    local shellCompDirectiveFilterFileExt=8
    local shellCompDirectiveFilterDirs=16

    local out requestComp lastParam lastChar comp directive args

    # Prepare the command to request completions for the program.
    # Calling ${words[0]} instead of directly dpfctl allows to handle aliases
    args=("${words[@]:1}")
    requestComp="${words[0]} __completeNoDesc ${args[*]}"

    lastParam=${words[$((${#words[@]}-1))]}
    lastChar=${lastParam:$((${#lastParam}-1)):1}
    __dpfctl_debug "${FUNCNAME[0]}: lastParam ${lastParam}, lastChar ${lastChar}"

    if [ -z "${cur}" ] && [ "${lastChar}" != "=" ]; then
        # If the last parameter is complete (there is a space following it)
        # We add an extra empty parameter so we can indicate this to the go method.
        __dpfctl_debug "${FUNCNAME[0]}: Adding extra empty parameter"
        requestComp="${requestComp} \"\""
    fi

    __dpfctl_debug "${FUNCNAME[0]}: calling ${requestComp}"
    # Use eval to handle any environment variables and such
    out=$(eval "${requestComp}" 2>/dev/null)

    # Extract the directive integer at the very end of the output following a colon (:)
    directive=${out##*:}
    # Remove the directive
    out=${out%:*}
    if [ "${directive}" = "${out}" ]; then
        # There is not directive specified
        directive=0
    fi
    __dpfctl_debug "${FUNCNAME[0]}: the completion directive is: ${directive}"
    __dpfctl_debug "${FUNCNAME[0]}: the completions are: ${out[*]}"

    if [ $((directive & shellCompDirectiveError)) -ne 0 ]; then
        # Error code.  No completion.
        __dpfctl_debug "${FUNCNAME[0]}: received error from custom completion go code"
        return
    else
        if [ $((directive & shellCompDirectiveNoSpace)) -ne 0 ]; then
            if [[ $(type -t compopt) = "builtin" ]]; then
                __dpfctl_debug "${FUNCNAME[0]}: activating no space"
                compopt -o nospace
            fi
        fi
        if [ $((directive & shellCompDirectiveNoFileComp)) -ne 0 ]; then
            if [[ $(type -t compopt) = "builtin" ]]; then
                __dpfctl_debug "${FUNCNAME[0]}: activating no file completion"
                compopt +o default
            fi
        fi
    fi

    if [ $((directive & shellCompDirectiveFilterFileExt)) -ne 0 ]; then
        # File extension filtering
        local fullFilter filter filteringCmd
        # Do not use quotes around the $out variable or else newline
        # characters will be kept.
        for filter in ${out[*]}; do
            fullFilter+="$filter|"
        done

        filteringCmd="_filedir $fullFilter"
        __dpfctl_debug "File filtering command: $filteringCmd"
        $filteringCmd
    elif [ $((directive & shellCompDirectiveFilterDirs)) -ne 0 ]; then
        # File completion for directories only
        local subDir
        # Use printf to strip any trailing newline
        subdir=$(printf "%s" "${out[0]}")
        if [ -n "$subdir" ]; then
            __dpfctl_debug "Listing directories in $subdir"
            __dpfctl_handle_subdirs_in_dir_flag "$subdir"
        else
            __dpfctl_debug "Listing directories in ."
            _filedir -d
        fi
    else
        while IFS='' read -r comp; do
            COMPREPLY+=("$comp")
        done < <(compgen -W "${out[*]}" -- "$cur")
    fi
}

__dpfctl_handle_reply()
{
    __dpfctl_debug "${FUNCNAME[0]}"
    local comp
    case $cur in
        -*)
            if [[ $(type -t compopt) = "builtin" ]]; then
                compopt -o nospace
            fi
            local allflags
            if [ ${#must_have_one_flag[@]} -ne 0 ]; then
                allflags=("${must_have_one_flag[@]}")
            else
                allflags=("${flags[*]} ${two_word_flags[*]}")
            fi
            while IFS='' read -r comp; do
                COMPREPLY+=("$comp")
            done < <(compgen -W "${allflags[*]}" -- "$cur")
            if [[ $(type -t compopt) = "builtin" ]]; then
                [[ "${COMPREPLY[0]}" == *= ]] || compopt +o nospace
            fi

            # complete after --flag=abc
            if [[ $cur == *=* ]]; then
                if [[ $(type -t compopt) = "builtin" ]]; then
                    compopt +o nospace
                fi

                local index flag
                flag="${cur%=*}"
                __dpfctl_index_of_word "${flag}" "${flags_with_completion[@]}"
                COMPREPLY=()
                if [[ ${index} -ge 0 ]]; then
                    PREFIX=""
                    cur="${cur#*=}"
                    ${flags_completion[${index}]}
                    if [ -n "${ZSH_VERSION}" ]; then
                        # zsh completion needs --flag= prefix
                        eval "COMPREPLY=( \"\${COMPREPLY[@]/#/${flag}=}\" )"
                    fi
                fi
            fi
            return 0;
            ;;
    esac

    # check if we are handling a flag with special work handling
    local index
    __dpfctl_index_of_word "${prev}" "${flags_with_completion[@]}"
    if [[ ${index} -ge 0 ]]; then
        ${flags_completion[${index}]}
        return
    fi

    # we are parsing a flag and don't have a special handler, no completion
    if [[ ${cur} != "${words[cword]}" ]]; then
        return
    fi

    local completions
    completions=("${commands[@]}")
    if [[ ${#must_have_one_noun[@]} -ne 0 ]]; then
        completions+=("${must_have_one_noun[@]}")
    elif [[ -n "${has_completion_function}" ]]; then
        # if a go completion function is provided, defer to that function
        __dpfctl_handle_go_custom_completion
    fi
    if [[ ${#must_have_one_flag[@]} -ne 0 ]]; then
        completions+=("${must_have_one_flag[@]}")
    fi
    while IFS='' read -r comp; do
        COMPREPLY+=("$comp")
    done < <(compgen -W "${completions[*]}" -- "$cur")

    if [[ ${#COMPREPLY[@]} -eq 0 && ${#noun_aliases[@]} -gt 0 && ${#must_have_one_noun[@]} -ne 0 ]]; then
        while IFS='' read -r comp; do
            COMPREPLY+=("$comp")
        done < <(compgen -W "${noun_aliases[*]}" -- "$cur")
    fi

    if [[ ${#COMPREPLY[@]} -eq 0 ]]; then
		if declare -F __dpfctl_custom_func >/dev/null; then
			# try command name qualified custom func
			__dpfctl_custom_func
		else
			# otherwise fall back to unqualified for compatibility
			declare -F __custom_func >/dev/null && __custom_func
		fi
    fi

    # available in bash-completion >= 2, not always present on macOS
    if declare -F __ltrim_colon_completions >/dev/null; then
        __ltrim_colon_completions "$cur"
    fi

    # If there is only 1 completion and it is a flag with an = it will be completed
    # but we don't want a space after the =
    if [[ "${#COMPREPLY[@]}" -eq "1" ]] && [[ $(type -t compopt) = "builtin" ]] && [[ "${COMPREPLY[0]}" == --*= ]]; then
       compopt -o nospace
    fi
}

# The arguments should be in the form "ext1|ext2|extn"
__dpfctl_handle_filename_extension_flag()
{
    local ext="$1"
    _filedir "@(${ext})"
}

__dpfctl_handle_subdirs_in_dir_flag()
{
    local dir="$1"
    pushd "${dir}" >/dev/null 2>&1 && _filedir -d && popd >/dev/null 2>&1 || return
}

__dpfctl_handle_flag()
{
    __dpfctl_debug "${FUNCNAME[0]}: c is $c words[c] is ${words[c]}"

    # if a command required a flag, and we found it, unset must_have_one_flag()
    local flagname=${words[c]}
    local flagvalue
    # if the word contained an =
    if [[ ${words[c]} == *"="* ]]; then
        flagvalue=${flagname#*=} # take in as flagvalue after the =
        flagname=${flagname%=*} # strip everything after the =
        flagname="${flagname}=" # but put the = back
    fi
    __dpfctl_debug "${FUNCNAME[0]}: looking for ${flagname}"
    if __dpfctl_contains_word "${flagname}" "${must_have_one_flag[@]}"; then
        must_have_one_flag=()
    fi

    # if you set a flag which only applies to this command, don't show subcommands
    if __dpfctl_contains_word "${flagname}" "${local_nonpersistent_flags[@]}"; then
      commands=()
    fi

    # keep flag value with flagname as flaghash
    # flaghash variable is an associative array which is only supported in bash > 3.
    if [[ -z "${BASH_VERSION}" || "${BASH_VERSINFO[0]}" -gt 3 ]]; then
        if [ -n "${flagvalue}" ] ; then
            flaghash[${flagname}]=${flagvalue}
        elif [ -n "${words[ $((c+1)) ]}" ] ; then
            flaghash[${flagname}]=${words[ $((c+1)) ]}
        else
            flaghash[${flagname}]="true" # pad "true" for bool flag
        fi
    fi

    # skip the argument to a two word flag
    if [[ ${words[c]} != *"="* ]] && __dpfctl_contains_word "${words[c]}" "${two_word_flags[@]}"; then
			  __dpfctl_debug "${FUNCNAME[0]}: found a flag ${words[c]}, skip the next argument"
        c=$((c+1))
        # if we are looking for a flags value, don't show commands
        if [[ $c -eq $cword ]]; then
            commands=()
        fi
    fi

    c=$((c+1))

}

__dpfctl_handle_noun()
{
    __dpfctl_debug "${FUNCNAME[0]}: c is $c words[c] is ${words[c]}"

    if __dpfctl_contains_word "${words[c]}" "${must_have_one_noun[@]}"; then
        must_have_one_noun=()
    elif __dpfctl_contains_word "${words[c]}" "${noun_aliases[@]}"; then
        must_have_one_noun=()
    fi

    nouns+=("${words[c]}")
    c=$((c+1))
}

__dpfctl_handle_command()
{
    __dpfctl_debug "${FUNCNAME[0]}: c is $c words[c] is ${words[c]}"

    local next_command
    if [[ -n ${last_command} ]]; then
        next_command="_${last_command}_${words[c]//:/__}"
    else
        if [[ $c -eq 0 ]]; then
            next_command="_dpfctl_root_command"
        else
            next_command="_${words[c]//:/__}"
        fi
    fi
    c=$((c+1))
    __dpfctl_debug "${FUNCNAME[0]}: looking for ${next_command}"
    declare -F "$next_command" >/dev/null && $next_command
}

__dpfctl_handle_word()
{
    if [[ $c -ge $cword ]]; then
        __dpfctl_handle_reply
        return
    fi
    __dpfctl_debug "${FUNCNAME[0]}: c is $c words[c] is ${words[c]}"
    if [[ "${words[c]}" == -* ]]; then
        __dpfctl_handle_flag
    elif __dpfctl_contains_word "${words[c]}" "${commands[@]}"; then
        __dpfctl_handle_command
    elif [[ $c -eq 0 ]]; then
        __dpfctl_handle_command
    elif __dpfctl_contains_word "${words[c]}" "${command_aliases[@]}"; then
        # aliashash variable is an associative array which is only supported in bash > 3.
        if [[ -z "${BASH_VERSION}" || "${BASH_VERSINFO[0]}" -gt 3 ]]; then
            words[c]=${aliashash[${words[c]}]}
            __dpfctl_handle_command
        else
            __dpfctl_handle_noun
        fi
    else
        __dpfctl_handle_noun
    fi
    __dpfctl_handle_word
}

_dpfctl_apply()
{
    last_command="dpfctl_apply"

    command_aliases=()

    commands=()

    flags=()
    two_word_flags=()
    local_nonpersistent_flags=()
    flags_with_completion=()
    flags_completion=()

    flags+=("--dry-run")
    flags+=("--filename=")
    two_word_flags+=("--filename")
    flags_with_completion+=("--filename")
    flags_completion+=("__dpfctl_handle_filename_extension_flag yaml|yml|json")
    flags_with_completion+=("--filename")
    flags_completion+=("__dpfctl_handle_go_custom_completion")
    two_word_flags+=("-f")
    flags_with_completion+=("-f")
    flags_completion+=("__dpfctl_handle_filename_extension_flag yaml|yml|json")
    flags_with_completion+=("-f")
    flags_completion+=("__dpfctl_handle_go_custom_completion")
    flags+=("--request-timeout=")
    two_word_flags+=("--request-timeout")
    flags+=("--wait")
    flags+=("--wait-timeout=")
    two_word_flags+=("--wait-timeout")
    flags+=("--config=")
    two_word_flags+=("--config")
    flags+=("--context=")
    two_word_flags+=("--context")
    flags+=("--debug=")
    two_word_flags+=("--debug")
    flags+=("--no-headers")
    flags+=("--output=")
    two_word_flags+=("--output")
    flags_with_completion+=("--output")
    flags_completion+=("__dpfctl_handle_go_custom_completion")
    two_word_flags+=("-o")
    flags_with_completion+=("-o")
    flags_completion+=("__dpfctl_handle_go_custom_completion")

    must_have_one_flag=()
    must_have_one_noun=()
    noun_aliases=()
}

_dpfctl_cancel()
{
    last_command="dpfctl_cancel"

    command_aliases=()

    commands=()

    flags=()
    two_word_flags=()
    local_nonpersistent_flags=()
    flags_with_completion=()
    flags_completion=()

    flags+=("--dry-run")
    flags+=("--filename=")
    two_word_flags+=("--filename")
    flags_with_completion+=("--filename")
    flags_completion+=("__dpfctl_handle_filename_extension_flag yaml|yml|json")
    flags_with_completion+=("--filename")
    flags_completion+=("__dpfctl_handle_go_custom_completion")
    two_word_flags+=("-f")
    flags_with_completion+=("-f")
    flags_completion+=("__dpfctl_handle_filename_extension_flag yaml|yml|json")
    flags_with_completion+=("-f")
    flags_completion+=("__dpfctl_handle_go_custom_completion")
    flags+=("--request-timeout=")
    two_word_flags+=("--request-timeout")
    flags+=("--wait")
    flags+=("--wait-timeout=")
    two_word_flags+=("--wait-timeout")
    flags+=("--config=")
    two_word_flags+=("--config")
    flags+=("--context=")
    two_word_flags+=("--context")
    flags+=("--debug=")
    two_word_flags+=("--debug")
    flags+=("--no-headers")
    flags+=("--output=")
    two_word_flags+=("--output")
    flags_with_completion+=("--output")
    flags_completion+=("__dpfctl_handle_go_custom_completion")
    two_word_flags+=("-o")
    flags_with_completion+=("-o")
    flags_completion+=("__dpfctl_handle_go_custom_completion")

    must_have_one_flag=()
    must_have_one_noun=()
    noun_aliases=()
}

_dpfctl_completion()
{
    last_command="dpfctl_completion"

    command_aliases=()

    commands=()

    flags=()
    two_word_flags=()
    local_nonpersistent_flags=()
    flags_with_completion=()
    flags_completion=()

    flags+=("--help")
    flags+=("-h")
    local_nonpersistent_flags+=("--help")
    local_nonpersistent_flags+=("-h")
    flags+=("--config=")
    two_word_flags+=("--config")
    flags+=("--context=")
    two_word_flags+=("--context")
    flags+=("--debug=")
    two_word_flags+=("--debug")
    flags+=("--no-headers")
    flags+=("--output=")
    two_word_flags+=("--output")
    flags_with_completion+=("--output")
    flags_completion+=("__dpfctl_handle_go_custom_completion")
    two_word_flags+=("-o")
    flags_with_completion+=("-o")
    flags_completion+=("__dpfctl_handle_go_custom_completion")

    must_have_one_flag=()
    must_have_one_noun=()
    must_have_one_noun+=("bash")
    must_have_one_noun+=("fish")
    must_have_one_noun+=("powershell")
    must_have_one_noun+=("zsh")
    noun_aliases=()
}

_dpfctl_config_get-context()
{
    last_command="dpfctl_config_get-context"

    command_aliases=()

    commands=()

    flags=()
    two_word_flags=()
    local_nonpersistent_flags=()
    flags_with_completion=()
    flags_completion=()

    flags+=("--config=")
    two_word_flags+=("--config")
    flags+=("--context=")
    two_word_flags+=("--context")
    flags+=("--debug=")
    two_word_flags+=("--debug")
    flags+=("--no-headers")
    flags+=("--output=")
    two_word_flags+=("--output")
    flags_with_completion+=("--output")
    flags_completion+=("__dpfctl_handle_go_custom_completion")
    two_word_flags+=("-o")
    flags_with_completion+=("-o")
    flags_completion+=("__dpfctl_handle_go_custom_completion")

    must_have_one_flag=()
    must_have_one_noun=()
    noun_aliases=()
}

_dpfctl_config_get-current-context()
{
    last_command="dpfctl_config_get-current-context"

    command_aliases=()

    commands=()

    flags=()
    two_word_flags=()
    local_nonpersistent_flags=()
    flags_with_completion=()
    flags_completion=()

    flags+=("--config=")
    two_word_flags+=("--config")
    flags+=("--context=")
    two_word_flags+=("--context")
    flags+=("--debug=")
    two_word_flags+=("--debug")
    flags+=("--no-headers")
    flags+=("--output=")
    two_word_flags+=("--output")
    flags_with_completion+=("--output")
    flags_completion+=("__dpfctl_handle_go_custom_completion")
    two_word_flags+=("-o")
    flags_with_completion+=("-o")
    flags_completion+=("__dpfctl_handle_go_custom_completion")

    must_have_one_flag=()
    must_have_one_noun=()
    noun_aliases=()
}

_dpfctl_config_set-context()
{
    last_command="dpfctl_config_set-context"

    command_aliases=()

    commands=()

    flags=()
    two_word_flags=()
    local_nonpersistent_flags=()
    flags_with_completion=()
    flags_completion=()

    flags+=("--config=")
    two_word_flags+=("--config")
    flags+=("--context=")
    two_word_flags+=("--context")
    flags+=("--debug=")
    two_word_flags+=("--debug")
    flags+=("--no-headers")
    flags+=("--output=")
    two_word_flags+=("--output")
    flags_with_completion+=("--output")
    flags_completion+=("__dpfctl_handle_go_custom_completion")
    two_word_flags+=("-o")
    flags_with_completion+=("-o")
    flags_completion+=("__dpfctl_handle_go_custom_completion")

    must_have_one_flag=()
    must_have_one_noun=()
    noun_aliases=()
}

_dpfctl_config_set-current-context()
{
    last_command="dpfctl_config_set-current-context"

    command_aliases=()

    commands=()

    flags=()
    two_word_flags=()
    local_nonpersistent_flags=()
    flags_with_completion=()
    flags_completion=()

    flags+=("--config=")
    two_word_flags+=("--config")
    flags+=("--context=")
    two_word_flags+=("--context")
    flags+=("--debug=")
    two_word_flags+=("--debug")
    flags+=("--no-headers")
    flags+=("--output=")
    two_word_flags+=("--output")
    flags_with_completion+=("--output")
    flags_completion+=("__dpfctl_handle_go_custom_completion")
    two_word_flags+=("-o")
    flags_with_completion+=("-o")
    flags_completion+=("__dpfctl_handle_go_custom_completion")

    must_have_one_flag=()
    must_have_one_noun=()
    noun_aliases=()
}

_dpfctl_config()
{
    last_command="dpfctl_config"

    command_aliases=()

    commands=()
    commands+=("get-context")
    commands+=("get-current-context")
    commands+=("set-context")
    commands+=("set-current-context")

    flags=()
    two_word_flags=()
    local_nonpersistent_flags=()
    flags_with_completion=()
    flags_completion=()

    flags+=("--config=")
    two_word_flags+=("--config")
    flags+=("--context=")
    two_word_flags+=("--context")
    flags+=("--debug=")
    two_word_flags+=("--debug")
    flags+=("--no-headers")
    flags+=("--output=")
    two_word_flags+=("--output")
    flags_with_completion+=("--output")
    flags_completion+=("__dpfctl_handle_go_custom_completion")
    two_word_flags+=("-o")
    flags_with_completion+=("-o")
    flags_completion+=("__dpfctl_handle_go_custom_completion")

    must_have_one_flag=()
    must_have_one_noun=()
    noun_aliases=()
}

_dpfctl_create()
{
    last_command="dpfctl_create"

    command_aliases=()

    commands=()

    flags=()
    two_word_flags=()
    local_nonpersistent_flags=()
    flags_with_completion=()
    flags_completion=()

    flags+=("--dry-run")
    flags+=("--filename=")
    two_word_flags+=("--filename")
    flags_with_completion+=("--filename")
    flags_completion+=("__dpfctl_handle_filename_extension_flag yaml|yml|json")
    flags_with_completion+=("--filename")
    flags_completion+=("__dpfctl_handle_go_custom_completion")
    two_word_flags+=("-f")
    flags_with_completion+=("-f")
    flags_completion+=("__dpfctl_handle_filename_extension_flag yaml|yml|json")
    flags_with_completion+=("-f")
    flags_completion+=("__dpfctl_handle_go_custom_completion")
    flags+=("--request-timeout=")
    two_word_flags+=("--request-timeout")
    flags+=("--wait")
    flags+=("--wait-timeout=")
    two_word_flags+=("--wait-timeout")
    flags+=("--config=")
    two_word_flags+=("--config")
    flags+=("--context=")
    two_word_flags+=("--context")
    flags+=("--debug=")
    two_word_flags+=("--debug")
    flags+=("--no-headers")
    flags+=("--output=")
    two_word_flags+=("--output")
    flags_with_completion+=("--output")
    flags_completion+=("__dpfctl_handle_go_custom_completion")
    two_word_flags+=("-o")
    flags_with_completion+=("-o")
    flags_completion+=("__dpfctl_handle_go_custom_completion")

    must_have_one_flag=()
    must_have_one_noun=()
    noun_aliases=()
}

_dpfctl_delete()
{
    last_command="dpfctl_delete"

    command_aliases=()

    commands=()

    flags=()
    two_word_flags=()
    local_nonpersistent_flags=()
    flags_with_completion=()
    flags_completion=()

    flags+=("--dry-run")
    flags+=("--filename=")
    two_word_flags+=("--filename")
    flags_with_completion+=("--filename")
    flags_completion+=("__dpfctl_handle_filename_extension_flag yaml|yml|json")
    flags_with_completion+=("--filename")
    flags_completion+=("__dpfctl_handle_go_custom_completion")
    two_word_flags+=("-f")
    flags_with_completion+=("-f")
    flags_completion+=("__dpfctl_handle_filename_extension_flag yaml|yml|json")
    flags_with_completion+=("-f")
    flags_completion+=("__dpfctl_handle_go_custom_completion")
    flags+=("--request-timeout=")
    two_word_flags+=("--request-timeout")
    flags+=("--wait")
    flags+=("--wait-timeout=")
    two_word_flags+=("--wait-timeout")
    flags+=("--config=")
    two_word_flags+=("--config")
    flags+=("--context=")
    two_word_flags+=("--context")
    flags+=("--debug=")
    two_word_flags+=("--debug")
    flags+=("--no-headers")
    flags+=("--output=")
    two_word_flags+=("--output")
    flags_with_completion+=("--output")
    flags_completion+=("__dpfctl_handle_go_custom_completion")
    two_word_flags+=("-o")
    flags_with_completion+=("-o")
    flags_completion+=("__dpfctl_handle_go_custom_completion")

    must_have_one_flag=()
    must_have_one_noun=()
    has_completion_function=1
    noun_aliases=()
}

_dpfctl_get()
{
    last_command="dpfctl_get"

    command_aliases=()

    commands=()

    flags=()
    two_word_flags=()
    local_nonpersistent_flags=()
    flags_with_completion=()
    flags_completion=()

    flags+=("--row-search-params=")
    two_word_flags+=("--row-search-params")
    flags+=("--config=")
    two_word_flags+=("--config")
    flags+=("--context=")
    two_word_flags+=("--context")
    flags+=("--debug=")
    two_word_flags+=("--debug")
    flags+=("--no-headers")
    flags+=("--output=")
    two_word_flags+=("--output")
    flags_with_completion+=("--output")
    flags_completion+=("__dpfctl_handle_go_custom_completion")
    two_word_flags+=("-o")
    flags_with_completion+=("-o")
    flags_completion+=("__dpfctl_handle_go_custom_completion")

    must_have_one_flag=()
    must_have_one_noun=()
    has_completion_function=1
    noun_aliases=()
}

_dpfctl_help()
{
    last_command="dpfctl_help"

    command_aliases=()

    commands=()

    flags=()
    two_word_flags=()
    local_nonpersistent_flags=()
    flags_with_completion=()
    flags_completion=()

    flags+=("--config=")
    two_word_flags+=("--config")
    flags+=("--context=")
    two_word_flags+=("--context")
    flags+=("--debug=")
    two_word_flags+=("--debug")
    flags+=("--no-headers")
    flags+=("--output=")
    two_word_flags+=("--output")
    flags_with_completion+=("--output")
    flags_completion+=("__dpfctl_handle_go_custom_completion")
    two_word_flags+=("-o")
    flags_with_completion+=("-o")
    flags_completion+=("__dpfctl_handle_go_custom_completion")

    must_have_one_flag=()
    must_have_one_noun=()
    has_completion_function=1
    noun_aliases=()
}

_dpfctl_update()
{
    last_command="dpfctl_update"

    command_aliases=()

    commands=()

    flags=()
    two_word_flags=()
    local_nonpersistent_flags=()
    flags_with_completion=()
    flags_completion=()

    flags+=("--dry-run")
    flags+=("--filename=")
    two_word_flags+=("--filename")
    flags_with_completion+=("--filename")
    flags_completion+=("__dpfctl_handle_go_custom_completion")
    flags_with_completion+=("--filename")
    flags_completion+=("__dpfctl_handle_filename_extension_flag yaml|yml|json")
    two_word_flags+=("-f")
    flags_with_completion+=("-f")
    flags_completion+=("__dpfctl_handle_filename_extension_flag yaml|yml|json")
    flags_with_completion+=("-f")
    flags_completion+=("__dpfctl_handle_go_custom_completion")
    flags+=("--request-timeout=")
    two_word_flags+=("--request-timeout")
    flags+=("--wait")
    flags+=("--wait-timeout=")
    two_word_flags+=("--wait-timeout")
    flags+=("--config=")
    two_word_flags+=("--config")
    flags+=("--context=")
    two_word_flags+=("--context")
    flags+=("--debug=")
    two_word_flags+=("--debug")
    flags+=("--no-headers")
    flags+=("--output=")
    two_word_flags+=("--output")
    flags_with_completion+=("--output")
    flags_completion+=("__dpfctl_handle_go_custom_completion")
    two_word_flags+=("-o")
    flags_with_completion+=("-o")
    flags_completion+=("__dpfctl_handle_go_custom_completion")

    must_have_one_flag=()
    must_have_one_noun=()
    noun_aliases=()
}

_dpfctl_update()
{
    last_command="dpfctl_update"

    command_aliases=()

    commands=()

    flags=()
    two_word_flags=()
    local_nonpersistent_flags=()
    flags_with_completion=()
    flags_completion=()

    flags+=("--dry-run")
    flags+=("--filename=")
    two_word_flags+=("--filename")
    flags_with_completion+=("--filename")
    flags_completion+=("__dpfctl_handle_filename_extension_flag yaml|yml|json")
    flags_with_completion+=("--filename")
    flags_completion+=("__dpfctl_handle_go_custom_completion")
    two_word_flags+=("-f")
    flags_with_completion+=("-f")
    flags_completion+=("__dpfctl_handle_filename_extension_flag yaml|yml|json")
    flags_with_completion+=("-f")
    flags_completion+=("__dpfctl_handle_go_custom_completion")
    flags+=("--request-timeout=")
    two_word_flags+=("--request-timeout")
    flags+=("--wait")
    flags+=("--wait-timeout=")
    two_word_flags+=("--wait-timeout")
    flags+=("--config=")
    two_word_flags+=("--config")
    flags+=("--context=")
    two_word_flags+=("--context")
    flags+=("--debug=")
    two_word_flags+=("--debug")
    flags+=("--no-headers")
    flags+=("--output=")
    two_word_flags+=("--output")
    flags_with_completion+=("--output")
    flags_completion+=("__dpfctl_handle_go_custom_completion")
    two_word_flags+=("-o")
    flags_with_completion+=("-o")
    flags_completion+=("__dpfctl_handle_go_custom_completion")

    must_have_one_flag=()
    must_have_one_noun=()
    noun_aliases=()
}

_dpfctl_root_command()
{
    last_command="dpfctl"

    command_aliases=()

    commands=()
    commands+=("apply")
    commands+=("cancel")
    commands+=("completion")
    commands+=("config")
    commands+=("create")
    commands+=("delete")
    commands+=("get")
    commands+=("help")
    commands+=("update")
    commands+=("update")

    flags=()
    two_word_flags=()
    local_nonpersistent_flags=()
    flags_with_completion=()
    flags_completion=()

    flags+=("--config=")
    two_word_flags+=("--config")
    flags+=("--context=")
    two_word_flags+=("--context")
    flags+=("--debug=")
    two_word_flags+=("--debug")
    flags+=("--no-headers")
    flags+=("--output=")
    two_word_flags+=("--output")
    flags_with_completion+=("--output")
    flags_completion+=("__dpfctl_handle_go_custom_completion")
    two_word_flags+=("-o")
    flags_with_completion+=("-o")
    flags_completion+=("__dpfctl_handle_go_custom_completion")

    must_have_one_flag=()
    must_have_one_noun=()
    noun_aliases=()
}

__start_dpfctl()
{
    local cur prev words cword split
    declare -A flaghash 2>/dev/null || :
    declare -A aliashash 2>/dev/null || :
    if declare -F _init_completion >/dev/null 2>&1; then
        _init_completion -s || return
    else
        __dpfctl_init_completion -n "=" || return
    fi

    local c=0
    local flags=()
    local two_word_flags=()
    local local_nonpersistent_flags=()
    local flags_with_completion=()
    local flags_completion=()
    local commands=("dpfctl")
    local command_aliases=()
    local must_have_one_flag=()
    local must_have_one_noun=()
    local has_completion_function
    local last_command
    local nouns=()
    local noun_aliases=()

    __dpfctl_handle_word
}

if [[ $(type -t compopt) = "builtin" ]]; then
    complete -o default -F __start_dpfctl dpfctl
else
    complete -o default -o nospace -F __start_dpfctl dpfctl
fi

# ex: ts=4 sw=4 et filetype=sh
