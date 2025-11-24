# ZC1000-ZC1099: Foundational Zsh Checks

This section provides detailed documentation for Katas ranging from ZC1001 to ZC1099, focusing on foundational Zsh best practices, common pitfalls, and efficiency improvements.

Click on any Kata ID below for a comprehensive explanation, including bad and good code examples, and configuration options.

*   [[Katas/ZC1000-ZC1099/ZC1001 | ZC1001]]: `Use ${} for array element access`
*   [[Katas/ZC1000-ZC1099/ZC1002 | ZC1002]]: `Use $(...) instead of backticks`
*   [[Katas/ZC1000-ZC1099/ZC1003 | ZC1003]]: `Use ((...))` for arithmetic comparisons instead of `[` or `test`
*   [[Katas/ZC1000-ZC1099/ZC1005 | ZC1005]]: `Use whence instead of which`
*   [[Katas/ZC1000-ZC1099/ZC1006 | ZC1006]]: `Prefer [[ over test for tests`
*   [[Katas/ZC1000-ZC1099/ZC1007 | ZC1007]]: `Avoid using chmod 777`
*   [[Katas/ZC1000-ZC1099/ZC1008 | ZC1008]]: `Use $(())` for arithmetic operations`
*   [[Katas/ZC1000-ZC1099/ZC1009 | ZC1009]]: `Use ((...)) for C-style arithmetic`
*   [[Katas/ZC1000-ZC1099/ZC1010 | ZC1010]]: `Use [[ ... ]] instead of [ ... ]`
*   [[Katas/ZC1000-ZC1099/ZC1011 | ZC1011]]: `Use git porcelain commands instead of plumbing commands`
*   [[Katas/ZC1000-ZC1099/ZC1012 | ZC1012]]: `Use read -r to prevent backslash escaping`
*   [[Katas/ZC1000-ZC1099/ZC1013 | ZC1013]]: `Use ((...))` for arithmetic operations instead of let`
*   [[Katas/ZC1000-ZC1099/ZC1014 | ZC1014]]: `Use git switch or git restore instead of git checkout`
*   [[Katas/ZC1000-ZC1099/ZC1017 | ZC1017]]: `Use print -r to print strings literally`
*   [[Katas/ZC1000-ZC1099/ZC1018 | ZC1018]]: `Use ((...)) for C-style arithmetic instead of expr`
*   [[Katas/ZC1000-ZC1099/ZC1020 | ZC1020]]: `Use [[ ... ]] for tests instead of test`
*   [[Katas/ZC1000-ZC1099/ZC1021 | ZC1021]]: `Use symbolic permissions with chmod instead of octal`
*   [[Katas/ZC1000-ZC1099/ZC1022 | ZC1022]]: `Use $((...))` for arithmetic expansion`
*   [[Katas/ZC1000-ZC1099/ZC1023 | ZC1023]]: `Use $((...))` for arithmetic expansion`
*   [[Katas/ZC1000-ZC1099/ZC1024 | ZC1024]]: `Use $((...))` for arithmetic expansion`
*   [[Katas/ZC1000-ZC1099/ZC1025 | ZC1025]]: `Use $((...))` for arithmetic expansion`
*   [[Katas/ZC1000-ZC1099/ZC1026 | ZC1026]]: `Use $((...))` for arithmetic expansion`
*   [[Katas/ZC1000-ZC1099/ZC1027 | ZC1027]]: `Use $((...))` for arithmetic expansion`
*   [[Katas/ZC1000-ZC1099/ZC1028 | ZC1028]]: `Use $((...))` for arithmetic expansion`
*   [[Katas/ZC1000-ZC1099/ZC1029 | ZC1029]]: `Use $((...))` for arithmetic expansion`
*   [[Katas/ZC1000-ZC1099/ZC1030 | ZC1030]]: `Use printf instead of echo`
*   [[Katas/ZC1000-ZC1099/ZC1031 | ZC1031]]: `Use #!/usr/bin/env zsh for portability`
*   [[Katas/ZC1000-ZC1099/ZC1032 | ZC1032]]: `Use ((...)) for C-style incrementing`
*   [[Katas/ZC1000-ZC1099/ZC1033 | ZC1033]]: `Use $((...))` for arithmetic expansion`
*   [[Katas/ZC1000-ZC1099/ZC1034 | ZC1034]]: `Use command -v instead of which`
*   [[Katas/ZC1000-ZC1099/ZC1035 | ZC1035]]: `Use $((...))` for arithmetic expansion`
*   [[Katas/ZC1000-ZC1099/ZC1036 | ZC1036]]: `Prefer [[ ... ]] over test command`
*   [[Katas/ZC1000-ZC1099/ZC1037 | ZC1037]]: `Use 'print -r --' for variable expansion`
*   [[Katas/ZC1000-ZC1099/ZC1038 | ZC1038]]: `Avoid useless use of cat`
*   [[Katas/ZC1000-ZC1099/ZC1039 | ZC1039]]: `Avoid rm with root path`
*   [[Katas/ZC1000-ZC1099/ZC1040 | ZC1040]]: `Use (N) nullglob qualifier for globs in loops`
*   [[Katas/ZC1000-ZC1099/ZC1041 | ZC1041]]: `Do not use variables in printf format string`
*   [[Katas/ZC1000-ZC1099/ZC1042 | ZC1042]]: `Use "$@" to iterate over arguments`
*   [[Katas/ZC1000-ZC1099/ZC1043 | ZC1043]]: `Use local for variables in functions`
*   [[Katas/ZC1000-ZC1099/ZC1044 | ZC1044]]: `Check for unchecked cd commands`
*   [[Katas/ZC1000-ZC1099/ZC1045 | ZC1045]]: `Declare and assign separately to avoid masking return values`
*   [[Katas/ZC1000-ZC1099/ZC1046 | ZC1046]]: `Avoid eval`
*   [[Katas/ZC1000-ZC1099/ZC1047 | ZC1047]]: `Avoid sudo in scripts`
*   [[Katas/ZC1000-ZC1099/ZC1048 | ZC1048]]: `Avoid source with relative paths`
*   [[Katas/ZC1000-ZC1099/ZC1049 | ZC1049]]: `Prefer functions over aliases`
*   [[Katas/ZC1000-ZC1099/ZC1050 | ZC1050]]: `Avoid iterating over ls output`
*   [[Katas/ZC1000-ZC1099/ZC1051 | ZC1051]]: `Quote variables in rm to avoid globbing`
*   [[Katas/ZC1000-ZC1099/ZC1052 | ZC1052]]: `Avoid sed -i for portability`
*   [[Katas/ZC1000-ZC1099/ZC1053 | ZC1053]]: `Silence grep output in conditions`
*   [[Katas/ZC1000-ZC1099/ZC1054 | ZC1054]]: `Use POSIX classes in regex/glob`
*   [[Katas/ZC1000-ZC1099/ZC1055 | ZC1055]]: `Use [[ -n/-z ]] for empty string checks`
*   [[Katas/ZC1000-ZC1099/ZC1056 | ZC1056]]: `Avoid $((...))` as a statement`
*   [[Katas/ZC1000-ZC1099/ZC1057 | ZC1057]]: `Avoid ls in assignments`
*   [[Katas/ZC1000-ZC1099/ZC1058 | ZC1058]]: `Avoid sudo with redirection`
*   [[Katas/ZC1000-ZC1099/ZC1059 | ZC1059]]: `Use ${var:?}` for `rm` arguments`
*   [[Katas/ZC1000-ZC1099/ZC1060 | ZC1060]]: `Avoid ps | grep without exclusion`
*   [[Katas/ZC1000-ZC1099/ZC1061 | ZC1061]]: `Prefer {start..end} over seq`
*   [[Katas/ZC1000-ZC1099/ZC1062 | ZC1062]]: `Prefer grep -E over egrep`
*   [[Katas/ZC1000-ZC1099/ZC1063 | ZC1063]]: `Prefer grep -F over fgrep`
*   [[Katas/ZC1000-ZC1099/ZC1064 | ZC1064]]: `Prefer command -v over type`
*   [[Katas/ZC1000-ZC1099/ZC1065 | ZC1065]]: `Ensure spaces around [ and [[`
*   [[Katas/ZC1000-ZC1099/ZC1066 | ZC1066]]: `Avoid iterating over cat output`
*   [[Katas/ZC1000-ZC1099/ZC1067 | ZC1067]]: `Separate export and assignment to avoid masking return codes`
*   [[Katas/ZC1000-ZC1099/ZC1068 | ZC1068]]: `Use add-zsh-hook instead of defining hook functions directly`
*   [[Katas/ZC1000-ZC1099/ZC1069 | ZC1069]]: `Avoid local outside of functions`
*   [[Katas/ZC1000-ZC1099/ZC1070 | ZC1070]]: `Use builtin or command to avoid infinite recursion in wrapper functions`
*   [[Katas/ZC1000-ZC1099/ZC1071 | ZC1071]]: `Use += for appending to arrays`
*   [[Katas/ZC1000-ZC1099/ZC1072 | ZC1072]]: `Use awk instead of grep | awk`
*   [[Katas/ZC1000-ZC1099/ZC1073 | ZC1073]]: `Unnecessary use of $ in arithmetic expressions`
