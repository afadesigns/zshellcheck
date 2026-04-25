# ZShellCheck Katas

Auto-generated list of all 1000 implemented checks. Do not edit by hand — regenerate via `go run ./internal/tools/gen-katas-md`.

## Summary

| Severity | Count |
| :--- | ---: |
| `error` | 220 |
| `warning` | 459 |
| `info` | 64 |
| `style` | 257 |
| **total** | **1000** |
| **with auto-fix** | **126** |

Auto-fix availability is marked per-entry below as **Auto-fix:** `yes` or `no`. Run `zshellcheck -fix path/...` to apply every available rewrite, or `-diff` to preview without writing.

## Table of Contents

- [ZC1001: Use ${} for array element access](#zc1001) · auto-fix
- [ZC1002: Use $(...) instead of backticks](#zc1002) · auto-fix
- [ZC1003: Use `((...))` for arithmetic comparisons instead of `\[` or `test`](#zc1003) · auto-fix
- [ZC1004: Use `return` instead of `exit` in functions](#zc1004) · auto-fix
- [ZC1005: Use whence instead of which](#zc1005) · auto-fix
- [ZC1006: Prefer \[\[ over test for tests](#zc1006)
- [ZC1007: Avoid using `chmod 777`](#zc1007)
- [ZC1008: Use `\$(())` for arithmetic operations](#zc1008)
- [ZC1009: Use `((...))` for C-style arithmetic](#zc1009)
- [ZC1010: Use \[\[ ... \]\] instead of \[ ... \]](#zc1010) · auto-fix
- [ZC1011: Use `git` porcelain commands instead of plumbing commands](#zc1011)
- [ZC1012: Use `read -r` to prevent backslash escaping](#zc1012) · auto-fix
- [ZC1013: Use `((...))` for arithmetic operations instead of `let`](#zc1013) · auto-fix
- [ZC1014: Use `git switch` or `git restore` instead of `git checkout`](#zc1014)
- [ZC1015: Use `$(...)` for command substitution instead of backticks](#zc1015) · auto-fix
- [ZC1016: Use `read -s` when reading sensitive information](#zc1016) · auto-fix
- [ZC1017: Use `print -r` to print strings literally](#zc1017) · auto-fix
- [ZC1018: Superseded by ZC1009 — retired duplicate](#zc1018)
- [ZC1019: Superseded by ZC1005 — retired duplicate](#zc1019)
- [ZC1020: Use `\[\[ ... \]\]` for tests instead of `test`](#zc1020)
- [ZC1021: Use symbolic permissions with `chmod` instead of octal](#zc1021)
- [ZC1022: Use `$((...))` for arithmetic expansion](#zc1022) · auto-fix
- [ZC1023: Superseded by ZC1022 — retired duplicate `let` detector](#zc1023)
- [ZC1024: Superseded by ZC1022 — retired duplicate `let` detector](#zc1024)
- [ZC1025: Superseded by ZC1022 — retired duplicate `let` detector](#zc1025)
- [ZC1026: Superseded by ZC1022 — retired duplicate `let` detector](#zc1026)
- [ZC1027: Superseded by ZC1022 — retired duplicate `let` detector](#zc1027)
- [ZC1028: Superseded by ZC1022 — retired duplicate `let` detector](#zc1028)
- [ZC1029: Superseded by ZC1022 — retired duplicate `let` detector](#zc1029)
- [ZC1030: Use `printf` instead of `echo`](#zc1030)
- [ZC1031: Use `#!/usr/bin/env zsh` for portability](#zc1031) · auto-fix
- [ZC1032: Use `((...))` for C-style incrementing](#zc1032) · auto-fix
- [ZC1033: Superseded by ZC1022 — retired duplicate `let` detector](#zc1033)
- [ZC1034: Use `command -v` instead of `which`](#zc1034) · auto-fix
- [ZC1035: Superseded by ZC1022 — retired duplicate `let` detector](#zc1035)
- [ZC1036: Prefer `\[\[ ... \]\]` over `test` command](#zc1036)
- [ZC1037: Use 'print -r --' for variable expansion](#zc1037)
- [ZC1038: Avoid useless use of cat](#zc1038)
- [ZC1039: Avoid `rm` with root path](#zc1039)
- [ZC1040: Use (N) nullglob qualifier for globs in loops](#zc1040) · auto-fix
- [ZC1041: Do not use variables in printf format string](#zc1041)
- [ZC1042: Use "$@" to iterate over arguments](#zc1042)
- [ZC1043: Use `local` for variables in functions](#zc1043) · auto-fix
- [ZC1044: Check for unchecked `cd` commands](#zc1044)
- [ZC1045: Declare and assign separately to avoid masking return values](#zc1045)
- [ZC1046: Avoid `eval`](#zc1046)
- [ZC1047: Avoid `sudo` in scripts](#zc1047)
- [ZC1048: Avoid `source` with relative paths](#zc1048)
- [ZC1049: Prefer functions over aliases](#zc1049)
- [ZC1050: Avoid iterating over `ls` output](#zc1050)
- [ZC1051: Quote variables in `rm` to avoid globbing](#zc1051) · auto-fix
- [ZC1052: Avoid `sed -i` for portability](#zc1052)
- [ZC1053: Silence `grep` output in conditions](#zc1053) · auto-fix
- [ZC1054: Use POSIX classes in regex/glob](#zc1054)
- [ZC1055: Use `\[\[ -n/-z \]\]` for empty string checks](#zc1055) · auto-fix
- [ZC1056: Avoid `$((...))` as a statement](#zc1056)
- [ZC1057: Avoid `ls` in assignments](#zc1057)
- [ZC1058: Avoid `sudo` with redirection](#zc1058)
- [ZC1059: Use `${var:?}` for `rm` arguments](#zc1059)
- [ZC1060: Avoid `ps \| grep` without exclusion](#zc1060)
- [ZC1061: Prefer `{start..end}` over `seq`](#zc1061) · auto-fix
- [ZC1062: Prefer `grep -E` over `egrep`](#zc1062) · auto-fix
- [ZC1063: Prefer `grep -F` over `fgrep`](#zc1063) · auto-fix
- [ZC1064: Prefer `command -v` over `type`](#zc1064) · auto-fix
- [ZC1065: Ensure spaces around `\[` and `\[\[`](#zc1065)
- [ZC1066: Avoid iterating over `cat` output](#zc1066)
- [ZC1067: Separate `export` and assignment to avoid masking return codes](#zc1067)
- [ZC1068: Use `add-zsh-hook` instead of defining hook functions directly](#zc1068)
- [ZC1069: Avoid `local` outside of functions](#zc1069)
- [ZC1070: Use `builtin` or `command` to avoid infinite recursion in wrapper functions](#zc1070)
- [ZC1071: Use `+=` for appending to arrays](#zc1071)
- [ZC1072: Use `awk` instead of `grep \| awk`](#zc1072)
- [ZC1073: Unnecessary use of `$` in arithmetic expressions](#zc1073) · auto-fix
- [ZC1074: Prefer modifiers :h/:t over dirname/basename](#zc1074)
- [ZC1075: Quote variable expansions to prevent globbing](#zc1075)
- [ZC1076: Use `autoload -Uz` for lazy loading](#zc1076) · auto-fix
- [ZC1077: Prefer `${var:u/l}` over `tr` for case conversion](#zc1077)
- [ZC1078: Quote `$@` and `$*` when passing arguments](#zc1078) · auto-fix
- [ZC1079: Quote RHS of `==` in `\[\[ ... \]\]` to prevent pattern matching](#zc1079) · auto-fix
- [ZC1080: Use `(N)` nullglob qualifier for globs in loops](#zc1080)
- [ZC1081: Use `${#var}` to get string length instead of `wc -c`](#zc1081)
- [ZC1082: Prefer `${var//old/new}` over `sed` for simple replacements](#zc1082)
- [ZC1083: Brace expansion limits cannot be variables](#zc1083)
- [ZC1084: Quote globs in `find` commands](#zc1084) · auto-fix
- [ZC1085: Quote variable expansions in `for` loops](#zc1085) · auto-fix
- [ZC1086: Prefer `func() { ... }` over `function func { ... }`](#zc1086) · auto-fix
- [ZC1087: Output redirection overwrites input file](#zc1087)
- [ZC1088: Subshell isolates state changes](#zc1088)
- [ZC1089: Redirection order matters (`2>&1 > file`)](#zc1089)
- [ZC1090: Quoted regex pattern in `=~`](#zc1090)
- [ZC1091: Use `((...))` for arithmetic comparisons in `\[\[...\]\]`](#zc1091) · auto-fix
- [ZC1092: Prefer `print` or `printf` over `echo` in Zsh](#zc1092) · auto-fix
- [ZC1093: Superseded by ZC1038 — retired duplicate](#zc1093)
- [ZC1094: Use parameter expansion instead of `sed` for simple substitutions](#zc1094)
- [ZC1095: Use `repeat N` for simple repetition](#zc1095) · auto-fix
- [ZC1096: Warn on `bc` for simple arithmetic](#zc1096)
- [ZC1097: Declare loop variables as `local` in functions](#zc1097)
- [ZC1098: Use `(q)` flag for quoting variables in eval](#zc1098)
- [ZC1099: Use `(f)` flag to split lines instead of `while read`](#zc1099)
- [ZC1100: Use parameter expansion instead of `dirname`/`basename`](#zc1100)
- [ZC1101: Use `$(( ))` instead of `bc` for simple arithmetic](#zc1101)
- [ZC1102: Redirecting output of `sudo` doesn't work as expected](#zc1102)
- [ZC1103: Suggest `path` array instead of `$PATH` string manipulation (direct assignment)](#zc1103)
- [ZC1104: Suggest `path` array instead of `export PATH` string manipulation](#zc1104)
- [ZC1105: Avoid nested arithmetic expansions for clarity](#zc1105)
- [ZC1106: Avoid `set -x` in production scripts for sensitive data exposure](#zc1106)
- [ZC1107: Use (( ... )) for arithmetic conditions](#zc1107)
- [ZC1108: Use Zsh case conversion instead of `tr`](#zc1108)
- [ZC1109: Use parameter expansion instead of `cut` for field extraction](#zc1109)
- [ZC1110: Use Zsh subscripts instead of `head -1` or `tail -1`](#zc1110)
- [ZC1111: Avoid `xargs` for simple command invocation](#zc1111)
- [ZC1112: Avoid `grep -c` — use Zsh pattern matching for counting](#zc1112)
- [ZC1113: Use `${var:A}` instead of `realpath` or `readlink -f`](#zc1113)
- [ZC1114: Consider Zsh `=(...)` for temporary files](#zc1114)
- [ZC1115: Use Zsh string manipulation instead of `rev`](#zc1115)
- [ZC1116: Use Zsh multios instead of `tee`](#zc1116)
- [ZC1117: Use `&!` or `disown` instead of `nohup`](#zc1117)
- [ZC1118: Use `print -rn` instead of `echo -n`](#zc1118) · auto-fix
- [ZC1119: Use `$EPOCHSECONDS` instead of `date +%s`](#zc1119)
- [ZC1120: Use `$PWD` instead of `pwd`](#zc1120)
- [ZC1121: Use `$HOST` instead of `hostname`](#zc1121)
- [ZC1122: Use `$USER` instead of `whoami`](#zc1122)
- [ZC1123: Use `$OSTYPE` instead of `uname`](#zc1123)
- [ZC1124: Use `: > file` instead of `cat /dev/null > file` to truncate](#zc1124) · auto-fix
- [ZC1125: Avoid `echo \| grep` for string matching](#zc1125)
- [ZC1126: Use `sort -u` instead of `sort \| uniq`](#zc1126) · auto-fix
- [ZC1127: Avoid `ls` for counting files](#zc1127)
- [ZC1128: Use `> file` instead of `touch file` for creation](#zc1128) · auto-fix
- [ZC1129: Use Zsh `stat` module instead of `wc -c` for file size](#zc1129)
- [ZC1131: Avoid `cat file \| while read` — use redirection](#zc1131)
- [ZC1132: Use Zsh pattern extraction instead of `grep -o`](#zc1132)
- [ZC1133: Avoid `kill -9` — use `kill` first, then escalate](#zc1133)
- [ZC1134: Avoid `sleep` in tight loops](#zc1134)
- [ZC1135: Avoid `env VAR=val cmd` — use inline assignment](#zc1135) · auto-fix
- [ZC1136: Avoid `rm -rf` without safeguard](#zc1136)
- [ZC1137: Avoid hardcoded `/tmp` paths](#zc1137)
- [ZC1139: Avoid `source` with URL — use local files](#zc1139)
- [ZC1140: Use `command -v` instead of `hash` for command existence](#zc1140) · auto-fix
- [ZC1141: Avoid `curl \| sh` pattern](#zc1141)
- [ZC1142: Avoid chained `grep \| grep` — combine patterns](#zc1142)
- [ZC1143: Avoid `set -e` — use explicit error handling](#zc1143)
- [ZC1144: Avoid `trap` with signal numbers — use names](#zc1144) · auto-fix
- [ZC1145: Avoid `tr -d` for character deletion — use parameter expansion](#zc1145)
- [ZC1146: Avoid `cat file \| awk` — pass file to awk directly](#zc1146) · auto-fix
- [ZC1147: Avoid `mkdir` without `-p` for nested paths](#zc1147) · auto-fix
- [ZC1148: Use `compdef` instead of `compctl` for completions](#zc1148)
- [ZC1149: Avoid `echo` for error messages — use `>&2`](#zc1149)
- [ZC1151: Avoid `cat -A` — use `print -v` or od for non-printable characters](#zc1151)
- [ZC1152: Use Zsh PCRE module instead of `grep -P`](#zc1152)
- [ZC1153: Use `cmp -s` instead of `diff` for equality check](#zc1153) · auto-fix
- [ZC1154: Use `find -exec {} +` instead of `find -exec {} \;`](#zc1154)
- [ZC1155: Use `whence -a` instead of `which -a`](#zc1155) · auto-fix
- [ZC1156: Avoid `ln` without `-s` for symlinks](#zc1156)
- [ZC1157: Avoid `strings` command — use Zsh `${(ps:\0:)var}`](#zc1157)
- [ZC1158: Avoid `chown -R` without `--no-dereference`](#zc1158)
- [ZC1159: Avoid `tar` without explicit compression flag](#zc1159)
- [ZC1160: Prefer `curl` over `wget` for portability](#zc1160)
- [ZC1161: Avoid `openssl` for simple hashing — use Zsh modules](#zc1161)
- [ZC1162: Use `cp -a` instead of `cp -r` to preserve attributes](#zc1162) · auto-fix
- [ZC1163: Use `grep -m 1` instead of `grep \| head -1`](#zc1163) · auto-fix
- [ZC1164: Avoid `sed -n 'Np'` — use Zsh array subscript](#zc1164)
- [ZC1165: Use Zsh parameter expansion for simple `awk` field extraction](#zc1165)
- [ZC1166: Avoid `grep -i` for case-insensitive match — use `(#i)` glob flag](#zc1166)
- [ZC1167: Avoid `timeout` command — use Zsh `TMOUT` or `zsh/sched`](#zc1167)
- [ZC1168: Use `${(f)...}` instead of `readarray`/`mapfile`](#zc1168)
- [ZC1169: Avoid `install` for simple copy+chmod — use `cp` then `chmod`](#zc1169)
- [ZC1170: Avoid `pushd`/`popd` without `-q` flag](#zc1170) · auto-fix
- [ZC1171: Use `print` instead of `echo -e` for escape sequences](#zc1171) · auto-fix
- [ZC1172: Use `read -A` instead of Bash `read -a` for arrays](#zc1172) · auto-fix
- [ZC1173: Avoid `column` command — use Zsh `print -C` for columnar output](#zc1173)
- [ZC1174: Use Zsh `${(j:delim:)}` instead of `paste -sd`](#zc1174)
- [ZC1175: Avoid `tput` for simple ANSI colors — use Zsh `%F{color}`](#zc1175)
- [ZC1176: Use `zparseopts` instead of `getopt`/`getopts`](#zc1176)
- [ZC1177: Avoid `id -u` — use Zsh `$UID` or `$EUID`](#zc1177)
- [ZC1178: Avoid `stty` for terminal size — use Zsh `$COLUMNS`/`$LINES`](#zc1178)
- [ZC1179: Use Zsh `strftime` instead of `date` for formatting](#zc1179)
- [ZC1180: Avoid `pgrep` for own background jobs — use Zsh job control](#zc1180)
- [ZC1181: Avoid `xdg-open`/`open` — use `$BROWSER` for portability](#zc1181)
- [ZC1182: Avoid `nc`/`netcat` for HTTP — use `curl` or `zsh/net/tcp`](#zc1182)
- [ZC1183: Use Zsh glob qualifiers instead of `ls -t` for file ordering](#zc1183)
- [ZC1184: Avoid `diff -u` for patch generation — use `git diff` when in a repo](#zc1184)
- [ZC1185: Use Zsh `${#${(z)var}}` instead of `wc -w` for word count](#zc1185)
- [ZC1186: Use `unset -v` or `unset -f` for explicit unsetting](#zc1186)
- [ZC1187: Avoid `notify-send` without fallback — check availability first](#zc1187)
- [ZC1188: Use Zsh `path+=()` instead of `export PATH=$PATH:dir`](#zc1188)
- [ZC1189: Avoid `source /dev/stdin` — use direct evaluation](#zc1189)
- [ZC1190: Combine chained `grep -v` into single invocation](#zc1190) · auto-fix
- [ZC1191: Avoid `clear` command — use ANSI escape sequences](#zc1191) · auto-fix
- [ZC1192: Avoid `sleep 0` — it is a no-op external process](#zc1192) · auto-fix
- [ZC1193: Avoid `rm -i` in non-interactive scripts](#zc1193)
- [ZC1194: Avoid `sed` with multiple `-e` — use a single script](#zc1194)
- [ZC1195: Avoid overly permissive `umask` values](#zc1195)
- [ZC1196: Avoid `cat` for reading single file into variable](#zc1196)
- [ZC1197: Avoid `more` in scripts — use `cat` or pager check](#zc1197)
- [ZC1198: Avoid interactive editors in scripts](#zc1198)
- [ZC1199: Avoid `telnet` in scripts — use `curl` or `zsh/net/tcp`](#zc1199)
- [ZC1200: Avoid `ftp` — use `sftp` or `curl` for secure transfers](#zc1200)
- [ZC1201: Avoid `rsh`/`rlogin`/`rcp` — use `ssh`/`scp`](#zc1201)
- [ZC1202: Avoid `ifconfig` — use `ip` for network configuration](#zc1202) · auto-fix
- [ZC1203: Avoid `netstat` — use `ss` for socket statistics](#zc1203) · auto-fix
- [ZC1204: Avoid `route` — use `ip route` for routing](#zc1204)
- [ZC1205: Avoid `arp` — use `ip neigh` for neighbor tables](#zc1205)
- [ZC1206: Avoid `crontab -e` in scripts — use `crontab file`](#zc1206)
- [ZC1207: Avoid `passwd` in scripts — use `chpasswd`](#zc1207)
- [ZC1208: Avoid `visudo` in scripts — use sudoers.d drop-in files](#zc1208)
- [ZC1209: Use `systemctl --no-pager` in scripts](#zc1209) · auto-fix
- [ZC1210: Use `journalctl --no-pager` in scripts](#zc1210) · auto-fix
- [ZC1211: Use `git stash push -m` instead of bare `git stash`](#zc1211)
- [ZC1212: Avoid `git add .` — use explicit paths or `git add -p`](#zc1212)
- [ZC1213: Use `apt-get -y` in scripts for non-interactive installs](#zc1213) · auto-fix
- [ZC1214: Avoid `su` in scripts — use `sudo -u` for user switching](#zc1214)
- [ZC1215: Source `/etc/os-release` instead of parsing with `cat`/`grep`](#zc1215) · auto-fix
- [ZC1216: Avoid `nslookup` — use `dig` or `host` for DNS queries](#zc1216) · auto-fix
- [ZC1217: Avoid `service` command — use `systemctl` on systemd](#zc1217) · auto-fix
- [ZC1218: Avoid `useradd` without `--shell /sbin/nologin` for service accounts](#zc1218)
- [ZC1219: Use `curl -fsSL` instead of `wget -O -` for piped downloads](#zc1219) · auto-fix
- [ZC1220: Use `chown :group` instead of `chgrp` for group changes](#zc1220)
- [ZC1221: Avoid `fdisk` in scripts — use `parted` or `sfdisk`](#zc1221)
- [ZC1222: Avoid `lsof -i` for port checks — use `ss -tlnp`](#zc1222)
- [ZC1223: Avoid `ip addr show` piped to `grep` — use `ip -br addr`](#zc1223)
- [ZC1224: Avoid parsing `free` output — read `/proc/meminfo` directly](#zc1224)
- [ZC1225: Avoid parsing `uptime` — read `/proc/uptime` directly](#zc1225)
- [ZC1226: Use `dmesg -T` or `--time-format=iso` for readable timestamps](#zc1226) · auto-fix
- [ZC1227: Use `curl -f` to fail on HTTP errors](#zc1227) · auto-fix
- [ZC1228: Avoid `ssh` without host key policy in scripts](#zc1228)
- [ZC1229: Prefer `rsync` over `scp` for file transfers](#zc1229)
- [ZC1230: Use `ping -c N` in scripts to limit ping count](#zc1230) · auto-fix
- [ZC1231: Use `git clone --depth 1` for CI and build scripts](#zc1231) · auto-fix
- [ZC1232: Avoid bare `pip install` — use `--user` or virtualenv](#zc1232)
- [ZC1233: Avoid `npm install -g` — use `npx` for one-off tools](#zc1233)
- [ZC1234: Use `docker run --rm` to auto-remove containers](#zc1234) · auto-fix
- [ZC1235: Use `git push --force-with-lease` instead of `--force`](#zc1235) · auto-fix
- [ZC1236: Avoid `git reset --hard` — irreversible data loss risk](#zc1236)
- [ZC1237: Use `git clean -n` before `git clean -fd`](#zc1237)
- [ZC1238: Avoid `docker exec -it` in scripts — drop `-it` for non-interactive](#zc1238) · auto-fix
- [ZC1239: Avoid `kubectl exec -it` in scripts](#zc1239) · auto-fix
- [ZC1240: Use `find -maxdepth` with `-delete` to limit scope](#zc1240)
- [ZC1241: Use `xargs -0` with null separators for safe argument passing](#zc1241) · auto-fix
- [ZC1242: Use `tar -C dir` to extract into a specific directory](#zc1242)
- [ZC1243: Use `grep -lZ` with `xargs -0` for safe file lists](#zc1243)
- [ZC1244: Consider `mv -n` to prevent overwriting existing files](#zc1244)
- [ZC1245: Avoid disabling TLS certificate verification](#zc1245)
- [ZC1246: Avoid hardcoded passwords in command arguments](#zc1246)
- [ZC1247: Avoid `chmod +s` — setuid/setgid bits are security risks](#zc1247)
- [ZC1248: Prefer `ufw`/`firewalld` over raw `iptables`](#zc1248)
- [ZC1249: Use `ssh-keygen -f` to specify key file in scripts](#zc1249)
- [ZC1250: Use `gpg --batch` in scripts for non-interactive operation](#zc1250)
- [ZC1251: Use `mount -o noexec,nosuid` for untrusted media](#zc1251)
- [ZC1252: Use `getent passwd` instead of `cat /etc/passwd`](#zc1252) · auto-fix
- [ZC1253: Use `docker build --no-cache` in CI for reproducible builds](#zc1253) · auto-fix
- [ZC1254: Avoid `git commit --amend` in shared branches](#zc1254)
- [ZC1255: Use `curl -L` to follow HTTP redirects](#zc1255) · auto-fix
- [ZC1256: Clean up `mkfifo` pipes with a trap on EXIT](#zc1256)
- [ZC1257: Use `docker stop -t` to set graceful shutdown timeout](#zc1257) · auto-fix
- [ZC1258: Consider `rsync --delete` for directory sync](#zc1258)
- [ZC1259: Avoid `docker pull` without explicit tag — pin image versions](#zc1259)
- [ZC1260: Use `git branch -d` instead of `-D` for safe deletion](#zc1260) · auto-fix
- [ZC1261: Avoid piping `base64 -d` output to shell execution](#zc1261)
- [ZC1262: Avoid `chmod -R 777` — recursive world-writable is critical](#zc1262)
- [ZC1263: Use `apt-get` instead of `apt` in scripts](#zc1263) · auto-fix
- [ZC1264: Use `dnf` instead of `yum` on modern Fedora/RHEL](#zc1264) · auto-fix
- [ZC1265: Use `systemctl enable --now` to enable and start together](#zc1265) · auto-fix
- [ZC1266: Use `nproc` instead of parsing `/proc/cpuinfo`](#zc1266)
- [ZC1267: Use `df -P` for POSIX-portable disk usage output](#zc1267) · auto-fix
- [ZC1268: Use `du -sh --` to handle filenames starting with dash](#zc1268) · auto-fix
- [ZC1269: Use `pgrep` instead of `ps aux \| grep` for process search](#zc1269)
- [ZC1270: Use `mktemp` instead of hardcoded `/tmp` paths](#zc1270)
- [ZC1271: Use `command -v` instead of `which` for command existence checks](#zc1271) · auto-fix
- [ZC1272: Use `install -m` instead of separate `cp` and `chmod`](#zc1272)
- [ZC1273: Use `grep -q` instead of redirecting grep output to `/dev/null`](#zc1273) · auto-fix
- [ZC1274: Use Zsh `${var:t}` instead of `basename`](#zc1274)
- [ZC1275: Use Zsh `${var:h}` instead of `dirname`](#zc1275)
- [ZC1276: Use Zsh `{start..end}` instead of `seq`](#zc1276) · auto-fix
- [ZC1277: Superseded by ZC1108 — retired duplicate](#zc1277)
- [ZC1278: Superseded by ZC1009 — retired duplicate](#zc1278)
- [ZC1279: Use `realpath` instead of `readlink -f` for canonical paths](#zc1279) · auto-fix
- [ZC1280: Use `Zsh ${var:e}` instead of shell expansion to extract file extension](#zc1280)
- [ZC1281: Use `sort -u` instead of `sort \| uniq` for deduplication](#zc1281)
- [ZC1282: Use Zsh `${var:r}` instead of `sed` to remove file extension](#zc1282)
- [ZC1283: Use `setopt` instead of `set -o` for Zsh options](#zc1283) · auto-fix
- [ZC1284: Use Zsh `${(s:sep:)var}` instead of `cut -d` for field splitting](#zc1284)
- [ZC1285: Use Zsh `${(o)array}` for sorting instead of piping to `sort`](#zc1285)
- [ZC1286: Use Zsh `${array:#pattern}` instead of `grep -v` for filtering](#zc1286)
- [ZC1287: Use `cat -v` alternative: Zsh `${(V)var}` for visible control characters](#zc1287)
- [ZC1288: Use `typeset` instead of `declare` in Zsh scripts](#zc1288) · auto-fix
- [ZC1289: Use Zsh `${(u)array}` for unique elements instead of `sort -u`](#zc1289)
- [ZC1290: Use Zsh `${(n)array}` for numeric sorting instead of `sort -n`](#zc1290)
- [ZC1291: Use Zsh `${(O)array}` for reverse sorting instead of `sort -r`](#zc1291)
- [ZC1292: Use Zsh `${var//old/new}` instead of `tr` for character translation](#zc1292)
- [ZC1293: Use `\[\[ \]\]` instead of `test` command in Zsh](#zc1293)
- [ZC1294: Use `bindkey` instead of `bind` for key bindings in Zsh](#zc1294)
- [ZC1295: Use `vared` instead of `read -e` for interactive editing in Zsh](#zc1295)
- [ZC1296: Avoid `shopt` in Zsh — use `setopt`/`unsetopt` instead](#zc1296)
- [ZC1297: Avoid `$BASH_SOURCE` — use `$0` or `${(%):-%x}` in Zsh](#zc1297) · auto-fix
- [ZC1298: Avoid `$FUNCNAME` — use `$funcstack` in Zsh](#zc1298) · auto-fix
- [ZC1299: Avoid `$BASH_LINENO` — use `$funcfiletrace` in Zsh](#zc1299)
- [ZC1300: Avoid `$BASH_VERSINFO` — use `$ZSH_VERSION` in Zsh](#zc1300) · auto-fix
- [ZC1301: Avoid `$PIPESTATUS` — use `$pipestatus` (lowercase) in Zsh](#zc1301) · auto-fix
- [ZC1302: Avoid `help` builtin — use `run-help` or `man` in Zsh](#zc1302)
- [ZC1303: Avoid `enable` command — use `zmodload` for Zsh modules](#zc1303)
- [ZC1304: Avoid `$BASH_SUBSHELL` — use `$ZSH_SUBSHELL` in Zsh](#zc1304) · auto-fix
- [ZC1305: Avoid `$COMP_WORDS` — use `$words` in Zsh completion](#zc1305) · auto-fix
- [ZC1306: Avoid `$COMP_CWORD` — use `$CURRENT` in Zsh completion](#zc1306) · auto-fix
- [ZC1307: Avoid `$DIRSTACK` — use `$dirstack` (lowercase) in Zsh](#zc1307) · auto-fix
- [ZC1308: Avoid `$COMP_LINE` — use `$BUFFER` in Zsh completion](#zc1308) · auto-fix
- [ZC1309: Avoid `$BASH_COMMAND` — not available in Zsh](#zc1309)
- [ZC1310: Avoid `$BASH_EXECUTION_STRING` — not available in Zsh](#zc1310)
- [ZC1311: Avoid `complete` command — use `compdef` in Zsh](#zc1311)
- [ZC1312: Avoid `compgen` command — use `compadd` in Zsh](#zc1312)
- [ZC1313: Avoid `$BASH_ALIASES` — use Zsh `aliases` hash](#zc1313) · auto-fix
- [ZC1314: Avoid `$BASH_LOADABLES_PATH` — not available in Zsh](#zc1314)
- [ZC1315: Avoid `$BASH_COMPAT` — use `emulate` for compatibility in Zsh](#zc1315)
- [ZC1316: Avoid `caller` builtin — use `$funcfiletrace` in Zsh](#zc1316)
- [ZC1317: Avoid `$BASH_ENV` — use `$ZDOTDIR` and `$ENV` in Zsh](#zc1317)
- [ZC1318: Avoid `$BASH_CMDS` — use `$commands` hash in Zsh](#zc1318) · auto-fix
- [ZC1319: Avoid `$BASH_ARGC` — use `$#` in Zsh](#zc1319) · auto-fix
- [ZC1320: Avoid `$BASH_ARGV` — use `$argv` in Zsh](#zc1320) · auto-fix
- [ZC1321: Avoid `$BASH_XTRACEFD` — not available in Zsh](#zc1321)
- [ZC1322: Avoid `$COPROC` — Zsh coproc uses different syntax](#zc1322)
- [ZC1323: Avoid `suspend` builtin — use `kill -STOP $$` in Zsh](#zc1323)
- [ZC1324: Avoid `$PROMPT_COMMAND` — use `precmd` hook in Zsh](#zc1324)
- [ZC1325: Avoid `$PS0` — use `preexec` hook in Zsh](#zc1325)
- [ZC1326: Avoid `$HISTTIMEFORMAT` — use `fc -li` in Zsh](#zc1326)
- [ZC1327: Avoid `history -c` — Zsh uses different history management](#zc1327)
- [ZC1328: Avoid `$HISTCONTROL` — use Zsh `setopt` history options](#zc1328)
- [ZC1329: Avoid `$HISTIGNORE` — use `zshaddhistory` hook in Zsh](#zc1329)
- [ZC1330: Avoid `$INPUTRC` — use `bindkey` in Zsh](#zc1330)
- [ZC1331: Avoid `$BASH_REMATCH` — use `$match` array in Zsh](#zc1331) · auto-fix
- [ZC1332: Avoid `$GLOBIGNORE` — use `setopt EXTENDED_GLOB` in Zsh](#zc1332)
- [ZC1333: Avoid `$TIMEFORMAT` — use `$TIMEFMT` in Zsh](#zc1333) · auto-fix
- [ZC1334: Avoid `type -p` — use `whence -p` in Zsh](#zc1334) · auto-fix
- [ZC1335: Use Zsh array reversal instead of `tac` for in-memory data](#zc1335)
- [ZC1336: Avoid `printenv` — use `typeset -x` or `export` in Zsh](#zc1336)
- [ZC1337: Avoid `fold` command — use Zsh `print -l` with `$COLUMNS`](#zc1337)
- [ZC1338: Avoid `seq -s` — use Zsh `${(j:sep:)${(s::)...}}` for joining](#zc1338)
- [ZC1339: Use Zsh `${#${(f)var}}` instead of `wc -l` for line count](#zc1339)
- [ZC1340: Avoid `shuf` for random array element — use Zsh `$RANDOM`](#zc1340)
- [ZC1341: Use Zsh `*(.x)` glob qualifier instead of `find -executable`](#zc1341)
- [ZC1342: Use Zsh `*(L0)` glob qualifier instead of `find -empty`](#zc1342)
- [ZC1343: Use Zsh `*(m±N)` glob qualifier instead of `find -mtime N`](#zc1343)
- [ZC1344: Use Zsh `*(L±Nk)` glob qualifier instead of `find -size`](#zc1344)
- [ZC1345: Use Zsh `*(f:mode:)` glob qualifier instead of `find -perm`](#zc1345)
- [ZC1346: Use Zsh `*(u:name:)` glob qualifier instead of `find -user`](#zc1346)
- [ZC1347: Use Zsh `*(g:name:)` glob qualifier instead of `find -group`](#zc1347)
- [ZC1348: Use Zsh glob type qualifiers instead of `find -type`](#zc1348)
- [ZC1349: Use `${#var}` instead of `expr length "$var"` for string length](#zc1349)
- [ZC1350: Use `${str:pos:len}` instead of `expr substr` for substring extraction](#zc1350)
- [ZC1351: Use `\[\[ $str =~ pattern \]\]` instead of `expr match` / `expr :` for regex](#zc1351)
- [ZC1352: Avoid `xargs -I{}` — use a Zsh `for` loop for per-item substitution](#zc1352)
- [ZC1353: Avoid `printf -v` — use `print -v` or command substitution in Zsh](#zc1353)
- [ZC1354: Use `whence -w` instead of Bash-specific `type -t` for command classification](#zc1354)
- [ZC1355: Use `print -r` instead of `echo -E` for raw output](#zc1355) · auto-fix
- [ZC1356: Use `read -A` instead of `read -a` for array read in Zsh](#zc1356) · auto-fix
- [ZC1357: Use Zsh `${(q)var}` instead of `printf '%q'` for shell-quoting](#zc1357)
- [ZC1358: Use `${PWD:P}` instead of `pwd -P` for physical current directory](#zc1358)
- [ZC1359: Avoid `id -Gn` — use Zsh `$groups` associative array](#zc1359)
- [ZC1360: Use Zsh `*(OL)` glob qualifier instead of `ls -S` for size-ordered listing](#zc1360)
- [ZC1361: Avoid `awk 'NR==N'` — use Zsh array subscript on `${(f)...}`](#zc1361)
- [ZC1362: Use `\[\[ -o option \]\]` instead of `test -o option` for Zsh option checks](#zc1362)
- [ZC1363: Use Zsh `*(e:...:)` eval qualifier instead of `find -newer`/`-older`](#zc1363)
- [ZC1364: Use Zsh `${var:pos:len}` instead of `cut -c` for character ranges](#zc1364)
- [ZC1365: Use Zsh `zstat` module instead of `stat -c` for file metadata](#zc1365)
- [ZC1366: Use Zsh `limit` instead of POSIX `ulimit` for idiomatic resource queries](#zc1366)
- [ZC1367: Use Zsh `strftime` instead of Bash `printf '%(fmt)T'`](#zc1367)
- [ZC1368: Avoid `sh -c` / `bash -c` inside a Zsh script — inline or use a function](#zc1368)
- [ZC1369: Prefer Zsh `${(V)var}` over `od -c` for printable-visible character output](#zc1369)
- [ZC1370: Prefer Zsh `repeat N { ... }` over `yes str \| head -n N` for finite output](#zc1370)
- [ZC1371: Use Zsh array `:t` modifier instead of `basename -a` for bulk path stripping](#zc1371)
- [ZC1372: Use Zsh `zmv` autoload function instead of `rename`/`rename.ul`](#zc1372)
- [ZC1373: Use Zsh `${(0)var}` flag for NUL-split parsing instead of `env -0`](#zc1373)
- [ZC1374: Avoid `$FUNCNEST` — Zsh uses `$FUNCNEST` as a limit, not a depth indicator](#zc1374) · auto-fix
- [ZC1375: Use `\[\[ -t fd \]\]` instead of `tty -s` for tty-check](#zc1375)
- [ZC1376: Avoid `BASH_XTRACEFD` — use Zsh `exec {fd}>file` + `setopt XTRACE`](#zc1376)
- [ZC1377: Avoid `$BASH_ALIASES` — use Zsh `$aliases` associative array](#zc1377) · auto-fix
- [ZC1378: Avoid uppercase `$DIRSTACK` — Zsh uses lowercase `$dirstack`](#zc1378) · auto-fix
- [ZC1379: Avoid `$PROMPT_COMMAND` — use Zsh `precmd` function](#zc1379)
- [ZC1380: Avoid `$HISTIGNORE` — use Zsh `$HISTORY_IGNORE`](#zc1380) · auto-fix
- [ZC1381: Avoid `$COMP_WORDS`/`$COMP_CWORD` — Zsh uses `words`/`$CURRENT`](#zc1381) · auto-fix
- [ZC1382: Avoid `$READLINE_LINE`/`$READLINE_POINT` — Zsh ZLE uses `$BUFFER`/`$CURSOR`](#zc1382) · auto-fix
- [ZC1383: Avoid `$TIMEFORMAT` — Zsh uses `$TIMEFMT`](#zc1383) · auto-fix
- [ZC1384: Avoid `$EXECIGNORE` — Bash-only; Zsh uses completion-system ignore patterns](#zc1384)
- [ZC1385: Avoid `$PS0` — Bash-only; Zsh uses `preexec` hook](#zc1385)
- [ZC1386: Avoid `$FIGNORE` — Bash-only; Zsh uses compsys tag patterns](#zc1386)
- [ZC1387: Avoid `$SHELLOPTS` — Zsh uses `$options` associative array](#zc1387)
- [ZC1388: Use Zsh lowercase `$mailpath` array instead of colon-separated `$MAILPATH`](#zc1388)
- [ZC1389: Avoid `$HOSTFILE` — Bash-only; Zsh uses `$hosts` array](#zc1389)
- [ZC1390: Avoid `$GROUPS\[@\]` — Zsh `$GROUPS` is a scalar, not an array](#zc1390)
- [ZC1391: Avoid `\[\[ -v VAR \]\]` for Bash set-check — use Zsh `(( ${+VAR} ))`](#zc1391)
- [ZC1392: Avoid `$CHILD_MAX` — Bash-only; Zsh uses `limit` / `ulimit -u`](#zc1392)
- [ZC1393: Avoid `$SRANDOM` — Bash 5.1+ only, read `/dev/urandom` in Zsh](#zc1393)
- [ZC1394: Avoid `$BASH` — Zsh uses `$ZSH_NAME` for the interpreter name](#zc1394) · auto-fix
- [ZC1395: Avoid `wait -n` — Bash 4.3+ only; Zsh `wait` on job IDs](#zc1395)
- [ZC1396: Avoid `unset -n` — Bash nameref semantics not in Zsh](#zc1396)
- [ZC1397: Avoid `$COMP_TYPE`/`$COMP_KEY` — Bash completion globals, not in Zsh](#zc1397)
- [ZC1398: Avoid `$PROMPT_DIRTRIM` — use Zsh `%N~` prompt modifier](#zc1398)
- [ZC1399: Use Zsh `$signals` array instead of `kill -l` for signal enumeration](#zc1399)
- [ZC1400: Use Zsh `$CPUTYPE` for architecture detection instead of parsing `$HOSTTYPE`](#zc1400)
- [ZC1401: Prefer Zsh `$VENDOR` over parsing `$MACHTYPE` for vendor detection](#zc1401)
- [ZC1402: Avoid `date -d @seconds` — use Zsh `strftime` for epoch formatting](#zc1402)
- [ZC1403: Setting `$HISTFILESIZE` alone is incomplete in Zsh — pair with `$SAVEHIST`](#zc1403) · auto-fix
- [ZC1404: Avoid `$BASH_CMDS` — Bash-specific hash-table mirror, use Zsh `$commands`](#zc1404) · auto-fix
- [ZC1405: Avoid `env -u VAR cmd` — use Zsh `(unset VAR; cmd)` subshell](#zc1405)
- [ZC1406: Prefer Zsh `zargs -P N` autoload over `xargs -P N` for parallel execution](#zc1406)
- [ZC1407: Avoid `/dev/tcp/...` — use Zsh `zsh/net/tcp` module](#zc1407)
- [ZC1408: Avoid `$BASH_FUNC_...%%` — Bash-specific exported-function envvar](#zc1408)
- [ZC1409: Avoid `\[ -N file \]` / `test -N file` — Bash-only, use Zsh `zstat` for mtime comparison](#zc1409)
- [ZC1410: Avoid `compopt` — Bash programmable-completion modifier, not in Zsh](#zc1410)
- [ZC1411: Use Zsh `disable` instead of Bash `enable -n` to hide builtins](#zc1411) · auto-fix
- [ZC1412: Avoid `$COMPREPLY` — Bash completion output, use Zsh `compadd`](#zc1412)
- [ZC1413: Use Zsh `whence -p cmd` instead of `hash -t cmd` for resolved path](#zc1413) · auto-fix
- [ZC1414: Beware `hash -d` — Bash deletes from hash table, Zsh defines named directory](#zc1414)
- [ZC1415: Prefer Zsh `TRAPZERR` function over `trap 'cmd' ERR`](#zc1415)
- [ZC1416: Prefer Zsh `preexec` hook over `trap 'cmd' DEBUG`](#zc1416)
- [ZC1417: Prefer Zsh `TRAPRETURN` function over `trap 'cmd' RETURN`](#zc1417)
- [ZC1418: Use Zsh `limit -h`/`-s` instead of `ulimit -H`/`-S` for hard/soft limits](#zc1418)
- [ZC1419: Avoid `chmod 777` — grants world-writable access](#zc1419)
- [ZC1420: Avoid `chmod +s` / `chmod u+s` — setuid/setgid is a security risk](#zc1420)
- [ZC1421: Avoid `chpasswd` / `passwd --stdin` — plaintext passwords in process tree](#zc1421)
- [ZC1422: Avoid `sudo -S` — reads password from stdin, exposes plaintext](#zc1422)
- [ZC1423: Dangerous: `iptables -F` / `nft flush ruleset` — drops all firewall rules](#zc1423)
- [ZC1424: Dangerous: `mkfs.*` / `mkfs -t` — formats a filesystem, destroys data](#zc1424)
- [ZC1425: `shutdown` / `reboot` / `halt` / `poweroff` — confirm before scripting](#zc1425)
- [ZC1426: Avoid `git clone http://` — unencrypted transport, use `https://` or `git://`+verify](#zc1426)
- [ZC1427: Dangerous: `nc -e` / `ncat -e` — spawns arbitrary command on network connect](#zc1427)
- [ZC1428: Avoid `curl -u user:pass` — credentials visible in process list](#zc1428)
- [ZC1429: Avoid `umount -f` / `-l` — force/lazy unmount masks real issues](#zc1429)
- [ZC1430: Prefer Zsh `zsh/sched` module over `at now` / `batch` for in-shell scheduling](#zc1430)
- [ZC1431: Dangerous: `crontab -r` — removes all the user's cron jobs without confirmation](#zc1431)
- [ZC1432: Dangerous: `passwd -d user` — deletes the password, leaving the account passwordless](#zc1432)
- [ZC1433: Caution with `userdel -f` / `-r` — removes home directory and kills processes](#zc1433)
- [ZC1434: Warn on `swapoff -a` — disables all swap, can OOM-kill](#zc1434)
- [ZC1435: Avoid `killall -9` / `killall -KILL` — force-kill by process name](#zc1435)
- [ZC1436: `sysctl -w` is ephemeral — persist in `/etc/sysctl.d/*.conf` for surviving reboots](#zc1436)
- [ZC1437: `dmesg -c` / `-C` clears the kernel ring buffer — destroys evidence](#zc1437)
- [ZC1438: `systemctl mask` permanently prevents service start — document the unmask path](#zc1438)
- [ZC1439: Enabling IP forwarding in a script — document firewall posture](#zc1439)
- [ZC1440: `usermod -G group user` replaces supplementary groups — use `-aG` to append](#zc1440)
- [ZC1441: Warn on `docker system prune -af` / `-a --force` (or similar podman/k8s)](#zc1441)
- [ZC1442: Dangerous: `kubectl delete --all` / `--all-namespaces` deletes cluster resources](#zc1442)
- [ZC1443: Dangerous: `terraform destroy` / `apply -destroy` without `-target`](#zc1443)
- [ZC1444: Dangerous: `redis-cli FLUSHALL` / `FLUSHDB` — wipes Redis data](#zc1444)
- [ZC1445: Dangerous: `dropdb` / `mysqladmin drop` — deletes a database](#zc1445)
- [ZC1446: Dangerous: `aws s3 rm --recursive` / `s3 rb --force` — bulk S3 deletion](#zc1446)
- [ZC1447: Avoid deprecated `ifconfig` / `netstat` — prefer `ip` / `ss`](#zc1447)
- [ZC1448: `apt-get install` / `apt install` without `-y` hangs in non-interactive scripts](#zc1448) · auto-fix
- [ZC1449: `dnf`/`yum` install without `-y` hangs in non-interactive scripts](#zc1449)
- [ZC1450: `pacman -S` / `zypper install` without non-interactive flag hangs in scripts](#zc1450)
- [ZC1451: Avoid `pip install` without `--user` or virtualenv](#zc1451)
- [ZC1452: Avoid `npm install -g` — global installs need root, break under multiple Node versions](#zc1452)
- [ZC1453: Avoid `sudo pip` / `sudo npm` / `sudo gem` — language package managers as root](#zc1453)
- [ZC1454: Avoid `docker/podman run --privileged` — disables most container isolation](#zc1454)
- [ZC1455: Avoid `docker run --net=host` / `--network=host` — disables network isolation](#zc1455)
- [ZC1456: Avoid `docker run -v /:...` — bind-mounts host root into container](#zc1456)
- [ZC1457: Warn on bind-mount of `/var/run/docker.sock` — container escape vector](#zc1457)
- [ZC1458: Warn on explicit `docker run --user root` / `--user 0`](#zc1458)
- [ZC1459: Warn on `docker run --cap-add=SYS_ADMIN` / other dangerous capabilities](#zc1459)
- [ZC1460: Warn on `docker run --security-opt seccomp=unconfined` / `apparmor=unconfined`](#zc1460)
- [ZC1461: Avoid `docker run --pid=host` — shares host PID namespace with the container](#zc1461)
- [ZC1462: Avoid `docker run --ipc=host` — shares host IPC namespace (/dev/shm, SysV IPC)](#zc1462)
- [ZC1463: Avoid `docker run --userns=host` — disables user-namespace remapping](#zc1463)
- [ZC1464: Warn on `iptables -F` / `-P INPUT ACCEPT` — flushes or opens the host firewall](#zc1464)
- [ZC1465: Warn on `setenforce 0` — disables SELinux enforcement](#zc1465)
- [ZC1466: Warn on disabling the host firewall (`ufw disable` / `systemctl stop firewalld`)](#zc1466)
- [ZC1467: Warn on `sysctl -w kernel.core_pattern=\|...` / `kernel.modprobe=...` (kernel hijack)](#zc1467)
- [ZC1468: Error on apt `--allow-unauthenticated` / `--force-yes` — installs unsigned packages](#zc1468)
- [ZC1469: Error on `dnf/yum --nogpgcheck` or `rpm --nosignature` (unsigned RPM install)](#zc1469)
- [ZC1470: Error on `git config http.sslVerify false` / `git -c http.sslVerify=false`](#zc1470)
- [ZC1471: Error on `kubectl/helm --insecure-skip-tls-verify` (cluster MITM)](#zc1471)
- [ZC1472: Error on `aws s3 --acl public-read` / `public-read-write` (public bucket)](#zc1472)
- [ZC1473: Warn on `openssl req ... -nodes` / `genrsa` without passphrase — unencrypted private key](#zc1473)
- [ZC1474: Warn on `ssh-keygen -N ""` — generates passwordless SSH key](#zc1474)
- [ZC1475: Warn on `setcap` granting dangerous capabilities to a binary (privesc)](#zc1475)
- [ZC1476: Warn on `apt-key add` — deprecated, trusts every repo system-wide](#zc1476)
- [ZC1477: Warn on `printf "$var"` — variable in format-string position (printf-fmt attack)](#zc1477)
- [ZC1478: Avoid `mktemp -u` — returns a name without creating the file (TOCTOU)](#zc1478)
- [ZC1479: Error on `ssh/scp -o StrictHostKeyChecking=no` / `UserKnownHostsFile=/dev/null`](#zc1479)
- [ZC1480: Warn on `terraform apply -auto-approve` / `destroy -auto-approve` in scripts](#zc1480)
- [ZC1481: Warn on `unset HISTFILE` / `export HISTFILE=/dev/null` — disables shell history](#zc1481)
- [ZC1482: Error on `docker login -p` / `--password=` — credential in process list](#zc1482)
- [ZC1483: Warn on `pip install --break-system-packages` — bypasses PEP 668 externally-managed guard](#zc1483)
- [ZC1484: Error on `npm/yarn/pnpm config set strict-ssl false` — disables registry TLS verification](#zc1484)
- [ZC1485: Warn on `openssl s_client -ssl3 / -tls1 / -tls1_1` — legacy TLS](#zc1485)
- [ZC1486: Warn on `curl -2` / `-3` — forces broken SSLv2 / SSLv3](#zc1486)
- [ZC1487: Warn on `history -c` — clears shell history (and is a Bash-ism under Zsh)](#zc1487)
- [ZC1488: Warn on `ssh -R 0.0.0.0:...` / `*:...` — reverse tunnel bound to all interfaces](#zc1488)
- [ZC1489: Error on `nc -e` / `ncat -e` — classic reverse-shell invocation](#zc1489)
- [ZC1490: Error on `socat ... EXEC:<shell>` / `SYSTEM:<shell>` — socat reverse-shell pattern](#zc1490)
- [ZC1491: Warn on `export LD_PRELOAD=...` / `LD_LIBRARY_PATH=...` — library injection](#zc1491)
- [ZC1492: Style: `at` / `batch` for deferred execution — prefer systemd timers for auditability](#zc1492)
- [ZC1493: Warn on `watch -n 0` — zero-interval watch spins CPU](#zc1493)
- [ZC1494: Warn on `tcpdump -w <file>` without `-Z <user>` — capture file owned by root](#zc1494)
- [ZC1495: Warn on `ulimit -c unlimited` — enables core dumps from setuid binaries](#zc1495)
- [ZC1496: Error on reading `/dev/mem` / `/dev/kmem` / `/dev/port` — leaks physical memory](#zc1496)
- [ZC1497: Error on `useradd -u 0` / `usermod -u 0` — creates a second root account](#zc1497)
- [ZC1498: Warn on `mount -o remount,rw /` — makes read-only root filesystem writable](#zc1498)
- [ZC1499: Style: `docker pull <image>` / `:latest` — unpinned image tag](#zc1499)
- [ZC1500: Warn on `systemctl edit <unit>` in scripts — requires interactive editor](#zc1500)
- [ZC1501: Style: `docker-compose` (hyphen) — use `docker compose` (space, built-in plugin)](#zc1501) · auto-fix
- [ZC1502: Warn on `grep "$var" file` without `--` — flag injection when `$var` starts with `-`](#zc1502)
- [ZC1503: Error on `groupadd -g 0` / `groupmod -g 0` — creates duplicate root group](#zc1503)
- [ZC1504: Warn on `git push --mirror` — overwrites every remote ref](#zc1504)
- [ZC1505: Warn on `dpkg --force-confnew` / `--force-confold` — silently overrides /etc changes](#zc1505)
- [ZC1506: Warn on `newgrp <group>` in scripts — spawns a new shell, breaks control flow](#zc1506)
- [ZC1507: Warn on `rsync -l` / default symlink handling — follows escaping symlinks](#zc1507)
- [ZC1508: Style: `ldd <binary>` may execute the binary — use `objdump -p` / `readelf -d` for untrusted files](#zc1508)
- [ZC1509: Warn on `trap '' TERM` / `trap - TERM` — ignores/resets fatal signal](#zc1509)
- [ZC1510: Error on `auditctl -e 0` / `auditctl -D` — disables kernel audit logging](#zc1510)
- [ZC1511: Error on `nmcli ... <wireless/vpn secret>` on command line](#zc1511)
- [ZC1512: Style: `service <unit> <verb>` — use `systemctl <verb> <unit>` on systemd hosts](#zc1512) · auto-fix
- [ZC1513: Style: `make install` without `DESTDIR=` — unmanaged system-wide install](#zc1513)
- [ZC1514: Error on `useradd -p <hash>` / `usermod -p <hash>` — password hash on cmdline](#zc1514)
- [ZC1515: Warn on `md5sum` / `sha1sum` for integrity check — collision-vulnerable](#zc1515)
- [ZC1516: Error on `umask 000` / `umask 0` — new files / directories world-writable](#zc1516)
- [ZC1517: Warn on `print -P "$var"` — prompt-escape injection via user-controlled string](#zc1517)
- [ZC1518: Warn on `bash -p` — privileged mode (skips env sanitisation on setuid)](#zc1518)
- [ZC1519: Warn on `ulimit -u unlimited` — removes user process cap, enables fork bombs](#zc1519)
- [ZC1520: Warn on `vared <var>` in scripts — reads interactively, hangs non-interactive](#zc1520)
- [ZC1521: Style: `strace` without `-e` filter — captures every syscall (incl. secrets, huge output)](#zc1521)
- [ZC1522: Warn on `ip route add default` / `route add default` — changes default gateway](#zc1522)
- [ZC1523: Error on `tar -C /` — extracting an archive into the filesystem root](#zc1523)
- [ZC1524: Warn on `sysctl -e` / `sysctl -q` — silently skip unknown keys, hide config drift](#zc1524)
- [ZC1525: Warn on `ping -f` — flood ping sends packets as fast as possible](#zc1525)
- [ZC1526: Error on `wipefs -a` / `wipefs -af` — erases filesystem signatures (unrecoverable)](#zc1526)
- [ZC1527: Warn on `crontab -` — replaces cron from stdin, overwrites without diff](#zc1527)
- [ZC1528: Warn on `chage -M 99999` / `-E -1` — disables password aging / expiry](#zc1528)
- [ZC1529: Warn on `fsck -y` / `fsck.<fs> -y` — auto-answer yes can corrupt](#zc1529)
- [ZC1530: Warn on `pkill -f <pattern>` — matches full command line, easy to over-kill](#zc1530)
- [ZC1531: Warn on `wget -t 0` — infinite retries, hangs on a dead endpoint](#zc1531)
- [ZC1532: Warn on `screen -dm` / `tmux new-session -d` — detached long-running session](#zc1532)
- [ZC1533: Warn on `setsid <cmd>` — detaches from controlling TTY, escapes supervision](#zc1533)
- [ZC1534: Warn on `dmesg -c` / `--clear` — wipes kernel ring buffer](#zc1534)
- [ZC1535: Warn on `ip link set <iface> promisc on` — enables packet capture](#zc1535)
- [ZC1536: Warn on `iptables -j DNAT` / `-j REDIRECT` — rewrites traffic destination](#zc1536)
- [ZC1537: Error on `lvremove -f` / `vgremove -f` / `pvremove -f` — force-destroys LVM metadata](#zc1537)
- [ZC1538: Error on `zpool destroy -f` / `zfs destroy -rR` — recursive ZFS destruction](#zc1538)
- [ZC1539: Warn on `parted -s <disk> <destructive-op>` — script mode bypasses confirmation](#zc1539)
- [ZC1540: Error on `cryptsetup erase` / `luksErase` — destroys LUKS header, data unrecoverable](#zc1540)
- [ZC1541: Error on `apk add --allow-untrusted` — installs unsigned Alpine packages](#zc1541)
- [ZC1542: Error on `snap install --dangerous` — installs unsigned snap](#zc1542)
- [ZC1543: Warn on `go install pkg@latest` / `cargo install --git <url>` without rev pin](#zc1543)
- [ZC1544: Warn on `dnf copr enable` / `add-apt-repository ppa:` — unvetted third-party repo](#zc1544)
- [ZC1545: Warn on `docker system prune -af --volumes` — drops unused volumes too](#zc1545)
- [ZC1546: Warn on `kubectl delete --force --grace-period=0` — skips PreStop, corrupts state](#zc1546)
- [ZC1547: Warn on `kubectl apply --prune --all` — deletes resources missing from manifest](#zc1547)
- [ZC1548: Warn on `helm install/upgrade --disable-openapi-validation` — skips schema check](#zc1548)
- [ZC1549: Error on `unzip -d /` / `unzip -o ... -d /` — extract archive into filesystem root](#zc1549)
- [ZC1550: Warn on `apt-mark hold <pkg>` — pins a package, blocks security updates](#zc1550)
- [ZC1551: Warn on `helm install/upgrade --skip-crds` — chart CRs land before their CRDs](#zc1551)
- [ZC1552: Warn on `openssl dhparam <2048` / `genrsa <2048` — weak key/parameter size](#zc1552)
- [ZC1553: Style: use Zsh `${(U)var}` / `${(L)var}` instead of `tr '\[:lower:\]' '\[:upper:\]'`](#zc1553)
- [ZC1554: Warn on `unzip -o` / `tar ... --overwrite` — silent overwrite during extract](#zc1554)
- [ZC1555: Error on `chmod` / `chown` on `/etc/shadow` or `/etc/sudoers` (managed files)](#zc1555)
- [ZC1556: Error on `openssl enc -des` / `-rc4` / `-3des` — broken symmetric cipher](#zc1556)
- [ZC1557: Error on `kubeadm reset -f` / `--force` — wipes Kubernetes control-plane state](#zc1557)
- [ZC1558: Warn on `usermod -aG wheel\|sudo\|root\|adm` — silent privilege group escalation](#zc1558)
- [ZC1559: Warn on `ssh-copy-id -f` / `-o StrictHostKeyChecking=no` — trust-on-first-use key push](#zc1559)
- [ZC1560: Error on `pip install --trusted-host` — accepts MITM / plain-HTTP PyPI index](#zc1560)
- [ZC1561: Error on `systemctl isolate rescue.target` / `emergency.target` from a script](#zc1561)
- [ZC1562: Warn on `env -u PATH` / `-u LD_LIBRARY_PATH` — clears security-relevant env](#zc1562)
- [ZC1563: Warn on `swapoff -a` — disables swap (memory pressure, potential OOM)](#zc1563)
- [ZC1564: Warn on `date -s` / `timedatectl set-time` — manual clock change breaks TLS / cron](#zc1564)
- [ZC1565: Style: use `command -v` instead of `whereis` / `locate` for command existence](#zc1565) · auto-fix
- [ZC1566: Error on `gem install -P NoSecurity\|LowSecurity` / `--trust-policy NoSecurity`](#zc1566)
- [ZC1567: Warn on `python -m http.server` without `--bind 127.0.0.1` — serves to all interfaces](#zc1567)
- [ZC1568: Error on `useradd -o` / `usermod -o` — allows non-unique UID (alias user)](#zc1568)
- [ZC1569: Error on `nvme format -s1` / `-s2` — cryptographic or full-block SSD erase](#zc1569)
- [ZC1570: Warn on `smbclient -N` / `mount.cifs guest` — anonymous SMB share access](#zc1570)
- [ZC1571: Style: `ntpdate` is deprecated — use `chronyc makestep` / `systemd-timesyncd`](#zc1571)
- [ZC1572: Warn on `docker run -e PASSWORD=<value>` — secret in container env / inspect](#zc1572)
- [ZC1573: Warn on `chattr -i` / `chattr -a` — removes immutable / append-only attribute](#zc1573)
- [ZC1574: Warn on `git config credential.helper store` — plaintext credentials on disk](#zc1574)
- [ZC1575: Error on `aws configure set aws_secret_access_key <value>` — secret on cmdline](#zc1575)
- [ZC1576: Warn on `terraform apply -target=...` — cherry-pick apply bypasses dependencies](#zc1576)
- [ZC1577: Warn on `dig <name> ANY` — deprecated query type (RFC 8482)](#zc1577)
- [ZC1578: Warn on `ssh-keygen -b <2048` for RSA / DSA — weak SSH key](#zc1578)
- [ZC1579: Warn on `curl --retry-all-errors` without `--max-time` — hammers endpoint on failure](#zc1579)
- [ZC1580: Warn on `go build -ldflags "-X main.<SECRET>=..."` — secret embedded in binary](#zc1580)
- [ZC1581: Warn on `ssh -o PubkeyAuthentication=no` / `-o PasswordAuthentication=yes`](#zc1581)
- [ZC1582: Warn on `bash -x` / `sh -x` / `zsh -x` — traces every command, leaks secrets](#zc1582)
- [ZC1583: Warn on `find ... -delete` without `-maxdepth` — unbounded recursive delete](#zc1583)
- [ZC1584: Warn on `sudo -E` / `--preserve-env` — carries caller env into root shell](#zc1584)
- [ZC1585: Warn on `ufw allow from any` / `firewall-cmd --add-source=0.0.0.0/0`](#zc1585)
- [ZC1586: Style: `chkconfig` / `update-rc.d` / `insserv` — SysV init relics, use `systemctl`](#zc1586)
- [ZC1587: Warn on `modprobe -r` / `rmmod` from scripts — unloading active kernel modules](#zc1587)
- [ZC1588: Error on `nsenter --target 1` — joins host init namespaces (container escape)](#zc1588)
- [ZC1589: Warn on `trap 'set -x' ERR/RETURN/EXIT/ZERR` — trace hook leaks env to stderr](#zc1589)
- [ZC1590: Error on `sshpass -p SECRET` — password in process list and history](#zc1590)
- [ZC1591: Style: use Zsh `print -l` / `${(F)array}` instead of `printf '%s\n' "${array\[@\]}"`](#zc1591)
- [ZC1592: Warn on `faillock --reset` / `pam_tally2 -r` — clears failed-auth counter](#zc1592)
- [ZC1593: Error on `blkdiscard` — issues TRIM/DISCARD across the whole device (data loss)](#zc1593)
- [ZC1594: Warn on `docker/podman run --security-opt=systempaths=unconfined` — unhides host kernel knobs](#zc1594)
- [ZC1595: Warn on `setfacl -m u:nobody:... / o::rwx` — ACL grants that bypass `chmod` scrutiny](#zc1595)
- [ZC1596: Style: `emulate sh/bash/ksh` without `-L` — flips options for the whole shell](#zc1596)
- [ZC1597: Warn on `systemd-run -p User=root` — launches arbitrary command with root privileges](#zc1597)
- [ZC1598: Error on `chmod` with world-write bit on a sensitive `/dev/` node](#zc1598)
- [ZC1599: Warn on `ldconfig -f PATH` outside `/etc/` — attacker-writable loader cache](#zc1599)
- [ZC1600: Warn on bare `chroot DIR CMD` — missing `--userspec=` keeps uid 0 inside the jail](#zc1600)
- [ZC1601: Warn on `ethtool -s $IF wol <g\|u\|m\|b\|a>` — enables remote Wake-on-LAN](#zc1601)
- [ZC1602: Warn on `setopt KSH_ARRAYS` / `SH_WORD_SPLIT` — flips Zsh core semantics shell-wide](#zc1602)
- [ZC1603: Warn on `gdb -p PID` / `ltrace -p PID` — live attach reads target memory](#zc1603)
- [ZC1604: Warn on `source <glob>` / `. <glob>` — loads every match; one bad file = code exec](#zc1604)
- [ZC1605: Error on `debugfs -w DEV` — write-mode filesystem debugger bypasses journal](#zc1605)
- [ZC1606: Warn on `mkdir -m NNN` / `install -m NNN` with world-write bit (no sticky)](#zc1606)
- [ZC1607: Warn on `git config safe.directory '*'` — disables CVE-2022-24765 protection](#zc1607)
- [ZC1608: Warn on `find -exec sh -c '... {} ...'` — filename in quoted script is injectable](#zc1608)
- [ZC1609: Warn on `aa-disable` / `aa-complain` / `apparmor_parser -R` — disables AppArmor enforcement](#zc1609)
- [ZC1610: Warn on `curl -o /etc/...` / `wget -O /etc/...` — direct download to a system path](#zc1610)
- [ZC1611: Style: `${var^^}` / `${var,,}` — prefer Zsh `${(U)var}` / `${(L)var}` for case change](#zc1611)
- [ZC1612: Warn on `sysctl -w` disabling kernel hardening knobs](#zc1612)
- [ZC1613: Warn on reading SSH private-key files with `cat` / `less` / `grep` / `head`](#zc1613)
- [ZC1614: Error on `expect` script containing `password` / `passphrase`](#zc1614)
- [ZC1615: Style: use Zsh `$EPOCHREALTIME` / `$epochtime` instead of `date "+%s.%N"`](#zc1615)
- [ZC1616: Warn on `fsfreeze -f MOUNTPOINT` — filesystem stays frozen until `-u` runs](#zc1616)
- [ZC1617: Warn on `xargs -P 0` — unbounded parallelism risks CPU / fd / memory exhaustion](#zc1617)
- [ZC1618: Warn on `git commit --no-verify` / `git push --no-verify` — bypasses hooks](#zc1618)
- [ZC1619: Warn on `mount -t nfs/cifs/smb/sshfs` missing `nosuid` or `nodev`](#zc1619)
- [ZC1620: Error on `tee /etc/sudoers` / `/etc/sudoers.d/*` — writes without `visudo -cf`](#zc1620)
- [ZC1621: Warn on `tmux -S /tmp/SOCKET` — shared-path socket invites session hijack](#zc1621)
- [ZC1622: Style: `${var@U/L/Q/...}` — prefer Zsh `${(U)var}` / `${(L)var}` / `${(Q)var}` flags](#zc1622)
- [ZC1623: Warn on `kill -STOP PID` / `pkill -STOP` — target halts until `kill -CONT` runs](#zc1623)
- [ZC1624: Error on `az login -p` / `--password` — service-principal secret in process list](#zc1624)
- [ZC1625: Error on `rm --no-preserve-root` — disables GNU rm safeguard against `rm -rf /`](#zc1625)
- [ZC1626: Error on `helm install/upgrade --set KEY=VALUE` with secret-shaped key](#zc1626)
- [ZC1627: Warn on `crontab /tmp/FILE` — attacker-writable path installed as a user's cron](#zc1627)
- [ZC1628: Warn on `insmod` / `modprobe -f` — loads modules bypassing blacklist / signature checks](#zc1628)
- [ZC1629: Warn on `rsync --rsync-path='sudo rsync'` — hidden remote privilege escalation](#zc1629)
- [ZC1630: Warn on `php -S 0.0.0.0:PORT` — PHP dev server exposes CWD to all interfaces](#zc1630)
- [ZC1631: Error on `openssl ... -passin pass:SECRET` / `-passout pass:SECRET`](#zc1631)
- [ZC1632: Warn on `shred` — unreliable on journaled / CoW filesystems (ext4, btrfs, zfs)](#zc1632)
- [ZC1633: Error on `gpg --passphrase SECRET` — passphrase on cmdline](#zc1633)
- [ZC1634: Warn on `umask NNN` that fails to mask world-write — mask-inversion footgun](#zc1634)
- [ZC1635: Error on `mysql -pSECRET` / `--password=SECRET` — password in process list](#zc1635)
- [ZC1636: Warn on `virsh destroy DOMAIN` — force-stops VM (no graceful shutdown)](#zc1636)
- [ZC1637: Style: prefer Zsh `typeset -r NAME=value` over POSIX `readonly NAME=value`](#zc1637) · auto-fix
- [ZC1638: Error on `docker/podman build --build-arg SECRET=VALUE` — secret baked into image layer](#zc1638)
- [ZC1639: Error on `curl -H 'Authorization: ...'` — credential header in process list](#zc1639)
- [ZC1640: Style: `${!var}` Bash indirect expansion — prefer Zsh `${(P)var}`](#zc1640)
- [ZC1641: Error on `kubectl create secret --from-literal=...` / `--docker-password=...`](#zc1641)
- [ZC1642: Warn on `tshark -w FILE` / `dumpcap -w FILE` without `-Z user` — capture file owned by root](#zc1642)
- [ZC1643: Style: `$(cat file)` — use `$(<file)` to skip the fork / exec](#zc1643) · auto-fix
- [ZC1644: Error on `unzip -P SECRET` / `zip -P SECRET` — archive password in process list](#zc1644)
- [ZC1645: Style: `lsb_release` — prefer sourcing `/etc/os-release` (no dependency, no fork)](#zc1645)
- [ZC1646: Warn on `btrfs check --repair` / `xfs_repair -L` — last-resort recovery, may worsen damage](#zc1646)
- [ZC1647: Warn on `kubectl apply -f URL` — remote manifest applied without digest verification](#zc1647)
- [ZC1648: Error on `cp /dev/null /var/log/...` / `truncate -s 0 /var/log/...` — audit-log wipe](#zc1648)
- [ZC1649: Warn on `openssl req -days N` with N > 825 — long-validity certificate](#zc1649)
- [ZC1650: Warn on `setopt RM_STAR_SILENT` / `unsetopt RM_STAR_WAIT` — removes `rm *` prompt](#zc1650)
- [ZC1651: Warn on `docker/podman run -p 0.0.0.0:PORT:PORT` — explicit all-interfaces publish](#zc1651)
- [ZC1652: Warn on `ssh -Y` — trusted X11 forwarding grants full X-server access to remote clients](#zc1652)
- [ZC1653: Avoid `$BASHPID` — Bash-only; Zsh uses `$sysparams\[pid\]` from `zsh/system`](#zc1653)
- [ZC1654: Warn on `sysctl -p /tmp/...` — loading kernel tunables from attacker-writable path](#zc1654)
- [ZC1655: Warn on `read -n N` — Bash reads N chars; Zsh's `-n` means "drop newline"](#zc1655)
- [ZC1656: Error on `rsync -e 'ssh -o StrictHostKeyChecking=no'` — host-key verify disabled](#zc1656)
- [ZC1657: Warn on `semanage permissive -a <type>` — puts SELinux domain in permissive mode](#zc1657)
- [ZC1658: Warn on `curl -OJ` / `-J -O` — server-controlled output filename](#zc1658)
- [ZC1659: Warn on `fuser -k <path>` — kills every process holding the subtree open](#zc1659)
- [ZC1660: Style: `printf '%0Nd' $n` — prefer Zsh `${(l:N::0:)n}` left-zero-pad](#zc1660)
- [ZC1661: Error on `curl --cacert /dev/null` — empty trust store, any cert passes](#zc1661)
- [ZC1662: Error on `pkexec env VAR=VAL CMD` — controlled env crossed into the root session](#zc1662)
- [ZC1663: Warn on `tune2fs -c 0` / `-i 0` — disables periodic filesystem checks](#zc1663)
- [ZC1664: Error on `systemctl set-default rescue.target\|emergency.target` — persistent single-user boot](#zc1664)
- [ZC1665: Warn on `chrt -r` / `-f` — real-time scheduling class from a shell script](#zc1665)
- [ZC1666: Warn on `kubectl patch --type=json` — bypasses strategic-merge defaults](#zc1666)
- [ZC1667: Warn on `openssl enc` without `-pbkdf2` — legacy MD5-based key derivation](#zc1667)
- [ZC1668: Error on `aws iam attach-*-policy ... AdministratorAccess` — grants full AWS admin](#zc1668)
- [ZC1669: Warn on `git gc --prune=now` / `git reflog expire --expire=now` — deletes recovery window](#zc1669)
- [ZC1670: Warn on `setsebool -P` enabling memory-protection-relaxing SELinux boolean](#zc1670)
- [ZC1671: Error on `install -m 777` / `mkdir -m 777` — creates world-writable target](#zc1671)
- [ZC1672: Info: `chcon` writes an ephemeral SELinux label — next `restorecon` wipes it](#zc1672)
- [ZC1673: Style: `stty -echo` around `read` — prefer Zsh `read -s`](#zc1673)
- [ZC1674: Warn on `docker/podman run --oom-kill-disable` or `--oom-score-adj <= -500`](#zc1674)
- [ZC1675: Avoid Bash-only `export -f` / `export -n` — use Zsh `typeset -fx` / `typeset +x`](#zc1675) · auto-fix
- [ZC1676: Warn on `helm rollback --force` — recreates in-flight resources, corrupts rolling updates](#zc1676)
- [ZC1677: Warn on `trap 'set -x' DEBUG` — xtrace on every command leaks secrets](#zc1677)
- [ZC1678: Error on `borg init --encryption=none` — unencrypted backup repository](#zc1678)
- [ZC1679: Error on `gcloud ... add-iam-policy-binding ... --role=roles/owner` — GCP primitive admin](#zc1679)
- [ZC1680: Error on `ansible-playbook --vault-password-file=/tmp/...` — world-traversable vault key](#zc1680)
- [ZC1681: Error on `tar -P` / `--absolute-names` — archive absolute paths, can overwrite host files](#zc1681)
- [ZC1682: Error on `npm install --unsafe-perm` — npm lifecycle scripts keep root privileges](#zc1682)
- [ZC1683: Error on `npm/yarn/pnpm config set registry http://...` — plaintext package index](#zc1683)
- [ZC1684: Error on `redis-cli -a PASSWORD` — authentication password in process list](#zc1684)
- [ZC1685: Info: `sleep infinity` — container keep-alive pattern that ignores SIGTERM](#zc1685)
- [ZC1686: Warn on `compinit -C` / `compinit -u` — skips / ignores `$fpath` integrity checks](#zc1686)
- [ZC1687: Warn on `snap install --classic` / `--devmode` — weakens snap confinement](#zc1687)
- [ZC1688: Warn on `aws s3 sync --delete` — destination objects deleted when source diverges](#zc1688)
- [ZC1689: Error on `borg delete --force` — forced deletion of backup archives or repository](#zc1689)
- [ZC1690: Warn on `pip install git+<URL>` without a commit / tag pin](#zc1690)
- [ZC1691: Warn on `rsync --remove-source-files` — SRC deletion tied to optimistic success](#zc1691)
- [ZC1692: Error on `kexec -e` — jumps into a new kernel without reboot, no audit trail](#zc1692)
- [ZC1693: Warn on `ionice -c 1` — real-time I/O class starves every other disk consumer](#zc1693)
- [ZC1694: Warn on `ssh -A` / `-o ForwardAgent=yes` — remote host can reuse local keys](#zc1694)
- [ZC1695: Warn on `terraform state rm` / `state push` — surgery on shared state outside plan/apply](#zc1695)
- [ZC1696: Warn on `pnpm install --no-frozen-lockfile` / `yarn install --no-immutable` — CI lockfile drift](#zc1696)
- [ZC1697: Info: `cryptsetup open --allow-discards` — TRIM pass-through leaks free-sector map](#zc1697)
- [ZC1698: Warn on `fail2ban-client unban --all` / `stop` — wipes the active brute-force ban list](#zc1698)
- [ZC1699: Warn on `kubectl drain --delete-emptydir-data` — pod-local scratch data lost](#zc1699)
- [ZC1700: Error on `ldapsearch -w PASSWORD` / `ldapmodify -w PASSWORD` — bind DN password in process list](#zc1700)
- [ZC1701: Info: `dpkg -i FILE.deb` installs without automatic signature verification](#zc1701)
- [ZC1702: Warn on `dpkg-reconfigure` without a noninteractive frontend — hangs in CI](#zc1702)
- [ZC1703: Warn on `sysctl -w` disabling network-hardening knobs](#zc1703)
- [ZC1704: Error on `aws ec2 authorize-security-group-ingress --cidr 0.0.0.0/0` — port open to the internet](#zc1704)
- [ZC1705: Info: `awk -i inplace` is gawk-only — script breaks on mawk / BSD awk](#zc1705)
- [ZC1706: Error on `lvresize -L -SIZE` without `-r` — shrink without filesystem resize corrupts data](#zc1706)
- [ZC1707: Warn on `gpg --keyserver hkp://…` — plaintext keyserver fetch](#zc1707)
- [ZC1708: Error on `find -L ... -delete` / `-exec rm` — symlink follow into unintended trees](#zc1708)
- [ZC1709: Error on `htpasswd -b USER PASSWORD` — basic-auth password in process list](#zc1709)
- [ZC1710: Error on `journalctl --vacuum-size=1` / `--vacuum-time=1s` — journal-wipe pattern](#zc1710)
- [ZC1711: Error on `etcdctl del --prefix ""` / `--from-key ""` — wipes the entire keyspace](#zc1711)
- [ZC1712: Error on `vault kv put PATH password=…` — secret value in process list](#zc1712)
- [ZC1713: Error on `consul kv delete -recurse /` — wipes the entire Consul KV store](#zc1713)
- [ZC1714: Error on `gh repo delete --yes` / `gh release delete --yes` — bypassed confirmation](#zc1714)
- [ZC1715: Error on `read -p "prompt"` — Zsh `-p` reads from coprocess, not a prompt](#zc1715)
- [ZC1716: Use Zsh `$CPUTYPE` / `$MACHTYPE` instead of `uname -m` / `-p`](#zc1716)
- [ZC1717: Warn on `docker pull/push --disable-content-trust` — bypasses image signature checks](#zc1717) · auto-fix
- [ZC1718: Error on `gh secret set --body SECRET` / `-b SECRET` — secret in process list](#zc1718)
- [ZC1719: Warn on `git filter-branch` — deprecated since Git 2.24, use `git filter-repo`](#zc1719)
- [ZC1720: Use Zsh `$COLUMNS` / `$LINES` instead of `tput cols` / `tput lines`](#zc1720)
- [ZC1721: Error on `chmod NNN /dev/<node>` — world-writable device node is local privilege escalation](#zc1721)
- [ZC1722: Warn on `ssh-keyscan HOST >> known_hosts` — TOFU bypass, blind-trust new host key](#zc1722)
- [ZC1723: Error on `gpg --delete-secret-keys` / `--delete-key` — irreversible key destruction](#zc1723)
- [ZC1724: Warn on `pacman -Sy <pkg>` — partial upgrade, breaks dependency closure](#zc1724)
- [ZC1725: Error on `cargo --token TOKEN` / `npm --otp CODE` — registry credential in process list](#zc1725)
- [ZC1726: Error on `gcloud ... delete --quiet` — silent destruction of GCP resources](#zc1726)
- [ZC1727: Error on `curl/wget --proxy http://USER:PASS@HOST` — proxy credentials in argv](#zc1727)
- [ZC1728: Error on `pip install --index-url http://...` — plaintext index allows MITM](#zc1728)
- [ZC1729: Error on `ip route flush all` / `ip route del default` — script loses network connectivity](#zc1729)
- [ZC1730: Warn on `brew install --HEAD <pkg>` — pulls upstream HEAD, no version stability](#zc1730)
- [ZC1731: Error on `curl -d 'password=…'` / `wget --post-data='token=…'` — secret in argv](#zc1731)
- [ZC1732: Warn on `flatpak override --filesystem=host` — removes Flatpak sandbox isolation](#zc1732)
- [ZC1733: Error on `docker plugin install --grant-all-permissions` — accepts every requested cap](#zc1733)
- [ZC1734: Error on `cp/mv/tee` overwriting `/etc/passwd\|shadow\|group\|gshadow`](#zc1734)
- [ZC1735: Error on `efibootmgr -B` — deletes UEFI boot entry, may brick boot](#zc1735)
- [ZC1736: Error on `pulumi destroy --yes` / `up --yes` — silent infra mutation in CI](#zc1736)
- [ZC1737: Error on `wpa_passphrase SSID PASSWORD` — Wi-Fi passphrase in process list](#zc1737)
- [ZC1738: Error on `aws rds delete-db-instance --skip-final-snapshot` — DB destroyed unrecoverable](#zc1738)
- [ZC1739: Warn on `git submodule update --remote` — pulls upstream HEAD, breaks reproducibility](#zc1739)
- [ZC1740: Warn on `gh release upload --clobber` — silent overwrite of release asset](#zc1740)
- [ZC1741: Error on `mkpasswd PASSWORD` — clear-text password in process list](#zc1741)
- [ZC1742: Error on `mc alias set NAME URL ACCESS_KEY SECRET_KEY` — S3 keys in process list](#zc1742)
- [ZC1743: Warn on `npm audit fix --force` — accepts major-version dependency bumps silently](#zc1743)
- [ZC1744: Warn on `kubectl port-forward --address 0.0.0.0` — cluster port exposed to every interface](#zc1744)
- [ZC1745: Error on `poetry publish --password PASS` / `twine upload -p PASS` — registry secret in argv](#zc1745)
- [ZC1746: Error on `sysctl -w kernel.randomize_va_space=0\|1` — weakens or disables ASLR](#zc1746)
- [ZC1747: Error on `npm/yarn/pnpm --registry http://...` — plaintext registry allows MITM](#zc1747)
- [ZC1748: Error on `helm repo add NAME http://...` — plaintext chart repo allows MITM](#zc1748)
- [ZC1749: Error on `virsh undefine DOMAIN --remove-all-storage` — wipes VM disk images](#zc1749)
- [ZC1750: Error on `kubectl proxy --address 0.0.0.0` — cluster API proxy on every interface](#zc1750)
- [ZC1751: Error on `rpm/dnf/yum remove --nodeps` — bypasses dependency check, breaks dependents](#zc1751)
- [ZC1752: Error on `pvcreate/vgcreate/lvcreate -ff\|--yes` — force-init LVM over existing data](#zc1752)
- [ZC1753: Error on `rclone purge REMOTE:PATH` — bulk delete of every object under the remote path](#zc1753)
- [ZC1754: Error on `gh auth status -t` / `--show-token` — prints OAuth token to stdout](#zc1754)
- [ZC1755: Error on `gcloud sql users {create,set-password} --password PASS` — DB password in argv](#zc1755)
- [ZC1756: Error on `chmod NNN /run/docker.sock` — world access is root-equivalent privesc](#zc1756)
- [ZC1757: Warn on `gh auth refresh --scopes delete_repo\|admin:*` — token escalated to destructive perms](#zc1757)
- [ZC1758: Warn on `gh codespace delete --force` — destroys codespace with uncommitted work](#zc1758)
- [ZC1759: Error on `vault login TOKEN` / `login -method=… password=…` — credential in process list](#zc1759)
- [ZC1760: Warn on `openssl rand -hex\|-base64 N` with N < 16 — generated value too short](#zc1760)
- [ZC1761: Warn on `gh gist create --public` — file becomes world-visible and indexed on GitHub](#zc1761)
- [ZC1762: Error on `kubeadm join --discovery-token-unsafe-skip-ca-verification` — cluster CA not checked](#zc1762)
- [ZC1763: Error on `docker compose down -v` / `--volumes` — wipes named volumes (data loss)](#zc1763)
- [ZC1764: Warn on `git commit --no-verify` / `-n` — skips pre-commit and commit-msg hooks](#zc1764)
- [ZC1765: Error on `snap remove --purge SNAP` — skips the automatic data snapshot](#zc1765)
- [ZC1766: Error on `memcached -l 0.0.0.0` — memcached exposed on every interface](#zc1766)
- [ZC1767: Error on `mongod --bind_ip 0.0.0.0` — MongoDB exposed on every interface](#zc1767)
- [ZC1768: Error on `sqlcmd -P PASSWORD` / `bcp -P PASSWORD` — SQL Server password in argv](#zc1768)
- [ZC1769: Warn on `vagrant destroy --force` — VM destroyed without confirmation](#zc1769)
- [ZC1770: Warn on `gpg --always-trust` / `--trust-model always` — bypasses Web-of-Trust](#zc1770)
- [ZC1771: Warn on `alias -g` / `alias -s` — global and suffix aliases surprise script readers](#zc1771)
- [ZC1772: Error on `hdparm --security-erase` / `--trim-sector-ranges` — ATA-level data destruction](#zc1772)
- [ZC1773: Warn on `xargs` without `-r` / `--no-run-if-empty` — runs once on empty input](#zc1773) · auto-fix
- [ZC1774: Warn on `setopt GLOB_SUBST` — `$var` starts glob-expanding, user data becomes a pattern](#zc1774)
- [ZC1775: Warn on `timeout DURATION cmd` without `--kill-after` / `-k` — hang on SIGTERM-resistant child](#zc1775)
- [ZC1776: Error on `psql postgresql://user:secret@host/db` — password in argv via connection URI](#zc1776)
- [ZC1777: Error on `tee/cp/mv/install/dd` writing `/etc/ld.so.preload` — classic rootkit persistence](#zc1777)
- [ZC1778: Warn on `systemctl link /path/to/unit` — persistence from a mutable source path](#zc1778)
- [ZC1779: Error on `az role assignment create --role Owner\|Contributor\|User Access Administrator`](#zc1779)
- [ZC1780: Warn on `sysctl -w fs.protected_symlinks=0\|protected_hardlinks=0\|…` — TOCTOU guard disabled](#zc1780)
- [ZC1781: Error on `git clone https://user:token@host/...` — PAT in argv and git config](#zc1781)
- [ZC1782: Error on `flatpak remote-add --no-gpg-verify` — trust chain disabled for the repo](#zc1782)
- [ZC1783: Error on `podman system reset` / `nerdctl system prune -af --volumes` — wipes every container artifact](#zc1783)
- [ZC1784: Warn on `git config core.hooksPath /tmp/...` — hook execution from a mutable path](#zc1784)
- [ZC1785: Error on `ufw default allow` — flips host firewall from deny-by-default to allow-by-default](#zc1785)
- [ZC1786: Error on `mount.cifs ... -o password=SECRET` — SMB password in argv](#zc1786)
- [ZC1787: Warn on `setopt AUTO_CD` — bare word that names a directory silently changes `$PWD`](#zc1787)
- [ZC1788: Warn on `ssh -F /tmp/config` — config from a mutable path can pin `ProxyCommand` to arbitrary code](#zc1788)
- [ZC1789: Warn on `setopt CORRECT` / `CORRECT_ALL` — Zsh spellcheck silently rewrites script tokens](#zc1789)
- [ZC1790: Warn on `unsetopt PIPE_FAIL` — pipeline exit status reverts to last-command-only](#zc1790)
- [ZC1791: Error on `curl --unix-socket /var/run/docker.sock` — direct container-daemon API access](#zc1791)
- [ZC1792: Warn on `btrfs subvolume delete` / `btrfs device remove` — unrecoverable btrfs data loss](#zc1792)
- [ZC1793: Warn on `kubectl certificate approve CSR` — signs the identity baked into the CSR](#zc1793)
- [ZC1794: Error on `cosign verify --insecure-ignore-tlog` / `--allow-insecure-registry` — signature chain disabled](#zc1794)
- [ZC1795: Error on `git remote add NAME https://user:token@host/repo` — credentials persisted in `.git/config`](#zc1795)
- [ZC1796: Warn on `pg_restore --clean` / `-c` — drops existing DB objects before recreating](#zc1796)
- [ZC1797: Warn on `ip link set <iface> down` / `ifdown <iface>` — locks out remote admin on that path](#zc1797)
- [ZC1798: Warn on `ufw reset` — wipes every configured firewall rule](#zc1798)
- [ZC1799: Warn on `rclone sync SRC DST` without `--dry-run` — one-way mirror can wipe DST](#zc1799)
- [ZC1800: Warn on `pg_ctl stop -m immediate` — abrupt shutdown skips checkpoint, forces WAL recovery](#zc1800)
- [ZC1801: Warn on `fwupdmgr update` / `install` — mid-flash interruption can brick firmware](#zc1801)
- [ZC1802: Warn on `dnf history undo N` / `rollback N` — reverses transactions without compat check](#zc1802)
- [ZC1803: Error on `mysql --skip-ssl` / `psql sslmode=disable` — plaintext credentials on the wire](#zc1803)
- [ZC1804: Warn on `aws ec2 terminate-instances` / `delete-volume` / `delete-snapshot` — destructive cloud state change](#zc1804)
- [ZC1805: Warn on `aws cloudformation delete-stack` / `dynamodb delete-table` / `logs delete-log-group` / `kms schedule-key-deletion` — destructive AWS state change](#zc1805)
- [ZC1806: Warn on `zmv 'PAT' 'REP'` without `-n` / `-i` — silent bulk rename](#zc1806)
- [ZC1807: Warn on `gh api -X DELETE` — raw GitHub DELETE bypasses `gh` command confirmations](#zc1807)
- [ZC1808: Warn on `kubectl replace --force` — deletes + recreates resource, drops running pods](#zc1808)
- [ZC1809: Error on `gsutil rm -r gs://…` / `gsutil rb -f gs://…` — bulk GCS deletion](#zc1809)
- [ZC1810: Warn on `wget -r` / `--mirror` without `--level=N` — unbounded recursive download](#zc1810)
- [ZC1811: Error on `chown/chmod/chgrp --no-preserve-root` — disables GNU safeguard against recursive `/`](#zc1811)
- [ZC1812: Error on `aws ssm put-parameter --type SecureString --value SECRET` — plaintext in argv](#zc1812)
- [ZC1813: Warn on `cryptsetup luksFormat` / `reencrypt` — destructive LUKS header write](#zc1813)
- [ZC1814: Error on `dpkg --force-all` — enables every single `--force-*` option at once](#zc1814)
- [ZC1815: Warn on `systemctl restart NetworkManager` / `systemd-networkd` — drops the SSH session](#zc1815)
- [ZC1816: Warn on `docker/podman commit` — produces un-reproducible image, bakes in runtime state](#zc1816)
- [ZC1817: Warn on `git push --delete` / `git push -d` / `git push origin :branch` — remote branch removal](#zc1817)
- [ZC1818: Warn on `rsync --delete` without `--dry-run` — empty or wrong SRC wipes DST](#zc1818)
- [ZC1819: Warn on `xattr -d com.apple.quarantine` / `xattr -cr` — removes macOS Gatekeeper quarantine](#zc1819)
- [ZC1820: Warn on `netplan apply` — applies network config immediately with no rollback timer](#zc1820)
- [ZC1821: Error on `diskutil eraseDisk` / `secureErase` / `partitionDisk` — macOS storage reformat](#zc1821)
- [ZC1822: Error on `csrutil disable` / `spctl --master-disable` — disables macOS system integrity / Gatekeeper](#zc1822)
- [ZC1823: Warn on `keytool -import -noprompt` — Java trust store imports without fingerprint check](#zc1823)
- [ZC1824: Warn on `kubectl drain --disable-eviction` — bypasses PodDisruptionBudget via raw DELETE](#zc1824)
- [ZC1825: Warn on `scp -O` — forces legacy SCP wire protocol exposed to filename-injection CVEs](#zc1825)
- [ZC1826: Warn on `install -m u+s` / `g+s` — symbolic setuid/setgid bit applied at install time](#zc1826)
- [ZC1827: Error on `npm unpublish` — breaks every downstream that pinned the version](#zc1827)
- [ZC1828: Warn on `gcore PID` / `strace -p PID` — live ptrace attach dumps target memory](#zc1828)
- [ZC1829: Warn on `tailscale down` / `wg-quick down` / `nmcli con down` — drops the VPN that may carry the SSH session](#zc1829)
- [ZC1830: Warn on `unsetopt NOMATCH` — unmatched glob becomes the literal pattern, silent bugs](#zc1830)
- [ZC1831: Error on `systemctl stop\|disable\|mask ssh/sshd` — locks out the next remote login](#zc1831)
- [ZC1832: Warn on Zsh `limit coredumpsize unlimited` — setuid memory landing in core files](#zc1832)
- [ZC1833: Warn on `unsetopt WARN_CREATE_GLOBAL` — silent accidental-global bugs inside functions](#zc1833)
- [ZC1834: Error on `tc qdisc … root netem loss 100%` — hard blackhole on a live interface](#zc1834)
- [ZC1835: Warn on `smartctl -s off` — drive self-monitoring (SMART) disabled, silent failure](#zc1835)
- [ZC1836: Error on `helm uninstall --no-hooks` — skips pre-delete cleanup, orphaned state](#zc1836)
- [ZC1837: Error on `chmod` granting non-owner access to `/dev/kvm` / `/dev/mem` / `/dev/kmem` / `/dev/port`](#zc1837)
- [ZC1838: Warn on `setopt GLOB_DOTS` — bare `*` silently starts matching hidden files](#zc1838)
- [ZC1839: Warn on `timedatectl set-ntp false` / disabling `systemd-timesyncd` / `chronyd`](#zc1839)
- [ZC1840: Error on `openssl enc -k PASSWORD` — legacy flag embeds secret in argv](#zc1840)
- [ZC1841: Error on `curl --proxy-insecure` — TLS verification disabled on the proxy hop](#zc1841)
- [ZC1842: Warn on `setopt CDABLE_VARS` — `cd NAME` silently falls back to `cd $NAME`](#zc1842)
- [ZC1843: Warn on `docker/podman run --cgroup-parent=/system.slice\|/init.scope\|/` — container escapes engine limits](#zc1843)
- [ZC1844: Warn on `logger -p local0.info\|local7.notice\|…` — unreserved facility often uncollected](#zc1844)
- [ZC1845: Warn on `setopt PATH_DIRS` — slash-bearing command names fall back to `$PATH` lookup](#zc1845)
- [ZC1846: Warn on `certbot … --force-renewal` — bypasses ACME rate-limit safety](#zc1846)
- [ZC1847: Warn on `setopt CHASE_LINKS` — every `cd` silently swaps symlink paths for the real inode](#zc1847)
- [ZC1848: Warn on `ssh -o CheckHostIP=no` — DNS-spoof warning for known hosts silenced](#zc1848)
- [ZC1849: Warn on `setopt ALL_EXPORT` — every later `var=value` silently becomes `export var=value`](#zc1849)
- [ZC1850: Warn on `ssh -o LogLevel=QUIET` — silences security-relevant ssh diagnostics](#zc1850)
- [ZC1851: Warn on `unsetopt FUNCTION_ARGZERO` — `$0` inside a function stops reporting the function name](#zc1851)
- [ZC1852: Error on `firewall-cmd --panic-on` — firewalld drops every packet, kills the SSH session](#zc1852)
- [ZC1853: Warn on `setopt MARK_DIRS` — glob-matched directories gain a silent trailing `/`](#zc1853)
- [ZC1854: Error on `yum-config-manager --add-repo http://…` / `zypper addrepo http://…` — plaintext repo allows MITM](#zc1854)
- [ZC1855: Avoid `$GROUPS` — Bash-only array; Zsh exposes supplementary groups as `$groups`](#zc1855)
- [ZC1856: Warn on `unset arr\[N\]` — Zsh does not delete the array element, the array keeps its length](#zc1856)
- [ZC1857: Error on `cloud-init clean` — wipes boot state, next reboot re-provisions the host](#zc1857)
- [ZC1858: Error on `ssh -c 3des-cbc\|arcfour\|blowfish-cbc` — weak cipher forced on the tunnel](#zc1858)
- [ZC1859: Warn on `unsetopt MULTIOS` — `cmd >a >b` silently keeps only the last redirection](#zc1859)
- [ZC1860: Warn on `hostnamectl set-hostname NEW` — caches and certs still reference the old name](#zc1860)
- [ZC1861: Warn on `setopt OCTAL_ZEROES` — leading-zero integers silently reinterpret as octal](#zc1861)
- [ZC1862: Warn on `ssh-keygen -R HOST` — deletes a known-hosts entry, next `ssh` re-trusts silently](#zc1862)
- [ZC1863: Warn on `unsetopt CASE_GLOB` — globs silently go case-insensitive across the script](#zc1863)
- [ZC1864: Error on `mount -o remount,exec` — re-enables exec on a previously `noexec` mount](#zc1864)
- [ZC1865: Warn on `unsetopt CASE_MATCH` — `\[\[ =~ \]\]` and pattern tests quietly fold case](#zc1865)
- [ZC1866: Warn on `docker exec -u 0` — bypasses the image's non-root `USER` directive](#zc1866)
- [ZC1867: Warn on `unsetopt GLOB` — pattern expansion turned off, `rm *.log` tries the literal filename](#zc1867)
- [ZC1868: Error on `gcloud config set auth/disable_ssl_validation true` — disables TLS on every later gcloud call](#zc1868)
- [ZC1869: Warn on `setopt RC_EXPAND_PARAM` — brace-adjacent array expansion silently distributes](#zc1869)
- [ZC1870: Warn on `setopt GLOB_ASSIGN` — RHS of `var=pattern` silently glob-expands](#zc1870)
- [ZC1871: Warn on `setopt IGNORE_BRACES` — brace expansion stops working script-wide](#zc1871)
- [ZC1872: Error on `badblocks -w` — destructive write-mode pattern test wipes the device](#zc1872)
- [ZC1873: Warn on `setopt ERR_RETURN` — functions silently bail out on the first non-zero exit](#zc1873)
- [ZC1874: Warn on `sshuttle -r HOST 0/0` — every outbound packet tunneled through the jump host](#zc1874)
- [ZC1875: Warn on `setopt RC_QUOTES` — `''` inside single quotes flips from empty-concat to literal apostrophe](#zc1875)
- [ZC1876: Warn on `cargo publish --allow-dirty` — publishes the crate with uncommitted local changes](#zc1876)
- [ZC1877: Warn on `unsetopt SHORT_LOOPS` — short-form `for`/`while` bodies stop parsing](#zc1877)
- [ZC1878: Warn on `kubectl apply --force-conflicts` — steals ownership of fields managed by other controllers](#zc1878)
- [ZC1879: Warn on `unsetopt BAD_PATTERN` — malformed glob patterns silently pass through as literals](#zc1879)
- [ZC1880: Warn on `kubectl annotate\|label --overwrite` — silently rewrites controller signals](#zc1880)
- [ZC1881: Warn on `unsetopt MULTIBYTE` — `${#str}`, substring, and `\[\[ =~ \]\]` stop counting characters](#zc1881)
- [ZC1882: Warn on `sudo -s` / `sudo su` / `sudo bash` — spawns an interactive root shell from a script](#zc1882)
- [ZC1883: Warn on `setopt PATH_SCRIPT` — `. ./script.sh` silently falls back to `$PATH` lookup](#zc1883)
- [ZC1884: Error on `curl/wget https://...?apikey=...` — credential in URL query string](#zc1884)
- [ZC1885: Warn on `setopt CSH_NULL_GLOB` — unmatched globs drop instead of erroring when any sibling matches](#zc1885)
- [ZC1886: Error on `tee/cp/mv/install/dd` writing system shell-init files — persistent privesc surface](#zc1886)
- [ZC1887: Warn on `setopt POSIX_TRAPS` — EXIT/ZERR traps change scope and no longer fire on function return](#zc1887)
- [ZC1888: Warn on `aws iam create-access-key` — mints long-lived static AWS credentials](#zc1888)
- [ZC1889: Error on `skopeo copy --src-tls-verify=false` / `--dest-tls-verify=false` — MITM on image copy](#zc1889)
- [ZC1890: Error on `kadmin -w PASS` / `kinit` with password arg — Kerberos password in argv](#zc1890)
- [ZC1891: Error on `kubectl config view --raw` — prints the full kubeconfig with client keys](#zc1891)
- [ZC1892: Error on `install -m 4755\|6755\|2755` — sets setuid/setgid bit at install time](#zc1892)
- [ZC1893: Warn on `unsetopt BARE_GLOB_QUAL` — `*(N)` null-glob qualifier stops being special](#zc1893)
- [ZC1894: Error on `conntrack -F` / `--flush` — every tracked connection (including SSH) is reset](#zc1894)
- [ZC1895: Warn on `setopt NUMERIC_GLOB_SORT` — glob output switches from lexicographic to numeric order](#zc1895)
- [ZC1896: Error on `docker/podman run -v /proc:…\|/sys:…` — bind-mounts host kernel interfaces into container](#zc1896)
- [ZC1897: Warn on `setopt SH_GLOB` — Zsh-specific glob patterns (`*(N)`, `<1-10>`, alternation) stop parsing](#zc1897)
- [ZC1898: Error on `gpg --export-secret-keys` — private-key material leaks to stdout](#zc1898)
- [ZC1899: Error on `mokutil --disable-validation` — turns UEFI Secure Boot off at the shim](#zc1899)
- [ZC1900: Warn on `curl --location-trusted` — Authorization/cookies forwarded across redirects](#zc1900)
- [ZC1901: Warn on `setopt POSIX_BUILTINS` — flips `command`/special-builtin semantics](#zc1901)
- [ZC1902: Error on `ln -s /dev/null <logfile>` — silently discards audit or history writes](#zc1902)
- [ZC1903: Error on `tee /etc/sudoers*` — appends a rule that bypasses `visudo -c` validation](#zc1903)
- [ZC1904: Warn on `setopt KSH_GLOB` — reinterprets `*(pattern)` and breaks Zsh glob qualifiers](#zc1904)
- [ZC1905: Warn on `ssh -g -L …` — local forward bound on all interfaces, not just loopback](#zc1905)
- [ZC1906: Warn on `setopt POSIX_CD` — changes when `cd` / `pushd` consult `CDPATH`](#zc1906)
- [ZC1907: Warn on `sysctl -w fs.protected_*=0` / `fs.suid_dumpable=2` — disables /tmp-race safeguards](#zc1907)
- [ZC1908: Warn on `setopt MAGIC_EQUAL_SUBST` — enables tilde/param expansion on `key=value` args](#zc1908)
- [ZC1909: Warn on `kexec -l` / `-e` — jumps to an alternate kernel, bypasses bootloader and Secure Boot](#zc1909)
- [ZC1910: Warn on `setopt GLOB_STAR_SHORT` — makes bare `**` recurse instead of matching literal](#zc1910)
- [ZC1911: Warn on `umount -l` / `--lazy` — detach now, leaves open fds pointing at a ghost mount](#zc1911)
- [ZC1912: Warn on `dhclient -r` / `dhclient -x` / `dhcpcd -k` — drops the lease and breaks network](#zc1912)
- [ZC1913: Warn on `setopt ALIAS_FUNC_DEF` — re-enables defining functions with aliased names](#zc1913)
- [ZC1914: Warn on `curl --doh-url …` / `--dns-servers …` — overrides system resolver per-request](#zc1914)
- [ZC1915: Error on `mdadm --zero-superblock` / `--stop` — drops RAID metadata or live array](#zc1915)
- [ZC1916: Warn on `setopt NULL_GLOB` — every unmatched glob silently expands to nothing](#zc1916)
- [ZC1917: Info on `iw dev $IF scan` / `iwlist $IF scan` — active WiFi scan from a script](#zc1917)
- [ZC1918: Warn on `setopt HIST_SUBST_PATTERN` — `!:s/old/new/` silently switches to pattern matching](#zc1918)
- [ZC1919: Warn on `ss -K` / `ss --kill` — terminates every socket that matches the filter](#zc1919)
- [ZC1920: Warn on `setopt VERBOSE` — every executed command is echoed to stderr](#zc1920)
- [ZC1921: Warn on `systemctl kill -s KILL` / `--signal=SIGKILL` — skips `ExecStop=`, leaks resources](#zc1921)
- [ZC1922: Error on `rpm --import http://…` / `rpmkeys --import ftp://…` — plaintext GPG key fetch](#zc1922)
- [ZC1923: Warn on `setopt PRINT_EXIT_VALUE` — every non-zero exit leaks a status line to stderr](#zc1923)
- [ZC1924: Warn on `virt-cat` / `virt-copy-out` / `guestfish` / `guestmount` — reads guest disk from host](#zc1924)
- [ZC1925: Warn on `unsetopt EQUALS` — disables `=cmd` path expansion and tilde-after-colon](#zc1925)
- [ZC1926: Warn on `telinit 0/1/6` / `init 0/1/6` — SysV runlevel change halts, reboots, or isolates the host](#zc1926)
- [ZC1927: Error on `xfreerdp /p:SECRET` / `rdesktop -p SECRET` — RDP password visible in argv](#zc1927)
- [ZC1928: Warn on `setopt SHARE_HISTORY` — every session writes its history into every sibling session](#zc1928)
- [ZC1929: Warn on `cpio -i` / `--extract` without `--no-absolute-filenames` — archive writes outside CWD](#zc1929)
- [ZC1930: Warn on `unsetopt HASH_CMDS` — every command invocation re-walks `$PATH`](#zc1930)
- [ZC1931: Warn on `ip netns delete $NS` / `netns del` — drops the whole network namespace](#zc1931)
- [ZC1932: Warn on `unsetopt GLOBAL_EXPORT` — `typeset -x` in a function stops leaking to outer scope](#zc1932)
- [ZC1933: Error on `ipvsadm -C` / `--clear` — wipes every IPVS virtual service, drops load balancer](#zc1933)
- [ZC1934: Warn on `setopt AUTO_NAME_DIRS` — any absolute-path parameter becomes a `~name` alias](#zc1934)
- [ZC1935: Warn on `apt autoremove --purge` / `dnf autoremove` — deletes auto-installed deps and their config](#zc1935)
- [ZC1936: Warn on `setopt POSIX_ALIASES` — aliases on reserved words (`if`, `for`, …) stop expanding](#zc1936)
- [ZC1937: Warn on `tmux kill-server` / `tmux kill-session` — tears down every detached process inside](#zc1937)
- [ZC1938: Warn on `setopt POSIX_JOBS` — flips job-control semantics and `%n` scope](#zc1938)
- [ZC1939: Error on `reboot -f` / `halt -f` / `poweroff -f` — skips shutdown sequence, no graceful service stop](#zc1939)
- [ZC1940: Warn on `setopt POSIX_ARGZERO` — `$0` no longer changes to the function name inside functions](#zc1940)
- [ZC1941: Error on `restic init --insecure-no-password` — creates an unencrypted backup repository](#zc1941)
- [ZC1942: Warn on `setopt CLOBBER_EMPTY` — `>file` still overwrites zero-length files under `NO_CLOBBER`](#zc1942)
- [ZC1943: Warn on `systemd-nspawn -b` / `--boot` — runs a full init inside a possibly untrusted rootfs](#zc1943)
- [ZC1944: Warn on `setopt IGNORE_EOF` — Ctrl-D no longer exits the shell, masking runaway pipelines](#zc1944)
- [ZC1945: Warn on `bpftrace -e` / `bpftool prog load` — loads in-kernel eBPF from a script](#zc1945)
- [ZC1946: Warn on `unsetopt HUP` — background jobs keep running after shell exit](#zc1946)
- [ZC1947: Error on `ip xfrm state flush` / `ip xfrm policy flush` — tears down every IPsec SA and policy](#zc1947)
- [ZC1948: Error on `ipmitool -P PASS` / `-E` — BMC password visible in argv](#zc1948)
- [ZC1949: Error on `rmmod -f` / `rmmod --force` — bypasses refcount, can panic the kernel](#zc1949)
- [ZC1950: Error on `tune2fs -O ^has_journal` / `-m 0` — removes journal or root reserve](#zc1950)
- [ZC1951: Error on `ceph osd pool delete … --yes-i-really-really-mean-it` — automates Ceph's double-safety phrase](#zc1951)
- [ZC1952: Error on `zfs set sync=disabled` — `fsync()` becomes a no-op, crash loses unflushed writes](#zc1952)
- [ZC1953: Warn on `mount --make-shared` / `--make-rshared` — flips propagation, container-escape vector](#zc1953)
- [ZC1954: Warn on `setfattr -n security.capability\|security.selinux\|security.ima` — bypasses `setcap`/`chcon`](#zc1954)
- [ZC1955: Warn on `rfkill block all` / `block wifi\|bluetooth\|wwan` — disables every radio, cuts wireless](#zc1955)
- [ZC1956: Error on `tailscale up --auth-key=SECRET` — single-use join key visible in argv](#zc1956)
- [ZC1957: Warn on `lvchange -an` / `vgchange -an` — deactivates a live LV/VG, risks mounted-fs corruption](#zc1957)
- [ZC1958: Warn on `helm upgrade --force` — delete-and-recreate resources, drops running pods](#zc1958)
- [ZC1959: Warn on `trivy … --skip-db-update` / `--skip-update` — scans against a stale vulnerability DB](#zc1959)
- [ZC1960: Warn on `az vm run-command invoke` / `aws ssm send-command` — arbitrary commands on remote VM](#zc1960)
- [ZC1961: Warn on `gcloud iam service-accounts keys create` — mints a long-lived service-account JSON key](#zc1961)
- [ZC1962: Warn on `kustomize build --load-restrictor=LoadRestrictionsNone` — path-traversal in overlays](#zc1962)
- [ZC1963: Warn on `npx pkg` / `pnpm dlx pkg` / `bunx pkg` without a version pin — runs latest registry code](#zc1963)
- [ZC1964: Warn on `uvx pkg` / `uv tool run pkg` / `pipx run pkg` without a version pin — runs latest PyPI release](#zc1964)
- [ZC1965: Error on `systemd-cryptenroll --wipe-slot=all` — wipes every LUKS key slot](#zc1965)
- [ZC1966: Error on `zpool import -f` / `zpool export -f` — forced ZFS pool op bypasses hostid/txg checks](#zc1966)
- [ZC1967: Warn on `setopt PROMPT_SUBST` — expansions inside `$PROMPT` evaluate command substitution every redraw](#zc1967)
- [ZC1968: Warn on `dnf versionlock add` / `yum versionlock add` — pins RPM, blocks CVE updates](#zc1968)
- [ZC1969: Warn on `zsh -f` / `zsh -d` — skips `/etc/zsh*` and `~/.zsh*` startup files](#zc1969)
- [ZC1970: Warn on `losetup -P` / `kpartx -a` / `partprobe` on untrusted image — runs kernel partition parser](#zc1970)
- [ZC1971: Warn on `unsetopt GLOBAL_RCS` / `setopt NO_GLOBAL_RCS` — skips `/etc/zprofile`, `/etc/zshrc`, `/etc/zlogin`, `/etc/zlogout`](#zc1971)
- [ZC1972: Error on `dmsetup remove_all` / `dmsetup remove -f` — tears down live LVM/LUKS/multipath mappings](#zc1972)
- [ZC1973: Warn on `setopt POSIX_IDENTIFIERS` — restricts parameter names to ASCII, breaks Unicode `$var`](#zc1973)
- [ZC1974: Error on `ipset flush` / `ipset destroy` — nukes named sets referenced by iptables/nft rules](#zc1974)
- [ZC1975: Warn on `unsetopt EXEC` / `setopt NO_EXEC` — parser keeps scanning, commands stop running](#zc1975)
- [ZC1976: Error on `exportfs -au` / `exportfs -u` — unexports live NFS shares, clients get `ESTALE`](#zc1976)
- [ZC1977: Warn on `setopt CHASE_DOTS` — `cd ..` physically resolves before walking up, breaking logical paths](#zc1977)
- [ZC1978: Warn on `tftp` — cleartext, unauthenticated UDP transfer](#zc1978)
- [ZC1979: Warn on `setopt HIST_FCNTL_LOCK` — `fcntl()` lock on NFS `$HISTFILE` stalls or deadlocks](#zc1979)
- [ZC1980: Error on `udevadm trigger --action=remove` — replays `remove` uevents, detaches live devices](#zc1980)
- [ZC1981: Warn on `exec -a NAME cmd` — replaces `argv\[0\]`, hides the real binary from `ps`](#zc1981)
- [ZC1982: Error on `ipcrm -a` — removes every SysV IPC object, breaks Postgres/Oracle/shm apps](#zc1982)
- [ZC1983: Warn on `setopt CSH_JUNKIE_QUOTES` — single/double-quoted strings that span lines become errors](#zc1983)
- [ZC1984: Error on `sgdisk -Z` / `sgdisk -o` — erases the GPT partition table on the target disk](#zc1984)
- [ZC1985: Warn on `setopt SH_FILE_EXPANSION` — expansion order flips from Zsh-native to sh/bash, `~` leaks](#zc1985)
- [ZC1986: Warn on `touch -d` / `-t` / `-r` — explicit timestamp write is a common antiforensics pattern](#zc1986)
- [ZC1987: Warn on `setopt BRACE_CCL` — `{a-z}` expands to each character instead of staying literal](#zc1987)
- [ZC1988: Error on `nsupdate -y HMAC:NAME:SECRET` — TSIG key visible in argv and shell history](#zc1988)
- [ZC1989: Warn on `setopt REMATCH_PCRE` — `\[\[ =~ \]\]` regex flips from POSIX ERE to PCRE, changes semantics](#zc1989)
- [ZC1990: Warn on `openssl passwd -crypt` / `-1` / `-apr1` — obsolete password hash formats](#zc1990)
- [ZC1991: Warn on `setopt CSH_NULLCMD` — bare `> file` raises an error instead of running `$NULLCMD`](#zc1991)
- [ZC1992: Warn on `pkexec cmd` — PolicyKit privilege elevation is historically bug-prone and hard to audit from scripts](#zc1992)
- [ZC1993: Warn on `setopt KSH_TYPESET` — `typeset var=$val` starts word-splitting the RHS](#zc1993)
- [ZC1994: Error on `lvreduce -f` / `lvreduce -y` — shrinks the LV without checking the filesystem above](#zc1994)
- [ZC1995: Warn on `unsetopt BGNICE` — background jobs run at full interactive priority, starve the foreground](#zc1995)
- [ZC1996: Warn on `unshare -U` / `-r` — unprivileged user namespace maps caller to root inside the NS](#zc1996)
- [ZC1997: Warn on `setopt HIST_NO_FUNCTIONS` — function definitions skipped from `$HISTFILE`, breaks forensic trail](#zc1997)
- [ZC1998: Error on `tpm2_clear` / `tpm2 clear` — wipes TPM storage hierarchy, kills every sealed key](#zc1998)
- [ZC1999: Error on `setopt AUTO_NAMED_DIRS` — unknown option, typo of `AUTO_NAME_DIRS`](#zc1999)
- [ZC2000: Error on `kubectl taint nodes $NODE key=value:NoExecute` — evicts every non-tolerating pod off the node](#zc2000)
- [ZC2001: Warn on `unsetopt EVAL_LINENO` — `$LINENO` inside `eval` stops tracking source, stack traces go blank](#zc2001)
- [ZC2002: Error on `crictl rmi -a` / `crictl rm -af` — wipes every image/container on the Kubernetes node](#zc2002)
- [ZC2003: Warn on `setopt KSH_ZERO_SUBSCRIPT` — `$arr\[0\]` stops aliasing the first element](#zc2003)

---

<a id="zc1001"></a>
### ZC1001 — Use ${} for array element access

**Severity:** `style`  
**Auto-fix:** `yes`

In Zsh, accessing array elements with `$my_array[1]` doesn't work as expected. It tries to access an element from an array named `my_array[1]`. The correct way to access an array element is to use `${my_array[1]}`.

Disable by adding `ZC1001` to `disabled_katas` in `.zshellcheckrc`.

---

<a id="zc1002"></a>
### ZC1002 — Use $(...) instead of backticks

**Severity:** `style`  
**Auto-fix:** `yes`

Backticks are the old-style command substitution. $(...) is nesting-safe, easier to read, and generally preferred.

Disable by adding `ZC1002` to `disabled_katas` in `.zshellcheckrc`.

---

<a id="zc1003"></a>
### ZC1003 — Use `((...))` for arithmetic comparisons instead of `[` or `test`

**Severity:** `style`  
**Auto-fix:** `yes`

Bash/Zsh have a dedicated arithmetic context `((...))` which is cleaner and faster than `[` or `test` for numeric comparisons.

Disable by adding `ZC1003` to `disabled_katas` in `.zshellcheckrc`.

---

<a id="zc1004"></a>
### ZC1004 — Use `return` instead of `exit` in functions

**Severity:** `warning`  
**Auto-fix:** `yes`

Using `exit` in a function terminates the entire shell, which is often unintended in interactive sessions or sourced scripts. Use `return` to exit the function.

Disable by adding `ZC1004` to `disabled_katas` in `.zshellcheckrc`.

---

<a id="zc1005"></a>
### ZC1005 — Use whence instead of which

**Severity:** `info`  
**Auto-fix:** `yes`

The `which` command is an external command and may not be available on all systems. The `whence` command is a built-in Zsh command that provides a more reliable and consistent way to find the location of a command.

Disable by adding `ZC1005` to `disabled_katas` in `.zshellcheckrc`.

---

<a id="zc1006"></a>
### ZC1006 — Prefer [[ over test for tests

**Severity:** `style`  
**Auto-fix:** `no`

The `test` command is an external command and may not be available on all systems. The `[[...]]` construct is a Zsh keyword, offering safer and more powerful conditional expressions than the traditional `test` command. It prevents word splitting and pathname expansion, and supports advanced features like regex matching.

Disable by adding `ZC1006` to `disabled_katas` in `.zshellcheckrc`.

---

<a id="zc1007"></a>
### ZC1007 — Avoid using `chmod 777`

**Severity:** `warning`  
**Auto-fix:** `no`

Using `chmod 777` is a security risk as it gives read, write, and execute permissions to everyone. It's better to use more restrictive permissions.

Disable by adding `ZC1007` to `disabled_katas` in `.zshellcheckrc`.

---

<a id="zc1008"></a>
### ZC1008 — Use `\$(())` for arithmetic operations

**Severity:** `style`  
**Auto-fix:** `no`

The `let` command is a shell builtin, but the `\$(())` syntax is more portable and generally preferred for arithmetic operations in Zsh. It's also more powerful as it can be used in more contexts.

Disable by adding `ZC1008` to `disabled_katas` in `.zshellcheckrc`.

---

<a id="zc1009"></a>
### ZC1009 — Use `((...))` for C-style arithmetic

**Severity:** `style`  
**Auto-fix:** `no`

The `((...))` construct in Zsh allows for C-style arithmetic. It is generally more efficient and readable than using `expr` or other external commands for arithmetic.

Disable by adding `ZC1009` to `disabled_katas` in `.zshellcheckrc`.

---

<a id="zc1010"></a>
### ZC1010 — Use [[ ... ]] instead of [ ... ]

**Severity:** `style`  
**Auto-fix:** `yes`

Zsh's [[ ... ]] is more powerful and safer than [ ... ]. It supports pattern matching, regex, and doesn't require quoting variables to prevent word splitting.

Disable by adding `ZC1010` to `disabled_katas` in `.zshellcheckrc`.

---

<a id="zc1011"></a>
### ZC1011 — Use `git` porcelain commands instead of plumbing commands

**Severity:** `info`  
**Auto-fix:** `no`

Plumbing commands in `git` are designed for scripting and can be unstable. Porcelain commands are designed for interactive use and are more stable.

Disable by adding `ZC1011` to `disabled_katas` in `.zshellcheckrc`.

---

<a id="zc1012"></a>
### ZC1012 — Use `read -r` to prevent backslash escaping

**Severity:** `style`  
**Auto-fix:** `yes`

By default, `read` interprets backslashes as escape characters. Use `read -r` to treat backslashes literally, which is usually what you want.

Disable by adding `ZC1012` to `disabled_katas` in `.zshellcheckrc`.

---

<a id="zc1013"></a>
### ZC1013 — Use `((...))` for arithmetic operations instead of `let`

**Severity:** `info`  
**Auto-fix:** `yes`

The `let` command is a shell builtin, but the `((...))` syntax is more portable and generally preferred for arithmetic operations in Zsh.

Disable by adding `ZC1013` to `disabled_katas` in `.zshellcheckrc`.

---

<a id="zc1014"></a>
### ZC1014 — Use `git switch` or `git restore` instead of `git checkout`

**Severity:** `info`  
**Auto-fix:** `no`

The `git checkout` command can be ambiguous. `git switch` is used for switching branches and `git restore` is used for restoring files. Using these more specific commands can make your scripts clearer and less error-prone.

Disable by adding `ZC1014` to `disabled_katas` in `.zshellcheckrc`.

---

<a id="zc1015"></a>
### ZC1015 — Use `$(...)` for command substitution instead of backticks

**Severity:** `style`  
**Auto-fix:** `yes`

The `$(...)` syntax is the modern, recommended way to perform command substitution. It is more readable and can be nested easily, unlike backticks.

Disable by adding `ZC1015` to `disabled_katas` in `.zshellcheckrc`.

---

<a id="zc1016"></a>
### ZC1016 — Use `read -s` when reading sensitive information

**Severity:** `style`  
**Auto-fix:** `yes`

When asking for passwords or secrets, use `read -s` to prevent the input from being echoed to the terminal.

Disable by adding `ZC1016` to `disabled_katas` in `.zshellcheckrc`.

---

<a id="zc1017"></a>
### ZC1017 — Use `print -r` to print strings literally

**Severity:** `style`  
**Auto-fix:** `yes`

The `print` command interprets backslash escape sequences by default. To print a string literally, use the `-r` option.

Disable by adding `ZC1017` to `disabled_katas` in `.zshellcheckrc`.

---

<a id="zc1018"></a>
### ZC1018 — Superseded by ZC1009 — retired duplicate

**Severity:** `style`  
**Auto-fix:** `no`

Retained as a no-op stub so legacy `.zshellcheckrc` files that disable this ID keep parsing. See https://github.com/afadesigns/zshellcheck/issues/343 for context; the canonical detection lives in ZC1009.

Disable by adding `ZC1018` to `disabled_katas` in `.zshellcheckrc`.

---

<a id="zc1019"></a>
### ZC1019 — Superseded by ZC1005 — retired duplicate

**Severity:** `style`  
**Auto-fix:** `no`

Retained as a no-op stub so legacy `.zshellcheckrc` files that disable this ID keep parsing. See https://github.com/afadesigns/zshellcheck/issues/342 for context; the canonical detection lives in ZC1005.

Disable by adding `ZC1019` to `disabled_katas` in `.zshellcheckrc`.

---

<a id="zc1020"></a>
### ZC1020 — Use `[[ ... ]]` for tests instead of `test`

**Severity:** `style`  
**Auto-fix:** `no`

The `test` command is an external command and may not be available on all systems. The `[[...]]` construct is a Zsh keyword, offering safer and more powerful conditional expressions than the traditional `test` command.

Disable by adding `ZC1020` to `disabled_katas` in `.zshellcheckrc`.

---

<a id="zc1021"></a>
### ZC1021 — Use symbolic permissions with `chmod` instead of octal

**Severity:** `style`  
**Auto-fix:** `no`

Symbolic permissions (e.g., `u+x`) are more readable and less error-prone than octal permissions (e.g., `755`).

Disable by adding `ZC1021` to `disabled_katas` in `.zshellcheckrc`.

---

<a id="zc1022"></a>
### ZC1022 — Use `$((...))` for arithmetic expansion

**Severity:** `style`  
**Auto-fix:** `yes`

The `$((...))` syntax is the modern, recommended way to perform arithmetic expansion. It is more readable and can be nested easily, unlike `let`.

Disable by adding `ZC1022` to `disabled_katas` in `.zshellcheckrc`.

---

<a id="zc1023"></a>
### ZC1023 — Superseded by ZC1022 — retired duplicate `let` detector

**Severity:** `style`  
**Auto-fix:** `no`

Retained as a no-op stub so legacy `.zshellcheckrc` files that disable this ID keep parsing. The canonical `let` → `$((...))` guidance lives in ZC1022; see https://github.com/afadesigns/zshellcheck/issues/345.

Disable by adding `ZC1023` to `disabled_katas` in `.zshellcheckrc`.

---

<a id="zc1024"></a>
### ZC1024 — Superseded by ZC1022 — retired duplicate `let` detector

**Severity:** `style`  
**Auto-fix:** `no`

Retained as a no-op stub so legacy `.zshellcheckrc` files that disable this ID keep parsing. The canonical `let` → `$((...))` guidance lives in ZC1022; see https://github.com/afadesigns/zshellcheck/issues/345.

Disable by adding `ZC1024` to `disabled_katas` in `.zshellcheckrc`.

---

<a id="zc1025"></a>
### ZC1025 — Superseded by ZC1022 — retired duplicate `let` detector

**Severity:** `style`  
**Auto-fix:** `no`

Retained as a no-op stub so legacy `.zshellcheckrc` files that disable this ID keep parsing. The canonical `let` → `$((...))` guidance lives in ZC1022; see https://github.com/afadesigns/zshellcheck/issues/345.

Disable by adding `ZC1025` to `disabled_katas` in `.zshellcheckrc`.

---

<a id="zc1026"></a>
### ZC1026 — Superseded by ZC1022 — retired duplicate `let` detector

**Severity:** `style`  
**Auto-fix:** `no`

Retained as a no-op stub so legacy `.zshellcheckrc` files that disable this ID keep parsing. The canonical `let` → `$((...))` guidance lives in ZC1022; see https://github.com/afadesigns/zshellcheck/issues/345.

Disable by adding `ZC1026` to `disabled_katas` in `.zshellcheckrc`.

---

<a id="zc1027"></a>
### ZC1027 — Superseded by ZC1022 — retired duplicate `let` detector

**Severity:** `style`  
**Auto-fix:** `no`

Retained as a no-op stub so legacy `.zshellcheckrc` files that disable this ID keep parsing. The canonical `let` → `$((...))` guidance lives in ZC1022; see https://github.com/afadesigns/zshellcheck/issues/345.

Disable by adding `ZC1027` to `disabled_katas` in `.zshellcheckrc`.

---

<a id="zc1028"></a>
### ZC1028 — Superseded by ZC1022 — retired duplicate `let` detector

**Severity:** `style`  
**Auto-fix:** `no`

Retained as a no-op stub so legacy `.zshellcheckrc` files that disable this ID keep parsing. The canonical `let` → `$((...))` guidance lives in ZC1022; see https://github.com/afadesigns/zshellcheck/issues/345.

Disable by adding `ZC1028` to `disabled_katas` in `.zshellcheckrc`.

---

<a id="zc1029"></a>
### ZC1029 — Superseded by ZC1022 — retired duplicate `let` detector

**Severity:** `style`  
**Auto-fix:** `no`

Retained as a no-op stub so legacy `.zshellcheckrc` files that disable this ID keep parsing. The canonical `let` → `$((...))` guidance lives in ZC1022; see https://github.com/afadesigns/zshellcheck/issues/345.

Disable by adding `ZC1029` to `disabled_katas` in `.zshellcheckrc`.

---

<a id="zc1030"></a>
### ZC1030 — Use `printf` instead of `echo`

**Severity:** `style`  
**Auto-fix:** `no`

The `echo` command's behavior can be inconsistent across different shells and environments, especially with flags and escape sequences. `printf` provides more reliable and portable string formatting.

Disable by adding `ZC1030` to `disabled_katas` in `.zshellcheckrc`.

---

<a id="zc1031"></a>
### ZC1031 — Use `#!/usr/bin/env zsh` for portability

**Severity:** `info`  
**Auto-fix:** `yes`

Using `#!/usr/bin/env zsh` is more portable than `#!/bin/zsh` because it searches for the `zsh` executable in the user's `PATH`.

Disable by adding `ZC1031` to `disabled_katas` in `.zshellcheckrc`.

---

<a id="zc1032"></a>
### ZC1032 — Use `((...))` for C-style incrementing

**Severity:** `style`  
**Auto-fix:** `yes`

Instead of `let i=i+1` or `let i=i-1`, you can use the more concise and idiomatic C-style increment `(( i++ ))` / decrement `(( i-- ))` in Zsh.

Disable by adding `ZC1032` to `disabled_katas` in `.zshellcheckrc`.

---

<a id="zc1033"></a>
### ZC1033 — Superseded by ZC1022 — retired duplicate `let` detector

**Severity:** `style`  
**Auto-fix:** `no`

Retained as a no-op stub so legacy `.zshellcheckrc` files that disable this ID keep parsing. The canonical `let` → `$((...))` guidance lives in ZC1022; see https://github.com/afadesigns/zshellcheck/issues/345.

Disable by adding `ZC1033` to `disabled_katas` in `.zshellcheckrc`.

---

<a id="zc1034"></a>
### ZC1034 — Use `command -v` instead of `which`

**Severity:** `style`  
**Auto-fix:** `yes`

`which` is an external command and may not be available or consistent across all systems. `command -v` is a POSIX standard and a shell builtin, making it more portable and reliable for checking if a command exists.

Disable by adding `ZC1034` to `disabled_katas` in `.zshellcheckrc`.

---

<a id="zc1035"></a>
### ZC1035 — Superseded by ZC1022 — retired duplicate `let` detector

**Severity:** `style`  
**Auto-fix:** `no`

Retained as a no-op stub so legacy `.zshellcheckrc` files that disable this ID keep parsing. The canonical `let` → `$((...))` guidance lives in ZC1022; see https://github.com/afadesigns/zshellcheck/issues/345.

Disable by adding `ZC1035` to `disabled_katas` in `.zshellcheckrc`.

---

<a id="zc1036"></a>
### ZC1036 — Prefer `[[ ... ]]` over `test` command

**Severity:** `style`  
**Auto-fix:** `no`

The `[[ ... ]]` construct is a more powerful and safer alternative to the `test` command (or `[ ... ]`) for conditional expressions in modern shells. It handles word splitting and globbing more intuitively and supports advanced features like regex matching.

Disable by adding `ZC1036` to `disabled_katas` in `.zshellcheckrc`.

---

<a id="zc1037"></a>
### ZC1037 — Use 'print -r --' for variable expansion

**Severity:** `style`  
**Auto-fix:** `no`

Using 'echo' to print strings containing variables can lead to unexpected behavior if the variable contains special characters or flags. A safer, more reliable alternative is 'print -r --'.

Disable by adding `ZC1037` to `disabled_katas` in `.zshellcheckrc`.

---

<a id="zc1038"></a>
### ZC1038 — Avoid useless use of cat

**Severity:** `style`  
**Auto-fix:** `no`

Using `cat file | command` is unnecessary and inefficient. Most commands can read from a file directly, e.g., `command file`. If not, you can use input redirection: `command < file`.

Disable by adding `ZC1038` to `disabled_katas` in `.zshellcheckrc`.

---

<a id="zc1039"></a>
### ZC1039 — Avoid `rm` with root path

**Severity:** `warning`  
**Auto-fix:** `no`

Running `rm` on the root directory `/` is dangerous. Ensure you are not deleting the entire filesystem.

Disable by adding `ZC1039` to `disabled_katas` in `.zshellcheckrc`.

---

<a id="zc1040"></a>
### ZC1040 — Use (N) nullglob qualifier for globs in loops

**Severity:** `style`  
**Auto-fix:** `yes`

In Zsh, a glob that matches nothing (e.g., `*.txt`) will cause an error by default. Use the `(N)` glob qualifier to make it null (empty) if no matches found, preventing the error.

Disable by adding `ZC1040` to `disabled_katas` in `.zshellcheckrc`.

---

<a id="zc1041"></a>
### ZC1041 — Do not use variables in printf format string

**Severity:** `style`  
**Auto-fix:** `no`

Using variables in `printf` format strings allows for format string attacks and unexpected behavior if the variable contains `%`. Use `printf '%s' "$var"` instead.

Disable by adding `ZC1041` to `disabled_katas` in `.zshellcheckrc`.

---

<a id="zc1042"></a>
### ZC1042 — Use "$@" to iterate over arguments

**Severity:** `style`  
**Auto-fix:** `no`

`$*` joins all arguments into a single string, which is rarely what you want in a loop. Use `"$@"` to iterate over each argument individually.

Disable by adding `ZC1042` to `disabled_katas` in `.zshellcheckrc`.

---

<a id="zc1043"></a>
### ZC1043 — Use `local` for variables in functions

**Severity:** `style`  
**Auto-fix:** `yes`

Variables defined in functions are global by default in Zsh. Use `local` to scope them to the function.

Disable by adding `ZC1043` to `disabled_katas` in `.zshellcheckrc`.

---

<a id="zc1044"></a>
### ZC1044 — Check for unchecked `cd` commands

**Severity:** `warning`  
**Auto-fix:** `no`

`cd` failures should be handled to avoid executing commands in the wrong directory. Use `cd ... || return` (or `exit`).

Disable by adding `ZC1044` to `disabled_katas` in `.zshellcheckrc`.

---

<a id="zc1045"></a>
### ZC1045 — Declare and assign separately to avoid masking return values

**Severity:** `info`  
**Auto-fix:** `no`

Declaring a variable with `local var=$(cmd)` masks the return value of `cmd`. The `local` command returns 0 (success) even if `cmd` fails. Declare the variable first (`local var`), then assign it (`var=$(cmd)`).

Disable by adding `ZC1045` to `disabled_katas` in `.zshellcheckrc`.

---

<a id="zc1046"></a>
### ZC1046 — Avoid `eval`

**Severity:** `warning`  
**Auto-fix:** `no`

`eval` is dangerous as it executes arbitrary code. Use arrays, parameter expansion, or other constructs instead.

Disable by adding `ZC1046` to `disabled_katas` in `.zshellcheckrc`.

---

<a id="zc1047"></a>
### ZC1047 — Avoid `sudo` in scripts

**Severity:** `warning`  
**Auto-fix:** `no`

Using `sudo` in scripts is generally discouraged. It makes the script interactive and less portable. Run the script as root or use `sudo` to invoke the script.

Disable by adding `ZC1047` to `disabled_katas` in `.zshellcheckrc`.

---

<a id="zc1048"></a>
### ZC1048 — Avoid `source` with relative paths

**Severity:** `style`  
**Auto-fix:** `no`

Sourcing a file with a relative path (e.g. `source ./lib.zsh`) depends on the current working directory. Use `${0:a:h}/lib.zsh` to source relative to the script location.

Disable by adding `ZC1048` to `disabled_katas` in `.zshellcheckrc`.

---

<a id="zc1049"></a>
### ZC1049 — Prefer functions over aliases

**Severity:** `style`  
**Auto-fix:** `no`

Aliases are expanded at parse time and can be confusing in scripts. Use functions for more predictable behavior.

Disable by adding `ZC1049` to `disabled_katas` in `.zshellcheckrc`.

---

<a id="zc1050"></a>
### ZC1050 — Avoid iterating over `ls` output

**Severity:** `style`  
**Auto-fix:** `no`

Iterating over `ls` output is fragile because filenames can contain spaces and newlines. Use globs (e.g. `for f in *.txt`) instead.

Disable by adding `ZC1050` to `disabled_katas` in `.zshellcheckrc`.

---

<a id="zc1051"></a>
### ZC1051 — Quote variables in `rm` to avoid globbing

**Severity:** `warning`  
**Auto-fix:** `yes`

`rm $VAR` is dangerous if `$VAR` contains spaces or glob characters. Quote the variable (`rm "$VAR"`) to ensure safe deletion.

Disable by adding `ZC1051` to `disabled_katas` in `.zshellcheckrc`.

---

<a id="zc1052"></a>
### ZC1052 — Avoid `sed -i` for portability

**Severity:** `style`  
**Auto-fix:** `no`

`sed -i` usage varies between GNU/Linux and macOS/BSD. macOS requires an extension argument (e.g. `sed -i ''`), while GNU does not. Use a temporary file and `mv`, or `perl -i`, for portability.

Disable by adding `ZC1052` to `disabled_katas` in `.zshellcheckrc`.

---

<a id="zc1053"></a>
### ZC1053 — Silence `grep` output in conditions

**Severity:** `style`  
**Auto-fix:** `yes`

Using `grep` in a condition prints matches to stdout. Use `grep -q` (or `> /dev/null`) to silence output if you only care about the exit code.

Disable by adding `ZC1053` to `disabled_katas` in `.zshellcheckrc`.

---

<a id="zc1054"></a>
### ZC1054 — Use POSIX classes in regex/glob

**Severity:** `style`  
**Auto-fix:** `no`

Ranges like `[a-z]` are locale-dependent. Use `[[:lower:]]` or `[a-z]` with `LC_ALL=C` to be explicit.

Disable by adding `ZC1054` to `disabled_katas` in `.zshellcheckrc`.

---

<a id="zc1055"></a>
### ZC1055 — Use `[[ -n/-z ]]` for empty string checks

**Severity:** `style`  
**Auto-fix:** `yes`

Comparing with empty string is less idiomatic than using `[[ -z $var ]]` (is empty) or `[[ -n $var ]]` (is not empty).

Disable by adding `ZC1055` to `disabled_katas` in `.zshellcheckrc`.

---

<a id="zc1056"></a>
### ZC1056 — Avoid `$((...))` as a statement

**Severity:** `style`  
**Auto-fix:** `no`

Using `$((...))` as a statement tries to execute the result as a command. Use `((...))` for arithmetic evaluation/assignment.

Disable by adding `ZC1056` to `disabled_katas` in `.zshellcheckrc`.

---

<a id="zc1057"></a>
### ZC1057 — Avoid `ls` in assignments

**Severity:** `style`  
**Auto-fix:** `no`

Assigning the output of `ls` to a variable is fragile. Use globs or arrays (e.g. `files=(*)`) to handle filenames correctly.

Disable by adding `ZC1057` to `disabled_katas` in `.zshellcheckrc`.

---

<a id="zc1058"></a>
### ZC1058 — Avoid `sudo` with redirection

**Severity:** `style`  
**Auto-fix:** `no`

Redirecting output of `sudo` (e.g. `sudo cmd > /file`) fails if the current user doesn't have permission. Use `| sudo tee /file` instead.

Disable by adding `ZC1058` to `disabled_katas` in `.zshellcheckrc`.

---

<a id="zc1059"></a>
### ZC1059 — Use `${var:?}` for `rm` arguments

**Severity:** `warning`  
**Auto-fix:** `no`

Deleting a directory based on a variable is dangerous if the variable is empty or unset. Use `${var:?}` to fail if empty, or check explicitly.

Disable by adding `ZC1059` to `disabled_katas` in `.zshellcheckrc`.

---

<a id="zc1060"></a>
### ZC1060 — Avoid `ps | grep` without exclusion

**Severity:** `style`  
**Auto-fix:** `no`

`ps | grep pattern` often matches the grep process itself. Use `grep [p]attern`, `pgrep`, or exclude grep with `grep -v grep`.

Disable by adding `ZC1060` to `disabled_katas` in `.zshellcheckrc`.

---

<a id="zc1061"></a>
### ZC1061 — Prefer `{start..end}` over `seq`

**Severity:** `style`  
**Auto-fix:** `yes`

Using `seq` creates an external process. Zsh supports integer range expansion natively: `{1..10}`.

Disable by adding `ZC1061` to `disabled_katas` in `.zshellcheckrc`.

---

<a id="zc1062"></a>
### ZC1062 — Prefer `grep -E` over `egrep`

**Severity:** `info`  
**Auto-fix:** `yes`

`egrep` is deprecated. Use `grep -E` instead.

Disable by adding `ZC1062` to `disabled_katas` in `.zshellcheckrc`.

---

<a id="zc1063"></a>
### ZC1063 — Prefer `grep -F` over `fgrep`

**Severity:** `info`  
**Auto-fix:** `yes`

`fgrep` is deprecated. Use `grep -F` instead.

Disable by adding `ZC1063` to `disabled_katas` in `.zshellcheckrc`.

---

<a id="zc1064"></a>
### ZC1064 — Prefer `command -v` over `type`

**Severity:** `info`  
**Auto-fix:** `yes`

`type` output format varies and is not POSIX standard for checking existence. `command -v` is quieter and standard.

Disable by adding `ZC1064` to `disabled_katas` in `.zshellcheckrc`.

---

<a id="zc1065"></a>
### ZC1065 — Ensure spaces around `[` and `[[`

**Severity:** `error`  
**Auto-fix:** `no`

`[[condition]]` is parsed incorrectly. Add spaces: `[[ condition ]]`.

Disable by adding `ZC1065` to `disabled_katas` in `.zshellcheckrc`.

---

<a id="zc1066"></a>
### ZC1066 — Avoid iterating over `cat` output

**Severity:** `style`  
**Auto-fix:** `no`

Iterating over `cat` output is fragile because lines can contain spaces. Use `while IFS= read -r line; do ... done < file` or `($(<file))` array expansion.

Disable by adding `ZC1066` to `disabled_katas` in `.zshellcheckrc`.

---

<a id="zc1067"></a>
### ZC1067 — Separate `export` and assignment to avoid masking return codes

**Severity:** `style`  
**Auto-fix:** `no`

Running `export var=$(cmd)` masks the return code of `cmd`. The exit status will be that of `export` (usually 0). Declare the variable first or export it after assignment.

Disable by adding `ZC1067` to `disabled_katas` in `.zshellcheckrc`.

---

<a id="zc1068"></a>
### ZC1068 — Use `add-zsh-hook` instead of defining hook functions directly

**Severity:** `info`  
**Auto-fix:** `no`

Defining special functions like `precmd`, `preexec`, `chpwd`, etc. directly overwrites any previously defined hooks. Use `autoload -Uz add-zsh-hook; add-zsh-hook <hook> <function>` to append to the hook list safely.

Disable by adding `ZC1068` to `disabled_katas` in `.zshellcheckrc`.

---

<a id="zc1069"></a>
### ZC1069 — Avoid `local` outside of functions

**Severity:** `info`  
**Auto-fix:** `no`

The `local` builtin can only be used inside functions. Using it in the global scope causes an error.

Disable by adding `ZC1069` to `disabled_katas` in `.zshellcheckrc`.

---

<a id="zc1070"></a>
### ZC1070 — Use `builtin` or `command` to avoid infinite recursion in wrapper functions

**Severity:** `warning`  
**Auto-fix:** `no`

When defining a wrapper function with the same name as a builtin or command (e.g., `cd`), calling the command directly inside the function causes infinite recursion. Use `builtin cd` or `command cd`.

Disable by adding `ZC1070` to `disabled_katas` in `.zshellcheckrc`.

---

<a id="zc1071"></a>
### ZC1071 — Use `+=` for appending to arrays

**Severity:** `warning`  
**Auto-fix:** `no`

Appending to an array using `arr=($arr ...)` is verbose and slower. Use `arr+=(...)` instead.

Disable by adding `ZC1071` to `disabled_katas` in `.zshellcheckrc`.

---

<a id="zc1072"></a>
### ZC1072 — Use `awk` instead of `grep | awk`

**Severity:** `style`  
**Auto-fix:** `no`

`grep pattern | awk '{...}'` is inefficient. Use `awk '/pattern/ {...}'` to combine matching and processing in a single process.

Disable by adding `ZC1072` to `disabled_katas` in `.zshellcheckrc`.

---

<a id="zc1073"></a>
### ZC1073 — Unnecessary use of `$` in arithmetic expressions

**Severity:** `style`  
**Auto-fix:** `yes`

Variables in `((...))` do not need `$` prefix. Use `(( var > 0 ))` instead of `(( $var > 0 ))`.

Disable by adding `ZC1073` to `disabled_katas` in `.zshellcheckrc`.

---

<a id="zc1074"></a>
### ZC1074 — Prefer modifiers :h/:t over dirname/basename

**Severity:** `style`  
**Auto-fix:** `no`

Zsh provides modifiers like `:h` (head/dirname) and `:t` (tail/basename) that are faster and more idiomatic than spawning external commands.

Disable by adding `ZC1074` to `disabled_katas` in `.zshellcheckrc`.

---

<a id="zc1075"></a>
### ZC1075 — Quote variable expansions to prevent globbing

**Severity:** `warning`  
**Auto-fix:** `no`

Unquoted variable expansions in Zsh are subject to globbing (filename generation). If the variable contains characters like `*` or `?`, it might match files unexpectedly. Use quotes `"$var"` to prevent this.

Disable by adding `ZC1075` to `disabled_katas` in `.zshellcheckrc`.

---

<a id="zc1076"></a>
### ZC1076 — Use `autoload -Uz` for lazy loading

**Severity:** `style`  
**Auto-fix:** `yes`

When using `autoload`, prefer `-Uz` to ensure standard Zsh behavior (no alias expansion, zsh style). `-U` prevents alias expansion, and `-z` ensures Zsh style autoloading.

Disable by adding `ZC1076` to `disabled_katas` in `.zshellcheckrc`.

---

<a id="zc1077"></a>
### ZC1077 — Prefer `${var:u/l}` over `tr` for case conversion

**Severity:** `style`  
**Auto-fix:** `no`

Using `tr` in a pipeline for simple case conversion is slower than using Zsh's built-in parameter expansion flags `:u` (upper) and `:l` (lower).

Disable by adding `ZC1077` to `disabled_katas` in `.zshellcheckrc`.

---

<a id="zc1078"></a>
### ZC1078 — Quote `$@` and `$*` when passing arguments

**Severity:** `warning`  
**Auto-fix:** `yes`

Using unquoted `$@` or `$*` splits arguments by IFS (usually space). Use `"$@"` to preserve the original argument grouping, or `"$*"` to join them into a single string.

Disable by adding `ZC1078` to `disabled_katas` in `.zshellcheckrc`.

---

<a id="zc1079"></a>
### ZC1079 — Quote RHS of `==` in `[[ ... ]]` to prevent pattern matching

**Severity:** `warning`  
**Auto-fix:** `yes`

In `[[ ... ]]`, unquoted variable expansions on the right-hand side of `==` or `!=` are treated as patterns (globbing). If you intend to compare strings literally, quote the variable.

Disable by adding `ZC1079` to `disabled_katas` in `.zshellcheckrc`.

---

<a id="zc1080"></a>
### ZC1080 — Use `(N)` nullglob qualifier for globs in loops

**Severity:** `style`  
**Auto-fix:** `no`

In Zsh, if a glob matches no files, it throws an error by default. When iterating over a glob in a `for` loop, use the `(N)` glob qualifier to allow it to match nothing (nullglob).

Disable by adding `ZC1080` to `disabled_katas` in `.zshellcheckrc`.

---

<a id="zc1081"></a>
### ZC1081 — Use `${#var}` to get string length instead of `wc -c`

**Severity:** `style`  
**Auto-fix:** `no`

Using `echo $var | wc -c` involves a subshell and external command overhead. Zsh has a built-in operator `${#var}` to get the length of a string instantly.

Disable by adding `ZC1081` to `disabled_katas` in `.zshellcheckrc`.

---

<a id="zc1082"></a>
### ZC1082 — Prefer `${var//old/new}` over `sed` for simple replacements

**Severity:** `style`  
**Auto-fix:** `no`

Using `sed` for simple string replacement is slower than Zsh's built-in parameter expansion. Use `${var/old/new}` (replace first) or `${var//old/new}` (replace all).

Disable by adding `ZC1082` to `disabled_katas` in `.zshellcheckrc`.

---

<a id="zc1083"></a>
### ZC1083 — Brace expansion limits cannot be variables

**Severity:** `error`  
**Auto-fix:** `no`

Brace expansion `{x..y}` happens before variable expansion. `{1..$n}` will not work. Use `seq` or `for ((...))`.

Disable by adding `ZC1083` to `disabled_katas` in `.zshellcheckrc`.

---

<a id="zc1084"></a>
### ZC1084 — Quote globs in `find` commands

**Severity:** `warning`  
**Auto-fix:** `yes`

Unquoted globs in `find` commands are expanded by the shell before `find` runs. If files match, `find` receives the list of files instead of the pattern. Quote arguments to `-name`, `-path`, etc.

Disable by adding `ZC1084` to `disabled_katas` in `.zshellcheckrc`.

---

<a id="zc1085"></a>
### ZC1085 — Quote variable expansions in `for` loops

**Severity:** `warning`  
**Auto-fix:** `yes`

Unquoted variable expansions in `for` loops are split by IFS (usually spaces). This often leads to iterating over words instead of lines or array elements. Quote the expansion to preserve structure.

Disable by adding `ZC1085` to `disabled_katas` in `.zshellcheckrc`.

---

<a id="zc1086"></a>
### ZC1086 — Prefer `func() { ... }` over `function func { ... }`

**Severity:** `style`  
**Auto-fix:** `yes`

The `function` keyword is optional in Zsh and non-standard in POSIX sh. Using `func() { ... }` is more portable and consistent.

Disable by adding `ZC1086` to `disabled_katas` in `.zshellcheckrc`.

---

<a id="zc1087"></a>
### ZC1087 — Output redirection overwrites input file

**Severity:** `error`  
**Auto-fix:** `no`

Redirecting output to a file that is also being read as input causes the file to be truncated before it is read. Use a temporary file or `sponge`.

Disable by adding `ZC1087` to `disabled_katas` in `.zshellcheckrc`.

---

<a id="zc1088"></a>
### ZC1088 — Subshell isolates state changes

**Severity:** `warning`  
**Auto-fix:** `no`

Commands inside `( ... )` run in a subshell. State changes like `cd`, `export`, or variable assignments are lost when the subshell exits. Use `{ ... }` for grouping if you want to preserve state changes.

Disable by adding `ZC1088` to `disabled_katas` in `.zshellcheckrc`.

---

<a id="zc1089"></a>
### ZC1089 — Redirection order matters (`2>&1 > file`)

**Severity:** `error`  
**Auto-fix:** `no`

Redirecting stderr to stdout (`2>&1`) before redirecting stdout to a file (`> file`) means stderr goes to the *original* stdout (usually tty), not the file. Use `> file 2>&1` or `&> file` to redirect both.

Disable by adding `ZC1089` to `disabled_katas` in `.zshellcheckrc`.

---

<a id="zc1090"></a>
### ZC1090 — Quoted regex pattern in `=~`

**Severity:** `warning`  
**Auto-fix:** `no`

Quoting the pattern on the right side of `=~` forces literal string matching in Zsh/Bash. Regex metacharacters inside quotes will be matched literally. Remove quotes to enable regex matching, or use `==` for literal string comparison.

Disable by adding `ZC1090` to `disabled_katas` in `.zshellcheckrc`.

---

<a id="zc1091"></a>
### ZC1091 — Use `((...))` for arithmetic comparisons in `[[...]]`

**Severity:** `style`  
**Auto-fix:** `yes`

The `[[ ... ]]` construct is primarily for string comparisons and file tests. For arithmetic comparisons (`-eq`, `-lt`, etc.), use the dedicated arithmetic context `(( ... ))`. It is cleaner and strictly numeric.

Disable by adding `ZC1091` to `disabled_katas` in `.zshellcheckrc`.

---

<a id="zc1092"></a>
### ZC1092 — Prefer `print` or `printf` over `echo` in Zsh

**Severity:** `warning`  
**Auto-fix:** `yes`

In Zsh, `echo` behavior can vary significantly based on options like `BSD_ECHO`. `print` is a builtin with consistent behavior and more features. For formatted output, `printf` is preferred.

Disable by adding `ZC1092` to `disabled_katas` in `.zshellcheckrc`.

---

<a id="zc1093"></a>
### ZC1093 — Superseded by ZC1038 — retired duplicate

**Severity:** `style`  
**Auto-fix:** `no`

Retained as a no-op stub so legacy `.zshellcheckrc` files that disable this ID keep parsing. See https://github.com/afadesigns/zshellcheck/issues/341 for context; the canonical detection lives in ZC1038.

Disable by adding `ZC1093` to `disabled_katas` in `.zshellcheckrc`.

---

<a id="zc1094"></a>
### ZC1094 — Use parameter expansion instead of `sed` for simple substitutions

**Severity:** `style`  
**Auto-fix:** `no`

For simple string substitutions on variables, use Zsh parameter expansion `${var//pattern/replacement}` instead of piping through `sed`. It avoids spawning an external process.

Disable by adding `ZC1094` to `disabled_katas` in `.zshellcheckrc`.

---

<a id="zc1095"></a>
### ZC1095 — Use `repeat N` for simple repetition

**Severity:** `style`  
**Auto-fix:** `yes`

Zsh provides `repeat N do ... done` for running a block a fixed number of times. It is cleaner than `for i in {1..N}` or C-style for loops when the iterator variable is unused.

Disable by adding `ZC1095` to `disabled_katas` in `.zshellcheckrc`.

---

<a id="zc1096"></a>
### ZC1096 — Warn on `bc` for simple arithmetic

**Severity:** `style`  
**Auto-fix:** `no`

Zsh has built-in support for floating point arithmetic using `(( ... ))` or `$(( ... ))`. Using `bc` is often unnecessary and slower.

Disable by adding `ZC1096` to `disabled_katas` in `.zshellcheckrc`.

---

<a id="zc1097"></a>
### ZC1097 — Declare loop variables as `local` in functions

**Severity:** `style`  
**Auto-fix:** `no`

Loop variables in `for` loops are global by default in Zsh functions. Use `local` to scope them to the function before the loop.

Disable by adding `ZC1097` to `disabled_katas` in `.zshellcheckrc`.

---

<a id="zc1098"></a>
### ZC1098 — Use `(q)` flag for quoting variables in eval

**Severity:** `style`  
**Auto-fix:** `no`

When constructing a command string for `eval`, use the `(q)` flag (or `(qq)`, `(q-)`) to safely quote variables and prevent command injection.

Disable by adding `ZC1098` to `disabled_katas` in `.zshellcheckrc`.

---

<a id="zc1099"></a>
### ZC1099 — Use `(f)` flag to split lines instead of `while read`

**Severity:** `style`  
**Auto-fix:** `no`

Zsh provides the `(f)` parameter expansion flag to split a string into lines. Iterating over `${(f)variable}` is often cleaner and faster than piping to `while read`.

Disable by adding `ZC1099` to `disabled_katas` in `.zshellcheckrc`.

---

<a id="zc1100"></a>
### ZC1100 — Use parameter expansion instead of `dirname`/`basename`

**Severity:** `style`  
**Auto-fix:** `no`

Zsh parameter expansion `${var%/*}` (dirname) and `${var##*/}` (basename) avoid spawning external processes for simple path manipulation.

Disable by adding `ZC1100` to `disabled_katas` in `.zshellcheckrc`.

---

<a id="zc1101"></a>
### ZC1101 — Use `$(( ))` instead of `bc` for simple arithmetic

**Severity:** `style`  
**Auto-fix:** `no`

Zsh supports arithmetic expansion with `$(( ))` and floating point via `zmodload zsh/mathfunc`. Avoid piping to `bc` for simple calculations.

Disable by adding `ZC1101` to `disabled_katas` in `.zshellcheckrc`.

---

<a id="zc1102"></a>
### ZC1102 — Redirecting output of `sudo` doesn't work as expected

**Severity:** `style`  
**Auto-fix:** `no`

Redirections are performed by the current shell before `sudo` is started. So `sudo echo > /root/file` will try to open `/root/file` as the current user, failing. Use `echo ... | sudo tee file` or `sudo sh -c 'echo ... > file'`.

Disable by adding `ZC1102` to `disabled_katas` in `.zshellcheckrc`.

---

<a id="zc1103"></a>
### ZC1103 — Suggest `path` array instead of `$PATH` string manipulation (direct assignment)

**Severity:** `style`  
**Auto-fix:** `no`

Zsh automatically maps the `$PATH` environment variable to the `$path` array. Modifying `$path` is cleaner and less error-prone than manipulating the colon-separated `$PATH` string.

Disable by adding `ZC1103` to `disabled_katas` in `.zshellcheckrc`.

---

<a id="zc1104"></a>
### ZC1104 — Suggest `path` array instead of `export PATH` string manipulation

**Severity:** `style`  
**Auto-fix:** `no`

Zsh automatically maps the `$PATH` environment variable to the `$path` array. Modifying `$path` is cleaner and less error-prone than manipulating the colon-separated `$PATH` string.

Disable by adding `ZC1104` to `disabled_katas` in `.zshellcheckrc`.

---

<a id="zc1105"></a>
### ZC1105 — Avoid nested arithmetic expansions for clarity

**Severity:** `style`  
**Auto-fix:** `no`

While Zsh supports nested arithmetic expansions like `(( $((...)) ))`, they can make code harder to read and reason about. Prefer flatter expressions or temporary variables for intermediate results to improve clarity.

Disable by adding `ZC1105` to `disabled_katas` in `.zshellcheckrc`.

---

<a id="zc1106"></a>
### ZC1106 — Avoid `set -x` in production scripts for sensitive data exposure

**Severity:** `style`  
**Auto-fix:** `no`

Using `set -x` (xtrace) in production environments can expose sensitive information, such as API keys or passwords, in logs. While useful for debugging, it should be avoided in production. Consider using targeted debugging or secure logging.

Disable by adding `ZC1106` to `disabled_katas` in `.zshellcheckrc`.

---

<a id="zc1107"></a>
### ZC1107 — Use (( ... )) for arithmetic conditions

**Severity:** `style`  
**Auto-fix:** `no`

Use `(( ... ))` for arithmetic comparisons instead of `[ ... -eq ... ]`. The double parenthesis syntax supports standard math operators (`>`, `<`, `==`, `!=`) and is optimized.

Disable by adding `ZC1107` to `disabled_katas` in `.zshellcheckrc`.

---

<a id="zc1108"></a>
### ZC1108 — Use Zsh case conversion instead of `tr`

**Severity:** `style`  
**Auto-fix:** `no`

Zsh provides `${(U)var}` for uppercase and `${(L)var}` for lowercase. Avoid piping through `tr '[:lower:]' '[:upper:]'` for simple case conversion.

Disable by adding `ZC1108` to `disabled_katas` in `.zshellcheckrc`.

---

<a id="zc1109"></a>
### ZC1109 — Use parameter expansion instead of `cut` for field extraction

**Severity:** `style`  
**Auto-fix:** `no`

For simple field extraction from variables, use Zsh parameter expansion like `${var%%:*}` or `${(s.:.)var}` instead of piping through `cut`.

Disable by adding `ZC1109` to `disabled_katas` in `.zshellcheckrc`.

---

<a id="zc1110"></a>
### ZC1110 — Use Zsh subscripts instead of `head -1` or `tail -1`

**Severity:** `style`  
**Auto-fix:** `no`

Zsh array subscripts `${lines[1]}` and `${lines[-1]}` can extract the first or last element without spawning `head` or `tail` as external processes.

Disable by adding `ZC1110` to `disabled_katas` in `.zshellcheckrc`.

---

<a id="zc1111"></a>
### ZC1111 — Avoid `xargs` for simple command invocation

**Severity:** `style`  
**Auto-fix:** `no`

Zsh can iterate arrays directly with `for` loops or use `${(f)...}` to split command output by newlines. Avoid `xargs` when processing lines one at a time.

Disable by adding `ZC1111` to `disabled_katas` in `.zshellcheckrc`.

---

<a id="zc1112"></a>
### ZC1112 — Avoid `grep -c` — use Zsh pattern matching for counting

**Severity:** `style`  
**Auto-fix:** `no`

For counting matches in a variable, use Zsh `${#${(f)...}}` or array filtering with `${(M)array:#pattern}` instead of piping through `grep -c`.

Disable by adding `ZC1112` to `disabled_katas` in `.zshellcheckrc`.

---

<a id="zc1113"></a>
### ZC1113 — Use `${var:A}` instead of `realpath` or `readlink -f`

**Severity:** `style`  
**Auto-fix:** `no`

Zsh provides the `:A` modifier to resolve a path to its absolute form, following symlinks. Avoid spawning `realpath` or `readlink -f` as external processes.

Disable by adding `ZC1113` to `disabled_katas` in `.zshellcheckrc`.

---

<a id="zc1114"></a>
### ZC1114 — Consider Zsh `=(...)` for temporary files

**Severity:** `style`  
**Auto-fix:** `no`

Zsh `=(cmd)` creates a temporary file with the command output that is automatically cleaned up. Consider this instead of manual `mktemp` and cleanup patterns.

Disable by adding `ZC1114` to `disabled_katas` in `.zshellcheckrc`.

---

<a id="zc1115"></a>
### ZC1115 — Use Zsh string manipulation instead of `rev`

**Severity:** `style`  
**Auto-fix:** `no`

Zsh can reverse strings using parameter expansion. Avoid spawning `rev` as an external process for simple string reversal.

Disable by adding `ZC1115` to `disabled_katas` in `.zshellcheckrc`.

---

<a id="zc1116"></a>
### ZC1116 — Use Zsh multios instead of `tee`

**Severity:** `style`  
**Auto-fix:** `no`

Zsh `setopt multios` allows redirecting output to multiple files with `cmd > file1 > file2`. Avoid spawning `tee` for simple output duplication.

Disable by adding `ZC1116` to `disabled_katas` in `.zshellcheckrc`.

---

<a id="zc1117"></a>
### ZC1117 — Use `&!` or `disown` instead of `nohup`

**Severity:** `style`  
**Auto-fix:** `no`

Zsh provides `&!` (shorthand for `& disown`) to run a command in the background immune to hangups. Avoid spawning `nohup` as an external process.

Disable by adding `ZC1117` to `disabled_katas` in `.zshellcheckrc`.

---

<a id="zc1118"></a>
### ZC1118 — Use `print -rn` instead of `echo -n`

**Severity:** `style`  
**Auto-fix:** `yes`

The behavior of `echo -n` varies across shells and platforms. In Zsh, `print -rn` is the reliable way to output text without a trailing newline.

Disable by adding `ZC1118` to `disabled_katas` in `.zshellcheckrc`.

---

<a id="zc1119"></a>
### ZC1119 — Use `$EPOCHSECONDS` instead of `date +%s`

**Severity:** `style`  
**Auto-fix:** `no`

Zsh provides `$EPOCHSECONDS` and `$EPOCHREALTIME` via `zsh/datetime` module. Avoid spawning `date` for simple Unix timestamp retrieval.

Disable by adding `ZC1119` to `disabled_katas` in `.zshellcheckrc`.

---

<a id="zc1120"></a>
### ZC1120 — Use `$PWD` instead of `pwd`

**Severity:** `style`  
**Auto-fix:** `no`

Zsh maintains `$PWD` as a built-in variable tracking the current directory. Avoid spawning `pwd` as an external process.

Disable by adding `ZC1120` to `disabled_katas` in `.zshellcheckrc`.

---

<a id="zc1121"></a>
### ZC1121 — Use `$HOST` instead of `hostname`

**Severity:** `style`  
**Auto-fix:** `no`

Zsh provides `$HOST` as a built-in variable containing the hostname. Avoid spawning `hostname` as an external process.

Disable by adding `ZC1121` to `disabled_katas` in `.zshellcheckrc`.

---

<a id="zc1122"></a>
### ZC1122 — Use `$USER` instead of `whoami`

**Severity:** `style`  
**Auto-fix:** `no`

Zsh provides `$USER` as a built-in variable containing the current username. Avoid spawning `whoami` as an external process.

Disable by adding `ZC1122` to `disabled_katas` in `.zshellcheckrc`.

---

<a id="zc1123"></a>
### ZC1123 — Use `$OSTYPE` instead of `uname`

**Severity:** `style`  
**Auto-fix:** `no`

Zsh provides `$OSTYPE` (e.g., `linux-gnu`, `darwin`) as a built-in variable. Avoid spawning `uname` for simple OS detection.

Disable by adding `ZC1123` to `disabled_katas` in `.zshellcheckrc`.

---

<a id="zc1124"></a>
### ZC1124 — Use `: > file` instead of `cat /dev/null > file` to truncate

**Severity:** `style`  
**Auto-fix:** `yes`

Truncating a file with `cat /dev/null > file` spawns an unnecessary process. Use `: > file` or simply `> file` in Zsh.

Disable by adding `ZC1124` to `disabled_katas` in `.zshellcheckrc`.

---

<a id="zc1125"></a>
### ZC1125 — Avoid `echo | grep` for string matching

**Severity:** `style`  
**Auto-fix:** `no`

Using `echo $var | grep pattern` spawns two unnecessary processes. Use Zsh `[[ $var =~ pattern ]]` or `[[ $var == *pattern* ]]` for string matching.

Disable by adding `ZC1125` to `disabled_katas` in `.zshellcheckrc`.

---

<a id="zc1126"></a>
### ZC1126 — Use `sort -u` instead of `sort | uniq`

**Severity:** `style`  
**Auto-fix:** `yes`

`sort | uniq` spawns two processes when `sort -u` does the same in one. Use `sort -u` to deduplicate sorted output efficiently.

Disable by adding `ZC1126` to `disabled_katas` in `.zshellcheckrc`.

---

<a id="zc1127"></a>
### ZC1127 — Avoid `ls` for counting files

**Severity:** `style`  
**Auto-fix:** `no`

Using `ls | wc -l` to count files spawns unnecessary processes. Use Zsh glob qualifiers: `files=(*(N)); echo ${#files}` for file counting.

Disable by adding `ZC1127` to `disabled_katas` in `.zshellcheckrc`.

---

<a id="zc1128"></a>
### ZC1128 — Use `> file` instead of `touch file` for creation

**Severity:** `style`  
**Auto-fix:** `yes`

If the goal is to create an empty file, `> file` does it without spawning `touch`. Use `touch` only when you need to update timestamps.

Disable by adding `ZC1128` to `disabled_katas` in `.zshellcheckrc`.

---

<a id="zc1129"></a>
### ZC1129 — Use Zsh `stat` module instead of `wc -c` for file size

**Severity:** `style`  
**Auto-fix:** `no`

Zsh's `zstat` (via `zmodload zsh/stat`) provides file size without spawning `wc`. Use `zstat +size file` for efficient file size queries.

Disable by adding `ZC1129` to `disabled_katas` in `.zshellcheckrc`.

---

<a id="zc1131"></a>
### ZC1131 — Avoid `cat file | while read` — use redirection

**Severity:** `style`  
**Auto-fix:** `no`

`cat file | while read line` spawns an unnecessary cat process and runs the loop in a subshell. Use `while read line; do ...; done < file` instead.

Disable by adding `ZC1131` to `disabled_katas` in `.zshellcheckrc`.

---

<a id="zc1132"></a>
### ZC1132 — Use Zsh pattern extraction instead of `grep -o`

**Severity:** `style`  
**Auto-fix:** `no`

For extracting matching parts from variables, use Zsh `${(M)var:#pattern}` or `${match[1]}` with `=~` instead of piping through `grep -o`.

Disable by adding `ZC1132` to `disabled_katas` in `.zshellcheckrc`.

---

<a id="zc1133"></a>
### ZC1133 — Avoid `kill -9` — use `kill` first, then escalate

**Severity:** `style`  
**Auto-fix:** `no`

`kill -9` (SIGKILL) cannot be caught or ignored. Always try `kill` (SIGTERM) first to allow the process to clean up, then use `kill -9` only as a last resort.

Disable by adding `ZC1133` to `disabled_katas` in `.zshellcheckrc`.

---

<a id="zc1134"></a>
### ZC1134 — Avoid `sleep` in tight loops

**Severity:** `style`  
**Auto-fix:** `no`

Using `sleep` inside a loop for polling creates busy-wait patterns. Consider `inotifywait`, `zle`, or event-driven approaches instead.

Disable by adding `ZC1134` to `disabled_katas` in `.zshellcheckrc`.

---

<a id="zc1135"></a>
### ZC1135 — Avoid `env VAR=val cmd` — use inline assignment

**Severity:** `style`  
**Auto-fix:** `yes`

Zsh supports inline environment variable assignment with `VAR=val cmd`. Avoid spawning `env` for simple variable-prefixed command execution.

Disable by adding `ZC1135` to `disabled_katas` in `.zshellcheckrc`.

---

<a id="zc1136"></a>
### ZC1136 — Avoid `rm -rf` without safeguard

**Severity:** `warning`  
**Auto-fix:** `no`

`rm -rf` with a variable path is dangerous if the variable is empty. Always validate the path or use `${var:?}` to fail on empty values.

Disable by adding `ZC1136` to `disabled_katas` in `.zshellcheckrc`.

---

<a id="zc1137"></a>
### ZC1137 — Avoid hardcoded `/tmp` paths

**Severity:** `style`  
**Auto-fix:** `no`

Hardcoded `/tmp` paths are predictable and may cause race conditions or symlink attacks. Use `mktemp` or Zsh `=(...)` for safe temp files.

Disable by adding `ZC1137` to `disabled_katas` in `.zshellcheckrc`.

---

<a id="zc1139"></a>
### ZC1139 — Avoid `source` with URL — use local files

**Severity:** `warning`  
**Auto-fix:** `no`

Sourcing scripts from URLs (curl | source) is a security risk. Download, verify, then source local files.

Disable by adding `ZC1139` to `disabled_katas` in `.zshellcheckrc`.

---

<a id="zc1140"></a>
### ZC1140 — Use `command -v` instead of `hash` for command existence

**Severity:** `style`  
**Auto-fix:** `yes`

`hash cmd` is a POSIX way to check command existence but provides poor error messages. Use `command -v cmd` for cleaner checks in Zsh.

Disable by adding `ZC1140` to `disabled_katas` in `.zshellcheckrc`.

---

<a id="zc1141"></a>
### ZC1141 — Avoid `curl | sh` pattern

**Severity:** `warning`  
**Auto-fix:** `no`

Piping curl output to sh/bash/zsh is a security risk. Download first, verify integrity (checksum or signature), then execute.

Disable by adding `ZC1141` to `disabled_katas` in `.zshellcheckrc`.

---

<a id="zc1142"></a>
### ZC1142 — Avoid chained `grep | grep` — combine patterns

**Severity:** `style`  
**Auto-fix:** `no`

Chaining `grep pattern1 | grep pattern2` spawns multiple processes. Use `grep -E 'p1.*p2|p2.*p1'` or `awk` for multi-pattern matching.

Disable by adding `ZC1142` to `disabled_katas` in `.zshellcheckrc`.

---

<a id="zc1143"></a>
### ZC1143 — Avoid `set -e` — use explicit error handling

**Severity:** `info`  
**Auto-fix:** `no`

`set -e` (errexit) has surprising behavior in Zsh with conditionals, pipes, and subshells. Use explicit `|| return` or `|| exit` for reliable error handling.

Disable by adding `ZC1143` to `disabled_katas` in `.zshellcheckrc`.

---

<a id="zc1144"></a>
### ZC1144 — Avoid `trap` with signal numbers — use names

**Severity:** `info`  
**Auto-fix:** `yes`

Signal numbers vary across platforms. Use signal names like `SIGTERM`, `SIGINT`, `EXIT` instead of numeric values for portability.

Disable by adding `ZC1144` to `disabled_katas` in `.zshellcheckrc`.

---

<a id="zc1145"></a>
### ZC1145 — Avoid `tr -d` for character deletion — use parameter expansion

**Severity:** `style`  
**Auto-fix:** `no`

For simple character deletion from variables, use Zsh `${var//char/}` instead of piping through `tr -d`.

Disable by adding `ZC1145` to `disabled_katas` in `.zshellcheckrc`.

---

<a id="zc1146"></a>
### ZC1146 — Avoid `cat file | awk` — pass file to awk directly

**Severity:** `style`  
**Auto-fix:** `yes`

`cat file | awk` spawns an unnecessary cat process. Pass the file directly as `awk '...' file`.

Disable by adding `ZC1146` to `disabled_katas` in `.zshellcheckrc`.

---

<a id="zc1147"></a>
### ZC1147 — Avoid `mkdir` without `-p` for nested paths

**Severity:** `info`  
**Auto-fix:** `yes`

Using `mkdir` without `-p` fails if parent directories don't exist. Use `mkdir -p` to create the full path safely.

Disable by adding `ZC1147` to `disabled_katas` in `.zshellcheckrc`.

---

<a id="zc1148"></a>
### ZC1148 — Use `compdef` instead of `compctl` for completions

**Severity:** `info`  
**Auto-fix:** `no`

`compctl` is the old Zsh completion system. Use `compdef` with the new completion system (`compsys`) for modern Zsh.

Disable by adding `ZC1148` to `disabled_katas` in `.zshellcheckrc`.

---

<a id="zc1149"></a>
### ZC1149 — Avoid `echo` for error messages — use `>&2`

**Severity:** `info`  
**Auto-fix:** `no`

Error messages should go to stderr, not stdout. Use `print -u2` or `echo ... >&2` to ensure errors are properly separated.

Disable by adding `ZC1149` to `disabled_katas` in `.zshellcheckrc`.

---

<a id="zc1151"></a>
### ZC1151 — Avoid `cat -A` — use `print -v` or od for non-printable characters

**Severity:** `style`  
**Auto-fix:** `no`

`cat -A` shows non-printable characters but varies across platforms. Use Zsh `print -v` or `od -c` for reliable non-printable character inspection.

Disable by adding `ZC1151` to `disabled_katas` in `.zshellcheckrc`.

---

<a id="zc1152"></a>
### ZC1152 — Use Zsh PCRE module instead of `grep -P`

**Severity:** `style`  
**Auto-fix:** `no`

`grep -P` (Perl regex) is not available on all platforms (e.g., macOS). Use `zmodload zsh/pcre` and `pcre_compile`/`pcre_match` for portable PCRE matching.

Disable by adding `ZC1152` to `disabled_katas` in `.zshellcheckrc`.

---

<a id="zc1153"></a>
### ZC1153 — Use `cmp -s` instead of `diff` for equality check

**Severity:** `style`  
**Auto-fix:** `yes`

When only checking if two files are identical (not viewing differences), `cmp -s` is faster than `diff` as it stops at the first difference.

Disable by adding `ZC1153` to `disabled_katas` in `.zshellcheckrc`.

---

<a id="zc1154"></a>
### ZC1154 — Use `find -exec {} +` instead of `find -exec {} \;`

**Severity:** `style`  
**Auto-fix:** `no`

`find -exec cmd {} \;` runs cmd once per file. `find -exec cmd {} +` batches files into fewer invocations, improving performance.

Disable by adding `ZC1154` to `disabled_katas` in `.zshellcheckrc`.

---

<a id="zc1155"></a>
### ZC1155 — Use `whence -a` instead of `which -a`

**Severity:** `info`  
**Auto-fix:** `yes`

`which -a` may be an external command on some systems. Zsh builtin `whence -a` reliably lists all command locations.

Disable by adding `ZC1155` to `disabled_katas` in `.zshellcheckrc`.

---

<a id="zc1156"></a>
### ZC1156 — Avoid `ln` without `-s` for symlinks

**Severity:** `info`  
**Auto-fix:** `no`

Hard links (`ln` without `-s`) share inodes and can cause confusion. Prefer symbolic links (`ln -s`) unless you specifically need hard links.

Disable by adding `ZC1156` to `disabled_katas` in `.zshellcheckrc`.

---

<a id="zc1157"></a>
### ZC1157 — Avoid `strings` command — use Zsh `${(ps:\0:)var}`

**Severity:** `style`  
**Auto-fix:** `no`

The `strings` command extracts printable strings from binaries. For simple filtering, Zsh parameter expansion with `(ps:\0:)` can split on null bytes.

Disable by adding `ZC1157` to `disabled_katas` in `.zshellcheckrc`.

---

<a id="zc1158"></a>
### ZC1158 — Avoid `chown -R` without `--no-dereference`

**Severity:** `warning`  
**Auto-fix:** `no`

`chown -R` follows symlinks by default, potentially changing ownership outside the intended directory. Use `--no-dereference` or `-h` to avoid this.

Disable by adding `ZC1158` to `disabled_katas` in `.zshellcheckrc`.

---

<a id="zc1159"></a>
### ZC1159 — Avoid `tar` without explicit compression flag

**Severity:** `info`  
**Auto-fix:** `no`

Use explicit compression flags (`-z` for gzip, `-j` for bzip2, `-J` for xz) instead of relying on `tar` auto-detection for clarity and portability.

Disable by adding `ZC1159` to `disabled_katas` in `.zshellcheckrc`.

---

<a id="zc1160"></a>
### ZC1160 — Prefer `curl` over `wget` for portability

**Severity:** `style`  
**Auto-fix:** `no`

`wget` is not installed by default on macOS. `curl` is available on virtually all Unix systems and is more portable.

Disable by adding `ZC1160` to `disabled_katas` in `.zshellcheckrc`.

---

<a id="zc1161"></a>
### ZC1161 — Avoid `openssl` for simple hashing — use Zsh modules

**Severity:** `style`  
**Auto-fix:** `no`

For simple SHA/MD5 hashing, Zsh provides `zmodload zsh/sha256` and `zmodload zsh/md5`. Avoid spawning `openssl` or `sha256sum` for basic hash operations.

Disable by adding `ZC1161` to `disabled_katas` in `.zshellcheckrc`.

---

<a id="zc1162"></a>
### ZC1162 — Use `cp -a` instead of `cp -r` to preserve attributes

**Severity:** `info`  
**Auto-fix:** `yes`

`cp -r` copies recursively but may not preserve permissions, timestamps, or symlinks. Use `cp -a` (archive mode) to preserve all attributes.

Disable by adding `ZC1162` to `disabled_katas` in `.zshellcheckrc`.

---

<a id="zc1163"></a>
### ZC1163 — Use `grep -m 1` instead of `grep | head -1`

**Severity:** `style`  
**Auto-fix:** `yes`

`grep pattern | head -1` spawns two processes when `grep -m 1` does the same. The `-m` flag stops after the first match, avoiding the pipeline.

Disable by adding `ZC1163` to `disabled_katas` in `.zshellcheckrc`.

---

<a id="zc1164"></a>
### ZC1164 — Avoid `sed -n 'Np'` — use Zsh array subscript

**Severity:** `style`  
**Auto-fix:** `no`

Extracting a specific line with `sed -n 'Np'` spawns a process. Use Zsh array subscript `${lines[N]}` after splitting with `${(f)...}`.

Disable by adding `ZC1164` to `disabled_katas` in `.zshellcheckrc`.

---

<a id="zc1165"></a>
### ZC1165 — Use Zsh parameter expansion for simple `awk` field extraction

**Severity:** `style`  
**Auto-fix:** `no`

Simple `awk '{print $1}'` or `awk '{print $NF}'` can often be replaced with Zsh parameter expansion `${var%% *}` (first field) or `${var##* }` (last field).

Disable by adding `ZC1165` to `disabled_katas` in `.zshellcheckrc`.

---

<a id="zc1166"></a>
### ZC1166 — Avoid `grep -i` for case-insensitive match — use `(#i)` glob flag

**Severity:** `style`  
**Auto-fix:** `no`

Zsh provides the `(#i)` glob flag for case-insensitive matching. For variable matching, use `[[ $var == (#i)pattern ]]` instead of piping through grep -i.

Disable by adding `ZC1166` to `disabled_katas` in `.zshellcheckrc`.

---

<a id="zc1167"></a>
### ZC1167 — Avoid `timeout` command — use Zsh `TMOUT` or `zsh/sched`

**Severity:** `style`  
**Auto-fix:** `no`

`timeout` is not available on all systems (macOS lacks it by default). Use Zsh `TMOUT` variable or `zmodload zsh/sched` for timeout functionality.

Disable by adding `ZC1167` to `disabled_katas` in `.zshellcheckrc`.

---

<a id="zc1168"></a>
### ZC1168 — Use `${(f)...}` instead of `readarray`/`mapfile`

**Severity:** `style`  
**Auto-fix:** `no`

`readarray` and `mapfile` are Bash builtins not available in Zsh. Use Zsh `${(f)...}` parameter expansion flag to split output into an array by newlines.

Disable by adding `ZC1168` to `disabled_katas` in `.zshellcheckrc`.

---

<a id="zc1169"></a>
### ZC1169 — Avoid `install` for simple copy+chmod — use `cp` then `chmod`

**Severity:** `style`  
**Auto-fix:** `no`

`install` command is less common and may confuse readers. For clarity, use separate `cp` and `chmod` commands or `install` only in Makefiles.

Disable by adding `ZC1169` to `disabled_katas` in `.zshellcheckrc`.

---

<a id="zc1170"></a>
### ZC1170 — Avoid `pushd`/`popd` without `-q` flag

**Severity:** `style`  
**Auto-fix:** `yes`

`pushd` and `popd` print the directory stack by default, cluttering output. Use `-q` flag to suppress output in scripts.

Disable by adding `ZC1170` to `disabled_katas` in `.zshellcheckrc`.

---

<a id="zc1171"></a>
### ZC1171 — Use `print` instead of `echo -e` for escape sequences

**Severity:** `style`  
**Auto-fix:** `yes`

`echo -e` behavior varies across shells and platforms. In Zsh, `print` natively interprets escape sequences and is more reliable.

Disable by adding `ZC1171` to `disabled_katas` in `.zshellcheckrc`.

---

<a id="zc1172"></a>
### ZC1172 — Use `read -A` instead of Bash `read -a` for arrays

**Severity:** `info`  
**Auto-fix:** `yes`

Bash uses `read -a` to read into an array, but Zsh uses `read -A`. Using `-a` in Zsh reads into a scalar, not an array.

Disable by adding `ZC1172` to `disabled_katas` in `.zshellcheckrc`.

---

<a id="zc1173"></a>
### ZC1173 — Avoid `column` command — use Zsh `print -C` for columnar output

**Severity:** `style`  
**Auto-fix:** `no`

Zsh `print -C N` formats output into N columns natively. Avoid spawning `column` as an external process for simple tabulation.

Disable by adding `ZC1173` to `disabled_katas` in `.zshellcheckrc`.

---

<a id="zc1174"></a>
### ZC1174 — Use Zsh `${(j:delim:)}` instead of `paste -sd`

**Severity:** `style`  
**Auto-fix:** `no`

Zsh `${(j:delim:)array}` joins array elements with a delimiter. Avoid spawning `paste` for simple field joining from variables.

Disable by adding `ZC1174` to `disabled_katas` in `.zshellcheckrc`.

---

<a id="zc1175"></a>
### ZC1175 — Avoid `tput` for simple ANSI colors — use Zsh `%F{color}`

**Severity:** `style`  
**Auto-fix:** `no`

Zsh prompt expansion `%F{red}` and `%f` handle colors natively. Avoid spawning `tput` for simple color output in prompts and scripts.

Disable by adding `ZC1175` to `disabled_katas` in `.zshellcheckrc`.

---

<a id="zc1176"></a>
### ZC1176 — Use `zparseopts` instead of `getopt`/`getopts`

**Severity:** `style`  
**Auto-fix:** `no`

Zsh provides `zparseopts` for powerful option parsing with long options, arrays, and defaults. Avoid `getopt`/`getopts` which are less capable in Zsh.

Disable by adding `ZC1176` to `disabled_katas` in `.zshellcheckrc`.

---

<a id="zc1177"></a>
### ZC1177 — Avoid `id -u` — use Zsh `$UID` or `$EUID`

**Severity:** `style`  
**Auto-fix:** `no`

Zsh provides `$UID` and `$EUID` as built-in variables for user/effective user ID. Avoid spawning `id` for simple UID checks.

Disable by adding `ZC1177` to `disabled_katas` in `.zshellcheckrc`.

---

<a id="zc1178"></a>
### ZC1178 — Avoid `stty` for terminal size — use Zsh `$COLUMNS`/`$LINES`

**Severity:** `style`  
**Auto-fix:** `no`

Zsh maintains `$COLUMNS` and `$LINES` as built-in variables tracking terminal dimensions. Avoid spawning `stty` or `tput` for size queries.

Disable by adding `ZC1178` to `disabled_katas` in `.zshellcheckrc`.

---

<a id="zc1179"></a>
### ZC1179 — Use Zsh `strftime` instead of `date` for formatting

**Severity:** `style`  
**Auto-fix:** `no`

Zsh provides `strftime` via `zmodload zsh/datetime` for date formatting. Avoid spawning `date` for simple timestamp formatting.

Disable by adding `ZC1179` to `disabled_katas` in `.zshellcheckrc`.

---

<a id="zc1180"></a>
### ZC1180 — Avoid `pgrep` for own background jobs — use Zsh job control

**Severity:** `info`  
**Auto-fix:** `no`

For managing your own background jobs, use Zsh job control (`jobs`, `kill %N`, `fg`, `bg`) instead of `pgrep`/`pkill` which search system-wide.

Disable by adding `ZC1180` to `disabled_katas` in `.zshellcheckrc`.

---

<a id="zc1181"></a>
### ZC1181 — Avoid `xdg-open`/`open` — use `$BROWSER` for portability

**Severity:** `info`  
**Auto-fix:** `no`

`xdg-open` is Linux-only, `open` is macOS-only. Use `$BROWSER` or check `$OSTYPE` for cross-platform URL/file opening.

Disable by adding `ZC1181` to `disabled_katas` in `.zshellcheckrc`.

---

<a id="zc1182"></a>
### ZC1182 — Avoid `nc`/`netcat` for HTTP — use `curl` or `zsh/net/tcp`

**Severity:** `warning`  
**Auto-fix:** `no`

`nc`/`netcat` for HTTP requests is fragile and lacks TLS support. Use `curl` or Zsh `zsh/net/tcp` module for reliable network operations.

Disable by adding `ZC1182` to `disabled_katas` in `.zshellcheckrc`.

---

<a id="zc1183"></a>
### ZC1183 — Use Zsh glob qualifiers instead of `ls -t` for file ordering

**Severity:** `style`  
**Auto-fix:** `no`

Zsh glob qualifiers like `*(om[1])` (newest) or `*(Om[1])` (oldest) order files without spawning `ls`. Avoid `ls -t | head` patterns.

Disable by adding `ZC1183` to `disabled_katas` in `.zshellcheckrc`.

---

<a id="zc1184"></a>
### ZC1184 — Avoid `diff -u` for patch generation — use `git diff` when in a repo

**Severity:** `style`  
**Auto-fix:** `no`

When working within a git repository, `git diff` provides better context, color output, and integration. Use `diff -u` only for non-repo file comparisons.

Disable by adding `ZC1184` to `disabled_katas` in `.zshellcheckrc`.

---

<a id="zc1185"></a>
### ZC1185 — Use Zsh `${#${(z)var}}` instead of `wc -w` for word count

**Severity:** `style`  
**Auto-fix:** `no`

Zsh `${(z)var}` splits a string into words and `${#...}` counts them. Avoid piping through `wc -w` for simple word counting from variables.

Disable by adding `ZC1185` to `disabled_katas` in `.zshellcheckrc`.

---

<a id="zc1186"></a>
### ZC1186 — Use `unset -v` or `unset -f` for explicit unsetting

**Severity:** `info`  
**Auto-fix:** `no`

Bare `unset name` is ambiguous — it unsets variables first, then functions. Use `unset -v` for variables or `unset -f` for functions to be explicit.

Disable by adding `ZC1186` to `disabled_katas` in `.zshellcheckrc`.

---

<a id="zc1187"></a>
### ZC1187 — Avoid `notify-send` without fallback — check availability first

**Severity:** `info`  
**Auto-fix:** `no`

`notify-send` is Linux-only (libnotify). For portable notifications, check `$OSTYPE` and fall back to `osascript` on macOS or `print` as default.

Disable by adding `ZC1187` to `disabled_katas` in `.zshellcheckrc`.

---

<a id="zc1188"></a>
### ZC1188 — Use Zsh `path+=()` instead of `export PATH=$PATH:dir`

**Severity:** `style`  
**Auto-fix:** `no`

Zsh ties the `path` array to `$PATH`. Use `path+=(dir)` to append directories cleanly instead of string manipulation with `export PATH=`.

Disable by adding `ZC1188` to `disabled_katas` in `.zshellcheckrc`.

---

<a id="zc1189"></a>
### ZC1189 — Avoid `source /dev/stdin` — use direct evaluation

**Severity:** `warning`  
**Auto-fix:** `no`

`source /dev/stdin` is fragile and platform-dependent. Use `eval "$(cmd)"` or direct command execution instead.

Disable by adding `ZC1189` to `disabled_katas` in `.zshellcheckrc`.

---

<a id="zc1190"></a>
### ZC1190 — Combine chained `grep -v` into single invocation

**Severity:** `style`  
**Auto-fix:** `yes`

`grep -v p1 | grep -v p2` spawns two processes. Use `grep -v -e p1 -e p2` to combine exclusions in one invocation.

Disable by adding `ZC1190` to `disabled_katas` in `.zshellcheckrc`.

---

<a id="zc1191"></a>
### ZC1191 — Avoid `clear` command — use ANSI escape sequences

**Severity:** `style`  
**Auto-fix:** `yes`

`clear` spawns an external process for screen clearing. Use `print -n '\e[2J\e[H'` for faster terminal clearing.

Disable by adding `ZC1191` to `disabled_katas` in `.zshellcheckrc`.

---

<a id="zc1192"></a>
### ZC1192 — Avoid `sleep 0` — it is a no-op external process

**Severity:** `info`  
**Auto-fix:** `yes`

`sleep 0` spawns an external process that does nothing. Remove it or use `:` if an explicit no-op is needed.

Disable by adding `ZC1192` to `disabled_katas` in `.zshellcheckrc`.

---

<a id="zc1193"></a>
### ZC1193 — Avoid `rm -i` in non-interactive scripts

**Severity:** `warning`  
**Auto-fix:** `no`

`rm -i` prompts for confirmation which hangs in non-interactive scripts. Remove the `-i` flag or use `rm -f` for scripts that run unattended.

Disable by adding `ZC1193` to `disabled_katas` in `.zshellcheckrc`.

---

<a id="zc1194"></a>
### ZC1194 — Avoid `sed` with multiple `-e` — use a single script

**Severity:** `style`  
**Auto-fix:** `no`

Multiple `sed -e 's/a/b/' -e 's/c/d/'` can be combined into `sed 's/a/b/; s/c/d/'` for cleaner syntax and fewer shell word splits.

Disable by adding `ZC1194` to `disabled_katas` in `.zshellcheckrc`.

---

<a id="zc1195"></a>
### ZC1195 — Avoid overly permissive `umask` values

**Severity:** `warning`  
**Auto-fix:** `no`

`umask 000` or `umask 0000` creates world-writable files by default. Use `umask 022` or more restrictive values for security.

Disable by adding `ZC1195` to `disabled_katas` in `.zshellcheckrc`.

---

<a id="zc1196"></a>
### ZC1196 — Avoid `cat` for reading single file into variable

**Severity:** `style`  
**Auto-fix:** `no`

Use Zsh `$(<file)` instead of `$(cat file)` to read file contents. `$(<file)` is a Zsh builtin that avoids spawning cat.

Disable by adding `ZC1196` to `disabled_katas` in `.zshellcheckrc`.

---

<a id="zc1197"></a>
### ZC1197 — Avoid `more` in scripts — use `cat` or pager check

**Severity:** `style`  
**Auto-fix:** `no`

`more` requires an interactive terminal and will hang in scripts. Use `cat` for output or check `$TERM` before invoking a pager.

Disable by adding `ZC1197` to `disabled_katas` in `.zshellcheckrc`.

---

<a id="zc1198"></a>
### ZC1198 — Avoid interactive editors in scripts

**Severity:** `warning`  
**Auto-fix:** `no`

`nano`, `vi`, and `vim` require interactive terminals and will hang in non-interactive scripts. Use `sed -i` or `ed` for scripted editing.

Disable by adding `ZC1198` to `disabled_katas` in `.zshellcheckrc`.

---

<a id="zc1199"></a>
### ZC1199 — Avoid `telnet` in scripts — use `curl` or `zsh/net/tcp`

**Severity:** `warning`  
**Auto-fix:** `no`

`telnet` is interactive and sends data in plain text. Use `curl` for HTTP or `zmodload zsh/net/tcp` for port checks in scripts.

Disable by adding `ZC1199` to `disabled_katas` in `.zshellcheckrc`.

---

<a id="zc1200"></a>
### ZC1200 — Avoid `ftp` — use `sftp` or `curl` for secure transfers

**Severity:** `warning`  
**Auto-fix:** `no`

`ftp` transmits credentials and data in plain text. Use `sftp`, `scp`, or `curl` with HTTPS/SFTP for secure file transfers.

Disable by adding `ZC1200` to `disabled_katas` in `.zshellcheckrc`.

---

<a id="zc1201"></a>
### ZC1201 — Avoid `rsh`/`rlogin`/`rcp` — use `ssh`/`scp`

**Severity:** `warning`  
**Auto-fix:** `no`

`rsh`, `rlogin`, and `rcp` are insecure legacy protocols. Use `ssh`, `scp`, or `rsync` over SSH for encrypted remote operations.

Disable by adding `ZC1201` to `disabled_katas` in `.zshellcheckrc`.

---

<a id="zc1202"></a>
### ZC1202 — Avoid `ifconfig` — use `ip` for network configuration

**Severity:** `info`  
**Auto-fix:** `yes`

`ifconfig` is deprecated on modern Linux. Use `ip addr`, `ip link`, or `ip route` from iproute2 for network operations.

Disable by adding `ZC1202` to `disabled_katas` in `.zshellcheckrc`.

---

<a id="zc1203"></a>
### ZC1203 — Avoid `netstat` — use `ss` for socket statistics

**Severity:** `info`  
**Auto-fix:** `yes`

`netstat` is deprecated on modern Linux in favor of `ss` from iproute2. `ss` is faster and provides more detailed socket information.

Disable by adding `ZC1203` to `disabled_katas` in `.zshellcheckrc`.

---

<a id="zc1204"></a>
### ZC1204 — Avoid `route` — use `ip route` for routing

**Severity:** `info`  
**Auto-fix:** `no`

`route` is deprecated on modern Linux in favor of `ip route` from iproute2. `ip route` provides consistent syntax with other `ip` subcommands.

Disable by adding `ZC1204` to `disabled_katas` in `.zshellcheckrc`.

---

<a id="zc1205"></a>
### ZC1205 — Avoid `arp` — use `ip neigh` for neighbor tables

**Severity:** `info`  
**Auto-fix:** `no`

`arp` is deprecated on modern Linux in favor of `ip neigh` from iproute2. `ip neigh` provides consistent syntax with other `ip` subcommands.

Disable by adding `ZC1205` to `disabled_katas` in `.zshellcheckrc`.

---

<a id="zc1206"></a>
### ZC1206 — Avoid `crontab -e` in scripts — use `crontab file`

**Severity:** `warning`  
**Auto-fix:** `no`

`crontab -e` opens an interactive editor which hangs in scripts. Use `crontab file` or pipe content with `crontab -` for programmatic cron management.

Disable by adding `ZC1206` to `disabled_katas` in `.zshellcheckrc`.

---

<a id="zc1207"></a>
### ZC1207 — Avoid `passwd` in scripts — use `chpasswd`

**Severity:** `warning`  
**Auto-fix:** `no`

`passwd` prompts interactively for password input. Use `chpasswd` or `usermod --password` for non-interactive password changes.

Disable by adding `ZC1207` to `disabled_katas` in `.zshellcheckrc`.

---

<a id="zc1208"></a>
### ZC1208 — Avoid `visudo` in scripts — use sudoers.d drop-in files

**Severity:** `warning`  
**Auto-fix:** `no`

`visudo` opens an interactive editor. For programmatic sudoers changes, write to `/etc/sudoers.d/` drop-in files with `visudo -c` for validation.

Disable by adding `ZC1208` to `disabled_katas` in `.zshellcheckrc`.

---

<a id="zc1209"></a>
### ZC1209 — Use `systemctl --no-pager` in scripts

**Severity:** `style`  
**Auto-fix:** `yes`

`systemctl` invokes a pager by default which hangs in non-interactive scripts. Use `--no-pager` or pipe to `cat` for reliable script output.

Disable by adding `ZC1209` to `disabled_katas` in `.zshellcheckrc`.

---

<a id="zc1210"></a>
### ZC1210 — Use `journalctl --no-pager` in scripts

**Severity:** `style`  
**Auto-fix:** `yes`

`journalctl` invokes a pager by default which hangs in non-interactive scripts. Use `--no-pager` for reliable script output.

Disable by adding `ZC1210` to `disabled_katas` in `.zshellcheckrc`.

---

<a id="zc1211"></a>
### ZC1211 — Use `git stash push -m` instead of bare `git stash`

**Severity:** `style`  
**Auto-fix:** `no`

Bare `git stash` creates unnamed stashes that are hard to identify later. Use `git stash push -m 'description'` for self-documenting stashes.

Disable by adding `ZC1211` to `disabled_katas` in `.zshellcheckrc`.

---

<a id="zc1212"></a>
### ZC1212 — Avoid `git add .` — use explicit paths or `git add -p`

**Severity:** `info`  
**Auto-fix:** `no`

`git add .` stages everything including unintended files. Use explicit file paths or `git add -p` for selective staging.

Disable by adding `ZC1212` to `disabled_katas` in `.zshellcheckrc`.

---

<a id="zc1213"></a>
### ZC1213 — Use `apt-get -y` in scripts for non-interactive installs

**Severity:** `warning`  
**Auto-fix:** `yes`

`apt-get install` without `-y` prompts for confirmation which hangs scripts. Use `-y` or set `DEBIAN_FRONTEND=noninteractive` for unattended installs.

Disable by adding `ZC1213` to `disabled_katas` in `.zshellcheckrc`.

---

<a id="zc1214"></a>
### ZC1214 — Avoid `su` in scripts — use `sudo -u` for user switching

**Severity:** `warning`  
**Auto-fix:** `no`

`su` prompts for a password interactively which hangs scripts. Use `sudo -u user cmd` for non-interactive privilege switching.

Disable by adding `ZC1214` to `disabled_katas` in `.zshellcheckrc`.

---

<a id="zc1215"></a>
### ZC1215 — Source `/etc/os-release` instead of parsing with `cat`/`grep`

**Severity:** `style`  
**Auto-fix:** `yes`

`/etc/os-release` is designed to be sourced directly. Use `. /etc/os-release` to get variables like `$ID`, `$VERSION_ID` without parsing.

Disable by adding `ZC1215` to `disabled_katas` in `.zshellcheckrc`.

---

<a id="zc1216"></a>
### ZC1216 — Avoid `nslookup` — use `dig` or `host` for DNS queries

**Severity:** `info`  
**Auto-fix:** `yes`

`nslookup` is deprecated in many distributions. `dig` provides more detailed output and `host` is simpler for basic lookups.

Disable by adding `ZC1216` to `disabled_katas` in `.zshellcheckrc`.

---

<a id="zc1217"></a>
### ZC1217 — Avoid `service` command — use `systemctl` on systemd

**Severity:** `info`  
**Auto-fix:** `yes`

`service` is a SysVinit compatibility wrapper. On systemd systems, use `systemctl start/stop/restart/status` directly.

Disable by adding `ZC1217` to `disabled_katas` in `.zshellcheckrc`.

---

<a id="zc1218"></a>
### ZC1218 — Avoid `useradd` without `--shell /sbin/nologin` for service accounts

**Severity:** `warning`  
**Auto-fix:** `no`

Service accounts created with `useradd` should use `--shell /sbin/nologin` and `--system` to prevent interactive login and use system UID ranges.

Disable by adding `ZC1218` to `disabled_katas` in `.zshellcheckrc`.

---

<a id="zc1219"></a>
### ZC1219 — Use `curl -fsSL` instead of `wget -O -` for piped downloads

**Severity:** `style`  
**Auto-fix:** `yes`

`wget -O -` outputs to stdout but lacks `curl`'s error handling. `curl -fsSL` fails on HTTP errors, is silent, follows redirects, and is more portable.

Disable by adding `ZC1219` to `disabled_katas` in `.zshellcheckrc`.

---

<a id="zc1220"></a>
### ZC1220 — Use `chown :group` instead of `chgrp` for group changes

**Severity:** `style`  
**Auto-fix:** `no`

`chgrp` is redundant when `chown :group file` does the same thing. Using `chown` for both user and group changes is more consistent.

Disable by adding `ZC1220` to `disabled_katas` in `.zshellcheckrc`.

---

<a id="zc1221"></a>
### ZC1221 — Avoid `fdisk` in scripts — use `parted` or `sfdisk`

**Severity:** `warning`  
**Auto-fix:** `no`

`fdisk` is interactive and not scriptable. Use `parted -s` or `sfdisk` for non-interactive disk partitioning.

Disable by adding `ZC1221` to `disabled_katas` in `.zshellcheckrc`.

---

<a id="zc1222"></a>
### ZC1222 — Avoid `lsof -i` for port checks — use `ss -tlnp`

**Severity:** `style`  
**Auto-fix:** `no`

`lsof -i` is slow and requires elevated permissions on some systems. `ss -tlnp` is faster and part of the standard iproute2 toolkit.

Disable by adding `ZC1222` to `disabled_katas` in `.zshellcheckrc`.

---

<a id="zc1223"></a>
### ZC1223 — Avoid `ip addr show` piped to `grep` — use `ip -br addr`

**Severity:** `style`  
**Auto-fix:** `no`

`ip addr show | grep` parses verbose output. `ip -br addr` provides machine-readable brief output without needing grep.

Disable by adding `ZC1223` to `disabled_katas` in `.zshellcheckrc`.

---

<a id="zc1224"></a>
### ZC1224 — Avoid parsing `free` output — read `/proc/meminfo` directly

**Severity:** `style`  
**Auto-fix:** `no`

`free` output format varies across versions and locales. Read `/proc/meminfo` directly for reliable memory information in scripts.

Disable by adding `ZC1224` to `disabled_katas` in `.zshellcheckrc`.

---

<a id="zc1225"></a>
### ZC1225 — Avoid parsing `uptime` — read `/proc/uptime` directly

**Severity:** `style`  
**Auto-fix:** `no`

`uptime` output is human-readable and varies by locale. Read `/proc/uptime` for machine-parseable uptime in seconds.

Disable by adding `ZC1225` to `disabled_katas` in `.zshellcheckrc`.

---

<a id="zc1226"></a>
### ZC1226 — Use `dmesg -T` or `--time-format=iso` for readable timestamps

**Severity:** `style`  
**Auto-fix:** `yes`

`dmesg` without `-T` shows raw kernel timestamps in seconds since boot. Use `-T` for human-readable timestamps or `--time-format=iso` for ISO 8601.

Disable by adding `ZC1226` to `disabled_katas` in `.zshellcheckrc`.

---

<a id="zc1227"></a>
### ZC1227 — Use `curl -f` to fail on HTTP errors

**Severity:** `warning`  
**Auto-fix:** `yes`

`curl` without `-f` silently returns error pages (404, 500) as success. Use `-f` or `--fail` to return exit code 22 on HTTP errors.

Disable by adding `ZC1227` to `disabled_katas` in `.zshellcheckrc`.

---

<a id="zc1228"></a>
### ZC1228 — Avoid `ssh` without host key policy in scripts

**Severity:** `warning`  
**Auto-fix:** `no`

`ssh` without `-o BatchMode=yes` or `-o StrictHostKeyChecking` prompts interactively for host key verification, hanging non-interactive scripts.

Disable by adding `ZC1228` to `disabled_katas` in `.zshellcheckrc`.

---

<a id="zc1229"></a>
### ZC1229 — Prefer `rsync` over `scp` for file transfers

**Severity:** `style`  
**Auto-fix:** `no`

`scp` uses a deprecated protocol and lacks delta transfer, resume, and progress features. `rsync` is more efficient and reliable for scripts.

Disable by adding `ZC1229` to `disabled_katas` in `.zshellcheckrc`.

---

<a id="zc1230"></a>
### ZC1230 — Use `ping -c N` in scripts to limit ping count

**Severity:** `warning`  
**Auto-fix:** `yes`

`ping` without `-c` runs indefinitely on Linux, hanging scripts. Always specify `-c N` to limit the number of packets.

Disable by adding `ZC1230` to `disabled_katas` in `.zshellcheckrc`.

---

<a id="zc1231"></a>
### ZC1231 — Use `git clone --depth 1` for CI and build scripts

**Severity:** `style`  
**Auto-fix:** `yes`

`git clone` without `--depth` downloads the entire history. Use `--depth 1` in CI/build scripts where only the latest commit is needed.

Disable by adding `ZC1231` to `disabled_katas` in `.zshellcheckrc`.

---

<a id="zc1232"></a>
### ZC1232 — Avoid bare `pip install` — use `--user` or virtualenv

**Severity:** `warning`  
**Auto-fix:** `no`

Bare `pip install` may modify system Python packages. Use `pip install --user`, `pipx`, or a virtualenv to isolate dependencies.

Disable by adding `ZC1232` to `disabled_katas` in `.zshellcheckrc`.

---

<a id="zc1233"></a>
### ZC1233 — Avoid `npm install -g` — use `npx` for one-off tools

**Severity:** `style`  
**Auto-fix:** `no`

Global npm installs pollute the system. Use `npx` to run tools without installing, or `npm install --save-dev` for project dependencies.

Disable by adding `ZC1233` to `disabled_katas` in `.zshellcheckrc`.

---

<a id="zc1234"></a>
### ZC1234 — Use `docker run --rm` to auto-remove containers

**Severity:** `style`  
**Auto-fix:** `yes`

`docker run` without `--rm` leaves stopped containers behind. Use `--rm` in scripts to automatically clean up after execution.

Disable by adding `ZC1234` to `disabled_katas` in `.zshellcheckrc`.

---

<a id="zc1235"></a>
### ZC1235 — Use `git push --force-with-lease` instead of `--force`

**Severity:** `warning`  
**Auto-fix:** `yes`

`git push --force` overwrites remote history unconditionally. `--force-with-lease` is safer as it fails if the remote has changed.

Disable by adding `ZC1235` to `disabled_katas` in `.zshellcheckrc`.

---

<a id="zc1236"></a>
### ZC1236 — Avoid `git reset --hard` — irreversible data loss risk

**Severity:** `warning`  
**Auto-fix:** `no`

`git reset --hard` discards all uncommitted changes irreversibly. Use `git stash` to save changes first, or `git reset --soft` to keep them staged.

Disable by adding `ZC1236` to `disabled_katas` in `.zshellcheckrc`.

---

<a id="zc1237"></a>
### ZC1237 — Use `git clean -n` before `git clean -fd`

**Severity:** `warning`  
**Auto-fix:** `no`

`git clean -fd` permanently deletes untracked files and directories. Use `-n` (dry run) first to preview what will be removed.

Disable by adding `ZC1237` to `disabled_katas` in `.zshellcheckrc`.

---

<a id="zc1238"></a>
### ZC1238 — Avoid `docker exec -it` in scripts — drop `-it` for non-interactive

**Severity:** `warning`  
**Auto-fix:** `yes`

`docker exec -it` allocates a TTY and attaches stdin, which hangs in non-interactive scripts. Use `docker exec` without `-it` for scripted commands.

Disable by adding `ZC1238` to `disabled_katas` in `.zshellcheckrc`.

---

<a id="zc1239"></a>
### ZC1239 — Avoid `kubectl exec -it` in scripts

**Severity:** `warning`  
**Auto-fix:** `yes`

`kubectl exec -it` allocates a TTY which hangs in non-interactive scripts. Use `kubectl exec` without `-it` or use `kubectl exec -- cmd` for scripted commands.

Disable by adding `ZC1239` to `disabled_katas` in `.zshellcheckrc`.

---

<a id="zc1240"></a>
### ZC1240 — Use `find -maxdepth` with `-delete` to limit scope

**Severity:** `warning`  
**Auto-fix:** `no`

`find -delete` without `-maxdepth` recurses infinitely and may delete more than intended. Always limit the search depth.

Disable by adding `ZC1240` to `disabled_katas` in `.zshellcheckrc`.

---

<a id="zc1241"></a>
### ZC1241 — Use `xargs -0` with null separators for safe argument passing

**Severity:** `warning`  
**Auto-fix:** `yes`

`xargs` without `-0` splits on whitespace, breaking on filenames with spaces. Use `xargs -0` paired with `find -print0` for safe handling.

Disable by adding `ZC1241` to `disabled_katas` in `.zshellcheckrc`.

---

<a id="zc1242"></a>
### ZC1242 — Use `tar -C dir` to extract into a specific directory

**Severity:** `info`  
**Auto-fix:** `no`

`tar xf` without `-C` extracts into the current directory which may overwrite files unexpectedly. Use `-C dir` to control the extraction target.

Disable by adding `ZC1242` to `disabled_katas` in `.zshellcheckrc`.

---

<a id="zc1243"></a>
### ZC1243 — Use `grep -lZ` with `xargs -0` for safe file lists

**Severity:** `warning`  
**Auto-fix:** `no`

`grep -l` outputs one filename per line, breaking on names with newlines. Use `grep -lZ` (null-terminated) paired with `xargs -0` for safe processing.

Disable by adding `ZC1243` to `disabled_katas` in `.zshellcheckrc`.

---

<a id="zc1244"></a>
### ZC1244 — Consider `mv -n` to prevent overwriting existing files

**Severity:** `info`  
**Auto-fix:** `no`

`mv` overwrites existing files without warning by default. Use `-n` (no-clobber) to prevent accidental overwrites in scripts.

Disable by adding `ZC1244` to `disabled_katas` in `.zshellcheckrc`.

---

<a id="zc1245"></a>
### ZC1245 — Avoid disabling TLS certificate verification

**Severity:** `error`  
**Auto-fix:** `no`

Flags like `--no-check-certificate` (wget) or `-k`/`--insecure` (curl) disable TLS verification, making connections vulnerable to MITM attacks.

Disable by adding `ZC1245` to `disabled_katas` in `.zshellcheckrc`.

---

<a id="zc1246"></a>
### ZC1246 — Avoid hardcoded passwords in command arguments

**Severity:** `error`  
**Auto-fix:** `no`

Passing passwords as command arguments exposes them in process lists and shell history. Use environment variables or credential files instead.

Disable by adding `ZC1246` to `disabled_katas` in `.zshellcheckrc`.

---

<a id="zc1247"></a>
### ZC1247 — Avoid `chmod +s` — setuid/setgid bits are security risks

**Severity:** `error`  
**Auto-fix:** `no`

Setting the setuid or setgid bit (`chmod +s` or `chmod u+s`) allows files to execute with the owner's privileges, creating privilege escalation risks.

Disable by adding `ZC1247` to `disabled_katas` in `.zshellcheckrc`.

---

<a id="zc1248"></a>
### ZC1248 — Prefer `ufw`/`firewalld` over raw `iptables`

**Severity:** `info`  
**Auto-fix:** `no`

Raw `iptables` rules are complex and non-persistent by default. Use `ufw` (Ubuntu) or `firewalld` (RHEL) for manageable, persistent firewall rules.

Disable by adding `ZC1248` to `disabled_katas` in `.zshellcheckrc`.

---

<a id="zc1249"></a>
### ZC1249 — Use `ssh-keygen -f` to specify key file in scripts

**Severity:** `warning`  
**Auto-fix:** `no`

`ssh-keygen` without `-f` prompts for a file path interactively. Use `-f /path/to/key` and `-N ''` for non-interactive key generation.

Disable by adding `ZC1249` to `disabled_katas` in `.zshellcheckrc`.

---

<a id="zc1250"></a>
### ZC1250 — Use `gpg --batch` in scripts for non-interactive operation

**Severity:** `warning`  
**Auto-fix:** `no`

`gpg` without `--batch` may prompt for passphrases or confirmations. Use `--batch` and `--yes` for fully non-interactive GPG operations in scripts.

Disable by adding `ZC1250` to `disabled_katas` in `.zshellcheckrc`.

---

<a id="zc1251"></a>
### ZC1251 — Use `mount -o noexec,nosuid` for untrusted media

**Severity:** `warning`  
**Auto-fix:** `no`

Mounting untrusted filesystems without `noexec,nosuid` allows execution of malicious binaries and setuid exploits. Always restrict mount options.

Disable by adding `ZC1251` to `disabled_katas` in `.zshellcheckrc`.

---

<a id="zc1252"></a>
### ZC1252 — Use `getent passwd` instead of `cat /etc/passwd`

**Severity:** `style`  
**Auto-fix:** `yes`

`cat /etc/passwd` misses users from LDAP, NIS, or SSSD sources. `getent passwd` queries NSS and returns all configured user databases.

Disable by adding `ZC1252` to `disabled_katas` in `.zshellcheckrc`.

---

<a id="zc1253"></a>
### ZC1253 — Use `docker build --no-cache` in CI for reproducible builds

**Severity:** `style`  
**Auto-fix:** `yes`

`docker build` uses layer caching which can mask dependency changes. Use `--no-cache` in CI pipelines to ensure fully reproducible builds.

Disable by adding `ZC1253` to `disabled_katas` in `.zshellcheckrc`.

---

<a id="zc1254"></a>
### ZC1254 — Avoid `git commit --amend` in shared branches

**Severity:** `warning`  
**Auto-fix:** `no`

`git commit --amend` rewrites the last commit which causes problems if already pushed. Use `git commit --fixup` or a new commit instead.

Disable by adding `ZC1254` to `disabled_katas` in `.zshellcheckrc`.

---

<a id="zc1255"></a>
### ZC1255 — Use `curl -L` to follow HTTP redirects

**Severity:** `info`  
**Auto-fix:** `yes`

`curl` without `-L` does not follow redirects, returning 301/302 responses instead of the actual content. Use `-L` to follow redirects automatically.

Disable by adding `ZC1255` to `disabled_katas` in `.zshellcheckrc`.

---

<a id="zc1256"></a>
### ZC1256 — Clean up `mkfifo` pipes with a trap on EXIT

**Severity:** `info`  
**Auto-fix:** `no`

`mkfifo` creates named pipes that persist on the filesystem. Set up a `trap` to remove them on EXIT to prevent leftover files.

Disable by adding `ZC1256` to `disabled_katas` in `.zshellcheckrc`.

---

<a id="zc1257"></a>
### ZC1257 — Use `docker stop -t` to set graceful shutdown timeout

**Severity:** `style`  
**Auto-fix:** `yes`

`docker stop` defaults to 10s before SIGKILL. In CI scripts, set an explicit timeout with `-t` to control shutdown behavior.

Disable by adding `ZC1257` to `disabled_katas` in `.zshellcheckrc`.

---

<a id="zc1258"></a>
### ZC1258 — Consider `rsync --delete` for directory sync

**Severity:** `warning`  
**Auto-fix:** `no`

`rsync` without `--delete` keeps files on the destination that were removed from the source. Use `--delete` for true directory mirroring.

Disable by adding `ZC1258` to `disabled_katas` in `.zshellcheckrc`.

---

<a id="zc1259"></a>
### ZC1259 — Avoid `docker pull` without explicit tag — pin image versions

**Severity:** `warning`  
**Auto-fix:** `no`

`docker pull image` without a tag defaults to `:latest` which is mutable and non-reproducible. Always pin to a specific version tag or digest.

Disable by adding `ZC1259` to `disabled_katas` in `.zshellcheckrc`.

---

<a id="zc1260"></a>
### ZC1260 — Use `git branch -d` instead of `-D` for safe deletion

**Severity:** `warning`  
**Auto-fix:** `yes`

`git branch -D` force-deletes branches even if unmerged. Use `-d` which refuses to delete unmerged branches, preventing data loss.

Disable by adding `ZC1260` to `disabled_katas` in `.zshellcheckrc`.

---

<a id="zc1261"></a>
### ZC1261 — Avoid piping `base64 -d` output to shell execution

**Severity:** `error`  
**Auto-fix:** `no`

Decoding base64 and piping to `sh`/`zsh`/`eval` is a code injection risk. Always inspect decoded content before execution.

Disable by adding `ZC1261` to `disabled_katas` in `.zshellcheckrc`.

---

<a id="zc1262"></a>
### ZC1262 — Avoid `chmod -R 777` — recursive world-writable is critical

**Severity:** `error`  
**Auto-fix:** `no`

`chmod -R 777` makes every file and directory world-writable and executable. Use specific permissions like `755` for directories and `644` for files.

Disable by adding `ZC1262` to `disabled_katas` in `.zshellcheckrc`.

---

<a id="zc1263"></a>
### ZC1263 — Use `apt-get` instead of `apt` in scripts

**Severity:** `style`  
**Auto-fix:** `yes`

`apt` is designed for interactive use and its output format may change. `apt-get` has a stable interface suitable for scripts and CI.

Disable by adding `ZC1263` to `disabled_katas` in `.zshellcheckrc`.

---

<a id="zc1264"></a>
### ZC1264 — Use `dnf` instead of `yum` on modern Fedora/RHEL

**Severity:** `style`  
**Auto-fix:** `yes`

`yum` is deprecated on Fedora 22+ and RHEL 8+. `dnf` is the modern replacement with better dependency resolution.

Disable by adding `ZC1264` to `disabled_katas` in `.zshellcheckrc`.

---

<a id="zc1265"></a>
### ZC1265 — Use `systemctl enable --now` to enable and start together

**Severity:** `style`  
**Auto-fix:** `yes`

`systemctl enable` without `--now` only enables on next boot. Use `--now` to enable and immediately start the service.

Disable by adding `ZC1265` to `disabled_katas` in `.zshellcheckrc`.

---

<a id="zc1266"></a>
### ZC1266 — Use `nproc` instead of parsing `/proc/cpuinfo`

**Severity:** `style`  
**Auto-fix:** `no`

Parsing `/proc/cpuinfo` for CPU count is fragile and platform-specific. `nproc` is a portable, dedicated tool for this purpose.

Disable by adding `ZC1266` to `disabled_katas` in `.zshellcheckrc`.

---

<a id="zc1267"></a>
### ZC1267 — Use `df -P` for POSIX-portable disk usage output

**Severity:** `style`  
**Auto-fix:** `yes`

`df -h` output format varies across systems and locales. Use `df -P` for single-line, fixed-format output safe for script parsing.

Disable by adding `ZC1267` to `disabled_katas` in `.zshellcheckrc`.

---

<a id="zc1268"></a>
### ZC1268 — Use `du -sh --` to handle filenames starting with dash

**Severity:** `info`  
**Auto-fix:** `yes`

`du -sh *` breaks if a filename starts with `-`. Use `--` to signal end of options and safely handle all filenames.

Disable by adding `ZC1268` to `disabled_katas` in `.zshellcheckrc`.

---

<a id="zc1269"></a>
### ZC1269 — Use `pgrep` instead of `ps aux | grep` for process search

**Severity:** `style`  
**Auto-fix:** `no`

`ps aux | grep` matches itself in the process list requiring workarounds. Use `pgrep` which is designed for process searching without self-matching.

Disable by adding `ZC1269` to `disabled_katas` in `.zshellcheckrc`.

---

<a id="zc1270"></a>
### ZC1270 — Use `mktemp` instead of hardcoded `/tmp` paths

**Severity:** `warning`  
**Auto-fix:** `no`

Hardcoding `/tmp/filename` is vulnerable to symlink attacks and race conditions. Use `mktemp` to create unique temporary files safely.

Disable by adding `ZC1270` to `disabled_katas` in `.zshellcheckrc`.

---

<a id="zc1271"></a>
### ZC1271 — Use `command -v` instead of `which` for command existence checks

**Severity:** `style`  
**Auto-fix:** `yes`

`which` is not POSIX-standard and behaves inconsistently across systems. Use `command -v` which is portable and built into Zsh.

Disable by adding `ZC1271` to `disabled_katas` in `.zshellcheckrc`.

---

<a id="zc1272"></a>
### ZC1272 — Use `install -m` instead of separate `cp` and `chmod`

**Severity:** `style`  
**Auto-fix:** `no`

`install` atomically copies a file and sets permissions in one step. Using separate `cp` and `chmod` creates a window where the file has wrong permissions.

Disable by adding `ZC1272` to `disabled_katas` in `.zshellcheckrc`.

---

<a id="zc1273"></a>
### ZC1273 — Use `grep -q` instead of redirecting grep output to `/dev/null`

**Severity:** `style`  
**Auto-fix:** `yes`

`grep -q` suppresses output and exits on first match, which is faster and more idiomatic than piping or redirecting to `/dev/null`.

Disable by adding `ZC1273` to `disabled_katas` in `.zshellcheckrc`.

---

<a id="zc1274"></a>
### ZC1274 — Use Zsh `${var:t}` instead of `basename`

**Severity:** `style`  
**Auto-fix:** `no`

Zsh provides the `:t` (tail) modifier for parameter expansion which extracts the filename component, avoiding the overhead of forking `basename`.

Disable by adding `ZC1274` to `disabled_katas` in `.zshellcheckrc`.

---

<a id="zc1275"></a>
### ZC1275 — Use Zsh `${var:h}` instead of `dirname`

**Severity:** `style`  
**Auto-fix:** `no`

Zsh provides the `:h` (head) modifier for parameter expansion which extracts the directory component, avoiding the overhead of forking `dirname`.

Disable by adding `ZC1275` to `disabled_katas` in `.zshellcheckrc`.

---

<a id="zc1276"></a>
### ZC1276 — Use Zsh `{start..end}` instead of `seq`

**Severity:** `style`  
**Auto-fix:** `yes`

Zsh natively supports `{start..end}` brace expansion for generating number sequences, avoiding the overhead of forking the external `seq` command.

Disable by adding `ZC1276` to `disabled_katas` in `.zshellcheckrc`.

---

<a id="zc1277"></a>
### ZC1277 — Superseded by ZC1108 — retired duplicate

**Severity:** `style`  
**Auto-fix:** `no`

Retained as a no-op stub so legacy `.zshellcheckrc` files that disable this ID keep parsing. See https://github.com/afadesigns/zshellcheck/issues/344 for context; the canonical detection lives in ZC1108.

Disable by adding `ZC1277` to `disabled_katas` in `.zshellcheckrc`.

---

<a id="zc1278"></a>
### ZC1278 — Superseded by ZC1009 — retired duplicate

**Severity:** `style`  
**Auto-fix:** `no`

Retained as a no-op stub so legacy `.zshellcheckrc` files that disable this ID keep parsing. See https://github.com/afadesigns/zshellcheck/issues/343 for context; the canonical detection lives in ZC1009.

Disable by adding `ZC1278` to `disabled_katas` in `.zshellcheckrc`.

---

<a id="zc1279"></a>
### ZC1279 — Use `realpath` instead of `readlink -f` for canonical paths

**Severity:** `info`  
**Auto-fix:** `yes`

`readlink -f` is not portable across all platforms (notably macOS). Use `realpath` which is POSIX-standard and available on modern systems.

Disable by adding `ZC1279` to `disabled_katas` in `.zshellcheckrc`.

---

<a id="zc1280"></a>
### ZC1280 — Use `Zsh ${var:e}` instead of shell expansion to extract file extension

**Severity:** `style`  
**Auto-fix:** `no`

Zsh provides the `:e` (extension) modifier for parameter expansion which extracts the file extension, avoiding complex shell patterns or external tools.

Disable by adding `ZC1280` to `disabled_katas` in `.zshellcheckrc`.

---

<a id="zc1281"></a>
### ZC1281 — Use `sort -u` instead of `sort | uniq` for deduplication

**Severity:** `style`  
**Auto-fix:** `no`

`sort -u` combines sorting and deduplication in a single pass, which is more efficient than piping `sort` into `uniq` as a separate process.

Disable by adding `ZC1281` to `disabled_katas` in `.zshellcheckrc`.

---

<a id="zc1282"></a>
### ZC1282 — Use Zsh `${var:r}` instead of `sed` to remove file extension

**Severity:** `style`  
**Auto-fix:** `no`

Zsh provides the `:r` modifier to remove a filename extension. Using `sed` or `cut` to strip the extension is unnecessary when the built-in parameter expansion handles it directly.

Disable by adding `ZC1282` to `disabled_katas` in `.zshellcheckrc`.

---

<a id="zc1283"></a>
### ZC1283 — Use `setopt` instead of `set -o` for Zsh options

**Severity:** `style`  
**Auto-fix:** `yes`

Zsh provides `setopt` and `unsetopt` as native builtins for managing shell options. Using `set -o` / `set +o` is a POSIX compatibility form that is less idiomatic in Zsh scripts.

Disable by adding `ZC1283` to `disabled_katas` in `.zshellcheckrc`.

---

<a id="zc1284"></a>
### ZC1284 — Use Zsh `${(s:sep:)var}` instead of `cut -d` for field splitting

**Severity:** `style`  
**Auto-fix:** `no`

Zsh provides the `(s:separator:)` parameter expansion flag to split strings into arrays by a delimiter. This is more idiomatic than invoking `cut -d` and avoids spawning an external process.

Disable by adding `ZC1284` to `disabled_katas` in `.zshellcheckrc`.

---

<a id="zc1285"></a>
### ZC1285 — Use Zsh `${(o)array}` for sorting instead of piping to `sort`

**Severity:** `style`  
**Auto-fix:** `no`

Zsh provides the `(o)` parameter expansion flag to sort array elements in ascending order and `(O)` for descending order. This avoids spawning an external `sort` process for simple array sorting.

Disable by adding `ZC1285` to `disabled_katas` in `.zshellcheckrc`.

---

<a id="zc1286"></a>
### ZC1286 — Use Zsh `${array:#pattern}` instead of `grep -v` for filtering

**Severity:** `style`  
**Auto-fix:** `no`

Zsh provides `${array:#pattern}` to remove matching elements from an array and `${(M)array:#pattern}` to keep only matching elements. This avoids spawning an external `grep` process for simple filtering tasks.

Disable by adding `ZC1286` to `disabled_katas` in `.zshellcheckrc`.

---

<a id="zc1287"></a>
### ZC1287 — Use `cat -v` alternative: Zsh `${(V)var}` for visible control characters

**Severity:** `style`  
**Auto-fix:** `no`

Zsh provides the `(V)` parameter expansion flag to make control characters visible in a variable. This avoids piping through `cat -v` for simple visibility of non-printable characters.

Disable by adding `ZC1287` to `disabled_katas` in `.zshellcheckrc`.

---

<a id="zc1288"></a>
### ZC1288 — Use `typeset` instead of `declare` in Zsh scripts

**Severity:** `style`  
**Auto-fix:** `yes`

`typeset` is the native Zsh builtin for variable declarations. `declare` is a Bash compatibility alias. Using `typeset` is more idiomatic and signals that the script is Zsh-native.

Disable by adding `ZC1288` to `disabled_katas` in `.zshellcheckrc`.

---

<a id="zc1289"></a>
### ZC1289 — Use Zsh `${(u)array}` for unique elements instead of `sort -u`

**Severity:** `style`  
**Auto-fix:** `no`

Zsh provides the `(u)` parameter expansion flag to remove duplicate elements from an array. This preserves original order and avoids spawning an external `sort -u` process.

Disable by adding `ZC1289` to `disabled_katas` in `.zshellcheckrc`.

---

<a id="zc1290"></a>
### ZC1290 — Use Zsh `${(n)array}` for numeric sorting instead of `sort -n`

**Severity:** `style`  
**Auto-fix:** `no`

Zsh provides the `(n)` parameter expansion flag to sort array elements numerically. This avoids spawning an external `sort -n` process for simple numeric sorting of array data.

Disable by adding `ZC1290` to `disabled_katas` in `.zshellcheckrc`.

---

<a id="zc1291"></a>
### ZC1291 — Use Zsh `${(O)array}` for reverse sorting instead of `sort -r`

**Severity:** `style`  
**Auto-fix:** `no`

Zsh provides the `(O)` parameter expansion flag to sort array elements in descending (reverse) order. This avoids spawning an external `sort -r` process for simple reverse sorting of array data.

Disable by adding `ZC1291` to `disabled_katas` in `.zshellcheckrc`.

---

<a id="zc1292"></a>
### ZC1292 — Use Zsh `${var//old/new}` instead of `tr` for character translation

**Severity:** `style`  
**Auto-fix:** `no`

Zsh provides `${var//old/new}` for global substitution within a variable. For simple single-character translation, this avoids spawning `tr` as an external process.

Disable by adding `ZC1292` to `disabled_katas` in `.zshellcheckrc`.

---

<a id="zc1293"></a>
### ZC1293 — Use `[[ ]]` instead of `test` command in Zsh

**Severity:** `style`  
**Auto-fix:** `no`

Zsh `[[ ]]` provides a more powerful conditional expression syntax than the `test` command. It supports pattern matching, regex, and does not require quoting of variable expansions to prevent word splitting.

Disable by adding `ZC1293` to `disabled_katas` in `.zshellcheckrc`.

---

<a id="zc1294"></a>
### ZC1294 — Use `bindkey` instead of `bind` for key bindings in Zsh

**Severity:** `warning`  
**Auto-fix:** `no`

`bind` is a Bash builtin for key bindings. Zsh uses `bindkey` for ZLE (Zsh Line Editor) key bindings. Using `bind` in a Zsh script will fail unless Bash compatibility is loaded.

Disable by adding `ZC1294` to `disabled_katas` in `.zshellcheckrc`.

---

<a id="zc1295"></a>
### ZC1295 — Use `vared` instead of `read -e` for interactive editing in Zsh

**Severity:** `style`  
**Auto-fix:** `no`

Zsh provides `vared` for interactive editing of variables with full ZLE support (tab completion, history, cursor movement). The `read -e` flag is a Bash extension; Zsh `vared` is the native equivalent.

Disable by adding `ZC1295` to `disabled_katas` in `.zshellcheckrc`.

---

<a id="zc1296"></a>
### ZC1296 — Avoid `shopt` in Zsh — use `setopt`/`unsetopt` instead

**Severity:** `warning`  
**Auto-fix:** `no`

`shopt` is a Bash builtin that does not exist in Zsh. Use `setopt` or `unsetopt` to control Zsh shell options. Common Bash `shopt` options have Zsh equivalents via `setopt`.

Disable by adding `ZC1296` to `disabled_katas` in `.zshellcheckrc`.

---

<a id="zc1297"></a>
### ZC1297 — Avoid `$BASH_SOURCE` — use `$0` or `${(%):-%x}` in Zsh

**Severity:** `warning`  
**Auto-fix:** `yes`

`$BASH_SOURCE` is a Bash-specific variable that does not exist in Zsh. In Zsh, use `$0` inside a sourced file to get the script path, or `${(%):-%x}` for the current file regardless of sourcing context.

Disable by adding `ZC1297` to `disabled_katas` in `.zshellcheckrc`.

---

<a id="zc1298"></a>
### ZC1298 — Avoid `$FUNCNAME` — use `$funcstack` in Zsh

**Severity:** `warning`  
**Auto-fix:** `yes`

`$FUNCNAME` is a Bash-specific array that does not exist in Zsh. Zsh provides `$funcstack` as the equivalent, containing the call stack of function names with the current function at index 1.

Disable by adding `ZC1298` to `disabled_katas` in `.zshellcheckrc`.

---

<a id="zc1299"></a>
### ZC1299 — Avoid `$BASH_LINENO` — use `$funcfiletrace` in Zsh

**Severity:** `warning`  
**Auto-fix:** `no`

`$BASH_LINENO` is a Bash-specific array that does not exist in Zsh. Zsh provides `$funcfiletrace` as the equivalent, containing file:line pairs for each call in the function stack.

Disable by adding `ZC1299` to `disabled_katas` in `.zshellcheckrc`.

---

<a id="zc1300"></a>
### ZC1300 — Avoid `$BASH_VERSINFO` — use `$ZSH_VERSION` in Zsh

**Severity:** `warning`  
**Auto-fix:** `yes`

`$BASH_VERSINFO` is a Bash-specific array containing version components. In Zsh, use `$ZSH_VERSION` (string) or `${(s:.:)ZSH_VERSION}` to split it into components for version comparison.

Disable by adding `ZC1300` to `disabled_katas` in `.zshellcheckrc`.

---

<a id="zc1301"></a>
### ZC1301 — Avoid `$PIPESTATUS` — use `$pipestatus` (lowercase) in Zsh

**Severity:** `warning`  
**Auto-fix:** `yes`

`$PIPESTATUS` is a Bash array containing exit statuses from the last pipeline. Zsh uses `$pipestatus` (lowercase) for the same purpose. The uppercase form is undefined in Zsh.

Disable by adding `ZC1301` to `disabled_katas` in `.zshellcheckrc`.

---

<a id="zc1302"></a>
### ZC1302 — Avoid `help` builtin — use `run-help` or `man` in Zsh

**Severity:** `info`  
**Auto-fix:** `no`

The `help` command is a Bash builtin that displays builtin help. Zsh does not have a `help` builtin. Use `run-help <command>` or `man zshbuiltins` for Zsh builtin documentation.

Disable by adding `ZC1302` to `disabled_katas` in `.zshellcheckrc`.

---

<a id="zc1303"></a>
### ZC1303 — Avoid `enable` command — use `zmodload` for Zsh modules

**Severity:** `warning`  
**Auto-fix:** `no`

The `enable` command is a Bash builtin for enabling/disabling builtins. Zsh uses `zmodload` to load and manage modules, and `disable`/`enable` have different semantics. Use `zmodload` for module management.

Disable by adding `ZC1303` to `disabled_katas` in `.zshellcheckrc`.

---

<a id="zc1304"></a>
### ZC1304 — Avoid `$BASH_SUBSHELL` — use `$ZSH_SUBSHELL` in Zsh

**Severity:** `warning`  
**Auto-fix:** `yes`

`$BASH_SUBSHELL` tracks subshell nesting depth in Bash. Zsh provides `$ZSH_SUBSHELL` as the native equivalent.

Disable by adding `ZC1304` to `disabled_katas` in `.zshellcheckrc`.

---

<a id="zc1305"></a>
### ZC1305 — Avoid `$COMP_WORDS` — use `$words` in Zsh completion

**Severity:** `warning`  
**Auto-fix:** `yes`

`$COMP_WORDS` is a Bash completion variable containing the words on the command line. Zsh completion uses `$words` array for the same purpose.

Disable by adding `ZC1305` to `disabled_katas` in `.zshellcheckrc`.

---

<a id="zc1306"></a>
### ZC1306 — Avoid `$COMP_CWORD` — use `$CURRENT` in Zsh completion

**Severity:** `warning`  
**Auto-fix:** `yes`

`$COMP_CWORD` is a Bash completion variable for the current cursor word index. Zsh completion uses `$CURRENT` for the same purpose.

Disable by adding `ZC1306` to `disabled_katas` in `.zshellcheckrc`.

---

<a id="zc1307"></a>
### ZC1307 — Avoid `$DIRSTACK` — use `$dirstack` (lowercase) in Zsh

**Severity:** `warning`  
**Auto-fix:** `yes`

`$DIRSTACK` is the Bash form of the directory stack array. Zsh uses `$dirstack` (lowercase) for the same purpose.

Disable by adding `ZC1307` to `disabled_katas` in `.zshellcheckrc`.

---

<a id="zc1308"></a>
### ZC1308 — Avoid `$COMP_LINE` — use `$BUFFER` in Zsh completion

**Severity:** `warning`  
**Auto-fix:** `yes`

`$COMP_LINE` is a Bash completion variable containing the full command line. Zsh completion uses `$BUFFER` for the current command line content.

Disable by adding `ZC1308` to `disabled_katas` in `.zshellcheckrc`.

---

<a id="zc1309"></a>
### ZC1309 — Avoid `$BASH_COMMAND` — not available in Zsh

**Severity:** `warning`  
**Auto-fix:** `no`

`$BASH_COMMAND` contains the currently executing command in Bash. Zsh does not provide a direct equivalent. Use `$ZSH_DEBUG_CMD` in debug traps or restructure the logic.

Disable by adding `ZC1309` to `disabled_katas` in `.zshellcheckrc`.

---

<a id="zc1310"></a>
### ZC1310 — Avoid `$BASH_EXECUTION_STRING` — not available in Zsh

**Severity:** `warning`  
**Auto-fix:** `no`

`$BASH_EXECUTION_STRING` contains the argument to `bash -c`. Zsh does not provide this variable. Access the script argument directly.

Disable by adding `ZC1310` to `disabled_katas` in `.zshellcheckrc`.

---

<a id="zc1311"></a>
### ZC1311 — Avoid `complete` command — use `compdef` in Zsh

**Severity:** `warning`  
**Auto-fix:** `no`

`complete` is a Bash builtin for registering tab completions. Zsh uses `compdef` for completion registration and the `compctl` legacy interface. Use `compdef` for the modern Zsh completion system.

Disable by adding `ZC1311` to `disabled_katas` in `.zshellcheckrc`.

---

<a id="zc1312"></a>
### ZC1312 — Avoid `compgen` command — use `compadd` in Zsh

**Severity:** `warning`  
**Auto-fix:** `no`

`compgen` is a Bash builtin for generating completions. Zsh uses `compadd` and the completion system functions for adding completion candidates.

Disable by adding `ZC1312` to `disabled_katas` in `.zshellcheckrc`.

---

<a id="zc1313"></a>
### ZC1313 — Avoid `$BASH_ALIASES` — use Zsh `aliases` hash

**Severity:** `warning`  
**Auto-fix:** `yes`

`$BASH_ALIASES` is a Bash associative array of defined aliases. Zsh provides the `aliases` associative array for the same purpose.

Disable by adding `ZC1313` to `disabled_katas` in `.zshellcheckrc`.

---

<a id="zc1314"></a>
### ZC1314 — Avoid `$BASH_LOADABLES_PATH` — not available in Zsh

**Severity:** `warning`  
**Auto-fix:** `no`

`$BASH_LOADABLES_PATH` is a Bash variable for loadable builtin search paths. Zsh has no equivalent; use `zmodload` with full module names instead.

Disable by adding `ZC1314` to `disabled_katas` in `.zshellcheckrc`.

---

<a id="zc1315"></a>
### ZC1315 — Avoid `$BASH_COMPAT` — use `emulate` for compatibility in Zsh

**Severity:** `warning`  
**Auto-fix:** `no`

`$BASH_COMPAT` sets Bash compatibility level. Zsh uses `emulate` to control compatibility mode (e.g., `emulate -L sh` for POSIX mode).

Disable by adding `ZC1315` to `disabled_katas` in `.zshellcheckrc`.

---

<a id="zc1316"></a>
### ZC1316 — Avoid `caller` builtin — use `$funcfiletrace` in Zsh

**Severity:** `warning`  
**Auto-fix:** `no`

`caller` is a Bash builtin that returns the call stack context. Zsh provides `$funcfiletrace`, `$funcstack`, and `$funcsourcetrace` for inspecting the call stack.

Disable by adding `ZC1316` to `disabled_katas` in `.zshellcheckrc`.

---

<a id="zc1317"></a>
### ZC1317 — Avoid `$BASH_ENV` — use `$ZDOTDIR` and `$ENV` in Zsh

**Severity:** `info`  
**Auto-fix:** `no`

`$BASH_ENV` specifies a startup file for non-interactive Bash shells. Zsh uses `$ZDOTDIR` to locate `.zshrc` and related files, and `$ENV` for POSIX-compatible startup.

Disable by adding `ZC1317` to `disabled_katas` in `.zshellcheckrc`.

---

<a id="zc1318"></a>
### ZC1318 — Avoid `$BASH_CMDS` — use `$commands` hash in Zsh

**Severity:** `warning`  
**Auto-fix:** `yes`

`$BASH_CMDS` is a Bash associative array caching command lookups. Zsh provides the `$commands` hash for the same purpose, mapping command names to their full paths.

Disable by adding `ZC1318` to `disabled_katas` in `.zshellcheckrc`.

---

<a id="zc1319"></a>
### ZC1319 — Avoid `$BASH_ARGC` — use `$#` in Zsh

**Severity:** `warning`  
**Auto-fix:** `yes`

`$BASH_ARGC` is a Bash array tracking argument counts per stack frame. Zsh uses `$#` for argument count and `$argv` for the argument array.

Disable by adding `ZC1319` to `disabled_katas` in `.zshellcheckrc`.

---

<a id="zc1320"></a>
### ZC1320 — Avoid `$BASH_ARGV` — use `$argv` in Zsh

**Severity:** `warning`  
**Auto-fix:** `yes`

`$BASH_ARGV` is a Bash array containing arguments in reverse order. Zsh provides `$argv` (or `$@`) for positional parameters.

Disable by adding `ZC1320` to `disabled_katas` in `.zshellcheckrc`.

---

<a id="zc1321"></a>
### ZC1321 — Avoid `$BASH_XTRACEFD` — not available in Zsh

**Severity:** `warning`  
**Auto-fix:** `no`

`$BASH_XTRACEFD` redirects Bash xtrace output to a file descriptor. Zsh does not have this variable. Use `exec 2>file` or redirect stderr directly for trace output redirection.

Disable by adding `ZC1321` to `disabled_katas` in `.zshellcheckrc`.

---

<a id="zc1322"></a>
### ZC1322 — Avoid `$COPROC` — Zsh coproc uses different syntax

**Severity:** `warning`  
**Auto-fix:** `no`

`$COPROC` is a Bash array for coprocess file descriptors. Zsh coprocesses use `coproc` keyword with different variable naming and `read -p`/`print -p` for I/O.

Disable by adding `ZC1322` to `disabled_katas` in `.zshellcheckrc`.

---

<a id="zc1323"></a>
### ZC1323 — Avoid `suspend` builtin — use `kill -STOP $$` in Zsh

**Severity:** `warning`  
**Auto-fix:** `no`

`suspend` is a Bash builtin that suspends the shell. Zsh does not have a `suspend` builtin. Use `kill -STOP $$` or Ctrl-Z for the same effect.

Disable by adding `ZC1323` to `disabled_katas` in `.zshellcheckrc`.

---

<a id="zc1324"></a>
### ZC1324 — Avoid `$PROMPT_COMMAND` — use `precmd` hook in Zsh

**Severity:** `warning`  
**Auto-fix:** `no`

`$PROMPT_COMMAND` is a Bash variable that executes a command before each prompt. Zsh uses the `precmd` hook function for the same purpose.

Disable by adding `ZC1324` to `disabled_katas` in `.zshellcheckrc`.

---

<a id="zc1325"></a>
### ZC1325 — Avoid `$PS0` — use `preexec` hook in Zsh

**Severity:** `warning`  
**Auto-fix:** `no`

`$PS0` is a Bash 4.4+ prompt string displayed before command execution. Zsh uses the `preexec` hook function for running code before each command.

Disable by adding `ZC1325` to `disabled_katas` in `.zshellcheckrc`.

---

<a id="zc1326"></a>
### ZC1326 — Avoid `$HISTTIMEFORMAT` — use `fc -li` in Zsh

**Severity:** `info`  
**Auto-fix:** `no`

`$HISTTIMEFORMAT` is a Bash variable for formatting history timestamps. Zsh stores timestamps automatically when `EXTENDED_HISTORY` is set, and displays them with `fc -li` or `history -i`.

Disable by adding `ZC1326` to `disabled_katas` in `.zshellcheckrc`.

---

<a id="zc1327"></a>
### ZC1327 — Avoid `history -c` — Zsh uses different history management

**Severity:** `warning`  
**Auto-fix:** `no`

`history -c` clears history in Bash. Zsh provides `fc -p` for pushing history to a new file and `fc -P` for popping. Use `fc -W` to write and `fc -R` to read history files.

Disable by adding `ZC1327` to `disabled_katas` in `.zshellcheckrc`.

---

<a id="zc1328"></a>
### ZC1328 — Avoid `$HISTCONTROL` — use Zsh `setopt` history options

**Severity:** `info`  
**Auto-fix:** `no`

`$HISTCONTROL` is a Bash variable controlling history deduplication. Zsh uses `setopt HIST_IGNORE_DUPS`, `HIST_IGNORE_ALL_DUPS`, and `HIST_IGNORE_SPACE` for the same functionality.

Disable by adding `ZC1328` to `disabled_katas` in `.zshellcheckrc`.

---

<a id="zc1329"></a>
### ZC1329 — Avoid `$HISTIGNORE` — use `zshaddhistory` hook in Zsh

**Severity:** `info`  
**Auto-fix:** `no`

`$HISTIGNORE` is a Bash variable for pattern-based history filtering. Zsh uses the `zshaddhistory` hook function and `setopt HIST_IGNORE_SPACE` for controlling which commands enter history.

Disable by adding `ZC1329` to `disabled_katas` in `.zshellcheckrc`.

---

<a id="zc1330"></a>
### ZC1330 — Avoid `$INPUTRC` — use `bindkey` in Zsh

**Severity:** `info`  
**Auto-fix:** `no`

`$INPUTRC` points to the readline configuration file in Bash. Zsh uses `bindkey` and ZLE widgets for key binding configuration, not readline.

Disable by adding `ZC1330` to `disabled_katas` in `.zshellcheckrc`.

---

<a id="zc1331"></a>
### ZC1331 — Avoid `$BASH_REMATCH` — use `$match` array in Zsh

**Severity:** `warning`  
**Auto-fix:** `yes`

`$BASH_REMATCH` holds regex capture groups in Bash. Zsh stores regex matches in the `$match` array (and `$MATCH` for the full match) when using `=~` with `setopt BASH_REMATCH` disabled.

Disable by adding `ZC1331` to `disabled_katas` in `.zshellcheckrc`.

---

<a id="zc1332"></a>
### ZC1332 — Avoid `$GLOBIGNORE` — use `setopt EXTENDED_GLOB` in Zsh

**Severity:** `info`  
**Auto-fix:** `no`

`$GLOBIGNORE` is a Bash variable for excluding patterns from glob expansion. Zsh uses `setopt EXTENDED_GLOB` with the `~` (exclusion) operator or `setopt NULL_GLOB` for different glob behavior.

Disable by adding `ZC1332` to `disabled_katas` in `.zshellcheckrc`.

---

<a id="zc1333"></a>
### ZC1333 — Avoid `$TIMEFORMAT` — use `$TIMEFMT` in Zsh

**Severity:** `info`  
**Auto-fix:** `yes`

`$TIMEFORMAT` is the Bash variable for customizing `time` output. Zsh uses `$TIMEFMT` for the same purpose, with different format specifiers.

Disable by adding `ZC1333` to `disabled_katas` in `.zshellcheckrc`.

---

<a id="zc1334"></a>
### ZC1334 — Avoid `type -p` — use `whence -p` in Zsh

**Severity:** `warning`  
**Auto-fix:** `yes`

`type -p` is a Bash flag that prints the path of a command. Zsh `type` does not support `-p`. Use `whence -p` to get the path of an external command in Zsh.

Disable by adding `ZC1334` to `disabled_katas` in `.zshellcheckrc`.

---

<a id="zc1335"></a>
### ZC1335 — Use Zsh array reversal instead of `tac` for in-memory data

**Severity:** `style`  
**Auto-fix:** `no`

`tac` reverses lines from a file or stdin. For in-memory array data, Zsh provides `${(Oa)array}` to reverse array element order without spawning an external process.

Disable by adding `ZC1335` to `disabled_katas` in `.zshellcheckrc`.

---

<a id="zc1336"></a>
### ZC1336 — Avoid `printenv` — use `typeset -x` or `export` in Zsh

**Severity:** `style`  
**Auto-fix:** `no`

`printenv` is an external command for listing environment variables. Zsh provides `typeset -x` to list exported variables and `export` to display them without spawning a subprocess.

Disable by adding `ZC1336` to `disabled_katas` in `.zshellcheckrc`.

---

<a id="zc1337"></a>
### ZC1337 — Avoid `fold` command — use Zsh `print -l` with `$COLUMNS`

**Severity:** `style`  
**Auto-fix:** `no`

`fold` wraps text to a specified width. Zsh provides `$COLUMNS` for terminal width and `print -l` for line-by-line output, reducing dependency on external commands.

Disable by adding `ZC1337` to `disabled_katas` in `.zshellcheckrc`.

---

<a id="zc1338"></a>
### ZC1338 — Avoid `seq -s` — use Zsh `${(j:sep:)${(s::)...}}` for joining

**Severity:** `style`  
**Auto-fix:** `no`

`seq -s` generates a sequence with a custom separator. Zsh provides native brace expansion with `{start..end}` and `${(j:sep:)array}` for joining, avoiding an external process.

Disable by adding `ZC1338` to `disabled_katas` in `.zshellcheckrc`.

---

<a id="zc1339"></a>
### ZC1339 — Use Zsh `${#${(f)var}}` instead of `wc -l` for line count

**Severity:** `style`  
**Auto-fix:** `no`

Zsh `${(f)var}` splits a string into lines and `${#...}` counts them. Avoid piping through `wc -l` for simple line counting from variables.

Disable by adding `ZC1339` to `disabled_katas` in `.zshellcheckrc`.

---

<a id="zc1340"></a>
### ZC1340 — Avoid `shuf` for random array element — use Zsh `$RANDOM`

**Severity:** `style`  
**Auto-fix:** `no`

Zsh provides `$RANDOM` and array subscripts to pick random elements without spawning `shuf`. For a single random array element, use `${array[RANDOM%$#array+1]}`.

Disable by adding `ZC1340` to `disabled_katas` in `.zshellcheckrc`.

---

<a id="zc1341"></a>
### ZC1341 — Use Zsh `*(.x)` glob qualifier instead of `find -executable`

**Severity:** `style`  
**Auto-fix:** `no`

Zsh's `*(.x)` glob qualifier matches regular files that are executable. Avoid shelling out to `find -executable` when the same selection is one glob away.

Disable by adding `ZC1341` to `disabled_katas` in `.zshellcheckrc`.

---

<a id="zc1342"></a>
### ZC1342 — Use Zsh `*(L0)` glob qualifier instead of `find -empty`

**Severity:** `style`  
**Auto-fix:** `no`

Zsh's `*(L0)` glob qualifier matches files with length 0. Combine with `.` or `/` to restrict to regular files or directories. Avoid shelling out to `find -empty` for the same result.

Disable by adding `ZC1342` to `disabled_katas` in `.zshellcheckrc`.

---

<a id="zc1343"></a>
### ZC1343 — Use Zsh `*(m±N)` glob qualifier instead of `find -mtime N`

**Severity:** `style`  
**Auto-fix:** `no`

Zsh's `*(mN)`, `*(m+N)`, `*(m-N)` glob qualifiers match files by age in days (exact / older / newer). For hours use `*(h±N)`, for minutes `*(M±N)`. Same expressive power as `find -mtime`, no external process.

Disable by adding `ZC1343` to `disabled_katas` in `.zshellcheckrc`.

---

<a id="zc1344"></a>
### ZC1344 — Use Zsh `*(L±Nk)` glob qualifier instead of `find -size`

**Severity:** `style`  
**Auto-fix:** `no`

Zsh's `*(LN)`, `*(L+N)`, `*(L-N)` match files by size in 512-byte blocks (or bytes with a unit suffix: `k`, `m`, `p`). Same expressive power as `find -size` without an external process.

Disable by adding `ZC1344` to `disabled_katas` in `.zshellcheckrc`.

---

<a id="zc1345"></a>
### ZC1345 — Use Zsh `*(f:mode:)` glob qualifier instead of `find -perm`

**Severity:** `style`  
**Auto-fix:** `no`

Zsh's `*(f:mode:)` glob qualifier matches files by permission mode. Use octal (`*(f:0755:)`) or symbolic (`*(f:u+x:)`) inside the colon-delimited form. Avoids spawning `find` for permission filters.

Disable by adding `ZC1345` to `disabled_katas` in `.zshellcheckrc`.

---

<a id="zc1346"></a>
### ZC1346 — Use Zsh `*(u:name:)` glob qualifier instead of `find -user`

**Severity:** `style`  
**Auto-fix:** `no`

Zsh's `*(u:name:)` and `*(u+uid)` glob qualifiers match files by owner (name or numeric uid). The `*(U)` shorthand matches files owned by the current user. Avoid `find -user` for the same selection.

Disable by adding `ZC1346` to `disabled_katas` in `.zshellcheckrc`.

---

<a id="zc1347"></a>
### ZC1347 — Use Zsh `*(g:name:)` glob qualifier instead of `find -group`

**Severity:** `style`  
**Auto-fix:** `no`

Zsh's `*(g:name:)` and `*(g+gid)` glob qualifiers match files by group (name or numeric gid). The `*(G)` shorthand matches files in the current user's group. Avoid `find -group`/`-gid`/`-nogroup` for the same selection.

Disable by adding `ZC1347` to `disabled_katas` in `.zshellcheckrc`.

---

<a id="zc1348"></a>
### ZC1348 — Use Zsh glob type qualifiers instead of `find -type`

**Severity:** `style`  
**Auto-fix:** `no`

Zsh glob qualifiers select node type directly: `*(/)` directories, `*(.)` regular files, `*(@)` symlinks, `*(=)` sockets, `*(p)` named pipes, `*(*)` executable regular files, `*(%)` char/block devices. Avoid `find -type X` for the same selection.

Disable by adding `ZC1348` to `disabled_katas` in `.zshellcheckrc`.

---

<a id="zc1349"></a>
### ZC1349 — Use `${#var}` instead of `expr length "$var"` for string length

**Severity:** `style`  
**Auto-fix:** `no`

Zsh (and POSIX) `${#var}` returns string length without spawning `expr`. Use it wherever you would reach for `expr length` or `expr STRING : '.*'`.

Disable by adding `ZC1349` to `disabled_katas` in `.zshellcheckrc`.

---

<a id="zc1350"></a>
### ZC1350 — Use `${str:pos:len}` instead of `expr substr` for substring extraction

**Severity:** `style`  
**Auto-fix:** `no`

Zsh parameter expansion `${str:pos:len}` extracts a substring starting at `pos` of length `len`. No external `expr` call, and the semantics are consistent with `${str:pos}` (to end) and negative positions.

Disable by adding `ZC1350` to `disabled_katas` in `.zshellcheckrc`.

---

<a id="zc1351"></a>
### ZC1351 — Use `[[ $str =~ pattern ]]` instead of `expr match` / `expr :` for regex

**Severity:** `style`  
**Auto-fix:** `no`

Zsh's `[[ $str =~ pattern ]]` evaluates regex natively and populates `$match` / `$MATCH` / `$mbegin` / `$mend` arrays. Avoid shelling out to `expr match` or the `expr STRING : REGEX` form.

Disable by adding `ZC1351` to `disabled_katas` in `.zshellcheckrc`.

---

<a id="zc1352"></a>
### ZC1352 — Avoid `xargs -I{}` — use a Zsh `for` loop for per-item substitution

**Severity:** `style`  
**Auto-fix:** `no`

`xargs -I{}` runs one command per item with `{}` substituted. A Zsh `for` loop over the same input (`for x in ${(f)"$(cmd)"}`) is clearer and keeps state in the current shell.

Disable by adding `ZC1352` to `disabled_katas` in `.zshellcheckrc`.

---

<a id="zc1353"></a>
### ZC1353 — Avoid `printf -v` — use `print -v` or command substitution in Zsh

**Severity:** `style`  
**Auto-fix:** `no`

`printf -v var fmt ...` is a Bash-ism. In Zsh use `print -v var -rf fmt ...` or plain command substitution `var=$(printf fmt ...)`. `-v` is silently ignored by POSIX printf, producing surprising bugs on portable scripts.

Disable by adding `ZC1353` to `disabled_katas` in `.zshellcheckrc`.

---

<a id="zc1354"></a>
### ZC1354 — Use `whence -w` instead of Bash-specific `type -t` for command classification

**Severity:** `style`  
**Auto-fix:** `no`

`type -t` returns the category (alias, keyword, function, builtin, file) of a command in Bash. Zsh's `whence -w` produces `name: category` output with the same information and without shelling out for the sub-field extraction.

Disable by adding `ZC1354` to `disabled_katas` in `.zshellcheckrc`.

---

<a id="zc1355"></a>
### ZC1355 — Use `print -r` instead of `echo -E` for raw output

**Severity:** `style`  
**Auto-fix:** `yes`

`echo -E` disables backslash interpretation, but the flag is Bash-ism and ignored by POSIX `echo`. Zsh's `print -r` is the idiomatic raw-printer; combine with `-n` (no newline), `-l` (one per line), `-u<fd>` (file descriptor), or `--` (end of flags) as needed.

Disable by adding `ZC1355` to `disabled_katas` in `.zshellcheckrc`.

---

<a id="zc1356"></a>
### ZC1356 — Use `read -A` instead of `read -a` for array read in Zsh

**Severity:** `error`  
**Auto-fix:** `yes`

Zsh's `read` uses `-A` (uppercase A) to read into an array. Bash uses `-a` (lowercase) for the same thing. In Zsh, `read -a` assigns a flag to a scalar variable — not what Bash users expect. Use `-A` for portable-Zsh behavior.

Disable by adding `ZC1356` to `disabled_katas` in `.zshellcheckrc`.

---

<a id="zc1357"></a>
### ZC1357 — Use Zsh `${(q)var}` instead of `printf '%q'` for shell-quoting

**Severity:** `style`  
**Auto-fix:** `no`

Bash's `printf '%q'` emits shell-quoted output. Zsh's `${(q)var}` parameter flag does the same in-shell, with variants `${(qq)var}`, `${(qqq)var}`, `${(qqqq)var}` for single-quote, double-quote, $'...', and POSIX ANSI-C styles respectively.

Disable by adding `ZC1357` to `disabled_katas` in `.zshellcheckrc`.

---

<a id="zc1358"></a>
### ZC1358 — Use `${PWD:P}` instead of `pwd -P` for physical current directory

**Severity:** `style`  
**Auto-fix:** `no`

`pwd -P` resolves symlinks to the physical path. Zsh's `${PWD:P}` modifier does the same without spawning the external — the `P` modifier returns the canonical (absolute, symlink-resolved) form.

Disable by adding `ZC1358` to `disabled_katas` in `.zshellcheckrc`.

---

<a id="zc1359"></a>
### ZC1359 — Avoid `id -Gn` — use Zsh `$groups` associative array

**Severity:** `style`  
**Auto-fix:** `no`

Zsh's `zsh/parameter` module exposes the `$groups` associative array mapping group names to GIDs for the current process. Load with `zmodload zsh/parameter` (often auto-loaded) and inspect `${(k)groups}` for names, avoiding an external `id -Gn`/`groups` call.

Disable by adding `ZC1359` to `disabled_katas` in `.zshellcheckrc`.

---

<a id="zc1360"></a>
### ZC1360 — Use Zsh `*(OL)` glob qualifier instead of `ls -S` for size-ordered listing

**Severity:** `style`  
**Auto-fix:** `no`

Zsh glob qualifier `*(OL)` orders results by size (descending). `*(oL)` is ascending. Combined with `[N]` subscript you get the N-th largest/smallest file without `ls -S` and piping.

Disable by adding `ZC1360` to `disabled_katas` in `.zshellcheckrc`.

---

<a id="zc1361"></a>
### ZC1361 — Avoid `awk 'NR==N'` — use Zsh array subscript on `${(f)...}`

**Severity:** `style`  
**Auto-fix:** `no`

Picking the N-th line with `awk 'NR==N'` spawns awk. Zsh can split file contents on newlines with `${(f)"$(<file)"}` and index directly: `lines=(${(f)"$(<f)"}); print $lines[N]`.

Disable by adding `ZC1361` to `disabled_katas` in `.zshellcheckrc`.

---

<a id="zc1362"></a>
### ZC1362 — Use `[[ -o option ]]` instead of `test -o option` for Zsh option checks

**Severity:** `info`  
**Auto-fix:** `no`

In Zsh, `[[ -o name ]]` tests whether a shell option is set. The `test` / `[` builtin interprets `-o` as a logical OR, not an option-query — so `test -o foo` is a syntax error or wrong behavior. Use the `[[ ... ]]` form for option tests.

Disable by adding `ZC1362` to `disabled_katas` in `.zshellcheckrc`.

---

<a id="zc1363"></a>
### ZC1363 — Use Zsh `*(e:...:)` eval qualifier instead of `find -newer`/`-older`

**Severity:** `style`  
**Auto-fix:** `no`

Zsh's `*(e:expr:)` glob qualifier evaluates an arbitrary expression per match — perfect for `-newer REF`-style predicates. Example: `*(e:'[[ $REPLY -nt reference ]]':)` selects files newer than `reference`.

Disable by adding `ZC1363` to `disabled_katas` in `.zshellcheckrc`.

---

<a id="zc1364"></a>
### ZC1364 — Use Zsh `${var:pos:len}` instead of `cut -c` for character ranges

**Severity:** `style`  
**Auto-fix:** `no`

`cut -c N-M` extracts characters N through M from each line. Zsh's `${var:pos:len}` (0-indexed position, length) does the same from a variable without spawning `cut`.

Disable by adding `ZC1364` to `disabled_katas` in `.zshellcheckrc`.

---

<a id="zc1365"></a>
### ZC1365 — Use Zsh `zstat` module instead of `stat -c` for file metadata

**Severity:** `style`  
**Auto-fix:** `no`

Zsh's `zsh/stat` module (loaded with `zmodload zsh/stat` — the command is named `zstat`) exposes every `stat(2)` field natively: mtime, size, owner, group, mode, links, etc. Avoid external `stat -c '%...'` invocations.

Disable by adding `ZC1365` to `disabled_katas` in `.zshellcheckrc`.

---

<a id="zc1366"></a>
### ZC1366 — Use Zsh `limit` instead of POSIX `ulimit` for idiomatic resource queries

**Severity:** `style`  
**Auto-fix:** `no`

Zsh provides both `ulimit` (POSIX compatibility) and `limit` (Zsh native). `limit` prints human-readable values (`cputime 10 seconds` vs `-t 10`) and accepts `unlimited` as a value. Prefer `limit` for Zsh-idiomatic scripts; keep `ulimit` only when the script must run under Bash as well.

Disable by adding `ZC1366` to `disabled_katas` in `.zshellcheckrc`.

---

<a id="zc1367"></a>
### ZC1367 — Use Zsh `strftime` instead of Bash `printf '%(fmt)T'`

**Severity:** `style`  
**Auto-fix:** `no`

Bash 4.2+ supports `printf '%(fmt)T\n' seconds` to format a timestamp. Zsh's `zsh/datetime` module provides `strftime` which is more readable and works consistently across versions: `strftime '%Y-%m-%d' $EPOCHSECONDS`.

Disable by adding `ZC1367` to `disabled_katas` in `.zshellcheckrc`.

---

<a id="zc1368"></a>
### ZC1368 — Avoid `sh -c` / `bash -c` inside a Zsh script — inline or use a function

**Severity:** `style`  
**Auto-fix:** `no`

Invoking `sh -c` or `bash -c` inside a Zsh script spawns a second shell, loses access to the parent script's functions, arrays, and associative arrays, and re-interprets POSIX-only syntax. Inline the code as a function or use `zsh -c` when a subshell is truly required.

Disable by adding `ZC1368` to `disabled_katas` in `.zshellcheckrc`.

---

<a id="zc1369"></a>
### ZC1369 — Prefer Zsh `${(V)var}` over `od -c` for printable-visible character output

**Severity:** `style`  
**Auto-fix:** `no`

Zsh's `${(V)var}` parameter flag renders non-printable characters in visible form (e.g. `\n` for newline). For simple inspection of a variable's contents, this avoids the `od -c` process entirely.

Disable by adding `ZC1369` to `disabled_katas` in `.zshellcheckrc`.

---

<a id="zc1370"></a>
### ZC1370 — Prefer Zsh `repeat N { ... }` over `yes str | head -n N` for finite output

**Severity:** `style`  
**Auto-fix:** `no`

`yes` plus `head` is a common idiom for producing N copies of a line. Zsh's `repeat N { print str }` does the same loop in-shell without spawning yes or the pipe, and without the SIGPIPE handshake.

Disable by adding `ZC1370` to `disabled_katas` in `.zshellcheckrc`.

---

<a id="zc1371"></a>
### ZC1371 — Use Zsh array `:t` modifier instead of `basename -a` for bulk path stripping

**Severity:** `style`  
**Auto-fix:** `no`

`basename -a a b c` returns the file name component of each path. Zsh's `${array:t}` parameter modifier applies the same tail-component extraction to every element of an array at once — no external process.

Disable by adding `ZC1371` to `disabled_katas` in `.zshellcheckrc`.

---

<a id="zc1372"></a>
### ZC1372 — Use Zsh `zmv` autoload function instead of `rename`/`rename.ul`

**Severity:** `style`  
**Auto-fix:** `no`

Zsh's `zmv` (autoloaded via `autoload -Uz zmv`) batch-renames files using glob patterns with capture groups. Safer than the various `rename`/`rename.ul`/`prename` utilities (perl-based vs util-linux) and does not depend on which one is installed.

Disable by adding `ZC1372` to `disabled_katas` in `.zshellcheckrc`.

---

<a id="zc1373"></a>
### ZC1373 — Use Zsh `${(0)var}` flag for NUL-split parsing instead of `env -0`

**Severity:** `style`  
**Auto-fix:** `no`

When reading NUL-terminated data (e.g. `/proc/*/environ`), Zsh's `${(0)var}` parameter flag splits on NUL into an array natively. Avoid `env -0 | xargs -0 ...` chains that require two additional processes.

Disable by adding `ZC1373` to `disabled_katas` in `.zshellcheckrc`.

---

<a id="zc1374"></a>
### ZC1374 — Avoid `$FUNCNEST` — Zsh uses `$FUNCNEST` as a limit, not a depth indicator

**Severity:** `warning`  
**Auto-fix:** `yes`

Bash's `$FUNCNEST` is both a writable limit and (implicitly) the current depth-query vehicle. Zsh's `$FUNCNEST` is only the limit — to read the current depth use `${#funcstack}`. Reading `$FUNCNEST` expecting depth returns the limit, not the current depth.

Disable by adding `ZC1374` to `disabled_katas` in `.zshellcheckrc`.

---

<a id="zc1375"></a>
### ZC1375 — Use `[[ -t fd ]]` instead of `tty -s` for tty-check

**Severity:** `style`  
**Auto-fix:** `no`

`tty -s` exits 0 if stdin is a terminal. Zsh's `[[ -t 0 ]]` (or `[[ -t 1 ]]` for stdout, `[[ -t 2 ]]` for stderr) does the same check without spawning `tty`.

Disable by adding `ZC1375` to `disabled_katas` in `.zshellcheckrc`.

---

<a id="zc1376"></a>
### ZC1376 — Avoid `BASH_XTRACEFD` — use Zsh `exec {fd}>file` + `setopt XTRACE`

**Severity:** `warning`  
**Auto-fix:** `no`

Bash's `BASH_XTRACEFD` redirects `set -x` output to a file descriptor. Zsh does not honor this variable; setting it is a silent no-op. To redirect trace output in Zsh, open a dedicated fd with `exec {fd}>file` and redirect fd 2 through it: `exec 2>&$fd; setopt XTRACE`.

Disable by adding `ZC1376` to `disabled_katas` in `.zshellcheckrc`.

---

<a id="zc1377"></a>
### ZC1377 — Avoid `$BASH_ALIASES` — use Zsh `$aliases` associative array

**Severity:** `warning`  
**Auto-fix:** `yes`

Bash's `$BASH_ALIASES` is an associative array of alias→value mappings. Zsh exposes the same information via `$aliases` (also an assoc array). `$BASH_ALIASES` is unset in Zsh; reading it yields nothing.

Disable by adding `ZC1377` to `disabled_katas` in `.zshellcheckrc`.

---

<a id="zc1378"></a>
### ZC1378 — Avoid uppercase `$DIRSTACK` — Zsh uses lowercase `$dirstack`

**Severity:** `error`  
**Auto-fix:** `yes`

Bash's `$DIRSTACK` is the `pushd`/`popd` directory stack. Zsh exposes the same stack as lowercase `$dirstack` (per zsh/parameter module). Using uppercase `$DIRSTACK` in Zsh accesses an unrelated (and usually empty) variable.

Disable by adding `ZC1378` to `disabled_katas` in `.zshellcheckrc`.

---

<a id="zc1379"></a>
### ZC1379 — Avoid `$PROMPT_COMMAND` — use Zsh `precmd` function

**Severity:** `warning`  
**Auto-fix:** `no`

Bash runs the command in `$PROMPT_COMMAND` before each prompt. Zsh does not honor this variable; the equivalent is a function named `precmd` (or registered via `add-zsh-hook precmd name`). Reading `$PROMPT_COMMAND` in Zsh is a no-op.

Disable by adding `ZC1379` to `disabled_katas` in `.zshellcheckrc`.

---

<a id="zc1380"></a>
### ZC1380 — Avoid `$HISTIGNORE` — use Zsh `$HISTORY_IGNORE`

**Severity:** `warning`  
**Auto-fix:** `yes`

Bash filters history entries matching `$HISTIGNORE` patterns. Zsh uses a parameter named `$HISTORY_IGNORE` (underscore in the middle). Setting `HISTIGNORE` in Zsh is a no-op.

Disable by adding `ZC1380` to `disabled_katas` in `.zshellcheckrc`.

---

<a id="zc1381"></a>
### ZC1381 — Avoid `$COMP_WORDS`/`$COMP_CWORD` — Zsh uses `words`/`$CURRENT`

**Severity:** `error`  
**Auto-fix:** `yes`

Bash programmable completion reads the partial command via `$COMP_WORDS` (array of tokens) and `$COMP_CWORD` (index of cursor). Zsh's completion system exposes the same via `words` (array) and `$CURRENT` (1-based cursor index). Using the Bash names in Zsh completion functions produces empty expansions.

Disable by adding `ZC1381` to `disabled_katas` in `.zshellcheckrc`.

---

<a id="zc1382"></a>
### ZC1382 — Avoid `$READLINE_LINE`/`$READLINE_POINT` — Zsh ZLE uses `$BUFFER`/`$CURSOR`

**Severity:** `error`  
**Auto-fix:** `yes`

Bash readline exposes the current input line as `$READLINE_LINE` and cursor offset as `$READLINE_POINT` inside `bind -x` handlers. Zsh's Line Editor (ZLE) uses `$BUFFER` (line text) and `$CURSOR` (1-based column) inside widget functions. The Bash names are unset in Zsh.

Disable by adding `ZC1382` to `disabled_katas` in `.zshellcheckrc`.

---

<a id="zc1383"></a>
### ZC1383 — Avoid `$TIMEFORMAT` — Zsh uses `$TIMEFMT`

**Severity:** `warning`  
**Auto-fix:** `yes`

Bash's `$TIMEFORMAT` controls the output of the `time` builtin. Zsh uses a shorter name, `$TIMEFMT`, for the same purpose. Setting `TIMEFORMAT` in a Zsh script has no effect; the Zsh `time` builtin reads `$TIMEFMT`.

Disable by adding `ZC1383` to `disabled_katas` in `.zshellcheckrc`.

---

<a id="zc1384"></a>
### ZC1384 — Avoid `$EXECIGNORE` — Bash-only; Zsh uses completion-system ignore patterns

**Severity:** `warning`  
**Auto-fix:** `no`

Bash's `$EXECIGNORE` excludes matching commands from PATH hashing. Zsh does not honor this variable; use the compsys tag-based filters (`zstyle ':completion:*' ignored-patterns ...`) for a similar effect on completion.

Disable by adding `ZC1384` to `disabled_katas` in `.zshellcheckrc`.

---

<a id="zc1385"></a>
### ZC1385 — Avoid `$PS0` — Bash-only; Zsh uses `preexec` hook

**Severity:** `warning`  
**Auto-fix:** `no`

Bash 4.4+ prints `$PS0` after reading a command and before executing it. Zsh does not honor `$PS0`; the equivalent is a `preexec` function (or `add-zsh-hook preexec funcname`) which receives the command line as `$1`.

Disable by adding `ZC1385` to `disabled_katas` in `.zshellcheckrc`.

---

<a id="zc1386"></a>
### ZC1386 — Avoid `$FIGNORE` — Bash-only; Zsh uses compsys tag patterns

**Severity:** `warning`  
**Auto-fix:** `no`

Bash's `$FIGNORE` hides filenames matching listed suffixes from completion. Zsh does not honor this variable; use `zstyle ':completion:*' ignored-patterns '*.o *.pyc'` or the file-patterns tag for equivalent filtering.

Disable by adding `ZC1386` to `disabled_katas` in `.zshellcheckrc`.

---

<a id="zc1387"></a>
### ZC1387 — Avoid `$SHELLOPTS` — Zsh uses `$options` associative array

**Severity:** `warning`  
**Auto-fix:** `no`

Bash's `$SHELLOPTS` is a colon-separated list of set options. Zsh exposes the same information via the `$options` associative array (keys are option names, values are `on`/`off`). `$SHELLOPTS` is unset in Zsh.

Disable by adding `ZC1387` to `disabled_katas` in `.zshellcheckrc`.

---

<a id="zc1388"></a>
### ZC1388 — Use Zsh lowercase `$mailpath` array instead of colon-separated `$MAILPATH`

**Severity:** `warning`  
**Auto-fix:** `no`

Bash uses `$MAILPATH` — a colon-separated string of mail files with optional `?message` suffixes. Zsh uses lowercase `$mailpath` as an array (each element: `file?message`), which is typed and parseable. Setting the uppercase name in Zsh is ignored.

Disable by adding `ZC1388` to `disabled_katas` in `.zshellcheckrc`.

---

<a id="zc1389"></a>
### ZC1389 — Avoid `$HOSTFILE` — Bash-only; Zsh uses `$hosts` array

**Severity:** `warning`  
**Auto-fix:** `no`

Bash reads `$HOSTFILE` to feed hostname completion. Zsh populates hostname completion from the `$hosts` array (lowercase). Setting `$HOSTFILE` in Zsh is ignored; extend `$hosts` instead.

Disable by adding `ZC1389` to `disabled_katas` in `.zshellcheckrc`.

---

<a id="zc1390"></a>
### ZC1390 — Avoid `$GROUPS[@]` — Zsh `$GROUPS` is a scalar, not an array

**Severity:** `error`  
**Auto-fix:** `no`

Bash's `$GROUPS` is an array of all group IDs the user belongs to, so `${GROUPS[@]}` iterates them. In Zsh, `$GROUPS` is a scalar (primary GID). The array of all group IDs is `$(groups)` output or `${(k)groups}` (if the `zsh/parameter` module is loaded, `$groups` is an assoc array name→gid).

Disable by adding `ZC1390` to `disabled_katas` in `.zshellcheckrc`.

---

<a id="zc1391"></a>
### ZC1391 — Avoid `[[ -v VAR ]]` for Bash set-check — use Zsh `(( ${+VAR} ))`

**Severity:** `warning`  
**Auto-fix:** `no`

Bash 4.2+ supports `[[ -v VAR ]]` to test whether a variable is set. Zsh `[[ -v VAR ]]` is parsed but not as the set-check — Zsh's canonical form is `(( ${+VAR} ))` which evaluates to 1 when set and 0 when unset, working reliably across Zsh versions.

Disable by adding `ZC1391` to `disabled_katas` in `.zshellcheckrc`.

---

<a id="zc1392"></a>
### ZC1392 — Avoid `$CHILD_MAX` — Bash-only; Zsh uses `limit` / `ulimit -u`

**Severity:** `info`  
**Auto-fix:** `no`

Bash's `$CHILD_MAX` reports the maximum number of exited child processes Bash remembers. Zsh does not export this var. For current process limits use `limit -s maxproc` or `ulimit -u` — but the exact Bash semantic is not mirrored.

Disable by adding `ZC1392` to `disabled_katas` in `.zshellcheckrc`.

---

<a id="zc1393"></a>
### ZC1393 — Avoid `$SRANDOM` — Bash 5.1+ only, read `/dev/urandom` in Zsh

**Severity:** `warning`  
**Auto-fix:** `no`

Bash 5.1 added `$SRANDOM` as a cryptographically secure 32-bit random value. Zsh does not have an equivalent variable. For secure random integers, read bytes from `/dev/urandom` (e.g. `(( n = 0x$(od -N4 -An -tx1 /dev/urandom | tr -d ' ') ))`) or use an external such as `openssl rand`.

Disable by adding `ZC1393` to `disabled_katas` in `.zshellcheckrc`.

---

<a id="zc1394"></a>
### ZC1394 — Avoid `$BASH` — Zsh uses `$ZSH_NAME` for the interpreter name

**Severity:** `info`  
**Auto-fix:** `yes`

Bash's `$BASH` holds the path to the running Bash executable. Zsh's equivalent is `$ZSH_NAME` (for the binary name) or `$0` (interactive shell). Using `$BASH` in a Zsh script yields empty output.

Disable by adding `ZC1394` to `disabled_katas` in `.zshellcheckrc`.

---

<a id="zc1395"></a>
### ZC1395 — Avoid `wait -n` — Bash 4.3+ only; Zsh `wait` on job IDs

**Severity:** `warning`  
**Auto-fix:** `no`

Bash 4.3+ added `wait -n` (wait for any job to finish). Zsh's `wait` does not accept `-n`; instead wait explicitly on job IDs or PIDs, or use `wait` with no args (waits for all). For any-of semantics use `wait $pid1 $pid2; ...` in a loop.

Disable by adding `ZC1395` to `disabled_katas` in `.zshellcheckrc`.

---

<a id="zc1396"></a>
### ZC1396 — Avoid `unset -n` — Bash nameref semantics not in Zsh

**Severity:** `error`  
**Auto-fix:** `no`

Bash's `unset -n NAME` unsets the nameref itself rather than the target variable it points to. Zsh does not implement namerefs; `unset -n` flags as an error or unsets something unintended. Use `unset -v` for variable unset and `unset -f` for function unset explicitly.

Disable by adding `ZC1396` to `disabled_katas` in `.zshellcheckrc`.

---

<a id="zc1397"></a>
### ZC1397 — Avoid `$COMP_TYPE`/`$COMP_KEY` — Bash completion globals, not in Zsh

**Severity:** `error`  
**Auto-fix:** `no`

Bash programmable completion exposes `$COMP_TYPE` (completion type) and `$COMP_KEY` (completion key pressed). Zsh's compsys does not use these variables; query completion context via `$compstate` assoc array or context keys from `_arguments`/`_values` instead.

Disable by adding `ZC1397` to `disabled_katas` in `.zshellcheckrc`.

---

<a id="zc1398"></a>
### ZC1398 — Avoid `$PROMPT_DIRTRIM` — use Zsh `%N~` prompt modifier

**Severity:** `warning`  
**Auto-fix:** `no`

Bash's `$PROMPT_DIRTRIM` limits the number of directory components shown in `\w`. Zsh has no such variable; use the `%N~` prompt escape (N is component count) or `%/` / `%~` with precmd adjustments for Zsh-native directory truncation.

Disable by adding `ZC1398` to `disabled_katas` in `.zshellcheckrc`.

---

<a id="zc1399"></a>
### ZC1399 — Use Zsh `$signals` array instead of `kill -l` for signal enumeration

**Severity:** `style`  
**Auto-fix:** `no`

Zsh exposes the `$signals` array (from `zsh/parameter`) holding all signal names indexed from 0. `print -l $signals` produces the same list as `kill -l` without spawning an external process.

Disable by adding `ZC1399` to `disabled_katas` in `.zshellcheckrc`.

---

<a id="zc1400"></a>
### ZC1400 — Use Zsh `$CPUTYPE` for architecture detection instead of parsing `$HOSTTYPE`

**Severity:** `info`  
**Auto-fix:** `no`

Bash's `$HOSTTYPE` is a combined architecture/vendor/OS string (e.g. `x86_64-pc-linux-gnu`). Zsh exposes the same as `$HOSTTYPE` but additionally splits out `$CPUTYPE` (e.g. `x86_64`) for pure architecture queries — no `awk -F-` needed to extract.

Disable by adding `ZC1400` to `disabled_katas` in `.zshellcheckrc`.

---

<a id="zc1401"></a>
### ZC1401 — Prefer Zsh `$VENDOR` over parsing `$MACHTYPE` for vendor detection

**Severity:** `info`  
**Auto-fix:** `no`

Both Bash and Zsh expose `$MACHTYPE` (e.g. `x86_64-pc-linux-gnu`). Zsh additionally pre-parses the vendor component into `$VENDOR` (e.g. `pc`, `apple`). Avoid `cut -d- -f2 <<< $MACHTYPE` when `$VENDOR` is available directly.

Disable by adding `ZC1401` to `disabled_katas` in `.zshellcheckrc`.

---

<a id="zc1402"></a>
### ZC1402 — Avoid `date -d @seconds` — use Zsh `strftime` for epoch formatting

**Severity:** `style`  
**Auto-fix:** `no`

`date -d @N -- '+fmt'` / `date --date=@N` converts epoch seconds to a formatted date. Zsh's `zsh/datetime` module provides `strftime fmt N` directly — a single builtin, no `date` spawn, and the `-d`/`@` form is GNU-specific (not portable to BSD `date`).

Disable by adding `ZC1402` to `disabled_katas` in `.zshellcheckrc`.

---

<a id="zc1403"></a>
### ZC1403 — Setting `$HISTFILESIZE` alone is incomplete in Zsh — pair with `$SAVEHIST`

**Severity:** `warning`  
**Auto-fix:** `yes`

Bash uses `$HISTSIZE` (in-memory) and `$HISTFILESIZE` (on disk). Zsh uses `$HISTSIZE` (in-memory) and `$SAVEHIST` (on disk). Setting only `$HISTFILESIZE` in Zsh has no effect on disk — `$SAVEHIST` must be set. Mixing both names leaves disk-history behavior undefined.

Disable by adding `ZC1403` to `disabled_katas` in `.zshellcheckrc`.

---

<a id="zc1404"></a>
### ZC1404 — Avoid `$BASH_CMDS` — Bash-specific hash-table mirror, use Zsh `$commands`

**Severity:** `warning`  
**Auto-fix:** `yes`

Bash's `$BASH_CMDS` associative array mirrors the hash-table of command names→paths. Zsh exposes the same via `$commands` (assoc array from `zsh/parameter`). `$BASH_CMDS` is unset in Zsh.

Disable by adding `ZC1404` to `disabled_katas` in `.zshellcheckrc`.

---

<a id="zc1405"></a>
### ZC1405 — Avoid `env -u VAR cmd` — use Zsh `(unset VAR; cmd)` subshell

**Severity:** `style`  
**Auto-fix:** `no`

`env -u VAR cmd` unsets a variable for a single command. In Zsh the idiomatic form is a subshell: `(unset VAR; cmd)` — no external `env` spawn, and the unset is naturally scoped to the subshell.

Disable by adding `ZC1405` to `disabled_katas` in `.zshellcheckrc`.

---

<a id="zc1406"></a>
### ZC1406 — Prefer Zsh `zargs -P N` autoload over `xargs -P N` for parallel execution

**Severity:** `style`  
**Auto-fix:** `no`

Zsh provides `zargs` (loaded via `autoload -Uz zargs`) — a native equivalent of `xargs` with parallel execution via `-P`. It keeps variables and functions in scope (unlike xargs) and avoids the utility-quoting surprises of `xargs`.

Disable by adding `ZC1406` to `disabled_katas` in `.zshellcheckrc`.

---

<a id="zc1407"></a>
### ZC1407 — Avoid `/dev/tcp/...` — use Zsh `zsh/net/tcp` module

**Severity:** `error`  
**Auto-fix:** `no`

`/dev/tcp/host/port` is a Bash-specific virtual-file interface for TCP connections; Zsh does not implement it. For TCP in Zsh, load `zmodload zsh/net/tcp` and use `ztcp host port` which exposes the connection as a regular file descriptor.

Disable by adding `ZC1407` to `disabled_katas` in `.zshellcheckrc`.

---

<a id="zc1408"></a>
### ZC1408 — Avoid `$BASH_FUNC_...%%` — Bash-specific exported-function envvar

**Severity:** `error`  
**Auto-fix:** `no`

Bash exports functions into environment variables named `BASH_FUNC_NAME%%`. These are consumed only by other Bash shells. Zsh does not recognize the format and will neither inherit the function nor clean these envvars.

Disable by adding `ZC1408` to `disabled_katas` in `.zshellcheckrc`.

---

<a id="zc1409"></a>
### ZC1409 — Avoid `[ -N file ]` / `test -N file` — Bash-only, use Zsh `zstat` for mtime comparison

**Severity:** `info`  
**Auto-fix:** `no`

`[ -N file ]` and `test -N file` test whether a file has been modified since last read (Bash extension). Zsh does not implement `-N`. Use the `zsh/stat` module to compare `atime` and `mtime` explicitly: `zstat -H s file; (( s[mtime] > s[atime] ))`.

Disable by adding `ZC1409` to `disabled_katas` in `.zshellcheckrc`.

---

<a id="zc1410"></a>
### ZC1410 — Avoid `compopt` — Bash programmable-completion modifier, not in Zsh

**Severity:** `error`  
**Auto-fix:** `no`

`compopt` tweaks Bash programmable-completion options for the current completion. Zsh's compsys does not implement `compopt`; completion options are set via `zstyle` / completion-function context instead.

Disable by adding `ZC1410` to `disabled_katas` in `.zshellcheckrc`.

---

<a id="zc1411"></a>
### ZC1411 — Use Zsh `disable` instead of Bash `enable -n` to hide builtins

**Severity:** `style`  
**Auto-fix:** `yes`

Bash's `enable -n name` disables a builtin so that the external of the same name is used. Zsh provides a dedicated `disable` builtin: `disable name` achieves the same in one verb. Re-enable later with `enable name`.

Disable by adding `ZC1411` to `disabled_katas` in `.zshellcheckrc`.

---

<a id="zc1412"></a>
### ZC1412 — Avoid `$COMPREPLY` — Bash completion output, use Zsh `compadd`

**Severity:** `error`  
**Auto-fix:** `no`

Bash completion functions populate the `$COMPREPLY` array to declare candidates. Zsh's compsys uses the `compadd` builtin: `compadd -- foo bar baz`. Setting `$COMPREPLY` in a Zsh completion does nothing.

Disable by adding `ZC1412` to `disabled_katas` in `.zshellcheckrc`.

---

<a id="zc1413"></a>
### ZC1413 — Use Zsh `whence -p cmd` instead of `hash -t cmd` for resolved path

**Severity:** `style`  
**Auto-fix:** `yes`

Bash's `hash -t cmd` prints the hashed path for `cmd` (or fails if not hashed). Zsh's `whence -p cmd` prints the PATH-resolved absolute path, whether hashed or not — more reliable and the native Zsh idiom.

Disable by adding `ZC1413` to `disabled_katas` in `.zshellcheckrc`.

---

<a id="zc1414"></a>
### ZC1414 — Beware `hash -d` — Bash deletes from hash table, Zsh defines named directory

**Severity:** `error`  
**Auto-fix:** `no`

The `-d` flag has opposite meanings across shells: Bash `hash -d NAME` removes `NAME` from the command-hash table. Zsh `hash -d NAME=PATH` **defines** a named directory (`~NAME` expansion). A Bash script ported to Zsh breaks silently when `hash -d ls` is interpreted as defining `~ls`.

Disable by adding `ZC1414` to `disabled_katas` in `.zshellcheckrc`.

---

<a id="zc1415"></a>
### ZC1415 — Prefer Zsh `TRAPZERR` function over `trap 'cmd' ERR`

**Severity:** `info`  
**Auto-fix:** `no`

Both Bash and Zsh accept `trap 'cmd' ERR`, but Zsh's idiomatic form is the named function `TRAPZERR`: `TRAPZERR() { echo "err at $LINENO"; }`. The named function receives `$1` = signal and is easier to compose than an inline string.

Disable by adding `ZC1415` to `disabled_katas` in `.zshellcheckrc`.

---

<a id="zc1416"></a>
### ZC1416 — Prefer Zsh `preexec` hook over `trap 'cmd' DEBUG`

**Severity:** `warning`  
**Auto-fix:** `no`

Bash's `trap 'cmd' DEBUG` runs `cmd` before each simple command. Zsh's equivalent is the `preexec` function (or `add-zsh-hook preexec name`) which receives the about-to-execute command line as `$1`, `$2`, `$3`. The DEBUG trap is not fired in Zsh the way it is in Bash — use preexec for portability.

Disable by adding `ZC1416` to `disabled_katas` in `.zshellcheckrc`.

---

<a id="zc1417"></a>
### ZC1417 — Prefer Zsh `TRAPRETURN` function over `trap 'cmd' RETURN`

**Severity:** `info`  
**Auto-fix:** `no`

Bash's `trap 'cmd' RETURN` runs `cmd` when a function returns. Zsh accepts the `RETURN` signal name but the idiomatic form is a function named `TRAPRETURN`: `TRAPRETURN() { print "returning $?"; }`.

Disable by adding `ZC1417` to `disabled_katas` in `.zshellcheckrc`.

---

<a id="zc1418"></a>
### ZC1418 — Use Zsh `limit -h`/`-s` instead of `ulimit -H`/`-S` for hard/soft limits

**Severity:** `style`  
**Auto-fix:** `no`

Bash's `ulimit` uses uppercase `-H` (hard) and `-S` (soft). Zsh's native `limit` builtin uses lowercase `-h` and `-s` for the same. The Zsh form is easier to remember and produces human-readable output.

Disable by adding `ZC1418` to `disabled_katas` in `.zshellcheckrc`.

---

<a id="zc1419"></a>
### ZC1419 — Avoid `chmod 777` — grants world-writable access

**Severity:** `warning`  
**Auto-fix:** `no`

Mode 777 (or 0777) grants read/write/execute to owner, group, and world. Files become world-writable, which on a multi-user system or inside a container with mapped UIDs is almost always wrong. Use 755 for executables, 644 for regular files, 700 for private directories, or `umask`-aware helpers.

Disable by adding `ZC1419` to `disabled_katas` in `.zshellcheckrc`.

---

<a id="zc1420"></a>
### ZC1420 — Avoid `chmod +s` / `chmod u+s` — setuid/setgid is a security risk

**Severity:** `warning`  
**Auto-fix:** `no`

Setuid (mode bit 4000) and setgid (2000) cause the program to run with the file-owner's (or group's) privileges, not the caller's. Any bug in such a program is a privilege-escalation vector. Reserve setuid for audited, minimal binaries; prefer sudo + policy, capabilities, or containers for less-trusted tooling.

Disable by adding `ZC1420` to `disabled_katas` in `.zshellcheckrc`.

---

<a id="zc1421"></a>
### ZC1421 — Avoid `chpasswd` / `passwd --stdin` — plaintext passwords in process tree

**Severity:** `error`  
**Auto-fix:** `no`

Passing passwords on stdin to `chpasswd` or `passwd --stdin` exposes the plaintext in the process command line or pipeline — visible to `ps`, logs, and environment. Use encrypted-hash input (`chpasswd -e`), `usermod -p` with a hash, or an IaC tool that handles credentials outside the process tree.

Disable by adding `ZC1421` to `disabled_katas` in `.zshellcheckrc`.

---

<a id="zc1422"></a>
### ZC1422 — Avoid `sudo -S` — reads password from stdin, exposes plaintext

**Severity:** `error`  
**Auto-fix:** `no`

`sudo -S` reads the password from stdin, enabling `echo $PW | sudo -S cmd` patterns that place the plaintext password in the process tree and shell history. Prefer `sudo -A` with a graphical askpass, `NOPASSWD:` in sudoers for specific commands, or `pkexec` for policy-based privilege elevation.

Disable by adding `ZC1422` to `disabled_katas` in `.zshellcheckrc`.

---

<a id="zc1423"></a>
### ZC1423 — Dangerous: `iptables -F` / `nft flush ruleset` — drops all firewall rules

**Severity:** `warning`  
**Auto-fix:** `no`

Flushing the firewall ruleset removes every existing rule, typically reverting to the default policy. On a remote machine with policy=DROP, this locks you out. Save existing rules first (`iptables-save > backup`) and consider `iptables-apply` with a rollback timer.

Disable by adding `ZC1423` to `disabled_katas` in `.zshellcheckrc`.

---

<a id="zc1424"></a>
### ZC1424 — Dangerous: `mkfs.*` / `mkfs -t` — formats a filesystem, destroys data

**Severity:** `error`  
**Auto-fix:** `no`

`mkfs.ext4 /dev/sda1`, `mkfs.xfs /dev/...`, `mkfs -t ...` all destroy the existing filesystem on the target device. A typo on the target path reformats the wrong disk. Validate the device path, use `blkid` / `lsblk` first, and consider a confirmation prompt.

Disable by adding `ZC1424` to `disabled_katas` in `.zshellcheckrc`.

---

<a id="zc1425"></a>
### ZC1425 — `shutdown` / `reboot` / `halt` / `poweroff` — confirm before scripting

**Severity:** `warning`  
**Auto-fix:** `no`

Scripts that invoke `shutdown`, `reboot`, `halt`, `poweroff`, or `systemctl poweroff` take down the system. Unattended invocation in automation is often wrong (e.g. leftover test step). Prefer `systemctl isolate rescue.target` for controlled scenarios, and require explicit confirmation for interactive scripts.

Disable by adding `ZC1425` to `disabled_katas` in `.zshellcheckrc`.

---

<a id="zc1426"></a>
### ZC1426 — Avoid `git clone http://` — unencrypted transport, use `https://` or `git://`+verify

**Severity:** `warning`  
**Auto-fix:** `no`

`git clone http://...` transfers repository content unencrypted and unauthenticated — susceptible to MITM insertion of malicious commits. Use `https://` for authenticated hosts (GitHub, GitLab) or SSH (`git@host:path`) with verified host keys. Plain `http://` has no integrity guarantee.

Disable by adding `ZC1426` to `disabled_katas` in `.zshellcheckrc`.

---

<a id="zc1427"></a>
### ZC1427 — Dangerous: `nc -e` / `ncat -e` — spawns arbitrary command on network connect

**Severity:** `error`  
**Auto-fix:** `no`

`nc -e cmd` and `ncat --exec cmd` pipe the network socket to an arbitrary command. Incoming connections get a shell or any command you specify — the classic reverse-shell pattern. Many distros ship `nc` compiled without `-e` for this reason. Remove `-e` from scripts except in audited, restricted contexts.

Disable by adding `ZC1427` to `disabled_katas` in `.zshellcheckrc`.

---

<a id="zc1428"></a>
### ZC1428 — Avoid `curl -u user:pass` — credentials visible in process list

**Severity:** `error`  
**Auto-fix:** `no`

`curl -u user:password` places the credentials in the command line, where they show up in `ps`, `/proc/*/cmdline`, shell history, and most audit logs. Use `-u user:` with an interactive password prompt, `--netrc`/`--netrc-file` for persistent credentials, or a credentials manager.

Disable by adding `ZC1428` to `disabled_katas` in `.zshellcheckrc`.

---

<a id="zc1429"></a>
### ZC1429 — Avoid `umount -f` / `-l` — force/lazy unmount masks real issues

**Severity:** `warning`  
**Auto-fix:** `no`

`umount -f` forces the unmount even if the FS is busy; `-l` (lazy) detaches immediately but keeps the FS in-use. Both can leave stale file handles and data loss. Fix the underlying 'target busy' (use `lsof` / `fuser -m` to find users) instead of forcing.

Disable by adding `ZC1429` to `disabled_katas` in `.zshellcheckrc`.

---

<a id="zc1430"></a>
### ZC1430 — Prefer Zsh `zsh/sched` module over `at now` / `batch` for in-shell scheduling

**Severity:** `style`  
**Auto-fix:** `no`

`at`/`batch` schedule commands via the atd daemon — requires daemon running, leaves a spool-file audit trail, and runs in a fresh environment. For in-shell scheduling the Zsh `zsh/sched` module (`sched +1:00 cmd`) runs the command from the current shell without the daemon dependency.

Disable by adding `ZC1430` to `disabled_katas` in `.zshellcheckrc`.

---

<a id="zc1431"></a>
### ZC1431 — Dangerous: `crontab -r` — removes all the user's cron jobs without confirmation

**Severity:** `warning`  
**Auto-fix:** `no`

`crontab -r` deletes the entire crontab for the current user (or the target user with `-u`). There is no `.bak` left behind, no `-i` prompt by default on most platforms. Back up first with `crontab -l > /tmp/cron.bak`, then use `crontab -ir` (interactive) to require confirmation.

Disable by adding `ZC1431` to `disabled_katas` in `.zshellcheckrc`.

---

<a id="zc1432"></a>
### ZC1432 — Dangerous: `passwd -d user` — deletes the password, leaving the account passwordless

**Severity:** `error`  
**Auto-fix:** `no`

`passwd -d user` removes the password entirely, making the account usable without any password (depending on PAM config). This is almost never what you want — use `passwd -l user` to lock the account, or `usermod -L` + delete the ssh keys to fully disable login.

Disable by adding `ZC1432` to `disabled_katas` in `.zshellcheckrc`.

---

<a id="zc1433"></a>
### ZC1433 — Caution with `userdel -f` / `-r` — removes home directory and kills processes

**Severity:** `warning`  
**Auto-fix:** `no`

`userdel -f` proceeds even when the user is logged in or has running processes, potentially killing unsaved work. `-r` additionally deletes the home directory and mail spool. Combined (`-rf`) these are destructive and often misused for 'clean up a user' without warning. Verify no active sessions first.

Disable by adding `ZC1433` to `disabled_katas` in `.zshellcheckrc`.

---

<a id="zc1434"></a>
### ZC1434 — Warn on `swapoff -a` — disables all swap, can OOM-kill

**Severity:** `warning`  
**Auto-fix:** `no`

`swapoff -a` disables every active swap. On a memory-constrained host this pushes data back into RAM, potentially triggering OOM-killer. Prefer disabling specific devices/files (`swapoff /swapfile`) and verify memory headroom with `free -m` first.

Disable by adding `ZC1434` to `disabled_katas` in `.zshellcheckrc`.

---

<a id="zc1435"></a>
### ZC1435 — Avoid `killall -9` / `killall -KILL` — force-kill by process name

**Severity:** `warning`  
**Auto-fix:** `no`

`killall -9 name` sends SIGKILL to every process matching `name` — in multi-user or containerized environments, this can hit unrelated processes that happen to share the name. Prefer `killall -TERM` first (graceful), or kill by PID after locating with `pgrep` / `pidof`.

Disable by adding `ZC1435` to `disabled_katas` in `.zshellcheckrc`.

---

<a id="zc1436"></a>
### ZC1436 — `sysctl -w` is ephemeral — persist in `/etc/sysctl.d/*.conf` for surviving reboots

**Severity:** `info`  
**Auto-fix:** `no`

`sysctl -w key=value` sets a kernel parameter until the next reboot. For configuration that must survive reboots, write a file in `/etc/sysctl.d/` and apply with `sysctl --system`. Using only `-w` in provisioning scripts creates silent drift.

Disable by adding `ZC1436` to `disabled_katas` in `.zshellcheckrc`.

---

<a id="zc1437"></a>
### ZC1437 — `dmesg -c` / `-C` clears the kernel ring buffer — destroys evidence

**Severity:** `warning`  
**Auto-fix:** `no`

`dmesg -c` prints the ring buffer and then **clears** it. `dmesg -C` clears without printing. Any later debugging loses the earlier messages. Prefer plain `dmesg` for read-only inspection, or `journalctl -k` with a time filter.

Disable by adding `ZC1437` to `disabled_katas` in `.zshellcheckrc`.

---

<a id="zc1438"></a>
### ZC1438 — `systemctl mask` permanently prevents service start — document the unmask path

**Severity:** `warning`  
**Auto-fix:** `no`

`systemctl mask unit` symlinks the unit to `/dev/null`, preventing any start (manual, dependency, or at boot). Even `systemctl start` fails with 'Unit is masked.'. The reverse `systemctl unmask` is easy to forget. Document the unmask in provisioning scripts or use `disable` (which still allows manual start).

Disable by adding `ZC1438` to `disabled_katas` in `.zshellcheckrc`.

---

<a id="zc1439"></a>
### ZC1439 — Enabling IP forwarding in a script — document firewall posture

**Severity:** `warning`  
**Auto-fix:** `no`

Setting `net.ipv4.ip_forward=1` (or `-w`-ing a sysctl to the same effect) turns the host into a router. Without matching iptables/nftables rules this can silently expose services between interfaces. If intentional (VPN, container host, NAT gateway), pair with explicit firewall rules and persist via `/etc/sysctl.d/`.

Disable by adding `ZC1439` to `disabled_katas` in `.zshellcheckrc`.

---

<a id="zc1440"></a>
### ZC1440 — `usermod -G group user` replaces supplementary groups — use `-aG` to append

**Severity:** `warning`  
**Auto-fix:** `no`

`usermod -G group user` overwrites the user's supplementary group list — any prior group memberships are removed. Users commonly add themselves to `docker` or `wheel` via `-G` and inadvertently lose `sudo`/`audio`/other memberships. Always pair with `-a` (`-aG`) to append instead of replace.

Disable by adding `ZC1440` to `disabled_katas` in `.zshellcheckrc`.

---

<a id="zc1441"></a>
### ZC1441 — Warn on `docker system prune -af` / `-a --force` (or similar podman/k8s)

**Severity:** `warning`  
**Auto-fix:** `no`

`docker system prune -af` deletes every unused image, container, network, and (with `--volumes`) volume. On shared CI runners or build hosts this obliterates cached layers and slows future builds. Scope prunes with `--filter "until=168h"` or target one resource type at a time.

Disable by adding `ZC1441` to `disabled_katas` in `.zshellcheckrc`.

---

<a id="zc1442"></a>
### ZC1442 — Dangerous: `kubectl delete --all` / `--all-namespaces` deletes cluster resources

**Severity:** `error`  
**Auto-fix:** `no`

`kubectl delete --all pods` (in the current namespace) or `-A`/`--all-namespaces` scopes delete operations across the whole cluster. A typo on the resource type can wipe deployments, services, secrets, or even CRDs. Always use `--dry-run=client` first, then apply with `-n` explicit namespace.

Disable by adding `ZC1442` to `disabled_katas` in `.zshellcheckrc`.

---

<a id="zc1443"></a>
### ZC1443 — Dangerous: `terraform destroy` / `apply -destroy` without `-target`

**Severity:** `warning`  
**Auto-fix:** `no`

`terraform destroy` (or `terraform apply -destroy`) without a `-target` removes every resource in state — entire environments, databases, volumes, DNS, everything. Always prefer targeted destroy or scope via workspaces. Consider guarding state-destroying commands behind an interactive confirmation.

Disable by adding `ZC1443` to `disabled_katas` in `.zshellcheckrc`.

---

<a id="zc1444"></a>
### ZC1444 — Dangerous: `redis-cli FLUSHALL` / `FLUSHDB` — wipes Redis data

**Severity:** `error`  
**Auto-fix:** `no`

`FLUSHALL` deletes every key in every database; `FLUSHDB` clears the current DB. Running against production is usually catastrophic. Either rename the command in `redis.conf` (`rename-command FLUSHALL ""`) or require an explicit confirmation in scripts.

Disable by adding `ZC1444` to `disabled_katas` in `.zshellcheckrc`.

---

<a id="zc1445"></a>
### ZC1445 — Dangerous: `dropdb` / `mysqladmin drop` — deletes a database

**Severity:** `error`  
**Auto-fix:** `no`

`dropdb NAME` removes a PostgreSQL database including all data and schemas. `mysqladmin drop NAME` does the same for MySQL. Always `pg_dump` / `mysqldump` first and consider requiring `-i`/`-y`-less forms so operators must type confirmation.

Disable by adding `ZC1445` to `disabled_katas` in `.zshellcheckrc`.

---

<a id="zc1446"></a>
### ZC1446 — Dangerous: `aws s3 rm --recursive` / `s3 rb --force` — bulk S3 deletion

**Severity:** `error`  
**Auto-fix:** `no`

`aws s3 rm s3://bucket/prefix --recursive` deletes every key under the prefix. `aws s3 rb --force` deletes the bucket along with its contents. Combine with a wrong prefix or bucket name and data loss is total. Enable versioning on production buckets and use `aws s3api list-object-versions` before bulk removals.

Disable by adding `ZC1446` to `disabled_katas` in `.zshellcheckrc`.

---

<a id="zc1447"></a>
### ZC1447 — Avoid deprecated `ifconfig` / `netstat` — prefer `ip` / `ss`

**Severity:** `style`  
**Auto-fix:** `no`

On modern Linux, `ifconfig` and `netstat` (from net-tools) are deprecated in favor of the iproute2 suite: `ip addr`, `ip link`, `ip route`, `ss`. net-tools is not installed by default on many distros (Alpine, Fedora Cloud, minimal images), so scripts break. Use iproute2 commands for portability.

Disable by adding `ZC1447` to `disabled_katas` in `.zshellcheckrc`.

---

<a id="zc1448"></a>
### ZC1448 — `apt-get install` / `apt install` without `-y` hangs in non-interactive scripts

**Severity:** `warning`  
**Auto-fix:** `yes`

In provisioning scripts, `apt-get install foo` (no `-y`) waits for interactive confirmation and stalls CI/Dockerfiles indefinitely. Always pass `-y` (or `--yes`), and for unattended upgrades also set `DEBIAN_FRONTEND=noninteractive` in the environment.

Disable by adding `ZC1448` to `disabled_katas` in `.zshellcheckrc`.

---

<a id="zc1449"></a>
### ZC1449 — `dnf`/`yum` install without `-y` hangs in non-interactive scripts

**Severity:** `warning`  
**Auto-fix:** `no`

In CI/Dockerfiles, `dnf install pkg` or `yum install pkg` prompts for confirmation and stalls. Always pass `-y` (or `--assumeyes`) for unattended runs. Also consider `--nodocs` and `--setopt=install_weak_deps=False` for slim images.

Disable by adding `ZC1449` to `disabled_katas` in `.zshellcheckrc`.

---

<a id="zc1450"></a>
### ZC1450 — `pacman -S` / `zypper install` without non-interactive flag hangs in scripts

**Severity:** `warning`  
**Auto-fix:** `no`

Arch's `pacman -S` waits on confirmation unless `--noconfirm` is passed. SUSE's `zypper install` needs `--non-interactive` (or `-n`). Both stall CI pipelines and Dockerfiles without these flags.

Disable by adding `ZC1450` to `disabled_katas` in `.zshellcheckrc`.

---

<a id="zc1451"></a>
### ZC1451 — Avoid `pip install` without `--user` or virtualenv

**Severity:** `warning`  
**Auto-fix:** `no`

`pip install pkg` (no `--user`, no active venv) targets the system Python, potentially breaking system tools or requiring sudo. On modern Linux this now fails with PEP 668 `externally-managed-environment`. Always use a virtualenv (`python -m venv`, `uv`, `poetry`) or `--user` for scoped installs.

Disable by adding `ZC1451` to `disabled_katas` in `.zshellcheckrc`.

---

<a id="zc1452"></a>
### ZC1452 — Avoid `npm install -g` — global installs need root, break under multiple Node versions

**Severity:** `style`  
**Auto-fix:** `no`

`npm install -g` places packages in a system-wide prefix (typically `/usr/local`). That requires sudo, conflicts with Node version managers (nvm, asdf, volta), and is rarely what you want in a project. Prefer project-local installs (`npm i`), or `pnpm dlx`/`npx` for one-off tools.

Disable by adding `ZC1452` to `disabled_katas` in `.zshellcheckrc`.

---

<a id="zc1453"></a>
### ZC1453 — Avoid `sudo pip` / `sudo npm` / `sudo gem` — language package managers as root

**Severity:** `warning`  
**Auto-fix:** `no`

Running a language package manager as root installs third-party code with full privileges, may overwrite distro-managed libs, and can execute arbitrary install-time hooks as root. Use `--user`, a virtualenv/venv, or a version manager (nvm, pyenv, rbenv) instead.

Disable by adding `ZC1453` to `disabled_katas` in `.zshellcheckrc`.

---

<a id="zc1454"></a>
### ZC1454 — Avoid `docker/podman run --privileged` — disables most container isolation

**Severity:** `error`  
**Auto-fix:** `no`

`--privileged` disables the seccomp profile, grants all Linux capabilities, and lets the container access all host devices. It is effectively equivalent to running the process as host root. Add specific capabilities with `--cap-add` and bind-mount specific devices with `--device` instead.

Disable by adding `ZC1454` to `disabled_katas` in `.zshellcheckrc`.

---

<a id="zc1455"></a>
### ZC1455 — Avoid `docker run --net=host` / `--network=host` — disables network isolation

**Severity:** `warning`  
**Auto-fix:** `no`

Host networking gives the container direct access to the host's network stack, including localhost services. A vulnerable container can reach services meant to be local-only. Use `-p hostport:containerport` for specific publishes and dedicated networks for inter-container traffic.

Disable by adding `ZC1455` to `disabled_katas` in `.zshellcheckrc`.

---

<a id="zc1456"></a>
### ZC1456 — Avoid `docker run -v /:...` — bind-mounts host root into container

**Severity:** `error`  
**Auto-fix:** `no`

Mounting `/` (host root) into a container gives the container read/write access to the entire host filesystem — a trivial container escape. Mount only the specific host paths the container needs, using `:ro` for read-only where possible.

Disable by adding `ZC1456` to `disabled_katas` in `.zshellcheckrc`.

---

<a id="zc1457"></a>
### ZC1457 — Warn on bind-mount of `/var/run/docker.sock` — container escape vector

**Severity:** `warning`  
**Auto-fix:** `no`

Mounting `/var/run/docker.sock` into a container lets the container start any privileged container, mount host filesystems, and effectively gain root on the host. Reserve this for trusted CI/tooling images; for general workloads use rootless containers or a dedicated orchestrator API.

Disable by adding `ZC1457` to `disabled_katas` in `.zshellcheckrc`.

---

<a id="zc1458"></a>
### ZC1458 — Warn on explicit `docker run --user root` / `--user 0`

**Severity:** `warning`  
**Auto-fix:** `no`

Running as UID 0 inside a container means a break-out bug leaves the attacker as root on the host (absent user namespaces). Build images with a non-root `USER` directive and avoid overriding to root at runtime.

Disable by adding `ZC1458` to `disabled_katas` in `.zshellcheckrc`.

---

<a id="zc1459"></a>
### ZC1459 — Warn on `docker run --cap-add=SYS_ADMIN` / other dangerous capabilities

**Severity:** `warning`  
**Auto-fix:** `no`

Granting `SYS_ADMIN`, `SYS_PTRACE`, `SYS_MODULE`, `NET_ADMIN`, or `ALL` capabilities effectively disables the container's security boundary — most container escapes rely on exactly these. Drop all capabilities and add back only the specific ones the workload needs (usually none).

Disable by adding `ZC1459` to `disabled_katas` in `.zshellcheckrc`.

---

<a id="zc1460"></a>
### ZC1460 — Warn on `docker run --security-opt seccomp=unconfined` / `apparmor=unconfined`

**Severity:** `warning`  
**Auto-fix:** `no`

Disabling seccomp or AppArmor removes the syscall / MAC filter that blocks most container escape exploits. Only disable these in a known-safe development context; production workloads should keep the default profile or ship a stricter custom profile.

Disable by adding `ZC1460` to `disabled_katas` in `.zshellcheckrc`.

---

<a id="zc1461"></a>
### ZC1461 — Avoid `docker run --pid=host` — shares host PID namespace with the container

**Severity:** `warning`  
**Auto-fix:** `no`

`--pid=host` lets the container see every host process and send signals to them, including sending SIGKILL to init-managed daemons or attaching a debugger to host-side processes. Use only for diagnostic tools (e.g. strace/perf containers) and never for general workloads.

Disable by adding `ZC1461` to `disabled_katas` in `.zshellcheckrc`.

---

<a id="zc1462"></a>
### ZC1462 — Avoid `docker run --ipc=host` — shares host IPC namespace (/dev/shm, SysV IPC)

**Severity:** `warning`  
**Auto-fix:** `no`

`--ipc=host` makes the container share `/dev/shm` and the SysV IPC keyspace with the host. Any process on the host can read/write the container's shared memory (and vice-versa), making side-channel and data-theft attacks trivial. Use the default private IPC namespace unless two containers explicitly need to share IPC.

Disable by adding `ZC1462` to `disabled_katas` in `.zshellcheckrc`.

---

<a id="zc1463"></a>
### ZC1463 — Avoid `docker run --userns=host` — disables user-namespace remapping

**Severity:** `warning`  
**Auto-fix:** `no`

`--userns=host` turns off the user-namespace remap, meaning UID 0 in the container maps to UID 0 on the host. Combined with any of the `--cap-add`, `--privileged`, or bind-mount footguns, this becomes a direct host-root escalation. Leave the default (container-side remap) enabled.

Disable by adding `ZC1463` to `disabled_katas` in `.zshellcheckrc`.

---

<a id="zc1464"></a>
### ZC1464 — Warn on `iptables -F` / `-P INPUT ACCEPT` — flushes or opens the host firewall

**Severity:** `warning`  
**Auto-fix:** `no`

Flushing all rules (`-F`) or setting the default INPUT/FORWARD policy to ACCEPT leaves the host with no network filter. This is rarely correct outside a first-boot provisioning script, and is a frequent post-compromise persistence step. Use `iptables-save`/`iptables-restore` for atomic reloads and keep a default-drop policy on all hook chains.

Disable by adding `ZC1464` to `disabled_katas` in `.zshellcheckrc`.

---

<a id="zc1465"></a>
### ZC1465 — Warn on `setenforce 0` — disables SELinux enforcement

**Severity:** `warning`  
**Auto-fix:** `no`

`setenforce 0` switches SELinux to permissive mode, silencing every policy decision into an audit log line instead of a deny. It is the textbook post-compromise persistence step and also a common "fix" that papers over an actual policy bug. Address the specific AVC with `audit2allow` instead, and leave `setenforce 1` (enforcing) in production.

Disable by adding `ZC1465` to `disabled_katas` in `.zshellcheckrc`.

---

<a id="zc1466"></a>
### ZC1466 — Warn on disabling the host firewall (`ufw disable` / `systemctl stop firewalld`)

**Severity:** `warning`  
**Auto-fix:** `no`

Disabling the host firewall leaves every listening port reachable from every network the host is on. This is a common "just make it work" shortcut that has shipped to production more than once. Keep the firewall running and open the specific port with `ufw allow <port>` / `firewall-cmd --add-port=<port>/tcp`.

Disable by adding `ZC1466` to `disabled_katas` in `.zshellcheckrc`.

---

<a id="zc1467"></a>
### ZC1467 — Warn on `sysctl -w kernel.core_pattern=|...` / `kernel.modprobe=...` (kernel hijack)

**Severity:** `error`  
**Auto-fix:** `no`

Writing `kernel.core_pattern` to a pipe handler or `kernel.modprobe` to a user-writable path is a textbook privilege-escalation trick: the next crashing setuid process (or the next auto-load of an absent module) executes the supplied binary as root. Keep `core_pattern` set to `core` or `systemd-coredump` and leave `kernel.modprobe` at the distro default (`/sbin/modprobe`).

Disable by adding `ZC1467` to `disabled_katas` in `.zshellcheckrc`.

---

<a id="zc1468"></a>
### ZC1468 — Error on apt `--allow-unauthenticated` / `--force-yes` — installs unsigned packages

**Severity:** `error`  
**Auto-fix:** `no`

`--allow-unauthenticated` and the deprecated `--force-yes` disable APT's package-signature verification, turning any MITM or typo-squat into arbitrary code execution as root. Always sign internal packages and leave verification on.

Disable by adding `ZC1468` to `disabled_katas` in `.zshellcheckrc`.

---

<a id="zc1469"></a>
### ZC1469 — Error on `dnf/yum --nogpgcheck` or `rpm --nosignature` (unsigned RPM install)

**Severity:** `error`  
**Auto-fix:** `no`

`--nogpgcheck` / `--nosignature` / `--nodigest` disable RPM package signature and digest verification. This turns every mirror, cache, or MITM into a direct root compromise. Always keep GPG/signature checking on; sign internal repos with your own key.

Disable by adding `ZC1469` to `disabled_katas` in `.zshellcheckrc`.

---

<a id="zc1470"></a>
### ZC1470 — Error on `git config http.sslVerify false` / `git -c http.sslVerify=false`

**Severity:** `error`  
**Auto-fix:** `no`

Disabling `http.sslVerify` in git means every subsequent fetch / clone accepts any TLS certificate — MITM trivially replaces the tree you are cloning with attacker-controlled code. Fix the broken CA instead: install the certificate, point at the right store with `GIT_SSL_CAINFO`, or use an SSH transport.

Disable by adding `ZC1470` to `disabled_katas` in `.zshellcheckrc`.

---

<a id="zc1471"></a>
### ZC1471 — Error on `kubectl/helm --insecure-skip-tls-verify` (cluster MITM)

**Severity:** `error`  
**Auto-fix:** `no`

`--insecure-skip-tls-verify` tells kubectl / helm to accept any certificate from the API server. Against a production cluster, this hands every secret and admission payload to a MITM. Fix the trust chain: point `--certificate-authority` at the right CA bundle, or restore `KUBECONFIG` with the cluster's embedded CA.

Disable by adding `ZC1471` to `disabled_katas` in `.zshellcheckrc`.

---

<a id="zc1472"></a>
### ZC1472 — Error on `aws s3 --acl public-read` / `public-read-write` (public bucket)

**Severity:** `error`  
**Auto-fix:** `no`

Using the `public-read` or `public-read-write` canned ACL when uploading, syncing, or setting a bucket policy makes the object (and often the bucket) readable by anyone on the internet. Prefer bucket policies scoped to specific principals, or CloudFront with Origin Access Identity if you truly need public read.

Disable by adding `ZC1472` to `disabled_katas` in `.zshellcheckrc`.

---

<a id="zc1473"></a>
### ZC1473 — Warn on `openssl req ... -nodes` / `genrsa` without passphrase — unencrypted private key

**Severity:** `warning`  
**Auto-fix:** `no`

`-nodes` tells OpenSSL not to encrypt the private key that is written to disk. The file ends up at whatever filesystem permissions the umask dictates, and any subsequent backup / container image / rsync picks up a usable key with no passphrase. Use `-aes256` / `-aes-256-cbc` and keep the passphrase in a secrets store, or rely on a hardware-backed key via PKCS#11 / TPM.

Disable by adding `ZC1473` to `disabled_katas` in `.zshellcheckrc`.

---

<a id="zc1474"></a>
### ZC1474 — Warn on `ssh-keygen -N ""` — generates passwordless SSH key

**Severity:** `warning`  
**Auto-fix:** `no`

Generating an SSH key with an empty passphrase (`-N ""`) leaves the key usable by anything that can read the file. Combined with a weak umask or a backup that follows the file, this is a common lateral-movement vector. Use a real passphrase, or delegate key storage to `ssh-agent` / a hardware token.

Disable by adding `ZC1474` to `disabled_katas` in `.zshellcheckrc`.

---

<a id="zc1475"></a>
### ZC1475 — Warn on `setcap` granting dangerous capabilities to a binary (privesc)

**Severity:** `warning`  
**Auto-fix:** `no`

Adding CAP_SYS_ADMIN, CAP_DAC_OVERRIDE, CAP_DAC_READ_SEARCH, CAP_SYS_PTRACE, or CAP_SETUID to a binary lets any user who can execute it perform operations roughly equivalent to root — read any file, change any UID, attach ptrace to root processes. Scope the capability as narrowly as possible (e.g. CAP_NET_BIND_SERVICE) or run the binary under a dedicated service user with a systemd unit.

Disable by adding `ZC1475` to `disabled_katas` in `.zshellcheckrc`.

---

<a id="zc1476"></a>
### ZC1476 — Warn on `apt-key add` — deprecated, trusts every repo system-wide

**Severity:** `warning`  
**Auto-fix:** `no`

`apt-key` was deprecated in APT 2.2 and removed from `apt` 2.5. Keys added with `apt-key add` end up in a global keyring that signs every repo on the system, so a typo-squatted third-party PPA can ship updates for `apt`, `libc6`, or `openssh-server`. Store the key in `/etc/apt/keyrings/<vendor>.gpg` and scope it in `signed-by=` on the specific `sources.list.d` entry.

Disable by adding `ZC1476` to `disabled_katas` in `.zshellcheckrc`.

---

<a id="zc1477"></a>
### ZC1477 — Warn on `printf "$var"` — variable in format-string position (printf-fmt attack)

**Severity:** `warning`  
**Auto-fix:** `no`

The first argument to `printf` is a format string. Interpolating a shell variable into it means any `%` sequence inside the variable is interpreted as a format specifier — at best producing garbage, at worst crashing with `%s`-out-of-bounds reads or writing attacker-controlled data with `%n`. Always use a literal format string: `printf '%s\n' "$var"`.

Disable by adding `ZC1477` to `disabled_katas` in `.zshellcheckrc`.

---

<a id="zc1478"></a>
### ZC1478 — Avoid `mktemp -u` — returns a name without creating the file (TOCTOU)

**Severity:** `warning`  
**Auto-fix:** `no`

`mktemp -u` allocates a unique name but does not create the file, leaving a classic time-of-check to time-of-use race: a second process (possibly attacker- controlled on a multi-user host or shared CI runner) can claim the name before you redirect into it. Drop `-u` and operate on the file `mktemp` creates for you, or use `mktemp -d` if you need a directory path.

Disable by adding `ZC1478` to `disabled_katas` in `.zshellcheckrc`.

---

<a id="zc1479"></a>
### ZC1479 — Error on `ssh/scp -o StrictHostKeyChecking=no` / `UserKnownHostsFile=/dev/null`

**Severity:** `error`  
**Auto-fix:** `no`

Setting `StrictHostKeyChecking=no` or pointing `UserKnownHostsFile` at `/dev/null` makes the client accept any server key on the first (and every) connection, stripping the protection against MITM that SSH is designed to provide. For ephemeral CI targets, pin the host key in `known_hosts` with `ssh-keyscan` and verify the fingerprint out of band, or use `StrictHostKeyChecking=accept-new` at most.

Disable by adding `ZC1479` to `disabled_katas` in `.zshellcheckrc`.

---

<a id="zc1480"></a>
### ZC1480 — Warn on `terraform apply -auto-approve` / `destroy -auto-approve` in scripts

**Severity:** `warning`  
**Auto-fix:** `no`

Running `terraform apply -auto-approve` or `destroy -auto-approve` from a shell script skips the plan-review step that exists to catch schema drift, accidental `-replace`, and resources being deleted. Fine for throwaway CI against a PR environment, but dangerous against shared state. Prefer running `plan` + `apply` with an out-file and human approval, or scope the auto-apply to specific branches/environments.

Disable by adding `ZC1480` to `disabled_katas` in `.zshellcheckrc`.

---

<a id="zc1481"></a>
### ZC1481 — Warn on `unset HISTFILE` / `export HISTFILE=/dev/null` — disables shell history

**Severity:** `warning`  
**Auto-fix:** `no`

Disabling shell history (`unset HISTFILE`, `HISTFILE=/dev/null`, `HISTSIZE=0`) is a classic stepping stone for hiding post-compromise activity. Legitimate scripts almost never need this — if you are pasting a secret on the command line, use `HISTCONTROL=ignorespace` and prefix the line with a space, or read the value from a file / stdin.

Disable by adding `ZC1481` to `disabled_katas` in `.zshellcheckrc`.

---

<a id="zc1482"></a>
### ZC1482 — Error on `docker login -p` / `--password=` — credential in process list

**Severity:** `error`  
**Auto-fix:** `no`

Passing the registry password on the command line puts it in the output of `ps`, `/proc/<pid>/cmdline`, and the shell history. On a shared CI runner or a host with unprivileged users, that is an immediate leak. Use `--password-stdin` and pipe the secret in from `cat /run/secrets/foo` or a credential helper.

Disable by adding `ZC1482` to `disabled_katas` in `.zshellcheckrc`.

---

<a id="zc1483"></a>
### ZC1483 — Warn on `pip install --break-system-packages` — bypasses PEP 668 externally-managed guard

**Severity:** `warning`  
**Auto-fix:** `no`

`--break-system-packages` tells pip to ignore the distro's PEP 668 marker and install into `/usr/lib/python*`, overwriting files the package manager owns. The next `apt`/`dnf` upgrade clobbers or gets clobbered by the pip-installed version, and you now have two sources of truth for Python dependencies. Install into a virtualenv (`python -m venv`), use `pipx` for application scripts, or use `uv` / `poetry` for project dependencies.

Disable by adding `ZC1483` to `disabled_katas` in `.zshellcheckrc`.

---

<a id="zc1484"></a>
### ZC1484 — Error on `npm/yarn/pnpm config set strict-ssl false` — disables registry TLS verification

**Severity:** `error`  
**Auto-fix:** `no`

Turning off `strict-ssl` for npm, yarn, or pnpm makes the client accept any TLS certificate from the registry — a MITM (corporate proxy, compromised WiFi, rogue BGP) can substitute any package, including new versions of `react` or `lodash`. If the registry uses a private CA, point `cafile` / `NODE_EXTRA_CA_CERTS` at the right bundle instead.

Disable by adding `ZC1484` to `disabled_katas` in `.zshellcheckrc`.

---

<a id="zc1485"></a>
### ZC1485 — Warn on `openssl s_client -ssl3 / -tls1 / -tls1_1` — legacy TLS

**Severity:** `warning`  
**Auto-fix:** `no`

Forcing SSLv3, TLSv1.0, or TLSv1.1 connects with protocols that have known downgrade and bit-flip attacks (POODLE, BEAST). These are disabled by default in every maintained OpenSSL build. If the remote only speaks an old protocol, the right fix is to update the remote, not downgrade your client.

Disable by adding `ZC1485` to `disabled_katas` in `.zshellcheckrc`.

---

<a id="zc1486"></a>
### ZC1486 — Warn on `curl -2` / `-3` — forces broken SSLv2 / SSLv3

**Severity:** `warning`  
**Auto-fix:** `no`

`curl -2` (SSLv2) and `-3` (SSLv3) force protocols that are removed from every current TLS library. `-2` matches no working server; `-3` leaves you open to POODLE. If the remote really needs an old protocol the fix is on the server, not the client.

Disable by adding `ZC1486` to `disabled_katas` in `.zshellcheckrc`.

---

<a id="zc1487"></a>
### ZC1487 — Warn on `history -c` — clears shell history (and is a Bash-ism under Zsh)

**Severity:** `warning`  
**Auto-fix:** `no`

`history -c` clears the in-memory history buffer in Bash. It is a standard post-compromise anti-forensics step. It is also a Bash-ism: in Zsh, `history` takes completely different arguments, so a copy-pasted `history -c` silently no-ops and leaves the author thinking history was cleared when it was not. If you really need to rotate history in a Zsh script, unset `HISTFILE` before the sensitive block or redirect to `/dev/null` explicitly.

Disable by adding `ZC1487` to `disabled_katas` in `.zshellcheckrc`.

---

<a id="zc1488"></a>
### ZC1488 — Warn on `ssh -R 0.0.0.0:...` / `*:...` — reverse tunnel bound to all interfaces

**Severity:** `warning`  
**Auto-fix:** `no`

The default for `ssh -R` binds the remote listener to `localhost`. Pointing it at `0.0.0.0` or `*` (or an explicit public IP) exposes the forwarded port to the whole network, including anything else that has reached the jump host. For persistent ops tunnels, pin the bind address to a specific private interface and require `GatewayPorts clientspecified` server-side.

Disable by adding `ZC1488` to `disabled_katas` in `.zshellcheckrc`.

---

<a id="zc1489"></a>
### ZC1489 — Error on `nc -e` / `ncat -e` — classic reverse-shell invocation

**Severity:** `error`  
**Auto-fix:** `no`

`nc -e <shell>` and `ncat -e <shell>` pipe a shell to a network socket. This is the canonical reverse-shell payload. Most distro builds of `nc` have `-e` disabled for precisely this reason, so seeing it in a script is either an attacker backdoor or a deployment time bomb waiting on a different packaging of netcat. If you need a bidirectional pipe, use `socat TCP:... EXEC:...,pty` with an explicit authorization check and document the use.

Disable by adding `ZC1489` to `disabled_katas` in `.zshellcheckrc`.

---

<a id="zc1490"></a>
### ZC1490 — Error on `socat ... EXEC:<shell>` / `SYSTEM:<shell>` — socat reverse-shell pattern

**Severity:** `error`  
**Auto-fix:** `no`

The `EXEC:` and `SYSTEM:` socat address types spawn a subprocess connected to the other socat endpoint. Paired with `TCP:` or `TCP-LISTEN:`, they form the second-most-common reverse/bind shell payload after `nc -e`. Legitimate uses exist (test harnesses, pty brokers) but should be gated behind explicit authorization and a non-shell command. Scan hits are worth a look.

Disable by adding `ZC1490` to `disabled_katas` in `.zshellcheckrc`.

---

<a id="zc1491"></a>
### ZC1491 — Warn on `export LD_PRELOAD=...` / `LD_LIBRARY_PATH=...` — library injection

**Severity:** `warning`  
**Auto-fix:** `no`

Setting `LD_PRELOAD` in a script forces every subsequent dynamically-linked command to load the specified shared object first, a classic post-compromise privesc and persistence technique. Setting `LD_LIBRARY_PATH` to a writable path is a gentler variant of the same class. Legitimate uses exist (perf profiling, asan instrumentation) but should be scoped to a single invocation and the path pinned to a read-only location.

Disable by adding `ZC1491` to `disabled_katas` in `.zshellcheckrc`.

---

<a id="zc1492"></a>
### ZC1492 — Style: `at` / `batch` for deferred execution — prefer systemd timers for auditability

**Severity:** `style`  
**Auto-fix:** `no`

`at` and `batch` schedule one-shot deferred jobs via `atd`. The job payload lands in `/var/spool/at*/` with no unit file or dependency graph, which makes it harder to review in fleet audits, easier to miss in a compromise triage, and one of the less-watched places adversaries stash persistence. Prefer `systemd-run --on-calendar=` or a proper `.timer` unit with a corresponding `.service`.

Disable by adding `ZC1492` to `disabled_katas` in `.zshellcheckrc`.

---

<a id="zc1493"></a>
### ZC1493 — Warn on `watch -n 0` — zero-interval watch spins CPU

**Severity:** `warning`  
**Auto-fix:** `no`

`watch -n 0` (or `-n 0.0` / `-n .0`) tells `watch` to re-run the command with no delay, which immediately pins a core to 100% and usually saturates the terminal emulator too. Pick a realistic interval (`-n 1`, `-n 2`, `-n 0.5`) — or if you truly want tight polling, use a dedicated event API (`inotifywait`, `systemd.path` unit, `journalctl -f`).

Disable by adding `ZC1493` to `disabled_katas` in `.zshellcheckrc`.

---

<a id="zc1494"></a>
### ZC1494 — Warn on `tcpdump -w <file>` without `-Z <user>` — capture file owned by root

**Severity:** `warning`  
**Auto-fix:** `no`

`tcpdump` needs root (or CAP_NET_RAW) to open the raw socket, but once the socket is open it should drop privileges with `-Z <user>` before writing the pcap. Without `-Z`, the capture file is owned by root, any bpf filter bug is exercised with root privileges, and on a shared host the pcap can land with permissions that leak sensitive traffic to other users. Pair `-w` with `-Z tcpdump` (or a dedicated capture user).

Disable by adding `ZC1494` to `disabled_katas` in `.zshellcheckrc`.

---

<a id="zc1495"></a>
### ZC1495 — Warn on `ulimit -c unlimited` — enables core dumps from setuid binaries

**Severity:** `warning`  
**Auto-fix:** `no`

`ulimit -c unlimited` enables unbounded core dumps for the current shell and its children. On a system with `fs.suid_dumpable=2` and a world-readable coredump directory, a setuid process that segfaults leaks its memory into a file any user can read — Dirty COW-class keys, TLS session material, kerberos tickets. Leave core dumps at the distro default (usually 0) and use systemd-coredump with access controls if you genuinely need post-mortems.

Disable by adding `ZC1495` to `disabled_katas` in `.zshellcheckrc`.

---

<a id="zc1496"></a>
### ZC1496 — Error on reading `/dev/mem` / `/dev/kmem` / `/dev/port` — leaks physical memory

**Severity:** `error`  
**Auto-fix:** `no`

These device nodes map physical memory, kernel memory, and x86 I/O ports. Reading them (with `strings`, `xxd`, `cat`, or `dd`) exposes kernel state, keys, and any other live secret on the box. Modern kernels gate `/dev/mem` behind `CONFIG_STRICT_DEVMEM` but most distros also carry `CAP_SYS_RAWIO` on installed debugging tools, so the protection is fragile. If you really need a memory dump, use `kdump` + `crash` on a proper crash-kernel image.

Disable by adding `ZC1496` to `disabled_katas` in `.zshellcheckrc`.

---

<a id="zc1497"></a>
### ZC1497 — Error on `useradd -u 0` / `usermod -u 0` — creates a second root account

**Severity:** `error`  
**Auto-fix:** `no`

Creating a user with UID 0 makes them a second root — indistinguishable from `root` for every access decision, but hiding behind a non-obvious username (`backup`, `service`, `svc-updater`). This is a textbook persistence technique. If you need privileged but auditable operations, grant sudo rules tied to a specific non-0 UID and log via sudo's session plugin.

Disable by adding `ZC1497` to `disabled_katas` in `.zshellcheckrc`.

---

<a id="zc1498"></a>
### ZC1498 — Warn on `mount -o remount,rw /` — makes read-only root filesystem writable

**Severity:** `warning`  
**Auto-fix:** `no`

Remounting the root filesystem read-write is either an intentional config change that belongs in `/etc/fstab` (in which case this script is the wrong place) or a post-compromise step for persisting changes on an immutable / verity-backed root. On distros that ship with RO root (Fedora Silverblue, Chrome OS, appliance images) this also breaks rollback guarantees. Use `systemd-sysext` or `ostree admin deploy` for legitimate modifications.

Disable by adding `ZC1498` to `disabled_katas` in `.zshellcheckrc`.

---

<a id="zc1499"></a>
### ZC1499 — Style: `docker pull <image>` / `:latest` — unpinned image tag

**Severity:** `style`  
**Auto-fix:** `no`

Pulling without a tag defaults to `:latest`, which is a moving label. That breaks CI reproducibility (yesterday's build passed, today's fails for no reason the author changed) and reintroduces supply-chain surface every pull. Pin to a specific tag for convenience or to an immutable `@sha256:` digest for production.

Disable by adding `ZC1499` to `disabled_katas` in `.zshellcheckrc`.

---

<a id="zc1500"></a>
### ZC1500 — Warn on `systemctl edit <unit>` in scripts — requires interactive editor

**Severity:** `warning`  
**Auto-fix:** `no`

`systemctl edit <unit>` (without `--no-edit` and without a piped `EDITOR`) opens `$EDITOR` on a tmpfile and waits for the user. In a non-interactive script this either hangs until timeout or silently succeeds with no change, depending on how the editor handles a closed stdin. For scripted unit tweaks, drop a `.conf` drop-in under `/etc/systemd/system/<unit>.d/` and call `systemctl daemon-reload`.

Disable by adding `ZC1500` to `disabled_katas` in `.zshellcheckrc`.

---

<a id="zc1501"></a>
### ZC1501 — Style: `docker-compose` (hyphen) — use `docker compose` (space, built-in plugin)

**Severity:** `style`  
**Auto-fix:** `yes`

`docker-compose` is the Python Compose V1 binary. Docker stopped shipping it with Docker Desktop in 2023 and Compose V2 is now the first-class `docker compose` subcommand. Scripts that invoke `docker-compose` silently degrade on fresh installs and miss V2-only options (`--profile`, `--wait`, richer env interpolation). Call `docker compose` (space) or pin the V2 binary explicitly.

Disable by adding `ZC1501` to `disabled_katas` in `.zshellcheckrc`.

---

<a id="zc1502"></a>
### ZC1502 — Warn on `grep "$var" file` without `--` — flag injection when `$var` starts with `-`

**Severity:** `warning`  
**Auto-fix:** `no`

Without a `--` end-of-flags marker, `grep` (and most POSIX tools) treats any argument that starts with `-` as a flag. If `$var` comes from user input or a fuzzed filename, an attacker can pass `--include=*secret*` or `-f /etc/shadow` and get grep to read paths the script author never intended. Always write `grep -- "$var" file` or use a grep-compatible library with explicit pattern API.

Disable by adding `ZC1502` to `disabled_katas` in `.zshellcheckrc`.

---

<a id="zc1503"></a>
### ZC1503 — Error on `groupadd -g 0` / `groupmod -g 0` — creates duplicate root group

**Severity:** `error`  
**Auto-fix:** `no`

Creating or renaming a group to GID 0 gives its members the same privileges as members of `root` for every file that grants permissions to GID 0. Combined with `usermod -G 0 <user>` it becomes an invisible privilege escalation path. Distro tooling already reserves GID 0 for `root`; pick a sensible unused GID (`getent group` gives the list) and scope access via sudoers or polkit.

Disable by adding `ZC1503` to `disabled_katas` in `.zshellcheckrc`.

---

<a id="zc1504"></a>
### ZC1504 — Warn on `git push --mirror` — overwrites every remote ref

**Severity:** `warning`  
**Auto-fix:** `no`

`git push --mirror` pushes every ref under `refs/` and deletes any remote ref that is not present locally. Running it against a shared origin instantly wipes everyone else's branches and tags. Legitimate uses are mirror-to-mirror replication where the source is the authoritative tree; for everyday pushes use an explicit refspec or `git push --all`.

Disable by adding `ZC1504` to `disabled_katas` in `.zshellcheckrc`.

---

<a id="zc1505"></a>
### ZC1505 — Warn on `dpkg --force-confnew` / `--force-confold` — silently overrides /etc changes

**Severity:** `warning`  
**Auto-fix:** `no`

`--force-confnew` replaces any locally-modified config file with the maintainer version; `--force-confold` keeps the local file and drops the new defaults on the floor. Either way dpkg silently picks a side without prompting, so a legitimate /etc tweak (hardening, compliance override) can vanish or a security-relevant config update can be ignored. Review the conffile diff per upgrade (`ucf` / `etckeeper`) rather than hard-coding the decision.

Disable by adding `ZC1505` to `disabled_katas` in `.zshellcheckrc`.

---

<a id="zc1506"></a>
### ZC1506 — Warn on `newgrp <group>` in scripts — spawns a new shell, breaks control flow

**Severity:** `warning`  
**Auto-fix:** `no`

`newgrp` starts a new login shell with the requested primary group. Inside a non-interactive script that shell inherits no commands, so the script either hangs waiting for the user or exits immediately depending on stdin. If the script genuinely needs temporarily-augmented group access, call `sg <group> -c <cmd>` or, in a service context, use `SupplementaryGroups=` in the unit file.

Disable by adding `ZC1506` to `disabled_katas` in `.zshellcheckrc`.

---

<a id="zc1507"></a>
### ZC1507 — Warn on `rsync -l` / default symlink handling — follows escaping symlinks

**Severity:** `warning`  
**Auto-fix:** `no`

By default rsync copies symlinks as-is but does not prevent one from pointing outside the source tree. When the destination is rooted elsewhere (or the receiver creates a file at the symlink's resolved path) this becomes a path traversal primitive. Use `--safe-links` to skip symlinks pointing outside the transfer set, or `--copy-unsafe-links` to materialise them as regular files.

Disable by adding `ZC1507` to `disabled_katas` in `.zshellcheckrc`.

---

<a id="zc1508"></a>
### ZC1508 — Style: `ldd <binary>` may execute the binary — use `objdump -p` / `readelf -d` for untrusted files

**Severity:** `style`  
**Auto-fix:** `no`

On glibc, `ldd` is implemented by setting `LD_TRACE_LOADED_OBJECTS=1` and invoking the binary. A malicious ELF with a custom interpreter (`PT_INTERP`) or constructors can therefore run code when `ldd` is pointed at it. `objdump -p <file> | grep NEEDED` or `readelf -d <file>` give the same shared-library list without executing the binary.

Disable by adding `ZC1508` to `disabled_katas` in `.zshellcheckrc`.

---

<a id="zc1509"></a>
### ZC1509 — Warn on `trap '' TERM` / `trap - TERM` — ignores/resets fatal signal

**Severity:** `warning`  
**Auto-fix:** `no`

`trap '' <signal>` makes the signal uninterruptible. `trap - <signal>` restores the default disposition, which on `TERM`/`INT`/`HUP` means the script exits without running any cleanup handler. Both forms are routinely used to harden long-running scripts against accidental `Ctrl-C`, but also to hide from `kill` during incident response. Keep the explicit cleanup handler on at least `EXIT` so state is always unwound.

Disable by adding `ZC1509` to `disabled_katas` in `.zshellcheckrc`.

---

<a id="zc1510"></a>
### ZC1510 — Error on `auditctl -e 0` / `auditctl -D` — disables kernel audit logging

**Severity:** `error`  
**Auto-fix:** `no`

`auditctl -e 0` switches the Linux audit subsystem off, and `auditctl -D` deletes every audit rule, including the ones that monitor `/etc/shadow`, `execve`, and privilege escalations. Both are textbook anti-forensics steps. If you need to temporarily quiet audit for a maintenance window, use `-e 2` (lock enabled + immutable) to require a reboot for any further change and document the action.

Disable by adding `ZC1510` to `disabled_katas` in `.zshellcheckrc`.

---

<a id="zc1511"></a>
### ZC1511 — Error on `nmcli ... <wireless/vpn secret>` on command line

**Severity:** `error`  
**Auto-fix:** `no`

Passing Wi-Fi pre-shared keys or VPN secrets as positional `nmcli` args puts them in `ps`, shell history, and `/proc/<pid>/cmdline`. Let NetworkManager store the secret for you via `--ask` (interactive prompt, no TTY echo) or use `keyfile` connection profiles under `/etc/NetworkManager/system-connections/` with mode 0600.

Disable by adding `ZC1511` to `disabled_katas` in `.zshellcheckrc`.

---

<a id="zc1512"></a>
### ZC1512 — Style: `service <unit> <verb>` — use `systemctl <verb> <unit>` on systemd hosts

**Severity:** `style`  
**Auto-fix:** `yes`

`service` is the SysV init compatibility wrapper. On a systemd-managed host (every mainstream distro since ~2016) it translates to `systemctl` anyway, but reverses argument order, loses `--user` scope, ignores unit templating, and can't restart sockets or timers. Prefer `systemctl start|stop|restart|reload <unit>` for consistency across scripts and interactive shells.

Disable by adding `ZC1512` to `disabled_katas` in `.zshellcheckrc`.

---

<a id="zc1513"></a>
### ZC1513 — Style: `make install` without `DESTDIR=` — unmanaged system-wide install

**Severity:** `style`  
**Auto-fix:** `no`

`make install` drops files directly into `$(prefix)` with no package manager tracking. Upgrades can leave stale files behind, uninstalls rely on `make uninstall` being accurate, and the operation typically needs `sudo`. For local builds, set `DESTDIR=/tmp/pkgroot` + wrap in `checkinstall` / `fpm` / distro packaging, or use `stow` / `xstow` to manage symlinks under `/usr/local`.

Disable by adding `ZC1513` to `disabled_katas` in `.zshellcheckrc`.

---

<a id="zc1514"></a>
### ZC1514 — Error on `useradd -p <hash>` / `usermod -p <hash>` — password hash on cmdline

**Severity:** `error`  
**Auto-fix:** `no`

`-p` takes an already-hashed password (crypt(3) format) and writes it to `/etc/shadow`. That hash is in `ps`, `/proc/<pid>/cmdline`, and history for as long as the process runs — enough time for a co-tenant to grab it and start an offline crack. Use `chpasswd` with `--crypt-method=SHA512` reading from stdin, or write `/etc/shadow` via a configuration-management tool with proper file permissions.

Disable by adding `ZC1514` to `disabled_katas` in `.zshellcheckrc`.

---

<a id="zc1515"></a>
### ZC1515 — Warn on `md5sum` / `sha1sum` for integrity check — collision-vulnerable

**Severity:** `warning`  
**Auto-fix:** `no`

MD5 and SHA-1 are broken for collision resistance: public attacks cheaply craft two different files with the same hash. For verifying a download against a published checksum, or for comparing archives against a manifest, use `sha256sum` / `sha512sum` / `b2sum` instead. MD5 is still fine for non-adversarial cache keys but almost every invocation in scripts is the integrity case.

Disable by adding `ZC1515` to `disabled_katas` in `.zshellcheckrc`.

---

<a id="zc1516"></a>
### ZC1516 — Error on `umask 000` / `umask 0` — new files / directories world-writable

**Severity:** `error`  
**Auto-fix:** `no`

`umask 000` means every file created after this line inherits mode 0666 and every directory inherits 0777 — world-readable, world-writable, no authorization layer. On a multi-user host (build runner, shared workstation) this leaks secrets through the filesystem and invites tampering. Pick a sensible umask (`022` for public software, `077` for secrets handling).

Disable by adding `ZC1516` to `disabled_katas` in `.zshellcheckrc`.

---

<a id="zc1517"></a>
### ZC1517 — Warn on `print -P "$var"` — prompt-escape injection via user-controlled string

**Severity:** `warning`  
**Auto-fix:** `no`

`print -P` enables prompt-escape expansion (`%F`, `%K`, `%B`, `%S`, plus arbitrary command substitution via `%{...%}`). Interpolating a shell variable means any of those sequences inside the variable are expanded — at best messing up terminal state, at worst running the attacker's command via `%(e:...)` or similar. Either drop `-P` or wrap the variable with `${(q-)var}` / `${(V)var}` to neutralize metacharacters before printing.

Disable by adding `ZC1517` to `disabled_katas` in `.zshellcheckrc`.

---

<a id="zc1518"></a>
### ZC1518 — Warn on `bash -p` — privileged mode (skips env sanitisation on setuid)

**Severity:** `warning`  
**Auto-fix:** `no`

`bash -p` (and `-o privileged`) tells bash not to drop its effective UID/GID and not to sanitize the environment when started on a setuid wrapper. It is explicitly the flag you use to keep `BASH_ENV`, `SHELLOPTS`, and similar attacker-controlled variables active while running as a more privileged user. Almost no legitimate script needs `-p`; audit and remove.

Disable by adding `ZC1518` to `disabled_katas` in `.zshellcheckrc`.

---

<a id="zc1519"></a>
### ZC1519 — Warn on `ulimit -u unlimited` — removes user process cap, enables fork bombs

**Severity:** `warning`  
**Auto-fix:** `no`

`ulimit -u` caps the number of processes a UID can run; `unlimited` removes that cap. Combined with a bug in a background loop (or a literal fork bomb via `:(){ :|:& };:`) it pegs the scheduler until the machine has to be cold-booted. Pick a realistic number (distro defaults around 4096 for interactive sessions) or set it in `/etc/security/limits.d/` so it is persistent and visible.

Disable by adding `ZC1519` to `disabled_katas` in `.zshellcheckrc`.

---

<a id="zc1520"></a>
### ZC1520 — Warn on `vared <var>` in scripts — reads interactively, hangs non-interactive

**Severity:** `warning`  
**Auto-fix:** `no`

`vared` is the Zsh interactive line-editor builtin that lets the user edit the value of a variable in place. In a non-interactive script (cron job, CI runner, ssh-with-command) `vared` has no TTY, so the script either errors out or hangs waiting for input that never arrives. For scripted input, read the value from stdin (`read varname`), a file, or an environment variable.

Disable by adding `ZC1520` to `disabled_katas` in `.zshellcheckrc`.

---

<a id="zc1521"></a>
### ZC1521 — Style: `strace` without `-e` filter — captures every syscall (incl. secrets, huge output)

**Severity:** `style`  
**Auto-fix:** `no`

Unfiltered `strace` records every syscall the process makes: every `read()`/`write()` buffer, every `connect()` sockaddr, every `open()` path. That includes passwords read from stdin, session tokens written to TLS sockets, and any memory a `write()` buffer happens to point at. Scope with `-e trace=<set>` (e.g. `trace=openat,connect`) and strip sensitive content with `-e abbrev=all`.

Disable by adding `ZC1521` to `disabled_katas` in `.zshellcheckrc`.

---

<a id="zc1522"></a>
### ZC1522 — Warn on `ip route add default` / `route add default` — changes default gateway

**Severity:** `warning`  
**Auto-fix:** `no`

Setting a new default route in a script silently redirects every non-local packet through the specified gateway. That is exactly the knob an attacker turns to MITM a whole host after a foothold, and it is also a common accidental foot- gun in CI runners (gateway in the runner network ≠ gateway in production). Use NetworkManager / systemd-networkd config files for persistent routes, and document any runtime change with a comment explaining why.

Disable by adding `ZC1522` to `disabled_katas` in `.zshellcheckrc`.

---

<a id="zc1523"></a>
### ZC1523 — Error on `tar -C /` — extracting an archive into the filesystem root

**Severity:** `error`  
**Auto-fix:** `no`

Extracting a tarball directly into `/` overwrites any file it carries a matching path for. Combined with a malicious tarball that contains entries like `etc/pam.d/sshd` or `usr/bin/ls`, this is a full system compromise disguised as a software install. Always extract into a staging directory, inspect contents, then copy specific files into place.

Disable by adding `ZC1523` to `disabled_katas` in `.zshellcheckrc`.

---

<a id="zc1524"></a>
### ZC1524 — Warn on `sysctl -e` / `sysctl -q` — silently skip unknown keys, hide config drift

**Severity:** `warning`  
**Auto-fix:** `no`

`sysctl -e` and `-q` suppress error output for unknown keys or failed writes. That is how a typo in `/etc/sysctl.d/99-hardening.conf` goes unnoticed for months — the hardening didn't actually take effect because the key name was wrong. Drop `-e`/`-q` in scripts and let errors bubble up; fix the offending conffile instead.

Disable by adding `ZC1524` to `disabled_katas` in `.zshellcheckrc`.

---

<a id="zc1525"></a>
### ZC1525 — Warn on `ping -f` — flood ping sends packets as fast as possible

**Severity:** `warning`  
**Auto-fix:** `no`

`ping -f` (flood mode) removes the one-per-second rate limit and sends ICMP echo requests in a tight loop. It's a root-only builtin specifically because it can saturate a slow link or overload a low-end host. Legitimate uses exist (latency benchmarking, stress testing known-internal targets), but in a script aimed at arbitrary hosts it is a noisy traffic generator. Scope tightly and document.

Disable by adding `ZC1525` to `disabled_katas` in `.zshellcheckrc`.

---

<a id="zc1526"></a>
### ZC1526 — Error on `wipefs -a` / `wipefs -af` — erases filesystem signatures (unrecoverable)

**Severity:** `error`  
**Auto-fix:** `no`

`wipefs -a` overwrites every filesystem, partition table, and RAID signature it finds on the target. Unlike `rm`, there is no retention anywhere — the only recovery path is a disk image backup taken beforehand. If the target variable is wrong (typo, empty, resolves to the wrong `/dev/sdX`), this bricks the disk. Always run with `--no-act` first or prefer `sgdisk --zap-all` for partition-table scope.

Disable by adding `ZC1526` to `disabled_katas` in `.zshellcheckrc`.

---

<a id="zc1527"></a>
### ZC1527 — Warn on `crontab -` — replaces cron from stdin, overwrites without diff

**Severity:** `warning`  
**Auto-fix:** `no`

`crontab -` (or `crontab -u <user> -`) reads a full crontab from stdin and replaces the user's existing entries wholesale. Any manual tweak, oncall override, or colleague's row is silently deleted. Paired with `curl | crontab -` it is a common persistence one-liner. Use `crontab -l > /tmp/old && ... crontab -e` with an explicit diff/merge, or ship cron entries via `/etc/cron.d/*` managed by config tooling.

Disable by adding `ZC1527` to `disabled_katas` in `.zshellcheckrc`.

---

<a id="zc1528"></a>
### ZC1528 — Warn on `chage -M 99999` / `-E -1` — disables password aging / expiry

**Severity:** `warning`  
**Auto-fix:** `no`

`chage -M 99999` sets the max password age to roughly 273 years (effectively never). `chage -E -1` clears the account expiration date. Both silently remove an automatic lockout mechanism a compromised credential would otherwise hit. If passwords genuinely should not expire (SSO, cert-based auth), encode that in a PAM profile rather than per-user `chage`.

Disable by adding `ZC1528` to `disabled_katas` in `.zshellcheckrc`.

---

<a id="zc1529"></a>
### ZC1529 — Warn on `fsck -y` / `fsck.<fs> -y` — auto-answer yes can corrupt

**Severity:** `warning`  
**Auto-fix:** `no`

`fsck -y` answers `yes` to every repair prompt. For the happy case it is a timesaver, but on a filesystem with unusual corruption (bad sector storm, mangled journal after power loss) the automatic answer can turn salvageable data into `lost+found` entries or zero it outright. In scripts, prefer `fsck -n` for a dry-run and let a human adjudicate a real repair, or run with `-p` (preen: only safe automatic fixes).

Disable by adding `ZC1529` to `disabled_katas` in `.zshellcheckrc`.

---

<a id="zc1530"></a>
### ZC1530 — Warn on `pkill -f <pattern>` — matches full command line, easy to over-kill

**Severity:** `warning`  
**Auto-fix:** `no`

`pkill -f` matches the pattern against the full command line, not just the process name. A pattern like `-f server` also matches the `grep -- server` in a user's shell history or any backup tool named `server-backup`. For routine use, drop `-f` (matches process name only) or scope with `-U <uid>` / `-G <gid>` / `-P <ppid>`. When you must match the command line, pin it with `^` / `$` anchors in the pattern.

Disable by adding `ZC1530` to `disabled_katas` in `.zshellcheckrc`.

---

<a id="zc1531"></a>
### ZC1531 — Warn on `wget -t 0` — infinite retries, hangs on a dead endpoint

**Severity:** `warning`  
**Auto-fix:** `no`

`wget -t 0` (or `--tries=0`) means retry forever. Paired with `-w` (wait between retries) and a dead endpoint, the script hangs until killed — in a cron job, every subsequent invocation piles up and eventually the UID's process limit trips. Use a finite retry count (`-t 5`) plus `--timeout=<seconds>` to cap total wall time.

Disable by adding `ZC1531` to `disabled_katas` in `.zshellcheckrc`.

---

<a id="zc1532"></a>
### ZC1532 — Warn on `screen -dm` / `tmux new-session -d` — detached long-running session

**Severity:** `warning`  
**Auto-fix:** `no`

Starting a detached screen/tmux session from a script puts a long-running process outside the systemd supervisory tree: no logs in the journal, no cgroup accounting, no restart-on-failure, no OOM scoring. It is also a common post- compromise persistence technique because the session survives the initial shell exit and hides in `ps -ef` as a short tmux/screen helper. For real long-running work, write a systemd unit (user or system) and start it with `systemctl [--user] start`.

Disable by adding `ZC1532` to `disabled_katas` in `.zshellcheckrc`.

---

<a id="zc1533"></a>
### ZC1533 — Warn on `setsid <cmd>` — detaches from controlling TTY, escapes supervision

**Severity:** `warning`  
**Auto-fix:** `no`

`setsid` starts a new session and process group. Combined with `-f` (`--fork`) the child is fully detached from the invoking shell: `SIGHUP` from logout does not reach it, the tty hang-up no longer terminates it, and it falls off the script's job table. That is legitimate for daemonising a long-running helper (though systemd does this better) and is also a standard persistence mechanism. Prefer a systemd unit; if you must detach, document why.

Disable by adding `ZC1533` to `disabled_katas` in `.zshellcheckrc`.

---

<a id="zc1534"></a>
### ZC1534 — Warn on `dmesg -c` / `--clear` — wipes kernel ring buffer

**Severity:** `warning`  
**Auto-fix:** `no`

`dmesg -c` reads and then clears the kernel ring buffer. Any subsequent reader sees an empty log, so OOM kills, driver panics, and audit messages that landed between the wipe and the incident response are gone. It is also an anti-forensics step in post-exploitation playbooks. Use `dmesg` (no flags) for a read, and let the journal retention policy handle rotation.

Disable by adding `ZC1534` to `disabled_katas` in `.zshellcheckrc`.

---

<a id="zc1535"></a>
### ZC1535 — Warn on `ip link set <iface> promisc on` — enables packet capture

**Severity:** `warning`  
**Auto-fix:** `no`

Putting an interface into promiscuous mode tells the NIC to deliver every frame to userspace, not just frames addressed to this host. Legitimate for tools like tcpdump/tshark (which turn it on themselves) but running it from a script and leaving it on is a sniffer-in-place — traffic from other hosts on the same broadcast domain lands in anyone's `tshark -i`. Re-disable as soon as capture is done, and prefer giving tcpdump `CAP_NET_RAW` so the mode is scoped to a single invocation.

Disable by adding `ZC1535` to `disabled_katas` in `.zshellcheckrc`.

---

<a id="zc1536"></a>
### ZC1536 — Warn on `iptables -j DNAT` / `-j REDIRECT` — rewrites traffic destination

**Severity:** `warning`  
**Auto-fix:** `no`

`-j DNAT` and `-j REDIRECT` in an iptables rule rewrite the destination address/port of matching packets. That is how you transparently proxy, but also how you silently redirect a victim's connections to an attacker-controlled port. Scripts that touch NAT rules should be carefully reviewed; prefer declarative network config (nftables ruleset, NetworkManager connection, firewalld service) and store rule provenance.

Disable by adding `ZC1536` to `disabled_katas` in `.zshellcheckrc`.

---

<a id="zc1537"></a>
### ZC1537 — Error on `lvremove -f` / `vgremove -f` / `pvremove -f` — force-destroys LVM metadata

**Severity:** `error`  
**Auto-fix:** `no`

The `-f`/`--force` flag on the LVM destructive commands skips the confirmation prompt that protects against a typo in the volume name. If the target variable resolves to the wrong VG/LV/PV (empty, unset, different host), a single line destroys every filesystem on top of that LVM stack. Leave the prompt in and pipe `yes` to it only when you have explicitly confirmed the target immediately beforehand.

Disable by adding `ZC1537` to `disabled_katas` in `.zshellcheckrc`.

---

<a id="zc1538"></a>
### ZC1538 — Error on `zpool destroy -f` / `zfs destroy -rR` — recursive ZFS destruction

**Severity:** `error`  
**Auto-fix:** `no`

`zpool destroy -f` nukes a whole ZFS pool including every dataset, snapshot, and clone on it. `zfs destroy -r` recurses into descendant datasets; `-R` additionally drops descendant clones. Unlike `rm`, the space is freed immediately and there is no recycle bin. Always require `zfs list`/`zpool list` + explicit target confirmation in the same script block, and prefer snapshot-based rollback for recoverable workflows.

Disable by adding `ZC1538` to `disabled_katas` in `.zshellcheckrc`.

---

<a id="zc1539"></a>
### ZC1539 — Warn on `parted -s <disk> <destructive-op>` — script mode bypasses confirmation

**Severity:** `warning`  
**Auto-fix:** `no`

`parted -s` (script mode) answers the `data will be destroyed` prompt with `yes`. Combined with `mklabel`, `mkpart`, `rm`, or `resizepart` on the wrong device variable it silently repartitions or zeros the partition table on a disk the author never intended. Require an explicit `parted <disk> print` check plus an out-of-band confirmation before the destructive call.

Disable by adding `ZC1539` to `disabled_katas` in `.zshellcheckrc`.

---

<a id="zc1540"></a>
### ZC1540 — Error on `cryptsetup erase` / `luksErase` — destroys LUKS header, data unrecoverable

**Severity:** `error`  
**Auto-fix:** `no`

`cryptsetup erase` (alias `luksErase`) overwrites the LUKS header and every key slot. Without the header the ciphertext on the device is unrecoverable — even the original passphrase cannot unlock it. Keep a `cryptsetup luksHeaderBackup` image somewhere safe before running erase, and prefer `luksRemoveKey`/`luksKillSlot` when only rotating one slot.

Disable by adding `ZC1540` to `disabled_katas` in `.zshellcheckrc`.

---

<a id="zc1541"></a>
### ZC1541 — Error on `apk add --allow-untrusted` — installs unsigned Alpine packages

**Severity:** `error`  
**Auto-fix:** `no`

`apk add --allow-untrusted` skips signature verification on the package being installed. On Alpine that is a direct MITM-to-root path: any mirror, cache, or typo-squat can slip a replacement `.apk` and the daemon starts running attacker code on next restart. Sign internal packages with your own key in `/etc/apk/keys/` and keep verification on.

Disable by adding `ZC1541` to `disabled_katas` in `.zshellcheckrc`.

---

<a id="zc1542"></a>
### ZC1542 — Error on `snap install --dangerous` — installs unsigned snap

**Severity:** `error`  
**Auto-fix:** `no`

`snap install --dangerous` tells snapd to install a snap that is not assertion-verified. That bypass is named after the risk: any `.snap` file on disk can register system services, confinement profiles, and hooks, running as whatever user the snap declares. Use `--devmode` for developer work (still verified) or ship the snap through the store / a private brand store for production rollouts.

Disable by adding `ZC1542` to `disabled_katas` in `.zshellcheckrc`.

---

<a id="zc1543"></a>
### ZC1543 — Warn on `go install pkg@latest` / `cargo install --git <url>` without rev pin

**Severity:** `warning`  
**Auto-fix:** `no`

`go install pkg@latest` and `cargo install --git <url>` without `--rev` / `--tag` / `--branch` resolve to whatever HEAD is at install time. The next CI run can pull a different commit — great for supply-chain attackers to inject post-breach, bad for reproducibility. Pin to a specific version tag (`pkg@v1.2.3`) or a commit hash (`cargo install --rev abc123 --git ...`).

Disable by adding `ZC1543` to `disabled_katas` in `.zshellcheckrc`.

---

<a id="zc1544"></a>
### ZC1544 — Warn on `dnf copr enable` / `add-apt-repository ppa:` — unvetted third-party repo

**Severity:** `warning`  
**Auto-fix:** `no`

Enabling a COPR project or an Ubuntu PPA pulls packages signed by a single community contributor — there is no distro security team or reproducible-builds guarantee behind that key. Any future compromise of that contributor's account ships a rootkit to every box that ran this line. If you need the package badly enough, pin to a specific `build-id`, verify the key fingerprint out of band, and mirror to an internal repository.

Disable by adding `ZC1544` to `disabled_katas` in `.zshellcheckrc`.

---

<a id="zc1545"></a>
### ZC1545 — Warn on `docker system prune -af --volumes` — drops unused volumes too

**Severity:** `warning`  
**Auto-fix:** `no`

`docker system prune -af --volumes` removes stopped containers, unused networks, dangling images — and every volume not currently attached to a running container. On a host where `docker-compose down` is used casually (shutdown before a laptop close, for example), the matching database volume looks "unused" to prune and goes with it. Drop `--volumes` from routine cleanup, or target specific prune scopes (`docker image prune`, `docker container prune`).

Disable by adding `ZC1545` to `disabled_katas` in `.zshellcheckrc`.

---

<a id="zc1546"></a>
### ZC1546 — Warn on `kubectl delete --force --grace-period=0` — skips PreStop, corrupts state

**Severity:** `warning`  
**Auto-fix:** `no`

`kubectl delete --force --grace-period=0` tells the API server to remove the resource from etcd without waiting for the kubelet to run PreStop hooks or drain the pod. For a StatefulSet pod this routinely corrupts the backing PV (database mid-flush, file lock left held) and the replacement pod refuses to start. Use standard delete and let the graceful shutdown run; only reach for `--force` when the node itself is gone.

Disable by adding `ZC1546` to `disabled_katas` in `.zshellcheckrc`.

---

<a id="zc1547"></a>
### ZC1547 — Warn on `kubectl apply --prune --all` — deletes resources missing from manifest

**Severity:** `warning`  
**Auto-fix:** `no`

`kubectl apply --prune --all` (or `--prune -l <selector>`) deletes every cluster resource whose label matches but which is not in the manifest you just applied. In a partial-repo deploy or a manifest typo, that can delete production Deployments, Services, or Secrets another team owns. Pair `--prune` with a narrow `-l` selector unique to your stack, or use a GitOps controller (Argo CD, Flux) that scopes prune to its own Application.

Disable by adding `ZC1547` to `disabled_katas` in `.zshellcheckrc`.

---

<a id="zc1548"></a>
### ZC1548 — Warn on `helm install/upgrade --disable-openapi-validation` — skips schema check

**Severity:** `warning`  
**Auto-fix:** `no`

`--disable-openapi-validation` tells Helm to skip the OpenAPI schema check the API server would apply. Malformed CRD instances or Deployments with invalid spec fields then silently land in etcd, only failing when the controller tries to reconcile — usually 3am, usually in prod. Keep the validation on; fix the schema deviation instead.

Disable by adding `ZC1548` to `disabled_katas` in `.zshellcheckrc`.

---

<a id="zc1549"></a>
### ZC1549 — Error on `unzip -d /` / `unzip -o ... -d /` — extract archive into filesystem root

**Severity:** `error`  
**Auto-fix:** `no`

Unzipping directly into `/` (or `/root`, `/boot`) overwrites any system file whose path matches an entry in the archive. A malicious zip that carries `etc/passwd`, `usr/bin/ls`, or `root/.ssh/authorized_keys` turns a seemingly harmless extract into full system compromise. Stage to a scratch directory, inspect contents, then copy or install specific files.

Disable by adding `ZC1549` to `disabled_katas` in `.zshellcheckrc`.

---

<a id="zc1550"></a>
### ZC1550 — Warn on `apt-mark hold <pkg>` — pins a package, blocks security updates

**Severity:** `warning`  
**Auto-fix:** `no`

`apt-mark hold` tells apt to leave the package at its current version on `apt upgrade` and `unattended-upgrades`. That is occasionally correct (pinning a kernel variant for a driver, or a broken-upstream version) but silently keeps the package vulnerable to every subsequent CVE. Document the reason in a comment, schedule a review, and prefer `apt-mark unhold` + `apt upgrade <pkg>` over leaving the pin in place indefinitely.

Disable by adding `ZC1550` to `disabled_katas` in `.zshellcheckrc`.

---

<a id="zc1551"></a>
### ZC1551 — Warn on `helm install/upgrade --skip-crds` — chart CRs land before their CRDs

**Severity:** `warning`  
**Auto-fix:** `no`

`--skip-crds` tells Helm to install only the `.Release` objects and skip the CustomResourceDefinition manifests under `crds/`. Without the CRDs present, any `.Release` object that references a custom resource is rejected by the API server at validation time, or — worse — fails later when a reconciler tries to watch a type that does not exist. Use the default (install CRDs) on first roll- out; if you need split lifecycle, install CRDs manually (`kubectl apply -f chart/crds/`) before the `helm install`.

Disable by adding `ZC1551` to `disabled_katas` in `.zshellcheckrc`.

---

<a id="zc1552"></a>
### ZC1552 — Warn on `openssl dhparam <2048` / `genrsa <2048` — weak key/parameter size

**Severity:** `warning`  
**Auto-fix:** `no`

Generating DH parameters or RSA keys shorter than 2048 bits is below every modern compliance baseline (NIST SP 800-57, BSI TR-02102, Mozilla Server Side TLS). A 1024-bit RSA modulus or DH group is within reach of academic precomputation (Logjam) and a 512-bit one was broken on commodity hardware in the 1990s. Use 2048 as a floor and 3072 / 4096 for long-lived keys.

Disable by adding `ZC1552` to `disabled_katas` in `.zshellcheckrc`.

---

<a id="zc1553"></a>
### ZC1553 — Style: use Zsh `${(U)var}` / `${(L)var}` instead of `tr '[:lower:]' '[:upper:]'`

**Severity:** `style`  
**Auto-fix:** `no`

Zsh provides `${(U)var}` and `${(L)var}` parameter-expansion flags for case conversion in-process. Spawning `tr` for this forks/execs per call (noticeable in a hot loop), relies on the external `tr` being POSIX-compliant (BusyBox and old macOS differ), and round-trips the data through a pipe. Drop `tr` for the built-in: `upper=${(U)lower}` / `lower=${(L)upper}`.

Disable by adding `ZC1553` to `disabled_katas` in `.zshellcheckrc`.

---

<a id="zc1554"></a>
### ZC1554 — Warn on `unzip -o` / `tar ... --overwrite` — silent overwrite during extract

**Severity:** `warning`  
**Auto-fix:** `no`

`unzip -o` overwrites existing files without prompting; `tar --overwrite` does the same for tarballs. In a directory that already contains user work or a previous release, a newer archive silently wins, discarding in-flight edits and custom config. Extract to a fresh staging directory, diff, then move specific files into place.

Disable by adding `ZC1554` to `disabled_katas` in `.zshellcheckrc`.

---

<a id="zc1555"></a>
### ZC1555 — Error on `chmod` / `chown` on `/etc/shadow` or `/etc/sudoers` (managed files)

**Severity:** `error`  
**Auto-fix:** `no`

`/etc/shadow`, `/etc/gshadow`, `/etc/sudoers`, and `/etc/passwd` have specific ownership and mode invariants that the distro `passwd`, `chage`, and `visudo` tools maintain atomically with file locking. Direct `chmod`/`chown` races those tools, can leave the file world-readable mid-modification (leaking the shadow file), and will be clobbered on the next `shadow -p` run. Use the proper wrapper, or ship a configuration-management drop-in.

Disable by adding `ZC1555` to `disabled_katas` in `.zshellcheckrc`.

---

<a id="zc1556"></a>
### ZC1556 — Error on `openssl enc -des` / `-rc4` / `-3des` — broken symmetric cipher

**Severity:** `error`  
**Auto-fix:** `no`

DES, RC4, and 3DES are all broken or on-deprecation-path: DES's 56-bit key fell to commodity brute-force decades ago, RC4 has practical biased-output attacks, and 3DES suffers the Sweet32 birthday collision when reused for more than ~32GB. None of them provide authenticity either. Use `-aes-256-gcm` or `-chacha20-poly1305`, or move up to a dedicated tool (`age`, `gpg`, `libsodium`).

Disable by adding `ZC1556` to `disabled_katas` in `.zshellcheckrc`.

---

<a id="zc1557"></a>
### ZC1557 — Error on `kubeadm reset -f` / `--force` — wipes Kubernetes control-plane state

**Severity:** `error`  
**Auto-fix:** `no`

`kubeadm reset` stops kubelet, tears down static-pod manifests, clears `/etc/kubernetes`, and (with `-f`) skips the confirmation that protects a mistyped target. On a control-plane node it also breaks every tenant that relied on that etcd quorum. Drain first, remove the node from the cluster, then run reset interactively to confirm.

Disable by adding `ZC1557` to `disabled_katas` in `.zshellcheckrc`.

---

<a id="zc1558"></a>
### ZC1558 — Warn on `usermod -aG wheel|sudo|root|adm` — silent privilege group escalation

**Severity:** `warning`  
**Auto-fix:** `no`

Adding a user to `wheel`, `sudo`, `root`, `adm`, `docker`, or `libvirt` from a script grants persistent admin-level access without the review a sudoers drop-in or PAM profile would get. `docker` and `libvirt` in particular are equivalent to root (spawn privileged containers / raw disk access). Use a sudoers.d file scoped to specific commands and audit changes in configuration management.

Disable by adding `ZC1558` to `disabled_katas` in `.zshellcheckrc`.

---

<a id="zc1559"></a>
### ZC1559 — Warn on `ssh-copy-id -f` / `-o StrictHostKeyChecking=no` — trust-on-first-use key push

**Severity:** `warning`  
**Auto-fix:** `no`

`ssh-copy-id` opens an SSH connection to deposit the caller's public key. With `-f` it overwrites existing `authorized_keys` without prompting; with `-o StrictHostKeyChecking=no` it does not verify the host key. Together they push a long-term credential at a host the script has never authenticated — a network MITM lands a permanent backdoor. Verify the target host's fingerprint out of band before pushing keys.

Disable by adding `ZC1559` to `disabled_katas` in `.zshellcheckrc`.

---

<a id="zc1560"></a>
### ZC1560 — Error on `pip install --trusted-host` — accepts MITM / plain-HTTP PyPI index

**Severity:** `error`  
**Auto-fix:** `no`

`--trusted-host` tells pip to skip TLS certificate verification for the specified host and to allow plain-HTTP URLs from that host. Any MITM on the path can substitute packages on install, and a typo in the host name means every subsequent `install` from the misspelled host is unauthenticated. Fix the CA trust (install the real corporate CA) instead of silencing pip, and keep the default `--index-url https://...` over the TLS-verified endpoint.

Disable by adding `ZC1560` to `disabled_katas` in `.zshellcheckrc`.

---

<a id="zc1561"></a>
### ZC1561 — Error on `systemctl isolate rescue.target` / `emergency.target` from a script

**Severity:** `error`  
**Auto-fix:** `no`

`systemctl isolate rescue.target` drops the host into single-user rescue mode; `emergency.target` goes even further, leaving only the root shell on the console. Both terminate networking, SSH sessions, and most services. On a remote host the script loses its own connection mid-run, and anyone relying on the box is cut off without warning. Reserve these for console recovery, not script flow.

Disable by adding `ZC1561` to `disabled_katas` in `.zshellcheckrc`.

---

<a id="zc1562"></a>
### ZC1562 — Warn on `env -u PATH` / `-u LD_LIBRARY_PATH` — clears security-relevant env

**Severity:** `warning`  
**Auto-fix:** `no`

`env -u PATH` unsets the caller's `PATH` before running the child, forcing the child to fall back to the hard-coded search list (`/bin:/usr/bin` on glibc). That bypasses PATH hardening done by the parent shell (e.g. a sanitised PATH under `sudo`). Unsetting `LD_PRELOAD` / `LD_LIBRARY_PATH` mid-stream is also usually the caller trying to shake off an earlier `export`. Either use `env -i` to sanitise completely, or explicitly set the variables the child should see.

Disable by adding `ZC1562` to `disabled_katas` in `.zshellcheckrc`.

---

<a id="zc1563"></a>
### ZC1563 — Warn on `swapoff -a` — disables swap (memory pressure, potential OOM)

**Severity:** `warning`  
**Auto-fix:** `no`

`swapoff -a` turns off every active swap device. Kubelet installers do this because kubelet refuses to run with swap, but leaving it in a general-purpose script means the next memory-hungry process on the host hits the OOM killer instead of paging. If the goal is kubelet-friendly, also remove the swap entry from `/etc/fstab` and document the trade-off; otherwise keep swap on.

Disable by adding `ZC1563` to `disabled_katas` in `.zshellcheckrc`.

---

<a id="zc1564"></a>
### ZC1564 — Warn on `date -s` / `timedatectl set-time` — manual clock change breaks TLS / cron

**Severity:** `warning`  
**Auto-fix:** `no`

Setting the system clock by hand (`date -s`, `timedatectl set-time`, `hwclock --set`) moves wall-clock time enough to invalidate short-lived TLS certificates, reset `cron`'s missed-job catch-up, and confuse `systemd.timer` units that depend on monotonic math. Use `systemd-timesyncd` / `chrony` / `ntpd` for routine correction; reserve manual set for first-boot bootstrap or air-gapped recovery and document the action.

Disable by adding `ZC1564` to `disabled_katas` in `.zshellcheckrc`.

---

<a id="zc1565"></a>
### ZC1565 — Style: use `command -v` instead of `whereis` / `locate` for command existence

**Severity:** `style`  
**Auto-fix:** `yes`

`whereis` searches a hard-coded list of binary/manual/source directories and returns everything it finds, including stale paths on custom `$PATH` layouts. `locate` relies on a cron-maintained index that may be hours or days stale. For a scripted "does this command exist?" check, `command -v <cmd>` respects the current `$PATH`, returns the selected resolution, and has no index-refresh coupling.

Disable by adding `ZC1565` to `disabled_katas` in `.zshellcheckrc`.

---

<a id="zc1566"></a>
### ZC1566 — Error on `gem install -P NoSecurity|LowSecurity` / `--trust-policy NoSecurity`

**Severity:** `error`  
**Auto-fix:** `no`

RubyGems' trust policy decides what signatures the installer accepts. `NoSecurity` skips signature verification entirely; `LowSecurity` warns but still installs unsigned gems. On a registry MITM or a hijacked maintainer account those policies turn into arbitrary code execution at gem-install time. Use `HighSecurity` (reject all but fully-signed) or `MediumSecurity` for hybrid repos.

Disable by adding `ZC1566` to `disabled_katas` in `.zshellcheckrc`.

---

<a id="zc1567"></a>
### ZC1567 — Warn on `python -m http.server` without `--bind 127.0.0.1` — serves to all interfaces

**Severity:** `warning`  
**Auto-fix:** `no`

`python -m http.server` (and the legacy `SimpleHTTPServer`) default to `0.0.0.0`, exposing the current directory's contents to every network the host is on. Tmp scratch files, `.env`, SSH keys, or a `node_modules` tree with private config all become reachable from anywhere on the LAN (or the internet, on a VPS). Pass `--bind 127.0.0.1` (or `--bind ::1`) unless you really need external access and know what is in the cwd.

Disable by adding `ZC1567` to `disabled_katas` in `.zshellcheckrc`.

---

<a id="zc1568"></a>
### ZC1568 — Error on `useradd -o` / `usermod -o` — allows non-unique UID (alias user)

**Severity:** `error`  
**Auto-fix:** `no`

`-o` (or `--non-unique`) lets `useradd` / `usermod` assign a UID that is already in use. The new account has the same kernel identity as the existing one but its own login name, password, shell, and home dir. It is indistinguishable in `ps` / audit / file ACLs, so a compromise of either account is a compromise of both. Pick a fresh UID instead.

Disable by adding `ZC1568` to `disabled_katas` in `.zshellcheckrc`.

---

<a id="zc1569"></a>
### ZC1569 — Error on `nvme format -s1` / `-s2` — cryptographic or full-block SSD erase

**Severity:** `error`  
**Auto-fix:** `no`

`nvme format -s1` does a cryptographic erase of the target namespace; `-s2` (or the full-NVMe sanitize) rewrites every block. Both are unrecoverable in seconds. On a typo in the device variable — or a script that iterates over `/dev/nvme*n*` and catches the wrong namespace — the wrong disk is gone by the time the operator notices. Run interactively on verified targets, or not at all from automation.

Disable by adding `ZC1569` to `disabled_katas` in `.zshellcheckrc`.

---

<a id="zc1570"></a>
### ZC1570 — Warn on `smbclient -N` / `mount.cifs guest` — anonymous SMB share access

**Severity:** `warning`  
**Auto-fix:** `no`

`smbclient -N` skips authentication entirely (anonymous / null session); `mount.cifs` with `guest,username=` or `-o guest` does the same at the mount layer. Any host on the network segment can then read the share. If the share is truly public (software mirror, build cache) wrap in a read-only filesystem and document it; otherwise require Kerberos (`-k`) or pass credentials via `credentials=<file>` with 0600 perms.

Disable by adding `ZC1570` to `disabled_katas` in `.zshellcheckrc`.

---

<a id="zc1571"></a>
### ZC1571 — Style: `ntpdate` is deprecated — use `chronyc makestep` / `systemd-timesyncd`

**Severity:** `style`  
**Auto-fix:** `no`

`ntpdate` was retired by the ntp.org project around 4.2.6. Distros increasingly ship without it, and packaging it breaks the invariant that only one program writes the clock at a time (if `chrony` or `timesyncd` is also running the two fight). Use `chronyc makestep` (if chrony is active) or `systemctl restart systemd-timesyncd` (if timesyncd is active) for a one-shot step, and leave the daemon to keep it synchronised.

Disable by adding `ZC1571` to `disabled_katas` in `.zshellcheckrc`.

---

<a id="zc1572"></a>
### ZC1572 — Warn on `docker run -e PASSWORD=<value>` — secret in container env / inspect

**Severity:** `warning`  
**Auto-fix:** `no`

Passing a secret through `docker run -e NAME=value` puts it in the output of `docker inspect`, the container's `/proc/1/environ` (readable by anything that shares the PID namespace), and the shell history of whoever launched the container. Use `--env-file` with 0600 perms, a secret-mount `--secret` via BuildKit / Swarm, or mount a tmpfs file the container reads at runtime.

Disable by adding `ZC1572` to `disabled_katas` in `.zshellcheckrc`.

---

<a id="zc1573"></a>
### ZC1573 — Warn on `chattr -i` / `chattr -a` — removes immutable / append-only attribute

**Severity:** `warning`  
**Auto-fix:** `no`

Removing the immutable (`-i`) or append-only (`-a`) attribute lets the file be overwritten or truncated again. When the target is a log file, shadow file, or hardened system binary, that flag was explicitly set to make tampering noisy. Removing it mid-script is either a one-shot upgrade (follow with the `chattr +i` restore) or an anti-forensics step. If it is the former, wrap the change in a function and re-set the attribute at the end.

Disable by adding `ZC1573` to `disabled_katas` in `.zshellcheckrc`.

---

<a id="zc1574"></a>
### ZC1574 — Warn on `git config credential.helper store` — plaintext credentials on disk

**Severity:** `warning`  
**Auto-fix:** `no`

`credential.helper store` writes the username and password to `~/.git-credentials` in plaintext. Anything that backs up that file (rsync, imaging, cloud sync) then carries the credential around. Use a platform helper instead: `manager` / `manager-core` on Windows / Mac, `libsecret` on Linux, or `cache --timeout=3600` for short-lived in-memory caching.

Disable by adding `ZC1574` to `disabled_katas` in `.zshellcheckrc`.

---

<a id="zc1575"></a>
### ZC1575 — Error on `aws configure set aws_secret_access_key <value>` — secret on cmdline

**Severity:** `error`  
**Auto-fix:** `no`

`aws configure set aws_secret_access_key …` writes the secret access key into `~/.aws/credentials` and leaves the raw value in `ps` / shell history until the process exits. On a shared CI runner or a multi-user host, that window is long enough for a co-tenant to snapshot the key. Use IAM-role-based auth (EC2 instance profile, IRSA on EKS, OIDC from GitHub / GitLab) or read the value from stdin / a 0600 file and let `aws configure` import it interactively.

Disable by adding `ZC1575` to `disabled_katas` in `.zshellcheckrc`.

---

<a id="zc1576"></a>
### ZC1576 — Warn on `terraform apply -target=...` — cherry-pick apply bypasses dependencies

**Severity:** `warning`  
**Auto-fix:** `no`

`-target` restricts `terraform apply` to a specific resource / module and everything it depends on. In theory that is a surgical fix; in practice it routinely skips changes the targeted resource actually depends on, leading to drift between state and configuration. HashiCorp documents `-target` as a tool for incident response, not routine operations. Re-run without `-target` or split the configuration into separate root modules.

Disable by adding `ZC1576` to `disabled_katas` in `.zshellcheckrc`.

---

<a id="zc1577"></a>
### ZC1577 — Warn on `dig <name> ANY` — deprecated query type (RFC 8482)

**Severity:** `warning`  
**Auto-fix:** `no`

ANY queries return whatever the authoritative server feels like sending back — or just the HINFO placeholder mandated by RFC 8482. Modern recursors filter ANY to avoid reflection-amplification abuse, so scripts that rely on ANY for enumeration get inconsistent or empty results. Query the specific record types you want (`dig A name`, `dig MX name`, `dig NS name`) and combine them.

Disable by adding `ZC1577` to `disabled_katas` in `.zshellcheckrc`.

---

<a id="zc1578"></a>
### ZC1578 — Warn on `ssh-keygen -b <2048` for RSA / DSA — weak SSH key

**Severity:** `warning`  
**Auto-fix:** `no`

Generating an SSH RSA or DSA key shorter than 2048 bits fails current OpenSSH baselines and is rejected by recent `ssh` versions when used for authentication. DSA was removed from OpenSSH 9.8 outright. Use `ssh-keygen -t ed25519` (compact, fast, modern defaults) or `ssh-keygen -t rsa -b 4096` if you need RSA for compatibility.

Disable by adding `ZC1578` to `disabled_katas` in `.zshellcheckrc`.

---

<a id="zc1579"></a>
### ZC1579 — Warn on `curl --retry-all-errors` without `--max-time` — hammers endpoint on failure

**Severity:** `warning`  
**Auto-fix:** `no`

`--retry-all-errors` (curl 7.71+) treats every HTTP error as retryable. Without `--max-time` capping total wall clock, a server that responds `500` quickly gets hit back-to-back until `--retry` exhausts — a mini-DoS against your own upstream, especially if the script itself is scheduled on many nodes. Pair with `--max-time <seconds>` or prefer `--retry-connrefused` (only retries connection-level failures).

Disable by adding `ZC1579` to `disabled_katas` in `.zshellcheckrc`.

---

<a id="zc1580"></a>
### ZC1580 — Warn on `go build -ldflags "-X main.<SECRET>=..."` — secret embedded in binary

**Severity:** `warning`  
**Auto-fix:** `no`

`-ldflags="-X pkg.Var=value"` sets a Go string variable at link time. Putting a secret here bakes it into the resulting binary (discoverable with `strings`, `objdump`, or simply opening the file). It also leaves the value on the build host's shell history and in any CI transcript. Read the value at runtime from `os.Getenv` / a mounted secret file / the cloud secret manager.

Disable by adding `ZC1580` to `disabled_katas` in `.zshellcheckrc`.

---

<a id="zc1581"></a>
### ZC1581 — Warn on `ssh -o PubkeyAuthentication=no` / `-o PasswordAuthentication=yes`

**Severity:** `warning`  
**Auto-fix:** `no`

Forcing password authentication on a connection that has a working key turns a strong (challenge-response, no password leaves the client) into a weak (password-in-the-clear-on-disk-or-prompt) authentication path. Similarly disabling pubkey skips the good path entirely. Leave the defaults, let the server's `PubkeyAuthentication yes` pick the key, and document any exception.

Disable by adding `ZC1581` to `disabled_katas` in `.zshellcheckrc`.

---

<a id="zc1582"></a>
### ZC1582 — Warn on `bash -x` / `sh -x` / `zsh -x` — traces every command, leaks secrets

**Severity:** `warning`  
**Auto-fix:** `no`

`-x` turns on xtrace, printing every command (expanded) to stderr before it runs. In a CI log that is indexed / shared / archived, any line that touches a secret leaks it verbatim — `curl` with a `Bearer` header, `psql` with a password, `echo $API_TOKEN > ...`. If you really need tracing, wrap the non-secret block with `set -x; ...; set +x` and exclude the secret-handling parts.

Disable by adding `ZC1582` to `disabled_katas` in `.zshellcheckrc`.

---

<a id="zc1583"></a>
### ZC1583 — Warn on `find ... -delete` without `-maxdepth` — unbounded recursive delete

**Severity:** `warning`  
**Auto-fix:** `no`

`find PATH -delete` walks the tree recursively and removes every match. Without `-maxdepth N` the walk crosses into every subtree, including symlinks that point outside the intended scope and mount points that expand the blast radius. Scope the depth (`-maxdepth 2`) and prefer a dry-run first (`find ... -print | head`).

Disable by adding `ZC1583` to `disabled_katas` in `.zshellcheckrc`.

---

<a id="zc1584"></a>
### ZC1584 — Warn on `sudo -E` / `--preserve-env` — carries caller env into root shell

**Severity:** `warning`  
**Auto-fix:** `no`

`sudo -E` preserves the invoking user's environment — `PATH`, `LD_PRELOAD`, `PYTHONPATH`, etc. On a workstation where the user has a personal `~/bin` early in `$PATH`, any wrapper named like a system binary gets executed by the privileged process. That is exactly the sudoers `secure_path` mechanic fails to protect against. Whitelist specific variables with `env_keep` in sudoers, or call `sudo env VAR=value cmd` with the minimum.

Disable by adding `ZC1584` to `disabled_katas` in `.zshellcheckrc`.

---

<a id="zc1585"></a>
### ZC1585 — Warn on `ufw allow from any` / `firewall-cmd --add-source=0.0.0.0/0`

**Severity:** `warning`  
**Auto-fix:** `no`

`ufw allow from any to any port …` (and its firewall-cmd sibling `--add-source=0.0.0.0/0`) opens the port to the whole internet. That is sometimes the point (public HTTP / HTTPS), but on management ports (22, 3306, 5432, 6379, 9200, 27017) it is a routine foot-gun when the script author assumed the host would only ever be reached via VPN. Scope the rule to a specific source CIDR.

Disable by adding `ZC1585` to `disabled_katas` in `.zshellcheckrc`.

---

<a id="zc1586"></a>
### ZC1586 — Style: `chkconfig` / `update-rc.d` / `insserv` — SysV init relics, use `systemctl`

**Severity:** `style`  
**Auto-fix:** `no`

`chkconfig` (Red Hat), `update-rc.d` (Debian), and `insserv` (SUSE) are SysV-init compatibility wrappers for enabling/disabling services at boot. On any distro that has used systemd for the last decade they are translated to `systemctl enable|disable`, but silently lose unit-template arguments, `[Install]` alias handling, and socket-activated services. Call `systemctl enable <unit>` directly.

Disable by adding `ZC1586` to `disabled_katas` in `.zshellcheckrc`.

---

<a id="zc1587"></a>
### ZC1587 — Warn on `modprobe -r` / `rmmod` from scripts — unloading active kernel modules

**Severity:** `warning`  
**Auto-fix:** `no`

Unloading a kernel module that is in use — `nvme` (storage), `nvidia` (GPU), `e1000`/`ixgbe` (network), `kvm` (virt) — instantly takes the backing subsystem offline. On a remote host the script loses its storage or network mid-run. Reserve `modprobe -r` / `rmmod` for console maintenance, and consider `systemctl stop <unit>` if you are trying to stop a service.

Disable by adding `ZC1587` to `disabled_katas` in `.zshellcheckrc`.

---

<a id="zc1588"></a>
### ZC1588 — Error on `nsenter --target 1` — joins host init namespaces (container escape)

**Severity:** `error`  
**Auto-fix:** `no`

`nsenter -t 1` attaches to the namespaces of pid 1. Inside a privileged container or one with `CAP_SYS_ADMIN`, pid 1 is the host init — joining its mount / pid / net / uts / ipc namespaces is the canonical escape primitive. From that new shell the caller sees and writes the host filesystem, kills host processes, and hijacks host network. Legit debugging runs from the host, not from inside the container. If you need to exec into a container, use `docker exec` / `kubectl exec`.

Disable by adding `ZC1588` to `disabled_katas` in `.zshellcheckrc`.

---

<a id="zc1589"></a>
### ZC1589 — Warn on `trap 'set -x' ERR/RETURN/EXIT/ZERR` — trace hook leaks env to stderr

**Severity:** `warning`  
**Auto-fix:** `no`

Installing a trap that enables `set -x` (or `set -o xtrace` / `set -v`) causes every subsequent expanded command to hit stderr. Expansions embed environment variables — API tokens, passwords, signed URLs — directly into the trace. In CI, that stderr lands in build logs and gets shipped to long-term log retention. Scope `set -x` to a `set -x ... set +x` block around the suspect code, or replace the trap with `trap 'safe_dump' ERR` that prints only non-sensitive state.

Disable by adding `ZC1589` to `disabled_katas` in `.zshellcheckrc`.

---

<a id="zc1590"></a>
### ZC1590 — Error on `sshpass -p SECRET` — password in process list and history

**Severity:** `error`  
**Auto-fix:** `no`

`sshpass -p SECRET` places the password in argv. It leaks into `ps`, `/proc/<pid>/cmdline`, shell history, and audit logs for every process on the box that can list processes. The `-f FILE` and `-e` (SSHPASS env) variants keep it off argv, but key-based auth is the real fix. Generate an SSH key, authorize it on the remote, and drop the password tool entirely.

Disable by adding `ZC1590` to `disabled_katas` in `.zshellcheckrc`.

---

<a id="zc1591"></a>
### ZC1591 — Style: use Zsh `print -l` / `${(F)array}` instead of `printf '%s\n' "${array[@]}"`

**Severity:** `style`  
**Auto-fix:** `no`

`printf '%s\n' "${array[@]}"` is the Bash-idiomatic way to print one element per line. Zsh has `print -l -r -- "${array[@]}"` (one element per line, raw, sentinel-safe) and the parameter-expansion flag `${(F)array}` (newline-join, fine for `$(...)`). Both are shorter than the printf incantation and avoid format-string surprises if the array ever contains a literal `%`.

Disable by adding `ZC1591` to `disabled_katas` in `.zshellcheckrc`.

---

<a id="zc1592"></a>
### ZC1592 — Warn on `faillock --reset` / `pam_tally2 -r` — clears failed-auth counter

**Severity:** `warning`  
**Auto-fix:** `no`

Both tools zero the PAM counter that triggers account lockout after too many failed logins. A script that resets lockouts — even legitimately, to recover locked users — also erases evidence of an ongoing brute-force attempt. Intrusion detection relies on those counters for alerting. Do not automate resets; if you must, log the prior count and page security on every invocation.

Disable by adding `ZC1592` to `disabled_katas` in `.zshellcheckrc`.

---

<a id="zc1593"></a>
### ZC1593 — Error on `blkdiscard` — issues TRIM/DISCARD across the whole device (data loss)

**Severity:** `error`  
**Auto-fix:** `no`

`blkdiscard $DEV` tells the underlying SSD controller to invalidate every block in the range. On most modern drives the data is unrecoverable the moment the controller acknowledges — even forensic recovery cannot pull it back. Scripts that reach this command from any codepath an attacker or typo can trigger destroy the drive. Gate it behind interactive confirmation, not shell flow control.

Disable by adding `ZC1593` to `disabled_katas` in `.zshellcheckrc`.

---

<a id="zc1594"></a>
### ZC1594 — Warn on `docker/podman run --security-opt=systempaths=unconfined` — unhides host kernel knobs

**Severity:** `warning`  
**Auto-fix:** `no`

`systempaths=unconfined` removes the container runtime's masking of `/proc/sys`, `/proc/sysrq-trigger`, `/sys/firmware`, and related kernel surfaces. Without the default shield a compromised process inside the container can write `/proc/sysrq-trigger` to panic the host, or edit `/proc/sys/kernel/*` to change kernel policy on the fly. Keep the default `systempaths=all` (masked) unless you have a specific kernel tunable you need, then mount only that path.

Disable by adding `ZC1594` to `disabled_katas` in `.zshellcheckrc`.

---

<a id="zc1595"></a>
### ZC1595 — Warn on `setfacl -m u:nobody:... / o::rwx` — ACL grants that bypass `chmod` scrutiny

**Severity:** `warning`  
**Auto-fix:** `no`

Filesystem ACLs live outside the mode bits that `chmod` / `ls -l` / `stat -c %a` surface. Granting `u:nobody:rwx` gives the daemon-fallback account write access to a file; `o::rwx` / `o::rw` world-writes via ACL even when the mode bits still look safe. Review scripts that check `stat -c %a` miss both. Prefer `chmod` for world perms, and for specific users name the real account with the minimum perm set.

Disable by adding `ZC1595` to `disabled_katas` in `.zshellcheckrc`.

---

<a id="zc1596"></a>
### ZC1596 — Style: `emulate sh/bash/ksh` without `-L` — flips options for the whole shell

**Severity:** `style`  
**Auto-fix:** `no`

`emulate MODE` without the `-L` flag changes Zsh options globally. After that line runs the shell is no longer in Zsh mode — `${(F)arr}`, 1-indexed arrays, glob qualifiers, and other Zsh-only constructs either error or silently behave differently. Wrap emulation in a function and use `emulate -L MODE` to scope it to that function. A `.zsh` script that starts with `emulate sh` likely belongs in a `.sh` file instead.

Disable by adding `ZC1596` to `disabled_katas` in `.zshellcheckrc`.

---

<a id="zc1597"></a>
### ZC1597 — Warn on `systemd-run -p User=root` — launches arbitrary command with root privileges

**Severity:** `warning`  
**Auto-fix:** `no`

`systemd-run` submits a transient unit to systemd. With `-p User=root` (or `User=0`) the unit runs as root — bypassing the usual `sudo` audit path in `/var/log/auth.log`. On hosts where the caller's polkit / dbus rules allow the operation, this is effectively privilege escalation by a different name. Prefer explicit `sudo` so the invocation is logged, or pre-provision a dedicated systemd unit that names the exact command it can run.

Disable by adding `ZC1597` to `disabled_katas` in `.zshellcheckrc`.

---

<a id="zc1598"></a>
### ZC1598 — Error on `chmod` with world-write bit on a sensitive `/dev/` node

**Severity:** `error`  
**Auto-fix:** `no`

Device nodes under `/dev/` are kernel interfaces. Making one world-writable ( last digit `2`, `3`, `6`, or `7` ) gives every local user a direct line into the kernel — `/dev/kvm` yields VM hypercalls, `/dev/mem` / `/dev/kmem` / `/dev/port` read and write physical memory, `/dev/sd*` and `/dev/nvme*` give raw block access, `/dev/input/*` sniffs keystrokes. Keep restrictive perms (600 / 660) and use udev rules (`GROUP=`, `MODE=`) to grant access declaratively.

Disable by adding `ZC1598` to `disabled_katas` in `.zshellcheckrc`.

---

<a id="zc1599"></a>
### ZC1599 — Warn on `ldconfig -f PATH` outside `/etc/` — attacker-writable loader cache

**Severity:** `warning`  
**Auto-fix:** `no`

`ldconfig -f PATH` rebuilds `/etc/ld.so.cache` using PATH instead of the system `/etc/ld.so.conf`. If PATH sits in `/tmp`, `/var/tmp`, `$HOME`, or any directory an attacker can create, they can inject an `include` line that points at their directory of malicious shared objects. After the cache rebuild, every subsequent executable on the host loads their library first. Keep the config under `/etc/ld.so.conf.d/` with root ownership and run `ldconfig` with no `-f`.

Disable by adding `ZC1599` to `disabled_katas` in `.zshellcheckrc`.

---

<a id="zc1600"></a>
### ZC1600 — Warn on bare `chroot DIR CMD` — missing `--userspec=` keeps uid 0 inside the jail

**Severity:** `warning`  
**Auto-fix:** `no`

`chroot` changes the filesystem root but does not drop privileges. The caller is almost always root (the syscall needs `CAP_SYS_CHROOT`), and without `--userspec=USER:GROUP` the command inside the chroot still runs as uid 0. It can write anywhere inside the tree, chmod binaries, and — if proc / sys / device nodes are bind-mounted in — escape. Pass `--userspec=` to run the command as a named unprivileged user, or drop to a dedicated helper (bubblewrap, firejail) that also unshares user namespaces.

Disable by adding `ZC1600` to `disabled_katas` in `.zshellcheckrc`.

---

<a id="zc1601"></a>
### ZC1601 — Warn on `ethtool -s $IF wol <g|u|m|b|a>` — enables remote Wake-on-LAN

**Severity:** `warning`  
**Auto-fix:** `no`

Wake-on-LAN powers the host on from a sleep / soft-off state when a matching packet reaches the NIC. The wake logic fires in a privileged firmware path long before the kernel boots and firewall rules are loaded — so any packet that reaches the interface (magic-packet, unicast, broadcast, ARP) triggers the power-on unfiltered. On a shared or public LAN attackers on the broadcast domain can wake hosts at will. Keep `wol d` (disable) unless a documented operational need requires one of the wake bits.

Disable by adding `ZC1601` to `disabled_katas` in `.zshellcheckrc`.

---

<a id="zc1602"></a>
### ZC1602 — Warn on `setopt KSH_ARRAYS` / `SH_WORD_SPLIT` — flips Zsh core semantics shell-wide

**Severity:** `warning`  
**Auto-fix:** `no`

`KSH_ARRAYS` makes arrays 0-indexed (the Bash / ksh convention), breaking every Zsh access that uses `[1]` for the first element. `SH_WORD_SPLIT` makes unquoted `$var` word-split on `IFS`, breaking the core Zsh promise that `echo $x` passes exactly one argument. Setting either globally is a bug-magnet — pre-existing code silently misbehaves from that line on. If you need the semantics only inside a function, scope it with `emulate -L ksh` or `emulate -L sh`.

Disable by adding `ZC1602` to `disabled_katas` in `.zshellcheckrc`.

---

<a id="zc1603"></a>
### ZC1603 — Warn on `gdb -p PID` / `ltrace -p PID` — live attach reads target memory

**Severity:** `warning`  
**Auto-fix:** `no`

`gdb -p PID` and `ltrace -p PID` attach via ptrace and hand the caller full read / write access to the target process: registers, heap, stack, open file descriptors, and every environment variable. Credentials in `$AWS_SECRET_ACCESS_KEY`, session tokens on the stack, TLS keys in memory — all readable. A root-run script that attaches to another user's process extracts everything that user has. Keep production scripts out of the debugger; if post-mortem diagnostics are needed, use `coredumpctl` against a captured core file instead.

Disable by adding `ZC1603` to `disabled_katas` in `.zshellcheckrc`.

---

<a id="zc1604"></a>
### ZC1604 — Warn on `source <glob>` / `. <glob>` — loads every match; one bad file = code exec

**Severity:** `warning`  
**Auto-fix:** `no`

`source /etc/profile.d/*.sh` and similar glob-sourcing patterns load every file that matches, in the order Zsh enumerates them. One attacker-writable file anywhere in the glob yields arbitrary code execution as whoever is running the script, with that caller's privileges. Prefer explicit filenames so review can enumerate exactly what gets loaded. If a directory of drop-ins is required, audit ownership and perms at install time and keep the directory root-owned.

Disable by adding `ZC1604` to `disabled_katas` in `.zshellcheckrc`.

---

<a id="zc1605"></a>
### ZC1605 — Error on `debugfs -w DEV` — write-mode filesystem debugger bypasses journal

**Severity:** `error`  
**Auto-fix:** `no`

`debugfs -w` opens the filesystem in write mode. It sidesteps the kernel's normal write path — the journal doesn't see the changes, filesystem locks are ignored, and inodes / blocks can be edited directly. On a mounted filesystem this corrupts state silently; even on an unmounted one, the operator can repoint a directory entry at an arbitrary inode. Scripts should never need this — keep `debugfs -w` as an interactive last-resort from a rescue environment.

Disable by adding `ZC1605` to `disabled_katas` in `.zshellcheckrc`.

---

<a id="zc1606"></a>
### ZC1606 — Warn on `mkdir -m NNN` / `install -m NNN` with world-write bit (no sticky)

**Severity:** `warning`  
**Auto-fix:** `no`

`mkdir -m 777 /path` and `install -m 777 src /dest` create a path that every local user can write and rename inside. If the script later creates files there, classic TOCTOU symlink attacks become trivial — the attacker drops a symlink named like the expected output file, redirecting the write wherever they choose. A sticky-bit mode (`1777`) mitigates this for shared temp dirs. Prefer `mkdir -m 700` (or 750), and scope access by group or ACL rather than everyone.

Disable by adding `ZC1606` to `disabled_katas` in `.zshellcheckrc`.

---

<a id="zc1607"></a>
### ZC1607 — Warn on `git config safe.directory '*'` — disables CVE-2022-24765 protection

**Severity:** `warning`  
**Auto-fix:** `no`

`safe.directory` is git's mitigation for CVE-2022-24765 (fake git dirs planted by another uid). Setting it to `'*'` trusts every directory on the host — an attacker who creates `/tmp/evil/.git` with a malicious `core.fsmonitor` hook gets arbitrary code execution the first time any user runs `git status` near that path. List the specific paths that need cross-owner git access instead, or fix the underlying ownership mismatch.

Disable by adding `ZC1607` to `disabled_katas` in `.zshellcheckrc`.

---

<a id="zc1608"></a>
### ZC1608 — Warn on `find -exec sh -c '... {} ...'` — filename in quoted script is injectable

**Severity:** `warning`  
**Auto-fix:** `no`

Substituting `{}` directly into the quoted command string of `find -exec sh -c` lets filenames with shell metacharacters break out. A file named `$(rm -rf ~)` invokes command substitution; a file named `foo; curl evil` chains a second command. Pass `{}` as a positional argument to `sh` so the filename arrives as a parameter, not as source: `find -exec sh -c 'grep pat "$1"' _ {} \;`.

Disable by adding `ZC1608` to `disabled_katas` in `.zshellcheckrc`.

---

<a id="zc1609"></a>
### ZC1609 — Warn on `aa-disable` / `aa-complain` / `apparmor_parser -R` — disables AppArmor enforcement

**Severity:** `warning`  
**Auto-fix:** `no`

`aa-disable` fully unloads the named AppArmor profile; `aa-complain` flips the profile from enforce to complain (violations are logged but allowed); `apparmor_parser -R` removes a profile from the running kernel. Each one lets the confined process run without its mandatory-access-control restrictions — if the profile existed for a reason, that reason is now unenforced. Interactive debugging is legitimate, but scripts that permanently disable profiles should be reviewed.

Disable by adding `ZC1609` to `disabled_katas` in `.zshellcheckrc`.

---

<a id="zc1610"></a>
### ZC1610 — Warn on `curl -o /etc/...` / `wget -O /etc/...` — direct download to a system path

**Severity:** `warning`  
**Auto-fix:** `no`

Writing the body of an HTTP response straight into `/etc/`, `/usr/`, `/bin/`, `/sbin/`, or `/lib/` skips every integrity check the system usually applies. If the URL is compromised or MITM'd, the attacker's content replaces a system config or binary the next command over. Download to a temp file, verify signature / checksum, and `install -m 0644` the final file into place. Package managers exist for a reason — prefer them for system files.

Disable by adding `ZC1610` to `disabled_katas` in `.zshellcheckrc`.

---

<a id="zc1611"></a>
### ZC1611 — Style: `${var^^}` / `${var,,}` — prefer Zsh `${(U)var}` / `${(L)var}` for case change

**Severity:** `style`  
**Auto-fix:** `no`

`${var^^}` (uppercase) and `${var,,}` (lowercase) came from Bash 4. Zsh accepts them for compatibility but the idiomatic form is the parameter-expansion flag: `${(U)var}` / `${(L)var}`. The flag is also available per-element in arrays (`${(U)array}`) and composes with other flags (`${(UL)array}` doesn't make sense, but `${(U)${(f)str}}` does). Prefer the Zsh-native form in a `.zsh` script; it keeps the codebase consistent with other `(X)var` patterns.

Disable by adding `ZC1611` to `disabled_katas` in `.zshellcheckrc`.

---

<a id="zc1612"></a>
### ZC1612 — Warn on `sysctl -w` disabling kernel hardening knobs

**Severity:** `warning`  
**Auto-fix:** `no`

Several sysctl knobs exist specifically to constrain what unprivileged users can do — `kernel.yama.ptrace_scope`, `kernel.kptr_restrict`, `kernel.dmesg_restrict`, `kernel.unprivileged_bpf_disabled`, `net.core.bpf_jit_harden`, and `kernel.perf_event_paranoid`. Setting any of them to the lowest-restriction value removes a distinct defense-in-depth layer: unrelated processes can ptrace each other, kernel pointers leak to `/proc`, unprivileged users read kernel ring buffers, BPF JIT-spray mitigations disappear. Leave these defaults alone unless a measured performance or debugging need justifies it.

Disable by adding `ZC1612` to `disabled_katas` in `.zshellcheckrc`.

---

<a id="zc1613"></a>
### ZC1613 — Warn on reading SSH private-key files with `cat` / `less` / `grep` / `head`

**Severity:** `warning`  
**Auto-fix:** `no`

Piping an SSH private key through a generic text tool copies the raw key material into the process and — if stdout is redirected or piped — often into logs, backup files, or a terminal scrollback buffer. Host keys under `/etc/ssh/ssh_host_*_key` impersonate the server; user keys under `~/.ssh/id_*` impersonate the user. Use `ssh-keygen -l -f KEY` for fingerprint / metadata, or pass the key path to the consumer directly (`ssh -i`, `git -c core.sshCommand`) without staging it through a shell tool.

Disable by adding `ZC1613` to `disabled_katas` in `.zshellcheckrc`.

---

<a id="zc1614"></a>
### ZC1614 — Error on `expect` script containing `password` / `passphrase`

**Severity:** `error`  
**Auto-fix:** `no`

`expect -c '... password ... send "..."'` puts the entire scripted dialog on the command line. Anything there — including the password or passphrase — is visible in `ps`, `/proc/<pid>/cmdline`, shell history, and audit logs. Use key-based authentication (SSH keys, GSSAPI) where possible. If password feeding is truly unavoidable, read it from a protected file with `spawn -o`, or source it from an environment variable the script does not print.

Disable by adding `ZC1614` to `disabled_katas` in `.zshellcheckrc`.

---

<a id="zc1615"></a>
### ZC1615 — Style: use Zsh `$EPOCHREALTIME` / `$epochtime` instead of `date "+%s.%N"`

**Severity:** `style`  
**Auto-fix:** `no`

Zsh's `zsh/datetime` module exposes `$EPOCHREALTIME` (scalar with fractional seconds) and `$epochtime` (two-element array of seconds and nanoseconds). Both read straight from `clock_gettime(CLOCK_REALTIME)` without forking `date`. On a hot path the builtin is dramatically faster and avoids subshell process-startup overhead. Autoload the module once with `zmodload zsh/datetime`.

Disable by adding `ZC1615` to `disabled_katas` in `.zshellcheckrc`.

---

<a id="zc1616"></a>
### ZC1616 — Warn on `fsfreeze -f MOUNTPOINT` — filesystem stays frozen until `-u` runs

**Severity:** `warning`  
**Auto-fix:** `no`

`fsfreeze -f` blocks every write on the mountpoint until `fsfreeze -u` thaws it. The intended use is a short window around a hypervisor or LVM snapshot. If the script errors between the freeze and the unfreeze (or is killed), the filesystem stays frozen — every subsequent write hangs forever until the admin manually thaws it, and a reboot may be the only way out on the root fs. Pair every freeze with `trap 'fsfreeze -u MOUNTPOINT' EXIT` and keep the window under a few seconds.

Disable by adding `ZC1616` to `disabled_katas` in `.zshellcheckrc`.

---

<a id="zc1617"></a>
### ZC1617 — Warn on `xargs -P 0` — unbounded parallelism risks CPU / fd / memory exhaustion

**Severity:** `warning`  
**Auto-fix:** `no`

`xargs -P 0` tells xargs to spawn as many concurrent children as input lines. On any non-trivial input that number can blow past `RLIMIT_NPROC`, saturate the downstream tool's file-descriptor limit, or drive the host OOM. Pick an explicit cap — `xargs -P $(nproc)` for CPU-bound work, `-P 4..8` for I/O-bound — so the failure mode is bounded and predictable.

Disable by adding `ZC1617` to `disabled_katas` in `.zshellcheckrc`.

---

<a id="zc1618"></a>
### ZC1618 — Warn on `git commit --no-verify` / `git push --no-verify` — bypasses hooks

**Severity:** `warning`  
**Auto-fix:** `no`

`--no-verify` skips pre-commit, commit-msg, and pre-push hooks. Those hooks are where projects run linting, type-checking, unit tests, and secret scanning before code lands. A commit or push with `--no-verify` ships code the project's own automation would have rejected. Reserve the flag for emergencies with a follow-up commit that passes the hooks; scripts should not use it routinely.

Disable by adding `ZC1618` to `disabled_katas` in `.zshellcheckrc`.

---

<a id="zc1619"></a>
### ZC1619 — Warn on `mount -t nfs/cifs/smb/sshfs` missing `nosuid` or `nodev`

**Severity:** `warning`  
**Auto-fix:** `no`

Network filesystems present files whose mode bits are controlled by a remote server. Without `nosuid` in the mount options, a compromised or hostile server can plant a setuid-root binary on the share; the client kernel honors the suid bit and the binary runs as root on the mounting host. Without `nodev`, the server can plant device nodes the kernel treats as real. Always mount network shares with `nosuid,nodev`; add `noexec` unless the export is intended to hold executables.

Disable by adding `ZC1619` to `disabled_katas` in `.zshellcheckrc`.

---

<a id="zc1620"></a>
### ZC1620 — Error on `tee /etc/sudoers` / `/etc/sudoers.d/*` — writes without `visudo -cf`

**Severity:** `error`  
**Auto-fix:** `no`

`tee` copies stdin to the file with no syntax check. A typo in a sudoers rule — a stray comma, a missing `ALL`, an unclosed alias — leaves the file unparseable. The next sudo call refuses to load it and on most systems nobody can become root until someone boots from rescue media. Pipe the content through `visudo -cf /dev/stdin` first, or write to a temp file, validate with `visudo -cf`, then atomically `mv` into `/etc/sudoers.d/`.

Disable by adding `ZC1620` to `disabled_katas` in `.zshellcheckrc`.

---

<a id="zc1621"></a>
### ZC1621 — Warn on `tmux -S /tmp/SOCKET` — shared-path socket invites session hijack

**Severity:** `warning`  
**Auto-fix:** `no`

`tmux -S PATH` overrides the default socket location (normally under `$XDG_RUNTIME_DIR/tmux-$UID/`, a 0700-mode directory). Paths under `/tmp/` or `/var/tmp/` are world-traversable; if the socket is created with loose permissions, any local user who can read it can `tmux -S /tmp/PATH attach` and see / drive the session — keystrokes, output, arbitrary commands in the attached pane. Keep the socket in `$XDG_RUNTIME_DIR` or another 0700-scoped directory.

Disable by adding `ZC1621` to `disabled_katas` in `.zshellcheckrc`.

---

<a id="zc1622"></a>
### ZC1622 — Style: `${var@U/L/Q/...}` — prefer Zsh `${(U)var}` / `${(L)var}` / `${(Q)var}` flags

**Severity:** `style`  
**Auto-fix:** `no`

The `@<op>` suffix came from Bash 5. Zsh 5.9+ compiles in compatibility for the common ones, but the idiomatic Zsh form is the `(X)var` parameter-expansion flag — `${(U)var}` uppercase, `${(L)var}` lowercase, `${(Q)var}` unquote, `${(k)var}` keys, `${(t)var}` type, `${(e)var}` re-evaluate. The flag form composes (`${(Uf)str}` works) and reads consistently across the Zsh documentation. Prefer the native flag over the Bash-compat form.

Disable by adding `ZC1622` to `disabled_katas` in `.zshellcheckrc`.

---

<a id="zc1623"></a>
### ZC1623 — Warn on `kill -STOP PID` / `pkill -STOP` — target halts until `kill -CONT` runs

**Severity:** `warning`  
**Auto-fix:** `no`

Sending SIGSTOP halts the target process until SIGCONT arrives. If the script fails, is killed, or exits before the resume, the target stays paused indefinitely — consuming memory, holding locks, blocking its dependents. Wrap every `kill -STOP $PID` with `trap "kill -CONT $PID" EXIT` (or an explicit cleanup path) so the resume fires even on failure. Prefer `kill -TSTP` if the target can handle it (the user-space tstop that the process can ignore).

Disable by adding `ZC1623` to `disabled_katas` in `.zshellcheckrc`.

---

<a id="zc1624"></a>
### ZC1624 — Error on `az login -p` / `--password` — service-principal secret in process list

**Severity:** `error`  
**Auto-fix:** `no`

`az login -p SECRET` passes the service-principal password as an argv element. The expanded value shows up in `ps`, `/proc/<pid>/cmdline`, shell history, and audit logs — readable by any local user who can list processes. Prefer federated-token OIDC (`--federated-token`), managed identity on the host, or interactive device-code flow. If a password is unavoidable, export it as `AZURE_PASSWORD` via a protected env var and call plain `az login --service-principal` (which reads from env).

Disable by adding `ZC1624` to `disabled_katas` in `.zshellcheckrc`.

---

<a id="zc1625"></a>
### ZC1625 — Error on `rm --no-preserve-root` — disables GNU rm safeguard against `rm -rf /`

**Severity:** `error`  
**Auto-fix:** `no`

GNU `rm` refuses to remove `/` by default — the `--preserve-root` safeguard added in coreutils 8.4. `--no-preserve-root` explicitly disables that check so `rm -rf /` actually recurses and wipes the filesystem. Scripts that pass the flag are asking `rm` to go ahead if the argument happens to evaluate to `/`. Remove the flag; if a specific path genuinely needs deletion, list it explicitly and leave the safeguard in place.

Disable by adding `ZC1625` to `disabled_katas` in `.zshellcheckrc`.

---

<a id="zc1626"></a>
### ZC1626 — Error on `helm install/upgrade --set KEY=VALUE` with secret-shaped key

**Severity:** `error`  
**Auto-fix:** `no`

`--set` and `--set-string` put the full `KEY=VALUE` pair on the helm command line. When the key name looks like a secret (`password`, `secret`, `token`, `apikey`, `access_key`, `private_key`), the expanded VALUE appears in `ps`, `/proc/<pid>/cmdline`, shell history, and audit logs — readable by any local user who can list processes. Put secrets in a protected values file (`helm install -f /secure/values.yaml`), or use `--set-file KEY=PATH` so helm reads the content from PATH at apply time.

Disable by adding `ZC1626` to `disabled_katas` in `.zshellcheckrc`.

---

<a id="zc1627"></a>
### ZC1627 — Warn on `crontab /tmp/FILE` — attacker-writable path installed as a user's cron

**Severity:** `warning`  
**Auto-fix:** `no`

`crontab PATH` replaces the user's cron with whatever PATH currently contains. A path under `/tmp/` or `/var/tmp/` is world-traversable; a concurrent local user can replace the file between the moment the script writes it and the moment `crontab` reads it, substituting their own cron rules. Keep the staging file in a 0700-scoped directory (e.g. `$XDG_RUNTIME_DIR/` or `mktemp -d`), or pipe the content via `crontab -` after generating it in-memory.

Disable by adding `ZC1627` to `disabled_katas` in `.zshellcheckrc`.

---

<a id="zc1628"></a>
### ZC1628 — Warn on `insmod` / `modprobe -f` — loads modules bypassing blacklist / signature checks

**Severity:** `warning`  
**Auto-fix:** `no`

`insmod PATH.ko` loads a kernel module from a file, skipping the depmod-built dependency graph and the `/etc/modprobe.d/*.conf` blacklist. `modprobe -f` instructs modprobe to ignore version-magic and kernel-mismatch checks. Either path lets a module enter the kernel that the administrator explicitly disabled, or one compiled against a different kernel — crash, privesc, or full kernel compromise. Use plain `modprobe MODNAME` so the system's policy and signature verification run.

Disable by adding `ZC1628` to `disabled_katas` in `.zshellcheckrc`.

---

<a id="zc1629"></a>
### ZC1629 — Warn on `rsync --rsync-path='sudo rsync'` — hidden remote privilege escalation

**Severity:** `warning`  
**Auto-fix:** `no`

`--rsync-path` normally overrides the path to the remote rsync binary. Setting it to `sudo rsync` (or `doas rsync` / `pkexec rsync`) instead makes the remote side run rsync as root. That is sometimes legitimate — copying into `/etc/` from a CI job — but the flag is easy to miss in review because it looks like a path override. Provision a scoped sudoers rule that names exactly which rsync invocation the remote user may run, and keep the path explicit (`--rsync-path=/usr/bin/rsync`).

Disable by adding `ZC1629` to `disabled_katas` in `.zshellcheckrc`.

---

<a id="zc1630"></a>
### ZC1630 — Warn on `php -S 0.0.0.0:PORT` — PHP dev server exposes CWD to all interfaces

**Severity:** `warning`  
**Auto-fix:** `no`

`php -S 0.0.0.0:PORT` starts PHP's built-in dev server listening on every interface the host has. It serves files from the working directory (or the docroot named after the bind) with no auth, no TLS, and minimal access logging. The PHP docs explicitly say not to use it in production. Bind to `127.0.0.1:PORT` for local testing and put nginx / caddy in front for anything externally exposed.

Disable by adding `ZC1630` to `disabled_katas` in `.zshellcheckrc`.

---

<a id="zc1631"></a>
### ZC1631 — Error on `openssl ... -passin pass:SECRET` / `-passout pass:SECRET`

**Severity:** `error`  
**Auto-fix:** `no`

OpenSSL's `-passin` / `-passout` accept a password source selector. The `pass:LITERAL` form embeds the password as an argv element — visible in `ps`, `/proc/<pid>/cmdline`, shell history, and audit logs. Use one of the safer sources: `env:VARNAME` reads from an env var, `file:PATH` reads the first line of PATH, `fd:N` reads from an open descriptor, `stdin` reads a line from stdin.

Disable by adding `ZC1631` to `disabled_katas` in `.zshellcheckrc`.

---

<a id="zc1632"></a>
### ZC1632 — Warn on `shred` — unreliable on journaled / CoW filesystems (ext4, btrfs, zfs)

**Severity:** `warning`  
**Auto-fix:** `no`

`shred` assumes in-place overwrites, which is how ext2 worked. On a journaled ext4 the overwrite passes go through the journal and may not hit the original data blocks. On CoW filesystems (btrfs, zfs, xfs with reflink) the overwrite lands in fresh blocks and leaves the old content intact until garbage collection decides otherwise. `shred`'s own man page warns about this. For modern secure deletion, use full-disk encryption with key destruction, or retire the device with `blkdiscard` on an SSD.

Disable by adding `ZC1632` to `disabled_katas` in `.zshellcheckrc`.

---

<a id="zc1633"></a>
### ZC1633 — Error on `gpg --passphrase SECRET` — passphrase on cmdline

**Severity:** `error`  
**Auto-fix:** `no`

`gpg --passphrase VALUE` passes the key passphrase as an argv element. Visible in `ps`, `/proc/<pid>/cmdline`, shell history, and audit logs for every local user who can list processes. Use `--passphrase-file PATH` (reads the first line of PATH), `--passphrase-fd N` (reads from file descriptor N), or `--pinentry-mode=loopback` with the passphrase piped on stdin. Pair with `--batch` for non-interactive runs.

Disable by adding `ZC1633` to `disabled_katas` in `.zshellcheckrc`.

---

<a id="zc1634"></a>
### ZC1634 — Warn on `umask NNN` that fails to mask world-write — mask-inversion footgun

**Severity:** `warning`  
**Auto-fix:** `no`

`umask` is a mask: bits that are set are removed from the default permission. The classic pitfall is reading it as "permissions I want" — `umask 111` feels tight ("no execute for anyone") but it does not mask the write bit, so every new file is `666` (rw-rw-rw-). The "other" digit must be one of `2/3/6/7` to strip world-write. Use `022` for publicly readable files, `077` for secrets-handling.

Disable by adding `ZC1634` to `disabled_katas` in `.zshellcheckrc`.

---

<a id="zc1635"></a>
### ZC1635 — Error on `mysql -pSECRET` / `--password=SECRET` — password in process list

**Severity:** `error`  
**Auto-fix:** `no`

MySQL / MariaDB clients accept the password concatenated with the `-p` flag (`-pSECRET`) or via `--password=SECRET`. Both forms put the secret in argv — visible in `ps`, `/proc/<pid>/cmdline`, shell history, and audit logs for every local user who can list processes. Use `-p` with no argument for an interactive prompt, `--login-path` for the credentials helper file, or a `~/.my.cnf` with `0600` perms and `[client] password=...` so the client reads it at startup.

Disable by adding `ZC1635` to `disabled_katas` in `.zshellcheckrc`.

---

<a id="zc1636"></a>
### ZC1636 — Warn on `virsh destroy DOMAIN` — force-stops VM (no graceful shutdown)

**Severity:** `warning`  
**Auto-fix:** `no`

`virsh destroy DOM` is the libvirt equivalent of pulling the plug on a running VM. The guest OS gets no chance to flush filesystems, close network connections, or run its own shutdown services — data corruption risk on any open file in the guest. For graceful shutdown use `virsh shutdown DOM` (ACPI event), wait for completion, and only fall back to `destroy` for a genuinely unresponsive guest. `virsh destroy --graceful DOM` attempts a timed graceful first, then forces — that variant is not flagged.

Disable by adding `ZC1636` to `disabled_katas` in `.zshellcheckrc`.

---

<a id="zc1637"></a>
### ZC1637 — Style: prefer Zsh `typeset -r NAME=value` over POSIX `readonly NAME=value`

**Severity:** `style`  
**Auto-fix:** `yes`

Both `readonly NAME` and `typeset -r NAME` create a read-only parameter. In Zsh the idiomatic form is `typeset -r` — it composes with other typeset flags (`-ir` for readonly integer, `-xr` for readonly export, `-gr` to pin a readonly global from inside a function). `readonly` works but reads as a Bash / POSIX-ism in a Zsh codebase.

Disable by adding `ZC1637` to `disabled_katas` in `.zshellcheckrc`.

---

<a id="zc1638"></a>
### ZC1638 — Error on `docker/podman build --build-arg SECRET=VALUE` — secret baked into image layer

**Severity:** `error`  
**Auto-fix:** `no`

`--build-arg KEY=VALUE` values land in the image metadata that `docker history` (and the analogous podman / buildah tooling) read back from the layer. Even if the Dockerfile only uses the arg to export as a build-time env var, the literal value is cached in the layer forever. A key-shaped name (`password`, `secret`, `token`, `apikey`, `access_key`, `private_key`) with a concrete value embeds that secret in every image pulled. Use BuildKit secrets (`--secret id=mysecret,src=path`) or a multi-stage build where the secret stays in a discarded stage.

Disable by adding `ZC1638` to `disabled_katas` in `.zshellcheckrc`.

---

<a id="zc1639"></a>
### ZC1639 — Error on `curl -H 'Authorization: ...'` — credential header in process list

**Severity:** `error`  
**Auto-fix:** `no`

`-H "Authorization: Bearer $TOKEN"` (and similar credential-bearing headers like `X-Api-Key`, `X-Auth-Token`, `Proxy-Authorization`, `Cookie`) put the expanded value in argv. It shows up in `ps`, `/proc/<pid>/cmdline`, shell history, and audit logs — every local user who can list processes reads the secret. Pass the header via a file with `-H @FILE` or use `--config FILE` so the value stays on disk (with 0600 perms), never on the command line.

Disable by adding `ZC1639` to `disabled_katas` in `.zshellcheckrc`.

---

<a id="zc1640"></a>
### ZC1640 — Style: `${!var}` Bash indirect expansion — prefer Zsh `${(P)var}`

**Severity:** `style`  
**Auto-fix:** `no`

`${!var}` is Bash indirect expansion — it reads the value of the parameter whose name is stored in `$var`. Zsh has the native flag form `${(P)var}` which does the same and composes with other parameter-expansion flags (`${(Pf)var}` to split the indirect value on newlines, for example). `${!prefix*}` / `${!array[@]}` have Zsh equivalents via the `$parameters` hash or `(k)` subscript flags. Prefer the native Zsh form in a Zsh codebase.

Disable by adding `ZC1640` to `disabled_katas` in `.zshellcheckrc`.

---

<a id="zc1641"></a>
### ZC1641 — Error on `kubectl create secret --from-literal=...` / `--docker-password=...`

**Severity:** `error`  
**Auto-fix:** `no`

`kubectl create secret generic --from-literal=KEY=VALUE` and `kubectl create secret docker-registry --docker-password=VALUE` put the secret content in argv. The expanded value shows up in `ps`, `/proc/<pid>/cmdline`, shell history, and audit logs — readable by any local user who can list processes. Use `--from-file=KEY=PATH` (reads from a 0600-protected file), `--from-env-file=PATH` (reads KEY=VALUE lines), or pipe a manifest into `kubectl apply -f -` with base64-encoded `data:` values staged on disk.

Disable by adding `ZC1641` to `disabled_katas` in `.zshellcheckrc`.

---

<a id="zc1642"></a>
### ZC1642 — Warn on `tshark -w FILE` / `dumpcap -w FILE` without `-Z user` — capture file owned by root

**Severity:** `warning`  
**Auto-fix:** `no`

Packet captures routinely need `CAP_NET_RAW`, so the capture process typically runs as root. Without `-Z USER` the resulting pcap is root-owned — a subsequent analyst who opens it with Wireshark (which can run helper scripts from the file) operates on a root-owned file and may unintentionally invoke things as root. `-Z USER` tells `tshark` / `dumpcap` to drop privileges for the actual capture and write the file as `USER`.

Disable by adding `ZC1642` to `disabled_katas` in `.zshellcheckrc`.

---

<a id="zc1643"></a>
### ZC1643 — Style: `$(cat file)` — use `$(<file)` to skip the fork / exec

**Severity:** `style`  
**Auto-fix:** `yes`

`$(cat FILE)` forks, execs `/usr/bin/cat`, reads FILE, writes the bytes to the pipe, waits for the child. `$(<FILE)` is a shell builtin — it reads FILE directly into the command-substitution buffer with no fork and no exec. In a hot path the speedup is dramatic, and even in cold paths it avoids one of the most common useless-use-of-cat patterns in review feedback.

Disable by adding `ZC1643` to `disabled_katas` in `.zshellcheckrc`.

---

<a id="zc1644"></a>
### ZC1644 — Error on `unzip -P SECRET` / `zip -P SECRET` — archive password in process list

**Severity:** `error`  
**Auto-fix:** `no`

`unzip -P PASSWORD` / `zip -P PASSWORD` (or the concatenated `-PPASSWORD` form) places the archive password in argv. The expanded value shows up in `ps`, `/proc/<pid>/cmdline`, shell history, and audit logs for every local user who can list processes. Both tools prompt interactively if `-P` is absent — use that for human workflows. For automation prefer an archive format with a real key-derivation story (for example `7z -p` piped over stdin, or `age` / `gpg` envelope encryption that reads keys from a protected file).

Disable by adding `ZC1644` to `disabled_katas` in `.zshellcheckrc`.

---

<a id="zc1645"></a>
### ZC1645 — Style: `lsb_release` — prefer sourcing `/etc/os-release` (no dependency, no fork)

**Severity:** `style`  
**Auto-fix:** `no`

`lsb_release` is provided by the `lsb-release` / `redhat-lsb-core` package, which is missing on most minimal / container images (Alpine does not ship it at all). Scripts that depend on `lsb_release` fail the moment they hit a stripped image. `/etc/os-release` is standardized by systemd and always present on modern Linux — `source /etc/os-release; print -r -- $ID $VERSION_ID` gives the same distribution info without the extra package, and without forking.

Disable by adding `ZC1645` to `disabled_katas` in `.zshellcheckrc`.

---

<a id="zc1646"></a>
### ZC1646 — Warn on `btrfs check --repair` / `xfs_repair -L` — last-resort recovery, may worsen damage

**Severity:** `warning`  
**Auto-fix:** `no`

Both commands are destructive last-resort recovery. `btrfs check --repair` explicitly warns in its man page that it "may cause additional filesystem damage" and the btrfs developers ask users to try `btrfs scrub` and read-only `btrfs check` first. `xfs_repair -L` zeroes the log, dropping any uncommitted transactions and the data they held. In both cases snapshot the underlying block device before running, so the attempt is reversible.

Disable by adding `ZC1646` to `disabled_katas` in `.zshellcheckrc`.

---

<a id="zc1647"></a>
### ZC1647 — Warn on `kubectl apply -f URL` — remote manifest applied without digest verification

**Severity:** `warning`  
**Auto-fix:** `no`

`kubectl apply -f https://...` fetches the manifest over the network and applies it to the cluster. TLS (when present) verifies transport but not authorship — if the URL is compromised or the content changes between reviews, the cluster picks up the new definition. Pin the content: download to disk, verify a known SHA256, then `kubectl apply -f local.yaml`. For plain HTTP the attacker controls the response directly — never acceptable.

Disable by adding `ZC1647` to `disabled_katas` in `.zshellcheckrc`.

---

<a id="zc1648"></a>
### ZC1648 — Error on `cp /dev/null /var/log/...` / `truncate -s 0 /var/log/...` — audit-log wipe

**Severity:** `error`  
**Auto-fix:** `no`

Replacing a file under `/var/log/` with `/dev/null` or truncating it to size zero erases audit evidence: failed login attempts from `auth.log`, sudo usage from `sudo.log`, kernel audit trail from `audit/audit.log`, console history from `wtmp` / `btmp`. Scripts that do this during "cleanup" are almost always misusing logrotate (which handles rotation safely via a `create` stage) or deliberately covering tracks. Use `logrotate -f /etc/logrotate.d/...` for rotation, `journalctl --vacuum-time=...` for journald.

Disable by adding `ZC1648` to `disabled_katas` in `.zshellcheckrc`.

---

<a id="zc1649"></a>
### ZC1649 — Warn on `openssl req -days N` with N > 825 — long-validity certificate

**Severity:** `warning`  
**Auto-fix:** `no`

CA/Browser Forum capped public TLS cert validity at 825 days in 2018 and major browsers tightened it to 398 days in 2020. A cert issued for 3650 days (10 years) can not be revoked effectively — once the private key leaks, the attacker keeps access until the cert expires naturally. For an internal root CA the long validity is defensible; for leaf / server certs keep it under 398 days and automate rotation. `-days` over 825 almost always means "I don't want to deal with renewal," which is a maintenance smell dressed up as security.

Disable by adding `ZC1649` to `disabled_katas` in `.zshellcheckrc`.

---

<a id="zc1650"></a>
### ZC1650 — Warn on `setopt RM_STAR_SILENT` / `unsetopt RM_STAR_WAIT` — removes `rm *` prompt

**Severity:** `warning`  
**Auto-fix:** `no`

Zsh's default behaviour on an interactive `rm *` (or `rm /path/*`) is to pause for 10 seconds and ask "do you really want to delete N files?" — the `RM_STAR_WAIT` option. `setopt RM_STAR_SILENT` or `unsetopt RM_STAR_WAIT` both disable the prompt. In a profile / dot file the option leaks to every future interactive shell and removes a safety net that has saved countless home directories.

Disable by adding `ZC1650` to `disabled_katas` in `.zshellcheckrc`.

---

<a id="zc1651"></a>
### ZC1651 — Warn on `docker/podman run -p 0.0.0.0:PORT:PORT` — explicit all-interfaces publish

**Severity:** `warning`  
**Auto-fix:** `no`

A port spec of `0.0.0.0:HOST:CONT`, `[::]:HOST:CONT`, or `*:HOST:CONT` publishes the container port to every interface the host has. On a multi-tenant LAN or a cloud host with a public IP the service is immediately reachable from anywhere. If the service needs only local reverse-proxy access, bind to `127.0.0.1:HOST:CONT` and let nginx / caddy handle external exposure.

Disable by adding `ZC1651` to `disabled_katas` in `.zshellcheckrc`.

---

<a id="zc1652"></a>
### ZC1652 — Warn on `ssh -Y` — trusted X11 forwarding grants full X-server access to remote clients

**Severity:** `warning`  
**Auto-fix:** `no`

`ssh -Y` enables trusted X11 forwarding. Remote X clients can read every keystroke on the local display, take screenshots, inject synthetic events, and otherwise drive the local session with no sandbox. `ssh -X` enables the untrusted variant, which routes X traffic through the X SECURITY extension so those capabilities are limited (some GUI features break, which is why people reach for `-Y` — usually at far higher risk than they realised). Prefer `-X` when X11 forwarding is genuinely needed; better yet drop it for Wayland tools or VNC-over-SSH with its own auth.

Disable by adding `ZC1652` to `disabled_katas` in `.zshellcheckrc`.

---

<a id="zc1653"></a>
### ZC1653 — Avoid `$BASHPID` — Bash-only; Zsh uses `$sysparams[pid]` from `zsh/system`

**Severity:** `warning`  
**Auto-fix:** `no`

`$BASHPID` returns the PID of the current subshell (while `$$` returns the parent shell's PID). In Zsh this parameter is not set — scripts that rely on `$BASHPID` silently get an empty string and misbehave. After `zmodload zsh/system`, Zsh exposes the current process PID as `$sysparams[pid]`, which updates inside subshells just like Bash's `$BASHPID`.

Disable by adding `ZC1653` to `disabled_katas` in `.zshellcheckrc`.

---

<a id="zc1654"></a>
### ZC1654 — Warn on `sysctl -p /tmp/...` — loading kernel tunables from attacker-writable path

**Severity:** `warning`  
**Auto-fix:** `no`

`sysctl -p PATH` reads `key=value` lines from PATH and applies them as kernel tunables. A PATH under `/tmp/` or `/var/tmp/` is world-traversable; a concurrent local user can substitute the file between write and read, injecting `kernel.core_pattern=|/tmp/evil`, `kernel.modprobe=/tmp/evil`, or disabling hardening knobs (`kernel.kptr_restrict=0`, `kernel.yama.ptrace_scope=0`). Keep sysctl configs under `/etc/sysctl.d/` with root ownership.

Disable by adding `ZC1654` to `disabled_katas` in `.zshellcheckrc`.

---

<a id="zc1655"></a>
### ZC1655 — Warn on `read -n N` — Bash reads N chars; Zsh's `-n` means "drop newline"

**Severity:** `warning`  
**Auto-fix:** `no`

In Bash, `read -n N var` reads exactly N characters (handy for single-keypress prompts). In Zsh, `-n` is the "don't append newline to the reply string" flag and doesn't take a count — `read -n 1 var` sets `var` to the whole line, not a single character. Use `read -k N var` in Zsh for N-character reads.

Disable by adding `ZC1655` to `disabled_katas` in `.zshellcheckrc`.

---

<a id="zc1656"></a>
### ZC1656 — Error on `rsync -e 'ssh -o StrictHostKeyChecking=no'` — host-key verify disabled

**Severity:** `error`  
**Auto-fix:** `no`

Disabling host-key verification through rsync's `-e` transport is the same attack surface as ZC1479 but easier to miss in review because the ssh flags sit inside a quoted string. A MITM on the network path can impersonate the remote host and the rsync stream goes straight through. Use `ssh-keyscan` or pre-provisioned `~/.ssh/known_hosts` to trust hosts deliberately, and keep `StrictHostKeyChecking=yes`.

Disable by adding `ZC1656` to `disabled_katas` in `.zshellcheckrc`.

---

<a id="zc1657"></a>
### ZC1657 — Warn on `semanage permissive -a <type>` — puts SELinux domain in permissive mode

**Severity:** `warning`  
**Auto-fix:** `no`

`semanage permissive -a DOMAIN` (or `--add`) marks an SELinux domain as permissive: policy violations are logged but not blocked. It is narrower than `setenforce 0` but still disables enforcement for whatever DOMAIN covers — often `httpd_t`, `container_t`, or `sshd_t` — and the override persists across reboots because it is written to policy. Fix the denial with an explicit allow rule built from `audit2allow` or ship a custom policy module, and remove the permissive mark with `semanage permissive -d DOMAIN` once the rule lands.

Disable by adding `ZC1657` to `disabled_katas` in `.zshellcheckrc`.

---

<a id="zc1658"></a>
### ZC1658 — Warn on `curl -OJ` / `-J -O` — server-controlled output filename

**Severity:** `warning`  
**Auto-fix:** `no`

`curl -J` (`--remote-header-name`) combined with `-O` (`--remote-name`) saves the response using the filename the server puts in the `Content-Disposition` header. The server — or anything on the path that can set headers, including a compromised CDN or an HTTP-serving reverse proxy — chooses the destination name. Paths like `../../etc/cron.d/evil` are rejected by curl's sanitizer, but benign-looking names still overwrite files in the current directory. Use `-o NAME` with a filename you control, and validate the payload before you act on it.

Disable by adding `ZC1658` to `disabled_katas` in `.zshellcheckrc`.

---

<a id="zc1659"></a>
### ZC1659 — Warn on `fuser -k <path>` — kills every process holding the subtree open

**Severity:** `warning`  
**Auto-fix:** `no`

`fuser -k PATH` sends a signal (SIGKILL by default) to every process that has any file under PATH open — not just the one you expected. On `/`, `/var`, `/tmp`, or any mount-root this reaches sshd, cron, dbus, and the caller's own shell; on a bind-mount it kills workloads that share the host inode. Target specific PIDs (`kill $(pidof app)`) or ports (`fuser -k PORT/tcp`), or use `systemctl stop UNIT` for services. `fuser -k` against a filesystem path is blast-radius that the caller rarely owns.

Disable by adding `ZC1659` to `disabled_katas` in `.zshellcheckrc`.

---

<a id="zc1660"></a>
### ZC1660 — Style: `printf '%0Nd' $n` — prefer Zsh `${(l:N::0:)n}` left-zero-pad

**Severity:** `style`  
**Auto-fix:** `no`

Zero-padding an integer through `printf '%0Nd'` forks a tiny sub-process and relies on printf's format-string parser — both things Zsh can avoid. `${(l:N::0:)n}` left-pads `$n` with `0` to width N using Zsh parameter expansion, no fork, and composes cleanly with other `(q)` / `(L)` / `(U)` flags. For right-pad use `${(r:N::0:)n}`; for space padding swap the fill character: `${(l:N:)n}` or `${(r:N:)n}`.

Disable by adding `ZC1660` to `disabled_katas` in `.zshellcheckrc`.

---

<a id="zc1661"></a>
### ZC1661 — Error on `curl --cacert /dev/null` — empty trust store, any cert passes

**Severity:** `error`  
**Auto-fix:** `no`

Pointing `--cacert` (or `--capath`) at `/dev/null` hands curl an empty trust anchor set. Counter-intuitively, curl treats the peer certificate as valid when no issuers are configured for the selected TLS backend (OpenSSL, wolfSSL, Schannel all accept any cert chain against an empty CA bundle). This is the TLS equivalent of `--insecure` with one more keystroke of plausible deniability. Use a real bundle (`/etc/ssl/certs/ca-certificates.crt`) or `--pinnedpubkey sha256//…` for known endpoints.

Disable by adding `ZC1661` to `disabled_katas` in `.zshellcheckrc`.

---

<a id="zc1662"></a>
### ZC1662 — Error on `pkexec env VAR=VAL CMD` — controlled env crossed into the root session

**Severity:** `error`  
**Auto-fix:** `no`

`pkexec env VAR=VALUE CMD` invokes `/usr/bin/env` as the target user (root by default) with a caller-controlled environment. Polkit sanitizes a short allow-list on its own, but once `env` takes over the remaining variables (`LD_PRELOAD`, `GCONV_PATH`, `PYTHONPATH`, `XDG_RUNTIME_DIR`, `LANGUAGE`) ride straight into root. CVE-2021-4034 (pwnkit) demonstrated the same primitive by abusing argv[0]; the `env` wrapper makes the bypass trivial. If the child needs specific variables, set them in a polkit rule or via `systemd-run --user` instead, not through `env`.

Disable by adding `ZC1662` to `disabled_katas` in `.zshellcheckrc`.

---

<a id="zc1663"></a>
### ZC1663 — Warn on `tune2fs -c 0` / `-i 0` — disables periodic filesystem checks

**Severity:** `warning`  
**Auto-fix:** `no`

`tune2fs -c 0` (mount count) and `tune2fs -i 0` (time interval) disable the ext2/3/4 periodic-fsck machinery so the filesystem only gets checked after a dirty unmount or a manual `fsck -f`. For desktops the nag is annoying; for long-lived servers it is the last line of defence against silent metadata corruption. Lower the cadence if the default is too aggressive (`tune2fs -c 30`, `-i 3m`) rather than turning it off, and schedule an offline `fsck` on a cadence you can defend.

Disable by adding `ZC1663` to `disabled_katas` in `.zshellcheckrc`.

---

<a id="zc1664"></a>
### ZC1664 — Error on `systemctl set-default rescue.target|emergency.target` — persistent single-user boot

**Severity:** `error`  
**Auto-fix:** `no`

`systemctl set-default` rewrites `/etc/systemd/system/default.target` as a symlink to the named target. Pointing it at `rescue.target` or `emergency.target` means every subsequent boot drops to single-user mode before networking, sshd, or any normal unit starts — you lose remote access to the box unless you have serial console / out-of-band management. Unlike `systemctl isolate` (one-shot, caught by ZC1561) this persists across reboots. Revert with `systemctl set-default multi-user.target` (servers) or `graphical.target` (desktops).

Disable by adding `ZC1664` to `disabled_katas` in `.zshellcheckrc`.

---

<a id="zc1665"></a>
### ZC1665 — Warn on `chrt -r` / `-f` — real-time scheduling class from a shell script

**Severity:** `warning`  
**Auto-fix:** `no`

`chrt -r PRIO CMD` (SCHED_RR) and `chrt -f PRIO CMD` (SCHED_FIFO) launch the child under a POSIX real-time scheduling class. An RT thread preempts every normal-priority task until it voluntarily yields; a busy-loop or a deadlock leaves the kernel with kworker, ksoftirqd, and sshd starved, often forcing a hard reboot. Unless the binary is known-bounded (audio glitch-free path, protocol timing loop), keep scripts on SCHED_OTHER — use `nice -n -5` or a systemd unit with `CPUWeight=` / `IOWeight=` instead of `chrt -r`.

Disable by adding `ZC1665` to `disabled_katas` in `.zshellcheckrc`.

---

<a id="zc1666"></a>
### ZC1666 — Warn on `kubectl patch --type=json` — bypasses strategic-merge defaults

**Severity:** `warning`  
**Auto-fix:** `no`

`kubectl patch --type=json` applies a raw RFC-6902 JSON patch: `remove`, `replace`, `add /spec/containers/0`, and `move` land verbatim on the resource. Unlike strategic-merge or merge-patch, Kubernetes does not reconcile the patch against field ownership or default values — so a mistyped `path` or an index that no longer exists fails silently or drops the wrong field. From a script this is a foot-gun for drift and supply-chain compromise: an attacker with write access to the patch file can slip `privileged: true` or `hostPath` mounts in. Prefer `--type=strategic` (the default) and hold JSON patches behind code review.

Disable by adding `ZC1666` to `disabled_katas` in `.zshellcheckrc`.

---

<a id="zc1667"></a>
### ZC1667 — Warn on `openssl enc` without `-pbkdf2` — legacy MD5-based key derivation

**Severity:** `warning`  
**Auto-fix:** `no`

Without `-pbkdf2`, `openssl enc` derives the symmetric key through EVP_BytesToKey, which is a single MD5 round over `password || salt`. A modern GPU cracks that at billions of guesses per second. Add `-pbkdf2 -iter 100000` (OpenSSL 1.1.1+) to switch to PBKDF2-HMAC-SHA256 with a real iteration count. Even better, stop using `openssl enc` for new code — it has no AEAD support and `-aes-256-gcm` silently drops the auth tag — and reach for `age`, `gpg --symmetric --cipher-algo AES256`, or `openssl smime` instead.

Disable by adding `ZC1667` to `disabled_katas` in `.zshellcheckrc`.

---

<a id="zc1668"></a>
### ZC1668 — Error on `aws iam attach-*-policy ... AdministratorAccess` — grants full AWS admin

**Severity:** `error`  
**Auto-fix:** `no`

Attaching the AWS-managed `AdministratorAccess` (or `PowerUserAccess`) policy gives the target principal `*:*` — create/delete IAM users, mutate KMS keys, rotate root passwords, exfiltrate every S3 bucket. Scripts rarely need full admin; the pattern usually means someone hit a permissions error and replaced the scoped policy with the blanket one. Write a least-privilege inline policy (`iam put-user-policy --policy-document`), or reference a customer-managed policy with only the `Action`/`Resource` pairs the workload needs. Admin attachment should land via change-reviewed Terraform, not a shell loop.

Disable by adding `ZC1668` to `disabled_katas` in `.zshellcheckrc`.

---

<a id="zc1669"></a>
### ZC1669 — Warn on `git gc --prune=now` / `git reflog expire --expire=now` — deletes recovery window

**Severity:** `warning`  
**Auto-fix:** `no`

Git keeps dropped commits and orphaned objects for `gc.reflogExpire` (default 90 days) and `gc.pruneExpire` (default two weeks) so a `git reflog` + `git reset` can still recover work you thought you threw away. `git gc --prune=now` and `git reflog expire --expire=now --all` bulldoze both windows in one go — a stray interactive rebase no longer has a safety net. Use the default cadence (`git gc`, no `--prune=now`) unless you are actively purging leaked secrets or proof-of-concept code; pair the destructive form with a stale mirror push so at least one copy of the dropped history remains.

Disable by adding `ZC1669` to `disabled_katas` in `.zshellcheckrc`.

---

<a id="zc1670"></a>
### ZC1670 — Warn on `setsebool -P` enabling memory-protection-relaxing SELinux boolean

**Severity:** `warning`  
**Auto-fix:** `no`

Specific SELinux policy booleans (`allow_execstack`, `allow_execmem`, `httpd_execmem`, `selinuxuser_execstack`, `domain_kernel_load_modules`, `mmap_low_allowed`, etc.) relax per-domain memory protections that the policy puts in place precisely because those domains should not need writable-and-executable pages. Persisting the flip with `-P` carries the regression across reboots. Fix the underlying binary (`execstack -c`, `chcon`, stop generating runtime-JIT code in the wrong domain) instead of loosening policy.

Disable by adding `ZC1670` to `disabled_katas` in `.zshellcheckrc`.

---

<a id="zc1671"></a>
### ZC1671 — Error on `install -m 777` / `mkdir -m 777` — creates world-writable target

**Severity:** `error`  
**Auto-fix:** `no`

`install -m MODE` / `mkdir -m MODE` applies MODE atomically at file or directory creation, so the world-writable window from a later `chmod 777` is not even needed — the path is wide-open from the moment it exists. Any local user can swap binaries under `/usr/local/bin`, write shell-completion hooks into `/etc/bash_completion.d`, or turn a shared directory into an LPE staging ground. Drop the world-write bit: `0755` for binaries, `0644` for files, `2770` with `chgrp` for shared directories.

Disable by adding `ZC1671` to `disabled_katas` in `.zshellcheckrc`.

---

<a id="zc1672"></a>
### ZC1672 — Info: `chcon` writes an ephemeral SELinux label — next `restorecon` wipes it

**Severity:** `info`  
**Auto-fix:** `no`

`chcon -t TYPE PATH` sets the file context out-of-band; it does not update the `file_contexts` policy database. As soon as `restorecon`, `semodule -n`, or a policy rebuild runs, the label snaps back to whatever the compiled policy says — often `default_t`, which can break a deployed workload or silently re-introduce a denial the script tried to fix. For anything long-lived use `semanage fcontext -a -t TYPE '<regex>'` then `restorecon -F <path>` so the mapping lives in policy.

Disable by adding `ZC1672` to `disabled_katas` in `.zshellcheckrc`.

---

<a id="zc1673"></a>
### ZC1673 — Style: `stty -echo` around `read` — prefer Zsh `read -s`

**Severity:** `style`  
**Auto-fix:** `no`

The classic `stty -echo; IFS= read -r password; stty echo` pattern has a serious failure mode: a crash or SIGINT between the two `stty` calls leaves the user's terminal stuck in echo-off, which is silent and confusing. Zsh's `read -s VAR` (also Bash 4+) disables echo only for that one `read`, restores it on return even if the read is interrupted, and avoids two external forks. Switch the prompt to `read -s` (or `read -ks` for single-key password) and drop the `stty` bracketing.

Disable by adding `ZC1673` to `disabled_katas` in `.zshellcheckrc`.

---

<a id="zc1674"></a>
### ZC1674 — Warn on `docker/podman run --oom-kill-disable` or `--oom-score-adj <= -500`

**Severity:** `warning`  
**Auto-fix:** `no`

`--oom-kill-disable` tells the kernel OOM killer to never touch the container's memory cgroup — a leak inside then drives the whole host into OOM reclaim until `sshd`, `systemd-journald`, or the init daemon itself gets killed. `--oom-score-adj <= -500` stops short of full immunity but still preferentially kills unrelated host processes under pressure. If the workload genuinely needs resilience, cap memory with `--memory=<limit>` and accept the container being killed on overrun; shift the heavy workload to a dedicated node instead of rigging OOM scores.

Disable by adding `ZC1674` to `disabled_katas` in `.zshellcheckrc`.

---

<a id="zc1675"></a>
### ZC1675 — Avoid Bash-only `export -f` / `export -n` — use Zsh `typeset -fx` / `typeset +x`

**Severity:** `info`  
**Auto-fix:** `yes`

`export -f FUNC` (export a function to child processes) and `export -n VAR` (strip the export flag while keeping the value) are Bash-only. Zsh's `export` ignores `-f` entirely and prints usage for `-n`, so scripts that depend on either silently break under Zsh. The Zsh equivalents are `typeset -fx FUNC` for function export (parameter-passing via `$FUNCTIONS` in a subshell) and `typeset +x VAR` to drop the export flag. Functions that must cross a subshell are usually better handled by `autoload -Uz` from an `fpath` directory than by serialisation.

Disable by adding `ZC1675` to `disabled_katas` in `.zshellcheckrc`.

---

<a id="zc1676"></a>
### ZC1676 — Warn on `helm rollback --force` — recreates in-flight resources, corrupts rolling updates

**Severity:** `warning`  
**Auto-fix:** `no`

`helm rollback RELEASE N --force` asks Helm to delete and recreate any resource that it cannot patch cleanly. If a deployment is mid-rollout, the `--force` flag takes out both the old and new ReplicaSets, kicks the pods, and forces a cold start — losing in-flight requests and any `PodDisruptionBudget` protections. Worse, rolling back to revision N brings back whatever CVEs or config regressions the later revisions had already fixed. Pin the target revision explicitly, omit `--force`, and gate the rollback behind a change-review ticket rather than a shell one-liner.

Disable by adding `ZC1676` to `disabled_katas` in `.zshellcheckrc`.

---

<a id="zc1677"></a>
### ZC1677 — Warn on `trap 'set -x' DEBUG` — xtrace on every command leaks secrets

**Severity:** `warning`  
**Auto-fix:** `no`

`trap 'set -x' DEBUG` runs the trap handler before every simple command, turning on xtrace for the remainder of the shell. Every subsequent `curl -H 'Authorization: Bearer …'`, `mysql -p<password>`, or `aws configure set …` then prints its full argv to stderr — commonly into a log file or CI artifact. The same antipattern shows up as `set -o xtrace` inside a DEBUG trap. Instrument selectively with `typeset -ft FUNC` (Zsh function-level xtrace), or add `exec 2>>"$log"; set -x` only around the part of the script you want traced.

Disable by adding `ZC1677` to `disabled_katas` in `.zshellcheckrc`.

---

<a id="zc1678"></a>
### ZC1678 — Error on `borg init --encryption=none` — unencrypted backup repository

**Severity:** `error`  
**Auto-fix:** `no`

`borg init --encryption=none REPO` creates a backup repository without client-side encryption or authentication. Anyone with read access to the repo gets every file in every archive, and no one can detect silent tampering — borg will happily extract a modified chunk. Even for local-only repos the cost of authenticated-encryption is tiny; use `--encryption=repokey-blake2` (or `--encryption=keyfile-blake2` when you want the key off the server), and store the passphrase in `BORG_PASSPHRASE_FILE` pointing at a mode-0400 file.

Disable by adding `ZC1678` to `disabled_katas` in `.zshellcheckrc`.

---

<a id="zc1679"></a>
### ZC1679 — Error on `gcloud ... add-iam-policy-binding ... --role=roles/owner` — GCP primitive admin

**Severity:** `error`  
**Auto-fix:** `no`

`gcloud projects|folders|organizations add-iam-policy-binding` with the primitive roles `roles/owner` or `roles/editor`, or with the IAM-escalation roles (`roles/iam.securityAdmin`, `roles/iam.serviceAccountTokenCreator`, `roles/iam.serviceAccountKeyAdmin`, `roles/resourcemanager.organizationAdmin`), hands the principal the ability to grant themselves any other permission. Scripts rarely need that scope; the pattern signals someone papering over a permissions error. Grant a specific predefined role (e.g. `roles/compute.viewer`) or build a custom role with only the `Action`s the workload needs, and apply admin changes via Terraform under change review.

Disable by adding `ZC1679` to `disabled_katas` in `.zshellcheckrc`.

---

<a id="zc1680"></a>
### ZC1680 — Error on `ansible-playbook --vault-password-file=/tmp/...` — world-traversable vault key

**Severity:** `error`  
**Auto-fix:** `no`

The Ansible Vault decryption key lives in the `--vault-password-file` path. `/tmp`, `/var/tmp`, and `/dev/shm` are world-traversable: a concurrent local user who guesses (or `inotifywait`s for) the filename opens it during the playbook run and dumps every secret the vault protects. Keep vault keys in a root-owned mode-0400 file under `/etc/ansible/` or `$HOME/.ansible/`, or supply the passphrase via a no-echo helper script (`vault-password-client`) that fetches from `pass` / `vault kv get`.

Disable by adding `ZC1680` to `disabled_katas` in `.zshellcheckrc`.

---

<a id="zc1681"></a>
### ZC1681 — Error on `tar -P` / `--absolute-names` — archive absolute paths, can overwrite host files

**Severity:** `error`  
**Auto-fix:** `no`

By default GNU tar strips the leading `/` from archive member paths so that extraction stays under the current directory. `-P` (or the long form `--absolute-names`) disables that strip: `tar -xPf evil.tar` happily writes to `/etc/cron.d/evil`, `/usr/local/bin/sshd`, or any other absolute path the archive mentions. Archives from untrusted sources should never be unpacked with `-P`. Drop the flag, extract with `-C <scratch-dir>`, audit the tree, then copy files into place with `install` or `cp`.

Disable by adding `ZC1681` to `disabled_katas` in `.zshellcheckrc`.

---

<a id="zc1682"></a>
### ZC1682 — Error on `npm install --unsafe-perm` — npm lifecycle scripts keep root privileges

**Severity:** `error`  
**Auto-fix:** `no`

npm normally drops to the UID that owns `package.json` before running `preinstall` / `install` / `postinstall` lifecycle scripts. `--unsafe-perm` (or `--unsafe-perm=true`) tells npm to skip that drop and run every script as the current UID — typically root when the install happens from a provisioning script. Any compromised or malicious dependency then executes as root. If a native addon truly needs privileges, scope them: drop them into a dedicated builder container, or use `sudo -u builduser npm install` from a non-root account that already owns `node_modules/`.

Disable by adding `ZC1682` to `disabled_katas` in `.zshellcheckrc`.

---

<a id="zc1683"></a>
### ZC1683 — Error on `npm/yarn/pnpm config set registry http://...` — plaintext package index

**Severity:** `error`  
**Auto-fix:** `no`

Pointing a JavaScript package manager at an `http://` registry disables TLS during fetch. Any host on the path (corporate proxy, hotel Wi-Fi, compromised CDN) can rewrite tarballs mid-flight; lockfile hashes catch the rewrite only if the user locks every dependency before the swap. Even on internal networks, pin to `https://` — reach for your own CA via `NODE_EXTRA_CA_CERTS` or `registry.cafile` rather than falling back to HTTP.

Disable by adding `ZC1683` to `disabled_katas` in `.zshellcheckrc`.

---

<a id="zc1684"></a>
### ZC1684 — Error on `redis-cli -a PASSWORD` — authentication password in process list

**Severity:** `error`  
**Auto-fix:** `no`

`redis-cli -a <password>` (and the joined form `-aPASSWORD`) puts the authentication password in the command line — visible to every user on the host through `ps`, `/proc/PID/cmdline`, audit logs, and shell history. redis-cli 6.0+ prints a warning to stderr but still connects. Use the `REDISCLI_AUTH` environment variable (read automatically by redis-cli), or `-askpass` to prompt from TTY; both keep the secret out of the argv tail.

Disable by adding `ZC1684` to `disabled_katas` in `.zshellcheckrc`.

---

<a id="zc1685"></a>
### ZC1685 — Info: `sleep infinity` — container keep-alive pattern that ignores SIGTERM

**Severity:** `info`  
**Auto-fix:** `no`

`sleep infinity` is most often used as a container or systemd-unit keep-alive. Problem: GNU `sleep` does not install a SIGTERM handler, so when `docker stop` / `systemctl stop` sends SIGTERM the process sits unresponsive until the grace period expires and SIGKILL lands. The orchestrator reports a hung stop, logs look wrong, and any cleanup registered on signal handlers in a wrapping shell never runs. Replace with `exec tail -f /dev/null` (signal-handles cleanly) or front with `tini` / `dumb-init` when PID 1 must stay.

Disable by adding `ZC1685` to `disabled_katas` in `.zshellcheckrc`.

---

<a id="zc1686"></a>
### ZC1686 — Warn on `compinit -C` / `compinit -u` — skips / ignores `$fpath` integrity checks

**Severity:** `warning`  
**Auto-fix:** `no`

Zsh's completion system loads every file from `$fpath` as shell code. `compinit` normally warns when an `$fpath` directory (or a file in one) is writable by someone other than the current user or root, and skips loading. `compinit -C` skips the security check entirely for speed; `compinit -u` acknowledges the warning and loads the insecure files anyway. Either way, a world-writable entry in `$fpath` becomes an execution primitive for any user on the host. Audit `$fpath` with `compaudit`, fix ownership / permissions, then run plain `compinit`.

Disable by adding `ZC1686` to `disabled_katas` in `.zshellcheckrc`.

---

<a id="zc1687"></a>
### ZC1687 — Warn on `snap install --classic` / `--devmode` — weakens snap confinement

**Severity:** `warning`  
**Auto-fix:** `no`

`snap install --classic` drops the AppArmor / cgroup / seccomp sandbox entirely — the snap behaves like a normal Debian package with full system access. `--devmode` keeps the sandbox wired up but logs violations instead of blocking them. Both modes are documented escape hatches for snaps that cannot yet fit the strict confinement (IDEs, compilers, some network tooling), but in provisioning scripts they usually mean "I could not be bothered to pick a strict snap." Find a strict alternative, or install from the distro repository with proper AppArmor profiles; if `--classic` is truly required, document the specific snap and the interface that needed elevation.

Disable by adding `ZC1687` to `disabled_katas` in `.zshellcheckrc`.

---

<a id="zc1688"></a>
### ZC1688 — Warn on `aws s3 sync --delete` — destination objects deleted when source diverges

**Severity:** `warning`  
**Auto-fix:** `no`

`aws s3 sync SRC DST --delete` removes every object in DST that does not exist under SRC. A misspelled SRC, an empty build directory, or a stale `cd` turns the sync into a full-bucket wipe with no second confirmation and no recovery unless the bucket had versioning enabled. Restrict deletion to the prefix that really changed (`aws s3 sync ./build s3://bucket/app/ --delete`), add `--dryrun` behind a gate, or enable versioning and MFA-delete before running the command from a pipeline.

Disable by adding `ZC1688` to `disabled_katas` in `.zshellcheckrc`.

---

<a id="zc1689"></a>
### ZC1689 — Error on `borg delete --force` — forced deletion of backup archives or repository

**Severity:** `error`  
**Auto-fix:** `no`

`borg delete --force REPO[::ARCHIVE]` bypasses the confirmation prompt and removes the archive (or the whole repository, if ARCHIVE is omitted) in one go. Unlike `borg prune`, which keeps a retention ladder and logs what it would drop, `--force` deletion leaves nothing to restore from if the target was typed wrong. Keep scripts to `borg prune --keep-daily` / `--keep-within` with an explicit retention policy and gate any outright `borg delete` behind a human `--checkpoint-interval` review.

Disable by adding `ZC1689` to `disabled_katas` in `.zshellcheckrc`.

---

<a id="zc1690"></a>
### ZC1690 — Warn on `pip install git+<URL>` without a commit / tag pin

**Severity:** `warning`  
**Auto-fix:** `no`

`pip install git+https://host/repo[@main]` checks out a moving ref (the repository's default branch when no `@` suffix is given, otherwise a branch name the attacker can rewrite). Every subsequent install pulls whatever HEAD the branch currently points at — no lockfile, no checksum, no reproducibility. Pin to a specific commit SHA (`@abc1234…`) or a signed tag (`@v1.2.3`). If a proper PyPI release is available, drop the `git+` form entirely and install the versioned package.

Disable by adding `ZC1690` to `disabled_katas` in `.zshellcheckrc`.

---

<a id="zc1691"></a>
### ZC1691 — Warn on `rsync --remove-source-files` — SRC deletion tied to optimistic success

**Severity:** `warning`  
**Auto-fix:** `no`

`rsync --remove-source-files` deletes each source file once rsync has transferred it. The delete is gated on rsync's per-file success, which is generous: a remote out-of-disk error after the partial write, a `--chmod` rejection, or a flaky network that drops after the data bytes but before metadata can still look like success. Couple that with a wrong DST path and the source is gone with nothing to recover. Prefer a two-step flow: `rsync -a SRC DST` first, verify DST (checksums / file count), then `rm` the source explicitly, or use `mv` for local moves.

Disable by adding `ZC1691` to `disabled_katas` in `.zshellcheckrc`.

---

<a id="zc1692"></a>
### ZC1692 — Error on `kexec -e` — jumps into a new kernel without reboot, no audit trail

**Severity:** `error`  
**Auto-fix:** `no`

`kexec -e` transfers control to whatever kernel image is currently loaded via `kexec -l` — there is no firmware reboot, no init re-run, no chance for PAM / auditd / systemd hooks to record the transition. Malware uses it to pivot into a rootkit kernel while the audit log shows no reboot. If the intent is a fast reboot, prefer `systemctl kexec` (writes a wtmp entry and flushes filesystems), or just `reboot` / `systemctl reboot` and take the firmware cost for the audit trail.

Disable by adding `ZC1692` to `disabled_katas` in `.zshellcheckrc`.

---

<a id="zc1693"></a>
### ZC1693 — Warn on `ionice -c 1` — real-time I/O class starves every other disk consumer

**Severity:** `warning`  
**Auto-fix:** `no`

`ionice -c 1` (real-time I/O scheduling class) promotes the child above every best-effort (class 2) and idle (class 3) task queued against the same device. A busy workload — `rsync`, `dd`, database backup — then blocks sshd reads, systemd journal writes, and every other process until it yields, which for sequential I/O is effectively never. If the intent is "fast I/O", stay on class 2 and let CFQ / BFQ handle it; reserve class 1 for latency-critical paths launched by a scheduler that knows how to cap duration.

Disable by adding `ZC1693` to `disabled_katas` in `.zshellcheckrc`.

---

<a id="zc1694"></a>
### ZC1694 — Warn on `ssh -A` / `-o ForwardAgent=yes` — remote host can reuse local keys

**Severity:** `warning`  
**Auto-fix:** `no`

`ssh -A` (and `-o ForwardAgent=yes`) forwards the caller's `SSH_AUTH_SOCK` into the remote session. Anyone with root on the remote (and any process that shares its uid) can read the socket and impersonate the caller against every host the caller's keys unlock. Prefer `ssh -J JUMP HOST` (ProxyJump) for multi-hop access — it keeps the keys on the local side — or configure a scoped key for the remote task and copy it in with `ssh-copy-id`. Save key-forwarding for interactive use on trusted hosts.

Disable by adding `ZC1694` to `disabled_katas` in `.zshellcheckrc`.

---

<a id="zc1695"></a>
### ZC1695 — Warn on `terraform state rm` / `state push` — surgery on shared state outside plan/apply

**Severity:** `warning`  
**Auto-fix:** `no`

`terraform state rm RESOURCE` drops the resource from Terraform's tracking without touching the real cloud object — the next `terraform apply` sees it as newly-created and tries to re-provision, often hitting name-collision errors. `terraform state push FILE` replaces the entire remote state with a local file, bypassing locking and overwriting any concurrent changes. Both commands skirt the usual plan/apply audit trail. Reach for `terraform import` / `terraform apply -replace=ADDR` instead, and only run `state rm|push` from a reviewed fix-up PR with state backup in hand.

Disable by adding `ZC1695` to `disabled_katas` in `.zshellcheckrc`.

---

<a id="zc1696"></a>
### ZC1696 — Warn on `pnpm install --no-frozen-lockfile` / `yarn install --no-immutable` — CI lockfile drift

**Severity:** `warning`  
**Auto-fix:** `no`

`pnpm install --no-frozen-lockfile` (pnpm) and `yarn install --no-immutable` (yarn 4+) tell the package manager that the lockfile is merely a suggestion — any dep resolution change since the lockfile was written gets picked up silently. Run that from CI and the artifact no longer matches the pinned dependency graph reviewers signed off on. Use `pnpm install --frozen-lockfile` (the CI default) or `yarn install --immutable`, and let lockfile regen happen only from a dev workstation PR.

Disable by adding `ZC1696` to `disabled_katas` in `.zshellcheckrc`.

---

<a id="zc1697"></a>
### ZC1697 — Info: `cryptsetup open --allow-discards` — TRIM pass-through leaks free-sector map

**Severity:** `info`  
**Auto-fix:** `no`

`--allow-discards` tells dm-crypt to forward TRIM/DISCARD commands from the filesystem to the underlying SSD. The performance and wear-levelling gains are real, but so is the side effect: an attacker with raw-device access can read the free-sector map and see which blocks are empty — enough to fingerprint partition layouts, distinguish encrypted-full-volume from encrypted-sparse-content cases, and defeat plausible-deniability scenarios. If the threat model includes offline-disk inspection, drop `--allow-discards` and accept the perf hit; otherwise keep the flag but state the trade-off in the runbook.

Disable by adding `ZC1697` to `disabled_katas` in `.zshellcheckrc`.

---

<a id="zc1698"></a>
### ZC1698 — Warn on `fail2ban-client unban --all` / `stop` — wipes the active brute-force ban list

**Severity:** `warning`  
**Auto-fix:** `no`

`fail2ban-client unban --all` clears every active ban across every jail; `fail2ban-client stop` shuts the service down and flushes its rules. Either command restores network access for the exact attacker IPs `fail2ban` has already flagged as hostile — usually hundreds of known bots. Target a single IP with `fail2ban-client set <jail> unbanip <ip>` or reload a jail with `reload <jail>` when you only need to pick up new filter rules.

Disable by adding `ZC1698` to `disabled_katas` in `.zshellcheckrc`.

---

<a id="zc1699"></a>
### ZC1699 — Warn on `kubectl drain --delete-emptydir-data` — pod-local scratch data lost

**Severity:** `warning`  
**Auto-fix:** `no`

`kubectl drain NODE --delete-emptydir-data` (older alias `--delete-local-data`) lets drain evict pods that mount an `emptyDir` volume — the volume is deleted along with the pod, destroying any scratch data it held. Production clusters use `emptyDir` for caches, write-ahead logs, and scratch state that takes hours to rebuild. Confirm the pods on the node tolerate the loss (or migrate to a `persistentVolumeClaim`) before adding the flag; otherwise plan a controlled drain without it and accept the stuck-drain warning for the affected pods.

Disable by adding `ZC1699` to `disabled_katas` in `.zshellcheckrc`.

---

<a id="zc1700"></a>
### ZC1700 — Error on `ldapsearch -w PASSWORD` / `ldapmodify -w PASSWORD` — bind DN password in process list

**Severity:** `error`  
**Auto-fix:** `no`

OpenLDAP client tools (`ldapsearch`, `ldapmodify`, `ldapadd`, `ldapdelete`, `ldapmodrdn`, `ldappasswd`, `ldapcompare`) accept the bind password via `-w STRING`. Once invoked, the password sits in `/proc/PID/cmdline`, shell history, audit records, and any `ps` output — typically granting cn=admin / service-account bind over the whole directory. Use `-W` (prompt), `-y FILEPATH` (read from a mode-0400 file), or `SASL` auth (`-Y GSSAPI` with Kerberos) to keep the secret out of argv.

Disable by adding `ZC1700` to `disabled_katas` in `.zshellcheckrc`.

---

<a id="zc1701"></a>
### ZC1701 — Info: `dpkg -i FILE.deb` installs without automatic signature verification

**Severity:** `info`  
**Auto-fix:** `no`

Unlike `apt install`, which verifies package signatures against the apt repository's `Release.gpg`, plain `dpkg -i FILE.deb` applies the package with no integrity check beyond Debian's own `.deb` format. In a provisioning pipeline that downloaded the file over HTTPS from a vendor, that is usually fine — the TLS chain vouches for the bytes. In scripts that pick the file up from `/tmp`, `/var/tmp`, `/dev/shm`, or a mutable cache, a local user could swap the file between download and install. Verify with `sha256sum -c`, `debsig-verify`, or `dpkg-sig --verify` before invoking `dpkg -i`.

Disable by adding `ZC1701` to `disabled_katas` in `.zshellcheckrc`.

---

<a id="zc1702"></a>
### ZC1702 — Warn on `dpkg-reconfigure` without a noninteractive frontend — hangs in CI

**Severity:** `warning`  
**Auto-fix:** `no`

`dpkg-reconfigure PACKAGE` opens the package's debconf questions in whatever frontend the caller's `DEBIAN_FRONTEND` resolves to — typically a terminal dialog that blocks until someone presses a key. Inside a non-interactive pipeline (Dockerfile, Ansible task, cloud-init) the call hangs until the build times out. Pass `-f noninteractive` (or export `DEBIAN_FRONTEND=noninteractive` at the top of the script) and accept the debconf defaults; pre-seed any non-default answer with `debconf-set-selections`.

Disable by adding `ZC1702` to `disabled_katas` in `.zshellcheckrc`.

---

<a id="zc1703"></a>
### ZC1703 — Warn on `sysctl -w` disabling network-hardening knobs

**Severity:** `warning`  
**Auto-fix:** `no`

Several `net.ipv4.*` / `net.ipv6.*` sysctl knobs exist specifically to harden the host against on-link spoofing, ICMP redirect tampering, smurf amplification, and source-routed packets — `rp_filter=1`, `accept_source_route=0`, `accept_redirects=0`, `send_redirects=0`, `icmp_echo_ignore_broadcasts=1`, `log_martians=1`. Flipping any of them to the lax value (rp_filter=0, accept_source_route=1, …) re-opens classic layer-3 attacks. Leave the protective defaults in place; if a niche workload really needs relaxed filtering, scope the change per-interface with a comment explaining why.

Disable by adding `ZC1703` to `disabled_katas` in `.zshellcheckrc`.

---

<a id="zc1704"></a>
### ZC1704 — Error on `aws ec2 authorize-security-group-ingress --cidr 0.0.0.0/0` — port open to the internet

**Severity:** `error`  
**Auto-fix:** `no`

`aws ec2 authorize-security-group-ingress --cidr 0.0.0.0/0` (or `::/0` for IPv6) adds a rule that accepts the specified protocol/port from any source — the exact shape shodan, automated login-probers, and every exploit-as-a-service customer scans for. Restrict the source to the office CIDR, a VPN range, or a named security-group (`--source-group sg-…`). If the workload genuinely needs public access, front it with an ALB / API Gateway / CloudFront with WAF — not a raw SG rule from a shell script.

Disable by adding `ZC1704` to `disabled_katas` in `.zshellcheckrc`.

---

<a id="zc1705"></a>
### ZC1705 — Info: `awk -i inplace` is gawk-only — script breaks on mawk / BSD awk

**Severity:** `info`  
**Auto-fix:** `no`

The `inplace` extension that powers `awk -i inplace` ships only with gawk. On Alpine (default `mawk`), Debian-busybox, macOS, FreeBSD, NetBSD, OpenBSD, or any container image without `gawk` installed the script aborts with `fatal: can't open extension 'inplace'`. If portability matters, write through a temporary file (`awk … input > tmp && mv tmp input`); if you really do need in-place edits in scripts that target gawk only, document the requirement and add `command -v gawk >/dev/null` at the top.

Disable by adding `ZC1705` to `disabled_katas` in `.zshellcheckrc`.

---

<a id="zc1706"></a>
### ZC1706 — Error on `lvresize -L -SIZE` without `-r` — shrink without filesystem resize corrupts data

**Severity:** `error`  
**Auto-fix:** `no`

`lvresize -L -SIZE` (or `--size -SIZE`) shrinks the logical volume by SIZE bytes/extents. The filesystem on top still thinks it owns the original range; reads beyond the new LV end now return zeros, and the next write corrupts metadata. The `-r` (`--resizefs`) flag tells lvresize to call `fsadm` (which calls `resize2fs` / `xfs_growfs` / etc.) so the filesystem shrinks first. For ext4, always shrink the FS before the LV; for XFS, online shrink is impossible — back up, recreate, restore.

Disable by adding `ZC1706` to `disabled_katas` in `.zshellcheckrc`.

---

<a id="zc1707"></a>
### ZC1707 — Warn on `gpg --keyserver hkp://…` — plaintext keyserver fetch

**Severity:** `warning`  
**Auto-fix:** `no`

`hkp://` is the unencrypted HKP keyserver protocol. A MITM on the path (corporate proxy, hotel Wi-Fi, hostile router) can swap key bytes during the fetch and `gpg --recv-keys` happily imports the substitute. Use `hkps://keys.openpgp.org` (TLS) or fetch the armored key over HTTPS and verify the fingerprint out-of-band before `gpg --import`.

Disable by adding `ZC1707` to `disabled_katas` in `.zshellcheckrc`.

---

<a id="zc1708"></a>
### ZC1708 — Error on `find -L ... -delete` / `-exec rm` — symlink follow into unintended trees

**Severity:** `error`  
**Auto-fix:** `no`

`find -L` follows symlinks during traversal. Combined with `-delete` (or `-exec rm`), a symlink under the start path that points outside the intended root steers `find` into / `unlink`s files in `/etc`, `/var/lib`, or any other directory the symlink target reaches. Drop `-L` (the default `-P` keeps symlinks as objects), or restrict the walk with `-xdev`, `-mount`, and an explicit `-type f` test. For log-rotation pipes, `logrotate` is safer than a `find` one-liner.

Disable by adding `ZC1708` to `disabled_katas` in `.zshellcheckrc`.

---

<a id="zc1709"></a>
### ZC1709 — Error on `htpasswd -b USER PASSWORD` — basic-auth password in process list

**Severity:** `error`  
**Auto-fix:** `no`

`htpasswd -b FILE USER PASSWORD` (batch mode) takes the password as an argv slot. The cleartext sits in `/proc/PID/cmdline`, shell history, audit records, and any `ps` output. Use `htpasswd -i FILE USER` and pipe the secret on stdin (`printf %s "$pw" | htpasswd -i FILE USER`), or omit `-b` and `-i` so htpasswd prompts on the controlling TTY.

Disable by adding `ZC1709` to `disabled_katas` in `.zshellcheckrc`.

---

<a id="zc1710"></a>
### ZC1710 — Error on `journalctl --vacuum-size=1` / `--vacuum-time=1s` — journal-wipe pattern

**Severity:** `error`  
**Auto-fix:** `no`

`journalctl --vacuum-size=1` (down to 1 byte / 1K), `--vacuum-time=1s` (retain only the last second), or `--vacuum-files=1` (keep one journal file) effectively flushes the entire systemd journal. The classic shape after a compromise — clear the audit trail before re-enabling logging. Real retention belongs in `/etc/systemd/journald.conf` (`SystemMaxUse=`, `MaxRetentionSec=`), not in an ad-hoc one-shot. If you genuinely need to bound disk use, set the limit to a meaningful value (`--vacuum-time=2weeks`, `--vacuum-size=200M`).

Disable by adding `ZC1710` to `disabled_katas` in `.zshellcheckrc`.

---

<a id="zc1711"></a>
### ZC1711 — Error on `etcdctl del --prefix ""` / `--from-key ""` — wipes the entire keyspace

**Severity:** `error`  
**Auto-fix:** `no`

`etcdctl del --prefix KEY` deletes every key under KEY's range. With KEY empty (`""` or `"\0"`) the range is `["", "\xFF")` — the whole etcd cluster, including kube-apiserver state if etcd is the Kubernetes datastore. `--from-key ""` has the same effect for the lower-bound form. Restrict the prefix to the namespace you actually own (`/app/staging/`), or wrap the call with an explicit `etcdctl get --prefix --keys-only` review step.

Disable by adding `ZC1711` to `disabled_katas` in `.zshellcheckrc`.

---

<a id="zc1712"></a>
### ZC1712 — Error on `vault kv put PATH password=…` — secret value in process list

**Severity:** `error`  
**Auto-fix:** `no`

`vault kv put PATH key=value` (and the older `vault write PATH key=value`) put the value on the command line. When the key name screams secret (`password`, `secret`, `token`, `apikey`, `access_key`, `private_key`), the cleartext shows up in `ps`, `/proc/<pid>/cmdline`, shell history, and the audit log of the calling host — exactly the surface Vault is meant to remove. Use `key=@path/to/file` to read from disk, `key=-` to take the value on stdin, or `vault kv put -mount=secret PATH @secret.json` for a JSON payload.

Disable by adding `ZC1712` to `disabled_katas` in `.zshellcheckrc`.

---

<a id="zc1713"></a>
### ZC1713 — Error on `consul kv delete -recurse /` — wipes the entire Consul KV store

**Severity:** `error`  
**Auto-fix:** `no`

`consul kv delete -recurse PREFIX` removes every key under PREFIX. With PREFIX `/` (or an empty string) the command nukes the whole KV store, including service-discovery payloads, ACL bootstrap tokens, and any application-level config the cluster relies on. Scope the prefix to the app namespace (`consul kv delete -recurse /app/staging/`), confirm the keys you are about to lose with `consul kv get -recurse -keys`, and snapshot the datacenter (`consul snapshot save snap.bin`) before any large delete.

Disable by adding `ZC1713` to `disabled_katas` in `.zshellcheckrc`.

---

<a id="zc1714"></a>
### ZC1714 — Error on `gh repo delete --yes` / `gh release delete --yes` — bypassed confirmation

**Severity:** `error`  
**Auto-fix:** `no`

`gh repo delete OWNER/REPO --yes` (and `gh release delete TAG --yes`) skip the interactive confirmation that protects against typos and broken variable expansion. A repository deletion is final — issues, PRs, releases, GitHub Actions history, and (for free accounts) any forks against it all disappear with no soft-delete window. From a script, run without `--yes` so a human reviews the target, or wrap deletion in a manually-triggered workflow with explicit input prompts.

Disable by adding `ZC1714` to `disabled_katas` in `.zshellcheckrc`.

---

<a id="zc1715"></a>
### ZC1715 — Error on `read -p "prompt"` — Zsh `-p` reads from coprocess, not a prompt

**Severity:** `error`  
**Auto-fix:** `no`

Bash's `read -p "Prompt: " var` prints the prompt before reading. Zsh's `read -p` means "read from the coprocess set up with `coproc`" — when no coprocess exists, `read` errors with `no coprocess` and leaves the variable empty, silently breaking the script. The Zsh idiom is `read "var?Prompt: "` — a `?` after the variable name introduces the prompt string, with the same behavior under `-r`, `-s`, etc.

Disable by adding `ZC1715` to `disabled_katas` in `.zshellcheckrc`.

---

<a id="zc1716"></a>
### ZC1716 — Use Zsh `$CPUTYPE` / `$MACHTYPE` instead of `uname -m` / `-p`

**Severity:** `style`  
**Auto-fix:** `no`

Zsh maintains `$CPUTYPE` (e.g. `x86_64`, `aarch64`) and `$MACHTYPE` (the GNU triplet) as built-in parameters. Reading them is a constant-time parameter expansion, while `uname -m` / `uname -p` forks an external for the same answer. The Zsh values are populated at shell start from the same `uname(2)` call, so they stay in lockstep with what `uname` would print.

Disable by adding `ZC1716` to `disabled_katas` in `.zshellcheckrc`.

---

<a id="zc1717"></a>
### ZC1717 — Warn on `docker pull/push --disable-content-trust` — bypasses image signature checks

**Severity:** `warning`  
**Auto-fix:** `yes`

When `DOCKER_CONTENT_TRUST=1` is enforced on a host (or set via `/etc/docker/daemon.json`), Docker rejects unsigned image pulls and signs every push. The `--disable-content-trust` flag overrides that per command: a `pull` accepts a replaced or unsigned image into local storage, a `push` lands an unsigned tag in the registry where downstream pulls cannot verify provenance. Drop the flag and sign the artifact (`docker trust sign IMAGE:TAG`) instead, or scope the bypass with a tight Notary signer policy.

Disable by adding `ZC1717` to `disabled_katas` in `.zshellcheckrc`.

---

<a id="zc1718"></a>
### ZC1718 — Error on `gh secret set --body SECRET` / `-b SECRET` — secret in process list

**Severity:** `error`  
**Auto-fix:** `no`

`gh secret set NAME --body VALUE` (or `-b VALUE`, `--body=VALUE`) puts the secret on the command line. The cleartext appears in `ps`, `/proc/<pid>/cmdline`, shell history, and the audit log of the host running `gh`. Pipe the value via stdin (`gh secret set NAME < file`, `printf %s "$SECRET" | gh secret set NAME --body -`) or use `--body-file PATH` so the value never lands in argv.

Disable by adding `ZC1718` to `disabled_katas` in `.zshellcheckrc`.

---

<a id="zc1719"></a>
### ZC1719 — Warn on `git filter-branch` — deprecated since Git 2.24, use `git filter-repo`

**Severity:** `warning`  
**Auto-fix:** `no`

`git filter-branch` is deprecated as of Git 2.24; its manpage opens with "WARNING: this command is deprecated" and points users at `git filter-repo`. `filter-branch` is single-process slow, mishandles common cases (tag rewrites, refs/notes/*, signed commits), and leaves orphaned objects behind. The modern replacement is `git filter-repo` (separate package; `apt/brew install git-filter-repo`) — much faster, safer defaults, and what GitHub / GitLab guidance recommends for secret-removal rewrites.

Disable by adding `ZC1719` to `disabled_katas` in `.zshellcheckrc`.

---

<a id="zc1720"></a>
### ZC1720 — Use Zsh `$COLUMNS` / `$LINES` instead of `tput cols` / `tput lines`

**Severity:** `style`  
**Auto-fix:** `no`

Zsh tracks the terminal width and height in `$COLUMNS` and `$LINES`, updated automatically on `SIGWINCH`. Reading them is a constant-time parameter expansion, while `tput cols` / `tput lines` forks the terminfo helper on every call. Use the parameters; reach for `tput` only for terminfo queries Zsh does not surface as parameters.

Disable by adding `ZC1720` to `disabled_katas` in `.zshellcheckrc`.

---

<a id="zc1721"></a>
### ZC1721 — Error on `chmod NNN /dev/<node>` — world-writable device node is local privilege escalation

**Severity:** `error`  
**Auto-fix:** `no`

Granting world-write to a device node hands every local user a primitive: `/dev/kvm` becomes a host-root VM-exit gadget, `/dev/uinput` lets any user inject keystrokes into the active session, `/dev/loop-control` forges loop devices, `/dev/dri/cardN` opens GPU shaders for code-exec, `/dev/mem` / `/dev/kmem` (where still permitted) leak kernel state. Keep the kernel-managed default permissions; if userspace needs access, add a udev rule that grants it to a specific group, never `666` to the world.

Disable by adding `ZC1721` to `disabled_katas` in `.zshellcheckrc`.

---

<a id="zc1722"></a>
### ZC1722 — Warn on `ssh-keyscan HOST >> known_hosts` — TOFU bypass, blind-trust new host key

**Severity:** `warning`  
**Auto-fix:** `no`

`ssh-keyscan` fetches whatever host key the remote serves on its first reply. Appending the result straight to `known_hosts` is the exact step the host-key check is meant to defend against: a man-in-the-middle on first contact wins permanently. Pin the expected fingerprint via a side channel (vendor docs, prior verified contact) and assert it matches `ssh-keyscan HOST | ssh-keygen -lf -` before the append.

Disable by adding `ZC1722` to `disabled_katas` in `.zshellcheckrc`.

---

<a id="zc1723"></a>
### ZC1723 — Error on `gpg --delete-secret-keys` / `--delete-key` — irreversible key destruction

**Severity:** `error`  
**Auto-fix:** `no`

GPG key deletion is permanent. Once `--delete-secret-keys`, `--delete-secret-and-public-keys`, `--delete-keys`, or `--delete-key` removes the keyring entry there is no recovery short of a separate backup or off-card reimport. Combined with `--batch --yes`, the confirmation prompt is bypassed and a single accidental KEYID resolves to a one-shot wipe. Export the key first (`gpg --export-secret-keys --armor KEYID > backup.asc`, store offline) and never pair the delete flag with `--batch --yes` in automation.

Disable by adding `ZC1723` to `disabled_katas` in `.zshellcheckrc`.

---

<a id="zc1724"></a>
### ZC1724 — Warn on `pacman -Sy <pkg>` — partial upgrade, breaks dependency closure

**Severity:** `warning`  
**Auto-fix:** `no`

Arch Linux is rolling-release on the invariant that the local package database and the installed package set move together. `pacman -Sy <pkg>` refreshes the database and installs ONE package against the new metadata while every other installed package stays at its old version. The new package's dependency closure pulls libraries newer than what the rest of the system has, leaving a half-upgraded state that often manifests as `error while loading shared libraries`. Run a full `pacman -Syu` first, then install (`pacman -S <pkg>`); for CI use `pacman -Syu --noconfirm <pkg>` so the upgrade and install are atomic.

Disable by adding `ZC1724` to `disabled_katas` in `.zshellcheckrc`.

---

<a id="zc1725"></a>
### ZC1725 — Error on `cargo --token TOKEN` / `npm --otp CODE` — registry credential in process list

**Severity:** `error`  
**Auto-fix:** `no`

`cargo publish --token TOKEN` (and `cargo login`, `cargo owner`, `cargo yank`) puts the crates.io API token in argv — visible in `ps`, `/proc/<pid>/cmdline`, shell history, and CI logs. `npm publish --otp CODE` leaks the one-time code the same way. Use environment variables (`CARGO_REGISTRY_TOKEN`, `NPM_TOKEN`) or pipe via stdin (`cargo login --token -` reads from stdin), and source credentials from a secrets manager instead of the command line.

Disable by adding `ZC1725` to `disabled_katas` in `.zshellcheckrc`.

---

<a id="zc1726"></a>
### ZC1726 — Error on `gcloud ... delete --quiet` — silent destruction of GCP resources

**Severity:** `error`  
**Auto-fix:** `no`

`gcloud` accepts `--quiet` (`-q`) globally to suppress every confirmation prompt. Combined with `delete` on projects, SQL instances, GKE clusters, compute VMs, secrets, or storage buckets, a single misresolved variable wipes the resource with no human-in-the-loop. Project deletion has a 30-day soft window but compute disks, secrets, and BigQuery tables are gone immediately. Drop `--quiet` from delete commands or route the bulk-destroy through a Terraform plan that surfaces the diff for review.

Disable by adding `ZC1726` to `disabled_katas` in `.zshellcheckrc`.

---

<a id="zc1727"></a>
### ZC1727 — Error on `curl/wget --proxy http://USER:PASS@HOST` — proxy credentials in argv

**Severity:** `error`  
**Auto-fix:** `no`

Embedding the proxy username and password in the URL passed to `--proxy` (curl), `-x` (curl short form), or `--proxy-password=` (wget) lands the credential in argv — visible in `ps`, `/proc/<pid>/cmdline`, shell history, and CI logs. Configure the proxy through `~/.curlrc` / `~/.netrc` (chmod 600) for curl, or `~/.wgetrc` for wget, so the secret never reaches the command line.

Disable by adding `ZC1727` to `disabled_katas` in `.zshellcheckrc`.

---

<a id="zc1728"></a>
### ZC1728 — Error on `pip install --index-url http://...` — plaintext index allows MITM

**Severity:** `error`  
**Auto-fix:** `no`

`pip install --index-url http://...`, `--extra-index-url http://...`, and `-i http://...` tell pip to fetch packages over plaintext HTTP. Any network-position attacker (open Wi-Fi, hostile transit, MITM proxy) can replace package metadata or wheel contents in flight — direct code execution on the install host. Switch to `https://`, or on internal networks terminate TLS at the mirror and only configure the `https://` URL.

Disable by adding `ZC1728` to `disabled_katas` in `.zshellcheckrc`.

---

<a id="zc1729"></a>
### ZC1729 — Error on `ip route flush all` / `ip route del default` — script loses network connectivity

**Severity:** `error`  
**Auto-fix:** `no`

`ip route flush all` (or `flush table main`) wipes every routing entry, including the default gateway. `ip route del default` removes only the default route — same outcome. The remote SSH session that just ran the command can no longer talk to the host, and any subsequent step that needs the network hangs until manual console intervention. Scope the flush (`flush dev <iface>`, `flush scope link`) or use `ip route replace default via <gw>` so the new route is in place before the old one disappears.

Disable by adding `ZC1729` to `disabled_katas` in `.zshellcheckrc`.

---

<a id="zc1730"></a>
### ZC1730 — Warn on `brew install --HEAD <pkg>` — pulls upstream HEAD, no version stability

**Severity:** `warning`  
**Auto-fix:** `no`

`brew install --HEAD <pkg>` (also `reinstall --HEAD`, `upgrade --HEAD`) builds the formula from the upstream source repository's HEAD branch. The build is unrepeatable — every run pulls a different commit — and any compromised upstream commit lands directly on the install host. Pin to a stable release of the formula, or if HEAD is genuinely required, vendor the build into a private tap that fixes a specific revision.

Disable by adding `ZC1730` to `disabled_katas` in `.zshellcheckrc`.

---

<a id="zc1731"></a>
### ZC1731 — Error on `curl -d 'password=…'` / `wget --post-data='token=…'` — secret in argv

**Severity:** `error`  
**Auto-fix:** `no`

`curl -d` / `--data` / `--data-raw` / `--data-urlencode` and `wget --post-data` / `--body-data` put the POST body in argv — visible in `ps`, `/proc/<pid>/cmdline`, shell history, and CI logs. When the body contains a credential-looking key (`password`, `secret`, `token`, `apikey`, `access_key`, `private_key`), the secret leaks the same way an inline `-u user:pass` would. Read the value from a file (`curl --data @secret.txt URL`, `--data-binary @-` piped from a secrets store) so the secret never reaches the command line.

Disable by adding `ZC1731` to `disabled_katas` in `.zshellcheckrc`.

---

<a id="zc1732"></a>
### ZC1732 — Warn on `flatpak override --filesystem=host` — removes Flatpak sandbox isolation

**Severity:** `warning`  
**Auto-fix:** `no`

Flatpak's primary security guarantee is filesystem sandboxing — apps see only their own data plus paths the user explicitly grants via portals. `flatpak override --filesystem=host` (also `host-os`, `host-etc`, `home`, `/`) persistently grants the app unrestricted read/write to the host filesystem at every subsequent run. Same risk applies to `flatpak run --filesystem=host`. Grant the specific subdirectory the app actually needs (`--filesystem=~/Documents:ro`) or rely on Filesystem portals so the user picks paths interactively per session.

Disable by adding `ZC1732` to `disabled_katas` in `.zshellcheckrc`.

---

<a id="zc1733"></a>
### ZC1733 — Error on `docker plugin install --grant-all-permissions` — accepts every requested cap

**Severity:** `error`  
**Auto-fix:** `no`

Docker plugins run as root with whatever privileges they ask for at install time — host networking, `/dev/*` mounts, arbitrary capability grants. The interactive prompt enumerates each request so the operator can refuse anything unexpected. `--grant-all-permissions` skips the prompt and accepts the whole list, so a compromised plugin author or a typo-squatted name owns the host on first install. Install plugins by name, walk the prompt manually, then pin the tag (`@sha256:...`) once vetted.

Disable by adding `ZC1733` to `disabled_katas` in `.zshellcheckrc`.

---

<a id="zc1734"></a>
### ZC1734 — Error on `cp/mv/tee` overwriting `/etc/passwd|shadow|group|gshadow`

**Severity:** `error`  
**Auto-fix:** `no`

The user-identity files are managed by `useradd` / `usermod` / `vipw` / `vigr`, which take a file lock and keep `passwd` / `shadow` (and `group` / `gshadow`) in sync. Replacing them with `cp`, `mv`, `tee`, or a redirect (`echo … > /etc/passwd`) bypasses the lock: concurrent edits race, malformed entries lock the whole system out, and the shadow file ends up pointing at users that no longer exist. Use `vipw -e` / `vigr -e` to edit, or `useradd` / `usermod` / `passwd` to mutate one entry at a time.

Disable by adding `ZC1734` to `disabled_katas` in `.zshellcheckrc`.

---

<a id="zc1735"></a>
### ZC1735 — Error on `efibootmgr -B` — deletes UEFI boot entry, may brick boot

**Severity:** `error`  
**Auto-fix:** `no`

`efibootmgr -B` deletes the currently-selected UEFI boot entry; combined with `-b BOOTNUM` it removes the specific entry instead. If that entry was the only viable bootloader (or the firmware's removable-media fallback is not present), the next reboot drops into the UEFI shell or picks an unexpected device — recovery needs console access. Run `efibootmgr -v` first to inspect `BootOrder`, ensure a fallback (`/EFI/BOOT/BOOTX64.EFI`) is in place, and prefer `efibootmgr -o NEW,ORDER` to demote rather than delete.

Disable by adding `ZC1735` to `disabled_katas` in `.zshellcheckrc`.

---

<a id="zc1736"></a>
### ZC1736 — Error on `pulumi destroy --yes` / `up --yes` — silent infra mutation in CI

**Severity:** `error`  
**Auto-fix:** `no`

`pulumi destroy --yes` (or `-y`) skips the preview-and-confirm step that normally surfaces every resource scheduled for deletion. A single misresolved stack name or wrong AWS credential resolves to a one-shot wipe of cloud infrastructure. `pulumi up --yes` and `pulumi refresh --yes` carry the same footgun for resource creation/replacement. Pipe `pulumi preview` output into a review step (manual approval, GitHub Actions environment protection rule) before applying, and never combine `--yes` with the `destroy` verb in automation.

Disable by adding `ZC1736` to `disabled_katas` in `.zshellcheckrc`.

---

<a id="zc1737"></a>
### ZC1737 — Error on `wpa_passphrase SSID PASSWORD` — Wi-Fi passphrase in process list

**Severity:** `error`  
**Auto-fix:** `no`

`wpa_passphrase SSID PASSPHRASE` generates `wpa_supplicant.conf` content on stdout. Putting PASSPHRASE on the command line lands it in `ps`, `/proc/<pid>/cmdline`, shell history, and the audit log of every local user that can list processes. Drop the second positional argument and let `wpa_passphrase SSID < /run/secrets/wifi` (or piped via stdin from a secrets store) read the passphrase from a file descriptor instead.

Disable by adding `ZC1737` to `disabled_katas` in `.zshellcheckrc`.

---

<a id="zc1738"></a>
### ZC1738 — Error on `aws rds delete-db-instance --skip-final-snapshot` — DB destroyed unrecoverable

**Severity:** `error`  
**Auto-fix:** `no`

RDS keeps a final snapshot when an instance or cluster is deleted — the only path back from a typo'd identifier or wrong account. `--skip-final-snapshot` opts out of that snapshot, so the database is gone the moment the API call returns; same applies to `aws rds delete-db-cluster --skip-final-snapshot`. Drop the flag (or pass `--final-db-snapshot-identifier <name>` so the snapshot name is explicit) and verify the snapshot lands before reusing the identifier.

Disable by adding `ZC1738` to `disabled_katas` in `.zshellcheckrc`.

---

<a id="zc1739"></a>
### ZC1739 — Warn on `git submodule update --remote` — pulls upstream HEAD, breaks reproducibility

**Severity:** `warning`  
**Auto-fix:** `no`

`git submodule update --remote` fetches each submodule's tracked branch HEAD instead of the commit pinned in the parent repo's index. Builds become non-reproducible — every CI run pulls a different commit — and any compromised upstream commit lands directly in the build. Use `git submodule update --init --recursive` (defaults to the pinned commit) and bump submodule pins through reviewed PRs.

Disable by adding `ZC1739` to `disabled_katas` in `.zshellcheckrc`.

---

<a id="zc1740"></a>
### ZC1740 — Warn on `gh release upload --clobber` — silent overwrite of release asset

**Severity:** `warning`  
**Auto-fix:** `no`

`gh release upload TAG FILE --clobber` replaces an existing asset with the same name without prompting. In production this is how a release artifact gets silently downgraded — a CI job re-runs with a stale build and the user-facing download moves backward without anyone noticing. Drop `--clobber` so the second upload errors out, or version the asset name (`mytool-1.2.3-linux.tar.gz` instead of `mytool-linux.tar.gz`) so each upload has a unique slot.

Disable by adding `ZC1740` to `disabled_katas` in `.zshellcheckrc`.

---

<a id="zc1741"></a>
### ZC1741 — Error on `mkpasswd PASSWORD` — clear-text password in process list

**Severity:** `error`  
**Auto-fix:** `no`

`mkpasswd PASSWORD` (whatwg/Debian `whois`-package version) and `mkpasswd -m METHOD PASSWORD` hash the password and print the crypt(3) string on stdout. Putting PASSWORD on the command line lands it in `ps`, `/proc/<pid>/cmdline`, shell history, and the host audit log. Drop the positional password and read from stdin (`mkpasswd -s` reads the password from stdin) — pipe the secret from a credentials file or vault: `printf %s "$PASSWORD" | mkpasswd -s`.

Disable by adding `ZC1741` to `disabled_katas` in `.zshellcheckrc`.

---

<a id="zc1742"></a>
### ZC1742 — Error on `mc alias set NAME URL ACCESS_KEY SECRET_KEY` — S3 keys in process list

**Severity:** `error`  
**Auto-fix:** `no`

MinIO's `mc alias set NAME URL ACCESS_KEY SECRET_KEY` (also `mc config host add ALIAS URL ACCESS SECRET` on legacy versions) accepts the S3 access and secret keys as positional arguments. Both land in argv — visible in `ps`, `/proc/<pid>/cmdline`, shell history, and CI logs. Drop the trailing keys and let `mc alias set NAME URL` prompt for them, or use the `MC_HOST_<alias>=https://ACCESS:SECRET@host` env-var form scoped to a single command and unset immediately after.

Disable by adding `ZC1742` to `disabled_katas` in `.zshellcheckrc`.

---

<a id="zc1743"></a>
### ZC1743 — Warn on `npm audit fix --force` — accepts major-version dependency bumps silently

**Severity:** `warning`  
**Auto-fix:** `no`

`npm audit fix --force` (and `pnpm audit --fix --force`) resolves advisories by upgrading dependencies past semver-major boundaries when no backward-compatible patch exists. The flag accepts every upgrade without surfacing the breaking changes — a build can silently move to a new major of a transitive dependency that removes APIs your code calls. Drop `--force` and triage each advisory individually; `npm audit fix` handles compatible patches, and the remaining advisory targets need a pin or a vendored fork.

Disable by adding `ZC1743` to `disabled_katas` in `.zshellcheckrc`.

---

<a id="zc1744"></a>
### ZC1744 — Warn on `kubectl port-forward --address 0.0.0.0` — cluster port exposed to every interface

**Severity:** `warning`  
**Auto-fix:** `no`

`kubectl port-forward` defaults to binding the local end of the tunnel on `127.0.0.1`. `--address 0.0.0.0` (or a specific non-loopback IP) exposes the target pod's port to every interface on the developer's workstation or the bastion host running the command. Anyone on the LAN / VPN can reach internal cluster services that never meant to be externally reachable. Drop the flag (loopback default), or pick a specific interface that is already scoped to a trusted network.

Disable by adding `ZC1744` to `disabled_katas` in `.zshellcheckrc`.

---

<a id="zc1745"></a>
### ZC1745 — Error on `poetry publish --password PASS` / `twine upload -p PASS` — registry secret in argv

**Severity:** `error`  
**Auto-fix:** `no`

Poetry's `publish --username USER --password PASS` and Twine's `upload --username USER --password PASS` (or the short `-u`/`-p` forms) put the PyPI / private-index password in argv — visible in `ps`, `/proc/<pid>/cmdline`, shell history, and CI logs. Use the `POETRY_PYPI_TOKEN_<NAME>` / `TWINE_USERNAME` + `TWINE_PASSWORD` environment variables (sourced from a secrets manager) or a `~/.pypirc` file with `0600` perms so the credential never reaches the command line.

Disable by adding `ZC1745` to `disabled_katas` in `.zshellcheckrc`.

---

<a id="zc1746"></a>
### ZC1746 — Error on `sysctl -w kernel.randomize_va_space=0|1` — weakens or disables ASLR

**Severity:** `error`  
**Auto-fix:** `no`

`kernel.randomize_va_space` controls Address Space Layout Randomization. Value `2` (default) randomizes stack, heap, VDSO, and mmap regions; value `1` omits the heap; value `0` disables ASLR entirely, making every memory layout deterministic. Exploits that rely on absolute addresses — stack overflows, ROP chains, kernel gadgets — become one-shot instead of brute-forceable. Never lower this below `2` outside a sandboxed kernel-debug context.

Disable by adding `ZC1746` to `disabled_katas` in `.zshellcheckrc`.

---

<a id="zc1747"></a>
### ZC1747 — Error on `npm/yarn/pnpm --registry http://...` — plaintext registry allows MITM

**Severity:** `error`  
**Auto-fix:** `no`

`npm install --registry http://...`, `pnpm --registry http://...`, and `yarn config set registry http://...` configure a plaintext HTTP package registry. Any network-position attacker (open Wi-Fi, hostile transit, MITM proxy) can replace tarball metadata or content in flight; npm install-time `postinstall` scripts then execute the swapped code on the build host. Switch the registry URL to `https://` (or terminate TLS at the internal mirror) and pair it with a lockfile to pin tarball integrity hashes.

Disable by adding `ZC1747` to `disabled_katas` in `.zshellcheckrc`.

---

<a id="zc1748"></a>
### ZC1748 — Error on `helm repo add NAME http://...` — plaintext chart repo allows MITM

**Severity:** `error`  
**Auto-fix:** `no`

`helm repo add NAME http://URL` registers a chart repository reached over plaintext HTTP. Any network-position attacker can swap `index.yaml` or a chart tarball in flight, and subsequent `helm install` pulls container images and Kubernetes manifests straight from the substituted content — fast path to cluster-wide code execution. Use `https://`, and pair it with chart provenance (`helm install --verify` or OCI signatures) to pin the digest.

Disable by adding `ZC1748` to `disabled_katas` in `.zshellcheckrc`.

---

<a id="zc1749"></a>
### ZC1749 — Error on `virsh undefine DOMAIN --remove-all-storage` — wipes VM disk images

**Severity:** `error`  
**Auto-fix:** `no`

`virsh undefine DOMAIN --remove-all-storage` (also `--wipe-storage` and the newer `--storage <vol,vol>`) removes the VM's configuration AND deletes every disk image the domain references. There is no soft-delete and no recycle bin — a misresolved DOMAIN or a shared storage pool turns one typo into data loss across VMs that happened to share a snapshot chain. Split the operation: back up the qcow2 images (`virsh vol-clone` or `qemu-img convert`), then `virsh undefine` without the storage flags, then delete volumes deliberately with `virsh vol-delete` after a review.

Disable by adding `ZC1749` to `disabled_katas` in `.zshellcheckrc`.

---

<a id="zc1750"></a>
### ZC1750 — Error on `kubectl proxy --address 0.0.0.0` — cluster API proxy on every interface

**Severity:** `error`  
**Auto-fix:** `no`

`kubectl proxy` tunnels Kubernetes API requests authenticated with the local kubeconfig's credentials. Defaults bind to `127.0.0.1` and accept only `localhost` hosts. `--address 0.0.0.0` (or a specific non-loopback IP) exposes that tunnel to every interface on the workstation / bastion, so anyone on the LAN or VPN gets the cluster admin the kubeconfig holds. Same risk applies to `--accept-hosts '.*'`. Keep the loopback default and scope with SSH port forwarding, or restrict `--address` to an interface behind a tight firewall.

Disable by adding `ZC1750` to `disabled_katas` in `.zshellcheckrc`.

---

<a id="zc1751"></a>
### ZC1751 — Error on `rpm/dnf/yum remove --nodeps` — bypasses dependency check, breaks dependents

**Severity:** `error`  
**Auto-fix:** `no`

`rpm -e --nodeps PKG` (also `dnf remove --nodeps`, `yum remove --nodeps`, `zypper remove --force`) removes the package while skipping the dependency solver. Anything transitively depending on the target immediately breaks — `libc`, `openssl`, `systemd` units, even `dnf` itself can get pulled out, leaving the host unbootable or unpackageable. Resolve the dependency conflict explicitly (`dnf swap`, `rpm -e --rebuilddb` never, pin the conflicting package) instead of bypassing the check.

Disable by adding `ZC1751` to `disabled_katas` in `.zshellcheckrc`.

---

<a id="zc1752"></a>
### ZC1752 — Error on `pvcreate/vgcreate/lvcreate -ff|--yes` — force-init LVM over existing data

**Severity:** `error`  
**Auto-fix:** `no`

LVM prompts before overwriting existing filesystem, RAID, or LVM signatures on a device — that prompt is the only thing saving you from a typo'd target destroying someone else's data. `pvcreate -ff`, `pvcreate --yes`, and the same flags on `vgcreate` / `lvcreate` skip the prompt. Drop the flag, inspect with `wipefs -n` and `lsblk -f` first, then confirm the target before re-running the create command.

Disable by adding `ZC1752` to `disabled_katas` in `.zshellcheckrc`.

---

<a id="zc1753"></a>
### ZC1753 — Error on `rclone purge REMOTE:PATH` — bulk delete of every object under the remote path

**Severity:** `error`  
**Auto-fix:** `no`

`rclone purge REMOTE:PATH` removes every object and empty directory under PATH on the remote — no dry-run gate, no confirmation, no soft-delete unless the backend happens to version. A typo'd path or a stale variable turns one line into a bucket-wide wipe (S3, GCS, Azure, Swift all honour the same API call). Preview with `rclone lsf REMOTE:PATH` or `rclone delete --dry-run`, then use `rclone delete` scoped narrower; enable object versioning on the backend so a bad run can roll back.

Disable by adding `ZC1753` to `disabled_katas` in `.zshellcheckrc`.

---

<a id="zc1754"></a>
### ZC1754 — Error on `gh auth status -t` / `--show-token` — prints OAuth token to stdout

**Severity:** `error`  
**Auto-fix:** `no`

`gh auth status -t` (alias `--show-token`) prints the stored GitHub OAuth token alongside the status summary. In CI logs, shared terminals, piped to `less`/`tee`, or captured via `script`, the token ends up on disk or in scrollback where anyone with log access becomes repo-admin. Never combine `-t` with `auth status` in automation; if a machine-readable token is needed, `gh auth token` prints only the token and makes the secret-handling path explicit.

Disable by adding `ZC1754` to `disabled_katas` in `.zshellcheckrc`.

---

<a id="zc1755"></a>
### ZC1755 — Error on `gcloud sql users {create,set-password} --password PASS` — DB password in argv

**Severity:** `error`  
**Auto-fix:** `no`

`gcloud sql users create USER --instance INST --password PASS` (and the `set-password` variant) place the Cloud SQL user password on the command line — visible in `ps`, `/proc/<pid>/cmdline`, shell history, and CI logs, and stored in Cloud Audit Logs' request payload. Use `--prompt-for-password` (interactive) or generate the password server-side in Secret Manager and post to the SQL Admin API via `gcloud auth print-access-token` piped to `curl` with the body sourced from a file.

Disable by adding `ZC1755` to `disabled_katas` in `.zshellcheckrc`.

---

<a id="zc1756"></a>
### ZC1756 — Error on `chmod NNN /run/docker.sock` — world access is root-equivalent privesc

**Severity:** `error`  
**Auto-fix:** `no`

Container-runtime sockets (`/var/run/docker.sock`, `/run/containerd/containerd.sock`, `/run/crio/crio.sock`, `/run/podman/podman.sock`) accept commands that run on the host with root privilege — starting privileged containers, mounting the host filesystem, reading every file on disk. Making the socket world-readable or world-writable (`chmod 644/660/666/777`) hands every local user that root-escalation primitive. Keep the socket `0660 root:docker` (or the equivalent runtime group) and add only trusted accounts to that group.

Disable by adding `ZC1756` to `disabled_katas` in `.zshellcheckrc`.

---

<a id="zc1757"></a>
### ZC1757 — Warn on `gh auth refresh --scopes delete_repo|admin:*` — token escalated to destructive perms

**Severity:** `warning`  
**Auto-fix:** `no`

`gh auth refresh --scopes <list>` (also `gh auth login --scopes`) rotates the stored OAuth token with additional scopes. `delete_repo`, `admin:org`, `admin:enterprise`, `admin:public_key`, and `admin:*_hook` give the token permanent destructive perms that outlast the script that asked for them — a compromised token now carries repo-deletion, org-membership, and SSH-key manipulation rights. Request the minimum scope the task needs (`repo`, `workflow`) and rotate the token off when the elevated operation completes.

Disable by adding `ZC1757` to `disabled_katas` in `.zshellcheckrc`.

---

<a id="zc1758"></a>
### ZC1758 — Warn on `gh codespace delete --force` — destroys codespace with uncommitted work

**Severity:** `warning`  
**Auto-fix:** `no`

`gh codespace delete --force` (alias `-f`) skips the confirmation prompt and deletes the target codespace along with any uncommitted, unpushed, or unstaged work inside it. Combined with `--all`, one line wipes every codespace on the account. Drop the flag, let the prompt enumerate what is about to go, and only confirm after verifying no local state would be lost — `git status` / `git stash list` inside the codespace first.

Disable by adding `ZC1758` to `disabled_katas` in `.zshellcheckrc`.

---

<a id="zc1759"></a>
### ZC1759 — Error on `vault login TOKEN` / `login -method=… password=…` — credential in process list

**Severity:** `error`  
**Auto-fix:** `no`

Vault accepts credentials on its `login` / `auth` subcommands in two argv-leaking shapes: a positional token (`vault login <TOKEN>`) and KEY=VALUE pairs for non-token methods (`vault login -method=userpass username=U password=P`). Both land the secret in `ps`, `/proc/<pid>/cmdline`, shell history, and Vault's audit log request payload. Read the token from stdin (`vault login -` with `printf %s "$TOKEN" |`) or source `VAULT_TOKEN` from a secrets file and run `vault login -method=token`.

Disable by adding `ZC1759` to `disabled_katas` in `.zshellcheckrc`.

---

<a id="zc1760"></a>
### ZC1760 — Warn on `openssl rand -hex|-base64 N` with N < 16 — generated value too short

**Severity:** `warning`  
**Auto-fix:** `no`

`openssl rand -hex N` (and `-base64 N`) outputs N random bytes encoded into the requested form. N below 16 (128 bits) produces a value short enough that an attacker with modest GPU resources can brute-force it offline — too weak for passwords, API tokens, reset URLs, or any other secret that sits at rest. Use `-hex 32` (256-bit) for secrets and long-lived tokens; `-hex 16` is acceptable only for short-validity nonces paired with rate-limited consumers.

Disable by adding `ZC1760` to `disabled_katas` in `.zshellcheckrc`.

---

<a id="zc1761"></a>
### ZC1761 — Warn on `gh gist create --public` — file becomes world-visible and indexed on GitHub

**Severity:** `warning`  
**Auto-fix:** `no`

`gh gist create --public FILE` (alias `-p`) creates the gist with `public: true`. Public gists are listed on `gist.github.com/discover`, crawled by search engines, and archived by secondary scrapers — a leaked secret, private company snippet, or unreleased note is effectively permanent the moment it lands. The default (`public: false`) keeps the gist unlisted and reachable only via its URL. Drop `--public` unless public exposure is the explicit goal.

Disable by adding `ZC1761` to `disabled_katas` in `.zshellcheckrc`.

---

<a id="zc1762"></a>
### ZC1762 — Error on `kubeadm join --discovery-token-unsafe-skip-ca-verification` — cluster CA not checked

**Severity:** `error`  
**Auto-fix:** `no`

`kubeadm join` verifies the control-plane API server's CA before accepting the kubelet bootstrap token. `--discovery-token-unsafe-skip-ca-verification` skips that check, so a network-position attacker can impersonate the API server, harvest the bootstrap token, and seed malicious workloads onto the joining node. Always pin the CA with `--discovery-token-ca-cert-hash sha256:<digest>` (emitted by `kubeadm token create --print-join-command`) or supply a kubeconfig discovery file that has the CA baked in.

Disable by adding `ZC1762` to `disabled_katas` in `.zshellcheckrc`.

---

<a id="zc1763"></a>
### ZC1763 — Error on `docker compose down -v` / `--volumes` — wipes named volumes (data loss)

**Severity:** `error`  
**Auto-fix:** `no`

`docker compose down -v` (alias `--volumes`, equivalent in `docker-compose down -v`) tears the stack down AND deletes every named volume declared in the compose file. Database contents, cache state, uploaded assets, and any other volume-backed data goes with them — there is no soft-delete. Drop the flag in CI and production scripts; keep it only for throwaway local testbeds where losing volume state is intentional.

Disable by adding `ZC1763` to `disabled_katas` in `.zshellcheckrc`.

---

<a id="zc1764"></a>
### ZC1764 — Warn on `git commit --no-verify` / `-n` — skips pre-commit and commit-msg hooks

**Severity:** `warning`  
**Auto-fix:** `no`

`git commit --no-verify` (alias `-n`) bypasses both the pre-commit and commit-msg hooks, which are often the last guardrail against leaked secrets, formatting drift, or failing tests. The flag is usually a symptom of a hook that needs fixing rather than silencing — the exception quickly becomes the rule. Fix the blocking hook, carve out a narrow per-file exemption in the hook itself, or file a tracked issue, instead of adding `--no-verify` to every commit in a script.

Disable by adding `ZC1764` to `disabled_katas` in `.zshellcheckrc`.

---

<a id="zc1765"></a>
### ZC1765 — Error on `snap remove --purge SNAP` — skips the automatic data snapshot

**Severity:** `error`  
**Auto-fix:** `no`

`snap remove SNAP` takes a snapshot of every writable area (`$SNAP_DATA`, `$SNAP_USER_DATA`, `$SNAP_COMMON`) before uninstalling, so the data can later be restored with `snap restore`. `--purge` skips that snapshot: the snap is gone along with every file it owned, and snapd has no record to roll back. Drop `--purge` unless the snap's data is genuinely disposable; otherwise `snap save SNAP` first, capture the set ID, and only then remove.

Disable by adding `ZC1765` to `disabled_katas` in `.zshellcheckrc`.

---

<a id="zc1766"></a>
### ZC1766 — Error on `memcached -l 0.0.0.0` — memcached exposed on every interface

**Severity:** `error`  
**Auto-fix:** `no`

`memcached -l 0.0.0.0` (or `::`, `--listen=0.0.0.0`) binds memcached's TCP listener to every interface on the host. Memcached has no authentication and, before `-U 0` became default, its UDP handler was the largest DDoS-amplification vector on the internet. Bind to `127.0.0.1` or a private-network IP only, and put memcached behind a firewall / security group scoped to the application that consumes it.

Disable by adding `ZC1766` to `disabled_katas` in `.zshellcheckrc`.

---

<a id="zc1767"></a>
### ZC1767 — Error on `mongod --bind_ip 0.0.0.0` — MongoDB exposed on every interface

**Severity:** `error`  
**Auto-fix:** `no`

`mongod --bind_ip 0.0.0.0` (or `::`) binds MongoDB's listener to every interface on the host. Combined with no-auth defaults (pre-3.4) or a wildcard database user, this was the source of the 2017 ransomware wave that wiped tens of thousands of public MongoDB instances. Bind to `127.0.0.1` or a private-network IP, enable authentication with `--auth`, and firewall port `27017`.

Disable by adding `ZC1767` to `disabled_katas` in `.zshellcheckrc`.

---

<a id="zc1768"></a>
### ZC1768 — Error on `sqlcmd -P PASSWORD` / `bcp -P PASSWORD` — SQL Server password in argv

**Severity:** `error`  
**Auto-fix:** `no`

Microsoft's SQL Server CLI tools (`sqlcmd`, `bcp`, `osql`) accept the password via `-P PASSWORD` as a positional argument value. The password lands in argv — visible in `ps`, `/proc/<pid>/cmdline`, shell history, CI logs, and SQL Server's audit trace for the session. Use `-P` with no value (prompts), or read the password from the environment variable `SQLCMDPASSWORD` (sourced from a secrets file). On modern sqlcmd, `-G` + Azure AD integration avoids the password altogether.

Disable by adding `ZC1768` to `disabled_katas` in `.zshellcheckrc`.

---

<a id="zc1769"></a>
### ZC1769 — Warn on `vagrant destroy --force` — VM destroyed without confirmation

**Severity:** `warning`  
**Auto-fix:** `no`

`vagrant destroy --force` (alias `-f`) tears every VM in the Vagrantfile down — and their ephemeral filesystem state — without prompting. Any data provisioned into the VM that was never exported back to the host (database seeds, build caches, local-only test fixtures) goes with it. In unattended scripts, drop the flag so the prompt still gates the destroy; for CI cycles, `vagrant halt` + `vagrant up` reuses the same box without losing state.

Disable by adding `ZC1769` to `disabled_katas` in `.zshellcheckrc`.

---

<a id="zc1770"></a>
### ZC1770 — Warn on `gpg --always-trust` / `--trust-model always` — bypasses Web-of-Trust

**Severity:** `warning`  
**Auto-fix:** `no`

`gpg --always-trust` (equivalent to `--trust-model always`) accepts every key in the keyring as fully trusted, regardless of signatures from the owner or any introducer. A signature made by an attacker-controlled key pair that was imported with no further vetting will verify cleanly. In automation this turns signature verification into a presence check — any key bundled with the payload satisfies `gpg --verify`. Remove the flag and build a proper trust path: either mark the expected signer key trusted once (`gpg --edit-key KEYID trust`), or pin the expected fingerprint and match it against the signer after `gpg --verify --status-fd 1`.

Disable by adding `ZC1770` to `disabled_katas` in `.zshellcheckrc`.

---

<a id="zc1771"></a>
### ZC1771 — Warn on `alias -g` / `alias -s` — global and suffix aliases surprise script readers

**Severity:** `warning`  
**Auto-fix:** `no`

`alias -g NAME=value` defines a global alias that expands anywhere on the command line, not just in command position. `alias -s ext=cmd` (suffix alias) runs `cmd file.ext` whenever a bare `file.ext` appears as a command. Both forms are Zsh-idiomatic interactive conveniences; in scripts they produce surprising substitutions that a reader cannot infer from local context — a bare word like `G` or `foo.log` stops meaning what it looks like. Use a function or a regular alias instead, and keep `alias -g` / `alias -s` in your `~/.zshrc` where the definition is discoverable.

Disable by adding `ZC1771` to `disabled_katas` in `.zshellcheckrc`.

---

<a id="zc1772"></a>
### ZC1772 — Error on `hdparm --security-erase` / `--trim-sector-ranges` — ATA-level data destruction

**Severity:** `error`  
**Auto-fix:** `no`

`hdparm --security-erase PASS $DISK` issues the ATA `SECURITY ERASE UNIT` command: the drive firmware wipes every block, ignoring filesystem or partition boundaries, and the operation cannot be interrupted or rolled back. `--security-erase-enhanced` is the same but also clears reallocated sectors, and `--trim-sector-ranges` discards the listed LBAs on any TRIM-capable device. `--security-set-pass`, `--security-disable`, `--security-unlock`, and `--security-freeze` alter the drive-level password state and, if misused in a script, lock the device out of future access. Keep these calls behind a guarded runbook with the exact disk pinned by `/dev/disk/by-id/…` and the password stored outside argv.

Disable by adding `ZC1772` to `disabled_katas` in `.zshellcheckrc`.

---

<a id="zc1773"></a>
### ZC1773 — Warn on `xargs` without `-r` / `--no-run-if-empty` — runs once on empty input

**Severity:** `warning`  
**Auto-fix:** `yes`

GNU `xargs` (the common default on Linux) invokes the child command once with no arguments when its stdin is empty. Paired with a destructive child (`xargs rm`, `xargs kill`, `xargs docker stop`) a pipeline that produces zero hits silently runs the command with no operand — usually an error at best and a footgun at worst. The flag `-r` (GNU) / `--no-run-if-empty` tells xargs to skip the call when no items arrive. Add `-r` to every `xargs` pipeline whose producer can return no results, or switch to `find ... -exec cmd {} +` which never runs the child on empty input. BSD xargs defaults to this behavior, but the portable and explicit choice is to pass `-r` and document the intent.

Disable by adding `ZC1773` to `disabled_katas` in `.zshellcheckrc`.

---

<a id="zc1774"></a>
### ZC1774 — Warn on `setopt GLOB_SUBST` — `$var` starts glob-expanding, user data becomes a pattern

**Severity:** `warning`  
**Auto-fix:** `no`

With `GLOB_SUBST` enabled, the result of any parameter expansion is rescanned for filename-generation metacharacters (`*`, `?`, `[`, `^`, `~`, brace ranges, qualifiers). Zsh's default — `NO_GLOB_SUBST` — keeps `$var` literal and matches the behavior most script authors expect after moving from Bash or POSIX sh. Turning `GLOB_SUBST` on globally means any unquoted `$var` that contains a metacharacter (environment, argv, file contents, user prompt) is expanded against the filesystem — an injection vector, and a subtle source of `no matches found` failures on empty variables. Keep `setopt GLOB_SUBST` inside a narrow subshell or function body, or use explicit `~` / `(e)` / `(P)` flags where you actually want the rescan.

Disable by adding `ZC1774` to `disabled_katas` in `.zshellcheckrc`.

---

<a id="zc1775"></a>
### ZC1775 — Warn on `timeout DURATION cmd` without `--kill-after` / `-k` — hang on SIGTERM-resistant child

**Severity:** `warning`  
**Auto-fix:** `no`

`timeout DURATION cmd` sends `SIGTERM` once the duration elapses and then waits for the child to exit. A child that blocks or ignores `SIGTERM` (long-running daemons, processes stuck in `D` state, a trapped / reset signal handler) never dies, so the entire pipeline hangs past the intended bound. Add `--kill-after=N` (`-k N`) so timeout escalates to `SIGKILL` after N seconds, guaranteeing exit. Typical choice: a few seconds shorter than your CI step budget, so the overall wait remains bounded.

Disable by adding `ZC1775` to `disabled_katas` in `.zshellcheckrc`.

---

<a id="zc1776"></a>
### ZC1776 — Error on `psql postgresql://user:secret@host/db` — password in argv via connection URI

**Severity:** `error`  
**Auto-fix:** `no`

Database and message-broker CLIs accept a single connection URI (`postgresql://`, `mysql://`, `mongodb://`, `redis://`, `amqp://`, `kafka://`, and friends). When the URI embeds a password — `scheme://user:secret@host/db` — the secret lands in argv, visible to every user via `ps`, `/proc/PID/cmdline`, process accounting, and audit trails, and it often survives in shell history. Keep the password out of argv: use the client's password-file / `.pgpass` / `PGPASSWORD` / `REDISCLI_AUTH` equivalent, or interpolate the URI from an environment variable so the secret is not on the command line that other users can see.

Disable by adding `ZC1776` to `disabled_katas` in `.zshellcheckrc`.

---

<a id="zc1777"></a>
### ZC1777 — Error on `tee/cp/mv/install/dd` writing `/etc/ld.so.preload` — classic rootkit persistence

**Severity:** `error`  
**Auto-fix:** `no`

`/etc/ld.so.preload` lists shared libraries that the dynamic linker forcibly loads into every dynamically-linked binary, root processes included. The file is almost never needed on a modern distribution — package managers do not touch it, and `LD_PRELOAD` handles the per-invocation case without persisting the change. A script that pipes content into `/etc/ld.so.preload` with `tee` / `cp` / `mv` / `install` / `dd` is a textbook rootkit persistence primitive (`libprocesshider`, `Azazel`, `Jynx`). Remove the line, audit `/etc/ld.so.preload` for unexpected entries (`sha256sum`, `diff` against a known-good backup), and if preloading is legitimately required, use a scoped `LD_PRELOAD=` on the specific invocation.

Disable by adding `ZC1777` to `disabled_katas` in `.zshellcheckrc`.

---

<a id="zc1778"></a>
### ZC1778 — Warn on `systemctl link /path/to/unit` — persistence from a mutable source path

**Severity:** `warning`  
**Auto-fix:** `no`

`systemctl link` symlinks the given unit file into `/etc/systemd/system/` so it can be `enable`d and `start`ed by name, but the unit definition lives at the original path forever. If that path is writable by any non-root user (`/tmp/*`, `/var/tmp/*`, `/home/*`, `/opt/` with wide perms, a build output directory), a later tamper of the source file silently changes what systemd runs the next time the unit starts. Copy the unit into `/etc/systemd/system/` with root-only permissions, or install a package that ships it under `/lib/systemd/system/`, rather than linking from a mutable location.

Disable by adding `ZC1778` to `disabled_katas` in `.zshellcheckrc`.

---

<a id="zc1779"></a>
### ZC1779 — Error on `az role assignment create --role Owner|Contributor|User Access Administrator`

**Severity:** `error`  
**Auto-fix:** `no`

`az role assignment create --role Owner` grants full control over the target scope (subscription, resource group, resource). `Contributor` grants everything except role assignment, and `User Access Administrator` grants the ability to assign any role — including Owner — elsewhere in the directory. Any of the three is effectively top-of-chain in the assigned scope. In provisioning automation this breaks least privilege, invites blast-radius escalations, and sidesteps any review that would flag the permission grant. Assign a narrower built-in role (Reader, specific-service Contributor) or a custom role whose permission list you can enumerate.

Disable by adding `ZC1779` to `disabled_katas` in `.zshellcheckrc`.

---

<a id="zc1780"></a>
### ZC1780 — Warn on `sysctl -w fs.protected_symlinks=0|protected_hardlinks=0|…` — TOCTOU guard disabled

**Severity:** `warning`  
**Auto-fix:** `no`

The `fs.protected_*` sysctls close a classic race: in a sticky directory (`/tmp`, `/var/tmp`, `/dev/shm`), a non-owner cannot follow a symlink, create a hardlink to a file they don't own, or open a FIFO / regular file they didn't create. Those four gates block the shape of attack where a privileged program predictably opens a `/tmp/NAME` that an attacker has already placed as a symlink to `/etc/shadow`. Setting any of them to `0` re-enables the race across the whole host. Leave the defaults (`1` / `2`) in place; if a specific application legitimately needs the old behavior, run it in a mount namespace where `/tmp` is not sticky-shared.

Disable by adding `ZC1780` to `disabled_katas` in `.zshellcheckrc`.

---

<a id="zc1781"></a>
### ZC1781 — Error on `git clone https://user:token@host/...` — PAT in argv and git config

**Severity:** `error`  
**Auto-fix:** `no`

A git remote URL in the form `https://user:token@host/path` puts the personal access token directly in argv — visible via `ps`, `/proc/PID/cmdline`, shell history, and process accounting. `git clone` additionally records the URL (including the credentials) in `.git/config` as the `origin` remote, so every later `git fetch` / `pull` re-exposes the same token to every user who can read that file. Use a credential helper (`git credential-store`, `git credential-osxkeychain`), `GIT_ASKPASS` with a secret pulled from an env var, HTTPS + an SSH deploy key, or set the token via the `Authorization: Bearer` header with `http.extraHeader` from an env-sourced value.

Disable by adding `ZC1781` to `disabled_katas` in `.zshellcheckrc`.

---

<a id="zc1782"></a>
### ZC1782 — Error on `flatpak remote-add --no-gpg-verify` — trust chain disabled for the repo

**Severity:** `error`  
**Auto-fix:** `no`

A Flatpak remote without GPG verification accepts any OSTree update that the server (or anyone on the path) cares to send. Signatures are what connect `flatpak install FOO` to the operator that actually built `FOO` — strip them and the install reduces to a plain HTTPS download with no identity attached. If you genuinely need a local / air-gapped repo, sign it yourself with `ostree gpg-sign` and add the key via `--gpg-import=KEYFILE`. Never leave `--no-gpg-verify` in provisioning scripts for production systems.

Disable by adding `ZC1782` to `disabled_katas` in `.zshellcheckrc`.

---

<a id="zc1783"></a>
### ZC1783 — Error on `podman system reset` / `nerdctl system prune -af --volumes` — wipes every container artifact

**Severity:** `error`  
**Auto-fix:** `no`

`podman system reset` removes every podman container, image, volume, network, pod, secret, and storage driver scratch area — a full factory reset of the local engine. `nerdctl system prune -af --volumes` achieves the same for containerd. On a developer workstation this wipes cached images for unrelated projects; on a CI runner or build host it invalidates every warm artifact the job relies on; on a prod host it drops the volumes the workload stores data in. Use narrower commands (`podman rmi`, `podman volume rm`, scoped `podman prune`) that only touch the resource you intend to remove, and never pair the reset with `--force`.

Disable by adding `ZC1783` to `disabled_katas` in `.zshellcheckrc`.

---

<a id="zc1784"></a>
### ZC1784 — Warn on `git config core.hooksPath /tmp/...` — hook execution from a mutable path

**Severity:** `warning`  
**Auto-fix:** `no`

`core.hooksPath` tells git which directory to run repository hooks from. Any file named `pre-commit`, `post-checkout`, `post-merge`, etc. under that directory becomes executable code invoked by routine git operations. Pointing `core.hooksPath` at `/tmp`, `/var/tmp`, `/dev/shm`, `/home/<other>`, `/opt`, `/srv`, or `/mnt` hands the git CLI an execution primitive from a path that a non-root (or another) user can write at will — a classic supply-chain entry point on shared hosts and CI runners. Keep hooks inside the repo's `.git/hooks/` (or a repo-owned `.githooks/` directory) and configure `core.hooksPath` only to paths that share the repo's owner and permissions.

Disable by adding `ZC1784` to `disabled_katas` in `.zshellcheckrc`.

---

<a id="zc1785"></a>
### ZC1785 — Error on `ufw default allow` — flips host firewall from deny-by-default to allow-by-default

**Severity:** `error`  
**Auto-fix:** `no`

`ufw default allow incoming` (or `allow outgoing`, `allow routed`) changes the chain's baseline verdict — instead of only what you explicitly opened, every port that does not have a matching `deny` rule is accepted. On an internet-facing host this is effectively "turn the firewall off", and the effect survives reboots because the default is persisted to `/etc/default/ufw`. Restore with `ufw default deny incoming` and add narrow `ufw allow <port>` rules for the services that actually need to be reachable.

Disable by adding `ZC1785` to `disabled_katas` in `.zshellcheckrc`.

---

<a id="zc1786"></a>
### ZC1786 — Error on `mount.cifs ... -o password=SECRET` — SMB password in argv

**Severity:** `error`  
**Auto-fix:** `no`

Passing `password=` (or `pass=`) inside `mount.cifs` / `mount -t cifs` options puts the SMB password in argv. Any local user who can read `ps`, `/proc/PID/cmdline`, or process-accounting records gets the cleartext, and the line also ends up in shell history and — if captured — in CI logs. Use a `credentials=/etc/cifs-creds` file (`0600`, `username=` and `password=` lines), the `$USER`/`$PASSWD` env vars `mount.cifs` reads when those options are missing, or `pam_mount` for login-time mounts.

Disable by adding `ZC1786` to `disabled_katas` in `.zshellcheckrc`.

---

<a id="zc1787"></a>
### ZC1787 — Warn on `setopt AUTO_CD` — bare word that names a directory silently changes `$PWD`

**Severity:** `warning`  
**Auto-fix:** `no`

With `AUTO_CD` on, any bare word that happens to name an existing directory is executed as `cd <word>` — no command name, no error. This is a pleasant interactive shortcut and an absolute footgun in scripts: a typo in a command name (`dockr` → a directory called `dockr` that was left lying around) or a user-controlled variable that expands to a path silently reshapes `$PWD` for every later relative path. Keep `AUTO_CD` inside `~/.zshrc` where it belongs, not in a `.zsh` script, and never turn it on inside a function that an external caller depends on.

Disable by adding `ZC1787` to `disabled_katas` in `.zshellcheckrc`.

---

<a id="zc1788"></a>
### ZC1788 — Warn on `ssh -F /tmp/config` — config from a mutable path can pin `ProxyCommand` to arbitrary code

**Severity:** `warning`  
**Auto-fix:** `no`

`ssh -F PATH` (and `scp -F PATH`, `sftp -F PATH`) loads a user-supplied config file. Anything in `/etc/ssh/ssh_config` can be overridden — notably `ProxyCommand`, `LocalCommand`, `PermitLocalCommand`, and `Include` — which means a mutable source path is an execution primitive: another local user flips `ProxyCommand` to `/tmp/pwn`, and the next `ssh` run launches it with the caller's credentials and forwarded agent. Keep the config in `~/.ssh/config` (or a repo-owned path with the same owner and `0600` perms) and never pass `-F` to `/tmp`, `/var/tmp`, `/dev/shm`, another user's `/home`, `/opt`, `/srv`, or `/mnt`.

Disable by adding `ZC1788` to `disabled_katas` in `.zshellcheckrc`.

---

<a id="zc1789"></a>
### ZC1789 — Warn on `setopt CORRECT` / `CORRECT_ALL` — Zsh spellcheck silently rewrites script tokens

**Severity:** `warning`  
**Auto-fix:** `no`

`setopt CORRECT` prompts to rewrite command names that look mistyped; `CORRECT_ALL` extends the check to every argument on the line. In an interactive shell this is a friendly nudge. In a script it becomes a footgun: a filename that is *close enough* to an existing file gets silently replaced with that other file, and the "nlh?" prompt reads from stdin — which may be the input the script was supposed to process. Keep `CORRECT` / `CORRECT_ALL` in `~/.zshrc` only and never toggle them inside a function a script calls.

Disable by adding `ZC1789` to `disabled_katas` in `.zshellcheckrc`.

---

<a id="zc1790"></a>
### ZC1790 — Warn on `unsetopt PIPE_FAIL` — pipeline exit status reverts to last-command-only

**Severity:** `warning`  
**Auto-fix:** `no`

With `PIPE_FAIL` off (the shell default), `cmd1 | cmd2 | cmd3` exits with `cmd3`'s status; failures in `cmd1` and `cmd2` are silently dropped. `unsetopt PIPE_FAIL` (or the equivalent `setopt NOPIPEFAIL`) mid-script turns a previously-enabled error check back off — typically because a known-flaky pipe stage was tripping `set -e`, and the author reached for the global off-switch. Undo the change in a subshell (`( unsetopt pipefail; …; )`) or a function with `emulate -L zsh; unsetopt pipefail` so the rest of the script keeps strict-pipe error propagation.

Disable by adding `ZC1790` to `disabled_katas` in `.zshellcheckrc`.

---

<a id="zc1791"></a>
### ZC1791 — Error on `curl --unix-socket /var/run/docker.sock` — direct container-daemon API access

**Severity:** `error`  
**Auto-fix:** `no`

A curl request to `docker.sock` / `containerd.sock` / `crio.sock` speaks the container-daemon HTTP API with no authentication beyond the socket's filesystem permissions. Anyone who can invoke curl as that uid can `POST /containers/create` with `HostConfig.Privileged=true` and a bind mount of `/` and land a root shell on the host — the primitive every "docker socket escape" write-up leans on. Use the real CLI (`docker`, `podman`, `nerdctl`) which enforces its own policy, or access the daemon over a TLS-protected TCP endpoint with mutual auth.

Disable by adding `ZC1791` to `disabled_katas` in `.zshellcheckrc`.

---

<a id="zc1792"></a>
### ZC1792 — Warn on `btrfs subvolume delete` / `btrfs device remove` — unrecoverable btrfs data loss

**Severity:** `warning`  
**Auto-fix:** `no`

`btrfs subvolume delete PATH` unlinks the subvolume and drops all of its extents once cleanup completes — on Snapper / Timeshift systems the argument is often a snapshot that is the only remaining copy of pre-incident state. `btrfs device remove DEV POOL` moves the stored chunks off DEV before detaching it; wrong device, mid-rebalance failure, or insufficient free space across the remaining members puts the filesystem into degraded mode with no automatic rollback. Keep a fresh `btrfs subvolume list`/`btrfs device usage` snapshot and confirm the target explicitly before running either command in automation.

Disable by adding `ZC1792` to `disabled_katas` in `.zshellcheckrc`.

---

<a id="zc1793"></a>
### ZC1793 — Warn on `kubectl certificate approve CSR` — signs the identity baked into the CSR

**Severity:** `warning`  
**Auto-fix:** `no`

`kubectl certificate approve NAME` tells the cluster signer to sign the pending CSR unchanged. The signer respects the Subject (CN, O) and the SubjectAltName extensions the caller put in the CSR — approve one that requests `system:masters` and you have handed the requester full admin on the cluster. In automation, review the CSR body first (`kubectl get csr NAME -o jsonpath='{.spec.request}' | base64 -d | openssl req -text`) and reject (`kubectl certificate deny`) any request that names a privileged group, kube-system service account, or hostname outside the intended scope.

Disable by adding `ZC1793` to `disabled_katas` in `.zshellcheckrc`.

---

<a id="zc1794"></a>
### ZC1794 — Error on `cosign verify --insecure-ignore-tlog` / `--allow-insecure-registry` — signature chain disabled

**Severity:** `error`  
**Auto-fix:** `no`

`cosign verify` with `--insecure-ignore-tlog` skips Rekor transparency-log verification, `--insecure-ignore-sct` skips Fulcio SCT verification, and `--insecure-skip-verify` turns off TLS certificate validation for the registry / Rekor / Fulcio endpoints. `cosign sign --allow-insecure-registry` and `--allow-http-registry` push signatures over plain HTTP. Each flag removes a distinct rung of the signature chain that `cosign` was built to enforce — a malicious registry or on-path attacker now passes verification without detection. Drop the flag, fix the underlying trust anchor (CA bundle, Rekor URL, Fulcio OIDC), and keep signature verification strict.

Disable by adding `ZC1794` to `disabled_katas` in `.zshellcheckrc`.

---

<a id="zc1795"></a>
### ZC1795 — Error on `git remote add NAME https://user:token@host/repo` — credentials persisted in `.git/config`

**Severity:** `error`  
**Auto-fix:** `no`

`git remote add NAME URL` and `git remote set-url NAME URL` write the URL into `.git/config` verbatim. When the URL embeds a `user:token@host` credential segment, every reader of the repo — other local users, a compromised backup, a CI cache, or anyone who runs `git config --list` — picks up the secret. It also shows up in argv at the moment of creation (visible via `ps` / `/proc/PID/cmdline`). Use a credential helper (`git credential-store`, `credential-osxkeychain`), `GIT_ASKPASS` sourced from an env var, or HTTPS + a deploy SSH key.

Disable by adding `ZC1795` to `disabled_katas` in `.zshellcheckrc`.

---

<a id="zc1796"></a>
### ZC1796 — Warn on `pg_restore --clean` / `-c` — drops existing DB objects before recreating

**Severity:** `warning`  
**Auto-fix:** `no`

`pg_restore -c` (also `--clean`) issues `DROP` for every table, index, function, and sequence in the target database before recreating them from the archive. If the backup is stale, incomplete, or points at the wrong database, the destination loses any object that isn't in the dump — including data added after the backup ran. Restore into a fresh empty database (`createdb new && pg_restore -d new`) or snapshot the target (`pg_dump -Fc > pre.dump`) before running `--clean`, and never pair it with `--if-exists` on a live production DB without a tested rollback path.

Disable by adding `ZC1796` to `disabled_katas` in `.zshellcheckrc`.

---

<a id="zc1797"></a>
### ZC1797 — Warn on `ip link set <iface> down` / `ifdown <iface>` — locks out remote admin on that path

**Severity:** `warning`  
**Auto-fix:** `no`

Taking a network interface down from an SSH session that rides on the same interface cuts the script off mid-run: the TCP connection freezes, any later step silently fails, and recovery requires console / out-of-band access. Common bugs are typos (`eth1` instead of `eth0`), scripts that target the only uplink on a cloud VM, or running the command without first confirming that the interface is not the one carrying the admin session. Wrap the `down` in a `systemd-run --on-active=30s --unit=recover ip link set <iface> up` rollback, or stage both `down` and `up` through `nmcli connection up/down` with a pinned fallback profile.

Disable by adding `ZC1797` to `disabled_katas` in `.zshellcheckrc`.

---

<a id="zc1798"></a>
### ZC1798 — Warn on `ufw reset` — wipes every configured firewall rule

**Severity:** `warning`  
**Auto-fix:** `no`

`ufw reset` returns the firewall to the distro default: every user-defined rule is removed, default incoming policy reverts (usually to `deny`, but the net effect is the loss of every allow-list entry the host relied on). Paired with `--force`, no prompt is issued. In a provisioning script the operation is sometimes desired to start from a clean slate, but running it mid-session or on a host that currently serves traffic drops connections without warning. Snapshot the rules first (`ufw status numbered > /tmp/ufw.bak`), and prefer removing specific rules with `ufw delete <num>` over a full reset.

Disable by adding `ZC1798` to `disabled_katas` in `.zshellcheckrc`.

---

<a id="zc1799"></a>
### ZC1799 — Warn on `rclone sync SRC DST` without `--dry-run` — one-way mirror can wipe DST

**Severity:** `warning`  
**Auto-fix:** `no`

`rclone sync` makes DST look exactly like SRC: anything in DST that isn't in SRC is deleted, including object versions on providers that support them. If SRC is accidentally empty (typo in path, unmounted drive, wrong credentials pointing at an empty bucket), the command silently wipes every object under DST without a confirmation prompt. Always preview the diff with `rclone sync --dry-run SRC DST` first; when you commit to the sync, keep `--backup-dir`, `--max-delete`, or `--min-age` guards so a bad SRC cannot cascade into unbounded deletion.

Disable by adding `ZC1799` to `disabled_katas` in `.zshellcheckrc`.

---

<a id="zc1800"></a>
### ZC1800 — Warn on `pg_ctl stop -m immediate` — abrupt shutdown skips checkpoint, forces WAL recovery

**Severity:** `warning`  
**Auto-fix:** `no`

`pg_ctl stop -m immediate` sends `SIGQUIT` to the postmaster. Server processes drop connections, no checkpoint is taken, and buffered changes are left in memory. Recovery on the next start has to replay every record since the last checkpoint; if WAL is corrupt, lost, or on different storage, committed transactions can be lost. Use `-m smart` (default) or `-m fast` so the server issues a shutdown checkpoint and closes cleanly; reserve `immediate` for the "the node is on fire" case and pair it with a tested PITR procedure.

Disable by adding `ZC1800` to `disabled_katas` in `.zshellcheckrc`.

---

<a id="zc1801"></a>
### ZC1801 — Warn on `fwupdmgr update` / `install` — mid-flash interruption can brick firmware

**Severity:** `warning`  
**Auto-fix:** `no`

`fwupdmgr update`, `fwupdmgr upgrade`, and `fwupdmgr install FIRMWARE` push new firmware into BIOS / UEFI, SSD, Thunderbolt controller, NIC, or dock microcontroller. Most of those devices have no A/B rollback — an interrupted flash (power cut, unexpected reboot, PSU toggle) leaves the chip in an unbootable state that needs vendor-recovery hardware. Run from a battery-backed session, mask reboot triggers with `systemd-inhibit`, pin the power supply, and verify the update history with `fwupdmgr get-history` once the device returns.

Disable by adding `ZC1801` to `disabled_katas` in `.zshellcheckrc`.

---

<a id="zc1802"></a>
### ZC1802 — Warn on `dnf history undo N` / `rollback N` — reverses transactions without compat check

**Severity:** `warning`  
**Auto-fix:** `no`

`dnf history undo N` reverts the exact package set of transaction N — every install turns into a remove, every remove into an install, every update into a downgrade. `dnf history rollback N` does the same for every transaction after N. Neither checks that the older versions still resolve cleanly against the current package graph: dependencies that moved forward for other reasons end up downgraded alongside, security patches get reverted, and services whose configuration was migrated fail to start on the older binary. Review the plan with `dnf history info N`, pin the rollback scope with `--exclude=` / `--assumeyes` only after review, or restore from a filesystem snapshot taken before the original transaction.

Disable by adding `ZC1802` to `disabled_katas` in `.zshellcheckrc`.

---

<a id="zc1803"></a>
### ZC1803 — Error on `mysql --skip-ssl` / `psql sslmode=disable` — plaintext credentials on the wire

**Severity:** `error`  
**Auto-fix:** `no`

Disabling TLS on a MySQL or PostgreSQL client pushes the login handshake (including the password or auth challenge) and every subsequent query and result over plaintext TCP. Anyone in the network path — the cloud VPC, the office LAN, a compromised router — can sniff or modify the stream. The flags vary (`--skip-ssl`, `--ssl=0`, `--ssl-mode=DISABLED` for MySQL / MariaDB; `sslmode=disable` in the connection URI or `PGSSLMODE=disable` env var for PostgreSQL) but the effect is the same. Prefer `--ssl-mode=VERIFY_IDENTITY` (MySQL 8+) and `sslmode=verify-full` (psql) with a pinned CA bundle.

Disable by adding `ZC1803` to `disabled_katas` in `.zshellcheckrc`.

---

<a id="zc1804"></a>
### ZC1804 — Warn on `aws ec2 terminate-instances` / `delete-volume` / `delete-snapshot` — destructive cloud state change

**Severity:** `warning`  
**Auto-fix:** `no`

AWS EC2 destructive actions (`terminate-instances`, `delete-volume`, `delete-snapshot`, `delete-vpc`, and friends) drop cloud state without any automatic backup: instance-store volumes vanish on terminate, EBS volumes and snapshots cannot be restored from the AWS side once deleted, and a wrong VPC / ENI / security-group ID can take down workloads in the same account. Review the target list with `aws ec2 describe-…`, pair destructive commands with `--dry-run`, and keep the IDs pinned in a file that `aws ... --cli-input-json` can consume rather than passing them inline.

Disable by adding `ZC1804` to `disabled_katas` in `.zshellcheckrc`.

---

<a id="zc1805"></a>
### ZC1805 — Warn on `aws cloudformation delete-stack` / `dynamodb delete-table` / `logs delete-log-group` / `kms schedule-key-deletion` — destructive AWS state change

**Severity:** `warning`  
**Auto-fix:** `no`

Each of these AWS actions drops state that AWS cannot restore: `cloudformation delete-stack` tears down every resource the stack manages in dependency order and has no rollback, `dynamodb delete-table` removes a table and its items, `logs delete-log-group` erases the CloudWatch audit trail, and `kms schedule-key-deletion` makes every ciphertext encrypted with the CMK unreadable after the grace window. Add `--dry-run` where supported, stage the call behind a typed confirmation, pin IDs through `--cli-input-json`, and export backups (`dynamodb export-table-to-point-in-time`, `logs create-export-task`) before pulling the trigger.

Disable by adding `ZC1805` to `disabled_katas` in `.zshellcheckrc`.

---

<a id="zc1806"></a>
### ZC1806 — Warn on `zmv 'PAT' 'REP'` without `-n` / `-i` — silent bulk rename

**Severity:** `warning`  
**Auto-fix:** `no`

`zmv` (autoloaded from Zsh's functions) rewrites every filename that matches the pattern in one shot. A small typo in the source pattern or replacement — `*.jpg` vs `*.JPG`, a misplaced `(..)`, forgetting `**` recursion — can collide names and silently overwrite files, since `zmv` aborts the batch only on its own conflict check, not on semantic errors. Use `zmv -n 'PAT' 'REP'` first to see the rename list, or `zmv -i` to prompt per file. Only drop the guard once the preview matches what you expect.

Disable by adding `ZC1806` to `disabled_katas` in `.zshellcheckrc`.

---

<a id="zc1807"></a>
### ZC1807 — Warn on `gh api -X DELETE` — raw GitHub DELETE bypasses `gh` command confirmations

**Severity:** `warning`  
**Auto-fix:** `no`

`gh api -X DELETE /repos/OWNER/REPO` (and `--method=DELETE` variants) sends a raw GitHub API request with the caller's token. There is no confirmation prompt, no `--yes` guard, and no friendly dry-run — a script that builds the path from a variable can wipe repos, releases, deploy keys, workflow runs, issue comments, or whole organisations in one call. Use the high-level `gh` subcommand for the target (`gh repo delete`, `gh release delete`, `gh workflow disable`) which still at least requires `--yes`, or wrap the raw call with a preflight `gh api -X GET /path` and an explicit confirmation in the script.

Disable by adding `ZC1807` to `disabled_katas` in `.zshellcheckrc`.

---

<a id="zc1808"></a>
### ZC1808 — Warn on `kubectl replace --force` — deletes + recreates resource, drops running pods

**Severity:** `warning`  
**Auto-fix:** `no`

`kubectl replace --force -f FILE` is `delete` followed by `create`: the existing resource (and every dependent pod / replicaset / endpoint) is removed before the new manifest is applied. In-flight requests drop, PodDisruptionBudget is ignored, and controllers that watch the object see it disappear and reappear. Prefer `kubectl apply -f FILE` — same manifest, server-side merge that preserves running pods — and reach for `replace --force` only when the resource schema has changed in a way `apply` cannot patch, with traffic drained beforehand.

Disable by adding `ZC1808` to `disabled_katas` in `.zshellcheckrc`.

---

<a id="zc1809"></a>
### ZC1809 — Error on `gsutil rm -r gs://…` / `gsutil rb -f gs://…` — bulk GCS deletion

**Severity:** `error`  
**Auto-fix:** `no`

`gsutil rm -r gs://bucket/prefix` and `gsutil rm -rf gs://bucket` delete every object under the prefix — with `-m` (parallel) they do it faster than any undo window. `gsutil rb -f gs://bucket` removes the bucket after force-deleting the contents. Neither soft-deletes; Object Versioning can help only if it is turned on in advance, and `gsutil rb` leaves no retention grace. Preview with `gsutil ls`, enable Object Versioning or retention locks before the fact, and prefer narrower `gsutil rm gs://bucket/specific-object` calls.

Disable by adding `ZC1809` to `disabled_katas` in `.zshellcheckrc`.

---

<a id="zc1810"></a>
### ZC1810 — Warn on `wget -r` / `--mirror` without `--level=N` — unbounded recursive download

**Severity:** `warning`  
**Auto-fix:** `no`

`wget -r` and `wget --mirror` (short `-m`) follow links to arbitrary depth. Without `--level=N` or `-l N` the crawl keeps going until `wget` hits the remote server's limits, fills the local disk, or climbs into a parent directory the author did not intend to mirror (add `--no-parent` to block that too). Pin a depth (`--level=3`), restrict siblings (`--no-parent`, `--accept=` / `--reject=`), and cap the byte budget (`--quota=1G`) before running a recursive wget in automation.

Disable by adding `ZC1810` to `disabled_katas` in `.zshellcheckrc`.

---

<a id="zc1811"></a>
### ZC1811 — Error on `chown/chmod/chgrp --no-preserve-root` — disables GNU safeguard against recursive `/`

**Severity:** `error`  
**Auto-fix:** `no`

GNU `chown`, `chmod`, and `chgrp` refuse to recurse into `/` by default (`--preserve-root` in coreutils). `--no-preserve-root` opts in to walking the entire filesystem, so a stray `$PATH` expansion or wrong variable combined with `-R` rewrites ownership or mode on every file on the host. The flag has no legitimate script use — if a specific top-level target genuinely needs recursion, list that path explicitly and keep the safeguard in place.

Disable by adding `ZC1811` to `disabled_katas` in `.zshellcheckrc`.

---

<a id="zc1812"></a>
### ZC1812 — Error on `aws ssm put-parameter --type SecureString --value SECRET` — plaintext in argv

**Severity:** `error`  
**Auto-fix:** `no`

`aws ssm put-parameter` stores the value as-is under the given parameter name; the whole point of `--type SecureString` is that the value is sensitive. Passing the plaintext with `--value SECRET` (or `--value=SECRET`) puts the secret in argv where `ps`, `/proc/PID/cmdline`, shell history, and AWS CLI debug logs (`--debug`) can read it. Pipe the value in from stdin with `--cli-input-json file://param.json` (mode 0600) or use `aws secretsmanager create-secret --secret-string file://secret` which supports `file://` in every code path.

Disable by adding `ZC1812` to `disabled_katas` in `.zshellcheckrc`.

---

<a id="zc1813"></a>
### ZC1813 — Warn on `cryptsetup luksFormat` / `reencrypt` — destructive LUKS header write

**Severity:** `warning`  
**Auto-fix:** `no`

`cryptsetup luksFormat DEV` writes a new LUKS2 header at the start of DEV and marks the remaining space as fresh ciphertext — any pre-existing filesystem or LUKS metadata is gone. `cryptsetup reencrypt DEV` rewrites the entire device in place, and an interruption mid-write leaves the volume partially re-encrypted and dependent on the `--resume-only` recovery path. Pair `luksFormat` with `--batch-mode` only after verifying DEV via `lsblk -o NAME,MODEL,SERIAL`, always back up the header (`cryptsetup luksHeaderBackup`) before touching it, and run `reencrypt` on an unmounted volume with UPS-backed power.

Disable by adding `ZC1813` to `disabled_katas` in `.zshellcheckrc`.

---

<a id="zc1814"></a>
### ZC1814 — Error on `dpkg --force-all` — enables every single `--force-*` option at once

**Severity:** `error`  
**Auto-fix:** `no`

`dpkg --force-all` is shorthand for ~18 distinct `--force-<option>` flags: overwrite existing files, install unsigned packages, downgrade, install depends-broken, remove essential, and more. The dpkg manual explicitly calls this "almost always a bad idea". In provisioning scripts it hides the specific constraint the author was trying to bypass, and when a later install re-triggers the same state the underlying dependency conflict just re-surfaces on the next unattended upgrade. Drop `--force-all` and spell out only the `--force-<option>` you genuinely need, or fix the upstream conflict.

Disable by adding `ZC1814` to `disabled_katas` in `.zshellcheckrc`.

---

<a id="zc1815"></a>
### ZC1815 — Warn on `systemctl restart NetworkManager` / `systemd-networkd` — drops the SSH session

**Severity:** `warning`  
**Auto-fix:** `no`

Restarting the network manager from an SSH session tears down every active connection the daemon supervises, including the one the script is running over. The script freezes, the client sees a broken pipe, and recovery usually requires console access. Route the change through `nmcli connection reload` + `nmcli connection up <name>` (NetworkManager), `networkctl reload` (systemd-networkd), or schedule the restart behind `systemd-run --on-active=30s` with a rollback timer that re-enables the previous config if SSH does not reconnect.

Disable by adding `ZC1815` to `disabled_katas` in `.zshellcheckrc`.

---

<a id="zc1816"></a>
### ZC1816 — Warn on `docker/podman commit` — produces un-reproducible image, bakes in runtime state

**Severity:** `warning`  
**Auto-fix:** `no`

`docker commit CONTAINER IMAGE` (and the podman / nerdctl equivalents) snapshots a running container's filesystem into a new image. There is no Dockerfile, so the build is not reproducible; the snapshot inherits whatever `/tmp` scratch, shell history, environment variables, and — frequently — credentials the container held at that moment; and the resulting image's layer metadata records only the container id, not what was actually installed. Build from a `Dockerfile` (or `docker buildx build`) so the image can be regenerated from source, and use `docker commit` only for one-off rescue work on a local image you are about to discard.

Disable by adding `ZC1816` to `disabled_katas` in `.zshellcheckrc`.

---

<a id="zc1817"></a>
### ZC1817 — Warn on `git push --delete` / `git push -d` / `git push origin :branch` — remote branch removal

**Severity:** `warning`  
**Auto-fix:** `no`

Deleting a branch on the remote is an irreversible server-side change the local reflog cannot rescue. `git push --delete REMOTE BRANCH`, the short `-d`, and the legacy `git push REMOTE :BRANCH` colon form all produce the same result: the ref vanishes from the server, open pull requests are orphaned, CI runners that pinned to the branch lose the target, and recovery needs the last commit SHA which may only live in somebody else's local clone. Confirm the remote name, check `git branch -r` / `gh pr list --head BRANCH` first, and prefer letting the hosting platform delete the branch after a PR merge (with the auto-delete setting) rather than scripting the push.

Disable by adding `ZC1817` to `disabled_katas` in `.zshellcheckrc`.

---

<a id="zc1818"></a>
### ZC1818 — Warn on `rsync --delete` without `--dry-run` — empty or wrong SRC wipes DST

**Severity:** `warning`  
**Auto-fix:** `no`

`rsync --delete` (plus `--delete-before/-during/-after/-excluded`) removes anything in DST that is not in SRC. If SRC is accidentally empty (typo in path, unmounted mount point, wrong credentials pointing at an empty remote), the destination loses every file that was there. The command has no undo. Always preview the diff with `rsync -av --delete --dry-run SRC DST` first, and cap the blast radius with `--max-delete=N` so the sync aborts if the plan removes more files than expected.

Disable by adding `ZC1818` to `disabled_katas` in `.zshellcheckrc`.

---

<a id="zc1819"></a>
### ZC1819 — Warn on `xattr -d com.apple.quarantine` / `xattr -cr` — removes macOS Gatekeeper quarantine

**Severity:** `warning`  
**Auto-fix:** `no`

macOS sets the `com.apple.quarantine` extended attribute on every file downloaded from the internet — Gatekeeper uses it to trigger the first-run notarization / signature check. `xattr -d com.apple.quarantine FILE` strips the attribute and lets the binary run with no prompt, and `xattr -cr DIR` does the same recursively for every file in the tree. In a script that processes downloaded artifacts this turns "we vetted the binary" into "we trust whatever landed in the download folder". Verify the signature (`codesign --verify`) and notarization (`spctl --assess --type execute`) first, or use `xip`/`installer` packages so Gatekeeper stays in the loop.

Disable by adding `ZC1819` to `disabled_katas` in `.zshellcheckrc`.

---

<a id="zc1820"></a>
### ZC1820 — Warn on `netplan apply` — applies network config immediately with no rollback timer

**Severity:** `warning`  
**Auto-fix:** `no`

`netplan apply` regenerates the rendered backend config (systemd-networkd or NetworkManager) and brings it live right away. A mistake in the YAML — wrong interface name, missing `dhcp4`, bad addresses, conflicting routes — drops the admin SSH session, and recovery needs console access. Run `netplan try` first: it applies the new config, waits for confirmation, and rolls back automatically if no keypress arrives within the timeout. Only fall through to `netplan apply` after the try window has elapsed successfully.

Disable by adding `ZC1820` to `disabled_katas` in `.zshellcheckrc`.

---

<a id="zc1821"></a>
### ZC1821 — Error on `diskutil eraseDisk` / `secureErase` / `partitionDisk` — macOS storage reformat

**Severity:** `error`  
**Auto-fix:** `no`

The `diskutil` subcommands `eraseDisk`, `eraseVolume`, `secureErase`, `zeroDisk`, `randomDisk`, `reformat`, `erasePartitions`, and `partitionDisk` all rewrite disk or volume state with no Time Machine snapshot or APFS preservation. A wrong `/dev/diskN` (especially after a reboot that reordered the BSD names) erases the wrong drive, and the only recovery is an offline backup. Always pair the call with a typed confirmation, resolve the target by `diskutil info -plist` / mount-point rather than by index, and run `diskutil list` right before the destructive call.

Disable by adding `ZC1821` to `disabled_katas` in `.zshellcheckrc`.

---

<a id="zc1822"></a>
### ZC1822 — Error on `csrutil disable` / `spctl --master-disable` — disables macOS system integrity / Gatekeeper

**Severity:** `error`  
**Auto-fix:** `no`

`csrutil disable` turns off System Integrity Protection: the kernel stops blocking writes under `/System`, `/bin`, `/sbin`, runtime attachment to protected processes becomes possible, and unsigned kexts can load. `spctl --master-disable` (and `--global-disable`, `kext-consent disable`) removes Gatekeeper / kext-consent enforcement, so any downloaded binary or kernel extension runs without the user being prompted. Neither has a legitimate provisioning use; both belong to ad-hoc developer workflows and are high-value persistence steps for malware. Re-enable with `csrutil enable` in recovery mode and `spctl --master-enable`.

Disable by adding `ZC1822` to `disabled_katas` in `.zshellcheckrc`.

---

<a id="zc1823"></a>
### ZC1823 — Warn on `keytool -import -noprompt` — Java trust store imports without fingerprint check

**Severity:** `warning`  
**Auto-fix:** `no`

`keytool -import -noprompt -trustcacerts -alias X -file CERT -keystore KS` adds CERT to the Java trust store without showing its SHA-256 fingerprint or asking the operator to confirm. If CERT came from an HTTP download, an attacker wrote it in a shared temp dir, or a provisioning step fetched the wrong file, the JVM will happily pin the attacker's CA as trusted and verify everything signed against it. Drop `-noprompt`, or pre-verify with `keytool -printcert -file CERT` and keep the alias+fingerprint pair in a versioned inventory before adding to any trust store.

Disable by adding `ZC1823` to `disabled_katas` in `.zshellcheckrc`.

---

<a id="zc1824"></a>
### ZC1824 — Warn on `kubectl drain --disable-eviction` — bypasses PodDisruptionBudget via raw DELETE

**Severity:** `warning`  
**Auto-fix:** `no`

`kubectl drain --disable-eviction` tells the client to delete pods directly via the API instead of issuing Eviction requests. The Eviction pathway is what honours PodDisruptionBudget — `--disable-eviction` drops pods regardless of the minAvailable / maxUnavailable contract the workload owner defined. On a multi-replica service this turns a rolling drain into a hard outage. Fix the blocking PDB (raise minAvailable, wait for replicas to reschedule, or negotiate with the owner) instead of flipping the flag off.

Disable by adding `ZC1824` to `disabled_katas` in `.zshellcheckrc`.

---

<a id="zc1825"></a>
### ZC1825 — Warn on `scp -O` — forces legacy SCP wire protocol exposed to filename-injection CVEs

**Severity:** `warning`  
**Auto-fix:** `no`

OpenSSH 9.0 switched `scp` to use the SFTP protocol by default — SFTP performs structured file transfer instead of piping a remote shell, and closes the filename-injection class that the old SCP wire protocol was vulnerable to (CVE-2020-15778 and friends). `scp -O` forces the legacy SCP protocol, putting the connection back on the old code path where a server (or a man-in-the-middle in the remote host's shell) can inject shell metacharacters into filenames. If a remote endpoint genuinely needs SCP, use `sftp` instead or upgrade the remote server. Drop `-O` unless you have a named compatibility bug that requires it.

Disable by adding `ZC1825` to `disabled_katas` in `.zshellcheckrc`.

---

<a id="zc1826"></a>
### ZC1826 — Warn on `install -m u+s` / `g+s` — symbolic setuid/setgid bit applied at install time

**Severity:** `warning`  
**Auto-fix:** `no`

`install -m u+s SRC DEST` (or `g+s` / `ug+s` / `u=rwxs` etc.) applies the setuid / setgid bit atomically at copy time — no intermediate `chmod` step where a tripwire would fire, no time window where the file exists without the special bit. Symbolic forms are easy to miss in review because they don't carry the tell-tale leading `4`/`2`/`6` digit that numeric-mode detection (see ZC1892) keys off. If DEST is on `$PATH`, every local user can invoke the elevated binary. Install setuid / setgid binaries only from trusted builds you have reviewed, and prefer narrow capabilities (`setcap cap_net_bind_service+ep`) over broad setuid.

Disable by adding `ZC1826` to `disabled_katas` in `.zshellcheckrc`.

---

<a id="zc1827"></a>
### ZC1827 — Error on `npm unpublish` — breaks every downstream that pinned the version

**Severity:** `error`  
**Auto-fix:** `no`

`npm unpublish PKG@VERSION` removes a published version from the registry. Every downstream that pinned to that version — directly or through a transitive lockfile entry — fails to install on the next `npm ci` / CI run. This is the exact mechanism behind the 2016 `left-pad` outage; npm has since limited unpublish to within 72 hours and added the `--force` gate, but within the window the blast radius is still the whole ecosystem that pulled the package. Use `npm deprecate PKG@VERSION 'reason'` instead — the version stays resolvable, but installs print a warning and users can pin forward on their own schedule.

Disable by adding `ZC1827` to `disabled_katas` in `.zshellcheckrc`.

---

<a id="zc1828"></a>
### ZC1828 — Warn on `gcore PID` / `strace -p PID` — live ptrace attach dumps target memory

**Severity:** `warning`  
**Auto-fix:** `no`

`gcore PID` writes a core dump of the running process to disk; `strace -p PID` streams every syscall the process makes. Both attach via ptrace and expose the target's memory, stack, environment variables, and argument buffers — credentials, TLS session keys, and `$AWS_SECRET_ACCESS_KEY`-style env vars are all readable. A root-run script that attaches to another user's process extracts whatever that user has. Keep production scripts off ptrace; reach for `coredumpctl` with a captured core or vendor-specific `perf` counters when you only need syscall statistics.

Disable by adding `ZC1828` to `disabled_katas` in `.zshellcheckrc`.

---

<a id="zc1829"></a>
### ZC1829 — Warn on `tailscale down` / `wg-quick down` / `nmcli con down` — drops the VPN that may carry the SSH session

**Severity:** `warning`  
**Auto-fix:** `no`

A script that closes the VPN tunnel from within a remote session cuts itself off whenever the admin SSH rides over that tunnel. `tailscale down`, `wg-quick down WG0`, `openvpn` teardown, and `nmcli connection down NAME` all tear the link down in place with no grace or rollback. Schedule the teardown behind `systemd-run --on-active=30s --unit=recover <cmd to bring it back up>` so the VPN is back before the unit expires, or run the command from the host's console / out-of-band path rather than over the VPN itself.

Disable by adding `ZC1829` to `disabled_katas` in `.zshellcheckrc`.

---

<a id="zc1830"></a>
### ZC1830 — Warn on `unsetopt NOMATCH` — unmatched glob becomes the literal pattern, silent bugs

**Severity:** `warning`  
**Auto-fix:** `no`

`NOMATCH` is on by default in Zsh — an unmatched glob (`*.log` with no matching files) errors out instead of silently passing through. Disabling it (`unsetopt NOMATCH` or the equivalent `setopt NO_NOMATCH`) reverts to POSIX-sh behaviour: the pattern is handed to the command verbatim, so `rm *.log` with no matches runs `rm '*.log'` — which fails noisily for `rm` but, for commands that accept arbitrary strings, silently processes the literal `*.log` instead of files. Prefer scoped `*(N)` null-glob qualifier or `setopt LOCAL_OPTIONS; setopt NULL_GLOB` inside a function, so the rest of the script keeps the default safety.

Disable by adding `ZC1830` to `disabled_katas` in `.zshellcheckrc`.

---

<a id="zc1831"></a>
### ZC1831 — Error on `systemctl stop|disable|mask ssh/sshd` — locks out the next remote login

**Severity:** `error`  
**Auto-fix:** `no`

Stopping, disabling, or masking the SSH daemon closes the door on the next remote login. Existing connections survive for a while because sshd's spawned per-session process keeps running, but any reconnect / CI follow-up step that needs to ssh back in gets `Connection refused`. `systemctl disable ssh` and `systemctl mask ssh` also survive reboots. Recovery requires console or out-of-band access. If the goal is config reload, use `systemctl reload sshd`; if the host is being retired, make sshd the last service you touch.

Disable by adding `ZC1831` to `disabled_katas` in `.zshellcheckrc`.

---

<a id="zc1832"></a>
### ZC1832 — Warn on Zsh `limit coredumpsize unlimited` — setuid memory landing in core files

**Severity:** `warning`  
**Auto-fix:** `no`

Zsh's `limit` builtin is the csh-style sibling of `ulimit`; `limit coredumpsize unlimited` is the Zsh equivalent of `ulimit -c unlimited` and has the same consequence: a crashing setuid or key-holding process leaves its address space on disk as a world-readable core file. Leave the coredump ceiling at the distro default (usually 0 for non-debug sessions), or use `systemd-coredump` with restricted permissions when you need post-mortem data. `ulimit -c unlimited` is covered by ZC1495; this kata catches the Zsh-specific `limit`/`unlimit coredumpsize` spelling.

Disable by adding `ZC1832` to `disabled_katas` in `.zshellcheckrc`.

---

<a id="zc1833"></a>
### ZC1833 — Warn on `unsetopt WARN_CREATE_GLOBAL` — silent accidental-global bugs inside functions

**Severity:** `warning`  
**Auto-fix:** `no`

`WARN_CREATE_GLOBAL` makes Zsh warn when a function assigns to a name that is not declared `local` / `typeset` in the current scope — the single highest-value guardrail against the classic Bash-ism where a helper function silently stomps on a caller's variable (`tmp=`, `i=`, `result=`). Disabling it (`unsetopt WARN_CREATE_GLOBAL` or the equivalent `setopt NO_WARN_CREATE_GLOBAL`) reverts to permissive behaviour: every unqualified assignment inside a function escapes to global scope with no diagnostic. Leave the option on and fix the offending function by adding `local` / `typeset` declarations, or — if you really must silence it for a specific block — use `setopt LOCAL_OPTIONS; unsetopt WARN_CREATE_GLOBAL` inside a function so the rest of the script keeps the safety.

Disable by adding `ZC1833` to `disabled_katas` in `.zshellcheckrc`.

---

<a id="zc1834"></a>
### ZC1834 — Error on `tc qdisc … root netem loss 100%` — hard blackhole on a live interface

**Severity:** `error`  
**Auto-fix:** `no`

`tc qdisc add/replace dev IFACE root netem loss 100%` (also `corrupt 100%` or `duplicate 100%` with no buffering) installs a Linux kernel qdisc that drops every outbound packet on the named interface. Running this on the interface that carries your SSH session is the canonical way to lock yourself out of a remote host — the `tc` command returns success, the kernel happily applies the rule, and the next TCP segment ACK never arrives. Even on the console it halts any process that depends on the interface. Stage netem experiments on a secondary interface, wrap them in `at now + 5 minutes` (or a `timeout … tc qdisc del …` recovery trap) so a partial failure does not leave the link permanently black-holed.

Disable by adding `ZC1834` to `disabled_katas` in `.zshellcheckrc`.

---

<a id="zc1835"></a>
### ZC1835 — Warn on `smartctl -s off` — drive self-monitoring (SMART) disabled, silent failure

**Severity:** `warning`  
**Auto-fix:** `no`

`smartctl -s off DEV` tells the drive firmware to stop recording the SMART attribute counters that warn operators about pending failure — reallocated sectors, pending sectors, uncorrectable errors, temperature excursions. Rotating disks and SSDs both ship with the monitoring on; disabling it keeps `smartctl -H` reporting PASSED right up until the drive falls off the bus, so the periodic fleet health scan never escalates until data loss is already happening. Use `smartctl -s on DEV` (default) and configure `smartd.conf` for proactive alerts instead of muting the source.

Disable by adding `ZC1835` to `disabled_katas` in `.zshellcheckrc`.

---

<a id="zc1836"></a>
### ZC1836 — Error on `helm uninstall --no-hooks` — skips pre-delete cleanup, orphaned state

**Severity:** `error`  
**Auto-fix:** `no`

`helm uninstall RELEASE --no-hooks` (also spelled `helm delete --no-hooks` on Helm v2 / `helm3 --no-hooks` interchangeably) tears down every chart-rendered resource but silently skips the release's `pre-delete` and `post-delete` Jobs / ConfigMap hooks. Those hooks are where production charts flush write-ahead logs, deregister service-discovery entries, back up PVC content before the PVC goes away, and release external locks — skipping them on a live release is one of the classic ways to leave the cluster in a partially deleted state with no way to replay the cleanup. Drop `--no-hooks` and let the chart run as designed; if a hook is genuinely wedged, disable it at the chart level with `helm.sh/hook-delete-policy: before-hook-creation,hook-succeeded`.

Disable by adding `ZC1836` to `disabled_katas` in `.zshellcheckrc`.

---

<a id="zc1837"></a>
### ZC1837 — Error on `chmod` granting non-owner access to `/dev/kvm` / `/dev/mem` / `/dev/kmem` / `/dev/port`

**Severity:** `error`  
**Auto-fix:** `no`

Distros ship `/dev/mem`, `/dev/kmem`, `/dev/port`, and `/dev/kvm` with tight owner-only or group-only permissions managed by udev rules — these nodes hand any process that can read or write them the keys to the kingdom (physical memory, kernel memory, raw I/O ports, full hypervisor API). Flipping the mode from a script (`chmod 666 /dev/kvm`, `chmod a+rw /dev/mem`) is a classic local privilege-escalation vector dressed up as a convenience fix for a permission error. Fix the actual problem: add the user to the `kvm` group, ship a proper udev rule (`/etc/udev/rules.d/*.rules`), or grant the specific capability the tool needs instead of blanket-chmod-ing the device.

Disable by adding `ZC1837` to `disabled_katas` in `.zshellcheckrc`.

---

<a id="zc1838"></a>
### ZC1838 — Warn on `setopt GLOB_DOTS` — bare `*` silently starts matching hidden files

**Severity:** `warning`  
**Auto-fix:** `no`

`GLOB_DOTS` off is the Zsh default: patterns like `*`, `*.log`, and recursive `**/*` skip filenames that begin with a dot (`.git/`, `.env`, `.ssh/`). Setting `setopt GLOB_DOTS` script-wide reverses that quietly — every subsequent glob now also matches hidden entries, which turns routine maintenance lines (`rm *`, `cp -r * /backup`, `chmod 644 *`) into repository-wiping, secret-copying, permission-flipping bugs. Leave the option alone at the script level and request dot-inclusion per-glob with the Zsh-native `*(D)` qualifier (or `.* *` when you explicitly want both), so the effect is scoped to the exact line that needs it.

Disable by adding `ZC1838` to `disabled_katas` in `.zshellcheckrc`.

---

<a id="zc1839"></a>
### ZC1839 — Warn on `timedatectl set-ntp false` / disabling `systemd-timesyncd` / `chronyd`

**Severity:** `warning`  
**Auto-fix:** `no`

`timedatectl set-ntp false` (also spelled `set-ntp no` / `set-ntp 0`) tells systemd to stop the network time client; `systemctl disable systemd-timesyncd` and `systemctl disable chronyd` / `ntpd` have the same effect. With no time source the hardware clock drifts, and within days TLS handshakes begin failing `notBefore`/`notAfter` checks, Kerberos tickets refuse to validate, time-based one-time passwords go out of sync, and log entries arrive in the wrong order — all silently, because the original command succeeded. Keep NTP enabled in production; if you really need a frozen clock for reproducibility, isolate it to a namespace or CI container rather than the host.

Disable by adding `ZC1839` to `disabled_katas` in `.zshellcheckrc`.

---

<a id="zc1840"></a>
### ZC1840 — Error on `openssl enc -k PASSWORD` — legacy flag embeds secret in argv

**Severity:** `error`  
**Auto-fix:** `no`

`openssl enc -k PASSWORD` (the pre-OpenSSL-3 short form of `-pass pass:PASSWORD`) takes the password directly as the next argv element — which makes it visible to every `ps` reader, every `/proc/<pid>/cmdline` consumer, shell history, and anything that logs command invocations. The same leak applies to `openssl rsa`, `openssl pkcs12`, and other subcommands that still accept the deprecated `-k` alias. Use `-pass env:VARNAME`, `-pass file:PATH`, or `-pass fd:N` (read from an open descriptor) so the secret never rides in the process argument vector.

Disable by adding `ZC1840` to `disabled_katas` in `.zshellcheckrc`.

---

<a id="zc1841"></a>
### ZC1841 — Error on `curl --proxy-insecure` — TLS verification disabled on the proxy hop

**Severity:** `error`  
**Auto-fix:** `no`

`curl --proxy-insecure` (alias of `-k` but scoped to the proxy leg, introduced alongside `--proxy-cacert` in curl 7.52) tells curl to accept any certificate presented by the HTTPS proxy that sits between the script and the origin server. The origin TLS handshake is still validated, which makes the issue easy to miss in review, but any box that can intercept traffic to the proxy — a captive portal, a rogue WPAD auto-discovery, an attacker on the same VLAN — can present its own cert and read or rewrite the tunnel contents, including any `Authorization:` header attached to the request. Install the proxy's CA bundle and point `--proxy-cacert` / `CURL_CA_BUNDLE` at it instead.

Disable by adding `ZC1841` to `disabled_katas` in `.zshellcheckrc`.

---

<a id="zc1842"></a>
### ZC1842 — Warn on `setopt CDABLE_VARS` — `cd NAME` silently falls back to `cd $NAME`

**Severity:** `warning`  
**Auto-fix:** `no`

With `CDABLE_VARS` on, any `cd NAME` whose `NAME` does not exist as a directory is retried as `cd ${NAME}` — if a parameter of the same name is set, the working directory silently jumps to wherever the variable points. A typo like `cd cinfig` (intent: `config`) suddenly lands inside `${cinfig}` when one exists, and every later relative path in the script is computed from the wrong root. Keep this option inside `~/.zshrc` where it is an interactive shortcut; in scripts, always `cd "$dir"` explicitly and pair with `|| exit` so a missed directory fails loudly instead of rewriting `$PWD`.

Disable by adding `ZC1842` to `disabled_katas` in `.zshellcheckrc`.

---

<a id="zc1843"></a>
### ZC1843 — Warn on `docker/podman run --cgroup-parent=/system.slice|/init.scope|/` — container escapes engine limits

**Severity:** `warning`  
**Auto-fix:** `no`

`--cgroup-parent=PATH` places the container under the given cgroup parent, which is normally `/docker` (or the engine's managed slice) and inherits the engine-wide memory/CPU/IO caps. Pointing the flag at `/`, `/system.slice`, or any host-managed slice puts the container side-by-side with systemd services — the engine's defaults no longer apply, and a runaway container can starve `sshd` or the kubelet for resources. Unless a specific orchestrator is supplying a managed cgroup path, drop the flag and let the engine choose; if you need custom limits, use `--memory` / `--cpus` / `--pids-limit` on the run itself.

Disable by adding `ZC1843` to `disabled_katas` in `.zshellcheckrc`.

---

<a id="zc1844"></a>
### ZC1844 — Warn on `logger -p local0.info|local7.notice|…` — unreserved facility often uncollected

**Severity:** `warning`  
**Auto-fix:** `no`

The eight `local0`–`local7` syslog facilities are reserved for site-specific use. Most distro `rsyslog` and `systemd-journald` defaults do not route them anywhere — they drop on the floor unless someone dropped a matching rule into `/etc/rsyslog.d/*.conf`. Scripts that call `logger -p local0.info 'audit: user added to wheel'` therefore log to nothing in the audit trail on a stock machine. For portable audit-style logging use the POSIX-reserved `auth.notice` or `authpriv.info` facility; for application events, pass `-t TAG` and use `user.notice` (the default) or a dedicated journald unit.

Disable by adding `ZC1844` to `disabled_katas` in `.zshellcheckrc`.

---

<a id="zc1845"></a>
### ZC1845 — Warn on `setopt PATH_DIRS` — slash-bearing command names fall back to `$PATH` lookup

**Severity:** `warning`  
**Auto-fix:** `no`

`PATH_DIRS` (off by default) changes how Zsh resolves a command that contains a `/`: instead of treating `./foo/bar` or `subdir/cmd` as a direct path, Zsh walks `$path` and retries `${path[i]}/subdir/cmd` until one is executable. The surface intent — run a local binary — is silently replaced by `/usr/local/bin/subdir/cmd` or any other same-shaped subtree that exists on `$PATH`. This gets even worse on shared build hosts where `$PATH` contains user-owned directories. Leave the option off and call local binaries with an explicit leading `./`, or hand the full absolute path to the shell.

Disable by adding `ZC1845` to `disabled_katas` in `.zshellcheckrc`.

---

<a id="zc1846"></a>
### ZC1846 — Warn on `certbot … --force-renewal` — bypasses ACME rate-limit safety

**Severity:** `warning`  
**Auto-fix:** `no`

`certbot renew --force-renewal` and `certbot certonly --force-renewal` reissue a certificate regardless of remaining validity. Placed in a daily cron, the same hostname burns through Let's Encrypt's per-domain rate limits (50 certificates per registered domain per 7 days, 5 duplicate certificates per domain per 7 days); once the limit trips, no cert for that host — fresh or renewal — can be issued until the rolling window expires, which often happens during an outage when you need it least. Drop `--force-renewal` and let certbot's default 30-days-before-expiry gate do its job, or if you really need a specific reissue, run it once manually.

Disable by adding `ZC1846` to `disabled_katas` in `.zshellcheckrc`.

---

<a id="zc1847"></a>
### ZC1847 — Warn on `setopt CHASE_LINKS` — every `cd` silently swaps symlink paths for the real inode

**Severity:** `warning`  
**Auto-fix:** `no`

`CHASE_LINKS` off is the Zsh default: `cd releases/current` leaves `$PWD` as the logical path the user typed, and `cd ..` steps back up through the symlink to where they came from. Turning the option on globally makes every `cd` resolve the target to its physical inode — so `cd releases/current` lands in `/srv/app/releases/20260415-deadbeef`, and the next `cd ../config` looks for `/srv/app/releases/config` instead of the `/srv/app/config` that the user expected. Scripts that rely on blue/green-style `current` symlinks break silently. Keep the option off at the script level and request one-shot physical resolution with `cd -P target` when a specific call needs it.

Disable by adding `ZC1847` to `disabled_katas` in `.zshellcheckrc`.

---

<a id="zc1848"></a>
### ZC1848 — Warn on `ssh -o CheckHostIP=no` — DNS-spoof warning for known hosts silenced

**Severity:** `warning`  
**Auto-fix:** `no`

`CheckHostIP` (on by default) stores the host's IP address alongside its host key in `~/.ssh/known_hosts`; if DNS later resolves the same name to a different IP but the key still matches, ssh warns you. Turning the check off with `-o CheckHostIP=no` keeps the host-key comparison but silences the IP-mismatch warning — which means a DNS-poisoning attacker who already holds the previously-seen host key (stolen, misplaced backup, leaked by a decommissioned box) can route the session through their box without a peep. Leave the default, and if you really need to skip the IP record (load-balanced pool with shared keys) document the risk and prefer `HostKeyAlias` instead.

Disable by adding `ZC1848` to `disabled_katas` in `.zshellcheckrc`.

---

<a id="zc1849"></a>
### ZC1849 — Warn on `setopt ALL_EXPORT` — every later `var=value` silently becomes `export var=value`

**Severity:** `warning`  
**Auto-fix:** `no`

`ALL_EXPORT` (POSIX `set -a` equivalent, off by default) tells Zsh to mark every parameter assignment for export as soon as it is created, so `password=$(cat secret)` immediately rides into the environment of every child process the script spawns — the `ps e`, `/proc/<pid>/environ`, and journal of any later `| tee`, `| mail`, or `logger` call. Enabling it script-wide to avoid a few `export` keywords leaks credentials and private config by default. Drop the `setopt`, scope exports explicitly with `export VAR=value`, or wrap a narrow section in `setopt LOCAL_OPTIONS; setopt ALL_EXPORT` inside a function so the effect cannot leak past the closing brace.

Disable by adding `ZC1849` to `disabled_katas` in `.zshellcheckrc`.

---

<a id="zc1850"></a>
### ZC1850 — Warn on `ssh -o LogLevel=QUIET` — silences security-relevant ssh diagnostics

**Severity:** `warning`  
**Auto-fix:** `no`

`LogLevel=QUIET` (aliased to the `-q` short flag) suppresses every informational or warning message ssh would otherwise print: host-key changes, key-exchange downgrades, agent-forwarding permission denials, canonical-hostname rewrites. In a script, that means the output looks clean even when ssh is shouting about a MITM on the other end. Keep the default `INFO` level (or raise to `VERBOSE` during debugging), capture stderr to a log if the noise bothers you, and never pair `LogLevel=QUIET` with `StrictHostKeyChecking=no` in the same call — that combination actively hides known-bad-key events.

Disable by adding `ZC1850` to `disabled_katas` in `.zshellcheckrc`.

---

<a id="zc1851"></a>
### ZC1851 — Warn on `unsetopt FUNCTION_ARGZERO` — `$0` inside a function stops reporting the function name

**Severity:** `warning`  
**Auto-fix:** `no`

`FUNCTION_ARGZERO` is Zsh's default: inside a function or `source`d file, `$0` holds the function/file name, which is what every `log_error "$0: ..."` helper, every self-reflecting `$funcfiletrace` fallback, and every `case $0` dispatcher expects. Turning it off reverts to POSIX-sh behaviour where `$0` always points at the outer script — so `my_func() { echo "${0}: bad input" }` silently starts logging `myscript.sh: bad input` for every function, which makes stack-trace logs unreadable and breaks dispatchers that branch on `$0`. Keep the option on at the script level and, if one specific helper needs the POSIX name, reach it explicitly with `$ZSH_ARGZERO` or `$ZSH_SCRIPT`.

Disable by adding `ZC1851` to `disabled_katas` in `.zshellcheckrc`.

---

<a id="zc1852"></a>
### ZC1852 — Error on `firewall-cmd --panic-on` — firewalld drops every packet, kills the SSH session

**Severity:** `error`  
**Auto-fix:** `no`

`firewall-cmd --panic-on` puts firewalld into panic mode, which drops every inbound and outbound packet regardless of zone or rule. Running this over a remote SSH session is the textbook way to lock yourself out: the command returns success, the TCP ACK for that reply never arrives, and nobody can reach the host until someone visits the console to `--panic-off`. Stage panic-mode experiments on a machine you can power-cycle, gate the call behind `at now + 5 minutes` with an auto-disable, or use targeted zone rules instead of the blanket switch.

Disable by adding `ZC1852` to `disabled_katas` in `.zshellcheckrc`.

---

<a id="zc1853"></a>
### ZC1853 — Warn on `setopt MARK_DIRS` — glob-matched directories gain a silent trailing `/`

**Severity:** `warning`  
**Auto-fix:** `no`

With `MARK_DIRS` on, every filename produced by a glob that resolves to a directory picks up a trailing `/`. Inside a shell it looks harmless, but scripts that pass the glob result to other tools break in quiet ways: `[[ -f "$f" ]]` rejects `dir/` because it is not a regular file, `rm -f *` sees `dir/` and silently skips it (GNU rm refuses to remove directories without `-r`), and downstream hash maps indexed on basenames suddenly carry two keys for what the user thinks is one entry. Keep the option off at the script level and request the trailing slash per-glob with the `(/)` qualifier (`dirs=( *(/) )`) when you really need directories only.

Disable by adding `ZC1853` to `disabled_katas` in `.zshellcheckrc`.

---

<a id="zc1854"></a>
### ZC1854 — Error on `yum-config-manager --add-repo http://…` / `zypper addrepo http://…` — plaintext repo allows MITM

**Severity:** `error`  
**Auto-fix:** `no`

Adding a package repository over plain HTTP (`yum-config-manager --add-repo http://…`, `dnf config-manager --add-repo http://…`, `zypper addrepo http://…`) tells the package manager to fetch metadata and RPMs without TLS — any on-path attacker can substitute packages, and even GPG signature checks do not help because the attacker can simply strip the `repo_gpgcheck=1` line from the unsigned `.repo` file. Use the `https://` mirror (every major distro now publishes one), or pin to a local mirror over TLS and drop a `gpgkey=file:///etc/pki/...` entry in the same `.repo` so signatures cannot be disabled mid-install.

Disable by adding `ZC1854` to `disabled_katas` in `.zshellcheckrc`.

---

<a id="zc1855"></a>
### ZC1855 — Avoid `$GROUPS` — Bash-only array; Zsh exposes supplementary groups as `$groups`

**Severity:** `warning`  
**Auto-fix:** `no`

`$GROUPS` is a Bash magic parameter that holds the caller's supplementary GIDs as a numeric array. Zsh does not populate `$GROUPS`; it has `$groups`, a lowercase associative array keyed by group *name* with the GID as value (`${(k)groups}` for names, `${(v)groups}` for IDs). Scripts ported from Bash that iterate `${GROUPS[@]}` therefore see an empty list under Zsh and silently skip group-membership checks. Use `${(k)groups}` for names or `${(v)groups}` for numeric GIDs; the Zsh `id -Gn` fallback keeps the script portable across shells.

Disable by adding `ZC1855` to `disabled_katas` in `.zshellcheckrc`.

---

<a id="zc1856"></a>
### ZC1856 — Warn on `unset arr[N]` — Zsh does not delete the array element, the array keeps its length

**Severity:** `warning`  
**Auto-fix:** `no`

In Bash, `unset arr[N]` removes the N-th element of the array (leaving a sparse hole). In Zsh the same invocation passes the literal string `arr[N]` to the `unset` builtin, which looks for a parameter with that name — finds nothing — and returns success. The array is left untouched, `${#arr[@]}` does not budge, and every downstream `for x in "${arr[@]}"` keeps iterating the element the script thought it had removed. Use Zsh's native assignment form `arr[N]=()` to delete an index, or `arr=("${(@)arr:#pattern}")` to filter by value.

Disable by adding `ZC1856` to `disabled_katas` in `.zshellcheckrc`.

---

<a id="zc1857"></a>
### ZC1857 — Error on `cloud-init clean` — wipes boot state, next reboot re-provisions the host

**Severity:** `error`  
**Auto-fix:** `no`

`cloud-init clean` (and variants `--logs`, `--reboot`, `--machine-id`) removes every marker under `/var/lib/cloud/` and `/var/log/cloud-init*`, which tells cloud-init to re-run from scratch on the next boot. That run re-imports the image-builder's user-data: regenerates SSH host keys, resets the hostname, replaces `/etc/fstab` entries the operator may have edited, and (with `--reboot`) triggers the replay immediately. In a maintenance script this silently erases everything the operator configured after first-boot. Keep the command out of automation; if you truly need to re-seed an instance, snapshot state first and run the command interactively.

Disable by adding `ZC1857` to `disabled_katas` in `.zshellcheckrc`.

---

<a id="zc1858"></a>
### ZC1858 — Error on `ssh -c 3des-cbc|arcfour|blowfish-cbc` — weak cipher forced on the tunnel

**Severity:** `error`  
**Auto-fix:** `no`

OpenSSH disables legacy ciphers by default; a script that explicitly forces one with `-c 3des-cbc`, `-c arcfour`, `-c blowfish-cbc`, or a matching entry in `-o Ciphers=...` downgrades the tunnel to an algorithm with known plaintext recovery, IV-reuse, or birthday-bound attacks. Typically this is done to reach an old appliance — but it drags every other session on the same invocation down with it. Leave cipher selection to OpenSSH's default; if a legacy device absolutely requires a weak cipher, isolate it in a `Host ...` block in `~/.ssh/config` with explicit `HostKeyAlgorithms` and keep the rest of the fleet on strong defaults.

Disable by adding `ZC1858` to `disabled_katas` in `.zshellcheckrc`.

---

<a id="zc1859"></a>
### ZC1859 — Warn on `unsetopt MULTIOS` — `cmd >a >b` silently keeps only the last redirection

**Severity:** `warning`  
**Auto-fix:** `no`

`MULTIOS` is on by default in Zsh: `cmd >out.log >>archive.log` sends stdout to both files via an implicit `tee`, and `cmd <a <b` concatenates the two inputs in order. Disabling it reverts to POSIX-sh semantics — Zsh opens each earlier redirection, closes it immediately, and only the last one in the direction wins. Any script that was written for Zsh suddenly starts dropping the `archive.log` tail, and log collectors that opened `archive.log` keep the fd but never receive new lines. Keep the option on at the script level; if one specific line really needs POSIX behaviour, wrap it in a function with `setopt LOCAL_OPTIONS; unsetopt MULTIOS`.

Disable by adding `ZC1859` to `disabled_katas` in `.zshellcheckrc`.

---

<a id="zc1860"></a>
### ZC1860 — Warn on `hostnamectl set-hostname NEW` — caches and certs still reference the old name

**Severity:** `warning`  
**Auto-fix:** `no`

`hostnamectl set-hostname NEW` (and the new-style `hostnamectl hostname NEW` and `hostname NEW`) updates `/etc/hostname` and `kernel.hostname` atomically, but every process that called `gethostname()` at startup keeps the old value until it restarts: syslog tags, Prometheus scrape labels, Docker daemons, and anything that populated a TLS `subjectAltName` with `$(hostname)` still speak as the previous host. Change the hostname interactively, then plan a restart window — in automation, prefer shipping the new hostname via cloud-init / Ignition so every service starts with it from boot.

Disable by adding `ZC1860` to `disabled_katas` in `.zshellcheckrc`.

---

<a id="zc1861"></a>
### ZC1861 — Warn on `setopt OCTAL_ZEROES` — leading-zero integers silently reinterpret as octal

**Severity:** `warning`  
**Auto-fix:** `no`

`OCTAL_ZEROES` is off in Zsh by default: arithmetic treats `0100` as the decimal integer one hundred, matching what every other scripting language does. Setting it on reverts to POSIX-shell semantics where the leading `0` flags the literal as octal — `(( n = 0100 ))` assigns 64, not 100. Scripts that read timestamps padded to `00:59`, CSVs of phone-number prefixes (`0049`), or file modes formatted as `0700` silently return the wrong integer. Keep the option off at script level; if you really want C-style octal literals, stay explicit with `(( n = 8#100 ))` or `$(( 8#$val ))` so the intent is obvious.

Disable by adding `ZC1861` to `disabled_katas` in `.zshellcheckrc`.

---

<a id="zc1862"></a>
### ZC1862 — Warn on `ssh-keygen -R HOST` — deletes a known-hosts entry, next `ssh` re-trusts silently

**Severity:** `warning`  
**Auto-fix:** `no`

`ssh-keygen -R HOST` scrubs the entry for `HOST` from `~/.ssh/known_hosts`. The legitimate trigger is a real key rotation (server reinstall, HSM replacement), but the flag is frequently dropped into automation to silence the REMOTE HOST IDENTIFICATION HAS CHANGED banner without ever confirming the new fingerprint. The very next `ssh` call then prompts once (or not at all under `StrictHostKeyChecking=no`) and blindly accepts whatever the network hands back — a MITM attacker who was waiting for a rebuild slips in without a trace. Fetch the new key out-of-band and `ssh-keyscan -t rsa,ed25519 HOST | ssh-keygen -lf -` before adding it, or pin fingerprints in a managed `known_hosts` file.

Disable by adding `ZC1862` to `disabled_katas` in `.zshellcheckrc`.

---

<a id="zc1863"></a>
### ZC1863 — Warn on `unsetopt CASE_GLOB` — globs silently go case-insensitive across the script

**Severity:** `warning`  
**Auto-fix:** `no`

`CASE_GLOB` on is the Zsh default: `*.log` matches `app.log` but not `APP.LOG`, `[A-Z]*` is a real case-sensitive range, and `[[ $f == Foo* ]]` keeps the distinction between `Foo1` and `foo1`. Turning it off (or equivalently `setopt NO_CASE_GLOB`) silently re-evaluates every subsequent pattern case-insensitively — `rm *.log` now sweeps `APP.LOG` up, pattern dispatchers that used to distinguish `README` from `readme` stop doing so, and hash maps keyed on glob-built labels start colliding. Keep the option on at script level; request case-folding per-pattern with the Zsh qualifier `(#i)*.log`.

Disable by adding `ZC1863` to `disabled_katas` in `.zshellcheckrc`.

---

<a id="zc1864"></a>
### ZC1864 — Error on `mount -o remount,exec` — re-enables exec on a previously `noexec` mount

**Severity:** `error`  
**Auto-fix:** `no`

Hardened systems mount `/tmp`, `/var/tmp`, `/dev/shm`, and `/home` with `noexec` so a dropper cannot chmod and launch a payload out of a world-writable directory. `mount -o remount,exec /tmp` (or the narrower `remount,suid`) removes that guardrail for the live kernel, and every shell that already had `cd /tmp` open picks it up immediately. Most legitimate uses come from install scripts that briefly relax `noexec`; those scripts should restore the flag in a `trap 'mount -o remount,noexec /tmp' EXIT`. Blanket `remount,exec` without a restore path leaves the system in a permanently weakened state until reboot.

Disable by adding `ZC1864` to `disabled_katas` in `.zshellcheckrc`.

---

<a id="zc1865"></a>
### ZC1865 — Warn on `unsetopt CASE_MATCH` — `[[ =~ ]]` and pattern tests quietly fold case

**Severity:** `warning`  
**Auto-fix:** `no`

`CASE_MATCH` on is Zsh's default: `[[ $x =~ ^FOO ]]`, `[[ $x == Foo* ]]`, and the subst-in-conditional forms honour letter case exactly as written. Turning the option off flips every later test to case-insensitive — `[[ $user == Admin ]]` also matches `admin`/`ADMIN`, regex dispatchers stop distinguishing `README` from `readme`, and log-pattern filters over-collect. Keep the option on at script level; if one specific regex really needs case-folding, request it per-pattern with the Zsh `(#i)` flag (e.g. `[[ $x =~ (#i)foo ]]`).

Disable by adding `ZC1865` to `disabled_katas` in `.zshellcheckrc`.

---

<a id="zc1866"></a>
### ZC1866 — Warn on `docker exec -u 0` — bypasses the image's non-root `USER` directive

**Severity:** `warning`  
**Auto-fix:** `no`

A hardened image runs with a non-root `USER` set in its Dockerfile so exploited processes inside the container are contained by the Linux user-namespace mapping. `docker exec -u 0` (and `-u root`, `--user=0`, the podman equivalent) overrides that choice on a per-exec basis and drops a shell back into uid 0 — every subsequent file write, cap check, and namespace test now runs as root inside the container, which on a default Docker setup is also root on the host via the shared mount namespace. Keep exec sessions as the container's configured user; if you genuinely need root for a one-off fix, document it in the ticket and consider rebuilding the image with the capability baked in so `-u 0` is never required.

Disable by adding `ZC1866` to `disabled_katas` in `.zshellcheckrc`.

---

<a id="zc1867"></a>
### ZC1867 — Warn on `unsetopt GLOB` — pattern expansion turned off, `rm *.log` tries the literal filename

**Severity:** `warning`  
**Auto-fix:** `no`

`GLOB` is on by default in Zsh: `*`, `?`, `[...]`, and `**/` expand against the filesystem before the command runs. Turning the option off script-wide (via `unsetopt GLOB` or the equivalent `setopt NO_GLOB`, same as POSIX `set -f`) means every later pattern is handed to the command verbatim, so `rm *.log` tries to remove a file literally named `*.log`, `for f in *.txt` iterates over the single literal string, and expected-array-length checks always return 1. Keep the option on at the script level; if one specific line needs the pattern as a literal, quote the argument (`'*.log'`) or scope with `setopt LOCAL_OPTIONS; setopt NO_GLOB` inside a function.

Disable by adding `ZC1867` to `disabled_katas` in `.zshellcheckrc`.

---

<a id="zc1868"></a>
### ZC1868 — Error on `gcloud config set auth/disable_ssl_validation true` — disables TLS on every later gcloud call

**Severity:** `error`  
**Auto-fix:** `no`

`gcloud config set auth/disable_ssl_validation true` writes the flag into the active configuration file, so every subsequent `gcloud` invocation on that machine stops verifying the Google API certificate until someone reverses it. A MITM holding a self-signed cert can then intercept service account tokens, project-level credentials, and every deploy that runs under the same user. Remove the setting (`gcloud config unset auth/disable_ssl_validation`), and if a corporate proxy really needs a custom CA use `core/custom_ca_certs_file` to pin it rather than disabling the check.

Disable by adding `ZC1868` to `disabled_katas` in `.zshellcheckrc`.

---

<a id="zc1869"></a>
### ZC1869 — Warn on `setopt RC_EXPAND_PARAM` — brace-adjacent array expansion silently distributes

**Severity:** `warning`  
**Auto-fix:** `no`

`RC_EXPAND_PARAM` is off in Zsh by default: `echo x${arr[@]}y` concatenates once, producing `xay xby xcy` only if you wrote the template carefully. Turning it on changes the rule — every adjacent literal is distributed across each array element, so `cp src/${files[@]}.bak /tmp` suddenly rewrites as `cp src/a.bak src/b.bak src/c.bak /tmp`. That is exactly what you want when you want it, and a nasty surprise anywhere else because the same syntax keeps working silently. Leave the option off at script level; if one specific line needs distributive expansion, request it per-use with `${^arr}` (the `^` flag scopes the behaviour to that parameter only).

Disable by adding `ZC1869` to `disabled_katas` in `.zshellcheckrc`.

---

<a id="zc1870"></a>
### ZC1870 — Warn on `setopt GLOB_ASSIGN` — RHS of `var=pattern` silently glob-expands

**Severity:** `warning`  
**Auto-fix:** `no`

`GLOB_ASSIGN` is off by default in Zsh: `logs=*.log` sets `$logs` to the literal string `*.log`, just like every other shell. Turning it on expands the right-hand side of unquoted assignments — `logs=*.log` silently becomes the first matching filename, `latest=backup-*` captures whatever sort-order the filesystem returns, and any empty-match case assigns an empty string. Scripts that port cleanly between Bash and Zsh suddenly diverge, and sensitive assignments like `cert=~/secrets/*` can grab attacker-dropped files. Keep the option off; use `set -A arr *.log` or explicit `arr=( *.log )` when you really want the expansion.

Disable by adding `ZC1870` to `disabled_katas` in `.zshellcheckrc`.

---

<a id="zc1871"></a>
### ZC1871 — Warn on `setopt IGNORE_BRACES` — brace expansion stops working script-wide

**Severity:** `warning`  
**Auto-fix:** `no`

`IGNORE_BRACES` is off by default in Zsh, which means `{1..10}`, `file.{log,bak}`, and nested combinations like `{a..z}{1..3}` all expand exactly as they do in Bash with `brace_expand` on. Turning it on disables every one of those — `for i in {1..10}` iterates over the single literal token `{1..10}`, and `cp app.{conf,conf.bak}` tries to copy a file literally called `app.{conf,conf.bak}`. Scripts that depend on either numeric or comma-list expansion silently become no-ops or fail with ENOENT. Keep the option off; if you really need a literal brace string, quote the specific argument (`'app.{conf,bak}'`) instead of flipping the shell-wide behaviour.

Disable by adding `ZC1871` to `disabled_katas` in `.zshellcheckrc`.

---

<a id="zc1872"></a>
### ZC1872 — Error on `badblocks -w` — destructive write-mode pattern test wipes the device

**Severity:** `error`  
**Auto-fix:** `no`

`badblocks -w` (alias `--write-mode`) runs the write-mode bad-block check, which overwrites every sector of the target device with a test pattern and reads it back. On a fresh drive about to be formatted that is exactly what you want; on an already-populated disk it is a silent data-wipe — the command returns success even as it bulldozes the filesystem. If only non-destructive checking is needed, use `badblocks -n` (read-test-restore) or `badblocks` without any mode flag (read-only). When a true destructive test is intended, gate the call behind a confirmation prompt and a freshly partitioned device.

Disable by adding `ZC1872` to `disabled_katas` in `.zshellcheckrc`.

---

<a id="zc1873"></a>
### ZC1873 — Warn on `setopt ERR_RETURN` — functions silently bail out on the first non-zero exit

**Severity:** `warning`  
**Auto-fix:** `no`

`ERR_RETURN` is the function-scoped cousin of `ERR_EXIT` and is off by default in Zsh. Turning it on script-wide makes every function `return` at the first command whose status is non-zero, which in practice means helpers that deliberately probe the environment (`test -f /some/file`, `grep -q PATTERN`, `id -u user`) will bail before they reach the branch that was meant to run when the probe failed. Callers see a success-or-nothing return and no stderr. Keep the option off at script level; inside one function that really wants fail-fast semantics, scope with `setopt LOCAL_OPTIONS; setopt ERR_RETURN` so the behaviour cannot leak to the rest of the shell.

Disable by adding `ZC1873` to `disabled_katas` in `.zshellcheckrc`.

---

<a id="zc1874"></a>
### ZC1874 — Warn on `sshuttle -r HOST 0/0` — every outbound packet tunneled through the jump host

**Severity:** `warning`  
**Auto-fix:** `no`

`sshuttle -r user@host 0/0` (or `0.0.0.0/0`, `::/0`) installs a VPN-like catch-all route: every TCP connection and DNS lookup on the local machine egresses through `user@host`, including traffic to corporate VPN endpoints, cloud APIs, and package mirrors that had been explicitly split-tunnel. If the jump host is compromised, misconfigured, or simply overloaded, every session on the workstation silently degrades or leaks to the wrong peer. Scope the subnet list to the networks you actually need (`10.0.0.0/8 172.16.0.0/12 192.168.0.0/16`), or prefer `ssh -D` with `--exclude` rules for a single browser profile.

Disable by adding `ZC1874` to `disabled_katas` in `.zshellcheckrc`.

---

<a id="zc1875"></a>
### ZC1875 — Warn on `setopt RC_QUOTES` — `''` inside single quotes flips from empty-concat to literal apostrophe

**Severity:** `warning`  
**Auto-fix:** `no`

`RC_QUOTES` is off by default in Zsh: inside a single-quoted string `'it''s'` parses as two adjacent single-quoted regions with an empty middle, producing the literal `its`. Turning the option on reinterprets the doubled apostrophe as one escaped quote, so `'it''s'` suddenly becomes `it's`. That is a source-level change to every already-written string literal in the file — password strings, SQL fragments, display text — so log lines, stored tokens, and API payloads silently diverge. Keep the option off; write a literal apostrophe with `\'` outside the quotes or with double-quoted wrapping.

Disable by adding `ZC1875` to `disabled_katas` in `.zshellcheckrc`.

---

<a id="zc1876"></a>
### ZC1876 — Warn on `cargo publish --allow-dirty` — publishes the crate with uncommitted local changes

**Severity:** `warning`  
**Auto-fix:** `no`

`cargo publish` by default refuses to upload when the working tree is dirty, because the published tarball is a snapshot of whatever is on disk — not whatever is committed. `--allow-dirty` skips that check, so a `println!` dropped in for debugging, an uncommitted `Cargo.toml` dep bump, or a `patch.crates-io` override that only exists locally ends up on crates.io under the same version users see on GitHub. This is irreversible — once a version is uploaded it cannot be replaced, only yanked. Commit first and publish from a clean checkout; if you truly must publish from a dirty tree, scope the flag to a one-off manual call with a `--dry-run` rehearsal first.

Disable by adding `ZC1876` to `disabled_katas` in `.zshellcheckrc`.

---

<a id="zc1877"></a>
### ZC1877 — Warn on `unsetopt SHORT_LOOPS` — short-form `for`/`while` bodies stop parsing

**Severity:** `warning`  
**Auto-fix:** `no`

`SHORT_LOOPS` is on in Zsh by default: the compact forms `for x in *.log; print $x`, `while true; print .`, and `repeat 3 sleep 1` parse with an implicit single-command body. Turning the option off reverts to POSIX-shell parsing, which demands an explicit `do ... done` or `{ ... }` block. Every subsequent short-form loop raises a parse error (`parse error near '\n'`), and the behaviour is global so even helper files sourced later fall over. Keep the option on; if you genuinely need POSIX-strict parsing, scope inside a function with `setopt LOCAL_OPTIONS; unsetopt SHORT_LOOPS`.

Disable by adding `ZC1877` to `disabled_katas` in `.zshellcheckrc`.

---

<a id="zc1878"></a>
### ZC1878 — Warn on `kubectl apply --force-conflicts` — steals ownership of fields managed by other controllers

**Severity:** `warning`  
**Auto-fix:** `no`

Server-side apply tracks every field of a resource by the applier that last set it (`metadata.managedFields`). When two appliers disagree, the default behaviour is to abort with `conflict` so you can reconcile deliberately. `kubectl apply --server-side --force-conflicts` overrides that: the current caller snatches ownership of every conflicting field — including fields set by operators, HPA, cert-manager, and webhook-injected sidecars — and those controllers will silently lose their reconcile pressure until their next write. Resolve the conflict instead: either drop the disputed fields from your manifest so the other owner can keep them, or coordinate a hand-off by first removing the managed-field entry (`kubectl apply --field-manager=... --subresource=...`).

Disable by adding `ZC1878` to `disabled_katas` in `.zshellcheckrc`.

---

<a id="zc1879"></a>
### ZC1879 — Warn on `unsetopt BAD_PATTERN` — malformed glob patterns silently pass through as literals

**Severity:** `warning`  
**Auto-fix:** `no`

`BAD_PATTERN` is on in Zsh by default: a syntactically broken glob (unbalanced `[`, stray `^` outside extended-glob context, runaway `(alt|…`) produces a `zsh: bad pattern` error so the script knows the filename filter is wrong. Turning the option off reverts to POSIX behaviour — the pattern is handed to the command verbatim, and `rm [abc` silently tries to remove a file literally called `[abc`. Malformed patterns routed to `find -name` or passed to `case` blocks likewise stop firing. Keep the option on at script level; if one particular line really needs POSIX pass-through, quote the pattern or scope with `setopt LOCAL_OPTIONS; unsetopt BAD_PATTERN` inside a function.

Disable by adding `ZC1879` to `disabled_katas` in `.zshellcheckrc`.

---

<a id="zc1880"></a>
### ZC1880 — Warn on `kubectl annotate|label --overwrite` — silently rewrites controller signals

**Severity:** `warning`  
**Auto-fix:** `no`

Kubernetes annotations and labels are not plain metadata — they are the protocol by which cert-manager, external-dns, ingress-nginx, the HorizontalPodAutoscaler, and most Helm-managed controllers decide what to do with a resource. `kubectl annotate --overwrite` and `kubectl label --overwrite` suppress the conflict check and replace whatever value was there, so the script silently rewrites `kubectl.kubernetes.io/last-applied-configuration`, `cert-manager.io/cluster-issuer`, or `prometheus.io/scrape`, triggering reissue / reconfiguration or breaking the next apply. Inspect the existing annotation with `kubectl get -o jsonpath='{.metadata.annotations}'` first, and drop `--overwrite` so a conflict surfaces as an error.

Disable by adding `ZC1880` to `disabled_katas` in `.zshellcheckrc`.

---

<a id="zc1881"></a>
### ZC1881 — Warn on `unsetopt MULTIBYTE` — `${#str}`, substring, and `[[ =~ ]]` stop counting characters

**Severity:** `warning`  
**Auto-fix:** `no`

`MULTIBYTE` is on in Zsh by default: `${#str}` returns character count, `${str:0:3}` extracts the first three characters, and `[[ $str =~ ... ]]` matches whole UTF-8 codepoints. Turning it off reverts every string operation to per-byte math, so an emoji that encodes to four bytes counts as four, a substring spanning a multi-byte character slices mid-codepoint and produces invalid UTF-8, and `[[ =~ ]]` regex ranges no longer cover Unicode blocks. Filenames containing non-ASCII, i18n log strings, and JSON snippets silently drift from their assumed layout. Keep the option on; if you truly need byte-level counting, use `${#${(%)str}}` or `wc -c <<< $str`.

Disable by adding `ZC1881` to `disabled_katas` in `.zshellcheckrc`.

---

<a id="zc1882"></a>
### ZC1882 — Warn on `sudo -s` / `sudo su` / `sudo bash` — spawns an interactive root shell from a script

**Severity:** `warning`  
**Auto-fix:** `no`

`sudo -s`, `sudo -i`, `sudo su [-]`, and `sudo bash` (or `zsh`/`sh`/`ksh`) with no trailing command hand you an interactive root shell. That is fine at a prompt, but in a non-interactive script the shell either hangs waiting for stdin or drains stdin into root's shell as if those lines were the shell's commands — neither is what the script author meant. Pass the actual command to sudo (`sudo /usr/local/bin/provision.sh`) so the elevation is scoped and audit logs capture the real work.

Disable by adding `ZC1882` to `disabled_katas` in `.zshellcheckrc`.

---

<a id="zc1883"></a>
### ZC1883 — Warn on `setopt PATH_SCRIPT` — `. ./script.sh` silently falls back to `$PATH` lookup

**Severity:** `warning`  
**Auto-fix:** `no`

`PATH_SCRIPT` (off by default) lets the `.` builtin and `source` fall back to a `$PATH` walk when the literal path resolves to no file. With it on, `. helper.sh` looks for `helper.sh` in every `$path` entry — including user-owned directories like `~/bin` or `./` — and silently sources whichever matches first. An attacker who can drop `helper.sh` into any `$PATH` component runs their code inside the current shell's process, with every parent env var and exported secret available. Keep the option off; always source scripts with an explicit path (`./helper.sh`, `/opt/…/helper.sh`) so the source cannot be redirected.

Disable by adding `ZC1883` to `disabled_katas` in `.zshellcheckrc`.

---

<a id="zc1884"></a>
### ZC1884 — Error on `curl/wget https://...?apikey=...` — credential in URL query string

**Severity:** `error`  
**Auto-fix:** `no`

Anything passed as an HTTP query parameter is logged by every intermediary: the server's access log, the transparent proxy, the CDN request-id trail, browser referrer headers, and any client-side observability tooling. A URL like `https://api.example/widgets?apikey=SECRET&token=xyz` therefore tattoos the credential into logs that live forever and are often shared with downstream teams. Move the secret into an HTTP header (`curl -H "Authorization: Bearer $TOKEN"`), a POST body with `--data-urlencode` + TLS, or an `-u user:` basic-auth combo — never the query string.

Disable by adding `ZC1884` to `disabled_katas` in `.zshellcheckrc`.

---

<a id="zc1885"></a>
### ZC1885 — Warn on `setopt CSH_NULL_GLOB` — unmatched globs drop instead of erroring when any sibling matches

**Severity:** `warning`  
**Auto-fix:** `no`

`CSH_NULL_GLOB` (off by default) mimics csh's rule: in a list like `rm *.log *.bak *.tmp`, if at least one pattern produces matches the remaining unmatched patterns are silently discarded, and only if every pattern produces nothing does the shell raise `no match`. That is a partial-failure concealer — a genuine typo `rm *.lg *.bak` can still delete the `.bak` files while hiding the `.lg` mismatch, and maintenance loops that relied on `NOMATCH` to stop on typos pass right through. Keep the option off at script level; use `*(N)` per-glob when you want null-glob behaviour.

Disable by adding `ZC1885` to `disabled_katas` in `.zshellcheckrc`.

---

<a id="zc1886"></a>
### ZC1886 — Error on `tee/cp/mv/install/dd` writing system shell-init files — persistent privesc surface

**Severity:** `error`  
**Auto-fix:** `no`

`/etc/profile`, `/etc/bash.bashrc`, `/etc/zshrc`, `/etc/zsh/zshenv`, `/etc/environment`, and every drop-in under `/etc/profile.d/` are sourced by every interactive shell (and `/etc/zshenv` by every Zsh invocation). A script that `tee`s, `cp`s, `mv`s, or `dd`s arbitrary content into any of those paths becomes a persistent foothold — the next root login runs the injected code. These files belong to the packaging system; hand-edit carefully, stage a temp file, validate it with a dry-run login, and move it into place with an atomic `install -m 644`.

Disable by adding `ZC1886` to `disabled_katas` in `.zshellcheckrc`.

---

<a id="zc1887"></a>
### ZC1887 — Warn on `setopt POSIX_TRAPS` — EXIT/ZERR traps change scope and no longer fire on function return

**Severity:** `warning`  
**Auto-fix:** `no`

`POSIX_TRAPS` is off by default in Zsh. With it off, `trap cleanup EXIT` inside a function fires when that function returns — the idiomatic Zsh way to scope cleanup to a scope. Turning the option on reverts to POSIX-sh semantics, where the EXIT trap only fires when the whole shell exits and is shared across the entire process. Scripts that installed a cleanup trap inside `do_work()` expecting it to run at each invocation now leak the first trap's handler into everything after, and helpers that counted on TRAPZERR / TRAPEXIT function-scoped behaviour silently skip. Keep the option off at script level; if a specific line really needs POSIX-scope, use `trap … EXIT` at top level and document it.

Disable by adding `ZC1887` to `disabled_katas` in `.zshellcheckrc`.

---

<a id="zc1888"></a>
### ZC1888 — Warn on `aws iam create-access-key` — mints long-lived static AWS credentials

**Severity:** `warning`  
**Auto-fix:** `no`

`aws iam create-access-key` hands out a static `AKIA.../secret` pair that is valid forever until someone rotates it; whoever gets the pair speaks for the IAM user on every API call AWS accepts. Most modern deploys no longer need these: EC2 instance profiles, EKS/IRSA, Lambda roles, GitHub OIDC, and IAM Identity Center all hand out short-lived session credentials on demand. Prefer those; if a static key is genuinely required (legacy third-party tooling), store it in AWS Secrets Manager, scope the user to the narrowest policy possible, and rotate on a schedule with `aws iam update-access-key --status Inactive` / `delete-access-key`.

Disable by adding `ZC1888` to `disabled_katas` in `.zshellcheckrc`.

---

<a id="zc1889"></a>
### ZC1889 — Error on `skopeo copy --src-tls-verify=false` / `--dest-tls-verify=false` — MITM on image copy

**Severity:** `error`  
**Auto-fix:** `no`

`skopeo copy` is the glue for promoting container images between registries in CI, mirroring upstream images into internal caches, and rehydrating images to an air-gapped registry. `--src-tls-verify=false` and `--dest-tls-verify=false` drop certificate verification on the respective leg, which means any on-path attacker can substitute a malicious manifest or layer and the copy completes without a warning. Use `--src-cert-dir`/`--dest-cert-dir` to pin a private CA if you are mirroring to or from an internal registry with self-signed certs, or fix the upstream's cert.

Disable by adding `ZC1889` to `disabled_katas` in `.zshellcheckrc`.

---

<a id="zc1890"></a>
### ZC1890 — Error on `kadmin -w PASS` / `kinit` with password arg — Kerberos password in argv

**Severity:** `error`  
**Auto-fix:** `no`

`kadmin -w PASS` and `kadmin.local -w PASS` pass the Kerberos admin principal's password directly as an argv element. Every `ps`, `/proc/<pid>/cmdline`, history file, and CI-pipeline log therefore sees it in plaintext, which is catastrophic for an account that can edit the realm's KDC. Use `-k -t /etc/krb5.keytab` for non-interactive auth (keytab permissioned to root only), or pipe the password through stdin with the `-q` batch form so it never rides in argv.

Disable by adding `ZC1890` to `disabled_katas` in `.zshellcheckrc`.

---

<a id="zc1891"></a>
### ZC1891 — Error on `kubectl config view --raw` — prints the full kubeconfig with client keys

**Severity:** `error`  
**Auto-fix:** `no`

`kubectl config view` by default redacts secrets: `client-certificate-data`, `client-key-data`, `token`, and `password` fields are replaced with `REDACTED`. Adding `--raw` (or the synonym `-R`) undoes every redaction and prints the client's base64-encoded private key, bearer tokens, and any embedded user password to stdout. In a script where stdout lands in CI log storage, a `journalctl` ring buffer, or a Slack paste, the entire kubeconfig walks out. Emit only the specific field you need (e.g. `kubectl config view -o jsonpath='{.current-context}'`) or decrypt once into a temp file and `shred` it.

Disable by adding `ZC1891` to `disabled_katas` in `.zshellcheckrc`.

---

<a id="zc1892"></a>
### ZC1892 — Error on `install -m 4755|6755|2755` — sets setuid/setgid bit at install time

**Severity:** `error`  
**Auto-fix:** `no`

`install -m <mode>` with the setuid (`4xxx`), setgid (`2xxx`), or combined (`6xxx`) octal prefix creates the target with those special bits set, which turns every execution into a privilege-elevation vector. An uninspected binary installed this way — especially from a build script or package post-install — becomes a persistent local-privesc primitive if the binary is writable, has command-injection, or links against attacker-influenced libraries. Drop the setuid/setgid bits from the mode (`install -m 0755`) and grant the narrow capability the program actually needs with `setcap cap_net_bind_service+ep`; audit the remaining setuid binaries with `find / -perm -4000`.

Disable by adding `ZC1892` to `disabled_katas` in `.zshellcheckrc`.

---

<a id="zc1893"></a>
### ZC1893 — Warn on `unsetopt BARE_GLOB_QUAL` — `*(N)` null-glob qualifier stops being special

**Severity:** `warning`  
**Auto-fix:** `no`

`BARE_GLOB_QUAL` is on by default in Zsh — that is what makes the per-glob qualifier syntax (`*(N)` for null-glob, `*(.x)` for executable, `*(Om)` for sort-by-mtime) work. Turning it off reverts to ksh-style parsing where `(...)` inside a glob is a pattern alternation, so `*(N)` stops being a null-glob and turns into "match zero-or-one N" — a completely different pattern. Scripts that relied on `for f in *.log(N)` to cope with empty directories then silently iterate the literal string or fail under NOMATCH. Keep the option on; if you really want ksh-style qualifiers, use `setopt LOCAL_OPTIONS; unsetopt BARE_GLOB_QUAL` inside a function.

Disable by adding `ZC1893` to `disabled_katas` in `.zshellcheckrc`.

---

<a id="zc1894"></a>
### ZC1894 — Error on `conntrack -F` / `--flush` — every tracked connection (including SSH) is reset

**Severity:** `error`  
**Auto-fix:** `no`

`conntrack -F` (alias `--flush`) wipes the netfilter connection-tracking table. Every established TCP flow that depended on conntrack (every stateful-NAT connection, every `-m conntrack --ctstate RELATED,ESTABLISHED` allowance, every MASQUERADE session) loses its entry and the next packet is matched from scratch; most firewall rulesets drop it as "new" and the session dies. Over SSH, that means the shell running the very command drops. Stage the flush behind `at now + 5 minutes` so the session can re-enter the table via a preceding rule, or narrow the scope with `conntrack -D -s <client-IP>` for a specific hung flow.

Disable by adding `ZC1894` to `disabled_katas` in `.zshellcheckrc`.

---

<a id="zc1895"></a>
### ZC1895 — Warn on `setopt NUMERIC_GLOB_SORT` — glob output switches from lexicographic to numeric order

**Severity:** `warning`  
**Auto-fix:** `no`

`NUMERIC_GLOB_SORT` is off by default: `ls *.log` returns filenames in the collation order the filesystem-iteration/sort step produces (lexicographic in the C locale, so `app-1.log`, `app-10.log`, `app-2.log`). Turning it on makes every subsequent glob and array expansion sort numeric runs numerically — the same glob now returns `app-1.log`, `app-2.log`, `app-10.log`. Scripts that tail the "latest" file by taking the last array element, pipelines that expect a specific stable order, and backup rotations built on `*[0-9].tar` silently shuffle. Keep the option off script-wide; request numeric sort per-glob with the `*(n)` qualifier when needed.

Disable by adding `ZC1895` to `disabled_katas` in `.zshellcheckrc`.

---

<a id="zc1896"></a>
### ZC1896 — Error on `docker/podman run -v /proc:…|/sys:…` — bind-mounts host kernel interfaces into container

**Severity:** `error`  
**Auto-fix:** `no`

`docker run -v /proc:/host/proc` (or `-v /sys:…`) bind-mounts the host's procfs / sysfs hierarchy into the container's mount namespace. From inside, the container can read every host process's `environ` (secrets passed via env), every `cmdline`, every `/proc/1/ns/` to open namespace fds for a breakout, and `/sys/fs/cgroup` to modify resource limits that affect host services. `:ro` does not help — `/proc/<pid>/ns/...` handles remain usable. If the container genuinely needs process / kernel visibility, grant the narrowest capability instead (`--cap-add=SYS_PTRACE`) or run the monitoring agent on the host rather than inside an untrusted image.

Disable by adding `ZC1896` to `disabled_katas` in `.zshellcheckrc`.

---

<a id="zc1897"></a>
### ZC1897 — Warn on `setopt SH_GLOB` — Zsh-specific glob patterns (`*(N)`, `<1-10>`, alternation) stop parsing

**Severity:** `warning`  
**Auto-fix:** `no`

`SH_GLOB` is off by default in Zsh. With it off, the shell recognises Zsh's extended patterns: `*(N)` null-glob qualifier, `<1-10>` numeric range globs, `(alt1|alt2)` in-glob alternation, and the whole `(#i)`/`(#c,m)` flag family. Turning the option on forces strict POSIX-sh parsing, so the parser re-interprets `(...)` as command grouping and the null-glob / range idioms raise parse errors. Every kata recommending `*(N)` (see ZC1830, ZC1893) silently breaks, and downstream helpers sourced after the setopt inherit the restricted pattern syntax. Keep the option off; scope inside a function with `setopt LOCAL_OPTIONS; setopt SH_GLOB` if a specific block genuinely needs POSIX patterns.

Disable by adding `ZC1897` to `disabled_katas` in `.zshellcheckrc`.

---

<a id="zc1898"></a>
### ZC1898 — Error on `gpg --export-secret-keys` — private-key material leaks to stdout

**Severity:** `error`  
**Auto-fix:** `no`

`gpg --export-secret-keys KEYID` and `--export-secret-subkeys` write the ASCII-armoured private key to stdout. In a script, that stream usually lands in a file the operator plans to move off-box — and any misstep (wrong `cd`, script-wide stdout captured by CI, tee to a world-readable log, piped into a remote unencrypted channel) permanently leaks the key. Backup the key interactively on an air-gapped machine; if automation is required, write the output to a `umask 077`-protected path and immediately encrypt with a second symmetric passphrase.

Disable by adding `ZC1898` to `disabled_katas` in `.zshellcheckrc`.

---

<a id="zc1899"></a>
### ZC1899 — Error on `mokutil --disable-validation` — turns UEFI Secure Boot off at the shim

**Severity:** `error`  
**Auto-fix:** `no`

`mokutil --disable-validation` queues a request for the shim to stop validating the kernel and modules against the enrolled MOK/PK certificates at next boot — Secure Boot silently becomes advisory. Any unsigned kernel or rootkit module then loads without prompt. Leave Secure Boot validation on; if you must load a custom module, enrol its key with `mokutil --import` and approve via the `MokManager` prompt at reboot.

Disable by adding `ZC1899` to `disabled_katas` in `.zshellcheckrc`.

---

<a id="zc1900"></a>
### ZC1900 — Warn on `curl --location-trusted` — Authorization/cookies forwarded across redirects

**Severity:** `warning`  
**Auto-fix:** `no`

`curl --location-trusted` (alias of `curl -L --location-trusted`) tells curl to replay the `Authorization` header, cookies, and `-u user:pass` credential on every redirect hop, even across hosts. A 302 to an attacker-controlled origin (or a compromised CDN edge) then receives the bearer token verbatim. Drop `--location-trusted`; if cross-origin auth is truly required, scope a short-lived token per destination and verify the final hostname before sending secrets.

Disable by adding `ZC1900` to `disabled_katas` in `.zshellcheckrc`.

---

<a id="zc1901"></a>
### ZC1901 — Warn on `setopt POSIX_BUILTINS` — flips `command`/special-builtin semantics

**Severity:** `warning`  
**Auto-fix:** `no`

`setopt POSIX_BUILTINS` switches Zsh to the POSIX rules for special builtins: assignments before `export`, `readonly`, `eval`, `.`, `trap`, `set`, etc. stay in the caller's scope, and `command builtin` can now resolve shell builtins. Mid-script Zsh code written against native semantics — where those assignments are local — silently leaks state. Leave the option off; scope any POSIX-specific block with `emulate -LR sh` instead of toggling globally.

Disable by adding `ZC1901` to `disabled_katas` in `.zshellcheckrc`.

---

<a id="zc1902"></a>
### ZC1902 — Error on `ln -s /dev/null <logfile>` — silently discards audit or history writes

**Severity:** `error`  
**Auto-fix:** `no`

A symlink from an audit or shell-history path to `/dev/null` turns every subsequent append into a no-op — `/var/log/auth.log`, `wtmp`, `~/.bash_history`, `~/.zsh_history` all stop recording. This is the textbook way to cover tracks on a compromised host and almost never appears in benign automation. If you really need to stop a log, disable the writer (rsyslog rule, `set +o history`) or rotate with `logrotate` — never redirect into `/dev/null`.

Disable by adding `ZC1902` to `disabled_katas` in `.zshellcheckrc`.

---

<a id="zc1903"></a>
### ZC1903 — Error on `tee /etc/sudoers*` — appends a rule that bypasses `visudo -c` validation

**Severity:** `error`  
**Auto-fix:** `no`

`tee /etc/sudoers` or `tee -a /etc/sudoers.d/<name>` is a common shortcut for adding a sudoers rule, but it skips the syntax check that `visudo -c` would perform. A malformed line (missing `ALL`, stray colon, unterminated `Cmnd_Alias`) makes sudo refuse every invocation — you lock yourself out of root recovery. Write the rule to a temporary file, run `visudo -cf /tmp/rule`, and only then `install -m 0440 /tmp/rule /etc/sudoers.d/<name>`.

Disable by adding `ZC1903` to `disabled_katas` in `.zshellcheckrc`.

---

<a id="zc1904"></a>
### ZC1904 — Warn on `setopt KSH_GLOB` — reinterprets `*(pattern)` and breaks Zsh glob qualifiers

**Severity:** `warning`  
**Auto-fix:** `no`

`setopt KSH_GLOB` turns `@(a|b)`, `*(x)`, `+(x)`, `?(x)`, `!(x)` into Korn-shell extended glob operators. The side effect is that `*(N)`, `*(D)`, `*(.)`, and every other Zsh glob qualifier stop working — `*(N)` becomes "zero or more `N` characters", silently shattering null-glob idioms across the script. If you need Korn-style patterns, prefer `setopt EXTENDED_GLOB` and its `(^...)` / `(#...)` forms, which coexist with the qualifier syntax. Otherwise scope the switch inside a function with `setopt LOCAL_OPTIONS KSH_GLOB`.

Disable by adding `ZC1904` to `disabled_katas` in `.zshellcheckrc`.

---

<a id="zc1905"></a>
### ZC1905 — Warn on `ssh -g -L …` — local forward bound on all interfaces, not just loopback

**Severity:** `warning`  
**Auto-fix:** `no`

`ssh -g` flips the default for `-L` (local forward) and `-D` (dynamic SOCKS) from `127.0.0.1:port` to `0.0.0.0:port`. Any host on the same LAN/VPN/WiFi segment can then use the tunnel without authenticating to the SSH session. Drop `-g`, pin the bind explicitly with `-L bind_address:port:target:port`, or use a firewall rule — never leave a forwarded port open to the network segment.

Disable by adding `ZC1905` to `disabled_katas` in `.zshellcheckrc`.

---

<a id="zc1906"></a>
### ZC1906 — Warn on `setopt POSIX_CD` — changes when `cd` / `pushd` consult `CDPATH`

**Severity:** `warning`  
**Auto-fix:** `no`

`setopt POSIX_CD` makes `cd`, `chdir`, and `pushd` skip `CDPATH` for any argument that starts with `/`, `.`, or `..`. Zsh's default — consulting `CDPATH` for anything that does not start with `/` — was exactly what made `cd foo` resolve the "project" dir via `CDPATH` even when a local `./foo` existed. Flipping the option globally makes scripts that relied on the Zsh behaviour silently enter different directories. Keep the option off; if POSIX parity is needed, wrap a single function with `emulate -LR sh`.

Disable by adding `ZC1906` to `disabled_katas` in `.zshellcheckrc`.

---

<a id="zc1907"></a>
### ZC1907 — Warn on `sysctl -w fs.protected_*=0` / `fs.suid_dumpable=2` — disables /tmp-race safeguards

**Severity:** `warning`  
**Auto-fix:** `no`

Linux ships `fs.protected_symlinks`, `fs.protected_hardlinks`, `fs.protected_fifos`, and `fs.protected_regular` enabled to stop classic `/tmp`-race escalation (dangling-symlink, hardlink-pivot, FIFO-open-owner). Setting any of them to `0`, or raising `fs.suid_dumpable` above `0`, hands unprivileged local users back the primitives. Keep the defaults; if a legacy tool genuinely needs them off, scope the change inside a namespace rather than flipping the host knob.

Disable by adding `ZC1907` to `disabled_katas` in `.zshellcheckrc`.

---

<a id="zc1908"></a>
### ZC1908 — Warn on `setopt MAGIC_EQUAL_SUBST` — enables tilde/param expansion on `key=value` args

**Severity:** `warning`  
**Auto-fix:** `no`

`MAGIC_EQUAL_SUBST` tells Zsh that every unquoted argument of the form `identifier=value` gets file expansion on the right-hand side, as if it were a parameter assignment. Under the default (option off), `rsync host:dst=~/backup` keeps the literal `~` — under the option on, the `~` expands to your home. Flipping the option globally makes a whole class of literal CLI arguments silently change meaning. Leave the option off; if a specific assignment truly needs expansion, wrap it in quotes or use a temporary variable.

Disable by adding `ZC1908` to `disabled_katas` in `.zshellcheckrc`.

---

<a id="zc1909"></a>
### ZC1909 — Warn on `kexec -l` / `-e` — jumps to an alternate kernel, bypasses bootloader and Secure Boot

**Severity:** `warning`  
**Auto-fix:** `no`

`kexec -l /path/to/vmlinuz …` stages a second kernel image, and `kexec -e` (or `kexec -f`) then transfers control to it without going through the firmware, GRUB, or shim. On a Secure-Boot system the staged kernel is never verified against the enrolled MOK/PK — an attacker who lands a root exec can boot a hostile kernel while leaving /boot untouched. Reserve `kexec` for the live-patching / crash-dump workflow it was designed for, gate the call behind `sudo` + audit, and prefer `systemctl kexec` or a normal reboot when possible.

Disable by adding `ZC1909` to `disabled_katas` in `.zshellcheckrc`.

---

<a id="zc1910"></a>
### ZC1910 — Warn on `setopt GLOB_STAR_SHORT` — makes bare `**` recurse instead of matching literal

**Severity:** `warning`  
**Auto-fix:** `no`

`GLOB_STAR_SHORT` teaches Zsh to expand bare `**` (not followed by `/`) as if it were `**/*` — suddenly `rm **` wipes every file under the current directory instead of erroring or matching the two-star literal. Scripts that pass `**` as a literal argument to `grep`, `sed`, or a logger call silently turn into deep directory recursions. Keep the option off; when you really need recursive globs, spell `**/*` explicitly so reviewers can see the intent.

Disable by adding `ZC1910` to `disabled_katas` in `.zshellcheckrc`.

---

<a id="zc1911"></a>
### ZC1911 — Warn on `umount -l` / `--lazy` — detach now, leaves open fds pointing at a ghost mount

**Severity:** `warning`  
**Auto-fix:** `no`

`umount -l` (lazy unmount) detaches the filesystem from the directory tree immediately but defers the real cleanup until every open file descriptor on it is closed. Any process still holding an fd keeps reading/writing into a mount that `mount | grep` no longer lists — cron jobs drop logs into a phantom directory, a re-mount of the same path stacks invisibly, and `lsof`/`fuser` often miss the stale handles. Find and stop the holder (`lsof`/`fuser`/`systemd-cgls`) first, then do a normal `umount`; reserve `-l` for break-glass recovery, not scripts.

Disable by adding `ZC1911` to `disabled_katas` in `.zshellcheckrc`.

---

<a id="zc1912"></a>
### ZC1912 — Warn on `dhclient -r` / `dhclient -x` / `dhcpcd -k` — drops the lease and breaks network

**Severity:** `warning`  
**Auto-fix:** `no`

`dhclient -r` releases the current DHCP lease (sending a DHCPRELEASE), `dhclient -x` terminates the daemon without releasing, and `dhcpcd -k` does the equivalent for dhcpcd. On a remote host the very next thing that happens is the SSH session drops, and in a VPC any automation waiting for a reply never sees one. Stage the release together with a re-acquire (`dhclient -1 $iface` or `nmcli device reapply $iface`) or schedule it via `systemd-run --on-active=` so the operator is not cut off mid-session.

Disable by adding `ZC1912` to `disabled_katas` in `.zshellcheckrc`.

---

<a id="zc1913"></a>
### ZC1913 — Warn on `setopt ALIAS_FUNC_DEF` — re-enables defining functions with aliased names

**Severity:** `warning`  
**Auto-fix:** `no`

Zsh's default refuses the syntax `ls () { … }` when `ls` is aliased — because the alias expands at definition time and the function the author meant to write never actually exists. `setopt ALIAS_FUNC_DEF` disables that guardrail: the alias is suppressed during definition, and the function silently shadows the alias afterwards. The combination is almost always a bug — one alias in a sourced rc file quietly replaces the function. Keep the option off and write `function \ls () { … }` (quoted) if you really need to override an aliased name.

Disable by adding `ZC1913` to `disabled_katas` in `.zshellcheckrc`.

---

<a id="zc1914"></a>
### ZC1914 — Warn on `curl --doh-url …` / `--dns-servers …` — overrides system resolver per-request

**Severity:** `warning`  
**Auto-fix:** `no`

`curl --doh-url https://doh.example/dns-query` routes the lookup through a caller-specified DNS-over-HTTPS endpoint; `curl --dns-servers 1.1.1.1,8.8.8.8` forces classic UDP to the listed servers. Both detour around the host's resolver chain — `/etc/hosts`, `systemd-resolved`, `nsswitch`, split-horizon DNS — so the request lands at an IP the operator did not vet. In production scripts that is usually a stray debug line left in; drop the flag or gate it behind an explicit `--doh-insecure` + `--resolve` pinning audit so reviewers can see the intent.

Disable by adding `ZC1914` to `disabled_katas` in `.zshellcheckrc`.

---

<a id="zc1915"></a>
### ZC1915 — Error on `mdadm --zero-superblock` / `--stop` — drops RAID metadata or live array

**Severity:** `error`  
**Auto-fix:** `no`

`mdadm --zero-superblock $DEV` wipes the MD superblock from a member — the array forgets the device exists and a subsequent `--create` with the wrong layout permanently scrambles the data. `mdadm --stop $MD` (or `-S`) halts a live array from underneath whatever is mounted on it; if root or `/boot` lives there the host panics on the next fsync. Run `mdadm --examine` first, snapshot the superblock with `mdadm --detail --export`, and keep both calls behind a runbook rather than an automated script.

Disable by adding `ZC1915` to `disabled_katas` in `.zshellcheckrc`.

---

<a id="zc1916"></a>
### ZC1916 — Warn on `setopt NULL_GLOB` — every unmatched glob silently expands to nothing

**Severity:** `warning`  
**Auto-fix:** `no`

`setopt NULL_GLOB` removes the Zsh default behaviour of erroring out when a pattern matches nothing. Every later glob becomes silently empty instead — `cp *.log /dest` when no `.log` files exist turns into `cp /dest` (wrong target), `rm *.tmp` into `rm` (argv too short), and `for f in *.json` into a no-op. Reach for the per-glob `*(N)` qualifier when you want a single pattern to tolerate a zero match, or scope the switch with `setopt LOCAL_OPTIONS NULL_GLOB` inside the one function that needs it.

Disable by adding `ZC1916` to `disabled_katas` in `.zshellcheckrc`.

---

<a id="zc1917"></a>
### ZC1917 — Info on `iw dev $IF scan` / `iwlist $IF scan` — active WiFi scan from a script

**Severity:** `info`  
**Auto-fix:** `no`

`iw dev wlan0 scan` (and the older `iwlist wlan0 scan`) performs an active probe-request sweep across every supported channel. It requires `CAP_NET_ADMIN`, briefly interrupts the current association, and announces the host's presence to every nearby access point — logs on the other side will show one MAC asking about every SSID. Use the cached `iw dev $IF link` / `iwctl station $IF show` for passive lookups, and reserve `scan` for diagnostic sessions with console approval rather than background scripts.

Disable by adding `ZC1917` to `disabled_katas` in `.zshellcheckrc`.

---

<a id="zc1918"></a>
### ZC1918 — Warn on `setopt HIST_SUBST_PATTERN` — `!:s/old/new/` silently switches to pattern matching

**Severity:** `warning`  
**Auto-fix:** `no`

`HIST_SUBST_PATTERN` makes the `:s` and `:&` history modifiers, as well as the identically-named parameter-expansion modifier `${foo:s/pat/rep/}`, match on patterns rather than literal strings. Text that looked safe as a constant (`#` comments, `^` anchors, `?`, `*`) suddenly gets interpreted as glob metacharacters, and replacements that always returned the original string now edit it in surprising ways. Keep the option off and use `${var//pat/rep}` explicitly when you do want glob substitution — that form declares the intent at the call site instead of via a shell-wide flag.

Disable by adding `ZC1918` to `disabled_katas` in `.zshellcheckrc`.

---

<a id="zc1919"></a>
### ZC1919 — Warn on `ss -K` / `ss --kill` — terminates every socket that matches the filter

**Severity:** `warning`  
**Auto-fix:** `no`

`ss -K` issues `SOCK_DESTROY` to every socket matching the filter (requires `CAP_NET_ADMIN`). With a broad filter — `ss -K state established`, `ss -K dport 22` — the command happily terminates the SSH session that is running it, along with every backend keep-alive that happens to match. Spell the filter tightly (`ss -K dst 10.0.0.5 dport 5432 state close-wait`), test it first without `-K` to confirm only the target sockets appear, and wrap the call in a review step rather than a scheduled job.

Disable by adding `ZC1919` to `disabled_katas` in `.zshellcheckrc`.

---

<a id="zc1920"></a>
### ZC1920 — Warn on `setopt VERBOSE` — every executed command is echoed to stderr

**Severity:** `warning`  
**Auto-fix:** `no`

`setopt VERBOSE` is Zsh's name for the POSIX `set -v` flag: the shell prints each command line to stderr immediately after reading it. In a script that processes secrets the stderr stream then carries every command that mentions them, including `mysql -pSECRET`, `curl -u user:pass`, `export DB_PASS=…`. Unlike `set -x` (which already has dedicated detectors) the `VERBOSE` flag is easy to leave on by accident because the output looks like normal command echo. Remove the call and rely on `printf` / a proper logger; if a debug trace is required, scope it in a function with `setopt LOCAL_OPTIONS VERBOSE` then `unsetopt VERBOSE`.

Disable by adding `ZC1920` to `disabled_katas` in `.zshellcheckrc`.

---

<a id="zc1921"></a>
### ZC1921 — Warn on `systemctl kill -s KILL` / `--signal=SIGKILL` — skips `ExecStop=`, leaks resources

**Severity:** `warning`  
**Auto-fix:** `no`

`systemctl kill UNIT -s KILL` (and `--signal=9` / `SIGKILL`) bypasses the unit's `ExecStop=` sequence and the `TimeoutStopSec=` budget. Any lockfile, socket, or shared-memory segment the service was supposed to unlink survives; the next restart often fails with "address already in use" or a corrupt journal. Default to `systemctl stop UNIT` (or `restart`) and let the stop sequence run. Reserve `-s KILL` for a last-resort recovery path with a runbook attached.

Disable by adding `ZC1921` to `disabled_katas` in `.zshellcheckrc`.

---

<a id="zc1922"></a>
### ZC1922 — Error on `rpm --import http://…` / `rpmkeys --import ftp://…` — plaintext GPG key fetch

**Severity:** `error`  
**Auto-fix:** `no`

`rpm --import` (and `rpmkeys --import`) add the supplied ASCII-armoured key to the system RPM trust store. When the source is a plain `http://` / `ftp://` URL an on-path attacker swaps the key, and every subsequent package they sign installs cleanly. Serve keys over HTTPS from a TLS-authenticated origin, pin the key's SHA-256 before import, or stage an offline copy verified out of band (`gpg --verify` against a known-good fingerprint).

Disable by adding `ZC1922` to `disabled_katas` in `.zshellcheckrc`.

---

<a id="zc1923"></a>
### ZC1923 — Warn on `setopt PRINT_EXIT_VALUE` — every non-zero exit leaks a status line to stderr

**Severity:** `warning`  
**Auto-fix:** `no`

`PRINT_EXIT_VALUE` makes Zsh emit `zsh: exit N` on stderr after every foreground command that returns a non-zero status. In a script the stream is typically captured by a supervisor or shipped to a log aggregator, and the extra line reveals which tool returned what — including grep / test / curl probes that were supposed to stay silent. Worse, tools that parse stderr for diagnostics (`git`, `ssh`, `rsync`) now see interleaved shell chatter. Remove the `setopt` call; if you actually want a per-command post-mortem, rely on `precmd`/`preexec` hooks or an explicit `|| printf …`.

Disable by adding `ZC1923` to `disabled_katas` in `.zshellcheckrc`.

---

<a id="zc1924"></a>
### ZC1924 — Warn on `virt-cat` / `virt-copy-out` / `guestfish` / `guestmount` — reads guest disk from host

**Severity:** `warning`  
**Auto-fix:** `no`

libguestfs tools (`virt-cat`, `virt-copy-out`, `virt-tar-out`, `virt-edit`, `virt-customize`, `guestfish`, `guestmount`) open a VM's disk image directly from the hypervisor and read or mutate its contents without going through the guest OS. That bypasses every in-guest permission, audit, and LUKS keyslot the VM was using, and — if the VM is live — risks filesystem corruption because two writers are now mounted on the same image. Snapshot the disk first, work on the clone, and prefer in-guest `ssh`/`scp`/`ansible` for anything that does not need out-of-band recovery.

Disable by adding `ZC1924` to `disabled_katas` in `.zshellcheckrc`.

---

<a id="zc1925"></a>
### ZC1925 — Warn on `unsetopt EQUALS` — disables `=cmd` path expansion and tilde-after-colon

**Severity:** `warning`  
**Auto-fix:** `no`

Zsh's `EQUALS` option (on by default) is what makes `=python`, `=ls`, and `=vim` expand to the absolute path of the command via `$PATH` lookup. It also drives the `PATH=~/bin:$PATH` tilde-after-colon expansion. `unsetopt EQUALS` turns both off: `=cmd` becomes a literal argument (breaking any idiom that relies on the short-path), and `PATH=~/bin:$PATH` stops expanding the tilde inside the colon-separated list. Keep the option on; if one function needs literal `=` arguments, scope via `setopt LOCAL_OPTIONS; unsetopt EQUALS` inside it.

Disable by adding `ZC1925` to `disabled_katas` in `.zshellcheckrc`.

---

<a id="zc1926"></a>
### ZC1926 — Warn on `telinit 0/1/6` / `init 0/1/6` — SysV runlevel change halts, reboots, or isolates the host

**Severity:** `warning`  
**Auto-fix:** `no`

`init 0`, `init 6`, `init 1`, and their `telinit` aliases ask systemd (or SysV) to switch runlevel: `0` → `poweroff.target`, `6` → `reboot.target`, `1`/`S` → `rescue.target`. From a script the side effect is a remote SSH disconnect, an immediate service teardown for every other session on the host, and — in the `1`/`S` case — dropping to single-user mode without a console to recover. Use `systemctl poweroff`/`reboot`/`rescue` (which are clearer in reviews) or schedule via `shutdown -h +N` so the operator has a cancel window.

Disable by adding `ZC1926` to `disabled_katas` in `.zshellcheckrc`.

---

<a id="zc1927"></a>
### ZC1927 — Error on `xfreerdp /p:SECRET` / `rdesktop -p SECRET` — RDP password visible in argv

**Severity:** `error`  
**Auto-fix:** `no`

`xfreerdp /p:<password>` and `rdesktop -p <password>` (plus the `-p -` stdin form when followed by an argv password) put the Windows credential into `ps`, `/proc/PID/cmdline`, shell history, and every `ps aux` captured by monitoring. Use `xfreerdp /from-stdin` + a piped credential, `freerdp-shadow-cli /sec:nla` with a cached credential, or drop the password into a protected `.rdp` file passed via `/load-config-file`. Never inline the password on the command line.

Disable by adding `ZC1927` to `disabled_katas` in `.zshellcheckrc`.

---

<a id="zc1928"></a>
### ZC1928 — Warn on `setopt SHARE_HISTORY` — every session writes its history into every sibling session

**Severity:** `warning`  
**Auto-fix:** `no`

`SHARE_HISTORY` flushes each command to `$HISTFILE` immediately and tells all other running zsh sessions to re-read the file. A secret typed in a one-off "private" terminal — `ssh user@host "$PASS"`, `aws sts ... --output text`, `git push https://user:token@…` — shows up in every other terminal's `fc -l` list seconds later. Prefer `setopt INC_APPEND_HISTORY` (append-only, per-session isolation) and `setopt HIST_IGNORE_SPACE` so a leading space keeps the line out of history altogether.

Disable by adding `ZC1928` to `disabled_katas` in `.zshellcheckrc`.

---

<a id="zc1929"></a>
### ZC1929 — Warn on `cpio -i` / `--extract` without `--no-absolute-filenames` — archive writes outside CWD

**Severity:** `warning`  
**Auto-fix:** `no`

`cpio -i` (and `--extract`) is the default copy-in mode: it materialises every path stored in the archive verbatim. Paths starting with `/` land where the archive told them to, and relative paths containing `..` slip out of the extraction directory entirely — so a rogue initramfs or firmware bundle can drop files into `/etc/cron.d/`, `/usr/lib/systemd/system/`, or the operator's `~/.ssh/authorized_keys`. Always pass `--no-absolute-filenames` and extract into a fresh scratch directory reviewed before `mv`.

Disable by adding `ZC1929` to `disabled_katas` in `.zshellcheckrc`.

---

<a id="zc1930"></a>
### ZC1930 — Warn on `unsetopt HASH_CMDS` — every command invocation re-walks `$PATH`

**Severity:** `warning`  
**Auto-fix:** `no`

`HASH_CMDS` (on by default) caches the resolved absolute path of every command after its first successful lookup. `unsetopt HASH_CMDS` disables the cache, so each invocation re-walks every `$PATH` entry and re-runs `stat()` on every candidate. On a slow filesystem (NFS home, encrypted volume, large `$PATH`) this adds tens to hundreds of milliseconds per command and can double the runtime of a long pipeline. Keep the option on; if you are changing a binary and want the cache invalidated, `rehash` (one-shot) or `hash -r` is the scoped fix.

Disable by adding `ZC1930` to `disabled_katas` in `.zshellcheckrc`.

---

<a id="zc1931"></a>
### ZC1931 — Warn on `ip netns delete $NS` / `netns del` — drops the whole network namespace

**Severity:** `warning`  
**Auto-fix:** `no`

`ip netns delete NAME` / `ip netns del NAME` unmounts the namespace and tears down every interface, veth pair, VXLAN, and WireGuard peer living inside. Processes still attached lose their network abruptly — container health checks fail, BGP sessions drop, and any other process using `ip netns exec NAME …` errors out with "No such file or directory". Stop the workloads first (`systemctl stop`, `pkill -SIGTERM -n $NS`), confirm `ip -n $NS link` is empty, then `delete` deliberately — or leave the namespace alone if it is managed by Docker/containerd/systemd-nspawn.

Disable by adding `ZC1931` to `disabled_katas` in `.zshellcheckrc`.

---

<a id="zc1932"></a>
### ZC1932 — Warn on `unsetopt GLOBAL_EXPORT` — `typeset -x` in a function stops leaking to outer scope

**Severity:** `warning`  
**Auto-fix:** `no`

`GLOBAL_EXPORT` (on by default) makes `typeset -x VAR=val` inside a function not only export `VAR` but also promote it to the outer scope, so callers and subsequent functions see the same value. Turning it off changes the meaning of every such assignment across the script: exports become function-local and vanish the moment the function returns. Scripts that rely on a helper to set up `PATH`, `VIRTUAL_ENV`, or `AWS_*` variables suddenly run commands under the old environment. Keep the option on; if you want a temporary export, scope it with a subshell instead of a shell-wide flip.

Disable by adding `ZC1932` to `disabled_katas` in `.zshellcheckrc`.

---

<a id="zc1933"></a>
### ZC1933 — Error on `ipvsadm -C` / `--clear` — wipes every IPVS virtual service, drops load balancer

**Severity:** `error`  
**Auto-fix:** `no`

`ipvsadm -C` (and the long form `--clear`) removes every virtual service, real server, and connection entry from the in-kernel IPVS table. Traffic that was being load-balanced to a backend farm now falls through to the host's local listen sockets (or drops), active keepalived/`ldirectord` states invert, and clients see 5xx until an operator replays the config. Save the current table first (`ipvsadm --save -n > /run/ipvs.bak`), drain specific services with `ipvsadm -D`, and keep `--clear` in break-glass-only runbooks.

Disable by adding `ZC1933` to `disabled_katas` in `.zshellcheckrc`.

---

<a id="zc1934"></a>
### ZC1934 — Warn on `setopt AUTO_NAME_DIRS` — any absolute-path parameter becomes a `~name` alias

**Severity:** `warning`  
**Auto-fix:** `no`

`AUTO_NAME_DIRS` (off by default) auto-registers any parameter whose value is an absolute directory path as a named directory — so `foo=/srv/data` immediately makes `~foo` resolve to `/srv/data` in later expansions and in `%~` prompt sequences. The option silently changes the meaning of `ls ~foo` across the script and surfaces directory names in `%~` prompts that the user never opted into. Keep the option off and call `hash -d name=/path` explicitly when a named directory is actually wanted.

Disable by adding `ZC1934` to `disabled_katas` in `.zshellcheckrc`.

---

<a id="zc1935"></a>
### ZC1935 — Warn on `apt autoremove --purge` / `dnf autoremove` — deletes auto-installed deps and their config

**Severity:** `warning`  
**Auto-fix:** `no`

`apt autoremove --purge` (and `apt-get autoremove --purge`, `dnf autoremove`, `zypper rm --clean-deps`) remove every package the resolver thinks is no longer required, plus — with `--purge` — their `/etc` config and data dirs. In CI this quietly uproots packages someone else installed manually but never `apt-mark manual`-ed, and `--purge` makes the removal irreversible. Run a plain `apt autoremove --dry-run` in review, mark the keepers with `apt-mark manual`, and drop `--purge` from unattended jobs.

Disable by adding `ZC1935` to `disabled_katas` in `.zshellcheckrc`.

---

<a id="zc1936"></a>
### ZC1936 — Warn on `setopt POSIX_ALIASES` — aliases on reserved words (`if`, `for`, …) stop expanding

**Severity:** `warning`  
**Auto-fix:** `no`

Zsh by default lets `alias if='…'`, `alias function='…'`, etc. expand when the reserved word appears in command position — the feature that makes oh-my-zsh plugins able to hook `if` into their `preexec` chain. `setopt POSIX_ALIASES` narrows alias expansion to plain identifiers, so any library that aliased a reserved word silently stops being picked up. Keep the option off for interactive Zsh; if you need POSIX parity for a specific block, wrap it with `emulate -LR sh` instead of flipping the flag script-wide.

Disable by adding `ZC1936` to `disabled_katas` in `.zshellcheckrc`.

---

<a id="zc1937"></a>
### ZC1937 — Warn on `tmux kill-server` / `tmux kill-session` — tears down every detached process inside

**Severity:** `warning`  
**Auto-fix:** `no`

`tmux kill-server` terminates the whole tmux daemon, `tmux kill-session -t NAME` drops one named session, and `screen -X quit` does the screen equivalent. Anything the operator parked inside — a long-running build, a `tail -F` on production logs, a held `sudo` token, a port-forward — dies with the session, and the detached processes get `SIGHUP`'d with no cleanup. Use `tmux kill-window -t …` for surgical removal, send `SIGTERM` to the specific backend, or rely on `systemd-run --scope` for workloads that should survive terminal churn.

Disable by adding `ZC1937` to `disabled_katas` in `.zshellcheckrc`.

---

<a id="zc1938"></a>
### ZC1938 — Warn on `setopt POSIX_JOBS` — flips job-control semantics and `%n` scope

**Severity:** `warning`  
**Auto-fix:** `no`

`POSIX_JOBS` makes Zsh's job-control spec follow POSIX: `%1` / `%n` refer only to jobs of the current shell (forked subshells get their own job table), `fg`/`bg` no longer accept a job ID from an outer shell, and `disown` on a subshell's job is a no-op. Scripts that launched a background job in the parent and then `wait %1`-ed from a `( subshell )` suddenly fail with "no such job". Leave the option off in Zsh; if POSIX job semantics are required, scope them via `emulate -LR sh` inside the single function that needs them.

Disable by adding `ZC1938` to `disabled_katas` in `.zshellcheckrc`.

---

<a id="zc1939"></a>
### ZC1939 — Error on `reboot -f` / `halt -f` / `poweroff -f` — skips shutdown sequence, no graceful service stop

**Severity:** `error`  
**Auto-fix:** `no`

`reboot -f`, `halt -f`, and `poweroff -f` short-circuit the systemd shutdown graph — no `ExecStop=`, no `DefaultDependencies=`, no filesystem sync, no Before/After ordering. The kernel's `reboot(2)` fires immediately and every dirty buffer that was not yet flushed is lost. Journal writes stop mid-line, databases on the host replay from the last checkpoint, and anything that needed a clean unmount (LUKS, NFS, cephfs) logs a dirty state. Use plain `systemctl reboot` / `shutdown -r +N`, and reserve `-f` for recovery when the normal path is already wedged.

Disable by adding `ZC1939` to `disabled_katas` in `.zshellcheckrc`.

---

<a id="zc1940"></a>
### ZC1940 — Warn on `setopt POSIX_ARGZERO` — `$0` no longer changes to the function name inside functions

**Severity:** `warning`  
**Auto-fix:** `no`

Zsh's default behaviour (option off) assigns `$0` to the name of the currently-running function, so a helper like `log() { printf '%s\n' "$0: $*"; }` prints `log: …`. `setopt POSIX_ARGZERO` keeps `$0` pointing at the outer script name (or the interpreter when sourced) — the logger instead prints the script path for every message and call-site context is lost. Every `case $0` dispatch inside an auto-loaded function also stops working. Leave the option off; if you need POSIX `$0`, scope it in a function with `emulate -LR sh`.

Disable by adding `ZC1940` to `disabled_katas` in `.zshellcheckrc`.

---

<a id="zc1941"></a>
### ZC1941 — Error on `restic init --insecure-no-password` — creates an unencrypted backup repository

**Severity:** `error`  
**Auto-fix:** `no`

`restic init --insecure-no-password` creates a repo whose data chunks are reachable without a key. Every later `backup` and `restore` round-trips plaintext blocks to the storage backend, so any operator with read access to the bucket / NFS share / SFTP directory can assemble the backed-up filesystem — including shell history, SSH keys, and database dumps. Pass a real passphrase via `--password-file` (mode `0400`, readable only by the backup user) or `--password-command`, and never use the `--insecure-*` family outside a local test repo.

Disable by adding `ZC1941` to `disabled_katas` in `.zshellcheckrc`.

---

<a id="zc1942"></a>
### ZC1942 — Warn on `setopt CLOBBER_EMPTY` — `>file` still overwrites zero-length files under `NO_CLOBBER`

**Severity:** `warning`  
**Auto-fix:** `no`

`setopt CLOBBER_EMPTY` relaxes `NO_CLOBBER`: a bare `>file` redirect still succeeds when the target is zero bytes. Scripts that rely on `setopt NO_CLOBBER` as a guard against accidental overwrite lose their safety net for every freshly-`touch`ed lock file, sentinel, or `install -D`-created placeholder — the next stray `>sentinel` quietly overwrites it. Keep the option off; use `>|file` explicitly when you do want to bypass the `NO_CLOBBER` guard for a specific write.

Disable by adding `ZC1942` to `disabled_katas` in `.zshellcheckrc`.

---

<a id="zc1943"></a>
### ZC1943 — Warn on `systemd-nspawn -b` / `--boot` — runs a full init inside a possibly untrusted rootfs

**Severity:** `warning`  
**Auto-fix:** `no`

`systemd-nspawn -b -D $ROOT` (and `--boot -D $ROOT`) launches the rootfs's `/sbin/init` inside a minimally-isolated namespace — by default the container inherits `CAP_AUDIT_CONTROL`, `CAP_NET_ADMIN`, and read-write access to the host's `/dev` nodes that match the container's cgroup. If `$ROOT` is an operator-supplied tarball, any init script it ships runs first and can probe the host. Use `-U` for user-namespace isolation, drop capabilities with `--capability=`, pair with `--private-network`, and prefer `machinectl start` on a reviewed image instead of ad-hoc boots.

Disable by adding `ZC1943` to `disabled_katas` in `.zshellcheckrc`.

---

<a id="zc1944"></a>
### ZC1944 — Warn on `setopt IGNORE_EOF` — Ctrl-D no longer exits the shell, masking runaway pipelines

**Severity:** `warning`  
**Auto-fix:** `no`

`IGNORE_EOF` tells the interactive shell to treat an end-of-file on stdin as if it were nothing, so `Ctrl-D` stops terminating a login. In an unattended `zsh -i -c` launch, or a sourced rc, this keeps a subshell alive that was supposed to wind down when the controlling terminal went away — sudo sessions, SSH tunnels, port-forwards, and build supervisors then linger long after the parent left. Keep the option off; if a stale-tty guard is truly wanted, set `TMOUT=NN` for a timed exit instead.

Disable by adding `ZC1944` to `disabled_katas` in `.zshellcheckrc`.

---

<a id="zc1945"></a>
### ZC1945 — Warn on `bpftrace -e` / `bpftool prog load` — loads in-kernel eBPF from a script

**Severity:** `warning`  
**Auto-fix:** `no`

`bpftrace -e '…'` compiles an inline script into an eBPF program and attaches to kprobes, tracepoints, or uprobes; `bpftool prog load FILE pinned /sys/fs/bpf/…` installs a pre-built program. Both require `CAP_BPF`/`CAP_SYS_ADMIN` and can read arbitrary kernel/userland memory — every command a sibling process runs, every syscall argument, every TCP payload. Pin the loaded program to a directory the operator owns, gate invocation behind a runbook, and prefer a short-lived `bpftrace -c CMD` window over long-running traces left on the host.

Disable by adding `ZC1945` to `disabled_katas` in `.zshellcheckrc`.

---

<a id="zc1946"></a>
### ZC1946 — Warn on `unsetopt HUP` — background jobs keep running after shell exit

**Severity:** `warning`  
**Auto-fix:** `no`

Zsh's `HUP` option (on by default) sends `SIGHUP` to each running child job when the shell exits, letting them wind down cleanly. `unsetopt HUP` / `setopt NO_HUP` disables that, so long pipelines, `sleep` loops, and user-spawned daemons live on — `ps aux` accumulates orphaned workers across logouts and resource consumption creeps up. If a specific job really needs to outlive the shell, use `disown` or `systemd-run --scope` on that one invocation; leave `HUP` on globally.

Disable by adding `ZC1946` to `disabled_katas` in `.zshellcheckrc`.

---

<a id="zc1947"></a>
### ZC1947 — Error on `ip xfrm state flush` / `ip xfrm policy flush` — tears down every IPsec SA and policy

**Severity:** `error`  
**Auto-fix:** `no`

`ip xfrm state flush` removes every IPsec Security Association; `ip xfrm policy flush` removes every policy that would have driven them. Strongswan, libreswan, FRR, and WireGuard-over-xfrm all lose their tunnels instantly — site-to-site VPNs drop, kernel packet paths stop encrypting, and peers renegotiate from scratch (with traffic leaking in plaintext during the gap on misconfigured hosts). Use `ip xfrm state deleteall src $A dst $B` to scope the change to a single tunnel, and pair flushes with a maintenance window.

Disable by adding `ZC1947` to `disabled_katas` in `.zshellcheckrc`.

---

<a id="zc1948"></a>
### ZC1948 — Error on `ipmitool -P PASS` / `-E` — BMC password visible in argv

**Severity:** `error`  
**Auto-fix:** `no`

`ipmitool -H <bmc> -U admin -P <password>` puts the BMC credential into `ps`, `/proc/PID/cmdline`, and every process-dump crash file. The BMC is a root-equivalent out-of-band controller (power, console, firmware update), so that password is one of the most sensitive tokens on the host. Use `-f <password_file>` (mode `0400`, owned by the automation user) or set `IPMI_PASSWORD` and pass `-E` — `ipmitool` reads the env var but never echoes it.

Disable by adding `ZC1948` to `disabled_katas` in `.zshellcheckrc`.

---

<a id="zc1949"></a>
### ZC1949 — Error on `rmmod -f` / `rmmod --force` — bypasses refcount, can panic the kernel

**Severity:** `error`  
**Auto-fix:** `no`

`rmmod -f` asks the kernel to tear down a module even if its reference count is non-zero. Any live `open("/dev/…")`, mounted filesystem, or in-flight network device driven by that module becomes a dangling pointer — the kernel oopses or outright panics as soon as the next callback fires. The feature is compiled out on most distros (`CONFIG_MODULE_FORCE_UNLOAD=n`), but when present it is strictly a break-glass recovery tool. Stop the holders first (`lsof /dev/FOO`, `umount`, `ip link set dev … down`), then use plain `rmmod` or `modprobe -r`.

Disable by adding `ZC1949` to `disabled_katas` in `.zshellcheckrc`.

---

<a id="zc1950"></a>
### ZC1950 — Error on `tune2fs -O ^has_journal` / `-m 0` — removes journal or root reserve

**Severity:** `error`  
**Auto-fix:** `no`

`tune2fs -O ^has_journal $DEV` strips the ext3/4 journal from the filesystem. Crash recovery drops from "replay the journal" to "scan the whole block device with `fsck -y`", which frequently truncates partially-written files. `tune2fs -m 0 $DEV` takes the reserved-for-root space down to zero; when the filesystem fills up there is no headroom for `journald`, `apt`, or even a root shell to clean up — recovery needs rescue media. Keep the journal on and leave `-m` at the distro default (5% is overkill on large disks, but `-m 1` is still safe).

Disable by adding `ZC1950` to `disabled_katas` in `.zshellcheckrc`.

---

<a id="zc1951"></a>
### ZC1951 — Error on `ceph osd pool delete … --yes-i-really-really-mean-it` — automates Ceph's double-safety phrase

**Severity:** `error`  
**Auto-fix:** `no`

Ceph intentionally requires both the pool name twice and the flag `--yes-i-really-really-mean-it` before it will delete a pool, so a typo during a live operation cannot drop production data. Baking the phrase into a script defeats the friction — a rebase of the wrong variable, a typo in the pool name, or a stale `for pool in $(…)` loop then silently deletes real pools. Remove the flag from scripts. Do the deletion interactively, or wrap it in a runbook that spells out the pool name in the commit message the operator acknowledges.

Disable by adding `ZC1951` to `disabled_katas` in `.zshellcheckrc`.

---

<a id="zc1952"></a>
### ZC1952 — Error on `zfs set sync=disabled` — `fsync()` becomes a no-op, crash loses unflushed writes

**Severity:** `error`  
**Auto-fix:** `no`

`zfs set sync=disabled POOL/DATASET` turns `fsync()`, `O_SYNC`, and `O_DSYNC` into no-ops on that dataset. PostgreSQL, MariaDB, etcd, and every application that relies on fsync for durability will report success for writes that are still in the ARC, so a panic or power cut loses minutes of committed transactions. The flag is a benchmarking knob, not a production setting. Leave sync at `standard` and, if latency is the concern, add a `log` vdev (SLOG) or tune `zfs_txg_timeout` instead.

Disable by adding `ZC1952` to `disabled_katas` in `.zshellcheckrc`.

---

<a id="zc1953"></a>
### ZC1953 — Warn on `mount --make-shared` / `--make-rshared` — flips propagation, container-escape vector

**Severity:** `warning`  
**Auto-fix:** `no`

`mount --make-shared /path` (and the recursive `--make-rshared`) turns the mount point into a peer in a shared-subtree group. Any later bind-mount that lands inside it propagates to every other peer, including containers and other namespaces. Combined with `CAP_SYS_ADMIN` inside a pod, that is one of the classic container-escape stepping stones — a hostile workload can mount into the host's `/` via the propagation group. Use `--make-private` on sensitive paths and mount containers with `--mount-propagation=private` / `slave` unless the app genuinely requires `shared`.

Disable by adding `ZC1953` to `disabled_katas` in `.zshellcheckrc`.

---

<a id="zc1954"></a>
### ZC1954 — Warn on `setfattr -n security.capability|security.selinux|security.ima` — bypasses `setcap`/`chcon`

**Severity:** `warning`  
**Auto-fix:** `no`

`setfattr -n security.capability -v …` writes the raw file-capability xattr that the kernel consults when a binary `execve()`s, bypassing the `setcap` wrapper's validation and audit trail. Similarly, `security.selinux` replaces the SELinux label without going through `chcon` / `semanage`, and `security.ima` overwrites the IMA hash that integrity-measurement trusts. These attributes are the raw kernel knobs behind purpose-built tools; script usage is almost always wrong. Use `setcap`, `chcon`/`semanage fcontext`, and `evmctl` instead.

Disable by adding `ZC1954` to `disabled_katas` in `.zshellcheckrc`.

---

<a id="zc1955"></a>
### ZC1955 — Warn on `rfkill block all` / `block wifi|bluetooth|wwan` — disables every radio, cuts wireless

**Severity:** `warning`  
**Auto-fix:** `no`

`rfkill block all` toggles the soft-kill switch on every radio the kernel registered — WiFi, Bluetooth, WWAN, NFC, GPS, UWB — so the host drops off the network in one call. A follow-up `rfkill unblock all` takes seconds to a minute on some drivers and requires the operator to be physically present or have a cellular fallback. Scope the block to a specific type (e.g. `rfkill block bluetooth`) and schedule via `at now + 5 minutes ... rfkill unblock all` so the host recovers on its own.

Disable by adding `ZC1955` to `disabled_katas` in `.zshellcheckrc`.

---

<a id="zc1956"></a>
### ZC1956 — Error on `tailscale up --auth-key=SECRET` — single-use join key visible in argv

**Severity:** `error`  
**Auto-fix:** `no`

`tailscale up --auth-key tskey-auth-…` (and the joined `--auth-key=…` form) passes the Tailscale pre-auth key as a command-line argument. Pre-auth keys grant full tailnet membership, and short-lived or not, the value ends up in `ps`, `/proc/PID/cmdline`, shell history, and any process dump taken before the join completes. Read the key from `TS_AUTHKEY` with `tailscale up --authkey-env=TS_AUTHKEY` (newer tailscaled), or from a file with `tailscale up --auth-key=file:/etc/ts.key` (mode `0400` owned by the provisioning user).

Disable by adding `ZC1956` to `disabled_katas` in `.zshellcheckrc`.

---

<a id="zc1957"></a>
### ZC1957 — Warn on `lvchange -an` / `vgchange -an` — deactivates a live LV/VG, risks mounted-fs corruption

**Severity:** `warning`  
**Auto-fix:** `no`

`lvchange -an VG/LV` (and `vgchange -an VG` for the whole group) deactivates a logical volume by removing its device-mapper entry. If the LV is mounted, writes that the kernel has buffered but not yet flushed may be lost, and any process holding an open fd on the filesystem gets EIO on the next syscall. `umount` the mount first, stop any service keeping files open, verify with `lsof` / `fuser`, and only then `lvchange -an`. For a scripted teardown, prefer `umount` + `lvremove` with a recovery snapshot in hand.

Disable by adding `ZC1957` to `disabled_katas` in `.zshellcheckrc`.

---

<a id="zc1958"></a>
### ZC1958 — Warn on `helm upgrade --force` — delete-and-recreate resources, drops running pods

**Severity:** `warning`  
**Auto-fix:** `no`

`helm upgrade RELEASE CHART --force` flips the upgrade strategy from three-way-merge to `delete + create` for every resource Helm owns. Deployments become new objects, Services lose their `clusterIP` for a beat, and any `PodDisruptionBudget` is bypassed because the resource is deleted, not rolled out. Use plain `helm upgrade` (three-way merge) or `--atomic` / `--wait` for a supervised roll. Reserve `--force` for recovery after a failed upgrade with a stuck resource, not routine deploys.

Disable by adding `ZC1958` to `disabled_katas` in `.zshellcheckrc`.

---

<a id="zc1959"></a>
### ZC1959 — Warn on `trivy … --skip-db-update` / `--skip-update` — scans against a stale vulnerability DB

**Severity:** `warning`  
**Auto-fix:** `no`

`trivy` embeds a vulnerability database that is rehydrated on every scan unless the operator passes `--skip-db-update` (or `--skip-update` on older releases). In CI the flag is tempting — each build then skips a 40 MB download — but the scan then misses every CVE disclosed since the cached DB was last refreshed. Keep the default download, or pre-populate the cache with `trivy image --download-db-only` once per day in a scheduled job, and only pass `--skip-db-update` inside the same job so every scan sees the fresh data.

Disable by adding `ZC1959` to `disabled_katas` in `.zshellcheckrc`.

---

<a id="zc1960"></a>
### ZC1960 — Warn on `az vm run-command invoke` / `aws ssm send-command` — arbitrary commands on remote VM

**Severity:** `warning`  
**Auto-fix:** `no`

`az vm run-command invoke --command-id RunShellScript --scripts "$CMD"` (and the AWS equivalent `aws ssm send-command --document-name AWS-RunShellScript --parameters "commands=['$CMD']"`) runs arbitrary shell on the target instance via the cloud control plane. The identity making the call is whatever role the script's credentials carry; if `$CMD` is composed from any operator or attacker input, the result is remote code execution through IAM. Gate the call behind a shell-escape-safe templater, pin the document version / script to a reviewed asset in blob / S3, and require MFA on the invoking role.

Disable by adding `ZC1960` to `disabled_katas` in `.zshellcheckrc`.

---

<a id="zc1961"></a>
### ZC1961 — Warn on `gcloud iam service-accounts keys create` — mints a long-lived service-account JSON key

**Severity:** `warning`  
**Auto-fix:** `no`

`gcloud iam service-accounts keys create key.json --iam-account=SA@PROJECT` exports an RSA key pair wrapped in a JSON file. Once written it is effectively a forever-valid bearer credential: no automatic rotation, no refresh, and a single "leaked by a `cat key.json`" is game-over. Prefer Workload Identity Federation (`gcloud iam workload-identity-pools …`), short-lived impersonation via `gcloud auth print-access-token --impersonate-service-account=SA`, or the key-less GCE/GKE attached service account. Reserve static JSON keys for provably off-platform callers.

Disable by adding `ZC1961` to `disabled_katas` in `.zshellcheckrc`.

---

<a id="zc1962"></a>
### ZC1962 — Warn on `kustomize build --load-restrictor=LoadRestrictionsNone` — path-traversal in overlays

**Severity:** `warning`  
**Auto-fix:** `no`

Kustomize's default `LoadRestrictionsRootOnly` limits every base, patch, configMapGenerator, and secretGenerator to paths under the current kustomization root. `kustomize build … --load-restrictor=LoadRestrictionsNone` (also the legacy spelling `--load_restrictor none` / `--load-restrictor=LoadRestrictionsNone_WarnForAll`) drops that guard, so an overlay from an untrusted remote base can reference `../../secrets/prod.env` or absolute paths and pull them into the render. Keep the default; if a legitimate overlay needs a sibling file, vendor it in.

Disable by adding `ZC1962` to `disabled_katas` in `.zshellcheckrc`.

---

<a id="zc1963"></a>
### ZC1963 — Warn on `npx pkg` / `pnpm dlx pkg` / `bunx pkg` without a version pin — runs latest registry code

**Severity:** `warning`  
**Auto-fix:** `no`

`npx PKG`, `pnpm dlx PKG`, `bunx PKG`, and `bun x PKG` fetch the named package from the npm registry and execute its `bin` entry. Without a version pin (`pkg@1.2.3`), each run resolves to the registry's `latest` tag — a compromised maintainer, squatted name, or even a mistyped package is enough to land attacker code in the build. Pin the exact version (`npx pkg@1.2.3`), cache the binary under `./node_modules/.bin/` via a regular `npm install`, or verify the tarball signature before execution.

Disable by adding `ZC1963` to `disabled_katas` in `.zshellcheckrc`.

---

<a id="zc1964"></a>
### ZC1964 — Warn on `uvx pkg` / `uv tool run pkg` / `pipx run pkg` without a version pin — runs latest PyPI release

**Severity:** `warning`  
**Auto-fix:** `no`

`uvx PKG`, `uv tool run PKG`, and `pipx run PKG` each resolve the package against PyPI and execute its entry point. Without a version constraint (`pkg==1.2.3` or `pkg@1.2.3` for uv), every run takes whatever the registry currently serves — a typosquatted lookalike, a compromised maintainer release, or a sudden major-version bump lands untested code in the pipeline. Pin the version at the call site or use `uv tool install pkg==X.Y.Z` + `uv tool run pkg` so the lockfile is the source of truth.

Disable by adding `ZC1964` to `disabled_katas` in `.zshellcheckrc`.

---

<a id="zc1965"></a>
### ZC1965 — Error on `systemd-cryptenroll --wipe-slot=all` — wipes every LUKS key slot

**Severity:** `error`  
**Auto-fix:** `no`

`systemd-cryptenroll --wipe-slot=all $DEV` removes every key slot on the LUKS volume — passphrase, recovery key, TPM2, FIDO2, PKCS#11 — in one call. `--wipe-slot=recovery` / `--wipe-slot=empty` are scoped; the `all` form is a one-shot brick with no confirmation. Either enrol the new slot first and then wipe the specific index you are retiring (`--wipe-slot=<n>`), or back up the header with `cryptsetup luksHeaderBackup` before the call so recovery is possible.

Disable by adding `ZC1965` to `disabled_katas` in `.zshellcheckrc`.

---

<a id="zc1966"></a>
### ZC1966 — Error on `zpool import -f` / `zpool export -f` — forced ZFS pool op bypasses hostid/txg checks

**Severity:** `error`  
**Auto-fix:** `no`

`zpool import -f $POOL` force-imports a pool even when the on-disk hostid differs — i.e. the pool is already imported on another host (multipath/SAN, shared JBOD, HA cluster). The second import writes to the same vdevs and silently corrupts the pool. `zpool export -f` skips the graceful-flush path and detaches vdevs with in-flight txgs, which can lose the tail of the ZIL. Export without `-f` after `zfs unmount -a`; import without `-f` after verifying `zpool import` (no target) reports the pool as `ONLINE` and the hostid matches.

Disable by adding `ZC1966` to `disabled_katas` in `.zshellcheckrc`.

---

<a id="zc1967"></a>
### ZC1967 — Warn on `setopt PROMPT_SUBST` — expansions inside `$PROMPT` evaluate command substitution every redraw

**Severity:** `warning`  
**Auto-fix:** `no`

`setopt PROMPT_SUBST` turns on parameter, command, and arithmetic substitution inside `$PS1`/`$PROMPT`/`$RPROMPT`. Any value that lands in the prompt from an untrusted source — a git branch name, a checkout path, a hostname in `/etc/hostname`, an env var set by a spawned tool — is reparsed as shell code on every redraw, so a branch like `$(id>/tmp/p)` runs each time the cursor returns. Prefer Zsh prompt escapes (`%n`, `%d`, `%~`, `%m`, `vcs_info`) which already sanitise their inputs, or scope with `setopt LOCAL_OPTIONS PROMPT_SUBST` inside the prompt-building function instead of flipping the option globally.

Disable by adding `ZC1967` to `disabled_katas` in `.zshellcheckrc`.

---

<a id="zc1968"></a>
### ZC1968 — Warn on `dnf versionlock add` / `yum versionlock add` — pins RPM, blocks CVE updates

**Severity:** `warning`  
**Auto-fix:** `no`

`dnf versionlock add pkg` (and the legacy `yum versionlock add pkg`) write an entry to `/etc/dnf/plugins/versionlock.list` that excludes the package from future `dnf update` / `dnf upgrade` runs. Mirrors `apt-mark hold` on Debian (ZC1550): the lock survives reboots and unattended-upgrades never sees the newer rpm, so kernel, openssl, or glibc CVEs pile up unseen. Document the exact reason in the commit, pair the lock with a scheduled `dnf versionlock delete` date, and prefer excluding the problematic transaction via `--exclude` or a one-shot `dnf update --setopt=exclude=pkg` rather than a persistent pin.

Disable by adding `ZC1968` to `disabled_katas` in `.zshellcheckrc`.

---

<a id="zc1969"></a>
### ZC1969 — Warn on `zsh -f` / `zsh -d` — skips `/etc/zsh*` and `~/.zsh*` startup files

**Severity:** `warning`  
**Auto-fix:** `no`

`zsh -f` is the short form of `--no-rcs`, which skips every personal and system-wide startup file: `/etc/zshenv`, `/etc/zprofile`, `/etc/zshrc`, `/etc/zlogin`, `~/.zshenv`, `~/.zshrc`, `~/.zlogin`. `zsh -d` (`--no-globalrcs`) drops only the `/etc/zsh*` set but keeps per-user ones. Either form strips corp-mandated settings — proxy/hosts overrides, audit hooks, umask, `HISTFILE` redirection, `PATH` hardening — silently. Use it deliberately only for a pristine test harness or a minimal repro; never as the shebang of a production script. When isolation is required, prefer `env -i zsh` with an explicit allow-list of variables.

Disable by adding `ZC1969` to `disabled_katas` in `.zshellcheckrc`.

---

<a id="zc1970"></a>
### ZC1970 — Warn on `losetup -P` / `kpartx -a` / `partprobe` on untrusted image — runs kernel partition parser

**Severity:** `warning`  
**Auto-fix:** `no`

`losetup -P $LOOP $IMG`, `kpartx -av $IMG`, and `partprobe $LOOP` all tell the kernel to rescan a block device's partition table and emit `/dev/loopNpX` (or dm-N) entries. When the image comes from an untrusted source — a customer-supplied VM disk, a downloaded installer, a forensic capture — the scan runs MBR/GPT/LVM parsers over attacker-controlled bytes and has historically triggered kernel CVEs (fsconfig heap overflow, ext4 mount bugs). Do the inspection in a throwaway VM or an offline parser like `fdisk -l $IMG` / `sfdisk --dump $IMG` that reads without kernel scan, and only attach partitions with `losetup -P` after the layout is known-good.

Disable by adding `ZC1970` to `disabled_katas` in `.zshellcheckrc`.

---

<a id="zc1971"></a>
### ZC1971 — Warn on `unsetopt GLOBAL_RCS` / `setopt NO_GLOBAL_RCS` — skips `/etc/zprofile`, `/etc/zshrc`, `/etc/zlogin`, `/etc/zlogout`

**Severity:** `warning`  
**Auto-fix:** `no`

`GLOBAL_RCS` is on by default; only `/etc/zshenv` is sourced before it can be toggled. Flipping the option off (either `unsetopt GLOBAL_RCS` or `setopt NO_GLOBAL_RCS`) tells Zsh to skip `/etc/zprofile`, `/etc/zshrc`, `/etc/zlogin`, and `/etc/zlogout` — which is exactly where admins put corp-wide `PATH` hardening, audit hooks, umask, `HISTFILE` redirection, and proxy variables. A login-shell script that disables the option in `/etc/zshenv` neutralises every downstream system rc without a trace. Keep the option on; if a specific helper needs pristine setup use `emulate -LR zsh` inside a function or spawn `env -i zsh -f` scoped to that helper.

Disable by adding `ZC1971` to `disabled_katas` in `.zshellcheckrc`.

---

<a id="zc1972"></a>
### ZC1972 — Error on `dmsetup remove_all` / `dmsetup remove -f` — tears down live LVM/LUKS/multipath mappings

**Severity:** `error`  
**Auto-fix:** `no`

`dmsetup remove_all` iterates every device-mapper node on the host — LVM logical volumes, LUKS containers, multipath aggregates, `cryptsetup` mappings — and asks the kernel to drop each one. `dmsetup remove --force $NAME` targets a single mapping but still evicts it with in-flight I/O. When any of those devices is mounted or backing a running VM, new I/O to it returns `ENXIO`, `fsck` is no longer possible, and LVM metadata needs a cold reboot to reappear. Use `dmsetup remove $NAME` without `--force` after `umount`/`vgchange -an`/`cryptsetup close`, and never `remove_all` on a host you care about.

Disable by adding `ZC1972` to `disabled_katas` in `.zshellcheckrc`.

---

<a id="zc1973"></a>
### ZC1973 — Warn on `setopt POSIX_IDENTIFIERS` — restricts parameter names to ASCII, breaks Unicode `$var`

**Severity:** `warning`  
**Auto-fix:** `no`

Zsh accepts Unicode parameter names by default: `$café`, `$π`, `$данные` all parse. `setopt POSIX_IDENTIFIERS` tightens that to the POSIX subset — ASCII letters, digits, underscore, not starting with a digit. Once the option is on, every later `${café}` or `café=1` is a parse error, and scripts/libraries that expose i18n-named vars stop loading. If you need POSIX identifiers for a specific helper, scope it inside a function with `emulate -LR sh`; leave the global option off so the rest of the shell keeps the Zsh behaviour the user expects.

Disable by adding `ZC1973` to `disabled_katas` in `.zshellcheckrc`.

---

<a id="zc1974"></a>
### ZC1974 — Error on `ipset flush` / `ipset destroy` — nukes named sets referenced by iptables/nft rules

**Severity:** `error`  
**Auto-fix:** `no`

`ipset flush` empties every entry from a named IP set; `ipset destroy` (no args) removes every set on the host. iptables/nft rules of the form `-m set --match-set $NAME src` then reference a set that is either empty or gone, so block-lists disappear instantly and allow-lists stop whitelisting — the ruleset falls through to its default policy. Target a specific set by name (`ipset destroy $NAME` after confirming no rule references it), or add new entries with `ipset add` instead of rebuilding from scratch. Reload atomically with `ipset restore -! < snapshot` if a full replace is genuinely needed.

Disable by adding `ZC1974` to `disabled_katas` in `.zshellcheckrc`.

---

<a id="zc1975"></a>
### ZC1975 — Warn on `unsetopt EXEC` / `setopt NO_EXEC` — parser keeps scanning, commands stop running

**Severity:** `warning`  
**Auto-fix:** `no`

`EXEC` is on by default; the shell both parses and runs each command. Turning it off (`unsetopt EXEC` or `setopt NO_EXEC`) tells Zsh to parse everything but silently skip the execution step — nothing fires, yet parameter assignments on the same line don't either, `$?` stays frozen, and functions that follow look defined but never run. That is the semantics behind `zsh -n script.zsh` for a pure syntax check; flipping the option in the middle of a production script converts every later line into a no-op without a visible error. Run syntax checks via `zsh -n` from the outside, never by flipping `EXEC` in-line.

Disable by adding `ZC1975` to `disabled_katas` in `.zshellcheckrc`.

---

<a id="zc1976"></a>
### ZC1976 — Error on `exportfs -au` / `exportfs -u` — unexports live NFS shares, clients get `ESTALE`

**Severity:** `error`  
**Auto-fix:** `no`

`exportfs -au` unexports every NFS share on the server; `exportfs -u HOST:/PATH` removes a single share. Any client that currently has the export mounted is not notified — the next read/write returns `ESTALE`, the mount looks live but every open fd fails, and the only recovery is a client-side `umount -l` + remount. `exportfs -f` (flush) is almost always what you actually want after an `/etc/exports` edit; keep `-u`/`-au` for planned shutdowns with a coordinated client `umount` first.

Disable by adding `ZC1976` to `disabled_katas` in `.zshellcheckrc`.

---

<a id="zc1977"></a>
### ZC1977 — Warn on `setopt CHASE_DOTS` — `cd ..` physically resolves before walking up, breaking logical paths

**Severity:** `warning`  
**Auto-fix:** `no`

Default Zsh keeps `..` logical: from `/app/current/lib` (where `/app/current` → `/app/releases/v5`), `cd ..` goes back to `/app/current`, matching the user's mental model and blue/green deployment symlinks. `setopt CHASE_DOTS` flips that — `..` first resolves the current directory to its physical inode, so the same `cd ..` lands in `/app/releases/v5` and the next `cd config` looks for `/app/releases/config` instead of `/app/config`. Scripts that rely on `${PWD}` staying logical or on `cd ../foo` matching the typed path break silently. Leave the option off; use `cd -P` one-shot when a specific call really needs physical resolution.

Disable by adding `ZC1977` to `disabled_katas` in `.zshellcheckrc`.

---

<a id="zc1978"></a>
### ZC1978 — Warn on `tftp` — cleartext, unauthenticated UDP transfer

**Severity:** `warning`  
**Auto-fix:** `no`

`tftp` has no authentication at all and moves the payload in plaintext over UDP/69 — any packet capture on the path recovers the full transfer and an attacker at the server can push an arbitrary file under the expected name without noticing a lack of credentials. The dual-channel design is also routinely mishandled by NAT/firewall gear. For PXE-style provisioning that historically used `tftp`, fetch a signed payload over HTTPS with `curl` and verify the signature locally before use. (See ZC1200 for `ftp`, the authenticated-but-plaintext sibling.)

Disable by adding `ZC1978` to `disabled_katas` in `.zshellcheckrc`.

---

<a id="zc1979"></a>
### ZC1979 — Warn on `setopt HIST_FCNTL_LOCK` — `fcntl()` lock on NFS `$HISTFILE` stalls or deadlocks

**Severity:** `warning`  
**Auto-fix:** `no`

Off by default, Zsh serialises writes to `$HISTFILE` with its own lock-file dance next to the history. `setopt HIST_FCNTL_LOCK` switches to POSIX `fcntl()` advisory locking — which is the safer primitive on local filesystems, but on NFS homes the lock is proxied through `rpc.lockd` and a single hung client or rebooted NFS server leaves every other shell blocked the next time it tries to write history. The interactive shell appears frozen on prompt return, and scripts that source user rc files hang in `zshaddhistory`. Keep the option off on NFS homes; only turn it on when `$HISTFILE` lives on a local filesystem (ext4, xfs, btrfs, zfs local pool) that implements `fcntl()` without network round-trips.

Disable by adding `ZC1979` to `disabled_katas` in `.zshellcheckrc`.

---

<a id="zc1980"></a>
### ZC1980 — Error on `udevadm trigger --action=remove` — replays `remove` uevents, detaches live devices

**Severity:** `error`  
**Auto-fix:** `no`

`udevadm trigger --action=remove` (also spelled `-c remove`) walks `/sys` and synthesises a `remove` uevent for every matching device. The kernel reacts as if every matched disk, NIC, GPU, or USB node was physically yanked — SATA controllers detach drives that back mounted filesystems, netdevs disappear mid-session, and `systemd-udevd` fires per-device cleanup rules it was never meant to run on a live host. The normal way to replay `add`/`change` events after a rules edit is `udevadm control --reload` followed by `udevadm trigger` with the default action (`change`); scope any `--action=remove` to a specific device subsystem with `--subsystem-match=` + `--attr-match=` and test on a non-production box first.

Disable by adding `ZC1980` to `disabled_katas` in `.zshellcheckrc`.

---

<a id="zc1981"></a>
### ZC1981 — Warn on `exec -a NAME cmd` — replaces `argv[0]`, hides the real binary from `ps`

**Severity:** `warning`  
**Auto-fix:** `no`

`exec -a NAME $BIN` tells Zsh to set `argv[0]` of the `exec`'d process to `NAME` instead of the actual program path. `ps`, `top`, `proc`-based audit tools, and systemd's unit accounting all see `NAME` — the real binary on disk is only discoverable from `/proc/PID/exe`, which most monitoring does not read. The feature has legitimate uses (login shells spelling themselves `-zsh` so tty/shell detection works) but also makes a great disguise for a reverse shell or a cron-triggered helper. Keep `exec -a` out of production scripts unless the intent is documented; prefer running the binary at its real path so operators can match process name to on-disk file.

Disable by adding `ZC1981` to `disabled_katas` in `.zshellcheckrc`.

---

<a id="zc1982"></a>
### ZC1982 — Error on `ipcrm -a` — removes every SysV IPC object, breaks Postgres/Oracle/shm apps

**Severity:** `error`  
**Auto-fix:** `no`

`ipcrm -a` deletes every System V shared-memory segment, semaphore set, and message queue owned by the caller (or, as root, every object on the host). Long-running services that rely on SysV IPC — PostgreSQL's shared buffers, Oracle's SGA, the `sysv` session store in several RDBMS test suites, shm-based mutexes in batch pipelines — lose their backing store mid-transaction and either SIGSEGV or return `EINVAL` on the next access. Scope the removal: `ipcrm -m ID`/`-s ID`/`-q ID` against the specific identifier reported by `ipcs -a`, after confirming no running process attached to it.

Disable by adding `ZC1982` to `disabled_katas` in `.zshellcheckrc`.

---

<a id="zc1983"></a>
### ZC1983 — Warn on `setopt CSH_JUNKIE_QUOTES` — single/double-quoted strings that span lines become errors

**Severity:** `warning`  
**Auto-fix:** `no`

With `CSH_JUNKIE_QUOTES` off (the default), Zsh lets `"foo\nbar"` and `'line1\nline2'` span physical lines. Setting the option on makes the parser emit an error on the first newline inside a quoted string — which breaks any existing multi-line SQL, JSON, or here-style payload that the script has been inlining up to this point. Functions that are autoloaded later or sourced from third-party helpers fail to parse, and the diagnostic points at the closing quote, not at the option toggle. Leave the option off; if csh-style strictness is genuinely required, scope with `emulate -LR csh` inside the single helper that needs it.

Disable by adding `ZC1983` to `disabled_katas` in `.zshellcheckrc`.

---

<a id="zc1984"></a>
### ZC1984 — Error on `sgdisk -Z` / `sgdisk -o` — erases the GPT partition table on the target disk

**Severity:** `error`  
**Auto-fix:** `no`

`sgdisk -Z $DISK` (`--zap-all`) wipes the primary GPT, the protective MBR, and the backup GPT at the end of the device. `sgdisk -o $DISK` (`--clear`) replaces the existing partition table with a fresh empty GPT. Either command detaches every partition, LVM PV, LUKS container, and filesystem header on the device — when the target variable resolves to a wrong path (tab completion, `$DISK` defaulted to `/dev/sda`), the host becomes unbootable. Require an `lsblk $DISK` + `blkid $DISK` preflight in the script, route the action through `--pretend` (`-t`) first, and keep a `sgdisk --backup=/root/$DISK.gpt $DISK` image before any zap.

Disable by adding `ZC1984` to `disabled_katas` in `.zshellcheckrc`.

---

<a id="zc1985"></a>
### ZC1985 — Warn on `setopt SH_FILE_EXPANSION` — expansion order flips from Zsh-native to sh/bash, `~` leaks

**Severity:** `warning`  
**Auto-fix:** `no`

Default Zsh runs parameter expansion first, then filename/`~` expansion — so a `VAR='~/cache'` keeps the tilde literal when you do `mkdir -p -- $VAR` because the `~` never leaves the value. `setopt SH_FILE_EXPANSION` (POSIX/sh ordering) flips the pass: filename expansion runs first on the raw text, then parameter expansion happens, so the same line suddenly makes the tilde resolve to `$HOME`, paths pointing at `~evil/.cache` resolve into another user's home, and `=cmd` spellings look up `$PATH` silently. Keep the option off; when a specific helper needs POSIX ordering use `emulate -LR sh` inside that function.

Disable by adding `ZC1985` to `disabled_katas` in `.zshellcheckrc`.

---

<a id="zc1986"></a>
### ZC1986 — Warn on `touch -d` / `-t` / `-r` — explicit timestamp write is a common antiforensics pattern

**Severity:** `warning`  
**Auto-fix:** `no`

`touch -d "2 years ago" $F`, `touch -t YYYYMMDDhhmm $F`, and `touch -r $REF $F` all write the atime/mtime to a specific value rather than the current clock. Legitimate uses exist — re-stamping a mirror to match upstream, generating deterministic tarballs for reproducible-build pipelines, `rsync --archive` edge cases — but the pattern also matches the classic "age the dropped file" antiforensics trick where an attacker normalises a new binary to look as old as its neighbours so `find -mtime`- based triage misses it. Audit rules should flag these forms in production scripts; in reproducible-build contexts, keep the timestamp derived from `SOURCE_DATE_EPOCH` via `touch -d @$SOURCE_DATE_EPOCH` so operators can recognise the intent at a glance.

Disable by adding `ZC1986` to `disabled_katas` in `.zshellcheckrc`.

---

<a id="zc1987"></a>
### ZC1987 — Warn on `setopt BRACE_CCL` — `{a-z}` expands to each character instead of staying literal

**Severity:** `warning`  
**Auto-fix:** `no`

`BRACE_CCL` is off by default: `echo {a-z}` stays literal `a-z` in Zsh, which is what most scripts that only want the numeric range form `{1..10}` actually expect. `setopt BRACE_CCL` promotes single-character ranges and enumerations inside braces to csh-style character-class expansion, so `echo {a-z}` suddenly prints every letter from `a` to `z` and `echo {ABC}` becomes `A B C`. Any later command line that embeds single-character ranges — regex fragments, hex masks, CI job names with stage suffixes — expands unexpectedly. Leave the option off; use `{a..z}` when a real range is wanted and quote literals that contain `{…}`.

Disable by adding `ZC1987` to `disabled_katas` in `.zshellcheckrc`.

---

<a id="zc1988"></a>
### ZC1988 — Error on `nsupdate -y HMAC:NAME:SECRET` — TSIG key visible in argv and shell history

**Severity:** `error`  
**Auto-fix:** `no`

`nsupdate -y [alg:]name:base64secret` hands the TSIG shared secret directly on the command line, so `ps auxf`, `/proc/PID/cmdline`, and `$HISTFILE` all capture the key — and whoever owns the key can rewrite any zone that trusts it (DNS hijack, MX hijack, ACME domain-validation bypass). `nsupdate -k /etc/named/KEY` (or `-k $KEYFILE` with `0600` perms) reads the key from disk without exposing it. If the secret must come from a secret store, pipe it through `nsupdate -k /dev/stdin <<<"$KEYFILE_CONTENTS"` so the raw material never lands in argv.

Disable by adding `ZC1988` to `disabled_katas` in `.zshellcheckrc`.

---

<a id="zc1989"></a>
### ZC1989 — Warn on `setopt REMATCH_PCRE` — `[[ =~ ]]` regex flips from POSIX ERE to PCRE, changes semantics

**Severity:** `warning`  
**Auto-fix:** `no`

By default Zsh's `[[ $str =~ pattern ]]` uses POSIX extended regex (ERE). `setopt REMATCH_PCRE` (after `zmodload zsh/pcre`) swaps the engine to PCRE for every later match. Patterns that pass through both engines change meaning subtly: `\b` is a word boundary in PCRE but a literal `b` in ERE, `\d`/`\s`/`\w` work in PCRE but not ERE, lookahead/lookbehind (`(?=…)`) parse in PCRE but error in ERE, and inline flags `(?i)` only exist in PCRE. Flipping the option globally silently rewrites the meaning of every existing regex — prefer an explicit `pcre_match`/`pcre_compile` call when PCRE is needed, or a `setopt LOCAL_OPTIONS REMATCH_PCRE` inside the single function that uses PCRE syntax.

Disable by adding `ZC1989` to `disabled_katas` in `.zshellcheckrc`.

---

<a id="zc1990"></a>
### ZC1990 — Warn on `openssl passwd -crypt` / `-1` / `-apr1` — obsolete password hash formats

**Severity:** `warning`  
**Auto-fix:** `no`

`openssl passwd -crypt` emits DES-crypt, 8-char truncated and crackable in seconds on modern hardware. `-1` is FreeBSD-style MD5, unsuitable for storage, long broken. `-apr1` is Apache's MD5-based variant with the same weakness. Any hash produced by these flags lands in `/etc/shadow`, an htpasswd file, or a database row where an attacker can offline-crack the whole batch with a single GPU. Use `-5` (SHA-256-crypt), `-6` (SHA-512-crypt), or prefer a dedicated KDF-based hasher — `mkpasswd -m yescrypt`, `htpasswd -B` (bcrypt), or `argon2` — so brute-force cost scales with hardware.

Disable by adding `ZC1990` to `disabled_katas` in `.zshellcheckrc`.

---

<a id="zc1991"></a>
### ZC1991 — Warn on `setopt CSH_NULLCMD` — bare `> file` raises an error instead of running `$NULLCMD`

**Severity:** `warning`  
**Auto-fix:** `no`

Default Zsh executes `$NULLCMD` (initially `cat`) when a line has redirections but no command, so `> file < input` copies input to file and `< file` pages through it with `$READNULLCMD` (initially `more`). `setopt CSH_NULLCMD` drops the Zsh convention and follows csh — any command line without an explicit command is a parse error, regardless of redirections. Scripts that rely on the bare-redirect idiom (log truncation via `> $LOG`, drop-in includes via `< file`, piped filters built from aliases) stop working with a confusing `parse error near '<'`. Keep the option off; write `: > file` (or `true > file`) explicitly when you mean to truncate.

Disable by adding `ZC1991` to `disabled_katas` in `.zshellcheckrc`.

---

<a id="zc1992"></a>
### ZC1992 — Warn on `pkexec cmd` — PolicyKit privilege elevation is historically bug-prone and hard to audit from scripts

**Severity:** `warning`  
**Auto-fix:** `no`

`pkexec` lifts a command to the UID configured in a PolicyKit `.policy` file — typically root — after consulting an authorisation agent. From a non-interactive script the agent has no way to prompt, so the call either depends on a pre-authorised `.policy` override or fails in a confusing manner. The binary also has a poor CVE track record (CVE-2021-4034 pwnkit, CVE-2017-16089, envvar handling bugs) and its audit trail is split across journald and `/var/log/auth.log`. Use `sudo` with a targeted `sudoers` drop-in for scripted privilege elevation, or run the script under a systemd unit with `User=` / `AmbientCapabilities=` when specific capabilities are needed.

Disable by adding `ZC1992` to `disabled_katas` in `.zshellcheckrc`.

---

<a id="zc1993"></a>
### ZC1993 — Warn on `setopt KSH_TYPESET` — `typeset var=$val` starts word-splitting the RHS

**Severity:** `warning`  
**Auto-fix:** `no`

Off by default, Zsh treats every `typeset`/`declare` assignment like a shell assignment: the whole RHS after `=` is one token, so `typeset msg="a b c"` produces a single-element string. `setopt KSH_TYPESET` follows ksh instead — each word on the `typeset` line is its own assignment or name, and the shell re-splits the RHS on whitespace. Functions that used to accept `typeset path=$HOME/My Files` suddenly treat `Files` as a second variable name, and `local` (an alias for `typeset` inside functions) inherits the same change. Keep the option off; if ksh compatibility is genuinely needed, scope with `emulate -LR ksh` inside the helper that needs it.

Disable by adding `ZC1993` to `disabled_katas` in `.zshellcheckrc`.

---

<a id="zc1994"></a>
### ZC1994 — Error on `lvreduce -f` / `lvreduce -y` — shrinks the LV without checking the filesystem above

**Severity:** `error`  
**Auto-fix:** `no`

`lvreduce -L SIZE $LV` cuts the block device below an existing filesystem. The confirmation prompt exists precisely because ext4/xfs/btrfs do not shrink themselves — LVM happily lops off the tail even though the filesystem still believes those blocks are allocated. `-f` / `-y` / `--force` / `--yes` skip the prompt, and the next mount returns corruption or missing files. Shrink the filesystem first with `resize2fs $LV $NEWSIZE` (or `xfs_growfs` equivalent — xfs cannot shrink, so offline backup + recreate), verify `df` / `fsck`, then `lvreduce --resizefs` (which performs both steps atomically) instead of bypassing the prompt.

Disable by adding `ZC1994` to `disabled_katas` in `.zshellcheckrc`.

---

<a id="zc1995"></a>
### ZC1995 — Warn on `unsetopt BGNICE` — background jobs run at full interactive priority, starve the foreground

**Severity:** `warning`  
**Auto-fix:** `no`

Default Zsh applies `nice +5` to every backgrounded job so long-running work does not starve the interactive session. `unsetopt BGNICE` (or `setopt NO_BGNICE`) turns that off and bg jobs compete at the same priority as the foreground shell — SSH keystroke handling, editor redraws, and `cmd &` batch fan-out all feel laggy, and a single CPU-bound bg job can peg every core of a container it shares with a human operator. Keep the option on; when a background job legitimately needs full priority (audio pipeline, realtime simulator), wrap just that one with `nice -n 0 -- cmd &` or a systemd unit with `Nice=` instead of flipping globally.

Disable by adding `ZC1995` to `disabled_katas` in `.zshellcheckrc`.

---

<a id="zc1996"></a>
### ZC1996 — Warn on `unshare -U` / `-r` — unprivileged user namespace maps caller to root inside the NS

**Severity:** `warning`  
**Auto-fix:** `no`

`unshare -U` opens a new user namespace and `-r` / `--map-root-user` maps the caller's UID to `0` inside it. That's the foundation of rootless containers (bubblewrap, podman rootless, flatpak) and is legitimate in that context. It is also the standard opening move for a long list of LPE chains — once you are uid `0` in a user namespace you can create additional mount/net/cgroup namespaces, run `mount -t overlay` against attacker-controlled dirs, and probe kernel attack surface that is normally gated on `CAP_SYS_ADMIN`. Audit rules should flag the pattern in production scripts; if a rootless runtime really needs it, route through the runtime binary (`bwrap`, `podman --rootless`) so the invocation is recognisable.

Disable by adding `ZC1996` to `disabled_katas` in `.zshellcheckrc`.

---

<a id="zc1997"></a>
### ZC1997 — Warn on `setopt HIST_NO_FUNCTIONS` — function definitions skipped from `$HISTFILE`, breaks forensic trail

**Severity:** `warning`  
**Auto-fix:** `no`

Default Zsh writes every command you type, including function definitions, to `$HISTFILE`. `setopt HIST_NO_FUNCTIONS` suppresses storage of commands that define a function. On a multi-admin box or a shared root account this breaks the forensic trail — the function the attacker just defined (or that an operator typed before running the destructive bit) vanishes from history while the invocation that used it still shows, leaving responders with a command that references a name that no longer exists on disk or in any log. Keep the option off and scope any hiding needs with the Zsh hook `zshaddhistory { return 1 }` inside a function where the secret actually lives.

Disable by adding `ZC1997` to `disabled_katas` in `.zshellcheckrc`.

---

<a id="zc1998"></a>
### ZC1998 — Error on `tpm2_clear` / `tpm2 clear` — wipes TPM storage hierarchy, kills every sealed key

**Severity:** `error`  
**Auto-fix:** `no`

`tpm2_clear -c p` (or `tpm2 clear -c p`) invokes the TPM 2.0 `TPM2_Clear` command, which invalidates every object sealed against the storage hierarchy — LUKS-TPM2 keyslots, systemd-cryptenroll's `--tpm2-device` slot, sshd TPM-backed host keys, and SecureBoot measured-boot state. The machine can still boot but any disk that unlocked through the TPM now needs a recovery passphrase, and every TLS cert issued from a TPM-sealed CA loses its anchor. There is no undo. Run `tpm2_clear` only under a documented recovery runbook with the recovery material in hand; never put it in an automated scheduled script.

Disable by adding `ZC1998` to `disabled_katas` in `.zshellcheckrc`.

---

<a id="zc1999"></a>
### ZC1999 — Error on `setopt AUTO_NAMED_DIRS` — unknown option, typo of `AUTO_NAME_DIRS`

**Severity:** `error`  
**Auto-fix:** `no`

`AUTO_NAMED_DIRS` (with the trailing `D`) is not a real Zsh option — `setopt AUTO_NAMED_DIRS` fails with `no such option` and the dir-to-`~name` auto-registration the author likely wanted is never enabled. The canonical spelling is `AUTO_NAME_DIRS` (see ZC1934 for its semantics and why flipping it on is usually wrong). Drop the typo and, if you actually want the behaviour, reach for `hash -d NAME=PATH` explicitly or scope `setopt LOCAL_OPTIONS AUTO_NAME_DIRS` inside the single helper that needs named-directory expansion.

Disable by adding `ZC1999` to `disabled_katas` in `.zshellcheckrc`.

---

<a id="zc2000"></a>
### ZC2000 — Error on `kubectl taint nodes $NODE key=value:NoExecute` — evicts every non-tolerating pod off the node

**Severity:** `error`  
**Auto-fix:** `no`

A `NoExecute` taint kicks every existing pod off the node unless the pod spec explicitly tolerates it. Draining one node during a rolling upgrade is one thing; a script that types the taint wrong (typoed toleration value, applying to `--all` nodes, or iterating a node list without a pause) can empty a whole cluster in seconds and trigger cascade reschedules that overwhelm the scheduler. Prefer `kubectl drain $NODE` (which respects PodDisruptionBudget and runs PreStop hooks) or a `NoSchedule` taint for gentle drain; reserve `NoExecute` for genuine incident response with a runbook and a safety countdown.

Disable by adding `ZC2000` to `disabled_katas` in `.zshellcheckrc`.

---

<a id="zc2001"></a>
### ZC2001 — Warn on `unsetopt EVAL_LINENO` — `$LINENO` inside `eval` stops tracking source, stack traces go blank

**Severity:** `warning`  
**Auto-fix:** `no`

On by default, Zsh's `EVAL_LINENO` keeps `$LINENO`, `$funcfiletrace`, and `$funcstack` pointing at the line inside the `eval`ed string where the error actually happened. Turning the option off (`unsetopt EVAL_LINENO` or `setopt NO_EVAL_LINENO`) reverts to pre-Zsh-4.3 behaviour: `$LINENO` collapses to the line that launched the `eval`, so every runtime error inside a generated config, a lazy-loaded function, or a `compile`d string reports the same line number and the stack trace loses every frame past the eval. Keep the option on; if strict POSIX-matching line numbers are needed inside one helper, scope with `emulate -LR sh` in that function.

Disable by adding `ZC2001` to `disabled_katas` in `.zshellcheckrc`.

---

<a id="zc2002"></a>
### ZC2002 — Error on `crictl rmi -a` / `crictl rm -af` — wipes every image/container on the Kubernetes node

**Severity:** `error`  
**Auto-fix:** `no`

`crictl` talks directly to the node's CRI runtime (containerd, CRI-O), below the kubelet and the cluster API. `crictl rmi -a` removes every cached image including the ones currently backing running pods — the kubelet must immediately re-pull from the registry, and image-pull rate limits or network blips turn the node Unready. `crictl rm -af` force-removes every container on the node, killing pods without running PreStop hooks or honoring PodDisruptionBudget. Route maintenance through `kubectl drain $NODE` + `kubectl delete pod --grace-period=30`; use `crictl` at most on a cordoned, drained node with a documented recovery plan.

Disable by adding `ZC2002` to `disabled_katas` in `.zshellcheckrc`.

---

<a id="zc2003"></a>
### ZC2003 — Warn on `setopt KSH_ZERO_SUBSCRIPT` — `$arr[0]` stops aliasing the first element

**Severity:** `warning`  
**Auto-fix:** `no`

Default Zsh treats `$arr[0]` as a quirk-compatibility alias for `$arr[1]` — `arr=(a b c); echo $arr[0]` prints `a`, and `arr[0]=new` rewrites the first element. `setopt KSH_ZERO_SUBSCRIPT` flips that to ksh semantics: `$arr[0]` becomes a distinct slot (the element just before the 1-indexed head, which Zsh stores separately), so reads silently switch to empty string and `arr[0]=new` no longer touches `$arr[1]`. Any Zsh code that intentionally used `$arr[0]` as a shortcut breaks, and ported Bash/ksh code that assumes 0-indexed access meets a split-world model. Leave the option off; use `$arr[1]` explicitly when you want the first element, and adopt `KSH_ARRAYS` scoped with `emulate -LR ksh` for ksh-style code paths.

Disable by adding `ZC2003` to `disabled_katas` in `.zshellcheckrc`.

---

