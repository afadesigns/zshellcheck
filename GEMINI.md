- **ZShellCheck Versioning:** Versioning follows `Major.Minor.Patch` based on the total number of implemented Katas.
  - **Major:** Represents each full thousand of Katas. (e.g., 1000 Katas = Major `1`, 2500 Katas = Major `2`).
  - **Minor:** Represents each full hundred of Katas *after* the Major component is accounted for. (e.g., for 120 Katas, after 0 Major, there's 1 hundred, so Minor `1`).
  - **Patch:** Represents the remaining Katas (tens and units) *after* Major and Minor components are accounted for. (e.g., for 120 Katas, after 0 Major and 1 Minor, there are 20 remaining, so Patch `20`).
  - Example 1: 120 Katas = `0.1.20` (0 thousands, 1 hundred, 20 remaining)
  - Example 2: 52 Katas = `0.0.52` (0 thousands, 0 hundreds, 52 remaining)
  - Example 3: 1005 Katas = `1.0.5` (1 thousand, 0 hundreds, 5 remaining)

  - **Release Strategy:** Draft releases can be created at any point. However, the **publication of Version `1.0.0` is a significant milestone, to be released to the GitHub Marketplace when 1000 Katas have been implemented (corresponding to reaching Kata ZC2000, assuming Kata numbering started at ZC1000).**