# A RESTful API for Indexed Tezos Data

[![CircleCI](https://circleci.com/gh/ecadlabs/tezos-indexer-api.svg?style=svg)](https://circleci.com/gh/ecadlabs/tezos-indexer-api)

_WARNING: This project is in early stage, active development. While we welcome users and feedback, please be warned that this project is a work in progress and users should proceed with caution._

## What is tezos-indexer-api?

`tezos-indexer-api` is a RESTful API that serves Tezos Blockchain data that has been indexed. This API is useful for efficient access to data such as balance history, or operation history over time.

Initially, this API serves data from Nomdaic Lab's [Tezos Indexer's][0] postgresql database. Other indexer backends may be added in the future, offering users a single API surface for indexed data with different backend options.

## Getting started

Docker images and pre-built binaries are available from the [releases](https://github.com/ecadlabs/tezos-indexer-api/releases) github page.

## Reporting Issues

### Security Issues

To report a security issue, please contact security@ecadlabs.com or via [keybase/jevonearth][1] on keybase.io.

Reports may be encrypted using keys published on keybase.io using [keybase/jevonearth][1].

### Other Issues & Feature Requests

Please use the [GitHub issue tracker](https://github.com/ecadlabs/tezos-indexer-api/issues) to report bugs or request features.

## Contributions

To contribute, please check the issue tracker to see if an existing issue exists for your planned contribution. If there's no Issue, please create one first, and then submit a pull request with your contribution.

For a contribution to be merged, it must be well documented, come with unit tests, and integration tests where appropriate. Submitting a "work in progress" pull request is welcome!

---

## Alternative Tezos indexers

At least two other indexers are available for Tezos.

We encourage bakers to, at a minimum, review these projects. We are eager to collaborate and be peers with these great projects.

* [Conseil](https://cryptonomic.github.io/Conseil/#/)
* [baking-soda-tezos](https://gitlab.com/9chapters/baking-soda-tezos)

## Disclaimer

THIS SOFTWARE IS PROVIDED "AS IS" AND ANY EXPRESSED OR IMPLIED WARRANTIES, INCLUDING, BUT NOT LIMITED TO, THE IMPLIED WARRANTIES OF MERCHANTABILITY AND FITNESS FOR A PARTICULAR PURPOSE ARE DISCLAIMED. IN NO EVENT SHALL THE REGENTS OR CONTRIBUTORS BE LIABLE FOR ANY DIRECT, INDIRECT, INCIDENTAL, SPECIAL, EXEMPLARY, OR CONSEQUENTIAL DAMAGES (INCLUDING, BUT NOT LIMITED TO, PROCUREMENT OF SUBSTITUTE GOODS OR SERVICES; LOSS OF USE, DATA, OR PROFITS; OR BUSINESS INTERRUPTION) HOWEVER CAUSED AND ON ANY THEORY OF LIABILITY, WHETHER IN CONTRACT, STRICT LIABILITY, OR TORT (INCLUDING NEGLIGENCE OR OTHERWISE) ARISING IN ANY WAY OUT OF THE USE OF THIS SOFTWARE, EVEN IF ADVISED OF THE POSSIBILITY OF SUCH DAMAGE.

[0]: https://gitlab.com/nomadic-labs/tezos-indexer 
[1]: https://keybase.io/jevonearth
