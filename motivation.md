## Encrypted Content Management

### Motivation
We start this project after many years of experience with 1Password solution.
Originally, there were few missing pieces with 1Password such as:
- command line interace (provided in version 8)
- web server (there is no publicly available server so far)
- support of multiple cloud providers and usage of on/off-site premises

But more importantly, the 1Password has changed their license based approach to
subscription model starting version 8, see their
[reply](https://1password.community/discussion/133705/does-1password8-support-non-subscription-mode/p1?new=1).
Even though it is profiable for tha
company we considered that over the time it is not sustaibable solution, e.g.
the current pricing model of $5/month leads to $60/year and since such manager
is mandatory in our digital life it can lead to substantial expenses over the
long run. Despite monthly fee the 1Password does not provide ability to use own
infrastructure, there is lack of support for different cloud storage providers
and 1password is closed source. All of these factors lead to seek alternative
solutions and idea of implementing password manager without aforementioned
limitation. The Go-language seems to be an excellent choice for such
implementation because:
- it is open-source
- it has solid crypto library as part of its Standard Library
- there are numerous native GUI solutions, e.g.
(via [fyne.io](https://fyne.io/))
- Go-based code can be ported moble platform,
e.g.  [gomobile](https://pkg.go.dev/golang.org/x/mobile/cmd/gomobile).
- it supportis various architectures, such as AMD, ARM, Power8, and
[WebAssembly](https://www.wikiwand.com/en/WebAssembly) which allows
to implement desired functionality in single language and ported it
across multiple hardware platforms.

Therefore, after few iterations the ECM toolkit was board and it is released
here as open-source solution.

