# Release Notes

## v1.0.0

- Update all non-major dependencies (#91)
- Merge pull request 'Configure Renovate' (#90) from renovate/configure into main
- Add renovate.json (a256793)
- Finish migration to Codeberg (#89)
- Bump go.uber.org/zap from 1.27.1 to 1.28.0 (#86)
- Bump github.com/onsi/ginkgo/v2 from 2.28.1 to 2.28.3 (#84)

## v0.7.1

- Bump github.com/onsi/gomega from 1.39.0 to 1.39.1 (#82)
- Bump github.com/onsi/ginkgo/v2 from 2.27.5 to 2.28.1 (#81)
- Bump github.com/onsi/ginkgo/v2 from 2.27.3 to 2.27.5 (#80)
- Bump github.com/onsi/gomega from 1.38.3 to 1.39.0 (#79)
- Bump github.com/onsi/gomega from 1.38.2 to 1.38.3 (#77)
- Bump github.com/onsi/ginkgo/v2 from 2.27.2 to 2.27.3 (#76)

## v0.7.0

- Add Mac ARM64 support (#75)
- Bump go.uber.org/zap from 1.27.0 to 1.27.1 (#74)
- Bump github.com/onsi/ginkgo/v2 from 2.26.0 to 2.27.2 (#73)
- Bump github.com/onsi/ginkgo/v2 from 2.25.3 to 2.26.0 (#72)
- Bump github.com/onsi/ginkgo/v2 from 2.25.2 to 2.25.3 (#71)
- Bump github.com/onsi/ginkgo/v2 from 2.24.0 to 2.25.1 (#69)
- Bump github.com/onsi/gomega from 1.38.0 to 1.38.1 (#70)
- Bump github.com/onsi/gomega from 1.37.0 to 1.38.0 (#67)
- Bump github.com/onsi/ginkgo/v2 from 2.23.4 to 2.24.0 (#68)
- Bump golang.org/x/oauth2 from 0.26.0 to 0.27.0 in the go_modules group (#66)
- Bump github.com/onsi/gomega from 1.36.2 to 1.37.0 (#63)
- Bump golang.org/x/net from 0.37.0 to 0.38.0 in the go_modules group (#64)
- Bump github.com/onsi/ginkgo/v2 from 2.23.3 to 2.23.4 (#62)
- Bump github.com/onsi/ginkgo/v2 from 2.23.0 to 2.23.3 (#61)
- Bump golang.org/x/net from 0.35.0 to 0.36.0 in the go_modules group (#59)
- Bump github.com/onsi/ginkgo/v2 from 2.22.2 to 2.23.0 (#58)
- Bump golang.org/x/sync from 0.10.0 to 0.11.0 (#57)
- Bump golang.org/x/oauth2 from 0.25.0 to 0.26.0 (#56)
- Add privacy policy (#55)

## v0.6.0

- Enhance OAuth doc, fix flag names (#54)
- Integrate oauth-cli (#53)
- Bump github.com/int128/oauth2cli from 1.14.1 to 1.15.1 (#52)
- Bump golang.org/x/oauth2 from 0.24.0 to 0.25.0 (#51)
- Bump github.com/onsi/ginkgo/v2 from 2.22.1 to 2.22.2 (#49)
- Bump github.com/onsi/ginkgo/v2 from 2.22.0 to 2.22.1 (#48)

## v0.5.0

- Update to use OAuth (#47)
- Update to golang/org/x/net v0.33.0 (#45)
- Bump github.com/onsi/gomega from 1.36.0 to 1.36.1 (#44)
- Bump github.com/onsi/gomega from 1.34.2 to 1.36.0 (#43)
- Bump github.com/onsi/ginkgo/v2 from 2.20.1 to 2.22.0 (#42)
- Bump github.com/onsi/ginkgo/v2 from 2.20.0 to 2.20.1 (#39)
- Bump github.com/onsi/ginkgo/v2 from 2.19.0 to 2.20.0 (#38)
- Bump github.com/onsi/gomega from 1.33.1 to 1.34.1 (#37)
- Bump github.com/onsi/ginkgo/v2 from 2.17.3 to 2.19.0 (#35)
- Bump github.com/onsi/ginkgo/v2 from 2.17.2 to 2.17.3 (#34)
- Bump github.com/onsi/gomega from 1.33.0 to 1.33.1 (#33)
- Bump github.com/onsi/ginkgo/v2 from 2.17.1 to 2.17.2 (#32)
- Bump github.com/onsi/gomega from 1.32.0 to 1.33.0 (#31)
- Bump golang.org/x/net from 0.20.0 to 0.23.0 in the go_modules group (#30)
- Bump github.com/onsi/ginkgo/v2 from 2.17.0 to 2.17.1 (#29)
- Bump github.com/onsi/ginkgo/v2 from 2.16.0 to 2.17.0 (#28)
- Bump github.com/onsi/gomega from 1.31.1 to 1.32.0 (#27)
- Bump github.com/onsi/ginkgo/v2 from 2.15.0 to 2.16.0 (#26)
- Bump go.uber.org/zap from 1.26.0 to 1.27.0 (#25)
- Bump github.com/onsi/gomega from 1.19.0 to 1.31.1 (#21)
- Bump go.uber.org/zap from 1.21.0 to 1.26.0 (#22)
- Bump github.com/onsi/ginkgo/v2 from 2.1.4 to 2.15.0 (#23)
- Fix typo in README.md (#24)
- Create dependabot.yml (#20)
- Bump the go_modules group across 1 directories with 2 updates (#18)
- Create go.yml workflow (#17)
- Un-ignore mocks, as this breaks builds (#19)

## v0.4.1

- Set budget to initial budget + delta, not desired end balance of budget (#15)

## v0.4.0

- Fix budget estimation logic (#13)

## v0.3.0

- Fixes #9: remove errorneous append to nonZeroChanges (#12)
- #10: implement read-only dry-run client (#11)

## v0.2.0

- Print out progress of application and add dry run mode (#7)
- Fix determination of 'ready to assign' (#6)
- Add files missed in resolving #3 (#5)
- Fixes #3 (#4)

## v0.1.1

- Fix module declaration to allow go install to work (be2cae3)
- Add ability to release artifacts (dd54405)

## v0.1.0

- Add release goal, missed some files (70dbda0)
- Remove resolved TODO (eb1ce67)
- Add support for average of monthly spent (f053d39)
- Add initial caching (d57105b)
- Add comparisons between input file and budget to prompt the user (cf3b17e)
- Add doc, make the main workflow a bit more reasonable (7fe3e1f)
- Got application working; now for README doc (5c4f670)
- Get budget application working (be66784)
- Get budget application working (4516249)
- Add PATCH support to HTTP client (d96d644)
- Add application confirmation (95bf5c5)
- Calculate deltas (ef4dd9b)
- Get delta calculation working (891be85)
- Add ability to read in files (bfec551)
- Added prompt to select budget (5e0cb78)
- Get budgets (1ed2c3b)
- Initial commit (3067d7b)
