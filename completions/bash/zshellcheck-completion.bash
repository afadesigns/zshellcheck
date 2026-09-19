# bash completion for zshellcheck                          -*- shell-script -*-

_zshellcheck()
{
    local cur prev words cword option='' prefix='' word i
    _init_completion -n = || return

    # Go flags stop at the first target or --. Consume string values first,
    # even when a value starts with a hyphen.
    for ((i = 1; i < cword; i++)); do
        word=${words[i]}
        if [[ -n $option ]]; then
            option=''
            continue
        fi
        case $word in
            --|-|[!-]*)
                _filedir '@(zsh|sh|zsh-theme)'
                return
                ;;
            -format|--format|-severity|--severity|-cpuprofile|--cpuprofile|\
            -baseline|--baseline|-baseline-write|--baseline-write|\
            -explain|--explain|-rule-severity|--rule-severity)
                option=$word
                ;;
        esac
    done

    if [[ -z $option && $cur == -*=* ]]; then
        option=${cur%%=*}
        # Bash's replacement word excludes text before an open quote or
        # a Readline word break, even when _init_completion rejoins it.
        if [[ ${2-} == "$option="* ]]; then
            prefix=$option=
        fi
        cur=${cur#*=}
    fi
    case $option in
        -format|--format)
            COMPREPLY=( $(compgen -W 'text json sarif' -- "$cur") )
            ;;
        -severity|--severity)
            if [[ $cur == *,* ]]; then
                prefix+=${cur%,*},
                cur=${cur##*,}
            fi
            COMPREPLY=( $(compgen -W 'error warning info style' -- "$cur") )
            ;;
        -cpuprofile|--cpuprofile|-baseline|--baseline|-baseline-write|--baseline-write)
            _filedir
            ;;
        -explain|--explain|-rule-severity|--rule-severity)
            COMPREPLY=()
            ;;
        '')
            if [[ $cur == -* ]]; then
                local flags='format severity no-color no-banner verbose cpuprofile
                    fix diff dry-run unsafe-fixes list-rules explain statistics
                    baseline baseline-write rule-severity add-noka detect-stale-noka
                    completions version h help'
                local options=''
                for word in $flags; do
                    options+=" -$word --$word"
                done
                COMPREPLY=( $(compgen -W "$options" -- "$cur") )
            else
                _filedir '@(zsh|sh|zsh-theme)'
            fi
            ;;
        *) COMPREPLY=() ;;
    esac
    if [[ -n $prefix ]]; then
        for i in "${!COMPREPLY[@]}"; do
            COMPREPLY[i]=$prefix${COMPREPLY[i]}
        done
    fi
} &&
complete -F _zshellcheck zshellcheck
