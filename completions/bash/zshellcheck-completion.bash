# bash completion for zshellcheck                          -*- shell-script -*-

_zshellcheck()
{
    local cur prev words cword
    _init_completion || return

    case $prev in
        -format)
            COMPREPLY=( $(compgen -W "text json sarif" -- "$cur") )
            return
            ;;
        -severity)
            COMPREPLY=( $(compgen -W "error warning info style" -- "$cur") )
            return
            ;;
        -cpuprofile)
            _filedir
            return
            ;;
    esac

    if [[ "$cur" == -* ]]; then
        COMPREPLY=( $(compgen -W \
            "-format -severity --no-color --verbose -cpuprofile -version -h --help" \
            -- "$cur") )
        return
    fi

    _filedir '@(zsh|sh|zsh-theme)'
} &&
complete -F _zshellcheck zshellcheck
