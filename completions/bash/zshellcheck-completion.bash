# bash completion for zshellcheck                          -*- shell-script -*-

_zshellcheck()
{
    local cur prev words cword
    _init_completion || return

    case $prev in
        -format)
            COMPREPLY=( $(compgen -W "text json" -- "$cur") )
            return
            ;;
    esac

    if [[ "$cur" == -* ]]; then
        COMPREPLY=( $(compgen -W "-format -help" -- "$cur") )
        return
    fi

    _filedir '@(zsh|sh|zsh-theme)'
} &&
complete -F _zshellcheck zshellcheck
